package main

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
	"log"
)

func rewardHandler(c *gin.Context) {
    var req RewardRequest

    if err := c.ShouldBindJSON(&req); err != nil || req.Quantity <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    err := processReward(req)
    if err != nil {
        log.Println("REWARD ERROR:", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success"})
}


func todayStocksHandler(c *gin.Context) {
    userId := c.Param("userId")

    rows, _ := DB.Query(`
        SELECT stock_symbol, SUM(quantity)
        FROM reward_events
        WHERE user_id=$1 AND DATE(rewarded_at)=CURRENT_DATE
        GROUP BY stock_symbol
    `, userId)

    result := map[string]float64{}
    for rows.Next() {
        var s string
        var q float64
        rows.Scan(&s, &q)
        result[s] = q
    }

    c.JSON(200, gin.H{"date": time.Now().Format("2006-01-02"), "rewards": result})
}

func historicalINRHandler(c *gin.Context) {
    userId := c.Param("userId")

    rows, _ := DB.Query(`
        SELECT DATE(rewarded_at), stock_symbol, SUM(quantity)
        FROM reward_events
        WHERE user_id=$1 AND DATE(rewarded_at) < CURRENT_DATE
        GROUP BY DATE(rewarded_at), stock_symbol
    `, userId)

    result := map[string]float64{}
    for rows.Next() {
        var d time.Time
        var s string
        var q float64
        rows.Scan(&d, &s, &q)
        result[d.Format("2006-01-02")] += getLatestPrice(s) * q
    }

    c.JSON(200, gin.H{"history": result})
}

func statsHandler(c *gin.Context) {
    userId := c.Param("userId")

    rows, _ := DB.Query(`
        SELECT stock_symbol, SUM(quantity)
        FROM reward_events
        WHERE user_id=$1
        GROUP BY stock_symbol
    `, userId)

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

    c.JSON(200, gin.H{
        "portfolio": portfolio,
        "total_value_inr": total,
    })
}

func portfolioHandler(c *gin.Context) {
    userId := c.Param("userId")

    rows, _ := DB.Query(`
        SELECT stock_symbol, SUM(quantity)
        FROM reward_events
        WHERE user_id=$1
        GROUP BY stock_symbol
    `, userId)

    holdings := []gin.H{}
    total := 0.0

    for rows.Next() {
        var s string
        var q float64
        rows.Scan(&s, &q)
        value := getLatestPrice(s) * q
        total += value

        holdings = append(holdings, gin.H{
            "stock": s,
            "quantity": q,
            "value_inr": value,
        })
    }

    c.JSON(200, gin.H{
        "holdings": holdings,
        "total_value_inr": total,
    })
}
