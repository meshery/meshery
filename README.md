<img align="right" src="https://layer5.io/assets/images/cube-sh-small.png" />

# Meshery

A service mesh playground to faciliate learning about functionality and performance of different service meshes. Meshery incorporates the collection and display of metrics from applications running in the playground.

- [Website](https://layer5.io/meshery)
- [Performance benchmark design document](https://docs.google.com/document/d/1nV8TunLmVC8j5cBELT42YfEXYmhG3ZqFtHxeG3-w9t0/edit?usp=sharing)
- [Architecture](https://docs.google.com/presentation/d/1UbuYMpn-e-mWVYwEASy4dzyZlrSgZX6MUfNtokraT9o/edit?usp=sharing)

![Service Mesh Playground](/public/static/img/meshery.png?raw=true "Service Mesh Playground")

## Functionality
1. Multi-mesh Performannce Benchmark
Meshery is intended to be a vendor and project-neutral utility for uniformly benchmarking the performance of service meshes. Between service mesh and proxy projects, a number of different tools *and results* exist. For example, Istio's [Performance and Scalability WG](https://github.com/istio/community/blob/master/WORKING-GROUPS.md#performance-and-scalability) currently uses a couple of different tools to measure Istio performance: [BluePerf](https://ibmcloud-perf.istio.io/regpatrol/) and [Fortio](https://fortio.istio.io).

1. Multi-mesh Functionalty Playground

### Running Meshery
## Service Mesh Playground
Sample applications will be included in Meshery. 

#### General Prerequisites
1. Docker engine (e.g. Docker for Desktop).
1. Kubernetes cluster (preferably version 1.10+).

#### Istio Playground Prerequisites
1. Istio version 1.0.3+ in `istio-system` namespace along with the Istio ingress gateway.


#### Run Meshery
To run Meshery:
- On Kubernetes
  - You can deploy Meshery to an existing kubernetes cluster using the provided yaml file into any namespace of your choice. For now let us deploy it to a namespace `meshery`: 

    ```
    kubectl create ns meshery
    kubectl -n meshery apply -f deployment_yamls/k8s

    # additional step for running on Istio
    kubectl -n meshery apply -f deployment_yamls/istio.yaml
    ```
    If you want to use a different namespace, please change the name of the namespace in the `ClusterRoleBinding` section appropriately.
  - Meshery can be deployed either on/off the mesh.
  - If deployed on the same Kubernetes cluster as the mesh, you dont have to provide a kubeconfig file.
  - please review the yaml and make necessary changes as needed for your cluster
- On Docker
  - We have a docker-compose.yaml file which can be used to spin up the services quickly
  - There a few requirements for running all the Meshery services on your local
    - SSO, which uses Twitter and/or Github
      - Instructions to setup Twitter for SSO can be found <a href="#twitter">here</a>
      - Instructions to setup Github for SSO can be found <a href="#github">here</a>
    - Add an entry for `meshery-saas` in your `/etc/hosts` file to point to 127.0.0.1 and save the file
    - After setting up SSO, store the respective key and secret as variables in the shell as shown below.
      - for Twitter:
      ```
      TWITTERKEY="PASTE TWITTER KEY"
      ```
      ```
      TWITTERSECRET="PASTE TWITTER SECRET"
      ```
      - for Github:
      ```
      GITHUBKEY="PASTE GITHUB KEY"
      ```
      ```
      GITHUBSECRET="PASTE GITHUB SECRET"
      ```
      __Note__: you can use Twitter and/or Github

      Now that the environment variables are setup, we can start the containers by running:
      ```
      docker-compose up
      ```
      Please add a `-d` flag to the above command if you want to run it in the background.
  - Now you should be able to access Meshery in your browser at `http://localhost:8080/play`

##### <a name="twitter">Using Twitter for SSO</a>
- Create an app in the Twitter developer console: [https://developer.twitter.com/en/apps](https://developer.twitter.com/en/apps) after logging in.
- Fill appropriate details in the presented form
  - Remember to enable to `Sign in with Twitter`
  - For the callback url, please use this value: `http://meshery-saas:9876/auth/twitter/callback`
- After creating the app you will be able to grab the API key and secret from the `Keys and tokens` section of the app.

##### <a name="github">Using Github for SSO</a>
- Create an OAuth app in the Github developer settings: [https://github.com/settings/developers](https://github.com/settings/developers) after logging in.
- Fill appropriate details in the presented form
  - For the callback url, please use this value: `http://meshery-saas:9876/auth/github/callback`
- After creating the app you will be able to grab the Client ID and Secret from the app page.


## Linkerd Playground App
_coming soon for Linkerd_
### Running Meshery
_coming soon for Linkerd_
#### Linkerd Prerequisites
_coming soon for Linkerd_
#### Run Linkerd Playground
_coming soon for Linkerd_

## Contributing
Please do! Contributions, updates, [discrepancy reports](/../../issues) and [pull requests](/../../pulls) are welcome. This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

### Building Meshery
Meshery is written in `Go` (Golang) and leverages Go Modules. The `deployment_yaml` folder contains the configuration yaml to deploy Meshery on Kubernetes, which includes a Deployment, Service, Service Entries and Virtual Services configurations.

A sample Makefile is included to build and package the app as a Docker image.
1. `Docker` to build the image.
1. `Go` version 1.11+ installed if you want to make changes to the existing code.
1. Clone this repository (`git clone https://github.com/layer5io/meshery.git`).
1. Build the Meshery Docker image (`docker build -t layer5/meshery .`).
    1. _pre-built images available: https://hub.docker.com/u/layer5/_

## License

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).

**About Layer5**
[Layer5.io](https://layer5.io) is a service mesh community, serving as a repository for information pertaining to the surrounding technology ecosystem (service meshes, api gateways, edge proxies, ingress and egress controllers) of microservice management in cloud native environments.
