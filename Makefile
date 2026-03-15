.PHONY: build build-windows build-linux build-all clean help
DEFAULT_TARGET := build
OUTPUT_DIR ?= ./bin
PROJECT_NAME := VirtualMemoryManagement
help:
	@echo "$(PROJECT_NAME) Build Targets"
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build for current platform"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make help           - Show this help"
	@echo "Usage:"
	@echo "  make build OUTPUT_DIR=./output"
build:
	@echo "Building for current platform..."
	@bash ./build-universal.sh $(OUTPUT_DIR)
clean:
	@echo "Cleaning build artifacts..."
	@go clean -cache
	@rm -rf $(OUTPUT_DIR)
	@echo "Clean completed"

