FROM golang:1.15.0-alpine3.12 as builder
ARG opts="GOARCH=arm GOARM=7"

COPY . /app/
WORKDIR /app
RUN env ${opts} CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rpi-soldat .

FROM scratch
COPY --from=builder /app/rpi-soldat /rpi-soldat
COPY --from=builder /app/views /views
ENTRYPOINT [ "/rpi-soldat"]
