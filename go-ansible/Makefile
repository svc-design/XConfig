# Go modules environment
GO ?= go

# 项目名称和可执行文件名
PROJECT_NAME := go-ansible
OUTPUT_FILE := bin/$(PROJECT_NAME)

.PHONY: all build run fmt lint clean

all: build

build:
	$(GO) build -o $(OUTPUT_FILE) ./cmd/main.go

run:
	$(OUTPUT_FILE)

fmt:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...

clean:
	rm -rf bin

