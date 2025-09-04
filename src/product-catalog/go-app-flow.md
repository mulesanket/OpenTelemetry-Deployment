# Go Application Development Flow (with Testing)

## 1) Project Setup
- Create project folder.
- Initialize module:
  go mod init my-service
- Result: creates `go.mod` (project identity + Go version).

## 2) Code Development
- Write source files (`main.go`, handlers, internal packages).
- Imports can be:
  - Standard library: fmt, net/http, os (no download needed).
  - External libs: e.g., github.com/gorilla/mux, github.com/sirupsen/logrus.

## 3) Dependency Management
- When you add/remove imports:
  - Option A: explicitly fetch:
      go get <module>
  - Option B: import in code and build/run; Go auto-fetches.
- Housekeeping (make it a habit before sharing/releasing):
    go mod tidy
  - Adds missing deps (that your code imports).
  - Removes unused deps.
  - Updates `go.sum` (checksums for reproducible builds).
- After tidy, `go.mod` + `go.sum` match exactly what your code uses.

## 4) Testing
- Put tests in files ending with `_test.go` in the same package.
- Run all tests:
    go test ./...
- Verbose:
    go test -v ./...
- Coverage:
    go test -cover ./...
- Pattern/one test:
    go test -run TestName -v

## 5) Build
- Compile + link into a single self-contained binary:
    go build -o my-app
- Binary contains:
  - Your code
  - Required dependencies
  - Minimal Go runtime glue

## 6) Run
- Run the binary directly (Go toolchain not required on target host):
    ./my-app

## 7) Deploy
- Ship the binary or containerize it.
- Typical Dockerfile (multi-stage) compiles statically, then runs on a small base.
- Advantages:
  - Fast startup
  - Portable artifact
  - Simple deploys (copy/run)

# Handy habits
- During quick local dev: go run main.go (fast feedback).
- Before commit/CI/release: go mod tidy && go test ./... && go build -o my-app
