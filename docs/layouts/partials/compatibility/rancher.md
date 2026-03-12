## Install Meshery in Rancher Kubernetes cluster

Once Docker is installed, the next step will be to install a Kubernetes cluster.
Under "Settings" in the Docker Desktop application, enable *Kubernetes*.
In this how-to, [K3d](https://github.com/rancher/k3d) will be used as it relies only on Docker.

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
k3d cluster create
export $env:KUBECONFIG = "$(k3d.exe kubeconfig get 'k3s-default')"</div></div>
  </pre>

If using Scoop, run the following in the PowerShell to install a Kubernetes cluster :

  <pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">scoop install k3d
k3d cluster create
export KUBECONFIG="$(k3d kubeconfig get 'k3s-default')"</div></div>
  </pre>
