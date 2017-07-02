FROM golang:1.7

WORKDIR /go/src/github.com/pboehm/ddns
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ."

CMD ["/bin/true"]
