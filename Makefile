APP_NAME=diary
GO=go
BUILD=$(GO) build
RUN=$(GO) run

clean:
	rm $(APP_NAME)

build:
	$(BUILD)

run: build
	./$(APP_NAME)
