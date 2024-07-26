.PHONY: build clean test run

APP   := migrate
BUILD := ./build

clean:
	rm -Rf $(BUILD)

install:
	go mod tidy
	go get github.com/mitchellh/gox

run:
	go run .

build: clean
	GOOS=windows go build -o build/windows/$(APP).exe
	GOOS=darwin  go build -o build/darwin/$(APP)
	GOOS=linux   go build -o build/linux/$(APP)
	cd build/darwin  && zip migrate-darwin.zip  migrate
	cd build/linux   && zip migrate-linux.zip   migrate
	cd build/windows && zip migrate-windows.zip migrate.exe

docker-up:
	docker compose up -d

docker-down:
	docker compose kill
	docker compose rm --force

test:
	$(MAKE) docker-up
	go test -v ./... -count=1
	$(MAKE) docker-down
