FROM golang:1.25-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/flower-doro-api ./cmd/server

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/flower-doro-api /app/flower-doro-api

EXPOSE 8080
CMD ["/app/flower-doro-api"]
