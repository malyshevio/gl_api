package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	// import driver

	_ "github.com/lib/pq"
	"gl_api.malyshev.io/internal/data"
	"gl_api.malyshev.io/internal/jsonlog"
	"gl_api.malyshev.io/internal/mailer"
)

var (
	version   string
	buildTime string // -X linker variable
)

// config hold all configuration settings
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

// application hold the dependencies for HTTP handlers, helpers, middleware
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment type (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN") // переменная теперь на уровне окружения

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Максимальное количество запросов в секунду")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Максимальное количество запросов одновнеменно")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Включить ограничение запросов")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "localhost", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "test@example.com", "SMTP sender")

	flag.Func("cors-trusted-origins", "Довереные CORS (через пробел)", func(s string) error {
		cfg.cors.trustedOrigins = strings.Fields(s)
		return nil
	})

	displayVersion := flag.Bool("version", false, "Отобразить текущую версию и выйти")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	//init new logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// connect to DB
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection successfully established", nil)

	// для роута отладки /debug/vars
	// версия
	expvar.NewString("version").Set(version)
	// количество горутин
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	// пул подключений к БД
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	// текущая метка времени unixtimestamp
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	// create App instance
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

// openDB() возвращаем подключеие к бд
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
