# 
FROM golang:alpine3.24 AS builder

WORKDIR /app

COPY . .

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine

RUN apk update && apk add ca-certificates --no-cache

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["./server"]