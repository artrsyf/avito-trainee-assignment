COMPOSE_DEV_FILE=./docker/docker-compose.yml
COMPOSE_TEST_FILE=./docker/docker-compose.integration.yml

run:
	docker-compose -f $(COMPOSE_DEV_FILE) up

rebuild:
	docker-compose -f $(COMPOSE_DEV_FILE) up --build

down:
	docker-compose -f $(COMPOSE_DEV_FILE) down

drop:
	docker-compose -f $(COMPOSE_DEV_FILE) down --volumes --remove-orphans

unit_test:
	go test -v \
		./internal/user/repository/postgres \
		./internal/user/usecase \
		./internal/session/repository/redis \
		./internal/session/usecase \
		./internal/transaction/repository/postgres \
		./internal/transaction/usecase \
		./internal/purchase/repository/postgres \
		./internal/purchase/usecase \
		-coverprofile=./docs/unit_coverage.out

unit_cover: unit_test
	go tool cover -html=./docs/unit_coverage.out -o ./docs/unit_coverage.html

integration_test:
	docker-compose -f $(COMPOSE_TEST_FILE) up --build -d
	go test -v ./tests/integration/
	docker-compose -f $(COMPOSE_TEST_FILE) down --volumes --remove-orphans

e2e_test:
	docker-compose -f $(COMPOSE_TEST_FILE) up --build -d
	go test -v ./tests/e2e/
	docker-compose -f $(COMPOSE_TEST_FILE) down --volumes --remove-orphans

load_test:
	docker-compose -f $(COMPOSE_DEV_FILE) up --build -d
	docker run -i --network=host \
		-v "$(CURDIR)/tests/load:/app" \
		grafana/k6 run /app/stress_test.js \
		-e API_URL=http://host.docker.internal:8080
	docker-compose -f $(COMPOSE_DEV_FILE) down --volumes --remove-orphans

lint: 
	golangci-lint run

help:
	@echo "Available commands:"
	@echo "  run                - Start containers"
	@echo "  rebuild            - Rebuild and start containers"
	@echo "  down               - Stop containers"
	@echo "  drop               - Stop containers and remove volumes"
	@echo "  unit_test          - Run unit tests with coverage"
	@echo "  unit_cover         - Generate HTML coverage report"
	@echo "  integration_test   - Run integration tests in containers"
	@echo "  e2e_test           - Run end-to-end tests"
	@echo "  load_test          - Run load testing with K6"
	@echo "  lint               - Run golangci-lint"
	@echo "  help               - Show available commands"

.PHONY: run rebuild down drop unit_test unit_cover integration_test e2e_test load_test lint help