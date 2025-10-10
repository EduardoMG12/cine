# Development Linting Guide

## Pre-commit Hooks (lint-staged)

The project uses `lint-staged` with `husky` to automatically format code before commits.

### Current Configuration

- **JavaScript/TypeScript**: Uses `biome format`
- **Go**: Uses `gofmt` for formatting
- **Dart/Flutter**: Uses `flutter format`

### Go Linting Notes

Due to module resolution issues with local packages, `go vet` is temporarily disabled in pre-commit hooks. Use the manual script for full linting during development.

## Manual Linting

### Go Projects

Use the provided script for comprehensive Go linting:

```bash
./scripts/go-lint.sh
```

This script will:
- Format code with `gofmt`
- Run `go vet` (may fail during development)
- Clean up modules with `go mod tidy`
- Check for build errors

### Individual Commands

```bash
# Format Go code
cd api_v2 && gofmt -w .

# Check formatting
cd api_v2 && gofmt -l .

# Run vet (may fail with local modules)
cd api_v2 && go vet ./...

# Build check
cd api_v2 && go build ./...
```

## Troubleshooting

### Go Module Issues

If you encounter module resolution errors:

1. Make sure you're in the correct directory (`api_v2/`)
2. Run `go mod tidy` to clean up dependencies
3. For local development, use the manual linting script

### Pre-commit Hook Failures

If commits fail due to linting:

1. Run the appropriate formatter manually
2. Stage the formatted files
3. Retry the commit

## Future Improvements

- [ ] Fix `go vet` integration with local modules
- [ ] Add `golangci-lint` for comprehensive Go linting
- [ ] Add test coverage checks to pre-commit hooks
