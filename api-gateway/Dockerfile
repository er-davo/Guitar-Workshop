FROM golang:1.25-alpine AS builder

RUN apk --no-cache add bash gcc musl-dev

WORKDIR /api-gateway

COPY api-gateway/app/go.mod /api-gateway/

RUN go mod download

COPY api-gateway/app/ /api-gateway/

RUN go build -o main cmd/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /api-gateway/main /app/main
COPY --from=builder /api-gateway/static /app/static

COPY api-gateway/config/config.yaml /app/config.yaml
COPY migrations /app/migrations

ENV CONFIG_PATH=/app/config.yaml
ENV APP_MIGRATIONS_PATH=../migrations

CMD [ "/app/main" ]