
# CryptoTrack

**CryptoTrack** is a sample Go project that demonstrates the architecture of a modern HTTP service. It features clean code, comprehensive unit and integration tests, Swagger documentation, and seamless interaction with both external and internal systems.

---

## What does this project do?

- Fetches real-time cryptocurrency prices from the CoinGecko API
- Saves price data in PostgreSQL
- Provides an HTTP API to view current and historical price data

---

## Key Features

- **HTTP API:** Endpoints to track, store, and retrieve cryptocurrency prices.
- **Swagger UI:** Auto-generated API documentation available out-of-the-box at `/swagger/index.html`.
- **External API integration:** Fetches live prices from CoinGecko.
- **Test coverage:** Includes both unit tests (for business logic) and integration tests (with testcontainers and mocks).
- **Project structure:** Clean separation of concerns (`config`, `db`, `service`, `client`, `handlers`, `validation`).
- **Database migrations:** Managed in code via golang-migrate.
- **Easy launch with Docker Compose.**

---

## Quick Start

1. **Clone the repository:**
   ```sh
   git clone https://github.com/kitbuilder587/cryptotrack.git
   cd cryptotrack
   ```

2. **Build and run the service:**
   ```sh
   docker-compose up --build
   ```

3. **Swagger UI:**  
   Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) in your browser.

---

## Main Endpoints

- `POST /track?coin=bitcoin` — Fetches and saves the current price for the specified coin.
- `GET /latest?coin=bitcoin` — Retrieves the latest saved price for the specified coin.
- `GET /history?coin=bitcoin&limit=10` — Retrieves the N most recent price entries for the specified coin.
- `GET /health` — Service health check.
- `GET /swagger/index.html` — Interactive OpenAPI (Swagger) documentation.

---

## How Swagger Works

- API documentation for all endpoints is auto-generated from Go comments using [swaggo](https://github.com/swaggo/swag).
- The interactive Swagger UI lets you view endpoint descriptions and try out API requests right from your browser.

---

## Testing

- Includes both unit and integration tests using the standard Go `testing` package and [testcontainers](https://github.com/testcontainers/testcontainers-go).
- Both data storage and CoinGecko integration logic are covered by tests.

---

## Who is this for?

- Developers who want a quick way to explore modern Go HTTP service structure, see how to organize code, add integration/unit tests, connect Swagger, and manage database migrations.

---

**This project is a great starting point for your pet project or as a base for a production-ready service!**
