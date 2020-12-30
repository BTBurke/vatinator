package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/BTBurke/vatinator"
	"github.com/BTBurke/vatinator/handlers"
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

type ServerConfig struct {
	DBPath     string
	DBCreateOK bool
}

func init() {
	// flag for creating the directories if they don't exist, otherwise it will error on start
	pflag.Bool("create", false, "whether it is ok to create directories and db on startup")
	pflag.Bool("dev", false, "run in development mode, using temp directories and deleting on shutdown")
}

func main() {
	// TODO: add viper config
	viper.SetDefault("Port", "8080")
	if version == "dev" {
		viper.SetDefault("DataDir", "/tmp/vat/data")
		viper.SetDefault("UploadDir", "/tmp/vat/upload")
		viper.SetDefault("ExportDir", "/tmp/vat/export")
		viper.SetDefault("CredentialFile", "./vatinator-f91ccb107c2c.json")
	} else {
		viper.SetDefault("DataDir", "/var/vat/data")
		viper.SetDefault("UploadDir", "/var/vat/upload")
		viper.SetDefault("ExportDir", "/var/vat/export")
		viper.SetDefault("CredentialFile", "/etc/vat/vatinator-f91ccb107c2c.json")
	}
	viper.SetEnvPrefix("vat")

	// preferences for config paths
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/vat")
	viper.AddConfigPath("/etc/vat")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("No config file found.  Running with defaults.")
		} else {
			log.Fatal(errors.Wrap(err, "configuration error"))
		}
	}

	// parse flags then check config
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
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

	log.Printf("Data directory: %s", viper.Get("DataDir"))
	log.Printf("Upload directory: %s", viper.Get("UploadDir"))
	log.Printf("Export directory: %s", viper.Get("ExportDir"))

	// set up account service
	db, err := vatinator.NewDB(filepath.Join(viper.GetString("DataDir"), "vat.db"))
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
	sessionSvc, err := vatinator.NewSessionService(filepath.Join(viper.GetString("DataDir"), "session.db"), keys...)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	port := viper.GetString("Port")
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Use(handlers.SessionMiddleware(sessionSvc))
		r.Post("/file", handlers.FileAddHandler(viper.GetString("UploadDir")))
		r.Get("/account", handlers.GetAccountHandler(accountSvc))
		r.Post("/account", handlers.UpdateAccountHandler(accountSvc))
		r.Post("/process", handlers.ProcessHandler())
	})
	// no auth routes
	r.Post("/create", handlers.CreateAccountHandler(accountSvc, sessionSvc))
	r.Post("/login", handlers.LoginHandler(accountSvc, sessionSvc))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-c
		log.Printf("Received signal: %+v", sig)
		cancel()
	}()

	if err := serve(ctx, r, port); err != nil {
		log.Fatal(err)
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
		default:
		}
	}
	return nil
}
