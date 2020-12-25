package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/BTBurke/vatinator/bundled"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	key := os.Getenv("FILLER_API_KEY")
	if key == "" {
		log.Fatal("No API key provided.  Set FILLER_API_KEY.")
	}

	http.HandleFunc("/", filler(key))
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func filler(wantKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if key != wantKey {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		start := time.Now()
		// TODO: add api key
		tmpdir, err := ioutil.TempDir("", "fill")
		if err != nil {
			handleError(w, err)
			return
		}

		fdfPath := filepath.Join(tmpdir, "data.fdf")
		f, err := os.Create(fdfPath)
		if err != nil {
			handleError(w, err)
			return
		}

		defer r.Body.Close()
		if _, err := io.Copy(f, r.Body); err != nil {
			handleError(w, err)
			return
		}

		t, err := bundled.Asset("assets/excise.pdf")
		if err != nil {
			handleError(w, err)
			return
		}
		templatePath := filepath.Join(tmpdir, "template.pdf")
		if err := ioutil.WriteFile(templatePath, t, 0644); err != nil {
			handleError(w, err)
			return
		}

		bin, err := exec.LookPath("pdftk")
		if err != nil {
			handleError(w, err)
			return
		}

		outPath := filepath.Join(tmpdir, "out.pdf")
		cmd := exec.Command(bin, templatePath, "fill_form", fdfPath, "output", outPath)
		stdouterr, err := cmd.CombinedOutput()
		if err != nil {
			handleError(w, fmt.Errorf("pdftk error: %s. Output: %s", err.Error(), stdouterr))
			return
		}

		out, err := os.Open(outPath)
		if err != nil {
			handleError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, out)
		log.Printf("form filled in %s", time.Since(start))
	}
}

func handleError(w http.ResponseWriter, e error) {
	log.Printf("server error: %s", e.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(e.Error()))
}
