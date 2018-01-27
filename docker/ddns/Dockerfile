FROM golang:1.9

WORKDIR /go/src/github.com/pboehm/ddns
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ."

ENV GIN_MODE release

CMD /go/bin/ddns --domain=${DDNS_DOMAIN} --soa_fqdn=${DDNS_SOA_DOMAIN} --redis=${DDNS_REDIS_HOST}
