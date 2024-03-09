---
layout: tutorials
title: Understanding Kubernetes ConfigMaps and Secrets with Meshery Playground
abstract: Learn how to work with ConfigMaps and Secrets in Kubernetes using Meshery Playground
permalink: guides/tutorials/kubernetes-configmaps-secrets
redirect_from: guides/tutorials/kubernetes-configmaps-secrets/
model: kubernetes
kind: configmaps
type: tutorials
category: tutorial
language: en
list: include
published: false
abstract: "In this tutorial, we will explore how to effectively use Kubernetes ConfigMaps and Secrets for managing configuration data and sensitive information. Leveraging Meshery Playground, an interactive live cluster environment, we'll perform hands-on labs to understand the practical aspects of working with ConfigMaps and Secrets in Kubernetes."
---

Introduction:
In this tutorial, we will delve into the realm of Kubernetes ConfigMaps and Secrets. ConfigMaps are used to manage configuration data, while Secrets handle sensitive information. Using Meshery Playground, an interactive live cluster environment, we'll perform hands-on labs to gain practical insights into working with ConfigMaps and Secrets in Kubernetes.

Prerequisites:
- Basic understanding of Kubernetes concepts.
- Meshery Playground access. If you don't have an account, sign up at [Meshery Playground](https://meshery.layer5.io/play).

Lab Scenario: Managing Application Configuration with ConfigMaps and Secrets

Objective:
Learn how to create and manage ConfigMaps and Secrets in Kubernetes to handle configuration data and sensitive information for your applications.

### Steps:

#### 1. **Accessing Meshery Playground:**
   - Log in to the [Meshery Playground](https://meshery.layer5.io/play) using your credentials.
   - Navigate to the Meshery Playground dashboard.

#### 2. **Creating a ConfigMap:**
   - Learn how to create a ConfigMap to store configuration data for your application.

\```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-configmap
data:
  app.properties: |
    key1=value1
    key2=value2
\```

Apply the ConfigMap configuration:

\```bash
kubectl apply -f configmap.yaml
\```

#### 3. **Using ConfigMaps in Pods:**
   - Explore how to use ConfigMaps in your Pod specifications to inject configuration data into your application.

\```yaml
# pod-with-configmap.yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
  - name: mycontainer
    image: nginx:latest
    volumeMounts:
    - name: config-volume
      mountPath: /etc/config
  volumes:
  - name: config-volume
    configMap:
      name: my-configmap
\```

Apply the Pod configuration:

\```bash
kubectl apply -f pod-with-configmap.yaml
\```

#### 4. **Creating a Secret:**
   - Learn how to create a Secret to store sensitive information, such as API keys or database credentials.

\```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
data:
  username: YWRtaW4=  # Base64 encoded username
  password: cGFzc3dvcmQ=  # Base64 encoded password
\```

Apply the Secret configuration:

\```bash
kubectl apply -f secret.yaml
\```

#### 5. **Using Secrets in Pods:**
   - Explore how to use Secrets in your Pod specifications to securely pass sensitive information to your application.

\```yaml
# pod-with-secret.yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod-secret
spec:
  containers:
  - name: mycontainer-secret
    image: nginx:latest
    env:
    - name: DB_USERNAME
      valueFrom:
        secretKeyRef:
          name: my-secret
          key: username
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: my-secret
          key: password
\```

Apply the Pod configuration:

\```bash
kubectl apply -f pod-with-secret.yaml
\```

#### 6. **Updating ConfigMaps and Secrets:**
   - Learn how to update ConfigMaps and Secrets dynamically and observe the changes in the associated Pods.

\```bash
kubectl edit configmap my-configmap
kubectl edit secret my-secret
\```

#### 7. **Clean-Up:**
   - Delete the ConfigMap, Secret, and associated resources after completing the lab.

\```bash
kubectl delete configmap my-configmap
kubectl delete secret my-secret
kubectl delete pod mypod mypod-secret
\```

#### 8. **Saving and Sharing:**
   - Save your scenario in Meshery Playground for future reference.
   - Share your ConfigMap and Secret scenarios with the Meshery community for collaborative learning.

### Conclusion
Congratulations! You've successfully completed the lab on Understanding Kubernetes ConfigMaps and Secrets using Meshery Playground. This hands-on experience has equipped you with practical knowledge on managing configuration data and sensitive information in Kubernetes. Continue exploring more scenarios in the Meshery Playground to enhance your skills in cloud-native technologies.
