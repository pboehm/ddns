FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /go/src/github.com/pboehm/ddns
COPY . .

RUN GO111MODULE=on go get -d -v ./...
RUN export CGO_ENABLED=0 && GO111MODULE=on go install -v ./...

ENV GIN_MODE release

FROM scratch

COPY --from=builder /go/bin/ddns /go/bin/ddns

ENTRYPOINT ["/go/bin/ddns"]
