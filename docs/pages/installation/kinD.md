---
layout: default
title: KinD
permalink: installation/platforms/kind
type: installation
display-title: "false"
language: en
list: include
image: /assets/img/platforms/kind.png
---

{% include installation_prerequisites.html %}

**To set up and run Meshery on KinD** 

- [**Installation**](#installation)
  - [**Create cluster using KinD**](#create-cluster-using-kind)
      - [1. **KinD on WSL2**](#1-kind-on-wsl2)
      - [2. **KinD on other systems**](#2-kind-on-other-systems)
  - [**Access the KinD cluster**](#access-the-kind-cluster)
  - [**Using Helm**](#using-helm)
    - [**Helm v3**](#helm-v3)

### **Installation**

- On Mac / Linux via Homebrew (Recommended):

{% capture code_content %}brew install kind{% endcapture %}
{% include code.html code=code_content %}

- On macOS / Linux via curl:

{% capture code_content %}curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.8.1/kind-$(uname)-amd64
chmod +x ./kind
mv ./kind /some-dir-in-your-PATH/kind{% endcapture %}
{% include code.html code=code_content %}

If you are running Ubuntu on WSL2, use `Docker Ubuntu` distro to install `Docker`.

#### **Create cluster using KinD**

In order to successfully build the Meshery server on your local server, follow the instructions specific to your Operating System to complete the creation of a KinD cluster.

###### 1. **KinD on WSL2**

First, we will get the ip address of your WSL2 distro by:

{% capture code_content %}ip addr | grep eth0{% endcapture %}
{% include code.html code=code_content %}

You will see the output like:

{% capture code_content %}4: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    inet 172.1.1.1/20 brd 172.1.1.255 scope global eth0{% endcapture %}
{% include code.html code=code_content %}

Copy the ip address, we will use that in the next step.

Then, create a file called `kind_cluster.yaml` and put the ip address under `apiServerAddress`:

{% capture code_content %}kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  apiServerAddress: "172.1.1.1"{% endcapture %}
{% include code.html code=code_content %}

Now create the KinD cluster with the config file `kind_cluster.yaml`:

{% capture code_content %}kind create cluster --config kind_cluster.yaml --name kind --wait 300s{% endcapture %}
{% include code.html code=code_content %}

You will see

```bash
Creating cluster "kind" ...
 • Ensuring node image (kindest/node:v1.17.0) 🖼  ...
 ✓ Ensuring node image (kindest/node:v1.17.0) 🖼
 • Preparing nodes 📦   ...
 ✓ Preparing nodes 📦
 • Writing configuration 📜  ...
 ✓ Writing configuration 📜
 • Starting control-plane 🕹️  ...
 ✓ Starting control-plane 🕹️
 • Installing CNI 🔌  ...
 ✓ Installing CNI 🔌
 • Installing StorageClass 💾  ...
 ✓ Installing StorageClass 💾
 • Waiting ≤ 5m0s for control-plane = Ready ⏳  ...
 ✓ Waiting ≤ 5m0s for control-plane = Ready ⏳
 • Ready after 59s 💚
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Not sure what to do next? 😅 Check out https://kind.sigs.k8s.io/docs/user/quick-start/
```

###### 2. **KinD on other systems**

Creating a Kubernetes cluster is as simple as `kind create cluster`.

For more configuration of installation, please refer to KinD official documentation.

#### **Access the KinD cluster**

By default, the cluster access configuration is stored in ${HOME}/.kube/config if $KUBECONFIG environment variable is not set. You can set the `KUBECONFIG` environment with the command below:

{% capture code_content %}export KUBECONFIG=${HOME}/.kube/config{% endcapture %}
{% include code.html code=code_content %}
Use the command below check the connection of the cluster and make sure the cluster you connected what's the cluster was created by KinD:

{% capture code_content %}kubectl cluster-info --context kind-kind{% endcapture %}
{% include code.html code=code_content %}

To delete your cluster use:

{% capture code_content %}kind delete cluster --name kind{% endcapture %}
{% include code.html code=code_content %}

#### **Using Helm**

##### **Helm v3**

We strongly recommend to use Helm v3, because of this version not including the Tiller(https://helm.sh/blog/helm-3-preview-pt2/#helm) component anymore. It’s lighter and safer.

Run the following:

{% capture code_content %} $ helm repo add meshery https://meshery.io/charts/
 $ helm install meshery meshery/meshery -n meshery --create-namespace{% endcapture %}
{% include code.html code=code_content %}
 - Meshery server supports customizing authentication flow callback URL, which can be configured in the following way
{% capture code_content %}$ helm install meshery --namespace meshery --set env.MESHERY_SERVER_CALLBACK_URL=https://custom-host meshery/meshery{% endcapture %}
{% include code.html code=code_content %}

- **NodePort** - If your cluster does not have an Ingress Controller or a load balancer, then use NodePort to expose Meshery and that can be modify under the chart `values.yaml`:

{% capture code_content %}service:
   type: NodePort
   port: 9081
   annotations: {}{% endcapture %}
{% include code.html code=code_content %}

Meshery should now be connected with your managed Kubernetes instance. Take a look at the [Meshery guides]({{ site.baseurl }}/guides) for advanced usage tips.
