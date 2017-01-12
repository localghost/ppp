# Packer Pre-Processor

## Use case

Filename `ubuntu-common.json`
```
{
    "iso_url": "http://releases.ubuntu.com/16.04/ubuntu-16.04.1-server-amd64.iso",
    "iso_checksum": "d2d939ca0e65816790375f6826e4032f",
    "iso_checksum_type": "md5",
    "boot_command": [
        "<enter><wait>",
        "<f6><esc>",
        "<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
        "<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
        "<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
        "<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
        "<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
        "<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
        "<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
        "<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
        "<bs><bs><bs>",
        "/install/vmlinuz ",
        "initrd=/install/initrd.gz ",
        "net.ifnames=0 ",
        "auto-install/enable=true ",
        "debconf/priority=critical ",
        "preseed/url=http://{{.HTTPIP}}:{{.HTTPPort}}/ubuntu-16.04/preseed.cfg ",
        "<enter>"
    ]
}
```

Filename `ubuntu-virtualbox.json`
```
{
    "builders": [{
        "type": "virtualbox-iso",
        "ppp-inline": "ubuntu-common.json"
        ...
    }]
}
```

Filename `ubuntu-vmware.json`
```
{
    "builders": [{
        "type": "vmware-iso",
        "ppp-inline": "ubuntu-common.json"
        ...
    }]
}
```

Command:
```
ppp ubuntu-virtualbox.json | packer build -
ppp ubuntu-vmware.json | packer build -
```

## Status

[![Build Status](https://travis-ci.org/localghost/ppp.svg?branch=master)](https://travis-ci.org/localghost/ppp)
