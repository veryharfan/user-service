FROM golang:1.24-alpine3.20
WORKDIR /go/src/user-service
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service app/main.go

FROM alpine:latest
EXPOSE 8081

RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata
WORKDIR /app/
COPY --from=0 /go/src/user-service/user-service .
CMD ["./user-service"]
