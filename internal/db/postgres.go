package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

func GetDBInstance() *sql.DB {
	once.Do(func() {
		dbInstance = InitPostgres()
	})
	return dbInstance
}

func InitPostgres() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err := migrateDB(db); err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	return db
}

func migrateDB(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("не удалось создать экземпляр базы данных: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations",
		"DeFiExplorer",
		driver,
	)
	if err != nil {
		return fmt.Errorf("не удалось создать экземпляр мигратора: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка при выполнении миграций: %w", err)
	}

	return nil
}
