export
BINARY_NAME=app

.PHONY: ${BINARY_NAME} deps mocks unit-test integration-test full-test

${BINARY_NAME}:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/${BINARY_NAME} -ldflags "$(LDFLAGS)" ./main.go

deps:
	go mod tidy
	go list -m all
mocks:
	./mockery --all --inpackage

clean-test-cache:
	go clean -testcache

unit-test: clean-test-cache
	go test -v -cover -short -parallel 1 ./...

integration-test: clean-test-cache
	go test -v -cover -run Integration ./...

full-test: clean-test-cache
	go test -v -cover ./...
