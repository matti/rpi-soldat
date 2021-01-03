FROM golang:1.15.0-alpine3.12 as builder
ARG opts="GOARCH=arm GOARM=7"

COPY . /build/
WORKDIR /build
RUN env ${opts} CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rpi-soldat .

FROM alpine:3.12

RUN apk add --no-cache \
  bash

COPY --from=builder /build/rpi-soldat /usr/bin/rpi-soldat

WORKDIR /app
COPY --from=builder /build/views /app/views
COPY --from=builder /build/public /app/public

ENTRYPOINT [ "/usr/bin/rpi-soldat"]
