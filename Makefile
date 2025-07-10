.PHONY: migrate migrate_down

migrate:
	migrate -source file://db/migration -database postgres://postgres:1234@localhost:5432/bank?sslmode=disable -verbose up

migrate-down:
	migrate -source file://db/migration -database postgres://postgres:1234@localhost:5432/bank?sslmode=disable -verbose down
