
APP_NAME=w
BIN_DIR=bin
BIN_PATH=$(BIN_DIR)/$(APP_NAME)
DB_PATH=./data.db
MIGRATIONS_DIR=./store/sqlite/migrations

.PHONY: build build-ui run dev migrate migrate-down clean

build-ui:
	@echo "📦 Building frontend..."
	cd web && npm install && npm run build

build: build-ui
	@echo "🔨 Building Go binary with embedded UI..."
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) ./cmd/w

dev:
	@echo "🚀 Starting backend and frontend in dev mode..."
	@trap 'kill 0' EXIT; \
	$(BIN_PATH) serve & \
	sleep 2 && \
	cd web && npm run dev

run: build
	@echo "🚀 Starting server with embedded UI..."
	$(BIN_PATH) serve

migrate:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) down

clean:
	rm -f $(BIN_PATH)
	rm -rf web/build
