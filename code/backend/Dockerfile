FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /xyfamily ./cmd/xyfamily/

FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql-client
WORKDIR /app
COPY --from=builder /xyfamily .
COPY configs/ ./configs/
COPY migrations/ ./migrations/
EXPOSE 8080
CMD ["./xyfamily"]
