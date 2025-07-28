FROM golang:1.24.5-alpine AS build
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
# RUN CGO_ENABLED=0 GOOS=linux go build -o /bank-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bank-api ./cmd/api

FROM alpine:latest
COPY --from=build /bank-api /usr/bin/bank-api
EXPOSE 3000
ENTRYPOINT ["bank-api"]