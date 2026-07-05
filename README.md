# FlowerDoro API

Go + Postgres API for storing FlowerDoro users, flower catalog data, gardens, and earned flower rewards.

## Stack

- Go 1.25
- Gin HTTP API
- GORM + PostgreSQL
- GoAdmin admin UI at `/admin`
- Docker Compose for development and production

## Local Development

```sh
cp .env.example .env
docker compose -f docker-compose.dev.yml up --build
```

Run migrations explicitly when needed:

```sh
make migrate-up
```

Open:

- API health: http://localhost:8080/healthz
- API root: http://localhost:8080/api/v1
- Admin UI: http://localhost:8080/admin

GoAdmin initializes its own admin tables in Postgres on first boot.
The development Postgres container is exposed on `localhost:5433` by default to avoid clashing with a local Postgres on `5432`.

Default GoAdmin login:

- Username: `admin`
- Password: `admin`

## Production

Create a production env file:

```sh
cp .env.example .env.prod
```

Set at least:

```sh
POSTGRES_PASSWORD=change-me
HTTP_PORT=8080
FLOWER_DORO_API_IMAGE=ghcr.io/aquariues/flower-doro-api:latest
```

Run:

```sh
docker compose --env-file .env.prod -f docker-compose.prod.yml up --build -d
```

## API

### Health

```http
GET /healthz
```

### Flower Catalog

```http
GET /api/v1/flowers
GET /api/v1/flowers/:kind
POST /api/v1/flowers
```

### Users

```http
POST /api/v1/users
GET /api/v1/users/:id/garden
POST /api/v1/users/:id/flowers
```

`POST /api/v1/users/:id/flowers` records an earned flower after a qualifying focus session. The API rejects reward records shorter than 30 focus minutes.

## Notes

- Postgres is the only supported runtime database.
- `AUTO_MIGRATE=true` creates and updates API tables at startup.
- Versioned SQL migrations live in `migrations/` and are embedded for app startup.
- GoAdmin also stores its admin metadata in the same Postgres database.

## Contributor Docs

- Agent/project guide: `CLAUDE.md`
- Git workflow: `docs/GIT_WORKFLOW.md`
- Release playbook: `docs/RELEASE_SKILL.md`
