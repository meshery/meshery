
# Meshery Playground Bare Metal Configuration 

![Deployment Topology](../docker/meshery-playground-deployment.png)


### DNS
playground.meshery.io - 147.28.141.9

### Hosts
- c3-medium-x86-01-meshery - docker host with Meshery Server
- c3-medium-x86-02-meshery - single node k8s cluster

#### Access
```
ssh -i ~/.ssh/equinix-metal root@c3-medium-x86-01-meshery
ssh -i ~/.ssh/equinix-metal root@c3-medium-x86-02-meshery
```

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

### Firewalling

cd /etc/network/interface/if-up.d/
cat 
```
#!/bin/sh

iptables -I INPUT -s 192.210.143.199 -j DROP
```
> protect-kubelet
chmod +x protect-kubelet

## Reinstalling Meshery after a clean Meshery uninstall
If mistakenly or for some reason Meshery is uninstalled (along with `meshery` namespace) then follow the below steps to bring it back up.
1. `helm install meshery meshery/meshery --namespace meshery --set env.PROVIDER=Meshery`
2. `kubectl create secret tls -n meshery tls-secret-meshery --cert=/etc/letsencrypt/live/playground.meshery.io-0001/fullchain.pem --key=/etc/letsencrypt/live/playground.meshery.io-0001/privkey.pem`  Make sure to have private and public keys generated for `playground.meshery.io` in appropriate directories.
3. kubectl apply -f contour-http-proxy.yaml 
## Prometheus deployment
For monitoring the playground deployment we have a node exporter running on both the nodes.
Node exporter can be installed by following the below commands after doing an SSH on the node.
```
wget https://github.com/prometheus/node_exporter/releases/download/v1.5.0/node_exporter-1.5.0.linux-amd64.tar.gz
tar xvfz  node_exporter-1.5.0.linux-amd64.tar.gz
cd node_exporter-1.5.0.linux-amd64/
./node_exporter &
```
Make sure to add nodeexporter.service inside /etc/systemd/system so that node exporter automatically starts when server boots up.

After this create the namespace `monitoring` and apply `prometheus.yaml`.
Notice the job_name: `prometheus` configured with two targets. It is the address of node exporter running on each node.

## Grafana deployment
Grafana uses prometheus as data source and expose it 3000. Prometheus and Node exporter should be deployed a long as with Grafana.

Create the namespace `monitoring` if not already done and apply `grafana.yaml`.
Notice that url in the data source should match with prometheus

### Nginx Configuration for websocket (until annotations issue is solved)
## When to do this: Whenever Ingress is reset directly or indirectly(via scaling deployments).
Nginx ingress controller by default (sometimes) doesn't pick up the custom configuration passed to it via annotations so after the nginx pod starts below steps are required to make sure that websocket requests go through.

- exec into the nginx pod and go to /etc/nginx/conf.d.
- There will be a .conf file present there. Take that file out in your local system to edit since the container doesn't have vi or any code editor.
- Add the below location blocks just above the location block for `/`

```
	location /api/system/graphql/query {
	   set $service "meshery"; 
	   proxy_set_header Upgrade $http_upgrade;
	   proxy_http_version 1.1;
	   proxy_set_header X-Forwarded-Host $http_host;
	   proxy_set_header X-Forwarded-Proto $scheme;
	   proxy_set_header X-Forwarded-For $remote_addr;
	   proxy_set_header Host $host;
	   proxy_set_header Connection "upgrade";
	   proxy_cache_bypass $http_upgrade;
	   proxy_pass http://meshery-meshery-playground.meshery.io-meshery-9082;
	 }
	location /api/provider/extension/server/graphql/query {
	   set $service "meshery"; 
	   proxy_set_header Upgrade $http_upgrade;
	   proxy_http_version 1.1;
	   proxy_set_header X-Forwarded-Host $http_host;
	   proxy_set_header X-Forwarded-Proto $scheme;
	   proxy_set_header X-Forwarded-For $remote_addr;
	   proxy_set_header Host $host;
	   proxy_set_header Connection "upgrade";
	   proxy_cache_bypass $http_upgrade;
	   proxy_pass http://meshery-meshery-playground.meshery.io-meshery-9082;
	 }	
```

- In the above configuration, make sure the proxy_pass directive is the same as the one present for the already existing `/` location block. If not, then replace it with the one already present

- Enter `cat > <name of the file>` to get the stdin and then paste the edit file on the console. Exit using cmd+c or cmd+d

- Perform an nginx reload with `nginx -s reload`

- Confirm that the websocket requests are passing through.



### Renewing SSL certificate for playground. (When cert manager is not used and certificates are generated manually)
NOTE: Make sure to renew certificates before they expire

- Run sudo certbot certonly --manual --preferred-challenges http -d playground.meshery.io and press 2 to renew the certificates. 
- New certs will be stored at /etc/letsencrypt/live/playground.meshery.io/fullchain.pem and /etc/letsencrypt/live/playground.meshery.io/privkey.pem
- Delete the previous secrets named tls-secret or tls-meshery-secret(check the name in the Ingress resource) and delete the the secret with `nginx-ingress-default-server-tls` in the name in meshery namespace
- Run `kubectl create secret tls tls-secret --cert=/etc/letsencrypt/live/playground.meshery.io/fullchain.pem --key=/etc/letsencrypt/live/playground.meshery.io/privkey.pem`  and `kubectl create secret tls nginx-ingress-default-server-tls --cert=/etc/letsencrypt/live/playground.meshery.io/fullchain.pem --key=/etc/letsencrypt/live/playground.meshery.io/privkey.pem`
- Delete the ingress pod so that it can restart and use the newly configured secrets.
