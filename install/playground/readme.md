
# Meshery Playground Bare Metal Configuration 

_Deployment Topology - see Meshery Architecture deck_

## Production (CNCF) Deployment

The live playground (`playground.meshery.io`) runs as the `playground` Deployment in the
`meshery` namespace on the CNCF cluster (kube-context `prod`). Its canonical spec lives at
[`kubernetes/playground-deployment.yaml`](kubernetes/playground-deployment.yaml) - apply it
with `kubectl -n meshery apply -f kubernetes/playground-deployment.yaml`.

Provider configuration is set there: the playground offers only the **Meshery** provider
(`PROVIDER_BASE_URLS=https://cloud.meshery.io`) and pre-selects it on load
(`PROVIDER=Meshery`). Releases reconcile the spec from this manifest and then roll the
container image forward via `.github/workflows/cncf-playground-deploy-meshery.yaml`
(`kubectl apply` then `kubectl set image`), so env changes ride in with the next release.
Out-of-band edits on the cluster will be overwritten on the next deploy - put the change
in this manifest instead.

### Sandbox cluster: Meshery Operator, MeshSync and Broker

The playground Pod runs on the `prod` cluster, but the kubeconfig it mounts from the
`playground` ConfigMap (at `/home/appuser/.kube`) names a **different** cluster - the
sandbox cluster (kube-context `staging`). Every cluster the playground shows its users,
and every operator status it reports, is the sandbox cluster's. The `meshery` namespace
on `prod` also carries a Meshery Operator release, but no playground connection uses it,
so check the **sandbox** cluster when the playground reports an unhealthy operator.

With `MESHSYNC_DEFAULT_DEPLOYMENT_MODE=operator`, Meshery Server installs the
`meshery-operator` chart into the sandbox cluster's `meshery` namespace (release
`meshery-operator`) and subscribes to its Broker. Two things the chart does not provide
today must be in place for discovery data to reach the playground:

| Requirement | Why | Manifest (apply to the sandbox cluster) |
| --- | --- | --- |
| Broker reachable from the `prod` cluster | The Server runs out of the sandbox cluster, and current releases skip the managed kube-API port-forward whenever the Server itself runs in any Kubernetes cluster, so a `ClusterIP` Broker is unreachable | [`sandbox-cluster/broker-service-exposure.yaml`](kubernetes/sandbox-cluster/broker-service-exposure.yaml) - `LoadBalancer`, restricted to the `prod` cluster's NAT egress address; NATS still requires its token, but the traffic is not TLS-encrypted: the range and the token limit who can connect, not what a network observer can read. The LoadBalancer also exposes the unauthenticated NATS monitoring port 8222 (`/connz`); the Broker resource offers no port selection, so keeping it off the LoadBalancer needs an operator-side change. Prefer the managed port-forward where it applies; expose the Broker only over a private path |
| MeshSync can read the cluster | MeshSync runs as ServiceAccount `meshery/meshery-operator`, which the v1.0.x chart scopes to the operator's own needs, so every watch fails with `forbidden` | [`sandbox-cluster/meshsync-discovery-interim-rbac.yaml`](kubernetes/sandbox-cluster/meshsync-discovery-interim-rbac.yaml) - read-only, interim: cluster-wide `get`/`list`/`watch`, including Secrets and their credentials, and `get` on subresources such as `nodes/proxy`, which Kubernetes documents as usable to execute commands in Pods through the kubelet API. Remove it once the chart gives MeshSync its own ServiceAccount |

Apply both from `install/playground` with
[`kubernetes/sandbox-cluster/apply.sh <sandbox-context> <egress-ip>`](kubernetes/sandbox-cluster/apply.sh).
It validates the address, applies with server-side apply under the
`meshery-playground-ops` field manager (which Meshery Server's own applies, field manager
`meshery-server`, never withdraw), restarts MeshSync and waits for the Broker's external
endpoint. The egress address is the public IP of the `prod` cluster VCN's NAT gateway
(`oci network nat-gateway list ... --query 'data[].{name:"display-name",ip:"nat-ip"}'`).

> [!WARNING]
> Do not switch a sandbox connection's MeshSync deployment mode, disconnect it, or delete
> it from the playground. In current releases each of those uninstalls the shared
> `meshery-operator` release for **every** connection to the sandbox cluster, and the
> uninstall leaves the Broker and MeshSync custom resources stuck `Terminating` (their
> finalizers need the operator that was just removed).

#### Verifying

```bash
kubectl --context staging -n meshery get deploy,sts,pods            # operator, meshery-nats, meshery-meshsync Ready
kubectl --context staging -n meshery get broker meshery-broker -o jsonpath='{.status.endpoint}'  # external endpoint present
kubectl --context staging -n meshery logs deploy/meshery-meshsync | grep -c forbidden            # 0
kubectl --context staging -n meshery exec meshery-nats-0 -c nats -- \
  wget -qO- 'http://127.0.0.1:8222/connz' | jq '[.connections[] | {ip, name}]'  # a client from the prod NAT address
kubectl --context prod -n meshery logs deploy/playground -c meshery | grep 'Connected to Meshery Broker'
```

#### Recovering the operator on the sandbox cluster

1. **Remove a stray operator only after confirming it is unmanaged.** Meshery manages only `meshery/meshery-operator`.
   A `meshery-operator` Deployment in any other namespace can be a manual install or a Helm release in another namespace.
   Run `helm --kube-context staging list -A` and confirm that no release owns the Deployment. Confirm that its objects carry no `meta.helm.sh/release-name` annotation and that no other installation uses it. Then delete it together with its ServiceAccount, Service, Role and RoleBinding.
2. **Clear stuck custom resources.** If `broker meshery-broker` or `meshsync meshery-meshsync`
   has a `deletionTimestamp` and no operator is running, remove their finalizers. Their
   owned objects (NATS StatefulSet, Services, ConfigMap, Secret, MeshSync Deployment) carry
   ownerReferences, so garbage collection removes them exactly as the finalizer would.
   Reinstalling first instead makes Helm adopt the terminating resources, and the new
   operator then finalizes them away, leaving an operator with no Broker.
3. **Reinstall the chart Meshery Server expects** (the version matching the playground
   release, from `https://meshery.github.io/meshery.io/charts`), using Helm v3 - the major
   version Meshery Server embeds:
   `helm install meshery-operator <chart> -n meshery --kube-context staging --set fullnameOverride=meshery-operator --wait`.
   Do not pass `--atomic`: a rollback-uninstall recreates the stuck state in step 2.
4. **Re-apply the sandbox settings** with `kubernetes/sandbox-cluster/apply.sh <sandbox-context> <egress-ip>`.
5. **Restart the playground Deployment** (`kubectl --context prod -n meshery rollout restart
   deploy/playground`). Its Broker connections keep the endpoint and token they resolved at
   connect time, so a recreated Broker is not picked up until they reconnect.

_The bare-metal/Equinix notes below describe an earlier topology and are retained for reference._

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


### Troubleshooting Steps:
1. Issues due to Cilium CNI:
	The networking is setup using WeaveNet, installation of Cilium updates the CNI conf causing issues.
	If the pods are not getting created/provisioned on the node, or there are any networking issues happening, most often installation of Cilium is the culprit.
	<br>
	First check if there any Cilium pods/deamonsets in the `kube-system` namespace.

2. Delete all the resources related to Cilium.
3. Run `cilium uninstall`.
4. Check if the issue is resolved, else proceed to next steps.
5. Run `cd /etc/cni/net.d/`, check for `05-cilium.conf` or any other conf files of pattern `*-cilium.conf`, remove those files.
6. You may also delete any cilium related network interfaces. Run `ifconfig interface_name down`, where interface_name is  all interfaces whose name starts with `cilium_*`.

