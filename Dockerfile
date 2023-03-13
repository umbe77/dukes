# build stage
FROM golang:1.20.2-alpine3.17 AS builder
WORKDIR /go/dukes
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go test ./...
RUN go build -v -o bin/dukes main.go

# final stage
FROM alpine:3.17
LABEL Name=dukes Version=0.0.1
WORKDIR /opt/dukes
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/dukes/bin .
ENTRYPOINT ./dukes
