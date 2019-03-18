FROM golang:1.12-alpine as builder

RUN apk update && apk add --no-cache git libc-dev gcc

COPY . /go/src/github.com/Nuance-Mobility/azure-request-limitometer

# Install dep
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/Nuance-Mobility/azure-request-limitometer
RUN dep ensure

WORKDIR /go
RUN go build github.com/Nuance-Mobility/azure-request-limitometer/cmd/limitometer

FROM alpine

RUN apk update && \
    apk add --no-cache \
      bash \
      ca-certificates

WORKDIR /root

COPY --from=builder /go/limitometer limitometer

ENTRYPOINT ["./limitometer"]
