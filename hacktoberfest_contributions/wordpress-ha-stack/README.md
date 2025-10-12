# WordPress High-Availability Stack

## Overview

A production-ready WordPress deployment designed for high traffic using Meshery. This design implements load balancing, caching, and persistent storage to ensure reliability and performance.

## Architecture

This design includes:

- **NGINX Ingress Controller**: Load balancer for incoming HTTP traffic
- **WordPress (2 replicas)**: High-availability web application with horizontal scaling
- **MySQL Database**: Persistent data storage with volume backup
- **Redis Cache**: Performance optimization through caching layer
- **Persistent Volumes**: Separate storage for WordPress files and MySQL data

## Components

### Load Balancing & Traffic
- **NGINX Ingress**: Routes external traffic to WordPress pods with load balancing

### Application Layer
- **WordPress Deployment**: 
  - 2 replicas for high availability
  - Connected to MySQL database
  - Integrated with Redis for caching
  - Persistent volume for uploaded files and themes

### Data Layer
- **MySQL Database**:
  - Single replica with persistent storage
  - 5Gi volume for database files
  - Configured for WordPress

### Caching Layer
- **Redis**:
  - In-memory caching
  - Reduces database load
  - Improves response times

### Storage
- **Persistent Volumes**:
  - `wordpress-pv`: 5Gi for WordPress files
  - `mysql-pv`: 5Gi for MySQL database
- **Persistent Volume Claims**:
  - Automatic binding to PVs
  - Ensures data persistence across pod restarts

## Relationships/Connections

```
NGINX Ingress → WordPress Service (HTTP traffic routing)
WordPress Pods → MySQL Database (Database connection)
WordPress Pods → Redis Cache (Cache connection)
WordPress Pods → WordPress PVC (Volume mount)
MySQL Database → MySQL PVC (Volume mount)
WordPress PVC → WordPress PV (Storage binding)
MySQL PVC → MySQL PV (Storage binding)
```

## Prerequisites

- Kubernetes cluster with:
  - NGINX Ingress Controller installed
  - Storage provisioner (standard storage class)
  - At least 2 CPU cores and 4GB RAM available

## Deployment Instructions

### Using Meshery UI

1. Open Kanvas in Meshery
2. Import this design file
3. Review the components and connections
4. Deploy to your connected Kubernetes cluster

### Using kubectl

```bash
kubectl apply -f wordpress-ha-stack.yaml
```

### Verify Deployment

```bash
# Check all resources
kubectl get all,pvc,pv,ingress

# Wait for pods to be ready
kubectl wait --for=condition=ready pod -l app=wordpress --timeout=300s

# Get ingress URL
kubectl get ingress wordpress-ingress
```

## Accessing WordPress

1. Get the Ingress IP/hostname:
   ```bash
   kubectl get ingress wordpress-ingress
   ```

2. For kind cluster (local development):
   ```bash
   kubectl port-forward -n ingress-nginx service/ingress-nginx-controller 8080:80
   ```
   Then access: http://localhost:8080

3. For cloud deployments, use the EXTERNAL-IP from the ingress

## Features

### High Availability
- Multiple WordPress replicas handle failures
- Load balancing distributes traffic evenly
- Persistent storage protects against data loss

### Performance
- Redis caching reduces database queries
- Multiple WordPress pods handle concurrent requests
- NGINX Ingress provides efficient load balancing

### Reliability
- Persistent volumes ensure data survives pod restarts
- Separate storage for application and database
- StatefulSet-like behavior for database

## Scaling

### Horizontal Scaling (WordPress)
```bash
kubectl scale deployment wordpress --replicas=3
```

### Vertical Scaling (Resources)
Edit the deployment and adjust resource limits/requests

## Monitoring

```bash
# Check pod status
kubectl get pods

# View WordPress logs
kubectl logs -l app=wordpress

# View MySQL logs
kubectl logs -l app=mysql

# View Redis logs
kubectl logs -l app=redis
```

## Troubleshooting

### WordPress pods not starting
- Check PVC binding: `kubectl get pvc`
- Verify MySQL is ready: `kubectl get pods -l app=mysql`

### Database connection errors
- Verify MySQL service: `kubectl get svc mysql`
- Check environment variables in WordPress deployment

### Ingress not working
- Ensure NGINX Ingress Controller is running:
  ```bash
  kubectl get pods -n ingress-nginx
  ```

## Design Metadata

- **Name**: WordPress High-Availability Stack
- **Category**: Web Applications
- **Technology**: WordPress, MySQL, Redis, Kubernetes
- **Use Case**: High-traffic website hosting
- **Difficulty**: Intermediate

## Contributing

This design is part of Hacktoberfest 2024 contributions to Meshery.

## License

Apache 2.0

