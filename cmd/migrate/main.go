package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Tutuacs/pkg/db"
	"github.com/Tutuacs/pkg/logs"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	dbConnection, err := db.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(dbConnection, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	} else if cmd == "force" {
		if len(os.Args) < 3 {
			log.Fatal("Please provide a version number to force")
		}
		version, err := strconv.Atoi(os.Args[len(os.Args)-2])
		if err != nil {
			log.Fatal("Invalid version number")
		}
		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}
	}
	logs.OkLog("Migration completed")
}
