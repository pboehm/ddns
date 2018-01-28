FROM buildpack-deps:jessie-scm

# the setup procedure according to https://repo.powerdns.com/ (Debian 8 Jessie)
RUN echo "deb http://repo.powerdns.com/debian jessie-auth-41 main" > /etc/apt/sources.list.d/pdns.list \
	&& echo "Package: pdns-*\nPin: origin repo.powerdns.com\nPin-Priority: 600\n" >> /etc/apt/preferences.d/pdns \
	&& curl https://repo.powerdns.com/FD380FBB-pub.asc | apt-key add - \
	&& apt-get -y update \
	&& apt-get install -y pdns-server pdns-backend-remote \
	&& rm -rf /var/lib/apt/lists/*

COPY pdns.conf /etc/powerdns/pdns.conf

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 53

CMD ["pdns_server", "--daemon=no"]