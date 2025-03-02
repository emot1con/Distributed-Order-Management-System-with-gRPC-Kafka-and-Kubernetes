package main

import (
	routes "broker/cmd/api"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	srv := http.Server{
		Addr:    ":8080",
		Handler: routes.Routes(),
	}
	logrus.Info("Start broker on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
