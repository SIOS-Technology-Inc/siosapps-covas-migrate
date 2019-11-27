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

build:
	gox -osarch "darwin/amd64 linux/amd64 windows/amd64" -output "build/$(APP)-{{.OS}}"

docker-up:
	docker-compose up -d

docker-down:
	docker-compose kill
	docker-compose rm --force

test:
	$(MAKE) docker-up
	go test -v ./... -count=1
#	$(MAKE) docker-down
