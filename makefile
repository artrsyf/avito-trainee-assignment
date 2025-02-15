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

help:
	@echo "Доступные команды:"

.PHONY: run rebuild down drop unit_test unit_cover integration_test e2e_test help
