# GoLang Boilerplate

## Dev Install

Clone the env `.env.model` file to `.env` and modify the variables as needed.

```bash
make install
```

## Start

```bash
make start
```

## Test

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
go test ./... --cover
```

## Docker Build

```bash
docker build -t slotfinder-back .
docker run -it --rm --name slotfinder-back slotfinder-back
```
