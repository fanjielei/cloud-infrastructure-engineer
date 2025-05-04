FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /build
COPY ./ ./

RUN CGO_ENABLED=0 go build -o ./bin/client cmd/client/main.go

FROM alpine:3.21.3
# -H nohome -D nopass
RUN addgroup -S group && \
    adduser -S user -G group -H -D

COPY --from=builder --chown=user:group /build/bin/client /client
USER user

ENV HOST="server"
ENV PORT="8080"

CMD [ "/client" ]
