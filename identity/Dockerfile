FROM golang:1.21.1-alpine3.18 as builder
WORKDIR /app
COPY identity identity
COPY common common
WORKDIR identity
RUN go build -o service .

FROM alpine:latest
COPY --from=builder /app/identity/service ./

CMD ["./service"]