#!/bin/sh

if [[ $# -eq 0 ]]; then

    if [[ "${DEBUG_MODE:-false}" == "true" ]]; then
        /go/bin/ddns \
            --domain=${DDNS_DOMAIN} \
            --soa_fqdn=${DDNS_SOA_DOMAIN} \
            --redis=${DDNS_REDIS_HOST} \
            --expiration-days=${DDNS_EXPIRATION_DAYS} \
            --verbose
    else
	 /go/bin/ddns \
            --domain=${DDNS_DOMAIN} \
            --soa_fqdn=${DDNS_SOA_DOMAIN} \
            --redis=${DDNS_REDIS_HOST} \
            --expiration-days=${DDNS_EXPIRATION_DAYS}
    fi

else
    "$@"
fi
