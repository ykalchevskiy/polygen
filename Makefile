.DEFAULT_GOAL := test
.PHONY: fmt lint generate test test-coverage

fmt:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 fmt
	GOEXPERIMENT=jsonv2 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 fmt

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run
	GOEXPERIMENT=jsonv2 go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run

generate:
	@go generate ./...

test: generate
	go test -v -count 1 ./...
	GOEXPERIMENT=jsonv2 go test -v -count 1 ./...

test-coverage: generate
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out
