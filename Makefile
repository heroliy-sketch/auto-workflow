APP_NAME := auto-workflow
BUILD_DIR := bin

.PHONY: build test run fmt

build:
	mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME) .

test:
	go test ./...

run: build
	AUTO_WORKFLOW_CONFIG=config.yaml $(BUILD_DIR)/$(APP_NAME)

fmt:
	gofmt -w main.go config/*.go internal/httpapi/*.go
