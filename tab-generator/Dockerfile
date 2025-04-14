FROM golang:alpine AS builder

WORKDIR /tab-generator

COPY app /tab-generator

RUN go mod download
RUN apk --no-cache add bash gcc musl-dev

RUN go build -o main cmd/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=builder /tab-generator/main /app/main

RUN apk --no-cache add curl

CMD [ "/app/main" ]