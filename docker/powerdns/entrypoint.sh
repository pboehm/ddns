#!/usr/bin/env bash

CONFIG_FILE=/etc/powerdns/pdns.conf

sed -i 's/{{PDNS_REMOTE_HTTP_HOST}}/'"${PDNS_REMOTE_HTTP_HOST}"'/g' ${CONFIG_FILE}

exec "$@"