# Release Skill

Use this release playbook when preparing a FlowerDoro API release.

## Goal

Ship a tested version from `develop` to `main`, tag it, and create a GitHub Release with clear notes.

## Required Inputs

- Release version, for example `v0.1.0`.
- Release notes file under `docs/releases/`.
- Confirmation that `develop` contains the intended changes.

## Preflight

```sh
git status --short --branch
git switch develop
git pull --ff-only
gofmt -w cmd internal
go test ./...
docker build -t flower-doro-api:local .
```

Run the dev stack:

```sh
docker compose -f docker-compose.dev.yml up --build -d
curl -fsS http://localhost:8080/healthz
curl -fsS http://localhost:8080/api/v1
curl -fsSI http://localhost:8080/admin
```

Smoke test reward validation:

- `focus_minutes < 30` should return `400`.
- `focus_minutes >= 30` should return `201`.

## Merge

Open a PR from `develop` to `main`.

Merge only after checks pass and the release diff is expected.

## Tag

```sh
git switch main
git pull --ff-only
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

## GitHub Release

```sh
gh release create v0.1.0 \
  --title "v0.1.0" \
  --notes-file docs/releases/v0.1.0.md
```

## Production Deploy

Build and run with production compose:

```sh
cp .env.example .env.prod
docker compose --env-file .env.prod -f docker-compose.prod.yml up --build -d
```

Required production env values:

- `POSTGRES_PASSWORD`
- `HTTP_PORT`
- `FLOWER_DORO_API_IMAGE`

## Rollback

If the release is bad:

```sh
gh release delete v0.1.0
git push origin :refs/tags/v0.1.0
```

Then redeploy the previous known-good image or tag.

