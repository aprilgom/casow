# Repository Guidelines

## Project Structure & Module Organization

This repository is initialized as a Go module named `github.com/aprilgom/casow`. The current executable entry point is `cmd/casow/main.go`. Add new command entry points under `cmd/<name>/` and keep shared code in focused packages as the project grows. Put tests next to the package they exercise using Go's standard `*_test.go` naming.

Avoid committing build outputs, local binaries, coverage files, or generated artifacts unless they are intentionally part of the source tree.

## Build, Test, and Development Commands

- `go test ./...`: run all Go packages and tests.
- `go fix ./...`: apply Go source updates for the active toolchain.
- `golangci-lint run`: run static analysis configured by golangci-lint defaults or future project config.
- `git config core.hooksPath .githooks`: enable the repository hooks in a fresh checkout.

## Coding Style & Naming Conventions

Use standard Go formatting and idioms. Keep package names short, lowercase, and focused on one responsibility. Prefer explicit, descriptive exported names and keep unexported helpers local to the package that uses them.

Run `go test ./...` before pushing. The pre-commit hook runs `go fix ./...` and `golangci-lint run`, but tests should still be run manually until CI is added.

## Testing Guidelines

Write tests with Go's built-in `testing` package unless the repository adopts a specific helper library later. Name tests after observable behavior, for example `TestParserRejectsEmptyInput`. Prefer table tests for related cases and keep fixtures small.

## Commit & Pull Request Guidelines

Use Conventional Commit subjects written in Korean: `<type>(<scope>): <description>`. Common types are `feat`, `fix`, `docs`, `test`, `refactor`, `build`, and `chore`.

Examples: `feat(api): 상태 조회 엔드포인트 추가`, `fix(cli): 빈 입력 처리`, `docs: 개발 명령 문서화`.

Use branch names such as `feat/api-status-endpoint` or `fix/cli-empty-input`. Pull request titles should match the main commit subject and include verification evidence such as `go test ./...` and `.githooks/pre-commit`.

## Quality Gates

The repository uses `.githooks/pre-commit`. Enable it with:

```sh
git config core.hooksPath .githooks
```

On commit, the hook runs `go fix ./...` and then `golangci-lint run`. If `go fix` modifies files, review and stage those changes before retrying the commit.
