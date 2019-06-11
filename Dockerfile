FROM golang:1.12 AS builder

ENV GOPATH /go/src/graylog-exporter

WORKDIR /go/src/graylog-exporter
COPY . .
RUN echo "> GOPATH: " $GOPATH
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w'

# Final image.
FROM quay.io/prometheus/busybox:latest

LABEL maintainer "Martin Weber <martin.weber@de.clara.net>"
LABEL version "0.2.0"

COPY --from=builder /go/src/graylog-exporter/graylog-exporter /usr/local/bin/graylog-exporter

ENTRYPOINT ["/usr/local/bin/graylog-exporter"]
EXPOSE 9404