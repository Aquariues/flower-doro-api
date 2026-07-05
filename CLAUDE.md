# CLAUDE.md

This file gives future coding agents and contributors the local rules for working on the FlowerDoro API.

## Project

FlowerDoro API stores the flower catalog, users, focus sessions, and earned garden flowers for the FlowerDoro macOS/iPhone app.

Core stack:

- Go 1.25
- Gin HTTP API
- GORM
- PostgreSQL
- GoAdmin at `/admin`
- Docker Compose for dev and prod

## Local Paths

Repository path:

```sh
/Users/aquariues/aquariues/projects/flower-doro-api
```

Main app repository:

```sh
/Users/aquariues/aquariues/projects/flower-doro
```

## Development Commands

Run the dev stack:

```sh
docker compose -f docker-compose.dev.yml up --build
```

Run checks:

```sh
gofmt -w cmd internal
go test ./...
docker build -t flower-doro-api:local .
```

Useful endpoints:

- API health: `http://localhost:8080/healthz`
- API root: `http://localhost:8080/api/v1`
- GoAdmin: `http://localhost:8080/admin`

Default local GoAdmin login:

- Username: `admin`
- Password: `admin`

## Product Rules

- Vietnamese is the default app locale.
- A user should only receive a flower reward after at least 30 completed focus minutes.
- The server must also enforce the 30-minute minimum so older or modified clients cannot farm flowers.
- Postgres is the source of truth for API data.
- Keep the API lightweight and boring before adding infrastructure.

## Engineering Rules

- Prefer small, readable handlers and models.
- Keep database schema changes explicit and compatible with existing data.
- Do not commit `.env`, logs, uploads, database volumes, or local generated files.
- Run `gofmt` before committing Go changes.
- Run `go test ./...` before pushing.
- Docker builds should continue to work after dependency or config changes.

## GitHub Notes

The active GitHub account should be `Aquariues`.

If pushing workflow files under `.github/workflows`, the GitHub auth token needs the `workflow` scope. Without that scope, GitHub rejects the push.

