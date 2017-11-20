
Prometheus zerto-exporter
==============================

Zerto is a powerfull software for desaster&recovery and long distance migration of virtual machines between multiple virtualization technologies. With Zerto you can failover and move VM's between VMWare vCenter, VMWare vCloud, HyperV, AWS and Azure.

Find out more about Zerto https://www.zerto.com/

This exporter can be used to scratch metrics from zerto api and make them available for prometheus monitoring.

Tested & Supported
=============

Currently we just tested Zerto for vSphere!

Usage
=======

Usage is really simple. Just need a user to login to Zerto ZVM and the url.

```
Usage of /zerto-exporter:
  -listen-address string
    	The address to lisiten on for HTTP requests. (default ":9403")
  -log.level value
    	Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal, panic].
  -zerto.password string
    	Zerto API User Password
  -zerto.session-age int
    	Zerto Session recreation Time (default 3600)
  -zerto.url string
    	Zerto URL to connect https://zvm.local.host:9669
  -zerto.username string
    	Zerto API User
```

Build with Docker
===============================

To build binariesd using docker you must have docker in version 17.06 or greater to support multiple build in one Dockerfile. To build the binaries just call `make`or to it manually with

```
$ docker build -t zerto-exporter .
```

HowTo Run
---------------

```
docker run --publish 9403:9403 zerto-exporter \
		-zerto.url "https://zvm.company.local:9669" \
		-zerto.username "prometheus@vsphere.local" \
		-zerto.password pASSwOrd
```

External Links
==================

* https://www.zerto.com
* https://prometheus.io
* [Zerto REST-API Documentation](http://s3.amazonaws.com/zertodownload_docs/Latest/Zerto%20Virtual%20Replication%20Zerto%20Virtual%20Manager%20%28ZVM%29%20-%20vSphere%20Online%20Help/index.html#page/RestfulAPIs/APIsIntro.2.1.html)
* [Grafana Dashboard](https://grafana.com/dashboards/3765)
