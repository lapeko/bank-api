.PHONY: \
	migrate migrate-down sqlc \
	docker-all-up docker-db-up docker-all-down docker-all-clear-down docker-build \
	local-api-up helm-local-up helm-remote-up helm-down helm-template \
	gen-mock-store update-kubeconfig


POSTGRES_URL=postgres://postgres:1234@localhost:5432/bank?sslmode=disable
MIGRATE_CMD=migrate -source file://internal/db/migration -database $(POSTGRES_URL) -verbose
DOCKER_COMPOSE_PATH=infra/docker/docker-compose.yaml
CHART_PATH=./infra/k8s
APP_NAME=bank-api
IMAGE_TAG=79ac383ebeb049ca0760264640425d0d897cdca6

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
	helm upgrade --install $(APP_NAME) $(CHART_PATH) -f $(CHART_PATH)/values.local.yaml --namespace $(APP_NAME) --create-namespace --set deploy.image.tag=$(IMAGE_TAG)
	@echo "wait for $(APP_NAME) to be run"
	@kubectl wait --for=condition=ready pod -l app=$(APP_NAME) -n $(APP_NAME) --timeout=60s
	@echo "$(APP_NAME) successfully run"
	kubectl port-forward svc/$(APP_NAME)-service 3000:80 -n $(APP_NAME)

helm-remote-up:
	@if ! kubectl get crd externalsecrets.external-secrets.io >/dev/null 2>&1; then \
		echo "CRD not found, applying..."; \
		kubectl apply --server-side -f https://raw.githubusercontent.com/external-secrets/external-secrets/v0.19.0/deploy/crds/bundle.yaml; \
	else \
		echo "CRD already exists, skipping apply."; \
	fi
	helm upgrade --install $(APP_NAME) $(CHART_PATH) --namespace $(APP_NAME) --create-namespace --set deploy.image.tag=$(IMAGE_TAG)

helm-down:
	helm uninstall $(APP_NAME) -n $(APP_NAME)

helm-template:
	 helm template bank-api $(CHART_PATH) -f .$(CHART_PATH)/values.local.yaml --namespace $(APP_NAME) --create-namespace

make k8s-ecr-register:
	kubectl delete secret ecr-secret --namespace=bank-api --ignore-not-found
	kubectl create secret docker-registry ecr-secret \
		--docker-server=539247467338.dkr.ecr.eu-central-1.amazonaws.com \
		--docker-username=AWS \
		--docker-password="$$(aws ecr get-login-password --region eu-central-1)" \
		--namespace=bank-api

gen-mock-store:
	mkdir -p internal/db/mockdb && \
	mockgen \
		-destination=internal/db/mockdb/store.go \
		-package=mockdb \
		github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/db/sqlc Store

update-kubeconfig:
	aws eks --region eu-central-1 update-kubeconfig --name bank-api-eks
