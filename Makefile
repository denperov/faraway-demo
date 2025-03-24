build-client:
	docker buildx build -f ./cmd/client/Dockerfile -t faraway-demo/client .

build-server:
	docker buildx build -f ./cmd/server/Dockerfile -t faraway-demo/server .

build: build-client build-server

run: build
	docker compose up