name = url-shortener

.PHONY: build
build:
	go build -o $(name)

.PHONY: test
test:
	go test ./...

.PHONY: test-short
test-short:
	go test -short ./...

.PHONY: clean
clean:
	rm $(name)

.PHONY: docker-build
docker-build:
	docker compose build

.PHONY: docker-run
docker-run: docker-build
	docker compose up
	
