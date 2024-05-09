package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dbHost, dbPort, dbName, dbUser, dbPassword, migrationsPath, migrationsTable string

	flag.StringVar(&dbHost, "db-host", "localhost", "Database host")
	flag.StringVar(&dbPort, "db-port", "5432", "Database port")
	flag.StringVar(&dbName, "db-name", "auth", "Database name")
	flag.StringVar(&dbUser, "db-user", "postgres", "Database user")
	flag.StringVar(&dbPassword, "db-password", "mysecret", "Database password")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if dbName == "" {
		panic("db-name is required")
	}
	if dbUser == "" {
		panic("db-user is required")
	}
	if dbPassword == "" {
		panic("db-password is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	// Формируем строку подключения к базе данных PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, migrationsTable)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		panic(err)
	}

	// Выполняем миграции до последней версии
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}
}
