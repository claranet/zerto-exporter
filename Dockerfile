# Dockerfile builds an image for a client_golang example.
#
# Builder image, where we build the example.

FROM golang:1.15.6 AS builder

ENV GOPATH /go/src/zerto-exporter

WORKDIR /go/src/zerto-exporter
COPY . .
RUN echo "> GOPATH: " $GOPATH
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w'

# Final image.
FROM quay.io/prometheus/busybox:latest

LABEL maintainer "Martin Weber <martin.weber@de.clara.net>"
LABEL version "0.1.3"

WORKDIR /
COPY --from=builder /go/src/zerto-exporter/zerto-exporter /usr/local/bin/zerto-exporter
EXPOSE 9403
ENTRYPOINT ["/usr/local/bin/zerto-exporter"]
