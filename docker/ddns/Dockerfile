FROM golang:alpine

RUN apk add --no-cache git

COPY docker/ddns/entrypoint.sh /entrypoint.sh
WORKDIR /go/src/github.com/pboehm/ddns
COPY . .

RUN GO111MODULE=on go get -d -v ./...
RUN GO111MODULE=on go install -v ./...

ENV GIN_MODE release
ENV DDNS_EXPIRATION_DAYS 10

ENTRYPOINT ["/entrypoint.sh"]
