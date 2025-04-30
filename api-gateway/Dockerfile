FROM golang:alpine AS builder

WORKDIR /api-gateway

COPY app /api-gateway

RUN go mod download
RUN apk --no-cache add bash gcc musl-dev

RUN go build -o main cmd/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /api-gateway/main /app/main
COPY --from=builder /api-gateway/static /app/static

CMD [ "/app/main" ]