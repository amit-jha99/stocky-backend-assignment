package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func rewardHandler(c *gin.Context) {
	var req RewardRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.Quantity <= 0 {
		log.WithError(err).Warn("invalid reward request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := processReward(req); err != nil {
		log.WithFields(log.Fields{
			"user_id": req.UserID,
			"stock":   req.StockSymbol,
		}).WithError(err).Error("reward processing failed")

		// ⚠️ Do NOT expose internal error details to client
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	log.WithFields(log.Fields{
		"user_id": req.UserID,
		"stock":   req.StockSymbol,
	}).Info("reward processed successfully")

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func todayStocksHandler(c *gin.Context) {
	userId := c.Param("userId")

	rows, err := DB.Query(`
		SELECT stock_symbol, SUM(quantity)
		FROM reward_events
		WHERE user_id=$1 AND DATE(rewarded_at)=CURRENT_DATE
		GROUP BY stock_symbol
	`, userId)
	if err != nil {
		log.WithError(err).Error("failed to fetch today's stocks")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	result := map[string]float64{}
	for rows.Next() {
		var s string
		var q float64
		rows.Scan(&s, &q)
		result[s] = q
	}

	c.JSON(http.StatusOK, gin.H{
		"date":    time.Now().Format("2006-01-02"),
		"rewards": result,
	})
}

func historicalINRHandler(c *gin.Context) {
	userId := c.Param("userId")

	rows, err := DB.Query(`
		SELECT DATE(rewarded_at), stock_symbol, SUM(quantity)
		FROM reward_events
		WHERE user_id=$1 AND DATE(rewarded_at) < CURRENT_DATE
		GROUP BY DATE(rewarded_at), stock_symbol
	`, userId)
	if err != nil {
		log.WithError(err).Error("failed to fetch historical rewards")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	result := map[string]float64{}
	for rows.Next() {
		var d time.Time
		var s string
		var q float64
		rows.Scan(&d, &s, &q)
		result[d.Format("2006-01-02")] += getLatestPrice(s) * q
	}

	c.JSON(http.StatusOK, gin.H{"history": result})
}

func statsHandler(c *gin.Context) {
	userId := c.Param("userId")

	rows, err := DB.Query(`
		SELECT stock_symbol, SUM(quantity)
		FROM reward_events
		WHERE user_id=$1
		GROUP BY stock_symbol
	`, userId)
	if err != nil {
		log.WithError(err).Error("failed to fetch stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	total := 0.0
	portfolio := map[string]float64{}

	for rows.Next() {
		var s string
		var q float64
		rows.Scan(&s, &q)
		value := getLatestPrice(s) * q
		portfolio[s] = value
		total += value
	}

	c.JSON(http.StatusOK, gin.H{
		"portfolio":       portfolio,
		"total_value_inr": total,
	})
}

func portfolioHandler(c *gin.Context) {
	userId := c.Param("userId")

	rows, err := DB.Query(`
		SELECT stock_symbol, SUM(quantity)
		FROM reward_events
		WHERE user_id=$1
		GROUP BY stock_symbol
	`, userId)
	if err != nil {
		log.WithError(err).Error("failed to fetch portfolio")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer rows.Close()

	holdings := []gin.H{}
	total := 0.0

	for rows.Next() {
		var s string
		var q float64
		rows.Scan(&s, &q)
		value := getLatestPrice(s) * q
		total += value

		holdings = append(holdings, gin.H{
			"stock":    s,
			"quantity": q,
			"value_inr": value,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"holdings":        holdings,
		"total_value_inr": total,
	})
}
