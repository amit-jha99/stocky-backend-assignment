# API Specifications

This document describes the HTTP APIs exposed by the Stocky backend service,
including request and response formats.

---

## POST /reward

Description:
Rewards a user with a stock. Duplicate or replayed requests are prevented
using an idempotency key.

Method: POST  
Endpoint: /reward

Request Body:
{
  "user_id": "user_1812",
  "stock_symbol": "INFY",
  "quantity": 1.25,
  "rewarded_at": "2025-12-24T10:00:00Z",
  "idempotency_key": "reward-004"
}

Responses:

200 OK
{
  "status": "success"
}

400 Bad Request
{
  "error": "invalid request"
}

409 Conflict
{
  "error": "duplicate reward"
}

---

## GET /today-stocks/{userId}

Description:
Returns stocks rewarded on the current date for the given user.

Method: GET  
Endpoint:
/today-stocks/user_1812

Response:
{
  "date": "2025-12-24",
  "rewards": {
    "INFY": 1.25
  }
}

If no rewards exist for today:
{
  "date": "2025-12-24",
  "rewards": {}
}

---

## GET /stats/{userId}

Description:
Returns aggregated portfolio value in INR.

Method: GET  
Endpoint:
/stats/user_1812

Response:
{
  "portfolio": {
    "INFY": 1875.50
  },
  "total_value_inr": 1875.50
}

---

## GET /historical-inr/{userId}

Description:
Returns historical INR valuation for rewards from previous dates.

Method: GET  
Endpoint:
/historical-inr/user_1812

Response:
{
  "history": {
    "2025-12-23": 3500.75
  }
}

If no historical data exists:
{
  "history": {}
}

---

## GET /portfolio/{userId} (Bonus)

Description:
Returns detailed stock holdings and valuation.

Method: GET  
Endpoint:
/portfolio/user_1812

Response:
{
  "holdings": [
    {
      "stock": "INFY",
      "quantity": 1.25,
      "value_inr": 1875.50
    }
  ],
  "total_value_inr": 1875.50
}
