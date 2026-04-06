# Contributing

## Development Setup

### Prerequisites

- Go 1.23+
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)

### Getting Started

```bash
git clone https://github.com/rolloverdotdev/rollover-go.git
cd rollover-go
go mod download
```

### Regenerating the Client

Generated from the OpenAPI spec using oapi-codegen.

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
oapi-codegen -generate types,client -package client ./openapi.json > client/client.gen.go
```

### Building

```bash
go build ./...
```

## Pull Requests

Keep changes focused and atomic, test locally before opening a PR.

### Commit Messages

Lowercase, start with a verb, single line.

```
add context timeout to client methods
fix response parsing for list endpoints
update openapi spec and regenerate client
```

## License

By contributing, you agree your contributions will be licensed under the [MIT License](../LICENSE).
