FROM golang:1.23 as builder
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y ca-certificates
RUN go mod tidy
RUN go build -o main .

FROM debian:stable-slim
WORKDIR /root/
COPY --from=builder /app/main .
RUN apt-get update && apt-get install -y ca-certificates
CMD ["./main"]
