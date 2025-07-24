.PHONY: migrate migrate_down sqlc

migrate:
	migrate -source file://internal/db/migration -database postgres://postgres:1234@localhost:5432/bank?sslmode=disable -verbose up

migrate-down:
	migrate -source file://internal/db/migration -database postgres://postgres:1234@localhost:5432/bank?sslmode=disable -verbose down

sqlc:
	sqlc generate

start:
	export \
		POSTGRES_URL=postgres://postgres:1234@localhost:5432/bank \
		JWT_SECRET_KEY=my_secret \
		APP_PORT=3000 \
		&& go run ./cmd/api/...