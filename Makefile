ifeq ($(OS),Windows_NT)
    GO_PATH=E:/go
else
    GO_PATH=/mnt/e/go
endif

MODULE=github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes

MIGRATIONS_PATH=storage/migrations
DB_CONTAINER_NAME=postgres_alpine

DB_USER=root
DB_PASSWORD=1234
DB_PORT=5432

DB_NAME=simple_bank
DB_TEST_NAME=simple_bank_test
DB_MANUAL_TEST_NAME=manual_test

DB_CONNECTION_URL=postgres://root:1234@127.0.0.1:5432/$(DB_NAME)?sslmode=disable
DB_TEST_CONNECTION_URL=postgres://root:1234@127.0.0.1:5432/$(DB_TEST_NAME)?sslmode=disable
DB_MANUAL_TEST_CONNECTION_URL=postgres://root:1234@127.0.0.1:5432/$(DB_MANUAL_TEST_NAME)?sslmode=disable

.PHONY:
	migrate docker-create docker-drop docker-start docker-stop db-create dp-drop migrate-gen migrate-up migrate-down migrate-up1 migrate-down1 sqlc test mockgen

docker-create:
	docker run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:alpine

docker-drop:
	docker rm $(DB_CONTAINER_NAME)

docker-start:
	docker start $(DB_CONTAINER_NAME)

docker-stop:
	docker stop $(DB_CONTAINER_NAME)

db-create:
	docker exec -i $(DB_CONTAINER_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)
	docker exec -i $(DB_CONTAINER_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_TEST_NAME)
	#docker exec -i $(DB_CONTAINER_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_MANUAL_TEST_NAME)

dp-drop:
	docker exec -i $(DB_CONTAINER_NAME) dropdb --username=$(DB_USER) $(DB_NAME)
	docker exec -i $(DB_CONTAINER_NAME) dropdb --username=$(DB_USER) $(DB_TEST_NAME)
	#docker exec -i $(DB_CONTAINER_NAME) dropdb --username=$(DB_USER) $(DB_MANUAL_TEST_NAME)

migrate-gen:
	$(GO_PATH)/bin/migrate.exe create -ext sql -dir $(MIGRATIONS_PATH) -seq add_users_table #last argument will be used as part of the name

migrate-up:
	$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_CONNECTION_URL) up
	$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_TEST_CONNECTION_URL) up
	#$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_MANUAL_TEST_CONNECTION_URL) up

migrate-down:
	$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_CONNECTION_URL) down
	$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_TEST_CONNECTION_URL) down
	#$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_MANUAL_TEST_CONNECTION_URL) down

migrate-up1:
	#$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_CONNECTION_URL) up 1
	$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_TEST_CONNECTION_URL) up 1
	#$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_MANUAL_TEST_CONNECTION_URL) up 1

migrate-down1:
	$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_CONNECTION_URL) down 1
	$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_TEST_CONNECTION_URL) down 1
	#$(GO_PATH)/bin/migrate.exe -path $(MIGRATIONS_PATH) -database $(DB_MANUAL_TEST_CONNECTION_URL) down 1

sqlc:
	cd ./storage/sqlc && $(GO_PATH)/bin/sqlc.exe generate

test:
	go test -v -cover ./...

mockgen:
	$(GO_PATH)/bin/mockgen.exe -build_flags=--mod=mod -package mockdb -destination ./storage/repository/mocks/store.go $(MODULE)/storage/repository Store