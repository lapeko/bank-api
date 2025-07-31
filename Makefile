.PHONY: \
	migrate migrate-down sqlc \
	docker-all-up docker-db-up docker-all-down docker-all-clear-down docker-build \
	local-api-up helm-local-up helm-local-down helm-template \
	gen-mock-store


POSTGRES_URL=postgres://postgres:1234@localhost:5432/bank?sslmode=disable
MIGRATE_CMD=migrate -source file://internal/db/migration -database $(POSTGRES_URL) -verbose
DOCKER_COMPOSE_PATH=infra/docker-compose.yaml
APP_NAME=bank-api

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
	docker build -f infra/Dockerfile -t $(APP_NAME) .

local-api-up:
	export \
		POSTGRES_URL=$(POSTGRES_URL) \
		JWT_SECRET_KEY=my_secret \
		APP_PORT=3000 \
		&& go run ./cmd/api/...

helm-local-up:
	helm install $(APP_NAME) ./infra/k8s -f ./infra/k8s/values.local.yaml --namespace $(APP_NAME) --create-namespace
	@echo "wait for $(APP_NAME) to be run"
	@kubectl wait --for=condition=ready pod -l app=$(APP_NAME) -n $(APP_NAME) --timeout=60s
	@echo "$(APP_NAME) successfully run"
	kubectl port-forward svc/$(APP_NAME)-service 3000:80 -n $(APP_NAME)

helm-local-down:
	helm uninstall $(APP_NAME) -n $(APP_NAME)

helm-template:
	 helm template bank-api ./infra/k8s -f ./infra/k8s/values.local.yaml --namespace $(APP_NAME) --create-namespace

gen-mock-store:
	mkdir -p internal/db/mockdb && \
	mockgen \
		-destination=internal/db/mockdb/store.go \
		-package=mockdb \
		github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc Store
