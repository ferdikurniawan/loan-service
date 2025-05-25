#!/bin/bash

all:
	@make build && make run

build:
	@echo "Building loan-service project..."
	@go build -o loan-service ./cmd/app/
	@echo "Done."

run:
	@echo "Running loan-service binary..."
	@./loan-service