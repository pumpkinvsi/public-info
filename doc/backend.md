# Backend for Frontend

## Purpose
Provides a REST API consumed exclusively by the React frontend. 
Owns all data and is responsible for all interactions with internal services (like sending emails)

## API Routes

| Method | Path                  | Description                              |
|--------|-----------------------|------------------------------------------|
| GET    | /api/v1/bio           | Personal info                            |
| GET    | /api/v1/skills        | All skills with proficiency levels       |
| GET    | /api/v1/projects      | Projects grouped by technology           |
| GET    | /api/v1/technologies  | Technologies list                        |
| GET    | /api/v1/contacts      | Contacts                                 |
| POST   | /api/v1/email         | Sends email from frontend                |
### Observability

| Method | Path           | Description                                                   |
|--------|----------------|---------------------------------------------------------------|
| GET    | /metrics       | Prometheus metrics (text exposition format)                   |
| GET    | /health/live   | Liveness probe — process is alive (no I/O)                    |
| GET    | /health/ready  | Readiness probe — all dependencies reachable (DB, Kafka)      |

#### Health response schema

```json
{
  "status": "ok | degraded",
  "checks": {
    "database": { "status": "ok | degraded", "error": "optional message" }
  }
}
```

HTTP 200 → ready. HTTP 503 → database unreachable.

Kafka doesn't affect readiness, because sending emails functionality is not crucial for the system.

#### Prometheus metrics exposed

| Metric                                        | Type      | Labels                        |
|-----------------------------------------------|-----------|-------------------------------|
| `personal_page_http_requests_total`           | Counter   | `method`, `path`, `status`    |
| `personal_page_http_request_duration_seconds` | Histogram | `method`, `path`              |
| `personal_page_http_requests_in_flight`       | Gauge     | —                             |

## Configuration (environment variables)

| Variable      | Default        | Required | Description                        |
|---------------|----------------|----------|------------------------------------|
| SERVER_HOST   | 0.0.0.0        | No       | Bind address                       |
| SERVER_PORT   | 8080           | No       | Bind port                          |
| KAFKA_BROKERS | localhost:9092 | No       | Comma-separated broker list        |
| KAFKA_TOPIC   | email          | No       | Topic for outbound email messages  |
| DB_HOST       | localhost      | No       | PostgreSQL host                    |
| DB_PORT       | 5432           | No       | PostgreSQL port                    |
| DB_USER       | postgres       | No       | PostgreSQL user                    |
| DB_PASSWORD   | —              | **Yes**  | PostgreSQL password                |
| DB_NAME       | personal_page  | No       | PostgreSQL database name           |