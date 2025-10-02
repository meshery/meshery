# New Relic Architecture Design

This directory contains a comprehensive New Relic architecture design for Meshery Catalog.

## Design Overview

The **New Relic Architecture** design demonstrates a complete observability stack using New Relic monitoring tools integrated with a Kubernetes-based application infrastructure.

### Architecture Components

#### Core Application Components
- **Application Server**: Multi-replica deployment running the main web application
- **Database**: PostgreSQL StatefulSet for persistent data storage
- **Cache Service**: Redis for in-memory caching
- **Load Balancer**: Kubernetes LoadBalancer service for traffic distribution

#### New Relic Observability Components
1. **New Relic APM**: Application Performance Monitoring for tracking application metrics, distributed tracing, and error analytics
2. **New Relic Infrastructure Agent**: DaemonSet deployment for host-level metrics collection
3. **New Relic Kubernetes Integration**: Kubernetes cluster monitoring for pods, containers, and control plane
4. **New Relic Browser Agent**: Real User Monitoring (RUM) for frontend performance
5. **New Relic Logs Integration**: Fluentbit-based log aggregation and forwarding
6. **New Relic Synthetic Monitor**: Uptime monitoring and synthetic API/browser testing

#### Configuration & Security
- **New Relic Secret**: Kubernetes Secret for storing New Relic license keys
- **Synthetic Monitor ConfigMap**: Configuration for synthetic monitoring tests

## Key Features

- **Complete Observability Stack**: Full-stack monitoring from infrastructure to application to user experience
- **Kubernetes Native**: All components deployed as Kubernetes resources
- **Distributed Tracing**: End-to-end transaction tracking across services
- **Real-time Monitoring**: Live metrics, logs, and traces
- **Proactive Monitoring**: Synthetic tests to catch issues before users do

## Component Statistics

- **Total Components**: 11 services
- **Relationship Edges**: 15 connections
- **Observability Components**: 6 New Relic integrations
- **Namespaces**: 2 (default, observability)

## Component Relationships

The design includes proper relationships (edges) between components:
- Load Balancer → Application Server
- Application Server → Database
- Application Server → Cache Service
- Application Server → New Relic Browser Agent
- New Relic APM → All Application Components
- New Relic Infrastructure Agent → New Relic APM
- New Relic Kubernetes Integration → New Relic APM
- New Relic Logs Integration → New Relic APM
- New Relic Synthetic Monitor → New Relic APM

## Deployment

To deploy this architecture:

1. Ensure you have a New Relic account and license key
2. Update the `new-relic-secret` with your actual license key (base64 encoded)
3. Import the design into Meshery
4. Deploy using Meshery's configuration management

## Configuration Requirements

Before deploying, update the following:
- **New Relic License Key**: Replace `<base64-encoded-license-key>` in the Secret
- **Cluster Name**: Update the cluster name in Kubernetes Integration
- **Application Name**: Customize the app name in APM configuration
- **Synthetic Monitor URLs**: Update URLs in the synthetic monitor configuration

## Security Considerations

- License keys are stored in Kubernetes Secrets
- All New Relic agents use secure HTTPS endpoints
- Infrastructure agent runs with read-only host access
- RBAC should be configured for the Kubernetes Integration service account

## Additional Resources

- [New Relic Documentation](https://docs.newrelic.com/)
- [Meshery Documentation](https://docs.meshery.io/)
- [Kubernetes Monitoring Best Practices](https://docs.newrelic.com/docs/kubernetes-pixie/kubernetes-integration/)

## Validation

This design has been validated using the included `validate_design.sh` script. All tests pass:
- ✓ YAML syntax validation
- ✓ File size check (16KB)
- ✓ Component verification
- ✓ Relationship validation
- ✓ Observability stack completeness

## Design Visualization

The design can be visualized in Meshery Kanvas at [https://kanvas.new](https://kanvas.new) or [Meshery Playground](https://playground.meshery.io/).

## Files

- `new-relic-architecture.yaml`: The main design file
- `validate_design.sh`: Validation script to test the design
- `README.md`: This documentation file

## Author

**Aviral Garg** ([@aviralgarg05](https://github.com/aviralgarg05))

## Created

October 2, 2025 - Hacktoberfest 2025

## Tags

`new-relic`, `observability`, `monitoring`, `apm`, `kubernetes`, `infrastructure`, `logging`, `tracing`, `synthetic-monitoring`, `meshery`, `hacktoberfest`

---

For questions or issues, please refer to the [Meshery Community](https://meshery.io/community) or [Discussion Forum](https://meshery.io/community#community-forums).
