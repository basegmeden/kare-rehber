FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server/main.go
RUN go build -o seed ./cmd/seed/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/seed .
COPY start.sh .
RUN chmod +x start.sh
CMD ["./start.sh"]
