ddns
====

A self-hosted Dynamic DNS solution similar to DynDNS or NO-IP.

You can use a hosted version at [ddns.pboehm.org](http://ddns.pboehm.org/) where you can register a host under the `d.pboehm.de` domain (e.g `test.d.pboehm.de`).

## How can I update my IP if it changes?

`ddns` is built around a small webservice, so that you can update your IP address, simply by calling an URL periodically through `curl` (using `cron`). Hosts that haven't been updated for 10 days will be automatically removed. This can be configured in your own instance.

An API similar to DynDNS/NO-IP hasn't been implemented yet.

## Requirements for self-hosting

* A custom domain where the registrar allows NS-Records for subdomains
* A global accessible Server running an OS which is supported by the tools listed below
* A running [Redis](http://redis.io) instance for data storage
* An installation of [PowerDNS](https://www.powerdns.com/) with the Pipe-Backend included
* [Go](http://golang.org/) 1.3
