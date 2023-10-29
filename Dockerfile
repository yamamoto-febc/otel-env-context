FROM golang:1.21 AS builder

ADD . /go/src/github.com/yamamoto-febc/otel-env-context
WORKDIR /go/src/github.com/yamamoto-febc/otel-env-context
ENV CGO_ENABLED 0
RUN go build -o otel-parent cmd/otel-parent/*
RUN go build -o otel-child cmd/otel-child/*
# ======

FROM alpine:3.16

COPY --from=builder /go/src/github.com/yamamoto-febc/otel-env-context/otel-parent /usr/bin/
COPY --from=builder /go/src/github.com/yamamoto-febc/otel-env-context/otel-child /usr/bin/
ENTRYPOINT ["/usr/bin/otel-parent"]