package main


import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
    initDB()

    r := gin.Default()
    registerRoutes(r)
  r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
    r.Run(":8080")
}
