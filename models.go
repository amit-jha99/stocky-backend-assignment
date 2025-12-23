package main

type RewardRequest struct {
    UserID         string  `json:"user_id"`
    StockSymbol    string  `json:"stock_symbol"`
    Quantity       float64 `json:"quantity"`
    RewardedAt     string  `json:"rewarded_at"`
    IdempotencyKey string  `json:"idempotency_key"`
}
