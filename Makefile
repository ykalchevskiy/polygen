.DEFAULT_GOAL := test
.PHONY: generate test test-coverage

generate:
	@go generate ./...

test: generate
	go test -v -count 1 ./...
	GOEXPERIMENT=jsonv2 go test -v -count 1 ./...

test-coverage: generate
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out
