<!--- app-name: Apache Kafka -->

# Bitnami package for Apache Kafka

Apache Kafka is a distributed streaming platform designed to build real-time pipelines and can be used as a message broker or as a replacement for a log aggregation solution for big data applications.

[Overview of Apache Kafka](http://kafka.apache.org/)

Trademarks: This software listing is packaged by Bitnami. The respective trademarks mentioned in the offering are owned by the respective companies, and use of them does not imply any affiliation or endorsement.

## TL;DR

```console
helm install my-release oci://registry-1.docker.io/bitnamicharts/kafka
```

> Tip: Did you know that this app is also available as a Kubernetes App on the Azure Marketplace? Kubernetes Apps are the easiest way to deploy Bitnami on AKS. Click [here](https://azuremarketplace.microsoft.com/en-us/marketplace/apps/bitnami.kafka-cnab) to see the listing on Azure Marketplace.

Looking to use Apache Kafka in production? Try [VMware Tanzu Application Catalog](https://bitnami.com/enterprise), the commercial edition of the Bitnami catalog.

## ⚠️ Important Notice: Upcoming changes to the Bitnami Catalog

Beginning August 28th, 2025, Bitnami will evolve its public catalog to offer a curated set of hardened, security-focused images under the new [Bitnami Secure Images initiative](https://news.broadcom.com/app-dev/broadcom-introduces-bitnami-secure-images-for-production-ready-containerized-applications). As part of this transition:

- Granting community users access for the first time to security-optimized versions of popular container images.
- Bitnami will begin deprecating support for non-hardened, Debian-based software images in its free tier and will gradually remove non-latest tags from the public catalog. As a result, community users will have access to a reduced number of hardened images. These images are published only under the “latest” tag and are intended for development purposes
- Starting August 28th, over two weeks, all existing container images, including older or versioned tags (e.g., 2.50.0, 10.6), will be migrated from the public catalog (docker.io/bitnami) to the “Bitnami Legacy” repository (docker.io/bitnamilegacy), where they will no longer receive updates.
- For production workloads and long-term support, users are encouraged to adopt Bitnami Secure Images, which include hardened containers, smaller attack surfaces, CVE transparency (via VEX/KEV), SBOMs, and enterprise support.

These changes aim to improve the security posture of all Bitnami users by promoting best practices for software supply chain integrity and up-to-date deployments. For more details, visit the [Bitnami Secure Images announcement](https://github.com/bitnami/containers/issues/83267).

## Introduction

This chart bootstraps a [Kafka](https://github.com/bitnami/containers/tree/main/bitnami/kafka) deployment on a [Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.23+
- Helm 3.8.0+
- PV provisioner support in the underlying infrastructure

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release oci://REGISTRY_NAME/REPOSITORY_NAME/kafka
```

> Note: You need to substitute the placeholders `REGISTRY_NAME` and `REPOSITORY_NAME` with a reference to your Helm chart registry and repository. For example, in the case of Bitnami, you need to use `REGISTRY_NAME=registry-1.docker.io` and `REPOSITORY_NAME=bitnamicharts`.

These commands deploy Kafka on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Configuration and installation details

### Listeners configuration

This chart allows you to automatically configure Kafka with 4 listeners:

- One for controller communications.
- A second one for inter-broker communications.
- A third one for communications with clients within the K8s cluster.
- (optional) a forth listener for communications with clients outside the K8s cluster. Check [this section](#accessing-kafka-brokers-from-outside-the-cluster) for more information.

For more complex configurations, set the `listeners`, `advertisedListeners` and `listenerSecurityProtocolMap` parameters as needed.

### Enable security for Kafka

You can configure different authentication protocols for each listener you configure in Kafka. For instance, you can use `sasl_tls` authentication for client communications, while using `tls` for controller and inter-broker communications. This table shows the available protocols and the security they provide:

| Method    | Authentication               | Encryption via TLS |
|-----------|------------------------------|--------------------|
| plaintext | None                         | No                 |
| tls       | None                         | Yes                |
| mtls      | Yes (two-way authentication) | Yes                |
| sasl      | Yes (via SASL)               | No                 |
| sasl_tls  | Yes (via SASL)               | Yes                |

Configure the authentication protocols for client, controller and inter-broker communications by setting the `listeners.client.protocol`, `listeners.controller.protocol` and `listeners.interbroker.protocol` parameters to the desired ones, respectively.

If you enabled SASL authentication on any listener, you can set the SASL credentials using the parameters below:

- `sasl.client.users`/`sasl.client.passwords`: when enabling SASL authentication for communications with clients.
- `sasl.interbroker.user`/`sasl.interbroker.password`:  when enabling SASL authentication for inter-broker communications.
- `sasl.controller.user`/`sasl.controller.password`:  when enabling SASL authentication for controller communications.

In order to configure TLS authentication/encryption, you **can** create a secret per Kafka node you have in the cluster containing the Java Key Stores (JKS) files: the truststore (`kafka.truststore.jks`) and the keystore (`kafka.keystore.jks`). Then, you need pass the secret names with the `tls.existingSecret` parameter when deploying the chart.

> **Note**: If the JKS files are password protected (recommended), you will need to provide the password to get access to the keystores. To do so, use the `tls.keystorePassword` and `tls.truststorePassword` parameters to provide your passwords.

For instance, to configure TLS authentication on a Kafka cluster with 2 Kafka nodes use the commands below to create the secrets:

```console
kubectl create secret generic kafka-jks-0 --from-file=kafka.truststore.jks=./kafka.truststore.jks --from-file=kafka.keystore.jks=./kafka-0.keystore.jks
kubectl create secret generic kafka-jks-1 --from-file=kafka.truststore.jks=./kafka.truststore.jks --from-file=kafka.keystore.jks=./kafka-1.keystore.jks
```

> **Note**: the command above assumes you already created the truststore and keystores files. This [script](https://raw.githubusercontent.com/confluentinc/confluent-platform-security-tools/master/kafka-generate-ssl.sh) can help you with the JKS files generation.

If, for some reason (like using CertManager) you can not use the default JKS secret scheme, you can use the additional parameters:

- `tls.jksTruststoreSecret` to define additional secret, where the `kafka.truststore.jks` is being kept. The truststore password **must** be the same as in `tls.truststorePassword`
- `tls.jksTruststoreKey` to overwrite the default value of the truststore key (`kafka.truststore.jks`).

> **Note**: If you are using CertManager, particularly when an ACME issuer is used, the `ca.crt` field is not put in the `Secret` that CertManager creates. To handle this, the `tls.pemChainIncluded` property can be set to `true` and the initContainer created by this Chart will attempt to extract the intermediate certs from the `tls.crt` field of the secret (which is a PEM chain)
> **Note**: The truststore/keystore from above **must** be protected with the same passwords set in the `tls.keystorePassword` and `tls.truststorePassword` parameters.

You can deploy the chart with authentication using the following parameters:

```console
replicaCount=2
listeners.client.protocol=SASL
listeners.interbroker.protocol=TLS
tls.existingSecret=kafka-jks
tls.keystorePassword=jksPassword
tls.truststorePassword=jksPassword
sasl.client.users[0]=brokerUser
sasl.client.passwords[0]=brokerPassword

```

By setting the following parameter: `listeners.client.protocol=SSL` and `listener.client.sslClientAuth=required`, Kafka will require the clients to authenticate to Kafka brokers via certificate.

As result, we will be able to see in `kafka-authorizer.log` the events specific Subject: `[...] Principal = User:CN=kafka,OU=...,O=...,L=...,C=..,ST=... is [...]`.

### Update credentials

The Bitnami Kafka chart, when upgrading, reuses the secret previously rendered by the chart or the one specified in `sasl.existingSecret`. To update credentials, use one of the following:

- Run `helm upgrade` specifying new credentials in the `sasl` section as explained in the [authentication section](#enable-security-for-kafka).
- Run `helm upgrade` specifying a new secret in `sasl.existingSecret`

### Accessing Kafka brokers from outside the cluster

In order to access Kafka Brokers from outside the cluster, an additional listener and advertised listener must be configured. Additionally, a specific service per kafka pod will be created.

There are three ways of configuring external access. Using LoadBalancer services, using NodePort services or using ClusterIP services.

#### Using LoadBalancer services

You have two alternatives to use LoadBalancer services:

- Option A) Use random load balancer IPs using an **initContainer** that waits for the IPs to be ready and discover them automatically.

```console
externalAccess.enabled=true
externalAccess.broker.service.type=LoadBalancer
externalAccess.controller.service.type=LoadBalancer
externalAccess.broker.service.ports.external=9094
externalAccess.controller.service.ports.external=9094
defaultInitContainers.autoDiscovery.enabled=true
serviceAccount.create=true
broker.automountServiceAccountToken=true
controller.automountServiceAccountToken=true
rbac.create=true
```

Note: This option requires creating RBAC rules on clusters where RBAC policies are enabled.

- Option B) Manually specify the load balancer IPs:

```console
externalAccess.enabled=true
externalAccess.controller.service.type=LoadBalancer
externalAccess.controller.service.containerPorts.external=9094
externalAccess.controller.service.loadBalancerIPs[0]='external-ip-1'
externalAccess.controller.service.loadBalancerIPs[1]='external-ip-2'
externalAccess.broker.service.type=LoadBalancer
externalAccess.broker.service.ports.external=9094
externalAccess.broker.service.loadBalancerIPs[0]='external-ip-3'
externalAccess.broker.service.loadBalancerIPs[1]='external-ip-4'
```

Note: You need to know in advance the load balancer IPs so each Kafka broker advertised listener is configured with it.

Following the aforementioned steps will also allow to connect the brokers from the outside using the cluster's default service (when `service.type` is `LoadBalancer` or `NodePort`). Use the property `service.externalPort` to specify the port used for external connections.

#### Using NodePort services

You have two alternatives to use NodePort services:

- Option A) Use random node ports using an **initContainer** that discover them automatically.

  ```console
  externalAccess.enabled=true
  externalAccess.controller.service.type=NodePort
  externalAccess.broker.service.type=NodePort
  defaultInitContainers.autoDiscovery.enabled=true
  serviceAccount.create=true
  rbac.create=true
  ```

  Note: This option requires creating RBAC rules on clusters where RBAC policies are enabled.

- Option B) Manually specify the node ports:

  ```console
  externalAccess.enabled=true
  externalAccess.controller.service.type=NodePort
  externalAccess.controller.service.nodePorts[0]='node-port-1'
  externalAccess.controller.service.nodePorts[1]='node-port-2'
  ```

  Note: You need to know in advance the node ports that will be exposed so each Kafka broker advertised listener is configured with it.

  The pod will try to get the external ip of the node using `curl -s https://ipinfo.io/ip` unless `externalAccess.<controller|broker>.service.domain` or `externalAccess.<controller|broker>.service.useHostIPs` is provided.

- Option C) Manually specify distinct external IPs (using controller+broker nodes)

  ```console
  externalAccess.enabled=true
  externalAccess.controller.service.type=NodePort
  externalAccess.controller.service.externalIPs[0]='172.16.0.20'
  externalAccess.controller.service.externalIPs[1]='172.16.0.21'
  externalAccess.controller.service.externalIPs[2]='172.16.0.22'
  ```

  Note: You need to know in advance the available IP of your cluster that will be exposed so each Kafka broker advertised listener is configured with it.

#### Using ClusterIP services

Note: This option requires that an ingress is deployed within your cluster

```console
externalAccess.enabled=true
externalAccess.controller.service.type=ClusterIP
externalAccess.controller.service.ports.external=9094
externalAccess.controller.service.domain='ingress-ip'
externalAccess.broker.service.type=ClusterIP
externalAccess.broker.service.ports.external=9094
externalAccess.broker.service.domain='ingress-ip'
```

Note: the deployed ingress must contain the following block:

```console
tcp:
  9094: "{{ include "common.names.namespace" . }}/{{ include "common.names.fullname" . }}-0-external:9094"
  9095: "{{ include "common.names.namespace" . }}/{{ include "common.names.fullname" . }}-1-external:9094"
  9096: "{{ include "common.names.namespace" . }}/{{ include "common.names.fullname" . }}-2-external:9094"
```

#### Name resolution with External-DNS

You can use the following values to generate External-DNS annotations which automatically creates DNS records for each ReplicaSet pod:

```yaml
externalAccess:
  controller:
    service:
      annotations:
        external-dns.alpha.kubernetes.io/hostname: "{{ .targetPod }}.example.com"
```

### Resource requests and limits

Bitnami charts allow setting resource requests and limits for all containers inside the chart deployment. These are inside the `resources` values (check parameters table). Setting requests is essential for production workloads and these should be adapted to your specific use case.

To make this process easier, the chart contains the `resourcesPreset` values, which automatically sets the `resources` section according to different presets. Check these presets in [the bitnami/common chart](https://github.com/bitnami/charts/blob/main/bitnami/common/templates/_resources.tpl#L15). However, in production workloads using `resourcesPreset` is discouraged as it may not fully adapt to your specific needs. Find more information on container resource management in the [official Kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/).

### Prometheus metrics

### Enable metrics

This chart can be integrated with Prometheus by setting `metrics.jmx.enabled` to `true`. This will deploy a sidecar container with [jmx_exporter](https://github.com/prometheus/jmx_exporter) in all pods and a `metrics` service, which can be configured under the `metrics.jmx.service` section. This service will have the necessary annotations to be automatically scraped by Prometheus.

#### Prometheus requirements

It is necessary to have a working installation of Prometheus or Prometheus Operator for the integration to work. Install the [Bitnami Prometheus helm chart](https://github.com/bitnami/charts/tree/main/bitnami/prometheus) or the [Bitnami Kube Prometheus helm chart](https://github.com/bitnami/charts/tree/main/bitnami/kube-prometheus) to easily have a working Prometheus in your cluster.

#### Integration with Prometheus Operator

The chart can deploy `ServiceMonitor` objects for integration with Prometheus Operator installations. To do so, set the value `metrics.serviceMonitor.enabled=true`. Ensure that the Prometheus Operator `CustomResourceDefinitions` are installed in the cluster or it will fail with the following error:

```text
no matches for kind "ServiceMonitor" in version "monitoring.coreos.com/v1"
```

Install the [Bitnami Kube Prometheus helm chart](https://github.com/bitnami/charts/tree/main/bitnami/kube-prometheus) for having the necessary CRDs and the Prometheus Operator.

### [Rolling VS Immutable tags](https://techdocs.broadcom.com/us/en/vmware-tanzu/application-catalog/tanzu-application-catalog/services/tac-doc/apps-tutorials-understand-rolling-tags-containers-index.html)

It is strongly recommended to use immutable tags in a production environment. This ensures your deployment does not change automatically if the same tag is updated with a different image.

Bitnami will release a new chart updating its containers if a new version of the main container, significant changes, or critical vulnerabilities exist.

### Sidecars

If you have a need for additional containers to run within the same pod as Kafka (e.g. an additional metrics or logging exporter), you can do so via the `sidecars` parameters. Simply define your container according to the Kubernetes container spec.

```yaml
sidecars:
  - name: your-image-name
    image: your-image
    imagePullPolicy: Always
    ports:
      - name: portname
       containerPort: 1234
```

### Setting Pod's affinity

This chart allows you to set your custom affinity using the `affinity` parameter. Find more information about Pod's affinity in the [kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity).

As an alternative, you can use of the preset configurations for pod affinity, pod anti-affinity, and node affinity available at the [bitnami/common](https://github.com/bitnami/charts/tree/main/bitnami/common#affinities) chart. To do so, set the `podAffinityPreset`, `podAntiAffinityPreset`, or `nodeAffinityPreset` parameters.

### Deploying extra resources

There are cases where you may want to deploy extra objects, such as Kafka Connect. For covering this case, the chart allows adding the full specification of other objects using the `extraDeploy` parameter. The following example would create a deployment including a Kafka Connect deployment so you can connect Kafka with MongoDB&reg;:

```yaml
extraDeploy:
  - |
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: {{ include "common.names.fullname" . }}-connect
      labels: {{- include "common.labels.standard" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 4 }}
        app.kubernetes.io/component: connector
    spec:
      replicas: 1
      selector:
        matchLabels: {{- include "common.labels.matchLabels" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 6 }}
          app.kubernetes.io/component: connector
      template:
        metadata:
          labels: {{- include "common.labels.standard" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 8 }}
            app.kubernetes.io/component: connector
        spec:
          containers:
            - name: connect
              image: KAFKA-CONNECT-IMAGE
              imagePullPolicy: IfNotPresent
              ports:
                - name: connector
                  containerPort: 8083
              volumeMounts:
                - name: configuration
                  mountPath: /bitnami/kafka/config
          volumes:
            - name: configuration
              configMap:
                name: {{ include "common.names.fullname" . }}-connect
  - |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: {{ include "common.names.fullname" . }}-connect
      labels: {{- include "common.labels.standard" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 4 }}
        app.kubernetes.io/component: connector
    data:
      connect-standalone.properties: |-
        bootstrap.servers = {{ include "common.names.fullname" . }}-controller-0.{{ include "common.names.fullname" . }}-controller-headless.{{ include "common.names.namespace" . }}.svc.{{ .Values.clusterDomain }}:{{ .Values.service.ports.client }}
        ...
      mongodb.properties: |-
        connection.uri=mongodb://root:password@mongodb-hostname:27017
        ...
  - |
    apiVersion: v1
    kind: Service
    metadata:
      name: {{ include "common.names.fullname" . }}-connect
      labels: {{- include "common.labels.standard" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 4 }}
        app.kubernetes.io/component: connector
    spec:
      ports:
        - protocol: TCP
          port: 8083
          targetPort: connector
      selector: {{- include "common.labels.matchLabels" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 4 }}
        app.kubernetes.io/component: connector
```

You can create the Kafka Connect image using the Dockerfile below:

```Dockerfile
FROM bitnami/kafka:latest
# Download MongoDB&reg; Connector for Apache Kafka https://www.confluent.io/hub/mongodb/kafka-connect-mongodb
RUN mkdir -p /opt/bitnami/kafka/plugins && \
    cd /opt/bitnami/kafka/plugins && \
    curl --remote-name --location --silent https://search.maven.org/remotecontent?filepath=org/mongodb/kafka/mongo-kafka-connect/1.2.0/mongo-kafka-connect-1.2.0-all.jar
CMD /opt/bitnami/kafka/bin/connect-standalone.sh /bitnami/kafka/config/connect-standalone.properties /bitnami/kafka/config/mongo.properties
```

### Persistence

The [Bitnami Kafka](https://github.com/bitnami/containers/tree/main/bitnami/kafka) image stores the Kafka data at the `/bitnami/kafka` path of the container. Persistent Volume Claims are used to keep the data across deployments. This is known to work in GCE, AWS, and minikube.

#### Adjust permissions of persistent volume mountpoint

As the image run as non-root by default, it is necessary to adjust the ownership of the persistent volume so that the container can write data into it.

By default, the chart is configured to use Kubernetes Security Context to automatically change the ownership of the volume. However, this feature does not work in all Kubernetes distributions. As an alternative, this chart supports using an initContainer to change the ownership of the volume before mounting it in the final destination.

You can enable this initContainer by setting `volumePermissions.enabled` to `true`.

#### Backup and restore

To back up and restore Helm chart deployments on Kubernetes, you need to back up the persistent volumes from the source deployment and attach them to a new deployment using [Velero](https://velero.io/), a Kubernetes backup/restore tool. Find the instructions for using Velero in [this guide](https://techdocs.broadcom.com/us/en/vmware-tanzu/application-catalog/tanzu-application-catalog/services/tac-doc/apps-tutorials-backup-restore-deployments-velero-index.html).

## Parameters

### Global parameters

| Name                                                  | Description                                                                                                                                                                                                                                                                                                                                                         | Value   |
| ----------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `global.imageRegistry`                                | Global Docker image registry                                                                                                                                                                                                                                                                                                                                        | `""`    |
| `global.imagePullSecrets`                             | Global Docker registry secret names as an array                                                                                                                                                                                                                                                                                                                     | `[]`    |
| `global.defaultStorageClass`                          | Global default StorageClass for Persistent Volume(s)                                                                                                                                                                                                                                                                                                                | `""`    |
| `global.security.allowInsecureImages`                 | Allows skipping image verification                                                                                                                                                                                                                                                                                                                                  | `false` |
| `global.compatibility.openshift.adaptSecurityContext` | Adapt the securityContext sections of the deployment to make them compatible with Openshift restricted-v2 SCC: remove runAsUser, runAsGroup and fsGroup and let the platform use their allowed default IDs. Possible values: auto (apply if the detected running cluster is Openshift), force (perform the adaptation always), disabled (do not perform adaptation) | `auto`  |

### Common parameters

| Name                      | Description                                                                             | Value           |
| ------------------------- | --------------------------------------------------------------------------------------- | --------------- |
| `kubeVersion`             | Override Kubernetes version                                                             | `""`            |
| `apiVersions`             | Override Kubernetes API versions reported by .Capabilities                              | `[]`            |
| `nameOverride`            | String to partially override common.names.fullname                                      | `""`            |
| `fullnameOverride`        | String to fully override common.names.fullname                                          | `""`            |
| `namespaceOverride`       | String to fully override common.names.namespace                                         | `""`            |
| `clusterDomain`           | Default Kubernetes cluster domain                                                       | `cluster.local` |
| `commonLabels`            | Labels to add to all deployed objects                                                   | `{}`            |
| `commonAnnotations`       | Annotations to add to all deployed objects                                              | `{}`            |
| `extraDeploy`             | Array of extra objects to deploy with the release                                       | `[]`            |
| `usePasswordFiles`        | Mount credentials as files instead of using environment variables                       | `true`          |
| `diagnosticMode.enabled`  | Enable diagnostic mode (all probes will be disabled and the command will be overridden) | `false`         |
| `diagnosticMode.command`  | Command to override all containers in the chart release                                 | `["sleep"]`     |
| `diagnosticMode.args`     | Args to override all containers in the chart release                                    | `["infinity"]`  |
| `serviceBindings.enabled` | Create secret for service binding (Experimental)                                        | `false`         |

### Kafka common parameters

| Name                                  | Description                                                                                                                                                                                                | Value                                                 |
| ------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| `image.registry`                      | Kafka image registry                                                                                                                                                                                       | `REGISTRY_NAME`                                       |
| `image.repository`                    | Kafka image repository                                                                                                                                                                                     | `REPOSITORY_NAME/kafka`                               |
| `image.digest`                        | Kafka image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                                                                                                      | `""`                                                  |
| `image.pullPolicy`                    | Kafka image pull policy                                                                                                                                                                                    | `IfNotPresent`                                        |
| `image.pullSecrets`                   | Specify docker-registry secret names as an array                                                                                                                                                           | `[]`                                                  |
| `image.debug`                         | Specify if debug values should be set                                                                                                                                                                      | `false`                                               |
| `clusterId`                           | Kafka Kraft cluster ID (ignored if existingKraftSecret is set). A random cluster ID will be generated the 1st time Kraft is initialized if not set.                                                        | `""`                                                  |
| `existingKraftSecret`                 | Name of the secret containing the Kafka KRaft Cluster ID and one directory ID per controller replica                                                                                                       | `""`                                                  |
| `kraftVersion`                        | Kraft version to be used. It determines whether static quorum (kraftVersion=0) or dynamic quorum (kraftVersion=1) will be used.                                                                            | `1`                                                   |
| `config`                              | Specify content for Kafka configuration (auto-generated based on other parameters otherwise)                                                                                                               | `{}`                                                  |
| `overrideConfiguration`               | Kafka common configuration override. Values defined here takes precedence over the ones defined at `config`                                                                                                | `{}`                                                  |
| `existingConfigmap`                   | Name of an existing ConfigMap with the Kafka configuration                                                                                                                                                 | `""`                                                  |
| `secretConfig`                        | Additional configuration to be appended at the end of the generated Kafka configuration (store in a secret)                                                                                                | `""`                                                  |
| `existingSecretConfig`                | Secret with additional configuration that will be appended to the end of the generated Kafka configuration                                                                                                 | `""`                                                  |
| `log4j2`                              | Specify content for Kafka log4j2 configuration (default one is used otherwise)                                                                                                                             | `""`                                                  |
| `existingLog4j2ConfigMap`             | The name of an existing ConfigMap containing the log4j2.yaml file                                                                                                                                          | `""`                                                  |
| `heapOpts`                            | Kafka Java Heap configuration                                                                                                                                                                              | `-XX:InitialRAMPercentage=75 -XX:MaxRAMPercentage=75` |
| `brokerRackAwareness.enabled`         | Enable Kafka Rack Awareness                                                                                                                                                                                | `false`                                               |
| `brokerRackAwareness.cloudProvider`   | Cloud provider to use to set Broker Rack Awareness. Allowed values: `aws-az`, `azure`                                                                                                                      | `""`                                                  |
| `brokerRackAwareness.azureApiVersion` | Metadata API version to use when brokerRackAwareness.cloudProvider is set to `azure`                                                                                                                       | `2023-11-15`                                          |
| `interBrokerProtocolVersion`          | Override the setting 'inter.broker.protocol.version' during the ZK migration.                                                                                                                              | `""`                                                  |
| `listeners.client.name`               | Name for the Kafka client listener                                                                                                                                                                         | `CLIENT`                                              |
| `listeners.client.containerPort`      | Port for the Kafka client listener                                                                                                                                                                         | `9092`                                                |
| `listeners.client.protocol`           | Security protocol for the Kafka client listener. Allowed values are 'PLAINTEXT', 'SASL_PLAINTEXT', 'SASL_SSL' and 'SSL'                                                                                    | `SASL_PLAINTEXT`                                      |
| `listeners.client.sslClientAuth`      | Optional. If SASL_SSL is enabled, configure mTLS TLS authentication type. If SSL protocol is enabled, overrides tls.authType for this listener. Allowed values are 'none', 'requested' and 'required'      | `""`                                                  |
| `listeners.controller.name`           | Name for the Kafka controller listener                                                                                                                                                                     | `CONTROLLER`                                          |
| `listeners.controller.containerPort`  | Port for the Kafka controller listener                                                                                                                                                                     | `9093`                                                |
| `listeners.controller.protocol`       | Security protocol for the Kafka controller listener. Allowed values are 'PLAINTEXT', 'SASL_PLAINTEXT', 'SASL_SSL' and 'SSL'                                                                                | `SASL_PLAINTEXT`                                      |
| `listeners.controller.sslClientAuth`  | Optional. If SASL_SSL is enabled, configure mTLS TLS authentication type. If SSL protocol is enabled, overrides tls.authType for this listener. Allowed values are 'none', 'requested' and 'required'      | `""`                                                  |
| `listeners.interbroker.name`          | Name for the Kafka inter-broker listener                                                                                                                                                                   | `INTERNAL`                                            |
| `listeners.interbroker.containerPort` | Port for the Kafka inter-broker listener                                                                                                                                                                   | `9094`                                                |
| `listeners.interbroker.protocol`      | Security protocol for the Kafka inter-broker listener. Allowed values are 'PLAINTEXT', 'SASL_PLAINTEXT', 'SASL_SSL' and 'SSL'                                                                              | `SASL_PLAINTEXT`                                      |
| `listeners.interbroker.sslClientAuth` | Optional. If SASL_SSL is enabled, configure mTLS TLS authentication type. If SSL protocol is enabled, overrides tls.authType for this listener. Allowed values are 'none', 'requested' and 'required'      | `""`                                                  |
| `listeners.external.containerPort`    | Port for the Kafka external listener                                                                                                                                                                       | `9095`                                                |
| `listeners.external.protocol`         | Security protocol for the Kafka external listener. . Allowed values are 'PLAINTEXT', 'SASL_PLAINTEXT', 'SASL_SSL' and 'SSL'                                                                                | `SASL_PLAINTEXT`                                      |
| `listeners.external.name`             | Name for the Kafka external listener                                                                                                                                                                       | `EXTERNAL`                                            |
| `listeners.external.sslClientAuth`    | Optional. If SASL_SSL is enabled, configure mTLS TLS authentication type. If SSL protocol is enabled, overrides tls.sslClientAuth for this listener. Allowed values are 'none', 'requested' and 'required' | `""`                                                  |
| `listeners.extraListeners`            | Array of listener objects to be appended to already existing listeners                                                                                                                                     | `[]`                                                  |
| `listeners.overrideListeners`         | Overrides the Kafka 'listeners' configuration setting.                                                                                                                                                     | `""`                                                  |
| `listeners.advertisedListeners`       | Overrides the Kafka 'advertised.listener' configuration setting.                                                                                                                                           | `""`                                                  |
| `listeners.securityProtocolMap`       | Overrides the Kafka 'security.protocol.map' configuration setting.                                                                                                                                         | `""`                                                  |

### Kafka SASL parameters

| Name                                | Description                                                                                                                                                                                   | Value                               |
| ----------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------- |
| `sasl.enabledMechanisms`            | Comma-separated list of allowed SASL mechanisms when SASL listeners are configured. Allowed types: `PLAIN`, `SCRAM-SHA-256`, `SCRAM-SHA-512`, `OAUTHBEARER`                                   | `PLAIN,SCRAM-SHA-256,SCRAM-SHA-512` |
| `sasl.interBrokerMechanism`         | SASL mechanism for inter broker communication.                                                                                                                                                | `PLAIN`                             |
| `sasl.controllerMechanism`          | SASL mechanism for controller communications.                                                                                                                                                 | `PLAIN`                             |
| `sasl.oauthbearer.tokenEndpointUrl` | The URL for the OAuth/OIDC identity provider                                                                                                                                                  | `""`                                |
| `sasl.oauthbearer.jwksEndpointUrl`  | The OAuth/OIDC provider URL from which the provider's JWKS (JSON Web Key Set) can be retrieved                                                                                                | `""`                                |
| `sasl.oauthbearer.expectedAudience` | The comma-delimited setting for the broker to use to verify that the JWT was issued for one of the expected audiences                                                                         | `""`                                |
| `sasl.oauthbearer.subClaimName`     | The OAuth claim name for the subject.                                                                                                                                                         | `sub`                               |
| `sasl.interbroker.user`             | Username for inter-broker communications when SASL is enabled                                                                                                                                 | `inter_broker_user`                 |
| `sasl.interbroker.password`         | Password for inter-broker communications when SASL is enabled. If not set and SASL is enabled for the controller listener, a random password will be generated.                               | `""`                                |
| `sasl.interbroker.clientId`         | Client ID for inter-broker communications when SASL is enabled with mechanism OAUTHBEARER                                                                                                     | `inter_broker_client`               |
| `sasl.interbroker.clientSecret`     | Client Secret for inter-broker communications when SASL is enabled with mechanism OAUTHBEARER. If not set and SASL is enabled for the controller listener, a random secret will be generated. | `""`                                |
| `sasl.controller.user`              | Username for controller communications when SASL is enabled                                                                                                                                   | `controller_user`                   |
| `sasl.controller.password`          | Password for controller communications when SASL is enabled. If not set and SASL is enabled for the inter-broker listener, a random password will be generated.                               | `""`                                |
| `sasl.controller.clientId`          | Client ID for controller communications when SASL is enabled with mechanism OAUTHBEARER                                                                                                       | `controller_broker_client`          |
| `sasl.controller.clientSecret`      | Client Secret for controller communications when SASL is enabled with mechanism OAUTHBEARER. If not set and SASL is enabled for the inter-broker listener, a random secret will be generated. | `""`                                |
| `sasl.client.users`                 | Comma-separated list of usernames for client communications when SASL is enabled                                                                                                              | `["user1"]`                         |
| `sasl.client.passwords`             | Comma-separated list of passwords for client communications when SASL is enabled, must match the number of client.users                                                                       | `""`                                |
| `sasl.existingSecret`               | Name of the existing secret containing credentials for client.users, interbroker.user and controller.user                                                                                     | `""`                                |

### Kafka TLS parameters

| Name                                                                                        | Description                                                                                                                                                                                                                                                                                                                | Value                      |
| ------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| `tls.type`                                                                                  | Format to use for TLS certificates. Allowed types: `JKS` and `PEM`                                                                                                                                                                                                                                                         | `JKS`                      |
| `tls.pemChainIncluded`                                                                      | Flag to denote that the Certificate Authority (CA) certificates are bundled with the endpoint cert.                                                                                                                                                                                                                        | `false`                    |
| `tls.autoGenerated.enabled`                                                                 | Enable automatic generation of TLS certificates (only supported if `tls.type` is `PEM`)                                                                                                                                                                                                                                    | `true`                     |
| `tls.autoGenerated.engine`                                                                  | Mechanism to generate the certificates (allowed values: helm, cert-manager)                                                                                                                                                                                                                                                | `helm`                     |
| `tls.autoGenerated.customAltNames`                                                          | List of additional subject alternative names (SANs) for the automatically generated TLS certificates.                                                                                                                                                                                                                      | `[]`                       |
| `tls.autoGenerated.certManager.existingIssuer`                                              | The name of an existing Issuer to use for generating the certificates (only for `cert-manager` engine)                                                                                                                                                                                                                     | `""`                       |
| `tls.autoGenerated.certManager.existingIssuerKind`                                          | Existing Issuer kind, defaults to Issuer (only for `cert-manager` engine)                                                                                                                                                                                                                                                  | `""`                       |
| `tls.autoGenerated.certManager.keyAlgorithm`                                                | Key algorithm for the certificates (only for `cert-manager` engine)                                                                                                                                                                                                                                                        | `RSA`                      |
| `tls.autoGenerated.certManager.keySize`                                                     | Key size for the certificates (only for `cert-manager` engine)                                                                                                                                                                                                                                                             | `2048`                     |
| `tls.autoGenerated.certManager.duration`                                                    | Duration for the certificates (only for `cert-manager` engine)                                                                                                                                                                                                                                                             | `2160h`                    |
| `tls.autoGenerated.certManager.renewBefore`                                                 | Renewal period for the certificates (only for `cert-manager` engine)                                                                                                                                                                                                                                                       | `360h`                     |
| `tls.existingSecret`                                                                        | Name of the existing secret containing the TLS certificates for the Kafka nodes.                                                                                                                                                                                                                                           | `""`                       |
| `tls.passwordsSecret`                                                                       | Name of the secret containing the password to access the JKS files or PEM key when they are password-protected. (`key`: `password`)                                                                                                                                                                                        | `""`                       |
| `tls.passwordsSecretKeystoreKey`                                                            | The secret key from the tls.passwordsSecret containing the password for the Keystore.                                                                                                                                                                                                                                      | `keystore-password`        |
| `tls.passwordsSecretTruststoreKey`                                                          | The secret key from the tls.passwordsSecret containing the password for the Truststore.                                                                                                                                                                                                                                    | `truststore-password`      |
| `tls.passwordsSecretPemPasswordKey`                                                         | The secret key from the tls.passwordsSecret containing the password for the PEM key inside 'tls.passwordsSecret'.                                                                                                                                                                                                          | `""`                       |
| `tls.keystorePassword`                                                                      | Password to access the JKS keystore when it is password-protected. Ignored when 'tls.passwordsSecret' is provided.                                                                                                                                                                                                         | `""`                       |
| `tls.truststorePassword`                                                                    | Password to access the JKS truststore when it is password-protected. Ignored when 'tls.passwordsSecret' is provided.                                                                                                                                                                                                       | `""`                       |
| `tls.keyPassword`                                                                           | Password to access the PEM key when it is password-protected.                                                                                                                                                                                                                                                              | `""`                       |
| `tls.jksKeystoreKey`                                                                        | The secret key from the `tls.existingSecret` containing the keystore                                                                                                                                                                                                                                                       | `""`                       |
| `tls.jksTruststoreSecret`                                                                   | Name of the existing secret containing your truststore if truststore not existing or different from the one in the `tls.existingSecret`                                                                                                                                                                                    | `""`                       |
| `tls.jksTruststoreKey`                                                                      | The secret key from the `tls.existingSecret` or `tls.jksTruststoreSecret` containing the truststore                                                                                                                                                                                                                        | `""`                       |
| `tls.endpointIdentificationAlgorithm`                                                       | The endpoint identification algorithm to validate server hostname using server certificate                                                                                                                                                                                                                                 | `https`                    |
| `tls.sslClientAuth`                                                                         | Sets the default value for the ssl.client.auth Kafka setting.                                                                                                                                                                                                                                                              | `required`                 |
| `extraEnvVars`                                                                              | Extra environment variables to add to Kafka pods                                                                                                                                                                                                                                                                           | `[]`                       |
| `extraEnvVarsCM`                                                                            | ConfigMap with extra environment variables                                                                                                                                                                                                                                                                                 | `""`                       |
| `extraEnvVarsSecret`                                                                        | Secret with extra environment variables                                                                                                                                                                                                                                                                                    | `""`                       |
| `extraVolumes`                                                                              | Optionally specify extra list of additional volumes for the Kafka pod(s)                                                                                                                                                                                                                                                   | `[]`                       |
| `extraVolumeMounts`                                                                         | Optionally specify extra list of additional volumeMounts for the Kafka container(s)                                                                                                                                                                                                                                        | `[]`                       |
| `sidecars`                                                                                  | Add additional sidecar containers to the Kafka pod(s)                                                                                                                                                                                                                                                                      | `[]`                       |
| `initContainers`                                                                            | Add additional Add init containers to the Kafka pod(s)                                                                                                                                                                                                                                                                     | `[]`                       |
| `dnsPolicy`                                                                                 | Specifies the DNS policy for the Kafka pods                                                                                                                                                                                                                                                                                | `""`                       |
| `dnsConfig`                                                                                 | allows users more control on the DNS settings for a Pod. Required if `dnsPolicy` is set to `None`                                                                                                                                                                                                                          | `{}`                       |
| `defaultInitContainers.volumePermissions.enabled`                                           | Enable init container that changes the owner and group of the persistent volume                                                                                                                                                                                                                                            | `false`                    |
| `defaultInitContainers.volumePermissions.image.registry`                                    | "volume-permissions" init-containers' image registry                                                                                                                                                                                                                                                                       | `REGISTRY_NAME`            |
| `defaultInitContainers.volumePermissions.image.repository`                                  | "volume-permissions" init-containers' image repository                                                                                                                                                                                                                                                                     | `REPOSITORY_NAME/os-shell` |
| `defaultInitContainers.volumePermissions.image.digest`                                      | "volume-permissions" init-containers' image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                                                                                                                                                                                      | `""`                       |
| `defaultInitContainers.volumePermissions.image.pullPolicy`                                  | "volume-permissions" init-containers' image pull policy                                                                                                                                                                                                                                                                    | `IfNotPresent`             |
| `defaultInitContainers.volumePermissions.image.pullSecrets`                                 | "volume-permissions" init-containers' image pull secrets                                                                                                                                                                                                                                                                   | `[]`                       |
| `defaultInitContainers.volumePermissions.containerSecurityContext.enabled`                  | Enabled "volume-permissions" init-containers' Security Context                                                                                                                                                                                                                                                             | `true`                     |
| `defaultInitContainers.volumePermissions.containerSecurityContext.seLinuxOptions`           | Set SELinux options in "volume-permissions" init-containers                                                                                                                                                                                                                                                                | `{}`                       |
| `defaultInitContainers.volumePermissions.containerSecurityContext.runAsUser`                | Set runAsUser in "volume-permissions" init-containers' Security Context                                                                                                                                                                                                                                                    | `0`                        |
| `defaultInitContainers.volumePermissions.containerSecurityContext.privileged`               | Set privileged in "volume-permissions" init-containers' Security Context                                                                                                                                                                                                                                                   | `false`                    |
| `defaultInitContainers.volumePermissions.containerSecurityContext.allowPrivilegeEscalation` | Set allowPrivilegeEscalation in "volume-permissions" init-containers' Security Context                                                                                                                                                                                                                                     | `false`                    |
| `defaultInitContainers.volumePermissions.containerSecurityContext.capabilities.add`         | List of capabilities to be added in "volume-permissions" init-containers                                                                                                                                                                                                                                                   | `[]`                       |
| `defaultInitContainers.volumePermissions.containerSecurityContext.capabilities.drop`        | List of capabilities to be dropped in "volume-permissions" init-containers                                                                                                                                                                                                                                                 | `["ALL"]`                  |
| `defaultInitContainers.volumePermissions.containerSecurityContext.seccompProfile.type`      | Set seccomp profile in "volume-permissions" init-containers                                                                                                                                                                                                                                                                | `RuntimeDefault`           |
| `defaultInitContainers.volumePermissions.resourcesPreset`                                   | Set Kafka "volume-permissions" init container resources according to one common preset (allowed values: none, nano, small, medium, large, xlarge, 2xlarge). This is ignored if defaultInitContainers.volumePermissions.resources is set (defaultInitContainers.volumePermissions.resources is recommended for production). | `nano`                     |
| `defaultInitContainers.volumePermissions.resources`                                         | Set Kafka "volume-permissions" init container requests and limits for different resources like CPU or memory (essential for production workloads)                                                                                                                                                                          | `{}`                       |
| `defaultInitContainers.prepareConfig.containerSecurityContext.enabled`                      | Enabled "prepare-config" init-containers' Security Context                                                                                                                                                                                                                                                                 | `true`                     |
| `defaultInitContainers.prepareConfig.containerSecurityContext.seLinuxOptions`               | Set SELinux options in "prepare-config" init-containers                                                                                                                                                                                                                                                                    | `{}`                       |
| `defaultInitContainers.prepareConfig.containerSecurityContext.runAsUser`                    | Set runAsUser in "prepare-config" init-containers' Security Context                                                                                                                                                                                                                                                        | `1001`                     |
| `defaultInitContainers.prepareConfig.containerSecurityContext.runAsGroup`                   | Set runAsUser in "prepare-config" init-containers' Security Context                                                                                                                                                                                                                                                        | `1001`                     |
| `defaultInitContainers.prepareConfig.containerSecurityContext.runAsNonRoot`                 | Set runAsNonRoot in "prepare-config" init-containers' Security Context                                                                                                                                                                                                                                                     | `true`                     |
| `defaultInitContainers.prepareConfig.containerSecurityContext.readOnlyRootFilesystem`       | Set readOnlyRootFilesystem in "prepare-config" init-containers' Security Context                                                                                                                                                                                                                                           | `true`                     |
| `defaultInitContainers.prepareConfig.containerSecurityContext.privileged`                   | Set privileged in "prepare-config" init-containers' Security Context                                                                                                                                                                                                                                                       | `false`                    |
| `defaultInitContainers.prepareConfig.containerSecurityContext.allowPrivilegeEscalation`     | Set allowPrivilegeEscalation in "prepare-config" init-containers' Security Context                                                                                                                                                                                                                                         | `false`                    |
| `defaultInitContainers.prepareConfig.containerSecurityContext.capabilities.add`             | List of capabilities to be added in "prepare-config" init-containers                                                                                                                                                                                                                                                       | `[]`                       |
| `defaultInitContainers.prepareConfig.containerSecurityContext.capabilities.drop`            | List of capabilities to be dropped in "prepare-config" init-containers                                                                                                                                                                                                                                                     | `["ALL"]`                  |
| `defaultInitContainers.prepareConfig.containerSecurityContext.seccompProfile.type`          | Set seccomp profile in "prepare-config" init-containers                                                                                                                                                                                                                                                                    | `RuntimeDefault`           |
| `defaultInitContainers.prepareConfig.resourcesPreset`                                       | Set Kafka "prepare-config" init container resources according to one common preset (allowed values: none, nano, small, medium, large, xlarge, 2xlarge). This is ignored if defaultInitContainers.prepareConfig.resources is set (defaultInitContainers.prepareConfig.resources is recommended for production).             | `nano`                     |
| `defaultInitContainers.prepareConfig.resources`                                             | Set Kafka "prepare-config" init container requests and limits for different resources like CPU or memory (essential for production workloads)                                                                                                                                                                              | `{}`                       |
| `defaultInitContainers.prepareConfig.extraInit`                                             | Additional content for the "prepare-config" init script, rendered as a template.                                                                                                                                                                                                                                           | `""`                       |
| `defaultInitContainers.autoDiscovery.enabled`                                               | Enable init container that auto-detects external IPs/ports by querying the K8s API                                                                                                                                                                                                                                         | `false`                    |
| `defaultInitContainers.autoDiscovery.image.registry`                                        | "auto-discovery" init-containers' image registry                                                                                                                                                                                                                                                                           | `REGISTRY_NAME`            |
| `defaultInitContainers.autoDiscovery.image.repository`                                      | "auto-discovery" init-containers' image repository                                                                                                                                                                                                                                                                         | `REPOSITORY_NAME/os-shell` |
| `defaultInitContainers.autoDiscovery.image.digest`                                          | "auto-discovery" init-containers' image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                                                                                                                                                                                          | `""`                       |
| `defaultInitContainers.autoDiscovery.image.pullPolicy`                                      | "auto-discovery" init-containers' image pull policy                                                                                                                                                                                                                                                                        | `IfNotPresent`             |
| `defaultInitContainers.autoDiscovery.image.pullSecrets`                                     | "auto-discovery" init-containers' image pull secrets                                                                                                                                                                                                                                                                       | `[]`                       |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.enabled`                      | Enabled "auto-discovery" init-containers' Security Context                                                                                                                                                                                                                                                                 | `true`                     |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.seLinuxOptions`               | Set SELinux options in "auto-discovery" init-containers                                                                                                                                                                                                                                                                    | `{}`                       |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.runAsUser`                    | Set runAsUser in "auto-discovery" init-containers' Security Context                                                                                                                                                                                                                                                        | `1001`                     |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.runAsGroup`                   | Set runAsUser in "auto-discovery" init-containers' Security Context                                                                                                                                                                                                                                                        | `1001`                     |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.runAsNonRoot`                 | Set runAsNonRoot in "auto-discovery" init-containers' Security Context                                                                                                                                                                                                                                                     | `true`                     |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.readOnlyRootFilesystem`       | Set readOnlyRootFilesystem in "auto-discovery" init-containers' Security Context                                                                                                                                                                                                                                           | `true`                     |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.privileged`                   | Set privileged in "auto-discovery" init-containers' Security Context                                                                                                                                                                                                                                                       | `false`                    |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.allowPrivilegeEscalation`     | Set allowPrivilegeEscalation in "auto-discovery" init-containers' Security Context                                                                                                                                                                                                                                         | `false`                    |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.capabilities.add`             | List of capabilities to be added in "auto-discovery" init-containers                                                                                                                                                                                                                                                       | `[]`                       |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.capabilities.drop`            | List of capabilities to be dropped in "auto-discovery" init-containers                                                                                                                                                                                                                                                     | `["ALL"]`                  |
| `defaultInitContainers.autoDiscovery.containerSecurityContext.seccompProfile.type`          | Set seccomp profile in "auto-discovery" init-containers                                                                                                                                                                                                                                                                    | `RuntimeDefault`           |
| `defaultInitContainers.autoDiscovery.resourcesPreset`                                       | Set Kafka "auto-discovery" init container resources according to one common preset (allowed values: none, nano, small, medium, large, xlarge, 2xlarge). This is ignored if defaultInitContainers.autoDiscovery.resources is set (defaultInitContainers.autoDiscovery.resources is recommended for production).             | `nano`                     |
| `defaultInitContainers.autoDiscovery.resources`                                             | Set Kafka "auto-discovery" init container requests and limits for different resources like CPU or memory (essential for production workloads)                                                                                                                                                                              | `{}`                       |

### Controller-eligible statefulset parameters

| Name                                                           | Description                                                                                                                                                                                                                             | Value                                                 |
| -------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| `controller.replicaCount`                                      | Number of Kafka controller-eligible nodes                                                                                                                                                                                               | `3`                                                   |
| `controller.controllerOnly`                                    | If set to true, controller nodes will be deployed as dedicated controllers, instead of controller+broker processes.                                                                                                                     | `false`                                               |
| `controller.quorumBootstrapServers`                            | Override the Kafka controller quorum bootstrap servers of the Kafka Kraft cluster. If not set, it will be automatically configured to use all controller-eligible nodes.                                                                | `""`                                                  |
| `controller.minId`                                             | Minimal node.id values for controller-eligible nodes. Do not change after first initialization.                                                                                                                                         | `0`                                                   |
| `controller.config`                                            | Specify content for Kafka configuration for Kafka controller-eligible nodes (auto-generated based on other parameters otherwise)                                                                                                        | `{}`                                                  |
| `controller.overrideConfiguration`                             | Kafka configuration override for Kafka controller-eligible nodes. Values defined here takes precedence over the ones defined at `controller.config`                                                                                     | `{}`                                                  |
| `controller.existingConfigmap`                                 | Name of an existing ConfigMap with the Kafka configuration for Kafka controller-eligible nodes                                                                                                                                          | `""`                                                  |
| `controller.secretConfig`                                      | Additional configuration to be appended at the end of the generated Kafka configuration for Kafka controller-eligible nodes (store in a secret)                                                                                         | `""`                                                  |
| `controller.existingSecretConfig`                              | Secret with additional configuration that will be appended to the end of the generated Kafka configuration for Kafka controller-eligible nodes                                                                                          | `""`                                                  |
| `controller.heapOpts`                                          | Kafka Java Heap configuration for controller-eligible nodes                                                                                                                                                                             | `-XX:InitialRAMPercentage=75 -XX:MaxRAMPercentage=75` |
| `controller.command`                                           | Override Kafka container command                                                                                                                                                                                                        | `[]`                                                  |
| `controller.args`                                              | Override Kafka container arguments                                                                                                                                                                                                      | `[]`                                                  |
| `controller.extraEnvVars`                                      | Extra environment variables to add to Kafka pods                                                                                                                                                                                        | `[]`                                                  |
| `controller.extraEnvVarsCM`                                    | ConfigMap with extra environment variables                                                                                                                                                                                              | `""`                                                  |
| `controller.extraEnvVarsSecret`                                | Secret with extra environment variables                                                                                                                                                                                                 | `""`                                                  |
| `controller.extraContainerPorts`                               | Kafka controller-eligible extra containerPorts.                                                                                                                                                                                         | `[]`                                                  |
| `controller.livenessProbe.enabled`                             | Enable livenessProbe on Kafka containers                                                                                                                                                                                                | `true`                                                |
| `controller.livenessProbe.initialDelaySeconds`                 | Initial delay seconds for livenessProbe                                                                                                                                                                                                 | `10`                                                  |
| `controller.livenessProbe.periodSeconds`                       | Period seconds for livenessProbe                                                                                                                                                                                                        | `10`                                                  |
| `controller.livenessProbe.timeoutSeconds`                      | Timeout seconds for livenessProbe                                                                                                                                                                                                       | `5`                                                   |
| `controller.livenessProbe.failureThreshold`                    | Failure threshold for livenessProbe                                                                                                                                                                                                     | `3`                                                   |
| `controller.livenessProbe.successThreshold`                    | Success threshold for livenessProbe                                                                                                                                                                                                     | `1`                                                   |
| `controller.readinessProbe.enabled`                            | Enable readinessProbe on Kafka containers                                                                                                                                                                                               | `true`                                                |
| `controller.readinessProbe.initialDelaySeconds`                | Initial delay seconds for readinessProbe                                                                                                                                                                                                | `5`                                                   |
| `controller.readinessProbe.periodSeconds`                      | Period seconds for readinessProbe                                                                                                                                                                                                       | `10`                                                  |
| `controller.readinessProbe.timeoutSeconds`                     | Timeout seconds for readinessProbe                                                                                                                                                                                                      | `5`                                                   |
| `controller.readinessProbe.failureThreshold`                   | Failure threshold for readinessProbe                                                                                                                                                                                                    | `6`                                                   |
| `controller.readinessProbe.successThreshold`                   | Success threshold for readinessProbe                                                                                                                                                                                                    | `1`                                                   |
| `controller.startupProbe.enabled`                              | Enable startupProbe on Kafka containers                                                                                                                                                                                                 | `false`                                               |
| `controller.startupProbe.initialDelaySeconds`                  | Initial delay seconds for startupProbe                                                                                                                                                                                                  | `30`                                                  |
| `controller.startupProbe.periodSeconds`                        | Period seconds for startupProbe                                                                                                                                                                                                         | `10`                                                  |
| `controller.startupProbe.timeoutSeconds`                       | Timeout seconds for startupProbe                                                                                                                                                                                                        | `1`                                                   |
| `controller.startupProbe.failureThreshold`                     | Failure threshold for startupProbe                                                                                                                                                                                                      | `15`                                                  |
| `controller.startupProbe.successThreshold`                     | Success threshold for startupProbe                                                                                                                                                                                                      | `1`                                                   |
| `controller.customLivenessProbe`                               | Custom livenessProbe that overrides the default one                                                                                                                                                                                     | `{}`                                                  |
| `controller.customReadinessProbe`                              | Custom readinessProbe that overrides the default one                                                                                                                                                                                    | `{}`                                                  |
| `controller.customStartupProbe`                                | Custom startupProbe that overrides the default one                                                                                                                                                                                      | `{}`                                                  |
| `controller.lifecycleHooks`                                    | lifecycleHooks for the Kafka container to automate configuration before or after startup                                                                                                                                                | `{}`                                                  |
| `controller.resourcesPreset`                                   | Set container resources according to one common preset (allowed values: none, nano, micro, small, medium, large, xlarge, 2xlarge). This is ignored if controller.resources is set (controller.resources is recommended for production). | `small`                                               |
| `controller.resources`                                         | Set container requests and limits for different resources like CPU or memory (essential for production workloads)                                                                                                                       | `{}`                                                  |
| `controller.podSecurityContext.enabled`                        | Enable security context for the pods                                                                                                                                                                                                    | `true`                                                |
| `controller.podSecurityContext.fsGroupChangePolicy`            | Set filesystem group change policy                                                                                                                                                                                                      | `Always`                                              |
| `controller.podSecurityContext.sysctls`                        | Set kernel settings using the sysctl interface                                                                                                                                                                                          | `[]`                                                  |
| `controller.podSecurityContext.supplementalGroups`             | Set filesystem extra groups                                                                                                                                                                                                             | `[]`                                                  |
| `controller.podSecurityContext.fsGroup`                        | Set Kafka pod's Security Context fsGroup                                                                                                                                                                                                | `1001`                                                |
| `controller.podSecurityContext.seccompProfile.type`            | Set Kafka pods's Security Context seccomp profile                                                                                                                                                                                       | `RuntimeDefault`                                      |
| `controller.containerSecurityContext.enabled`                  | Enable Kafka containers' Security Context                                                                                                                                                                                               | `true`                                                |
| `controller.containerSecurityContext.seLinuxOptions`           | Set SELinux options in container                                                                                                                                                                                                        | `{}`                                                  |
| `controller.containerSecurityContext.runAsUser`                | Set containers' Security Context runAsUser                                                                                                                                                                                              | `1001`                                                |
| `controller.containerSecurityContext.runAsGroup`               | Set containers' Security Context runAsGroup                                                                                                                                                                                             | `1001`                                                |
| `controller.containerSecurityContext.runAsGroup`               | Set Kafka containers' Security Context runAsGroup                                                                                                                                                                                       | `1001`                                                |
| `controller.containerSecurityContext.runAsNonRoot`             | Set Kafka containers' Security Context runAsNonRoot                                                                                                                                                                                     | `true`                                                |
| `controller.containerSecurityContext.allowPrivilegeEscalation` | Force the child process to be run as non-privileged                                                                                                                                                                                     | `false`                                               |
| `controller.containerSecurityContext.readOnlyRootFilesystem`   | Allows the pod to mount the RootFS as ReadOnly only                                                                                                                                                                                     | `true`                                                |
| `controller.containerSecurityContext.capabilities.drop`        | Set Kafka containers' server Security Context capabilities to be dropped                                                                                                                                                                | `["ALL"]`                                             |
| `controller.automountServiceAccountToken`                      | Mount Service Account token in pod                                                                                                                                                                                                      | `false`                                               |
| `controller.hostAliases`                                       | Kafka pods host aliases                                                                                                                                                                                                                 | `[]`                                                  |
| `controller.hostNetwork`                                       | Specify if host network should be enabled for Kafka pods                                                                                                                                                                                | `false`                                               |
| `controller.hostIPC`                                           | Specify if host IPC should be enabled for Kafka pods                                                                                                                                                                                    | `false`                                               |
| `controller.podLabels`                                         | Extra labels for Kafka pods                                                                                                                                                                                                             | `{}`                                                  |
| `controller.podAnnotations`                                    | Extra annotations for Kafka pods                                                                                                                                                                                                        | `{}`                                                  |
| `controller.topologyKey`                                       | Override common lib default topology key. If empty - "kubernetes.io/hostname" is used                                                                                                                                                   | `""`                                                  |
| `controller.podAffinityPreset`                                 | Pod affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                                                                     | `""`                                                  |
| `controller.podAntiAffinityPreset`                             | Pod anti-affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                                                                | `soft`                                                |
| `controller.nodeAffinityPreset.type`                           | Node affinity preset type. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                                                               | `""`                                                  |
| `controller.nodeAffinityPreset.key`                            | Node label key to match Ignored if `affinity` is set.                                                                                                                                                                                   | `""`                                                  |
| `controller.nodeAffinityPreset.values`                         | Node label values to match. Ignored if `affinity` is set.                                                                                                                                                                               | `[]`                                                  |
| `controller.affinity`                                          | Affinity for pod assignment                                                                                                                                                                                                             | `{}`                                                  |
| `controller.nodeSelector`                                      | Node labels for pod assignment                                                                                                                                                                                                          | `{}`                                                  |
| `controller.tolerations`                                       | Tolerations for pod assignment                                                                                                                                                                                                          | `[]`                                                  |
| `controller.topologySpreadConstraints`                         | Topology Spread Constraints for pod assignment spread across your cluster among failure-domains. Evaluated as a template                                                                                                                | `[]`                                                  |
| `controller.terminationGracePeriodSeconds`                     | Seconds the pod needs to gracefully terminate                                                                                                                                                                                           | `""`                                                  |
| `controller.podManagementPolicy`                               | StatefulSet controller supports relax its ordering guarantees while preserving its uniqueness and identity guarantees. There are two valid pod management policies: OrderedReady and Parallel                                           | `Parallel`                                            |
| `controller.minReadySeconds`                                   | How many seconds a pod needs to be ready before killing the next, during update                                                                                                                                                         | `0`                                                   |
| `controller.priorityClassName`                                 | Name of the existing priority class to be used by kafka pods                                                                                                                                                                            | `""`                                                  |
| `controller.runtimeClassName`                                  | Name of the runtime class to be used by pod(s)                                                                                                                                                                                          | `""`                                                  |
| `controller.enableServiceLinks`                                | Whether information about services should be injected into pod's environment variable                                                                                                                                                   | `true`                                                |
| `controller.schedulerName`                                     | Name of the k8s scheduler (other than default)                                                                                                                                                                                          | `""`                                                  |
| `controller.updateStrategy.type`                               | Kafka statefulset strategy type                                                                                                                                                                                                         | `RollingUpdate`                                       |
| `controller.extraVolumes`                                      | Optionally specify extra list of additional volumes for the Kafka pod(s)                                                                                                                                                                | `[]`                                                  |
| `controller.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts for the Kafka container(s)                                                                                                                                                     | `[]`                                                  |
| `controller.sidecars`                                          | Add additional sidecar containers to the Kafka pod(s)                                                                                                                                                                                   | `[]`                                                  |
| `controller.initContainers`                                    | Add additional Add init containers to the Kafka pod(s)                                                                                                                                                                                  | `[]`                                                  |

### Experimental: Kafka Controller Autoscaling configuration

| Name                                                          | Description                                                                                                                                                            | Value                     |
| ------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------- |
| `controller.autoscaling.vpa.enabled`                          | Enable VPA                                                                                                                                                             | `false`                   |
| `controller.autoscaling.vpa.annotations`                      | Annotations for VPA resource                                                                                                                                           | `{}`                      |
| `controller.autoscaling.vpa.controlledResources`              | VPA List of resources that the vertical pod autoscaler can control. Defaults to cpu and memory                                                                         | `[]`                      |
| `controller.autoscaling.vpa.maxAllowed`                       | VPA Max allowed resources for the pod                                                                                                                                  | `{}`                      |
| `controller.autoscaling.vpa.minAllowed`                       | VPA Min allowed resources for the pod                                                                                                                                  | `{}`                      |
| `controller.autoscaling.vpa.updatePolicy.updateMode`          | Autoscaling update policy Specifies whether recommended updates are applied when a Pod is started and whether recommended updates are applied during the life of a Pod | `Auto`                    |
| `controller.autoscaling.hpa.enabled`                          | Enable HPA for Kafka Controller                                                                                                                                        | `false`                   |
| `controller.autoscaling.hpa.annotations`                      | Annotations for HPA resource                                                                                                                                           | `{}`                      |
| `controller.autoscaling.hpa.minReplicas`                      | Minimum number of Kafka Controller replicas                                                                                                                            | `""`                      |
| `controller.autoscaling.hpa.maxReplicas`                      | Maximum number of Kafka Controller replicas                                                                                                                            | `""`                      |
| `controller.autoscaling.hpa.targetCPU`                        | Target CPU utilization percentage                                                                                                                                      | `""`                      |
| `controller.autoscaling.hpa.targetMemory`                     | Target Memory utilization percentage                                                                                                                                   | `""`                      |
| `controller.pdb.create`                                       | Deploy a pdb object for the Kafka pod                                                                                                                                  | `true`                    |
| `controller.pdb.minAvailable`                                 | Minimum number/percentage of available Kafka replicas                                                                                                                  | `""`                      |
| `controller.pdb.maxUnavailable`                               | Maximum number/percentage of unavailable Kafka replicas                                                                                                                | `""`                      |
| `controller.persistentVolumeClaimRetentionPolicy.enabled`     | Controls if and how PVCs are deleted during the lifecycle of a StatefulSet                                                                                             | `false`                   |
| `controller.persistentVolumeClaimRetentionPolicy.whenScaled`  | Volume retention behavior when the replica count of the StatefulSet is reduced                                                                                         | `Retain`                  |
| `controller.persistentVolumeClaimRetentionPolicy.whenDeleted` | Volume retention behavior that applies when the StatefulSet is deleted                                                                                                 | `Retain`                  |
| `controller.persistence.enabled`                              | Enable Kafka data persistence using PVC                                                                                                                                | `true`                    |
| `controller.persistence.existingClaim`                        | A manually managed Persistent Volume and Claim                                                                                                                         | `""`                      |
| `controller.persistence.storageClass`                         | PVC Storage Class for Kafka data volume                                                                                                                                | `""`                      |
| `controller.persistence.accessModes`                          | Persistent Volume Access Modes                                                                                                                                         | `["ReadWriteOnce"]`       |
| `controller.persistence.size`                                 | PVC Storage Request for Kafka data volume                                                                                                                              | `8Gi`                     |
| `controller.persistence.annotations`                          | Annotations for the PVC                                                                                                                                                | `{}`                      |
| `controller.persistence.labels`                               | Labels for the PVC                                                                                                                                                     | `{}`                      |
| `controller.persistence.selector`                             | Selector to match an existing Persistent Volume for Kafka data PVC. If set, the PVC can't have a PV dynamically provisioned for it                                     | `{}`                      |
| `controller.persistence.mountPath`                            | Mount path of the Kafka data volume                                                                                                                                    | `/bitnami/kafka`          |
| `controller.logPersistence.enabled`                           | Enable Kafka logs persistence using PVC                                                                                                                                | `false`                   |
| `controller.logPersistence.existingClaim`                     | A manually managed Persistent Volume and Claim                                                                                                                         | `""`                      |
| `controller.logPersistence.storageClass`                      | PVC Storage Class for Kafka logs volume                                                                                                                                | `""`                      |
| `controller.logPersistence.accessModes`                       | Persistent Volume Access Modes                                                                                                                                         | `["ReadWriteOnce"]`       |
| `controller.logPersistence.size`                              | PVC Storage Request for Kafka logs volume                                                                                                                              | `8Gi`                     |
| `controller.logPersistence.annotations`                       | Annotations for the PVC                                                                                                                                                | `{}`                      |
| `controller.logPersistence.selector`                          | Selector to match an existing Persistent Volume for Kafka log data PVC. If set, the PVC can't have a PV dynamically provisioned for it                                 | `{}`                      |
| `controller.logPersistence.mountPath`                         | Mount path of the Kafka logs volume                                                                                                                                    | `/opt/bitnami/kafka/logs` |

### Broker-only statefulset parameters

| Name                                                       | Description                                                                                                                                                                                                                     | Value                                                 |
| ---------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| `broker.replicaCount`                                      | Number of Kafka broker-only nodes                                                                                                                                                                                               | `0`                                                   |
| `broker.minId`                                             | Minimal node.id values for broker-only nodes. Do not change after first initialization.                                                                                                                                         | `100`                                                 |
| `broker.config`                                            | Specify content for Kafka configuration for Kafka broker-only nodes (auto-generated based on other parameters otherwise)                                                                                                        | `{}`                                                  |
| `broker.overrideConfiguration`                             | Kafka configuration override for Kafka broker-only nodes. Values defined here takes precedence over the ones defined at `broker.config`                                                                                         | `{}`                                                  |
| `broker.existingConfigmap`                                 | Name of an existing ConfigMap with the Kafka configuration for Kafka broker-only nodes                                                                                                                                          | `""`                                                  |
| `broker.secretConfig`                                      | Additional configuration to be appended at the end of the generated Kafka configuration for Kafka broker-only nodes (store in a secret)                                                                                         | `""`                                                  |
| `broker.existingSecretConfig`                              | Secret with additional configuration that will be appended to the end of the generated Kafka configuration for Kafka broker-only nodes                                                                                          | `""`                                                  |
| `broker.heapOpts`                                          | Kafka Java Heap configuration for broker-only nodes                                                                                                                                                                             | `-XX:InitialRAMPercentage=75 -XX:MaxRAMPercentage=75` |
| `broker.command`                                           | Override Kafka container command                                                                                                                                                                                                | `[]`                                                  |
| `broker.args`                                              | Override Kafka container arguments                                                                                                                                                                                              | `[]`                                                  |
| `broker.extraEnvVars`                                      | Extra environment variables to add to Kafka pods                                                                                                                                                                                | `[]`                                                  |
| `broker.extraEnvVarsCM`                                    | ConfigMap with extra environment variables                                                                                                                                                                                      | `""`                                                  |
| `broker.extraEnvVarsSecret`                                | Secret with extra environment variables                                                                                                                                                                                         | `""`                                                  |
| `broker.extraContainerPorts`                               | Kafka broker-only extra containerPorts.                                                                                                                                                                                         | `[]`                                                  |
| `broker.livenessProbe.enabled`                             | Enable livenessProbe on Kafka containers                                                                                                                                                                                        | `true`                                                |
| `broker.livenessProbe.initialDelaySeconds`                 | Initial delay seconds for livenessProbe                                                                                                                                                                                         | `10`                                                  |
| `broker.livenessProbe.periodSeconds`                       | Period seconds for livenessProbe                                                                                                                                                                                                | `10`                                                  |
| `broker.livenessProbe.timeoutSeconds`                      | Timeout seconds for livenessProbe                                                                                                                                                                                               | `5`                                                   |
| `broker.livenessProbe.failureThreshold`                    | Failure threshold for livenessProbe                                                                                                                                                                                             | `3`                                                   |
| `broker.livenessProbe.successThreshold`                    | Success threshold for livenessProbe                                                                                                                                                                                             | `1`                                                   |
| `broker.readinessProbe.enabled`                            | Enable readinessProbe on Kafka containers                                                                                                                                                                                       | `true`                                                |
| `broker.readinessProbe.initialDelaySeconds`                | Initial delay seconds for readinessProbe                                                                                                                                                                                        | `5`                                                   |
| `broker.readinessProbe.periodSeconds`                      | Period seconds for readinessProbe                                                                                                                                                                                               | `10`                                                  |
| `broker.readinessProbe.timeoutSeconds`                     | Timeout seconds for readinessProbe                                                                                                                                                                                              | `5`                                                   |
| `broker.readinessProbe.failureThreshold`                   | Failure threshold for readinessProbe                                                                                                                                                                                            | `6`                                                   |
| `broker.readinessProbe.successThreshold`                   | Success threshold for readinessProbe                                                                                                                                                                                            | `1`                                                   |
| `broker.startupProbe.enabled`                              | Enable startupProbe on Kafka containers                                                                                                                                                                                         | `false`                                               |
| `broker.startupProbe.initialDelaySeconds`                  | Initial delay seconds for startupProbe                                                                                                                                                                                          | `30`                                                  |
| `broker.startupProbe.periodSeconds`                        | Period seconds for startupProbe                                                                                                                                                                                                 | `10`                                                  |
| `broker.startupProbe.timeoutSeconds`                       | Timeout seconds for startupProbe                                                                                                                                                                                                | `1`                                                   |
| `broker.startupProbe.failureThreshold`                     | Failure threshold for startupProbe                                                                                                                                                                                              | `15`                                                  |
| `broker.startupProbe.successThreshold`                     | Success threshold for startupProbe                                                                                                                                                                                              | `1`                                                   |
| `broker.customLivenessProbe`                               | Custom livenessProbe that overrides the default one                                                                                                                                                                             | `{}`                                                  |
| `broker.customReadinessProbe`                              | Custom readinessProbe that overrides the default one                                                                                                                                                                            | `{}`                                                  |
| `broker.customStartupProbe`                                | Custom startupProbe that overrides the default one                                                                                                                                                                              | `{}`                                                  |
| `broker.lifecycleHooks`                                    | lifecycleHooks for the Kafka container to automate configuration before or after startup                                                                                                                                        | `{}`                                                  |
| `broker.resourcesPreset`                                   | Set container resources according to one common preset (allowed values: none, nano, micro, small, medium, large, xlarge, 2xlarge). This is ignored if broker.resources is set (broker.resources is recommended for production). | `small`                                               |
| `broker.resources`                                         | Set container requests and limits for different resources like CPU or memory (essential for production workloads)                                                                                                               | `{}`                                                  |
| `broker.podSecurityContext.enabled`                        | Enable security context for the pods                                                                                                                                                                                            | `true`                                                |
| `broker.podSecurityContext.fsGroupChangePolicy`            | Set filesystem group change policy                                                                                                                                                                                              | `Always`                                              |
| `broker.podSecurityContext.sysctls`                        | Set kernel settings using the sysctl interface                                                                                                                                                                                  | `[]`                                                  |
| `broker.podSecurityContext.supplementalGroups`             | Set filesystem extra groups                                                                                                                                                                                                     | `[]`                                                  |
| `broker.podSecurityContext.fsGroup`                        | Set Kafka pod's Security Context fsGroup                                                                                                                                                                                        | `1001`                                                |
| `broker.podSecurityContext.seccompProfile.type`            | Set Kafka pod's Security Context seccomp profile                                                                                                                                                                                | `RuntimeDefault`                                      |
| `broker.containerSecurityContext.enabled`                  | Enable Kafka containers' Security Context                                                                                                                                                                                       | `true`                                                |
| `broker.containerSecurityContext.seLinuxOptions`           | Set SELinux options in container                                                                                                                                                                                                | `{}`                                                  |
| `broker.containerSecurityContext.runAsUser`                | Set containers' Security Context runAsUser                                                                                                                                                                                      | `1001`                                                |
| `broker.containerSecurityContext.runAsGroup`               | Set containers' Security Context runAsGroup                                                                                                                                                                                     | `1001`                                                |
| `broker.containerSecurityContext.runAsNonRoot`             | Set Kafka containers' Security Context runAsNonRoot                                                                                                                                                                             | `true`                                                |
| `broker.containerSecurityContext.allowPrivilegeEscalation` | Force the child process to be run as non-privileged                                                                                                                                                                             | `false`                                               |
| `broker.containerSecurityContext.readOnlyRootFilesystem`   | Allows the pod to mount the RootFS as ReadOnly only                                                                                                                                                                             | `true`                                                |
| `broker.containerSecurityContext.capabilities.drop`        | Set Kafka containers' server Security Context capabilities to be dropped                                                                                                                                                        | `["ALL"]`                                             |
| `broker.automountServiceAccountToken`                      | Mount Service Account token in pod                                                                                                                                                                                              | `false`                                               |
| `broker.hostAliases`                                       | Kafka pods host aliases                                                                                                                                                                                                         | `[]`                                                  |
| `broker.hostNetwork`                                       | Specify if host network should be enabled for Kafka pods                                                                                                                                                                        | `false`                                               |
| `broker.hostIPC`                                           | Specify if host IPC should be enabled for Kafka pods                                                                                                                                                                            | `false`                                               |
| `broker.podLabels`                                         | Extra labels for Kafka pods                                                                                                                                                                                                     | `{}`                                                  |
| `broker.podAnnotations`                                    | Extra annotations for Kafka pods                                                                                                                                                                                                | `{}`                                                  |
| `broker.topologyKey`                                       | Override common lib default topology key. If empty - "kubernetes.io/hostname" is used                                                                                                                                           | `""`                                                  |
| `broker.podAffinityPreset`                                 | Pod affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                                                             | `""`                                                  |
| `broker.podAntiAffinityPreset`                             | Pod anti-affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                                                        | `soft`                                                |
| `broker.nodeAffinityPreset.type`                           | Node affinity preset type. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                                                       | `""`                                                  |
| `broker.nodeAffinityPreset.key`                            | Node label key to match Ignored if `affinity` is set.                                                                                                                                                                           | `""`                                                  |
| `broker.nodeAffinityPreset.values`                         | Node label values to match. Ignored if `affinity` is set.                                                                                                                                                                       | `[]`                                                  |
| `broker.affinity`                                          | Affinity for pod assignment                                                                                                                                                                                                     | `{}`                                                  |
| `broker.nodeSelector`                                      | Node labels for pod assignment                                                                                                                                                                                                  | `{}`                                                  |
| `broker.tolerations`                                       | Tolerations for pod assignment                                                                                                                                                                                                  | `[]`                                                  |
| `broker.topologySpreadConstraints`                         | Topology Spread Constraints for pod assignment spread across your cluster among failure-domains. Evaluated as a template                                                                                                        | `[]`                                                  |
| `broker.terminationGracePeriodSeconds`                     | Seconds the pod needs to gracefully terminate                                                                                                                                                                                   | `""`                                                  |
| `broker.podManagementPolicy`                               | StatefulSet controller supports relax its ordering guarantees while preserving its uniqueness and identity guarantees. There are two valid pod management policies: OrderedReady and Parallel                                   | `Parallel`                                            |
| `broker.minReadySeconds`                                   | How many seconds a pod needs to be ready before killing the next, during update                                                                                                                                                 | `0`                                                   |
| `broker.priorityClassName`                                 | Name of the existing priority class to be used by kafka pods                                                                                                                                                                    | `""`                                                  |
| `broker.runtimeClassName`                                  | Name of the runtime class to be used by pod(s)                                                                                                                                                                                  | `""`                                                  |
| `broker.enableServiceLinks`                                | Whether information about services should be injected into pod's environment variable                                                                                                                                           | `true`                                                |
| `broker.schedulerName`                                     | Name of the k8s scheduler (other than default)                                                                                                                                                                                  | `""`                                                  |
| `broker.updateStrategy.type`                               | Kafka statefulset strategy type                                                                                                                                                                                                 | `RollingUpdate`                                       |
| `broker.extraVolumes`                                      | Optionally specify extra list of additional volumes for the Kafka pod(s)                                                                                                                                                        | `[]`                                                  |
| `broker.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts for the Kafka container(s)                                                                                                                                             | `[]`                                                  |
| `broker.sidecars`                                          | Add additional sidecar containers to the Kafka pod(s)                                                                                                                                                                           | `[]`                                                  |
| `broker.initContainers`                                    | Add additional Add init containers to the Kafka pod(s)                                                                                                                                                                          | `[]`                                                  |
| `broker.pdb.create`                                        | Deploy a pdb object for the Kafka pod                                                                                                                                                                                           | `true`                                                |
| `broker.pdb.minAvailable`                                  | Maximum number/percentage of unavailable Kafka replicas                                                                                                                                                                         | `""`                                                  |
| `broker.pdb.maxUnavailable`                                | Maximum number/percentage of unavailable Kafka replicas                                                                                                                                                                         | `""`                                                  |

### Experimental: Kafka Broker Autoscaling configuration

| Name                                                      | Description                                                                                                                                                            | Value                     |
| --------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------- |
| `broker.autoscaling.vpa.enabled`                          | Enable VPA                                                                                                                                                             | `false`                   |
| `broker.autoscaling.vpa.annotations`                      | Annotations for VPA resource                                                                                                                                           | `{}`                      |
| `broker.autoscaling.vpa.controlledResources`              | VPA List of resources that the vertical pod autoscaler can control. Defaults to cpu and memory                                                                         | `[]`                      |
| `broker.autoscaling.vpa.maxAllowed`                       | VPA Max allowed resources for the pod                                                                                                                                  | `{}`                      |
| `broker.autoscaling.vpa.minAllowed`                       | VPA Min allowed resources for the pod                                                                                                                                  | `{}`                      |
| `broker.autoscaling.vpa.updatePolicy.updateMode`          | Autoscaling update policy Specifies whether recommended updates are applied when a Pod is started and whether recommended updates are applied during the life of a Pod | `Auto`                    |
| `broker.autoscaling.hpa.enabled`                          | Enable HPA for Kafka Broker                                                                                                                                            | `false`                   |
| `broker.autoscaling.hpa.annotations`                      | Annotations for HPA resource                                                                                                                                           | `{}`                      |
| `broker.autoscaling.hpa.minReplicas`                      | Minimum number of Kafka Broker replicas                                                                                                                                | `""`                      |
| `broker.autoscaling.hpa.maxReplicas`                      | Maximum number of Kafka Broker replicas                                                                                                                                | `""`                      |
| `broker.autoscaling.hpa.targetCPU`                        | Target CPU utilization percentage                                                                                                                                      | `""`                      |
| `broker.autoscaling.hpa.targetMemory`                     | Target Memory utilization percentage                                                                                                                                   | `""`                      |
| `broker.persistentVolumeClaimRetentionPolicy.enabled`     | Controls if and how PVCs are deleted during the lifecycle of a StatefulSet                                                                                             | `false`                   |
| `broker.persistentVolumeClaimRetentionPolicy.whenScaled`  | Volume retention behavior when the replica count of the StatefulSet is reduced                                                                                         | `Retain`                  |
| `broker.persistentVolumeClaimRetentionPolicy.whenDeleted` | Volume retention behavior that applies when the StatefulSet is deleted                                                                                                 | `Retain`                  |
| `broker.persistence.enabled`                              | Enable Kafka data persistence using PVC                                                                                                                                | `true`                    |
| `broker.persistence.existingClaim`                        | A manually managed Persistent Volume and Claim                                                                                                                         | `""`                      |
| `broker.persistence.storageClass`                         | PVC Storage Class for Kafka data volume                                                                                                                                | `""`                      |
| `broker.persistence.accessModes`                          | Persistent Volume Access Modes                                                                                                                                         | `["ReadWriteOnce"]`       |
| `broker.persistence.size`                                 | PVC Storage Request for Kafka data volume                                                                                                                              | `8Gi`                     |
| `broker.persistence.annotations`                          | Annotations for the PVC                                                                                                                                                | `{}`                      |
| `broker.persistence.labels`                               | Labels for the PVC                                                                                                                                                     | `{}`                      |
| `broker.persistence.selector`                             | Selector to match an existing Persistent Volume for Kafka data PVC. If set, the PVC can't have a PV dynamically provisioned for it                                     | `{}`                      |
| `broker.persistence.mountPath`                            | Mount path of the Kafka data volume                                                                                                                                    | `/bitnami/kafka`          |
| `broker.logPersistence.enabled`                           | Enable Kafka logs persistence using PVC                                                                                                                                | `false`                   |
| `broker.logPersistence.existingClaim`                     | A manually managed Persistent Volume and Claim                                                                                                                         | `""`                      |
| `broker.logPersistence.storageClass`                      | PVC Storage Class for Kafka logs volume                                                                                                                                | `""`                      |
| `broker.logPersistence.accessModes`                       | Persistent Volume Access Modes                                                                                                                                         | `["ReadWriteOnce"]`       |
| `broker.logPersistence.size`                              | PVC Storage Request for Kafka logs volume                                                                                                                              | `8Gi`                     |
| `broker.logPersistence.annotations`                       | Annotations for the PVC                                                                                                                                                | `{}`                      |
| `broker.logPersistence.selector`                          | Selector to match an existing Persistent Volume for Kafka log data PVC. If set, the PVC can't have a PV dynamically provisioned for it                                 | `{}`                      |
| `broker.logPersistence.mountPath`                         | Mount path of the Kafka logs volume                                                                                                                                    | `/opt/bitnami/kafka/logs` |

### Traffic Exposure parameters

| Name                                                              | Description                                                                                                                               | Value          |
| ----------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- | -------------- |
| `service.type`                                                    | Kubernetes Service type                                                                                                                   | `ClusterIP`    |
| `service.ports.client`                                            | Kafka svc port for client connections                                                                                                     | `9092`         |
| `service.ports.controller`                                        | Kafka svc port for controller connections                                                                                                 | `9093`         |
| `service.ports.interbroker`                                       | Kafka svc port for inter-broker connections                                                                                               | `9094`         |
| `service.ports.external`                                          | Kafka svc port for external connections                                                                                                   | `9095`         |
| `service.extraPorts`                                              | Extra ports to expose in the Kafka service (normally used with the `sidecar` value)                                                       | `[]`           |
| `service.nodePorts.client`                                        | Node port for the Kafka client connections                                                                                                | `""`           |
| `service.nodePorts.external`                                      | Node port for the Kafka external connections                                                                                              | `""`           |
| `service.sessionAffinity`                                         | Control where client requests go, to the same pod or round-robin                                                                          | `None`         |
| `service.sessionAffinityConfig`                                   | Additional settings for the sessionAffinity                                                                                               | `{}`           |
| `service.clusterIP`                                               | Kafka service Cluster IP                                                                                                                  | `""`           |
| `service.loadBalancerIP`                                          | Kafka service Load Balancer IP                                                                                                            | `""`           |
| `service.loadBalancerClass`                                       | Kafka service Load Balancer Class                                                                                                         | `""`           |
| `service.loadBalancerSourceRanges`                                | Kafka service Load Balancer sources                                                                                                       | `[]`           |
| `service.allocateLoadBalancerNodePorts`                           | Whether to allocate node ports when service type is LoadBalancer                                                                          | `true`         |
| `service.externalTrafficPolicy`                                   | Kafka service external traffic policy                                                                                                     | `Cluster`      |
| `service.annotations`                                             | Additional custom annotations for Kafka service                                                                                           | `{}`           |
| `service.headless.controller.annotations`                         | Annotations for the controller-eligible headless service.                                                                                 | `{}`           |
| `service.headless.controller.labels`                              | Labels for the controller-eligible headless service.                                                                                      | `{}`           |
| `service.headless.broker.annotations`                             | Annotations for the broker-only headless service.                                                                                         | `{}`           |
| `service.headless.broker.labels`                                  | Labels for the broker-only headless service.                                                                                              | `{}`           |
| `service.headless.ipFamilies`                                     | IP families for the headless service                                                                                                      | `[]`           |
| `service.headless.ipFamilyPolicy`                                 | IP family policy for the headless service                                                                                                 | `""`           |
| `externalAccess.enabled`                                          | Enable Kubernetes external cluster access to Kafka brokers                                                                                | `false`        |
| `externalAccess.controller.forceExpose`                           | If set to true, force exposing controller-eligible nodes although they are configured as controller-only nodes                            | `false`        |
| `externalAccess.controller.service.type`                          | Kubernetes Service type for external access. It can be NodePort, LoadBalancer or ClusterIP                                                | `LoadBalancer` |
| `externalAccess.controller.service.ports.external`                | Kafka port used for external access when service type is LoadBalancer                                                                     | `9094`         |
| `externalAccess.controller.service.loadBalancerClass`             | Kubernetes Service Load Balancer class for external access when service type is LoadBalancer                                              | `""`           |
| `externalAccess.controller.service.loadBalancerIPs`               | Array of load balancer IPs for each Kafka broker. Length must be the same as replicaCount                                                 | `[]`           |
| `externalAccess.controller.service.loadBalancerNames`             | Array of load balancer Names for each Kafka broker. Length must be the same as replicaCount                                               | `[]`           |
| `externalAccess.controller.service.loadBalancerAnnotations`       | Array of load balancer annotations for each Kafka broker. Length must be the same as replicaCount                                         | `[]`           |
| `externalAccess.controller.service.loadBalancerSourceRanges`      | Address(es) that are allowed when service is LoadBalancer                                                                                 | `[]`           |
| `externalAccess.controller.service.allocateLoadBalancerNodePorts` | Whether to allocate node ports when service type is LoadBalancer                                                                          | `true`         |
| `externalAccess.controller.service.nodePorts`                     | Array of node ports used for each Kafka broker. Length must be the same as replicaCount                                                   | `[]`           |
| `externalAccess.controller.service.externalIPs`                   | Use distinct service host IPs to configure Kafka external listener when service type is NodePort. Length must be the same as replicaCount | `[]`           |
| `externalAccess.controller.service.useHostIPs`                    | Use service host IPs to configure Kafka external listener when service type is NodePort                                                   | `false`        |
| `externalAccess.controller.service.usePodIPs`                     | using the MY_POD_IP address for external access.                                                                                          | `false`        |
| `externalAccess.controller.service.domain`                        | Domain or external ip used to configure Kafka external listener when service type is NodePort or ClusterIP                                | `""`           |
| `externalAccess.controller.service.publishNotReadyAddresses`      | Indicates that any agent which deals with endpoints for this Service should disregard any indications of ready/not-ready                  | `false`        |
| `externalAccess.controller.service.labels`                        | Service labels for external access                                                                                                        | `{}`           |
| `externalAccess.controller.service.annotations`                   | Service annotations for external access                                                                                                   | `{}`           |
| `externalAccess.controller.service.extraPorts`                    | Extra ports to expose in the Kafka external service                                                                                       | `[]`           |
| `externalAccess.controller.service.ipFamilies`                    | IP families for the external controller service                                                                                           | `[]`           |
| `externalAccess.controller.service.ipFamilyPolicy`                | IP family policy for the external controller service                                                                                      | `""`           |
| `externalAccess.broker.service.type`                              | Kubernetes Service type for external access. It can be NodePort, LoadBalancer or ClusterIP                                                | `LoadBalancer` |
| `externalAccess.broker.service.ports.external`                    | Kafka port used for external access when service type is LoadBalancer                                                                     | `9094`         |
| `externalAccess.broker.service.loadBalancerClass`                 | Kubernetes Service Load Balancer class for external access when service type is LoadBalancer                                              | `""`           |
| `externalAccess.broker.service.loadBalancerIPs`                   | Array of load balancer IPs for each Kafka broker. Length must be the same as replicaCount                                                 | `[]`           |
| `externalAccess.broker.service.loadBalancerNames`                 | Array of load balancer Names for each Kafka broker. Length must be the same as replicaCount                                               | `[]`           |
| `externalAccess.broker.service.loadBalancerAnnotations`           | Array of load balancer annotations for each Kafka broker. Length must be the same as replicaCount                                         | `[]`           |
| `externalAccess.broker.service.loadBalancerSourceRanges`          | Address(es) that are allowed when service is LoadBalancer                                                                                 | `[]`           |
| `externalAccess.broker.service.allocateLoadBalancerNodePorts`     | Whether to allocate node ports when service type is LoadBalancer                                                                          | `true`         |
| `externalAccess.broker.service.nodePorts`                         | Array of node ports used for each Kafka broker. Length must be the same as replicaCount                                                   | `[]`           |
| `externalAccess.broker.service.externalIPs`                       | Use distinct service host IPs to configure Kafka external listener when service type is NodePort. Length must be the same as replicaCount | `[]`           |
| `externalAccess.broker.service.useHostIPs`                        | Use service host IPs to configure Kafka external listener when service type is NodePort                                                   | `false`        |
| `externalAccess.broker.service.usePodIPs`                         | using the MY_POD_IP address for external access.                                                                                          | `false`        |
| `externalAccess.broker.service.domain`                            | Domain or external ip used to configure Kafka external listener when service type is NodePort or ClusterIP                                | `""`           |
| `externalAccess.broker.service.publishNotReadyAddresses`          | Indicates that any agent which deals with endpoints for this Service should disregard any indications of ready/not-ready                  | `false`        |
| `externalAccess.broker.service.labels`                            | Service labels for external access                                                                                                        | `{}`           |
| `externalAccess.broker.service.annotations`                       | Service annotations for external access                                                                                                   | `{}`           |
| `externalAccess.broker.service.extraPorts`                        | Extra ports to expose in the Kafka external service                                                                                       | `[]`           |
| `externalAccess.broker.service.ipFamilies`                        | IP families for the external broker service                                                                                               | `[]`           |
| `externalAccess.broker.service.ipFamilyPolicy`                    | IP family policy for the external broker service                                                                                          | `""`           |
| `networkPolicy.enabled`                                           | Specifies whether a NetworkPolicy should be created                                                                                       | `true`         |
| `networkPolicy.allowExternal`                                     | Don't require client label for connections                                                                                                | `true`         |
| `networkPolicy.allowExternalEgress`                               | Allow the pod to access any range of port and all destinations.                                                                           | `true`         |
| `networkPolicy.addExternalClientAccess`                           | Allow access from pods with client label set to "true". Ignored if `networkPolicy.allowExternal` is true.                                 | `true`         |
| `networkPolicy.extraIngress`                                      | Add extra ingress rules to the NetworkPolicy                                                                                              | `[]`           |
| `networkPolicy.extraEgress`                                       | Add extra ingress rules to the NetworkPolicy                                                                                              | `[]`           |
| `networkPolicy.ingressPodMatchLabels`                             | Labels to match to allow traffic from other pods. Ignored if `networkPolicy.allowExternal` is true.                                       | `{}`           |
| `networkPolicy.ingressNSMatchLabels`                              | Labels to match to allow traffic from other namespaces. Ignored if `networkPolicy.allowExternal` is true.                                 | `{}`           |
| `networkPolicy.ingressNSPodMatchLabels`                           | Pod labels to match to allow traffic from other namespaces. Ignored if `networkPolicy.allowExternal` is true.                             | `{}`           |

### Other Parameters

| Name                                          | Description                                                                                    | Value   |
| --------------------------------------------- | ---------------------------------------------------------------------------------------------- | ------- |
| `serviceAccount.create`                       | Enable creation of ServiceAccount for Kafka pods                                               | `true`  |
| `serviceAccount.name`                         | The name of the service account to use. If not set and `create` is `true`, a name is generated | `""`    |
| `serviceAccount.automountServiceAccountToken` | Allows auto mount of ServiceAccountToken on the serviceAccount created                         | `false` |
| `serviceAccount.annotations`                  | Additional custom annotations for the ServiceAccount                                           | `{}`    |
| `rbac.create`                                 | Whether to create & use RBAC resources or not                                                  | `false` |

### Metrics parameters

| Name                                                            | Description                                                                                                                                                                                                                               | Value                                                                                   |
| --------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `metrics.jmx.enabled`                                           | Whether or not to expose JMX metrics to Prometheus                                                                                                                                                                                        | `false`                                                                                 |
| `metrics.jmx.kafkaJmxPort`                                      | JMX port where the exporter will collect metrics, exposed in the Kafka container.                                                                                                                                                         | `5555`                                                                                  |
| `metrics.jmx.image.registry`                                    | JMX exporter image registry                                                                                                                                                                                                               | `REGISTRY_NAME`                                                                         |
| `metrics.jmx.image.repository`                                  | JMX exporter image repository                                                                                                                                                                                                             | `REPOSITORY_NAME/jmx-exporter`                                                          |
| `metrics.jmx.image.digest`                                      | JMX exporter image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                                                                                                                              | `""`                                                                                    |
| `metrics.jmx.image.pullPolicy`                                  | JMX exporter image pull policy                                                                                                                                                                                                            | `IfNotPresent`                                                                          |
| `metrics.jmx.image.pullSecrets`                                 | Specify docker-registry secret names as an array                                                                                                                                                                                          | `[]`                                                                                    |
| `metrics.jmx.containerSecurityContext.enabled`                  | Enable Prometheus JMX exporter containers' Security Context                                                                                                                                                                               | `true`                                                                                  |
| `metrics.jmx.containerSecurityContext.seLinuxOptions`           | Set SELinux options in container                                                                                                                                                                                                          | `{}`                                                                                    |
| `metrics.jmx.containerSecurityContext.runAsUser`                | Set containers' Security Context runAsUser                                                                                                                                                                                                | `1001`                                                                                  |
| `metrics.jmx.containerSecurityContext.runAsGroup`               | Set containers' Security Context runAsGroup                                                                                                                                                                                               | `1001`                                                                                  |
| `metrics.jmx.containerSecurityContext.runAsNonRoot`             | Set Prometheus JMX exporter containers' Security Context runAsNonRoot                                                                                                                                                                     | `true`                                                                                  |
| `metrics.jmx.containerSecurityContext.allowPrivilegeEscalation` | Set Prometheus JMX exporter containers' Security Context allowPrivilegeEscalation                                                                                                                                                         | `false`                                                                                 |
| `metrics.jmx.containerSecurityContext.readOnlyRootFilesystem`   | Set Prometheus JMX exporter containers' Security Context readOnlyRootFilesystem                                                                                                                                                           | `true`                                                                                  |
| `metrics.jmx.containerSecurityContext.capabilities.drop`        | Set Prometheus JMX exporter containers' Security Context capabilities to be dropped                                                                                                                                                       | `["ALL"]`                                                                               |
| `metrics.jmx.containerPorts.metrics`                            | Prometheus JMX exporter metrics container port                                                                                                                                                                                            | `5556`                                                                                  |
| `metrics.jmx.resourcesPreset`                                   | Set container resources according to one common preset (allowed values: none, nano, micro, small, medium, large, xlarge, 2xlarge). This is ignored if metrics.jmx.resources is set (metrics.jmx.resources is recommended for production). | `micro`                                                                                 |
| `metrics.jmx.resources`                                         | Set container requests and limits for different resources like CPU or memory (essential for production workloads)                                                                                                                         | `{}`                                                                                    |
| `metrics.jmx.livenessProbe.enabled`                             | Enable livenessProbe                                                                                                                                                                                                                      | `true`                                                                                  |
| `metrics.jmx.livenessProbe.initialDelaySeconds`                 | Initial delay seconds for livenessProbe                                                                                                                                                                                                   | `60`                                                                                    |
| `metrics.jmx.livenessProbe.periodSeconds`                       | Period seconds for livenessProbe                                                                                                                                                                                                          | `10`                                                                                    |
| `metrics.jmx.livenessProbe.timeoutSeconds`                      | Timeout seconds for livenessProbe                                                                                                                                                                                                         | `10`                                                                                    |
| `metrics.jmx.livenessProbe.failureThreshold`                    | Failure threshold for livenessProbe                                                                                                                                                                                                       | `3`                                                                                     |
| `metrics.jmx.livenessProbe.successThreshold`                    | Success threshold for livenessProbe                                                                                                                                                                                                       | `1`                                                                                     |
| `metrics.jmx.readinessProbe.enabled`                            | Enable readinessProbe                                                                                                                                                                                                                     | `true`                                                                                  |
| `metrics.jmx.readinessProbe.initialDelaySeconds`                | Initial delay seconds for readinessProbe                                                                                                                                                                                                  | `30`                                                                                    |
| `metrics.jmx.readinessProbe.periodSeconds`                      | Period seconds for readinessProbe                                                                                                                                                                                                         | `10`                                                                                    |
| `metrics.jmx.readinessProbe.timeoutSeconds`                     | Timeout seconds for readinessProbe                                                                                                                                                                                                        | `10`                                                                                    |
| `metrics.jmx.readinessProbe.failureThreshold`                   | Failure threshold for readinessProbe                                                                                                                                                                                                      | `3`                                                                                     |
| `metrics.jmx.readinessProbe.successThreshold`                   | Success threshold for readinessProbe                                                                                                                                                                                                      | `1`                                                                                     |
| `metrics.jmx.service.ports.metrics`                             | Prometheus JMX exporter metrics service port                                                                                                                                                                                              | `5556`                                                                                  |
| `metrics.jmx.service.clusterIP`                                 | Static clusterIP or None for headless services                                                                                                                                                                                            | `""`                                                                                    |
| `metrics.jmx.service.sessionAffinity`                           | Control where client requests go, to the same pod or round-robin                                                                                                                                                                          | `None`                                                                                  |
| `metrics.jmx.service.annotations`                               | Annotations for the Prometheus JMX exporter service                                                                                                                                                                                       | `{}`                                                                                    |
| `metrics.jmx.service.ipFamilies`                                | IP families for the jmx metrics service                                                                                                                                                                                                   | `[]`                                                                                    |
| `metrics.jmx.service.ipFamilyPolicy`                            | IP family policy for the jmx metrics service                                                                                                                                                                                              | `""`                                                                                    |
| `metrics.jmx.whitelistObjectNames`                              | Allows setting which JMX objects you want to expose to via JMX stats to JMX exporter                                                                                                                                                      | `["kafka.controller:*","kafka.server:*","java.lang:*","kafka.network:*","kafka.log:*"]` |
| `metrics.jmx.config`                                            | Configuration file for JMX exporter                                                                                                                                                                                                       | `""`                                                                                    |
| `metrics.jmx.existingConfigmap`                                 | Name of existing ConfigMap with JMX exporter configuration                                                                                                                                                                                | `""`                                                                                    |
| `metrics.jmx.extraRules`                                        | Add extra rules to JMX exporter configuration                                                                                                                                                                                             | `""`                                                                                    |
| `metrics.serviceMonitor.enabled`                                | if `true`, creates a Prometheus Operator ServiceMonitor (requires `metrics.jmx.enabled` to be `true`)                                                                                                                                     | `false`                                                                                 |
| `metrics.serviceMonitor.namespace`                              | Namespace in which Prometheus is running                                                                                                                                                                                                  | `""`                                                                                    |
| `metrics.serviceMonitor.path`                                   | Path where JMX exporter serves metrics                                                                                                                                                                                                    | `/metrics`                                                                              |
| `metrics.serviceMonitor.interval`                               | Interval at which metrics should be scraped                                                                                                                                                                                               | `""`                                                                                    |
| `metrics.serviceMonitor.scrapeTimeout`                          | Timeout after which the scrape is ended                                                                                                                                                                                                   | `""`                                                                                    |
| `metrics.serviceMonitor.labels`                                 | Additional labels that can be used so ServiceMonitor will be discovered by Prometheus                                                                                                                                                     | `{}`                                                                                    |
| `metrics.serviceMonitor.selector`                               | Prometheus instance selector labels                                                                                                                                                                                                       | `{}`                                                                                    |
| `metrics.serviceMonitor.relabelings`                            | RelabelConfigs to apply to samples before scraping                                                                                                                                                                                        | `[]`                                                                                    |
| `metrics.serviceMonitor.metricRelabelings`                      | MetricRelabelConfigs to apply to samples before ingestion                                                                                                                                                                                 | `[]`                                                                                    |
| `metrics.serviceMonitor.honorLabels`                            | Specify honorLabels parameter to add the scrape endpoint                                                                                                                                                                                  | `false`                                                                                 |
| `metrics.serviceMonitor.jobLabel`                               | The name of the label on the target service to use as the job name in prometheus.                                                                                                                                                         | `""`                                                                                    |
| `metrics.prometheusRule.enabled`                                | if `true`, creates a Prometheus Operator PrometheusRule (requires `metrics.jmx.enabled` to be `true`)                                                                                                                                     | `false`                                                                                 |
| `metrics.prometheusRule.namespace`                              | Namespace in which Prometheus is running                                                                                                                                                                                                  | `""`                                                                                    |
| `metrics.prometheusRule.labels`                                 | Additional labels that can be used so PrometheusRule will be discovered by Prometheus                                                                                                                                                     | `{}`                                                                                    |
| `metrics.prometheusRule.groups`                                 | Prometheus Rule Groups for Kafka                                                                                                                                                                                                          | `[]`                                                                                    |

### Kafka provisioning parameters

| Name                                                             | Description                                                                                                                                                                                                                                 | Value                 |
| ---------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| `provisioning.enabled`                                           | Enable Kafka provisioning Job                                                                                                                                                                                                               | `false`               |
| `provisioning.waitForKafka`                                      | Whether an init container should be created to wait until Kafka is ready before provisioning                                                                                                                                                | `true`                |
| `provisioning.useHelmHooks`                                      | Flag to indicate usage of helm hooks                                                                                                                                                                                                        | `true`                |
| `provisioning.automountServiceAccountToken`                      | Mount Service Account token in pod                                                                                                                                                                                                          | `false`               |
| `provisioning.numPartitions`                                     | Default number of partitions for topics when unspecified                                                                                                                                                                                    | `1`                   |
| `provisioning.replicationFactor`                                 | Default replication factor for topics when unspecified                                                                                                                                                                                      | `1`                   |
| `provisioning.topics`                                            | Kafka topics to provision                                                                                                                                                                                                                   | `[]`                  |
| `provisioning.nodeSelector`                                      | Node labels for pod assignment                                                                                                                                                                                                              | `{}`                  |
| `provisioning.tolerations`                                       | Tolerations for pod assignment                                                                                                                                                                                                              | `[]`                  |
| `provisioning.extraProvisioningCommands`                         | Extra commands to run to provision cluster resources                                                                                                                                                                                        | `[]`                  |
| `provisioning.parallel`                                          | Number of provisioning commands to run at the same time                                                                                                                                                                                     | `1`                   |
| `provisioning.preScript`                                         | Extra bash script to run before topic provisioning. /shared/client.properties is path to properties file with most needed configurations                                                                                                    | `""`                  |
| `provisioning.postScript`                                        | Extra bash script to run after topic provisioning. /shared/client.properties is path to properties file with most needed configurations                                                                                                     | `""`                  |
| `provisioning.auth.tls.type`                                     | Format to use for TLS certificates. Allowed types: `JKS` and `PEM`.                                                                                                                                                                         | `jks`                 |
| `provisioning.auth.tls.certificatesSecret`                       | Existing secret containing the TLS certificates for the Kafka provisioning Job.                                                                                                                                                             | `""`                  |
| `provisioning.auth.tls.cert`                                     | The secret key from the certificatesSecret if 'cert' key different from the default (tls.crt)                                                                                                                                               | `tls.crt`             |
| `provisioning.auth.tls.key`                                      | The secret key from the certificatesSecret if 'key' key different from the default (tls.key)                                                                                                                                                | `tls.key`             |
| `provisioning.auth.tls.caCert`                                   | The secret key from the certificatesSecret if 'caCert' key different from the default (ca.crt)                                                                                                                                              | `ca.crt`              |
| `provisioning.auth.tls.keystore`                                 | The secret key from the certificatesSecret if 'keystore' key different from the default (keystore.jks)                                                                                                                                      | `keystore.jks`        |
| `provisioning.auth.tls.truststore`                               | The secret key from the certificatesSecret if 'truststore' key different from the default (truststore.jks)                                                                                                                                  | `truststore.jks`      |
| `provisioning.auth.tls.passwordsSecret`                          | Name of the secret containing passwords to access the JKS files or PEM key when they are password-protected.                                                                                                                                | `""`                  |
| `provisioning.auth.tls.keyPasswordSecretKey`                     | The secret key from the passwordsSecret if 'keyPasswordSecretKey' key different from the default (key-password)                                                                                                                             | `key-password`        |
| `provisioning.auth.tls.keystorePasswordSecretKey`                | The secret key from the passwordsSecret if 'keystorePasswordSecretKey' key different from the default (keystore-password)                                                                                                                   | `keystore-password`   |
| `provisioning.auth.tls.truststorePasswordSecretKey`              | The secret key from the passwordsSecret if 'truststorePasswordSecretKey' key different from the default (truststore-password)                                                                                                               | `truststore-password` |
| `provisioning.auth.tls.keyPassword`                              | Password to access the password-protected PEM key if necessary. Ignored if 'passwordsSecret' is provided.                                                                                                                                   | `""`                  |
| `provisioning.auth.tls.keystorePassword`                         | Password to access the JKS keystore. Ignored if 'passwordsSecret' is provided.                                                                                                                                                              | `""`                  |
| `provisioning.auth.tls.truststorePassword`                       | Password to access the JKS truststore. Ignored if 'passwordsSecret' is provided.                                                                                                                                                            | `""`                  |
| `provisioning.command`                                           | Override provisioning container command                                                                                                                                                                                                     | `[]`                  |
| `provisioning.args`                                              | Override provisioning container arguments                                                                                                                                                                                                   | `[]`                  |
| `provisioning.extraEnvVars`                                      | Extra environment variables to add to the provisioning pod                                                                                                                                                                                  | `[]`                  |
| `provisioning.extraEnvVarsCM`                                    | ConfigMap with extra environment variables                                                                                                                                                                                                  | `""`                  |
| `provisioning.extraEnvVarsSecret`                                | Secret with extra environment variables                                                                                                                                                                                                     | `""`                  |
| `provisioning.podAnnotations`                                    | Extra annotations for Kafka provisioning pods                                                                                                                                                                                               | `{}`                  |
| `provisioning.podLabels`                                         | Extra labels for Kafka provisioning pods                                                                                                                                                                                                    | `{}`                  |
| `provisioning.serviceAccount.create`                             | Enable creation of ServiceAccount for Kafka provisioning pods                                                                                                                                                                               | `true`                |
| `provisioning.serviceAccount.name`                               | The name of the service account to use. If not set and `create` is `true`, a name is generated                                                                                                                                              | `""`                  |
| `provisioning.serviceAccount.automountServiceAccountToken`       | Allows auto mount of ServiceAccountToken on the serviceAccount created                                                                                                                                                                      | `false`               |
| `provisioning.resourcesPreset`                                   | Set container resources according to one common preset (allowed values: none, nano, micro, small, medium, large, xlarge, 2xlarge). This is ignored if provisioning.resources is set (provisioning.resources is recommended for production). | `micro`               |
| `provisioning.resources`                                         | Set container requests and limits for different resources like CPU or memory (essential for production workloads)                                                                                                                           | `{}`                  |
| `provisioning.podSecurityContext.enabled`                        | Enable security context for the pods                                                                                                                                                                                                        | `true`                |
| `provisioning.podSecurityContext.fsGroupChangePolicy`            | Set filesystem group change policy                                                                                                                                                                                                          | `Always`              |
| `provisioning.podSecurityContext.sysctls`                        | Set kernel settings using the sysctl interface                                                                                                                                                                                              | `[]`                  |
| `provisioning.podSecurityContext.supplementalGroups`             | Set filesystem extra groups                                                                                                                                                                                                                 | `[]`                  |
| `provisioning.podSecurityContext.fsGroup`                        | Set Kafka provisioning pod's Security Context fsGroup                                                                                                                                                                                       | `1001`                |
| `provisioning.podSecurityContext.seccompProfile.type`            | Set Kafka provisioning pod's Security Context seccomp profile                                                                                                                                                                               | `RuntimeDefault`      |
| `provisioning.containerSecurityContext.enabled`                  | Enable Kafka provisioning containers' Security Context                                                                                                                                                                                      | `true`                |
| `provisioning.containerSecurityContext.seLinuxOptions`           | Set SELinux options in container                                                                                                                                                                                                            | `{}`                  |
| `provisioning.containerSecurityContext.runAsUser`                | Set containers' Security Context runAsUser                                                                                                                                                                                                  | `1001`                |
| `provisioning.containerSecurityContext.runAsGroup`               | Set containers' Security Context runAsGroup                                                                                                                                                                                                 | `1001`                |
| `provisioning.containerSecurityContext.runAsNonRoot`             | Set Kafka provisioning containers' Security Context runAsNonRoot                                                                                                                                                                            | `true`                |
| `provisioning.containerSecurityContext.allowPrivilegeEscalation` | Set Kafka provisioning containers' Security Context allowPrivilegeEscalation                                                                                                                                                                | `false`               |
| `provisioning.containerSecurityContext.readOnlyRootFilesystem`   | Set Kafka provisioning containers' Security Context readOnlyRootFilesystem                                                                                                                                                                  | `true`                |
| `provisioning.containerSecurityContext.capabilities.drop`        | Set Kafka provisioning containers' Security Context capabilities to be dropped                                                                                                                                                              | `["ALL"]`             |
| `provisioning.schedulerName`                                     | Name of the k8s scheduler (other than default) for kafka provisioning                                                                                                                                                                       | `""`                  |
| `provisioning.enableServiceLinks`                                | Whether information about services should be injected into pod's environment variable                                                                                                                                                       | `true`                |
| `provisioning.extraVolumes`                                      | Optionally specify extra list of additional volumes for the Kafka provisioning pod(s)                                                                                                                                                       | `[]`                  |
| `provisioning.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts for the Kafka provisioning container(s)                                                                                                                                            | `[]`                  |
| `provisioning.sidecars`                                          | Add additional sidecar containers to the Kafka provisioning pod(s)                                                                                                                                                                          | `[]`                  |
| `provisioning.initContainers`                                    | Add additional Add init containers to the Kafka provisioning pod(s)                                                                                                                                                                         | `[]`                  |

```console
helm install my-release \
  --set controller.replicaCount=3 \
  oci://REGISTRY_NAME/REPOSITORY_NAME/kafka
```

> Note: You need to substitute the placeholders `REGISTRY_NAME` and `REPOSITORY_NAME` with a reference to your Helm chart registry and repository. For example, in the case of Bitnami, you need to use `REGISTRY_NAME=registry-1.docker.io` and `REPOSITORY_NAME=bitnamicharts`.

The above command deploys Kafka with 3 Kafka controller-eligible nodes.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```console
helm install my-release -f values.yaml oci://REGISTRY_NAME/REPOSITORY_NAME/kafka
```

> Note: You need to substitute the placeholders `REGISTRY_NAME` and `REPOSITORY_NAME` with a reference to your Helm chart registry and repository. For example, in the case of Bitnami, you need to use `REGISTRY_NAME=registry-1.docker.io` and `REPOSITORY_NAME=bitnamicharts`.
> **Tip**: You can use the default [values.yaml](https://github.com/bitnami/charts/tree/main/bitnami/kafka/values.yaml)

## Troubleshooting

Find more information about how to deal with common errors related to Bitnami's Helm charts in [this troubleshooting guide](https://docs.bitnami.com/general/how-to/troubleshoot-helm-chart-issues).

## Upgrading

### To 32.3.0

We have introduced the value `kraftVersion` to help control the change from static quorum to dynamic quorum.

By default, new clusters will be deployed with the new dynamic quorum (kraftVersion=1), but users upgrading from Kafka 3.x may need to modify the default values to continue using static quorum (kraftVersion=0).
That is because Kafka 4.0 does not yet support switching from static quorum (controller.quorum.voters) to dynamic quorum (controller.quorum.bootstrap.servers), causing upgrades to fail (#34015).

For more information please check [Kafka documentation](https://kafka.apache.org/documentation/#static_versus_dynamic_kraft_quorums).

### To 32.0.0

This major release bumps Kafka major version to `4.y.z` series. This version implies a significant milestone given now Kafka operates operate entirely without Apache ZooKeeper, running in KRaft mode by default. As a consequence, **ZooKeeper is no longer a chart dependency and every related parameter has been removed.**. Upgrading from `31.y.z` chart version is not supported unless KRaft mode was already enabled.

Also, some KRaft-related parameters have been renamed or removed:

- `kraft.enabled` has been removed. Kafka now operates in KRaft mode by default.
- `kraft.controllerQuorumVoters` has been renamed to `controller.quorumVoters`.
- `kraft.clusterId` and `kraft.existingClusterIdSecret` have been renamed to `clusterId` and `existingKraftSecret`, respectively.

Other notable changes:

- `log4j` and `existingLog4jConfig` parameters have been renamed to `log4j2` and `existingLog4j2ConfigMap`, respectively.
- `controller.quorumVoters` has been removed in favor of `controller.quorumBootstrapServers`.
- `brokerRackAssignment` and `brokerRackAssignmentApiVersion` are deprecated in favor of `brokerRackAwareness.*` parameters.
- `tls.autoGenerated` boolean is now an object with extended configuration options.
- `volumePermissions` parameters have been moved under `defaultInitContainers` parameter.
- `externalAccess.autoDiscovery` parameters have been moved under `defaultInitContainers` parameter.
- `controller.initContainerResources` and `broker.initContainerResources` have been removed. Use `defaultInitContainers.prepareConfig.resources` instead.
- `extraInit` has been renamed to `defaultInitContainers.prepareConfig.extraInit`.
- `extraConfig` has been replaced with `overrideConfiguration`.
- provisioning variable `$CLIENT_CONF` has been removed use `/shared/client.properties`.
- If using your own SSL certs you will need to disable `tls.autoGenerated.enabled`

### To 31.1.0

This version introduces image verification for security purposes. To disable it, set `global.security.allowInsecureImages` to `true`. More details at [GitHub issue](https://github.com/bitnami/charts/issues/30850).

### To 31.0.0

This major release bumps the Kafka version to 3.9. Find notable changes in [kafka upgrade notes](https://kafka.apache.org/39/documentation.html#upgrade).

### To 30.0.0

This major release bumps the Kafka version to 3.8. Find notable changes in [kafka upgrade notes](https://kafka.apache.org/38/documentation.html#upgrade).

### To 29.0.0

This major version of Kafka deprecates Kafka Exporter component.

### To 28.0.0

This major bump changes the following security defaults:

- `runAsGroup` is changed from `0` to `1001`
- `readOnlyRootFilesystem` is set to `true`
- `resourcesPreset` is changed from `none` to the minimum size working in our test suites (NOTE: `resourcesPreset` is not meant for production usage, but `resources` adapted to your use case).
- `global.compatibility.openshift.adaptSecurityContext` is changed from `disabled` to `auto`.
- The `networkPolicy` section has been normalized amongst all Bitnami charts. Compared to the previous approach, the values section has been simplified (check the Parameters section) and now it set to `enabled=true` by default. Egress traffic is allowed by default and ingress traffic is allowed by all pods but only to the ports set in `containerPorts` and `extraContainerPorts`.

This could potentially break any customization or init scripts used in your deployment. If this is the case, change the default values to the previous ones.

### To 26.0.0

This major release bumps the Kafka version to 3.6 [kafka upgrade notes](https://kafka.apache.org/36/documentation.html#upgrade).

### To 25.0.0

This major updates the Zookeeper subchart to it newest major, 12.0.0. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-1200).

### To 24.0.0

This major version is a refactor of the Kafka chart and its architecture, to better adapt to Kraft features introduced in version 22.0.0.

The changes introduced in this version are:

- New architecture. The chart now has two statefulsets, one for controller-eligible nodes (controller or controller+broker) and another one for broker-only nodes. Please take a look at the subsections [Upgrading from Kraft mode](#upgrading-from-kraft-mode) and [Upgrading from Zookeeper mode](#upgrading-from-zookeeper-mode) for more information about how to upgrade this chart depending on which mode you were using.

  The new architecture is designed to support two main features:
  - Deployment of dedicated nodes
  - Support for Zookeeper to Kraft migration

- Adds compatibility with `securityContext.readOnlyRootFs=true`, which is now the execution default.
  - The Kafka configuration is now mounted as a ConfigMap instead of generated at runtime.
  - Due to the implementation of readOnlyRootFs support, the following settings have been removed and will now rely on Kafka defaults. To override them, please use `extraConfig` to extend your Kafka configuration instead.
    - `deleteTopicEnable`
    - `autoCreateTopicsEnable`
    - `logFlushIntervalMessages`
    - `logFlushIntervalMs`
    - `logRetentionBytes`
    - `logRetentionCheckIntervalMs`
    - `logRetentionHours`
    - `logSegmentBytes`
    - `logsDirs`
    - `maxMessageBytes`
    - `defaultReplicationFactor`
    - `offsetsTopicReplicationFactor`
    - `transactionStateLogReplicationFactor`
    - `transactionStateLogMinIsr`
    - `numIoThreads`
    - `numNetworkThreads`
    - `numPartitions`
    - `numRecoveryThreadsPerDataDir`
    - `socketReceiveBufferBytes`
    - `socketRequestMaxBytes`
    - `socketSendBufferBytes`
    - `zookeeperConnectionTimeoutMs`
    - `authorizerClassName`
    - `allowEveryoneIfNoAclFound`
    - `superUsers`
- All listeners are configured with protocol 'SASL_PLAINTEXT' by default.
- Support for SCRAM authentication in KRaft mode
- All statefulset settings have been moved from values' root to `controller.*` and `broker.*`.
- Refactor of listeners configuration:
  - Settings `listeners`, `advertisedListeners` and `listenerSecurityProtocolMap` have been replaced with `listeners.*` object, which includes default listeners and each listener can be configured individually and extended using `listeners.extraListeners`.
  - Values `interBrokerListenerName`, `allowPlaintextListener` have been removed.
- Refactor of SASL, SSL and ACL settings:
  - Authentication nomenclature `plaintext,tls,mtls,sasl,sasl_tls` has been removed. Listeners are now configured using Kafka nomenclature `PLAINTEXT,SASL_PLAINTEXT,SASL_SSL,SSL` in `listeners.*.protocol`.
  - mTLS is configured by default for SSL protocol listeners, while it can now also be configured for SASL_SSL listeners if `listener.*.sslClientAuth` is set.
  - All SASL settings are now grouped under `sasl.*`.
    - `auth.sasl.mechanisms` -> `sasl.enabledMechanisms`
    - `auth.interBrokerMechanism` -> `sasl.interBrokerMechanism`
    - `auth.sasl.jaas.clientUSers` -> `sasl.client.users`
    - `auth.sasl.jaas.clientPasswords` -> `sasl.client.passwords`
    - `auth.sasl.jaas.interBrokerUser` -> `sasl.interbroker.user`
    - `auth.sasl.jaas.interBrokerPassword` -> `sasl.interbroker.password`
    - `auth.sasl.jaas.zookeeperUser` -> `sasl.zookeeper.user`
    - `auth.sasl.jaas.zookeeperPassword` -> `sasl.zookeeper.password`
    - `auth.sasl.jaas.existingSecret` -> `sasl.existingSecret`
  - Added support for Controller listener protocols other than PLAINTEXT.
  - TLS settings have been moved from `auth.tls.*` to `tls.*`.
  - Zookeeper TLS settings have been moved from `auth.zookeeper*` to `tls.zookeeper.*`
- Refactor externalAccess to support the new architecture:
  - `externalAccess.service.*` have been renamed to `externalAccess.controller.service.*` and `externalAccess.broker.service.*`.
  - Controller pods will not configure externalAccess unless either:
    - `controller.controllerOnly=false` (default), meaning the pods are running as 'controller+broker' nodes; or
    - `externalAccess.controller.service.forceExpose=true`, for use cases where controller-only nodes want to be exposed externally.
- TLS certificates value `tls.existingSecret` no longer supports an array of secrets (1 secret per node). It now accepts a single secret containing multiple certificates named `kafka-<role>-<pod-number>` for each Kafka pod, or alternatively, a single certificate shared by all Kafka nodes using wildcard CN and/or SubjectAltNames. **NOTE**: If using CertManager to automatically generate the certificate secrets, only the single certificate approach would be supported.

#### Upgrading from Kraft mode

If upgrading from Kraft mode, existing PVCs from Kafka containers should be reattached to 'controller' pods.

#### Upgrading from Zookeeper mode

If upgrading from Zookeeper mode, make sure you set 'controller.replicaCount=0' and reattach the existing PVCs to 'broker' pods.
This will allow you to perform a migration to Kraft mode in the future by following the following section.

##### Migrating from Zookeeper (Early access)

This guide is an adaptation from upstream documentation: [Migrate from ZooKeeper to KRaft](https://docs.confluent.io/platform/current/installation/migrate-zk-kraft.html)

1. Retrieve the cluster ID from Zookeeper:

    ```console
    $ kubectl exec -it <your-zookeeper-pod> -- zkCli.sh get /cluster/id
    /opt/bitnami/java/bin/java
    Connecting to localhost:2181

    WATCHER::

    WatchedEvent state:SyncConnected type:None path:null
    {"version":"1","id":"TEr3HVPvTqSWixWRHngP5g"}
    ```

2. Deploy at least one Kraft controller-only in your deployment and enable `zookeeperMigrationMode=true`. The Kraft controllers will migrate the data from your Kafka ZkBroker to Kraft mode.

    To do so add the following values to your Zookeeper deployment when upgrading:

    ```yaml
    controller:
      replicaCount: 1
      controllerOnly: true
      zookeeperMigrationMode: true
      # If needed, set controllers minID to avoid conflict with your ZK brokers' ids.
      # minID: 0
    broker:
      zookeeperMigrationMode: true
    kraft:
      enabled: true
      clusterId: "<your_cluster_id>"
    ```

3. Wait until until all brokers are ready. You should see the following log in the broker logs:

    ```console
    INFO [KafkaServer id=100] Finished catching up on KRaft metadata log, requesting that the KRaft controller unfence this broker (kafka.server.KafkaServer)
    INFO [BrokerLifecycleManager id=100 isZkBroker=true] The broker has been unfenced. Transitioning from RECOVERY to RUNNING. (kafka.server.BrokerLifecycleManager)
    ```

    In the controllers, the following message should show up:

    ```console
    Transitioning ZK migration state from PRE_MIGRATION to MIGRATION (org.apache.kafka.controller.FeatureControlManager)
    ```

4. Once all brokers have been successfully migrated, set `broker.zookeeperMigrationMode=false` to fully migrate them.

    ```yaml
    broker:
      zookeeperMigrationMode: false
    ```

5. To conclude the migration, switch off migration mode on controllers and stop Zookeeper:

    ```yaml
    controller:
      zookeeperMigrationMode: false
    zookeeper:
      enabled: false
    ```

    After migration is complete, you should see the following message in your controllers:

    ```console
    [2023-07-13 13:07:45,226] INFO [QuorumController id=1] Transitioning ZK migration state from MIGRATION to POST_MIGRATION (org.apache.kafka.controller.FeatureControlManager)
    ```

6. (**Optional**) If you would like to switch to a non-dedicated cluster, set `controller.controllerOnly=false`. This will cause controller-only nodes to switch to controller+broker nodes.

    At that point, you could manually decommission broker-only nodes by reassigning its partitions to controller-eligible nodes.

    For more information about decommissioning kafka broker check the [Kafka documentation](https://www.confluent.io/blog/remove-kafka-brokers-from-any-cluster-the-easy-way/).

#### Retaining PersistentVolumes

When upgrading the Kafka chart, you may want to retain your existing data. To do so, we recommend following this guide:

**NOTE**: This guide requires the binaries 'kubectl' and 'jq'.

```console
# Env variables
REPLICA=0
OLD_PVC="data-<your_release_name>-kafka-${REPLICA}"
NEW_PVC="data-<your_release_name>-kafka-<your_kafka_role>-${REPLICA}"
PV_NAME=$(kubectl get pvc $OLD_PVC -o jsonpath="{.spec.volumeName}")
NEW_PVC_MANIFEST_FILE="$NEW_PVC.yaml"

# Modify PV reclaim policy
kubectl patch pv $PV_NAME -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
# Manually check field 'RECLAIM POLICY'
kubectl get pv $PV_NAME

# Create new PVC manifest
kubectl get pvc $OLD_PVC -o json | jq "
  .metadata.name = \"$NEW_PVC\"
  | with_entries(
      select([.key] |
        inside([\"metadata\", \"spec\", \"apiVersion\", \"kind\"]))
    )
  | del(
      .metadata.annotations, .metadata.creationTimestamp,
      .metadata.finalizers, .metadata.resourceVersion,
      .metadata.selfLink, .metadata.uid
    )
  " > $NEW_PVC_MANIFEST_FILE
# Check manifest
cat $NEW_PVC_MANIFEST_FILE

# Delete your old Statefulset and PVC
kubectl delete sts "<your_release_name>-kafka"
kubectl delete pvc $OLD_PVC
# Make PV available again and create the new PVC
kubectl patch pv $PV_NAME -p '{"spec":{"claimRef": null}}'
kubectl apply -f $NEW_PVC_MANIFEST_FILE
```

Repeat this process for each replica you had in your Kafka cluster. Once completed, upgrade the cluster and the new Statefulset should reuse the existing PVCs.

### To 23.0.0

This major updates Kafka to its newest version, 3.5.x. For more information, please refer to [kafka upgrade notes](https://kafka.apache.org/35/documentation.html#upgrade).

### To 22.0.0

This major updates the Kafka's configuration to use Kraft by default. You can learn more about this configuration [here](https://developer.confluent.io/learn/kraft). Apart from seting the `kraft.enabled` parameter to `true`, we also made the following changes:

- Renamed `minBrokerId` parameter to `minId` to set the minimum ID to use when configuring the node.id or broker.id parameter depending on the Kafka's configuration. This parameter sets the `KAFKA_CFG_NODE_ID` env var in the container.
- Updated the `containerPorts` and `service.ports` parameters to include the new controller port.

### To 21.0.0

This major updates Kafka to its newest version, 3.4.x. For more information, please refer to [kafka upgrade notes](https://kafka.apache.org/34/documentation.html#upgrade).

### To 20.0.0

This major updates the Zookeeper subchart to it newest major, 11.0.0. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-1100).

### To 19.0.0

This major updates Kafka to its newest version, 3.3.x. For more information, please refer to [kafka upgrade notes](https://kafka.apache.org/33/documentation.html#upgrade).

### To 18.0.0

This major updates the Zookeeper subchart to it newest major, 10.0.0. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-1000).

### To 16.0.0

This major updates the Zookeeper subchart to it newest major, 9.0.0. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-900).

### To 15.0.0

This major release bumps Kafka major version to `3.x` series.
It also renames several values in this chart and adds missing features, in order to be inline with the rest of assets in the Bitnami charts repository. Some affected values are:

- `service.port`, `service.internalPort` and `service.externalPort` have been regrouped under the `service.ports` map.
- `metrics.kafka.service.port` has been regrouped under the `metrics.kafka.service.ports` map.
- `metrics.jmx.service.port` has been regrouped under the `metrics.jmx.service.ports` map.
- `updateStrategy` (string) and `rollingUpdatePartition` are regrouped under the `updateStrategy` map.
- Several parameters marked as deprecated `14.x.x` are not supported anymore.

Additionally updates the ZooKeeper subchart to it newest major, `8.0.0`, which contains similar changes.

### To 14.0.0

In this version, the `image` block is defined once and is used in the different templates, while in the previous version, the `image` block was duplicated for the main container and the provisioning one

```yaml
image:
  registry: docker.io
  repository: bitnami/kafka
  tag: 2.8.0
```

VS

```yaml
image:
  registry: docker.io
  repository: bitnami/kafka
  tag: 2.8.0
...
provisioning:
  image:
    registry: docker.io
    repository: bitnami/kafka
    tag: 2.8.0
```

See [PR#7114](https://github.com/bitnami/charts/pull/7114) for more info about the implemented changes

### To 13.0.0

This major updates the Zookeeper subchart to it newest major, 7.0.0, which renames all TLS-related settings. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-700).

### To 12.2.0

This version also introduces `bitnami/common`, a [library chart](https://helm.sh/docs/topics/library_charts/#helm) as a dependency. More documentation about this new utility could be found [here](https://github.com/bitnami/charts/tree/main/bitnami/common#bitnami-common-library-chart). Please, make sure that you have updated the chart dependencies before executing any upgrade.

### To 12.0.0

[On November 13, 2020, Helm v2 support was formally finished](https://github.com/helm/charts#status-of-the-project), this major version is the result of the required changes applied to the Helm Chart to be able to incorporate the different features added in Helm v3 and to be consistent with the Helm project itself regarding the Helm v2 EOL.

#### What changes were introduced in this major version?

- Previous versions of this Helm Chart use `apiVersion: v1` (installable by both Helm 2 and 3), this Helm Chart was updated to `apiVersion: v2` (installable by Helm 3 only). [Here](https://helm.sh/docs/topics/charts/#the-apiversion-field) you can find more information about the `apiVersion` field.
- Move dependency information from the *requirements.yaml* to the *Chart.yaml*
- After running `helm dependency update`, a *Chart.lock* file is generated containing the same structure used in the previous *requirements.lock*
- The different fields present in the *Chart.yaml* file has been ordered alphabetically in a homogeneous way for all the Bitnami Helm Charts

#### Considerations when upgrading to this version

- If you want to upgrade to this version from a previous one installed with Helm v3, you shouldn't face any issues
- If you want to upgrade to this version using Helm v2, this scenario is not supported as this version doesn't support Helm v2 anymore
- If you installed the previous version with Helm v2 and wants to upgrade to this version with Helm v3, please refer to the [official Helm documentation](https://helm.sh/docs/topics/v2_v3_migration/#migration-use-cases) about migrating from Helm v2 to v3

#### Useful links

- <https://techdocs.broadcom.com/us/en/vmware-tanzu/application-catalog/tanzu-application-catalog/services/tac-doc/apps-tutorials-resolve-helm2-helm3-post-migration-issues-index.html>
- <https://helm.sh/docs/topics/v2_v3_migration/>
- <https://helm.sh/blog/migrate-from-helm-v2-to-helm-v3/>

### To 11.8.0

External access to brokers can now be achieved through the cluster's Kafka service.

- `service.nodePort` -> deprecated  in favor of `service.nodePorts.client` and `service.nodePorts.external`

### To 11.7.0

The way to configure the users and passwords changed. Now it is allowed to create multiple users during the installation by providing the list of users and passwords.

- `auth.jaas.clientUser` (string) -> deprecated  in favor of `auth.jaas.clientUsers` (array).
- `auth.jaas.clientPassword` (string) -> deprecated  in favor of `auth.jaas.clientPasswords` (array).

### To 11.0.0

The way to configure listeners and athentication on Kafka is totally refactored allowing users to configure different authentication protocols on different listeners. Please check the [Listeners Configuration](#listeners-configuration) section for more information.

Backwards compatibility is not guaranteed you adapt your values.yaml to the new format. Here you can find some parameters that were renamed or disappeared in favor of new ones on this major version:

- `auth.enabled` -> deprecated in favor of `auth.clientProtocol` and `auth.interBrokerProtocol` parameters.
- `auth.ssl` -> deprecated in favor of `auth.clientProtocol` and `auth.interBrokerProtocol` parameters.
- `auth.certificatesSecret` -> renamed to `auth.jksSecret`.
- `auth.certificatesPassword` -> renamed to `auth.jksPassword`.
- `sslEndpointIdentificationAlgorithm` -> renamedo to `auth.tlsEndpointIdentificationAlgorithm`.
- `auth.interBrokerUser` -> renamed to `auth.jaas.interBrokerUser`
- `auth.interBrokerPassword` -> renamed to `auth.jaas.interBrokerPassword`
- `auth.zookeeperUser` -> renamed to `auth.jaas.zookeeperUser`
- `auth.zookeeperPassword` -> renamed to `auth.jaas.zookeeperPassword`
- `auth.existingSecret` -> renamed to `auth.jaas.existingSecret`
- `service.sslPort` -> deprecated in favor of `service.internalPort`
- `service.nodePorts.kafka` and `service.nodePorts.ssl` -> deprecated in favor of `service.nodePort`
- `metrics.kafka.extraFlag` -> new parameter
- `metrics.kafka.certificatesSecret` -> new parameter

### To 10.0.0

If you are setting the `config` or `log4j` parameter, backwards compatibility is not guaranteed, because the `KAFKA_MOUNTED_CONFDIR` has moved from `/opt/bitnami/kafka/conf` to `/bitnami/kafka/config`. In order to continue using these parameters, you must also upgrade your image to `docker.io/bitnami/kafka:2.4.1-debian-10-r38` or later.

### To 9.0.0

Backwards compatibility is not guaranteed you adapt your values.yaml to the new format. Here you can find some parameters that were renamed on this major version:

```diff
- securityContext.enabled
- securityContext.fsGroup
- securityContext.fsGroup
+ podSecurityContext
- externalAccess.service.loadBalancerIP
+ externalAccess.service.loadBalancerIPs
- externalAccess.service.nodePort
+ externalAccess.service.nodePorts
- metrics.jmx.configMap.enabled
- metrics.jmx.configMap.overrideConfig
+ metrics.jmx.config
- metrics.jmx.configMap.overrideName
+ metrics.jmx.existingConfigmap
```

Ports names were prefixed with the protocol to comply with Istio (see <https://istio.io/docs/ops/deployment/requirements/>).

### To 8.0.0

There is not backwards compatibility since the brokerID changes to the POD_NAME. For more information see [this PR](https://github.com/bitnami/charts/pull/2028).

### To 7.0.0

Backwards compatibility is not guaranteed when Kafka metrics are enabled, unless you modify the labels used on the exporter deployments.
Use the workaround below to upgrade from versions previous to 7.0.0. The following example assumes that the release name is kafka:

```console
helm upgrade kafka oci://REGISTRY_NAME/REPOSITORY_NAME/kafka --version 6.1.8 --set metrics.kafka.enabled=false
helm upgrade kafka oci://REGISTRY_NAME/REPOSITORY_NAME/kafka --version 7.0.0 --set metrics.kafka.enabled=true
```

> Note: You need to substitute the placeholders `REGISTRY_NAME` and `REPOSITORY_NAME` with a reference to your Helm chart registry and repository. For example, in the case of Bitnami, you need to use `REGISTRY_NAME=registry-1.docker.io` and `REPOSITORY_NAME=bitnamicharts`.

### To 2.0.0

Backwards compatibility is not guaranteed unless you modify the labels used on the chart's deployments.
Use the workaround below to upgrade from versions previous to 2.0.0. The following example assumes that the release name is kafka:

```console
kubectl delete statefulset kafka-kafka --cascade=false
kubectl delete statefulset kafka-zookeeper --cascade=false
```

### To 1.0.0

Backwards compatibility is not guaranteed unless you modify the labels used on the chart's deployments.
Use the workaround below to upgrade from versions previous to 1.0.0. The following example assumes that the release name is kafka:

```console
kubectl delete statefulset kafka-kafka --cascade=false
kubectl delete statefulset kafka-zookeeper --cascade=false
```

## License

Copyright &copy; 2025 Broadcom. The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
