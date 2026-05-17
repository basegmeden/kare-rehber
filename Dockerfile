FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/seed .
CMD ["sh", "-c", "if [ ! -f \"${DB_PATH:-kare_rehber.db}\" ]; then echo 'Seeding...'; ./seed; fi && exec ./server"]
