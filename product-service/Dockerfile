FROM golang:1.23-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o product-service ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/product-service .
COPY .env.example .env
EXPOSE 8082
CMD ["./product-service"]
