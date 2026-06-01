
APP_NAME=w
BIN_DIR=bin
BIN_PATH=$(BIN_DIR)/$(APP_NAME)
DB_PATH=./data.db
MIGRATIONS_DIR=./store/sqlite/migrations

.PHONY: build run migrate migrate-down clean

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) ./cmd/w

run: build
	$(BIN_PATH) serve
migrate:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) sqlite3 $(DB_PATH) down

clean:
	rm -f $(BIN_PATH)
