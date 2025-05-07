package main

import (
	"log"

	"github.com/ra-khalish/gosocial/internal/db"
	"github.com/ra-khalish/gosocial/internal/env"
	"github.com/ra-khalish/gosocial/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConn: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConn: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "Development"),
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConn,
		cfg.db.maxIdleConn,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Printf("database connection pool established\n")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	// os.LookupEnv("PATH")

	mux := app.mount()
	log.Fatal(app.run(mux))
}
