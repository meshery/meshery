
# Meshery Playground Bare Metal Configuration 

### DNS
playground.meshery.io - 147.28.141.9

### Host: meshery02

#### Static IP address configuration
File: `/etc/netplan/00-elastic.yaml`
```
network:
  version: 2
  renderer: networkd
  ethernets:
    lo:
      addresses:
        - 127.0.0.1/8
        - 147.28.141.9/32
```

Affect changes:
`netplan apply`

Needed?
`systemctl restart networking`
