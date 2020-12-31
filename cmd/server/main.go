package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/BTBurke/vatinator"
	"github.com/BTBurke/vatinator/handlers"
	magic "github.com/caddyserver/certmagic"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version string = "dev"
	commit  string
	date    string
)

func init() {
	// flag for creating the directories if they don't exist, otherwise it will error on start
	pflag.Bool("create", false, "whether it is ok to create directories and db on startup")
	pflag.Bool("dev", false, "run in development mode, using temp directories and deleting on shutdown")
	pflag.StringP("config", "c", "config", "specifies a config file to use")
}

func main() {
	serverStart := time.Now()
	viper.SetDefault("port", "8080")
	viper.SetDefault("data_dir", "/var/vat/data")
	viper.SetDefault("upload_dir", "/var/vat/upload")
	viper.SetDefault("export_dir", "/var/vat/export")
	viper.SetDefault("credential_file", "/etc/vat/vatinator-f91ccb107c2c.json")
	viper.SetEnvPrefix("vat")

	// parse flags then get config
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}

	// preferences for config paths
	viper.SetConfigName(viper.GetString("config"))
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/vat")
	viper.AddConfigPath("/etc/vat")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(errors.Wrap(err, "configuration error"))
		}
	}

	// in dev mode, auto create temp directories and then delete them all on shutdown
	if viper.GetBool("dev") {
		log.Printf("Running in development mode")
		if err := checkConfig(viper.AllSettings(), true); err != nil {
			log.Fatal(err)
		}
		defer func() {
			log.Printf("Shutting down dev server, deleting data")
			if err := os.RemoveAll("/tmp/vat"); err != nil {
				log.Fatal(err)
			}
		}()
	} else {
		if err := checkConfig(viper.AllSettings(), viper.GetBool("create")); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Data directory: %s", viper.Get("data_dir"))
	log.Printf("Upload directory: %s", viper.Get("upload_dir"))
	log.Printf("Export directory: %s", viper.Get("export_dir"))
	log.Printf("Using postmark for transactional emails")

	// set up account service
	db, err := vatinator.NewDB(filepath.Join(viper.GetString("data_dir"), "vat.db"))
	if err != nil {
		log.Fatal(err)
	}
	accountSvc := vatinator.NewAccountService(db)
	if err := vatinator.Migrate(db, "1.sql"); err != nil {
		log.Fatal(err)
	}

	// session service
	keys, err := vatinator.GetSessionKeys(db)
	if err != nil {
		log.Fatal(err)
	}
	sessionSvc, err := vatinator.NewSessionService(filepath.Join(viper.GetString("data_dir"), "session.db"), keys...)
	if err != nil {
		log.Fatal(err)
	}

	// token service
	tokenSvc := vatinator.NewTokenService(keys[0])

	// email service
	emailSvc := vatinator.NewEmailService(viper.GetString("postmark_server_token"), viper.GetString("postmark_api_token"))

	// process service
	processSvc := vatinator.NewProcessService(viper.GetString("upload_dir"),
		viper.GetString("export_dir"),
		viper.GetString("credential_file"),
		accountSvc,
		tokenSvc,
		emailSvc)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:3000", "https://vatinator.com", "https://www.vatinator.com", "https://vatinator.vercel.app"},
		//AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	port := viper.GetString("port")
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Use(handlers.SessionMiddleware(sessionSvc))
		r.Post("/file", handlers.FileAddHandler(viper.GetString("upload_dir")))
		r.Get("/account", handlers.GetAccountHandler(accountSvc))
		r.Post("/account", handlers.UpdateAccountHandler(accountSvc))
		r.Post("/process", handlers.ProcessHandler(processSvc))
	})
	fs := http.FileServer(http.Dir(viper.GetString("export_dir")))
	r.With(handlers.TokenMiddleware(tokenSvc)).Handle("/export/*", http.StripPrefix("/export", fs))
	// no auth routes
	r.Post("/create", handlers.CreateAccountHandler(accountSvc, sessionSvc))
	r.Post("/login", handlers.LoginHandler(accountSvc, sessionSvc))
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		resp := []byte(fmt.Sprintf("Version: %s\nCommit: %s\nDate: %s\nUptime: %s\n", version, commit, date, time.Since(serverStart)))
		w.Write(resp)
	})
	r.With(handlers.SessionMiddleware(sessionSvc)).Get("/session", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-c
		log.Printf("Received signal: %+v", sig)
		cancel()
	}()

	if port != "443" {
		// serve http on port
		if err := serve(ctx, r, port); err != nil {
			log.Fatal(err)
		}
	} else {
		// serve TLS on 443 with auto certs
		if err := serveTLS(ctx, r); err != nil {
			log.Fatal(err)
		}
	}
}

// runs server in go routine to allow safe shutdown and clean up of deferred funcs
func serve(ctx context.Context, h http.Handler, port string) (err error) {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("Server running on port %s", port)
	<-ctx.Done()
	log.Printf("Server shutting down")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed: %s", err)
	}

	if err == http.ErrServerClosed {
		return nil
	}
	return
}

func serveTLS(ctx context.Context, h http.Handler) (err error) {
	srv := &http.Server{
		Addr:    ":443",
		Handler: h,
	}

	cfg, err := magic.TLS([]string{"api.vatinator.com"})
	if err != nil {
		return err
	}

	ln, err := tls.Listen("tcp", ":443", cfg)
	if err != nil {
		return err
	}

	go func() {
		if err = srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("Server running TLS on port 443")
	<-ctx.Done()
	log.Printf("Server shutting down")
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed: %s", err)
	}

	if err == http.ErrServerClosed {
		return nil
	}
	return
}

// checks config values.  If note createOK, then the directories should already exist or it will error
// out.  The first run should have --create flag to create the dirs with the permission of the server UID/GID.
func checkConfig(cfg map[string]interface{}, createOK bool) error {
	for key, value := range cfg {
		switch {
		case createOK && strings.HasSuffix(key, "dir"):
			path := value.(string)
			if err := os.MkdirAll(path, 0700); err != nil {
				return errors.Wrapf(err, "error creating directory %s", path)
			}
		case !createOK && strings.HasSuffix(key, "dir"):
			path := value.(string)
			finfo, err := os.Stat(path)
			if os.IsNotExist(err) || !finfo.IsDir() {
				return errors.Wrapf(err, "expected %s to exist", path)
			}
		case strings.HasSuffix(key, "file"):
			path := value.(string)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return errors.Wrapf(err, "expected %s to exist", path)
			}
		case strings.HasSuffix(key, "token"):
			token := value.(string)
			if len(token) == 0 {
				return fmt.Errorf("%s must be set", key)
			}
		default:
		}
	}
	return nil
}
