---
layout: default
title: Deploying Apache Cassandra with a StatefulSet 
abstract: |
  This tutorial shows you how to run Apache Cassandra on Kubernetes. Cassandra, a database, needs persistent storage to provide data durability (application state). In this example, a custom Cassandra seed provider lets the database discover new Cassandra instances as they join the Cassandra cluster.
permalink: guides/tutorials/Deploy-Cassandra-with-StatefulSet
redirect_from: guides/tutorials/deploy-cassandra-with-statefulset
model: kubernetes
type: guides
category: tutorials
language: en
list: include
---
### Introduction

This tutorial demonstrates how to run Apache Cassandra on Kubernetes. Cassandra is a highly scalable and distributed database designed to handle large amounts of data across many commodity servers, providing high availability with no single point of failure. A key requirement for Cassandra is persistent storage to ensure data durability and application state continuity.

In this example, a custom Cassandra seed provider is used to enable the database to discover new Cassandra instances as they join the cluster, ensuring seamless scalability and integration.

StatefulSets in Kubernetes facilitate the deployment and management of stateful applications by maintaining a unique identity for each pod and stable storage. This tutorial leverages the StatefulSet feature to simplify the deployment of Cassandra. For more detailed information about StatefulSets and their benefits, please refer to the [StatefulSet documentation](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/).

### Prerequisites:
- Basic understanding of containerization , Kubernetes concepts and basic knowledge of Apache Cassandra.
- Access to the _Meshery Playground_. If you don't have an account, sign up at [Meshery Playground](https://meshery.layer5.io/).

### Lab Scenario:
In this lab, you will learn how to deploy Apache Cassandra on a Kubernetes cluster using StatefulSets. The scenario will guide you through setting up the necessary resources and configurations to run Cassandra with persistent storage, ensuring data durability and scalability. By the end of this lab, you will have a working Cassandra cluster capable of automatically discovering and integrating new instances.

### Learning Objectives
You will be able to deploy Apache Cassandra in Meshery Playground using StatefulSets. This lab will provide you with a comprehensive understanding of StatefulSets and their application in managing stateful workloads in Kubernetes.

### Steps:

#### Access Meshery Playground
   - Log in to the [Meshery Playground](https://meshery.layer5.io/) using your credentials. On successful login, you should be at the dashboard. Press the **X** on the _Where do you want to start?_ popup to close it (if required).
   - Click **Explore** in the Cloud Native Playground tile to navigate to _MeshMap_.

> **_NOTE:_**  MeshMap is still in beta.

#### Step 1 -Creating a headless Service for Cassandra

In Kubernetes, a Service describes a set of Pods that perform the same task.
The following Service is used for DNS lookups between Cassandra Pods and clients within your cluster:
```\yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: cassandra
  name: cassandra
spec:
  clusterIP: None
  ports:
  - port: 9042
  selector:
    app: cassandra
```
1.Import this file into meshery playground.

2.once succesfully imported make sure sure to deploy in meshery playground.

#### Step2 - Using a StatefulSet to create a Cassandra ring 

The StatefulSet manifest, included below, creates a Cassandra ring that consists of three Pods.
### Note:
This example uses the default provisioner for Minikube. Please update the following StatefulSet for the Kubernetes cluster you are working with in meshery playground.
### StatefullSet  manifest :
```\yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: cassandra
  labels:
    app: cassandra
spec:
  serviceName: cassandra
  replicas: 3
  selector:
    matchLabels:
      app: cassandra
  template:
    metadata:
      labels:
        app: cassandra
    spec:
      terminationGracePeriodSeconds: 1800
      containers:
      - name: cassandra
        image: gcr.io/google-samples/cassandra:v13
        imagePullPolicy: Always
        ports:
        - containerPort: 7000
          name: intra-node
        - containerPort: 7001
          name: tls-intra-node
        - containerPort: 7199
          name: jmx
        - containerPort: 9042
          name: cql
        resources:
          limits:
            cpu: "500m"
            memory: 1Gi
          requests:
            cpu: "500m"
            memory: 1Gi
        securityContext:
          capabilities:
            add:
              - IPC_LOCK
        lifecycle:
          preStop:
            exec:
              command: 
              - /bin/sh
              - -c
              - nodetool drain
        env:
          - name: MAX_HEAP_SIZE
            value: 512M
          - name: HEAP_NEWSIZE
            value: 100M
          - name: CASSANDRA_SEEDS
            value: "cassandra-0.cassandra.default.svc.cluster.local"
          - name: CASSANDRA_CLUSTER_NAME
            value: "K8Demo"
          - name: CASSANDRA_DC
            value: "DC1-K8Demo"
          - name: CASSANDRA_RACK
            value: "Rack1-K8Demo"
          - name: POD_IP
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
        readinessProbe:
          exec:
            command:
            - /bin/bash
            - -c
            - /ready-probe.sh
          initialDelaySeconds: 15
          timeoutSeconds: 5
        # These volume mounts are persistent. They are like inline claims,
        # but not exactly because the names need to match exactly one of
        # the stateful pod volumes.
        volumeMounts:
        - name: cassandra-data
          mountPath: /cassandra_data
  # These are converted to volume claims by the controller
  # and mounted at the paths mentioned above.
  # do not use these in production until ssd GCEPersistentDisk or other ssd pd
  volumeClaimTemplates:
  - metadata:
      name: cassandra-data
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: fast
      resources:
        requests:
          storage: 1Gi
---
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: fast
provisioner: k8s.io/minikube-hostpath
parameters:
  type: pd-ssd
```
1. Import this file into meshery playground .
2. once succesfully imported make sure sure to deploy in meshery playground.
3. merge both service (which is done in 1st step) and StatefulSet designs with help of merge designs feature in meshery playground.

#### Step3 -Modifying the Cassandra StatefulSet

1.Click the StatefullSet component to load the configuration modal.
2.A StatefullSet consits of replicas. modify the replicas according to the need.
3.for an example scenraio ,let's increase replicas to 5 in statefullset .
4.Save the design and Click on Deploy the design again to apply the modified  configuration.

#### Step4 - Vizualizing  the Cassandra StatefulSet

To view the resources deployed we will use the **Visualize** section of the _MeshMap_. A view is created with necessary filters to show the relevant resources.

1.  Click **Visualize** to begin.
2.  Give the view a name (rename).
3.  Click the filter icon.
4.  Choose appropriate filters to limited displayed resources in the view. For example, here we want to display StatefullSets,Service,StorageClass etc  
    Additionally, we will also add a label filter i.e. `tutorial=deployment` in this case. This should show a filtered view with only your resources something 
    similar to the screenshot below:

### Viewing Statefullset and Service information

Select the Statefullset from the _View_ to load the Statefullset details to the right. Ensure the _Details_ tab is selected.

Now, select one of the Service to display the service details. 

####  Step 5-Deleting the Deployment

To delete the deployment, use the **Undeploy** option from the _Design_ view.
#### Note: 
Deleting or scaling a StatefulSet down does not delete the volumes associated with the StatefulSet. This setting is for your safety because your data is more valuable than automatically purging all related StatefulSet resources.


#### Conclusion

This tutorial demonstrated deploying Apache Cassandra on Kubernetes using StatefulSets, ensuring data durability and seamless scalability. You learned to manage Cassandra instances with persistent storage and automatic discovery, leveraging Kubernetes' StatefulSet feature for robust application state management and high availability.


