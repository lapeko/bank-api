FROM golang:1.24.5-alpine AS build
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bank-api ./cmd/api
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xz

FROM alpine:latest
COPY --from=build /app/infra/scripts/entrypoint.sh /app/entrypoint.sh
COPY --from=build /app/internal/db/migration /app/internal/db/migration
COPY --from=build /app/migrate /usr/bin/migrate
COPY --from=build /bank-api /usr/bin/bank-api

EXPOSE 3000
ENTRYPOINT ["/app/entrypoint.sh"]