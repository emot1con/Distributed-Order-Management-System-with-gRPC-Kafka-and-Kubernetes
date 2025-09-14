package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

func Connect() (*sql.DB, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Jakarta", dbUser, dbPass, dbName)

	for i := 1; i < 5; i++ {
		DB, err := sql.Open("pgx", dsn)
		if err != nil {
			logrus.Warnf("Database connection failed. Retrying in 5 seconds... (%d/5)", i+1)
			time.Sleep(5 * time.Second)
			logrus.Error(err)
			continue
		}

		if err := DB.Ping(); err != nil {
			DB.Close()
			logrus.Warnf("Ping database. Retrying in 5 seconds... (%d/5)", i+1)
			time.Sleep(5 * time.Second)
			logrus.Error(err)
			continue
		}

		DB.SetMaxOpenConns(20)
		DB.SetMaxIdleConns(5)
		DB.SetConnMaxLifetime(1 * time.Hour)
		DB.SetConnMaxIdleTime(15 * time.Minute)

		logrus.Info("Database connected successfully")

		return DB, nil
	}

	return nil, fmt.Errorf("failed connect to database")
}
