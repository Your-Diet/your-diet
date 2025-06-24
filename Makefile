.PHONY: help run

help:
	@echo "Available commands:"
	@echo "  run   - Run the API server (requires environment variables)"
	@echo "  help  - Show this help message"

run:
	@docker-compose up -d && go run cmd/api/main.go 