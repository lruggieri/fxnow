FROM golang:1.21.1-alpine3.18 as builder
WORKDIR /app
COPY fxupdate fxupdate
COPY common common
WORKDIR fxupdate
RUN go build -o service .

FROM alpine:latest
COPY --from=builder /app/fxupdate/service ./

CMD ["./service"]