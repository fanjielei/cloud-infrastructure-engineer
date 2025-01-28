FROM golang:latest AS builder
WORKDIR /go/
COPY ./ ./
RUN go build -o /usr/local/bin/server-otel cmd/server/main.go

FROM debian:stable-slim
COPY --from=builder /usr/local/bin/server-otel /usr/local/bin/server-otel
EXPOSE 8080 9090
CMD [ "/usr/local/bin/server-otel" ]
