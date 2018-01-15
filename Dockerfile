##
FROM golang:alpine as builder

RUN apk update && apk add --no-cache git libc-dev gcc

WORKDIR /go

COPY common /go/src/azure-request-limitometer/common
COPY querier /go/src/azure-request-limitometer/querier

RUN go get azure-request-limitometer/querier && \
    go build -a azure-request-limitometer/querier

##
FROM alpine

RUN apk update && \
    apk add --no-cache \
      bash \
      curl

WORKDIR /root

COPY pusher/pusher.sh pusher.sh
COPY entrypoint.sh entrypoint.sh
COPY --from=builder /go/querier querier

RUN chmod +x pusher.sh && \
    chmod +x entrypoint.sh && \
    mkdir -p /etc/kubernetes

# expects /etc/kubernetes/azure.json to be mounted
# expects node hostname to be passed as argument ($1)
ENTRYPOINT ["/root/entrypoint.sh"]
CMD ["kn-edge-0"] # default arg
