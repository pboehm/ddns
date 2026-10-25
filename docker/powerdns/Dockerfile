FROM buildpack-deps:trixie-scm

# the setup procedure according to https://repo.powerdns.com/ (Debian 9 Stretch)
RUN echo "deb [signed-by=/etc/apt/keyrings/auth-51-pub.asc] http://repo.powerdns.com/debian trixie-auth-51 main" > /etc/apt/sources.list.d/pdns.list \
	&& echo "Package: pdns-*\nPin: origin repo.powerdns.com\nPin-Priority: 600\n" >> /etc/apt/preferences.d/pdns \
	&& install -d /etc/apt/keyrings \
	&& curl https://repo.powerdns.com/FD380FBB-pub.asc | tee /etc/apt/keyrings/auth-51-pub.asc \
	&& apt-get -y update \
	&& apt-get install -y pdns-server pdns-backend-remote \
	&& rm -rf /var/lib/apt/lists/*

COPY pdns.conf /etc/powerdns/pdns.conf

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 53

CMD ["pdns_server", "--daemon=no"]