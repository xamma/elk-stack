FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o goapp

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/goapp .

EXPOSE 9090

CMD ["./goapp"]
