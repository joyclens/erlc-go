# Contributing to erlc-go

Thank you for contributing to erlc-go.

## Development

### Code Style

- Follow Go conventions from [Effective Go](https://golang.org/doc/effective_go)
- Format with `gofmt`: `go fmt ./...`
- Organize imports with `goimports`
- Run tests before submitting: `go test ./...`

### Testing

Write tests for all new functionality:

```go
func TestFeature(t *testing.T) {
	// Test implementation
}
```

Run tests with coverage:

```bash
go test -cover ./...
```

### Commit Messages

Use clear, concise messages:

```
feat: add rate limiting to API
fix: correct cache eviction logic
docs: update README with examples
test: add tests for error handling
```

Format: `type: description`

Types:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `test` - Tests
- `refactor` - Code restructuring
- `perf` - Performance improvement
- `chore` - Dependencies, build, etc.

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/name`
3. Make changes and test: `go test ./...`
4. Format code: `go fmt ./...`
5. Commit with clear messages
6. Push to your fork
7. Open a pull request

## Reporting Issues

Include:
- Go version: `go version`
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior
- Error messages or logs

## Code Review

Reviewed PRs typically respond within 24-48 hours. Maintain respectful discussion and address feedback constructively.

## License

By contributing, you agree your changes are licensed under Apache License 2.0.
