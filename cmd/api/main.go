package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// import driver
	_ "github.com/lib/pq"
)

const version = "1.0.0" // just in case I don't generate it et

// config hold all configuration settings
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

// application hold the dependencies for HTTP handlers, helpers, middleware
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment type (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://gl_api:pa55word@localhost/gl_api?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	//init new logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// connect to DB
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection successfully established")

	// create App instance
	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// openDB() возвращаем подключеие к бд
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
