# GoLang Boilerplate

## Dev Install

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

## API Documentation

This project uses [swaggo/swag v2](https://github.com/swaggo/swag) to generate OpenAPI 3.1.0 documentation from Go annotations.

### Generate Swagger Documentation

The swagger documentation is automatically generated during development with Air. To manually regenerate:

```bash
swag init --v3.1
```

This will generate:
- `docs/docs.go` - Go package with embedded documentation
- `docs/swagger.json` - OpenAPI 3.1.0 JSON specification
- `docs/swagger.yaml` - OpenAPI 3.1.0 YAML specification

### Access Swagger UI

Once the server is running, access the interactive API documentation at:

```
http://localhost:3001/swagger/index.html
```

### Documentation Guidelines

- Add swagger annotations to your handler functions
- Use `@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, and `@Security` tags
- Document DTOs and models with proper struct tags and comments
- Run `swag init --v3.1` after making changes to API annotations

Example:
```go
// CreateAccount godoc
// @Summary Create an account
// @Description Create a new account with the provided parameters.
// @Tags Account
// @Accept json
// @Produce json
// @Param data body account.AccountCreateDto true "Account parameters"
// @Success 200
// @Failure 400 {object} helpers.ApiError
// @Router /v1/account [post]
func (c *AccountController) Create(ctx *gin.Context) {
    // implementation
}
```
