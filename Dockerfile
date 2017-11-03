# Dockerfile builds an image for a client_golang example.
#
# Builder image, where we build the example.

FROM golang:1.9.0 AS builder

ENV GOPATH /go/src/zerto-exporter

WORKDIR /go/src/zerto-exporter
COPY . .
RUN echo "> GOPATH: " $GOPATH
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w'

# Final image.
FROM quay.io/prometheus/busybox:latest

LABEL maintainer "Martin Weber <martin.weber@de.clara.net>"

WORKDIR /
COPY --from=builder /go/src/zerto-exporter/zerto-exporter .
EXPOSE 9403
ENTRYPOINT ["/zerto-exporter"]
