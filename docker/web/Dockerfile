FROM buildpack-deps:jessie-scm
MAINTAINER Philipp Böhm <philipp@pboehm.org>

ENV GIN_MODE release

CMD /go/bin/app --domain=${DDNS_DOMAIN} --redis=redis:6379 web
