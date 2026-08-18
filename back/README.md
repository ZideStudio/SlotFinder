# GoLang Boilerplate

## Dev Install

Clone the env `.env.model` file to `.env` and modify the variables as needed.

Install the toolchain from the repository root with [mise](https://mise.jdx.dev/):

```bash
mise install
```

## Start

```bash
go tool air
```

## Test

Clone `.env.test.model` to `.env.test` to point tests at `slotfinder_test` instead of your dev database (falls back to `.env` if absent):

```bash
cp .env.test.model .env.test
```

`slotfinder_test` is auto-created on a fresh `docker-compose.dev.yml` volume via `POSTGRES_MULTIPLE_DATABASES`. Otherwise, create it once: `docker exec -it <postgres_container> createdb -U slotfinder slotfinder_test`.

### All tests

```bash
go test ./... -v
```

### Specific package

```bash
go test ./[package]
```

### All tests with coverage

```bash
go test ./... -coverpkg=./... --cover
```

## Docker Build

```bash
docker build -t slotfinder-back .
docker run -it --rm --name slotfinder-back slotfinder-back
```
