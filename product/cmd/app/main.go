package main

import (
	"net/http"
	"product_service/cmd/db"

	"github.com/sirupsen/logrus"
)

func main() {
	if _, err := db.Connect(); err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}
	srv := http.Server{
		Addr:    ":40001",
		Handler: nil,
	}

	logrus.Info("Starting server on port 40001")
	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
