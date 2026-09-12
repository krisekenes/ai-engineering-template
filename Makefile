SHELL := /bin/sh
.PHONY: doctor check test dev fmt

doctor:
	@sh scripts/doctor.sh
check:
	@sh scripts/check.sh
test:
	cd backend && go test -race -count=1 ./...
dev:
	cd backend && go run ./cmd/server
fmt:
	cd backend && gofmt -w cmd internal
	cd examples/repair-loop && gofmt -w main.go main_test.go

.PHONY: frontier-demo
frontier-demo:
	cd examples/repair-loop && go run .
