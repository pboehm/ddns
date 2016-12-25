#!/usr/bin/env bash

CONFIG_FILE=/etc/powerdns/pdns.conf

sed -i 's/{{DDNS_DOMAIN}}/'"${DDNS_DOMAIN}"'/g' ${CONFIG_FILE}
sed -i 's/{{DDNS_SOA_DOMAIN}}/'"${DDNS_SOA_DOMAIN}"'/g' ${CONFIG_FILE}
sed -i 's/{{DDNS_VERBOSE}}/'"${DDNS_VERBOSE}"'/g' ${CONFIG_FILE}

exec "$@"