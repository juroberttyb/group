export
BINARY_NAME=app

.PHONY: deps mocks unit-test integration-test full-test

deps:
	go mod tidy
	go list -m all

mocks:
	./mockery --all --inpackage

clean-test-cache:
	go clean -testcache

unit-test: clean-test-cache
	go test -v -cover -short ./...

integration-test: clean-test-cache
	go test -v -cover -run Integration ./...

full-test: clean-test-cache
	go test -v -cover ./...
