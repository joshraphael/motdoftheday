APP_NAME=diary
GO=go
BUILD=$(GO) build
RUN=$(GO) run

clean:
	rm -rf vendor/
	rm -f go.sum
	go clean --modcache
	rm -f $(APP_NAME)

build: export GO111MODULE=on
build:
	$(GO) mod vendor
	$(BUILD) -o $(APP_NAME) ./cmd/diary/main.go

run: build
	./$(APP_NAME)

db-reset:
	rm -f diary.db
	sqlite3 diary.db < sql/schema.sql
	sqlite3 diary.db < sql/data.sql

db-clean:
	sqlite3 diary.db < sql/clean.sql

db-populate:
	sqlite3 diary.db < sql/data.sql