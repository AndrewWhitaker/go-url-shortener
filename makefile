docker-build:
	docker compose build --progress=plain

docker-up: docker-build
	docker compose up
