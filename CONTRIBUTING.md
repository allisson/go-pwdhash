# Contributing to go-pwdhash

Thanks for taking the time to contribute!

The following is a set of guidelines for contributing to `go-pwdhash`, which is hosted on GitHub. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Local Environment Setup

1.  **Fork and Clone:** Fork the repository on GitHub and clone it locally.
2.  **Go Version:** Ensure you are using Go 1.24 or later (as specified in `go.mod`).

## Development Workflow

Before submitting a Pull Request, please ensure your changes pass the test suite and linter. We use `Makefile` targets to simplify this.

### Running Tests

Execute the unit tests using standard Go tooling or the provided `make` target:

```bash
make test
# OR
go test -covermode=count -coverprofile=count.out -v ./...
```

### Running Linters

We use `golangci-lint` to ensure code quality.

```bash
make lint
# OR
golangci-lint run -v --fix
```

## Adding New Algorithms

`go-pwdhash` currently focuses explicitly on Argon2id.

Proposals for new algorithms will be considered but must meet strict criteria to avoid bloating the library or weakening the security guarantees:
*   The algorithm must be a widely recognized cryptographic standard (e.g., PHC winners).
*   The proposal must include a strong rationale for why Argon2id is no longer sufficient.
*   The implementation must include end-to-end test coverage.
*   The new hasher must integrate cleanly with the existing registry and options system.

We will generally not accept pull requests that add legacy algorithms (e.g., bcrypt, scrypt, PBKDF2).

## Pull Request Process

1.  Ensure any exported types, functions, and methods are well-documented following Go standards.
2.  Add tests. Your patch should maintain or increase test coverage. We prefer table-driven tests.
3.  Ensure your code is formatted with `gofmt` or `goimports`.
4.  Push your changes and submit a Pull Request.

Your PR will be reviewed as soon as someone is available!
