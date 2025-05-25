#!/bin/bash

all:
	@make build && make run

build:
	@echo "Building reconcile-service project..."
	@go build -o reconcile-service ./cmd/app/
	@echo "Done."

run:
	@echo "Running reconcile-service binary..."
	@./reconcile-service