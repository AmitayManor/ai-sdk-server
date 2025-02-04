FROM golang:1.23 as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o main .

FROM debian:stable-slim
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
