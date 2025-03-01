package main

import (
	"net/http"
	"user_service/cmd/api"
	"user_service/cmd/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logrus.Info("connect DB")
	gormDB, err := db.Connect(&gin.Context{})
	if err != nil {
		logrus.Fatalf("Error connect DB: %v", err)
	}

	srv := http.Server{
		Addr:    ":5000",
		Handler: api.Router(gormDB),
	}
	logrus.Info("Starting server on port 5000")

	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatalf("Error listen and serve: %v", err)
	}
}
