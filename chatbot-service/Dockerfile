FROM golang:1.23-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o chatbot-service ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/chatbot-service .
COPY .env .env
EXPOSE 8084
CMD ["./chatbot-service"]
