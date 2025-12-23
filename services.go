package main

import (
    "math/rand"
    "time"

    "github.com/google/uuid"
)

// Fetch price from cache or generate mock and store
func getLatestPrice(stock string) float64 {
    var price float64

    err := DB.QueryRow(`
        SELECT price_inr
        FROM stock_prices
        WHERE stock_symbol = $1
    `, stock).Scan(&price)

    if err == nil {
        return price
    }

    // fallback: mocked price
    rand.Seed(time.Now().UnixNano())
    price = 500 + rand.Float64()*1500

    DB.Exec(`
        INSERT INTO stock_prices (stock_symbol, price_inr, fetched_at)
        VALUES ($1, $2, now())
        ON CONFLICT (stock_symbol)
        DO UPDATE SET price_inr = $2, fetched_at = now()
    `, stock, price)

    return price
}

func processReward(req RewardRequest) error {
    eventID := uuid.New()
    rewardedAt, _ := time.Parse(time.RFC3339, req.RewardedAt)

    _, err := DB.Exec(`
        INSERT INTO reward_events
        (id, user_id, stock_symbol, quantity, rewarded_at, idempotency_key)
        VALUES ($1, $2, $3, $4, $5, $6)
    `,
        eventID,
        req.UserID,
        req.StockSymbol,
        req.Quantity,
        rewardedAt,
        req.IdempotencyKey,
    )
    if err != nil {
        return err
    }

    price := getLatestPrice(req.StockSymbol)
    totalValue := price * req.Quantity
    fee := totalValue * 0.01

    entries := []struct {
        entryType string
        amount    float64
        direction string
    }{
        {"STOCK", req.Quantity, "CREDIT"},
        {"CASH", totalValue, "DEBIT"},
        {"FEE", fee, "DEBIT"},
    }

    for _, e := range entries {
        DB.Exec(`
            INSERT INTO ledger_entries
            (id, event_id, entry_type, symbol, amount, direction)
            VALUES ($1, $2, $3, $4, $5, $6)
        `,
            uuid.New(),
            eventID,
            e.entryType,
            req.StockSymbol,
            e.amount,
            e.direction,
        )
    }

    return nil
}
