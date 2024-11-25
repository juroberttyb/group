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

unit-test:
	go test -v -cover -short ./...

integration-test:
	go test -v -cover -run Integration ./...

full-test:
	go test -v -cover ./...
