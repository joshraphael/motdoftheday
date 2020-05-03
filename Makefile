APP_NAME=motdoftheday
GO=go
BUILD=$(GO) build
RUN=$(GO) run
CONFIG_ENV:=.config/local.yml
.EXPORT_ALL_VARIABLES:

clean:
	rm -rf vendor/
	rm -f go.sum
	go clean --modcache
	rm -f $(APP_NAME)

build: export GO111MODULE=on
build:
	$(GO) mod vendor
	$(BUILD) -o $(APP_NAME) ./cmd/$(APP_NAME)/main.go

run: build
	./$(APP_NAME)

db-reset:
	rm -f $(APP_NAME).db
	sqlite3 $(APP_NAME).db < sql/schema.sql
	sqlite3 $(APP_NAME).db < sql/data.sql

db-clean:
	sqlite3 $(APP_NAME).db < sql/clean.sql

db-populate:
	sqlite3 $(APP_NAME).db < sql/data.sql