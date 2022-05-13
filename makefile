name = url-shortener

.PHONY: build
build:
	go build -o $(name)

.PHONY: clean
clean:
	rm $(name)

.PHONY: docker-build
docker-build:
	docker compose build

.PHONY: docker-run
docker-run: docker-build
	docker compose up
	
