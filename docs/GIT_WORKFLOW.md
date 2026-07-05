# Git Workflow

FlowerDoro API uses a branch workflow with `develop` as the integration branch and `main` as the release branch.

## Branches

- `main`: production-ready release branch.
- `develop`: integration branch for tested work before release.
- `feature/*`: new product/API features.
- `fix/*`: bug fixes.
- `chore/*`: maintenance, docs, dependency updates, CI, infra.
- `release/*`: optional stabilization branch before merging `develop` to `main`.

Use short descriptive names:

```sh
feature/flower-catalog-seed
fix/reward-minutes-validation
chore/docker-prod-hardening
```

## Daily Flow

Start from latest `develop`:

```sh
git switch develop
git pull --ff-only
git switch -c feature/my-change
```

Before opening a PR:

```sh
gofmt -w cmd internal
go test ./...
docker build -t flower-doro-api:local .
git status --short
```

Open PR into `develop`.

## Release Flow

When `develop` is ready:

```sh
git switch develop
git pull --ff-only
gofmt -w cmd internal
go test ./...
docker build -t flower-doro-api:local .
```

Open PR from `develop` into `main`.

After merge to `main`, create a release tag:

```sh
git switch main
git pull --ff-only
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
gh release create v0.1.0 --title "v0.1.0" --notes-file docs/releases/v0.1.0.md
```

## Commit Style

Use concise conventional commit style:

```text
feat: add flower catalog endpoint
fix: reject short focus rewards
chore: add prod compose file
docs: document release workflow
```

## Protection Rules To Add In GitHub

Recommended branch protection:

- Require PR before merging to `main`.
- Require PR before merging to `develop`.
- Require status checks once CI is enabled.
- Disallow force pushes on `main` and `develop`.
- Require linear history if the team wants a clean release log.

