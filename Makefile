include dev.env

# CONFIG 
MIGRATIONS_PATH=./pkg/database/migrations
DOCKER_CONTAINER_NAME=ecom_api
DOCKER_IMAGES=postgres:alpine

# DOCKER
docker-create:
	@docker run --name $(DOCKER_CONTAINER_NAME) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p $(DB_PORT):5432 -d $(DOCKER_IMAGES)

# MIGRATE
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_URL) up

migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_URL) down $(filter-out $@,$(MAKECMDGOALS))
	
migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_URL) force 1