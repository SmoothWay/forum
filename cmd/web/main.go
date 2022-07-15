package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SmoothWay/forum/internal"
	"github.com/SmoothWay/forum/pkg/models"
	"github.com/SmoothWay/forum/pkg/models/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "forum.db", "Database source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := models.OpenDB(*dsn, "sqlite3", "./pkg/models/sqlite/init.sql")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := internal.NewTemplateCache("./ui/templates/")
	if err != nil {
		errorLog.Fatal(err)
	}
	app := &internal.Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		TemplateCache: templateCache,
		User:          &sqlite.UserModel{DB: db},
		Posts:         &sqlite.PostModel{DB: db},
		Session:       &sqlite.Session{DB: db},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on http://localhost%s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
