# Start from the official Golang image
FROM golang:1.23-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o user-service ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/user-service .
COPY .env.example .env
EXPOSE 8081
CMD ["./user-service"]
