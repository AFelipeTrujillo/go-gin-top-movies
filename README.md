# Go-Gin-Top-Movies

A **Clean Architecture** REST API built with **Go**, **Gin**, and **GORM** that serves IMDB top-rated movie data from a SQLite database.

## Quick Start

```bash
go run cmd/app/main.go
```

Server starts at `http://localhost:8080`.

## API Endpoints

| Method | Path                 | Description                      | Query Params                 |
|--------|----------------------|----------------------------------|------------------------------|
| GET    | `/imdb`              | List movies (paginated)          | `page`, `page_size`          |
| GET    | `/imdb/top`          | Top rated movies                 | `n` (default 10)             |
| GET    | `/imdb/search`       | Search by title or overview      | `q`                          |
| GET    | `/imdb/genre`        | Filter by genre                  | `g` (e.g. `Action`)          |
| GET    | `/imdb/year`         | Filter by release year           | `y` (e.g. `2020`)            |
| GET    | `/imdb/:title`       | Get movie by exact title         | —                            |

## Project Structure

```
cmd/                    → Entry point & dependency injection
internal/
  Application/          → Use cases & DTOs
  Domain/               → Entities, interfaces, errors (pure Go, no frameworks)
  Infrastructure/       → GORM repository, HTTP handlers, config
tests/                  → Unit & integration tests
```

- **Domain** has zero external dependencies.
- **Use cases** only depend on domain interfaces.
- **Handlers** are thin — parse requests, call use cases, return responses.

## Database

Pre-populated SQLite database with IMDB top movies data at `internal/Infrastructure/Persistence/movies.sqlite`. No migration needed — the table already exists.

## Configuration

Environment variables (with defaults):

| Variable       | Default | Description                    |
|----------------|---------|--------------------------------|
| `PORT`         | `8080`  | HTTP server port               |
| `DB_PATH`      | *(path to movies.sqlite)* | Database file path |
| `PAGE_SIZE`    | `10`    | Default page size              |
| `MAX_PAGE_SIZE`| `100`   | Max items per page             |

## Tests

```bash
go test ./tests/...
```