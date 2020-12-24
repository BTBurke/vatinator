package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BTBurke/vatinator/bundled"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", filler)
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func filler(w http.ResponseWriter, r *http.Request) {
	// TODO: add api key
	tmpdir, err := ioutil.TempDir("", "fill")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fdfPath := filepath.Join(tmpdir, "data.fdf")
	f, err := os.Create(fdfPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()
	if _, err := io.Copy(f, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	t, err := bundled.Asset("assets/excise.pdf")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	templatePath := filepath.Join(tmpdir, "template.pdf")
	if err := ioutil.WriteFile(templatePath, t, 0644); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	bin, err := exec.LookPath("pdftk")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	outPath := filepath.Join(tmpdir, "out.pdf")
	cmd := exec.Command(bin, templatePath, "fill_form", fdfPath, "output", outPath)
	stdouterr, err := cmd.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(stdouterr)
		return
	}

	out, err := os.Open(outPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	io.Copy(w, out)
}
