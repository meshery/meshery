---
title: AWS EC2
---

<h1> How to Deploy AWS EC2 Instances with Meshery</h1>
<p>
    Learn how to configure, deploy, and manage AWS EC2 instances using Meshery's intuitive <a href="/extensions/kanvas">Kanvas</a> extension. In order to manage AWS resources, Meshery uses <a href="https://aws.amazon.com/blogs/containers/aws-controllers-for-kubernetes-ack/">AWS Controllers for Kubernetes (ACK)</a>. ACK facilitates the bridge between Kubernetes and AWS services, enabling Meshery to manage AWS resources and Meshery enabling you to benefit from the enhanced experience that Meshery and its extensions offer.
</p>

<p>
    To get started, you'll connect your Kubernetes cluster to Meshery. Using pre-configured designs available in the <a href="https://cloud.layer5.io/catalog">Meshery Catalog</a>, you'll then configure the EC2 controller with AWS credentials, deploy essential resources (like VPCs, subnets, and security groups) to create the network topology for your EC2 instance. Finally, you'll deploy the instance itself and verify its status through the AWS Management Console or Meshery's Kanvas operator mode.
</p>

<p>
    <a href="/guides/tutorials/deploy-aws-ec2-instances-with-meshery">Read the full guide on deploying AWS EC2 instances with Meshery ➔</a>
</p>