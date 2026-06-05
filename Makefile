.PHONY: help dev local stop docker-up docker-down compose-up compose-down migrate-up migrate-down migrate-status dev-backend dev-frontend test test-coverage build clean lint install-goose

DB_URL ?= postgres://forte:forte123@localhost:14045/forte_commerce?sslmode=disable

help:
	@echo "forte-commerce Makefile"
	@echo ""
	@echo "  make dev               Stop everything, then start all services (infra -> BE -> FE) in Docker"
	@echo "  make local             Start infra in Docker, run BE + FE locally"
	@echo "  make stop              Stop all containers and kill processes on :8080 and :3000"
	@echo "  make docker-up         Start postgres + rabbitmq (infra only)"
	@echo "  make docker-down       Stop all containers"
	@echo "  make compose-up        Start all services via docker-compose"
	@echo "  make compose-down      Stop all containers and remove volumes"
	@echo "  make install-goose     Install goose migration tool"
	@echo "  make migrate-up        Apply all pending migrations"
	@echo "  make migrate-down      Roll back last migration"
	@echo "  make migrate-status    Show migration status"
	@echo "  make dev-backend       Run backend in dev mode (local)"
	@echo "  make dev-frontend      Run frontend in dev mode (local)"
	@echo "  make test              Run backend tests"
	@echo "  make test-coverage     Run tests with coverage report (requires ≥90%)"
	@echo "  make build             Build backend binary"
	@echo "  make clean             Clean build artifacts"
	@echo "  make lint              Run go vet"

stop:
	-docker compose down
	-lsof -ti:8080 | xargs kill -9 2>/dev/null
	-lsof -ti:3000 | xargs kill -9 2>/dev/null
	-pkill -f "go run ./cmd/server" 2>/dev/null || true
	-pkill -f "next dev" 2>/dev/null || true

dev: stop
	@echo "==> Starting infra (postgres + rabbitmq)..."
	docker compose up -d postgres rabbitmq
	@echo "==> Waiting for postgres to be ready..."
	@until docker compose exec -T postgres pg_isready -U forte -p 14045 >/dev/null 2>&1; do printf '.'; sleep 1; done
	@echo " postgres ready."
	@echo "==> Starting backend..."
	docker compose up -d --build backend
	@echo "==> Waiting for backend on :8080..."
	@until nc -z localhost 8080 2>/dev/null; do printf '.'; sleep 1; done
	@echo " backend ready."
	@echo "==> Starting frontend..."
	docker compose up -d --build frontend
	@echo ""
	@echo "All services up:"
	@echo "  Frontend : http://localhost:3000"
	@echo "  Backend  : http://localhost:8080"
	@echo "  RabbitMQ : http://localhost:15672"

local: stop
	@echo "==> Starting infra (postgres + rabbitmq)..."
	docker compose up -d postgres rabbitmq
	@echo "==> Waiting for postgres to be ready..."
	@until docker compose exec -T postgres pg_isready -U forte -p 14045 >/dev/null 2>&1; do printf '.'; sleep 1; done
	@echo " postgres ready."
	@echo "==> Waiting for RabbitMQ to be ready..."
	@until [ "$$(docker inspect --format='{{.State.Health.Status}}' forte_rabbitmq 2>/dev/null)" = "healthy" ]; do printf '.'; sleep 1; done
	@echo " rabbitmq ready."
	@echo "==> Starting backend locally..."
	cd backend && go run ./cmd/server &
	@echo "==> Waiting for backend on :8080..."
	@until nc -z localhost 8080 2>/dev/null; do printf '.'; sleep 1; done
	@echo " backend ready."
	@echo "==> Starting frontend locally..."
	cd frontend && npm run dev &
	@echo ""
	@echo "All services up:"
	@echo "  Frontend : http://localhost:3000"
	@echo "  Backend  : http://localhost:8080"
	@echo "  RabbitMQ : http://localhost:15672"
	@echo ""
	@echo "Press Ctrl+C to stop. Run 'make stop' to clean up."

docker-up:
	docker compose up -d postgres rabbitmq

docker-down:
	docker compose down

compose-up:
	docker compose up -d --build

compose-down:
	docker compose down -v

install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate-up: install-goose
	goose -dir backend/migrations postgres "$(DB_URL)" up

migrate-down: install-goose
	goose -dir backend/migrations postgres "$(DB_URL)" down

migrate-status: install-goose
	goose -dir backend/migrations postgres "$(DB_URL)" status

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev

test:
	cd backend && go test -race ./... -v -count=1

test-coverage:
	cd backend && go test -race ./... -coverprofile=coverage.out -covermode=atomic
	@COVERAGE=$$(go tool cover -func=backend/coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $${COVERAGE}%"; \
	if [ $$(echo "$${COVERAGE} < 90" | bc -l) -eq 1 ]; then \
		echo "ERROR: Coverage $${COVERAGE}% is below the required 90%"; exit 1; \
	fi
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: backend/coverage.html"

build:
	mkdir -p bin
	cd backend && go build -o ../bin/server ./cmd/server

clean:
	rm -rf bin/
	cd backend && go clean ./...

lint:
	cd backend && go vet ./...
