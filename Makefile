start: 
	go run cmd/server/main.go

#Build image with docker-compose and run injecting .env file
docker-dev:
	docker compose -f docker-compose-dev.yaml --env-file .env up --build 

build:
	docker compose build --no-cache 
run:
	docker compose up -d
