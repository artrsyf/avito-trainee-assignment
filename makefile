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
		-v C:\projects\avito-trainee-assignment\tests\load:/app \
		loadimpact/k6 run /app/stress_test.js \
		-e API_URL=http://host.docker.internal:8080 \
		-e K6_OUT=influxdb=http://localhost:8086/k6
	docker-compose -f $(COMPOSE_DEV_FILE) down --volumes --remove-orphans

lint: 
	golangci-lint run

help:
	@echo "Доступные команды:"
	@echo "  run                - Запустить контейнеры"
	@echo "  rebuild            - Пересобрать контейнеры и запустить"
	@echo "  down               - Остановить контейнеры"
	@echo "  drop               - Остановить контейнеры и удалить тома"
	@echo "  unit_test          - Запустить unit-тесты с покрытием кода"
	@echo "  unit_cover         - Сгенерировать HTML-отчет по покрытию тестами"
	@echo "  integration_test   - Запустить интеграционные тесты в контейнерах"
	@echo "  e2e_test           - Запустить end-to-end тесты"
	@echo "  load_test          - Запустить нагрузочное тестирование с K6"
	@echo "  lint               - Запустить линтер golangci-lint"
	@echo "  help               - Показать доступные команды"

.PHONY: run rebuild down drop unit_test unit_cover integration_test e2e_test load_test lint help