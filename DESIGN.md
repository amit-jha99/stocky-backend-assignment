# Design Notes and Edge Case Handling

This document describes the key design decisions, edge case handling,
and scalability considerations for the Stocky backend service.

---

## Idempotency and Duplicate Requests

Each reward request includes an idempotency_key which is enforced as unique
at the database level. If a request is replayed with the same key, the system
rejects it with a 409 Conflict response. This ensures that rewards are processed
exactly once and protects against duplicate submissions or replay attacks.

---

## Stock Price Handling

Reward records store stock quantities rather than fixed INR values. Portfolio
valuation is calculated dynamically using the latest available stock price.
Prices are cached in the database to reduce dependency on external price sources
during API execution.

---

## Stock Splits, Mergers, and Delisting

Because portfolio valuation is computed dynamically, changes in stock prices
such as splits, mergers, or delisting events are automatically reflected in
portfolio value without modifying historical reward records.

---

## Monetary Precision

All monetary values are stored using fixed-precision numeric types in PostgreSQL.
This avoids floating-point rounding errors and ensures consistent INR valuation
across calculations.

---

## Adjustments and Refunds

The system uses a ledger-based design that supports compensating entries.
Adjustments or refunds can be implemented as debit or credit ledger entries
without mutating historical reward data, preserving auditability.

---

## Price Downtime and Stale Data

Stock prices are cached in the database. If fresh pricing data is unavailable,
the system can continue operating using cached values, ensuring API availability
and graceful degradation.

---

## Scaling Considerations

The service is stateless and can be horizontally scaled. Database constraints
ensure correctness under concurrent requests. Read-heavy endpoints such as
portfolio and stats can be optimized using caching if required. The schema
supports future extensions such as scheduled price refresh jobs or integration
with external price providers.
