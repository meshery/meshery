# Meshery Health Check Configuration

This document provides guidance on configuring Kubernetes health checks for Meshery deployments using the enhanced `/healthz/live` and `/healthz/ready` endpoints.

## Overview

Meshery implements Kubernetes-compliant health check endpoints that follow best practices from the Kubernetes API server:

- **`/healthz/live`** - Liveness probe endpoint
- **`/healthz/ready`** - Readiness probe endpoint

Both endpoints support a `?verbose=1` query parameter for detailed health check information.

## Health Check Behavior

### Liveness Probe (`/healthz/live`)
The liveness probe checks if Meshery server is running and responsive. Returns:
- **200 OK** with `"ok"` when the server is alive
- **503 Service Unavailable** when provider capabilities are not loaded

### Readiness Probe (`/healthz/ready`)
The readiness probe checks if Meshery is ready to accept traffic. It performs:
1. **Capability Check** (affects readiness): Verifies provider capabilities are loaded
2. **Extension Check** (informational only): Reports extension package status

Returns:
- **200 OK** when capabilities are loaded (ready to serve traffic)
- **503 Service Unavailable** when capabilities are not loaded

**Note**: Extension status is informational and does not affect the readiness state.

## Default Configuration

The Meshery Helm chart includes pre-configured health checks in `values.yaml`:

```yaml
probe:
  livenessProbe:
    enabled: true
    initialDelaySeconds: 80
    periodSeconds: 12
    failureThreshold: 4

  readinessProbe:
    enabled: true
    initialDelaySeconds: 10
    periodSeconds: 4
    failureThreshold: 4
```

## Installation Considerations

### Initial Installation

When installing Meshery for the first time:

1. **Liveness Probe**: Set `initialDelaySeconds` to allow time for the server to start and load configurations
   - Recommended: `80-120` seconds for first startup
   - This accounts for image pull, provider initialization, and capability loading

2. **Readiness Probe**: Can have a shorter delay as it will automatically retry
   - Recommended: `10-30` seconds
   - Lower `periodSeconds` (e.g., `4-5`) for faster readiness detection

Example for new installations:
```yaml
probe:
  livenessProbe:
    enabled: true
    initialDelaySeconds: 120  # Allow more time for initial setup
    periodSeconds: 15
    failureThreshold: 5

  readinessProbe:
    enabled: true
    initialDelaySeconds: 20
    periodSeconds: 5
    failureThreshold: 4
```

### Upgrade Considerations

When upgrading an existing Meshery deployment:

1. **Rolling Update Strategy**: The default deployment uses rolling updates which rely on health checks
2. **Capability Loading**: Provider capabilities must reload during upgrade
3. **Extension Availability**: Extension packages may need to be re-downloaded

Recommended settings for upgrades:
```yaml
probe:
  livenessProbe:
    enabled: true
    initialDelaySeconds: 60   # Shorter than initial install
    periodSeconds: 12
    failureThreshold: 4

  readinessProbe:
    enabled: true
    initialDelaySeconds: 15   # Allow time for capability reload
    periodSeconds: 5
    failureThreshold: 6       # Higher threshold during upgrades
```

## Advanced Configuration

### Verbose Mode for Debugging

During troubleshooting, you can manually check health status with verbose output:

```bash
# Check liveness with details
kubectl exec -n meshery deployment/meshery -- curl -s "http://localhost:8080/healthz/live?verbose=1"

# Check readiness with details
kubectl exec -n meshery deployment/meshery -- curl -s "http://localhost:8080/healthz/ready?verbose=1"
```

Example verbose output:
```
[+]capabilities ok
[i]extension extension package found
healthz check passed
```

Where:
- `[+]` indicates a passing health check
- `[-]` indicates a failing health check
- `[i]` indicates informational status (doesn't affect health)

### Custom Probe Configuration

For specific deployment scenarios, you can customize probe settings:

#### High Availability Setup
```yaml
probe:
  livenessProbe:
    enabled: true
    initialDelaySeconds: 90
    periodSeconds: 10
    failureThreshold: 3
    timeoutSeconds: 5

  readinessProbe:
    enabled: true
    initialDelaySeconds: 15
    periodSeconds: 3
    failureThreshold: 3
    successThreshold: 1
    timeoutSeconds: 3
```

#### Resource-Constrained Environments
```yaml
probe:
  livenessProbe:
    enabled: true
    initialDelaySeconds: 150  # More time for slower startup
    periodSeconds: 20
    failureThreshold: 5

  readinessProbe:
    enabled: true
    initialDelaySeconds: 30
    periodSeconds: 10
    failureThreshold: 6
```

### Startup Probes (Kubernetes 1.18+)

For better handling of slow-starting containers, consider adding a startup probe:

```yaml
# Add to deployment.yaml template
{{- if .Values.probe.startupProbe.enabled }}
startupProbe:
  httpGet:
    path: /healthz/live
    port: http
  initialDelaySeconds: {{ .Values.probe.startupProbe.initialDelaySeconds }}
  periodSeconds: {{ .Values.probe.startupProbe.periodSeconds }}
  failureThreshold: {{ .Values.probe.startupProbe.failureThreshold }}
{{- end }}
```

And in `values.yaml`:
```yaml
probe:
  startupProbe:
    enabled: true
    initialDelaySeconds: 0
    periodSeconds: 10
    failureThreshold: 30  # 30 * 10s = 5 minutes max startup time
```

## Monitoring and Troubleshooting

### Common Issues

1. **Pods stuck in `CrashLoopBackOff`**
   - Check if `initialDelaySeconds` is too short
   - Verify provider capabilities are being loaded
   - Use verbose mode to identify which check is failing

2. **Pods not becoming ready**
   - Check readiness probe logs with verbose mode
   - Ensure provider configuration is correct
   - Verify network connectivity to provider endpoints

3. **Frequent restarts**
   - Increase `failureThreshold` to tolerate temporary issues
   - Increase `periodSeconds` to reduce probe frequency
   - Check for resource constraints (CPU/memory)

### Debug Commands

```bash
# Check pod status
kubectl get pods -n meshery

# View pod events
kubectl describe pod -n meshery <pod-name>

# Check health endpoint directly
kubectl port-forward -n meshery deployment/meshery 8080:8080
curl http://localhost:8080/healthz/ready?verbose=1

# View container logs
kubectl logs -n meshery deployment/meshery -f
```

## Migration from Previous Versions

If upgrading from a version without enhanced health checks:

1. The endpoints maintain backward compatibility
2. Default probe configurations will work without changes
3. Consider adjusting timing parameters based on your environment
4. Test the upgrade in a staging environment first

## Best Practices

1. **Always enable both probes** - Liveness and readiness serve different purposes
2. **Set appropriate delays** - Account for provider initialization time
3. **Use startup probes** - For Kubernetes 1.18+ to handle initial startup separately
4. **Monitor probe metrics** - Use Kubernetes events and metrics to tune settings
5. **Test in staging** - Validate probe settings before production deployment
6. **Document custom settings** - If you modify defaults, document why

## References

- [Kubernetes Liveness, Readiness and Startup Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [Kubernetes API Health Endpoints](https://kubernetes.io/docs/reference/using-api/health-checks/)
- [Meshery Documentation](https://docs.meshery.io/)
