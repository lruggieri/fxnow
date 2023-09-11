FROM golang:1.21.1-alpine3.18 as builder
WORKDIR /app
COPY fxrate fxrate
COPY common common
WORKDIR fxrate
RUN go build -o service .

FROM alpine:latest
COPY --from=builder /app/fxrate/service ./

CMD ["./service"]