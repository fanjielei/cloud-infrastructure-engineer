FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /build
COPY ./ ./

RUN CGO_ENABLED=0 go build -o /usr/local/bin/server-otel cmd/server/main.go

FROM alpine:3.21.3
# -H nohome -D nopass
RUN addgroup -S group && \
    adduser -S user -G group -H -D

COPY --from=builder --chown=user:group /usr/local/bin/server-otel /usr/local/bin/server-otel
USER user
EXPOSE 8080 8090

CMD [ "/usr/local/bin/server-otel" ]
