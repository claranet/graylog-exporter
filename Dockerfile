# Dockerfile builds an image for a client_golang example.
#
# Use as (from the root for the client_golang repository):
#    docker build -f examples/$name/Dockerfile -t prometheus/golang-example-$name .

# Builder image, where we build the example.

FROM golang:1.9.0 AS builder

ENV GOPATH /go/src/graylog-exporter

WORKDIR /go/src/graylog-exporter
COPY . .
RUN echo "> GOPATH: " $GOPATH
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w'

# Final image.
FROM quay.io/prometheus/busybox:latest

LABEL maintainer "Martin Weber <martin.weber@de.clara.net>"

WORKDIR /
COPY --from=builder /go/src/graylog-exporter/graylog-exporter .
EXPOSE 9404
ENTRYPOINT ["/graylog-exporter"]
