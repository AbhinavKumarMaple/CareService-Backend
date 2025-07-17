# Makefile for Caregiver Microservices Application

.PHONY: test test-unit test-integration test-coverage test-verbose clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  test           - Run all unit tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-integration - Run integration tests"
	@echo "  test-schedule  - Run schedule module tests only"
	@echo "  clean          - Clean test artifacts"
	@echo "  help           - Show this help message"

# Run all unit tests
test:
	go test ./src/...

# Run unit tests only (same as test for this project)
test-unit:
	go test ./src/...

# Run tests with verbose output
test-verbose:
	go test -v ./src/...

# Run tests with coverage (uses the existing script)
test-coverage:
ifeq ($(OS),Windows_NT)
	powershell -ExecutionPolicy Bypass -File ./scripts/run-tests.ps1
else
	./scripts/run-tests.sh
endif

# Run integration tests
test-integration:
	./scripts/run-integration-test.bash

# Run schedule module tests only
test-schedule:
	go test -v ./src/application/usecases/schedule ./src/infrastructure/rest/controllers/schedule

# Clean test artifacts
clean:
	rm -rf coverage/
	go clean -testcache

# Quick test with timeout
test-quick:
	go test -timeout 30s ./src/...