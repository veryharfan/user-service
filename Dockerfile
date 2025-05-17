# Builder stage
FROM golang:1.24-alpine3.20 AS builder
WORKDIR /go/src/user-service
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service cmd/main.go

# Final image
FROM alpine:latest
EXPOSE 80

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /go/src/user-service/user-service .

CMD ["./user-service"]
