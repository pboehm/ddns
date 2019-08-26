FROM golang:alpine

RUN apk add --no-cache git

WORKDIR /go/src/github.com/pboehm/ddns
COPY . .

RUN GO111MODULE=on go get -d -v ./...
RUN GO111MODULE=on go install -v ./...

ENV GIN_MODE release

CMD /go/bin/ddns --domain=${DDNS_DOMAIN} --soa_fqdn=${DDNS_SOA_DOMAIN} --redis=${DDNS_REDIS_HOST}
