ddns
====

A self-hosted Dynamic DNS solution similar to DynDNS or NO-IP.

You can use a hosted version at [ddns.pboehm.org](http://ddns.pboehm.org/) where you can register a host under the `d.pboehm.de` domain (e.g `test.d.pboehm.de`).

## How can I update my IP if it changes?

`ddns` is built around a small webservice, so that you can update your IP address, simply by calling an URL periodically through `curl` (using `cron`). Hosts that haven't been updated for 10 days will be automatically removed. This can be configured in your own instance.

An API similar to DynDNS/NO-IP hasn't been implemented yet.

## Self-Hosting

### Requirements

* A custom domain where the registrar allows NS-Records for subdomains
* A global accessible Server running an OS which is supported by the tools listed below
* A running [Redis](http://redis.io) instance for data storage
* An installation of [PowerDNS](https://www.powerdns.com/) with the Pipe-Backend included
* [Go](http://golang.org/) 1.3

### Installation

The following instructions are valid for Ubuntu/Debian. Some files/packages
could have other names/locations, please search for it.

You should have a working Go environment (correct `$GOPATH`).

    $ go version # check that you have go 1.3 installed
    go version go1.3 linux/amd64

Then install `ddns` via:

    $ go get github.com/pboehm/ddns
    $ ls $GOPATH/bin/ddns # the displayed path will be used later
    /home/user/gocode/bin/ddns

#### Backend

Install `pdns` and `redis-server`:

    $ sudo apt-get install redis-server pdns-server pdns-backend-pipe

Both services should start at boot automatically. You should open `udp/53` and
`tcp/53` on your Firewall so that `pdns` can be be used from outside of your
host.

    $ sudo vim /etc/powerdns/pdns.d/pipe.conf

`pipe.conf` should have the following content. Please adjust the path of `ddns`
and the values supplied to `--domain` and `--soa_fqdn`:

    launch=pipe
    pipebackend-abi-version=1
    pipe-command=/home/user/gocode/bin/ddns --soa_fqdn=dns.example.com --domain=sub.example.com backend

Then restart `pdns`:

    $ sudo service pdns restart

#### Frontend

`ddns` includes a webservice which is used for creating new hosts and updating
ip addresses. I prefer using `nginx` as a reverse proxy and not running `ddns`
on port 80. As a process manager, I prefer using `supervisord` so it is
described here.

    $ sudo apt-get install nginx supervisor

Create a supervisor config file for ddns:

    $ sudo vim /etc/supervisor/conf.d/ddns.conf
    $ cat /etc/supervisor/conf.d/ddns.conf
    [program:ddns]
    directory = /tmp/
    user = user
    command = /home/user/gocode/bin/ddns --domain=sub.example.com web
    autostart = True
    autorestart = True
    redirect_stderr = True

Restart the `supervisor` daemon and `ddns` listens on Port 8080 (can be
changed by adding `--listen=:1234`).

    $ sudo service supervisor restart

Now you have to add a nginx virtual host for `ddns`:

    $ sudo vim /etc/nginx/sites-enabled/default
    $ cat /etc/nginx/sites-enabled/default
    server {
        listen   80;
        server_name  ddns.example.com;

        location / {
            proxy_pass http://127.0.0.1:8080;
            proxy_set_header    Host            $host;
            proxy_set_header    X-Real-IP       $remote_addr;
            proxy_set_header    X-Forwarded-for $proxy_add_x_forwarded_for;
            proxy_connect_timeout 300;
        }
    }

Please adjust the `server_name` with a valid FQDN. Now we only have to restart
`nginx`:

    $ sudo service nginx restart
