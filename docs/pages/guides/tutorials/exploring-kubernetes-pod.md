---
layout: tutorials
title: Exploring Kubernetes Pods with Meshery Playground
abstract: Learn Kubernetes Pods with Meshery Playground
permalink: guides/tutorials/kubernetes-pods
redirect_from: guides/tutorials/kubernetes-pods/
kind: pods
model: kubernetes
type: tutorials
category: tutorial
language: en
list: include
abstract: "In this tutorial, we will explore Kubernetes Pods using Meshery Playground, an interactive live cluster environment, to perform hands-on labs for deploying and managing containerized applications at the Pod level."
---

Introduction:
In this tutorial, we'll dive into the fundamentals of Kubernetes Pods, the smallest deployable units in the Kubernetes ecosystem. Using Meshery Playground, an interactive live cluster environment, we'll conduct hands-on labs to gain practical experience in deploying, managing, and understanding the concepts related to Kubernetes Pods.

Prerequisites:
- Basic understanding of containerization and Kubernetes concepts.
- Meshery Playground access. If you don't have an account, sign up at [Meshery Playground](https://meshery.layer5.io/play).

Lab Scenario: Deploying and Exploring Pods in a Microservices Architecture

Objective:
Learn how to create, manage, and explore Kubernetes Pods within the context of a microservices architecture.

### Steps:

#### 1. **Accessing Meshery Playground:**
   - Log in to the [Meshery Playground](https://meshery.layer5.io/play) using your credentials.
   - Navigate to the Meshery Playground dashboard.

#### 2. **Creating a Simple Pod:**
   - Deploy a basic Pod with a single container using a YAML configuration file.

\```yaml
# simple-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
  - name: mycontainer
    image: nginx:latest
\```

Apply the Pod configuration:

\```bash
kubectl apply -f simple-pod.yaml
\```

#### 3. **Exploring Pod States and Information:**
   - Use Meshery Playground to explore the different states a Pod can be in and gather information about the deployed Pod.

\```bash
kubectl get pod mypod
kubectl describe pod mypod
\```

#### 4. **Interacting with Pods:**
   - Gain hands-on experience by interacting with the Pod using commands and Meshery Playground's interactive features.

\```bash
kubectl exec -it mypod -- /bin/bash
\```

#### 5. **Pod Networking:**
   - Explore networking aspects related to Pods, including accessing a Pod from another Pod.

\```bash
kubectl expose pod mypod --type=NodePort --port=80
\```

#### 6. **Deleting and Recreating Pods:**
   - Understand the impact of deleting and recreating Pods and observe how Meshery Playground reflects these changes.

\```bash
kubectl delete pod mypod
\```

#### 7. **Pod Annotations and Labels:**
   - Learn how to add annotations and labels to Pods and leverage them for better organization and management.

\```yaml
# annotated-labeled-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
  annotations:
    environment: production
  labels:
    app: myapp
spec:
  containers:
  - name: mycontainer
    image: nginx:latest
\```

Apply the annotated and labeled Pod configuration:

\```bash
kubectl apply -f annotated-labeled-pod.yaml
\```

#### 8. **Clean-Up:**
   - Delete the Pods and associated resources after completing the lab.

\```bash
kubectl delete pod mypod
\```

#### 9. **Saving and Sharing:**
   - Save your scenario in Meshery Playground for future reference.
   - Share your Pod scenarios with the Meshery community for collaborative learning.

### Conclusion
Congratulations! You've successfully completed the lab on exploring Kubernetes Pods using Meshery Playground. This hands-on experience has provided valuable insights into the deployment, management, and interaction with Pods in a Kubernetes environment. Continue exploring more scenarios in the Meshery Playground to enhance your skills in container orchestration.



{% include suggested-reading.html language="en" %}