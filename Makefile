.PHONY: \
	migrate migrate-down sqlc \
	docker-all-up docker-db-up docker-all-down docker-all-clear-down docker-build \
	local-api-up helm-local-up helm-local-down \
	gen-mock-store


POSTGRES_URL=postgres://postgres:1234@localhost:5432/bank?sslmode=disable
MIGRATE_CMD=migrate -source file://internal/db/migration -database $(POSTGRES_URL) -verbose
DOCKER_COMPOSE_PATH=infra/docker-compose.yaml

migrate:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

sqlc:
	sqlc generate

docker-all-up:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d

docker-db-up:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d db

docker-all-down:
	docker compose -f $(DOCKER_COMPOSE_PATH) down

docker-all-clear-down:
	docker compose -f $(DOCKER_COMPOSE_PATH) down -v

# FOR EXAMPLE
docker-build:
	docker build -f infra/Dockerfile -t bank-api .

local-api-up:
	export \
		POSTGRES_URL=$(POSTGRES_URL) \
		JWT_SECRET_KEY=my_secret \
		APP_PORT=3000 \
		&& go run ./cmd/api/...

helm-local-up:
	helm install bank-api ./infra/k8s
	@echo "wait for bank-api running"
	@kubectl wait --for=condition=ready pod -l app=bank-api --timeout=60s
	@echo "bank-api successfully run"
	kubectl port-forward svc/bank-api-service 3000:80

helm-local-down:
	helm uninstall bank-api

gen-mock-store:
	mkdir -p internal/db/mockdb && \
	mockgen \
		-destination=internal/db/mockdb/store.go \
		-package=mockdb \
		github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc Store
