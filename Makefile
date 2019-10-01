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
