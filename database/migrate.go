package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found,using os env")
	}
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	dsn := os.Getenv("DB_DSN")
	if migrationsPath == "" || dsn == "" {
		log.Fatal("missing MIGRATION_PATH or DB_DSN")
	}
	if err := RunMigrations(migrationsPath, dsn); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("migrations applied succesfully!")
}

func RunMigrations(path, dsn string) error {
	m, err := migrate.New(path, dsn)
	if err != nil {
		return fmt.Errorf("Create migrate instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("migration source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("migration source close error: %v", dbErr)
		}
	}()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	return nil
}
