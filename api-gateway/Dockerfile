FROM golang:alpine AS builder

RUN apk --no-cache add bash gcc musl-dev

WORKDIR /api-gateway

COPY app/go.mod /api-gateway/

RUN go mod download

COPY app/ /api-gateway/

RUN go build -o main cmd/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /api-gateway/main /app/main
COPY --from=builder /api-gateway/static /app/static
COPY app/config /app/config

ENV CONFIG_PATH=/app/config

CMD [ "/app/main" ]