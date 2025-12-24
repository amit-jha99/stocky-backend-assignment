package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Logrus setup
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	godotenv.Load()
	initDB()

	log.Info("starting stocky backend server")

	r := gin.Default()
	registerRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.WithError(err).Fatal("failed to start server")
	}
}
