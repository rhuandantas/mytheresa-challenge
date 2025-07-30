# Etapa de build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api ./cmd/server/main.go

# Etapa final
FROM scratch
COPY --from=builder /app/api /api
EXPOSE 8484
ENTRYPOINT ["/api"]