# Kafka Streaming Platform - Meshery Design

## Overview

This Meshery design imports and configures Apache Kafka, a distributed streaming platform, based on the Bitnami Kafka Helm chart (version 32.4.3) with Kafka version 4.0.0.

## Design Description

Apache Kafka is a distributed streaming platform designed to build real-time pipelines and can be used as a message broker or as a replacement for a log aggregation solution for big data applications.

## Key Components

This design includes the following Kubernetes resources:

### StatefulSets
- **kafka-controller**: Controller-eligible Kafka nodes that handle metadata and partition management
- **kafka-broker**: Kafka broker nodes that handle message storage and processing

### Services  
- **kafka-controller-headless**: Headless service for controller communication
- **kafka**: Main Kafka service for client connections

### Additional Resources
- **NetworkPolicy**: Controls network traffic to/from Kafka pods
- **PodDisruptionBudget**: Ensures high availability during voluntary disruptions
- **ServiceAccount**: Kubernetes service account for Kafka pods
- **ConfigMaps**: Configuration for Kafka server and JVM settings

## Deployment Information

- **Helm Chart**: bitnami/kafka:32.4.3
- **Application Version**: 4.0.0
- **Namespace**: default
- **Architecture**: KRaft mode (no Zookeeper dependency)

## Key Features

1. **KRaft Mode**: Uses Kafka's native metadata handling without Zookeeper
2. **High Availability**: Configured with pod disruption budgets and anti-affinity rules
3. **Security**: Includes network policies and service accounts
4. **Monitoring**: JMX metrics exposure for monitoring integration
5. **Persistent Storage**: StatefulSets ensure data persistence across pod restarts

## Port Configuration

- **9092**: Client connections (PLAINTEXT)
- **9093**: Controller communication
- **9094**: Inter-broker communication

## Environment Variables

Key environment variables configured:
- `KAFKA_CFG_PROCESS_ROLES`: controller,broker
- `KAFKA_CFG_NODE_ID`: Unique node identifier
- `KAFKA_CFG_CONTROLLER_QUORUM_VOTERS`: Controller quorum configuration
- `KAFKA_CFG_LISTENERS`: Listener configuration
- `KAFKA_CFG_ADVERTISED_LISTENERS`: Advertised listener configuration

## Usage

1. Deploy this design through Meshery
2. Access Kafka at `kafka:9092` from within the cluster
3. Use Kafka clients to produce and consume messages
4. Monitor through JMX metrics on exposed ports

## Files Generated

- `kafka-32.4.3.tgz`: Original Bitnami Helm chart
- `kafka-manifests.yaml`: Rendered Kubernetes manifests
- `kafka-meshery-design.json`: Meshery design file
- `kafka-design.yaml`: Simplified design format

## Original Helm Chart Information

- **Repository**: https://charts.bitnami.com/bitnami
- **Chart**: kafka
- **Version**: 32.4.3
- **App Version**: 4.0.0
- **Source**: https://github.com/bitnami/charts/tree/main/bitnami/kafka

## Publishing to Meshery Catalog

This design can be published to the Meshery Catalog following these steps:

1. Access Meshery UI at http://localhost:9081
2. Navigate to Configuration > Designs
3. Import the design file
4. Fill in the design metadata:
   - **Name**: Kafka Streaming Platform
   - **Type**: Deployment
   - **Description**: Apache Kafka distributed streaming platform
   - **Tags**: kafka, streaming, messaging, bitnami
   - **Category**: Infrastructure
5. Click "Publish" to submit for review

## Contributing

This design follows the contribution guidelines from issue #15790 for importing Helm charts into the Meshery catalog.

## License

This design is based on the Apache Kafka project and Bitnami Helm charts, both under Apache 2.0 license.