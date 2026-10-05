# Pismo Technical Case

[![CI](https://github.com/MateusBezerraa/pismo-technical-case/actions/workflows/ci.yml/badge.svg)](https://github.com/MateusBezerraa/pismo-technical-case/actions/workflows/ci.yml)


Simple transactions service with 3 endpoints, built with Go and clean architecture.

---
## Quick start

The project ships with a `Makefile` that wraps every common task.
Run `make` or `make help` to see the full list of commands.

### With Docker (recommended for a clean environment)

```bash
make docker        # build the image
make docker-run    # start the container on :8080
```

> PS. If you used a previous version of this project, wipe stale data first:
> `make docker-clean`

Or with persistence to a named volume:

```bash
docker run --rm -p 8080:8080 -v pismo-data:/data pismo-technical-case
```

### Locally (requires Go 1.27+)

```bash
./run.sh
# or
make run
```

The API will be available at `http://localhost:8080`.

---

## Endpoints

### `GET /health` — Liveness probe

```bash
curl http://localhost:8080/health
```
**Response** `200 OK`:

```json
{
  "status": "ok"
}
```
---

### `POST /accounts` — Create an account

```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"document_number":"12345678900"}'
```
**Response** `201 Created`:

```json
{"account_id": 1, "document_number": "12345678900"}
```
---

### `GET /accounts/{accountId}` — Retrieve an account

```bash
curl http://localhost:8080/accounts/1
```
**Response** `200 OK`:

```json
{"account_id": 1, "document_number": "12345678900"}
```
**Errors:**

- `404 Not Found` — account does not exist
    
- `400 Bad Request` — invalid account ID
    
- `409 Conflict` — document number already exists


---

### `POST /transactions` — Create a transaction

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{"account_id":1,"operation_type_id":4,"amount":123.45}'
```
**Response** `201 Created`:

```json
{
  "transaction_id": 1,
  "account_id": 1,
  "operation_type_id": 4,
  "amount": 123.45,
  "event_date": "2026-10-02T12:00:00Z"
}
```
**Errors:**

- `404 Not Found` — account does not exist
    
- `400 Bad Request` — invalid operation type or non-positive amount
    

---

## Business rules

|ID|Operation|Stored sign|
|---|---|---|
|1|Normal Purchase|Negative|
|2|Purchase with installments|Negative|
|3|Withdrawal|Negative|
|4|Credit Voucher|Positive|

The **sign is applied by the domain** based on operation type. Clients always send a positive `amount`.

---

## Architecture

**Clean/Hexagonal Architecture** — dependencies point inward.

```text
┌──────────────────────────────────────┐
│   cmd/api (wiring, HTTP server)      │
├──────────────────────────────────────┤
│   adapter/httpx   │  adapter/repo    │
│   (inbound HTTP)  │  (outbound DB)   │
├──────────────────────────────────────┤
│   usecase (orchestration, ports)     │
├──────────────────────────────────────┤
│   domain (entities, rules, errors)   │
└──────────────────────────────────────┘
```
### Layer responsibilities

|Layer|Responsibility|Depends on|
|---|---|---|
|`domain`|Entities, business rules, domain errors|nothing|
|`usecase`|Orchestration; defines repository **ports**|`domain`|
|`adapter/httpx`|HTTP handlers, DTOs, middlewares|`domain`, `usecase`|
|`adapter/repository`|Persistence (SQLite, in-memory)|`domain`|
|`cmd/api`|Wiring, bootstrap|everything|

**Key point:** the `domain` and `usecase` packages **never import** `adapter`. Swapping SQLite for PostgreSQL means writing a new adapter — no change in business logic.

---

## Project structure

```text
.
├── cmd/api/                      # entrypoint
├── internal/
│   ├── domain/                   # entities + rules (pure)
│   ├── usecase/                  # orchestration + ports
│   └── adapter/
│       ├── httpx/                # HTTP: handlers, dto, middleware
│       └── repository/
│           ├── memory/           # in-memory adapter (alternative)
│           └── sqlite/           # SQLite adapter (default)
├── Dockerfile
├── Makefile
├── run.sh
└── README.md
```
---

## Tech stack

|Concern|Choice|Why|
|---|---|---|
|Language|Go 1.27|Performance, simplicity, great concurrency|
|HTTP router|`net/http` ServeMux (stdlib)|Go 1.22+ supports method + path params natively — no external dependency|
|Persistence|SQLite via `modernc.org/sqlite`|Pure Go, zero CGO, zero setup|
|Logging|`log/slog` (stdlib)|Structured logging, JSON-ready|
|Testing|`testing` + `httptest` (stdlib)|Standard, no external deps|

**Zero non-SQLite dependencies.** The only external package is the SQLite driver.

---

## Running tests

```bash
make test           # runs all tests with race detector
make test-cover     # runs tests + prints coverage summary
make test-html      # opens coverage report in browser
```
Coverage:

- `internal/domain` — 100%
- `internal/usecase` — 88%
- `internal/adapter/httpx/handler` — 91%

---

## Design decisions

### Why clean architecture?

The case criteria are **maintainability, simplicity, testability, documentation**. Clean architecture gives all four:

- Business logic is isolated from infrastructure.
    
- Repository ports allow swapping SQLite for Postgres with zero changes to use cases.
    
- Testing: domain is pure (no mocks), use cases use simple mocks, handlers use SQLite `:memory:`.
    

### Why SQLite?

- **Zero setup:** `go run .` just works. No `docker-compose up` needed.
    
- **Realistic:** tests use the real SQL engine, not a fake.
    
- **Extensible:** Next steps are easier when queries are declarative (`SELECT ... WHERE` vs iterating maps).

### Why only 2 middlewares (Recovery + Logging)?

**Simplicity** is an explicit criterion. `Recovery` prevents panics from killing the server. `Logging` gives observability. Adding `RequestID` or a `Chain` abstraction for 3 endpoints would be over-engineering — I opted to add them only if the system grows.

### Why `domain.NewTransaction` instead of validating in the use case?

The **sign rule** and **operation validity** are business invariants. They belong to the domain, not to the application layer. This means:

- Any entry point (HTTP, CLI, messaging) enforces the same rules.
    
- The use case has zero business logic — it only orchestrates.
    

### Why no exported `IsValid` / `IsDebit`?

They're internal details, only used by `NewTransaction`. Keeping them lowercase (`isValid`, `isDebit`) prevents external code from duplicating the rules. **Small public surface = fewer ways to misuse the package.**

### Why `:memory:` in tests?

Each test gets an isolated, fresh database in RAM. Fast, no cleanup, realistic. Requires `SetMaxOpenConns(1)` (SQLite `:memory:` creates a DB per connection).

### Why `document_number` is unique?

The spec implies a cardholder has one account, so the same document can't create two accounts. Enforced via **unique index** on SQLite (atomic) instead of check-then-insert in the use case (which has a race condition). The adapter translates the constraint violation into `domain.ErrDocumentAlreadyExists`, mapped to `409 Conflict`

---

## What I'd add next

If this were a production system:

- **Persistent storage:** PostgreSQL via `pgx` with migrations (`golang-migrate`).
    
- **Request ID middleware:** for distributed tracing (`X-Request-ID` propagation).
    
- **OpenTelemetry:** traces + metrics export.
    
- **Pagination:** for `GET /accounts/{id}/transactions`.
    
- **Authentication:** JWT validation via OIDC.
    
- **Contract tests:** between consumers and producers (Pact).
    
- **OpenAPI spec:** generated from handlers.
    
- **Rate limiting:** middleware or API gateway.
    

---

## Author

Mateus Vinicius Bezerra — [LinkedIn](https://www.linkedin.com/in/mateus-vinicius-bezerra/)