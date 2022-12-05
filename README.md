
Prometheus zerto-exporter
==============================

Zerto is a powerfull software for desaster&recovery and long distance migration of virtual machines between multiple virtualization technologies. With Zerto you can failover and move VM's between VMWare vCenter, VMWare vCloud, HyperV, AWS and Azure.

Find out more about Zerto https://www.zerto.com/

This exporter can be used to scratch metrics from zerto api and make them available for prometheus monitoring.

Tested & Supported
=============

Currently we just tested Zerto for vSphere!

Tested on __Windows ZVM__ and __Linux ZVM__ with ZVM Version 9.7

Usage
=======

Usage is really simple. Just need a user to login to Zerto ZVM and the url.

```
Usage of /zerto-exporter:
  -listen.address string
        The address to lisiten on for HTTP requests. (default ":9403")
  -log.level string
        Log-Level (debug, warn, error) (default "info")
  -version
        Prints current version
  -zerto.password string
        Zerto API User Password
  -zerto.url string
        Zerto URL to connect https://zvm.local.host:9669
  -zerto.username string
        Zerto API User
```

HowTo Run
---------------

```
docker run --publish 9403:9403 claranet/zerto-exporter \
		-zerto.url "https://zvm.company.local:9669" \
		-zerto.username "prometheus@vsphere.local" \
		-zerto.password pASSwOrd
```

External Links
==================

* https://www.zerto.com
* https://prometheus.io
* [Zerto REST-API Documentation](https://help.zerto.com/bundle/API.ZVR.HTML)
* [Grafana Dashboard](https://grafana.com/dashboards/3765)

ToDo's
=======

* Add metrics for errors
