package main

import "github.com/gin-gonic/gin"

func registerRoutes(r *gin.Engine) {
    r.POST("/reward", rewardHandler)
    r.GET("/today-stocks/:userId", todayStocksHandler)
    r.GET("/historical-inr/:userId", historicalINRHandler)
    r.GET("/stats/:userId", statsHandler)
    r.GET("/portfolio/:userId", portfolioHandler) // bonus
}
