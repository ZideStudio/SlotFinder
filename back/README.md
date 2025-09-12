# GoLang Boilerplate

## Dev Install

```bash
go install github.com/cosmtrek/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

## Start

```bash
air;
```

## Test

### All tests

```bash
go test ./... -v
```

### Specific package

```bash
go test ./[nom du package]
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
