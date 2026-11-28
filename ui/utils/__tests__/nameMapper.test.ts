import { describe, expect, it } from 'vitest';
import { podNameMapper, versionMapper, TelemetryComps, isTelemetryComponent } from '../nameMapper';

describe('podNameMapper', () => {
  it('maps known istio pod prefixes to their human-readable names', () => {
    expect(podNameMapper('istio', 'istio-galley-xyz-1')).toBe('Istio Galley');
    expect(podNameMapper('istio', 'istio-ingressgateway-xyz-1')).toBe('Istio Ingress Gateway');
    expect(podNameMapper('istio', 'istio-egressgateway-xyz-1')).toBe('Istio Egress Gateway');
    expect(podNameMapper('istio', 'istio-policy-xyz-1')).toBe('Istio Policy');
    expect(podNameMapper('istio', 'istio-citadel-xyz-1')).toBe('Istio Citadel');
    expect(podNameMapper('istio', 'istio-telemetry-xyz-1')).toBe('Istio Telemetry');
    expect(podNameMapper('istio', 'istio-pilot-xyz-1')).toBe('Istio Pilot');
    expect(podNameMapper('istio', 'istio-tracing-xyz-1')).toBe('Istio Tracing');
    expect(podNameMapper('istio', 'istio-sidecar-injector-xyz-1')).toBe('Istio Sidecar Injector');
  });

  it('returns a prettified name when istio prefix is unknown', () => {
    // Falls through to the generic prettifier: drops the last two dash segments
    // and capitalises remaining tokens
    expect(podNameMapper('istio', 'my-custom-pod-abc-123')).toBe('My Custom Pod');
  });

  it('returns simple labels for kiali / grafana / prometheus', () => {
    expect(podNameMapper('kiali', 'irrelevant')).toBe('Kiali');
    expect(podNameMapper('grafana', 'irrelevant')).toBe('Grafana');
    expect(podNameMapper('prometheus', 'irrelevant')).toBe('Prometheus');
  });

  it('prettifies pod names for unknown service meshes', () => {
    expect(podNameMapper('linkerd', 'linkerd-controller-pod-abc-123')).toBe(
      'Linkerd Controller Pod',
    );
  });

  it('handles a pod name with too few dash segments by returning an empty string after slicing', () => {
    // `'pod-1'.split('-').slice(0,-2)` yields [], which joins to ''
    expect(podNameMapper('unknown', 'pod-1')).toBe('');
  });

  it('handles pod names with consecutive dashes without throwing', () => {
    expect(podNameMapper('linkerd', 'linkerd--controller--pod-abc-123')).toBe(
      'Linkerd Controller Pod',
    );
  });

  it('handles nullish or non-string pod names safely', () => {
    expect(podNameMapper('istio', undefined)).toBe('');
    expect(podNameMapper('istio', null)).toBe('');
    expect(podNameMapper('unknown', undefined)).toBe('');
    expect(podNameMapper('unknown', null)).toBe('');
    expect(podNameMapper('unknown', 12345)).toBe('');
    expect(podNameMapper('kiali', undefined)).toBe('Kiali');
    expect(podNameMapper('grafana', null)).toBe('Grafana');
    expect(podNameMapper('prometheus', undefined)).toBe('Prometheus');
  });
});

describe('versionMapper', () => {
  it('prefixes with v if missing', () => {
    expect(versionMapper('1.0.0')).toBe('v1.0.0');
  });

  it('leaves a v-prefixed version unchanged', () => {
    expect(versionMapper('v1.2.3')).toBe('v1.2.3');
  });

  it('treats a string that already starts with lowercase v as having a v prefix', () => {
    expect(versionMapper('valpha')).toBe('valpha');
  });

  it('prefixes uppercase V versions (since the check is for lowercase v only)', () => {
    expect(versionMapper('V1')).toBe('vV1');
  });

  it('returns empty string for nullish or empty versions', () => {
    expect(versionMapper('')).toBe('');
    expect(versionMapper(undefined)).toBe('');
    expect(versionMapper(null)).toBe('');
  });

  it('returns empty string for non-string versions', () => {
    expect(versionMapper(1.0)).toBe('');
    expect(versionMapper({})).toBe('');
  });
});

describe('TelemetryComps', () => {
  it('exposes grafana and prometheus identifiers', () => {
    expect(TelemetryComps).toEqual({ GRAFANA: 'grafana', PROMETHEUS: 'prometheus' });
  });
});

describe('isTelemetryComponent', () => {
  it('recognises Grafana variants case-insensitively', () => {
    expect(isTelemetryComponent('Grafana')).toBe('grafana');
    expect(isTelemetryComponent('grafana-instance-1')).toBe('grafana');
    expect(isTelemetryComponent('my-grafana')).toBe('grafana');
  });

  it('recognises Prometheus variants case-insensitively', () => {
    expect(isTelemetryComponent('Prometheus')).toBe('prometheus');
    expect(isTelemetryComponent('prom-prometheus-server')).toBe('prometheus');
  });

  it('returns empty string for non-telemetry components', () => {
    expect(isTelemetryComponent('istio')).toBe('');
    expect(isTelemetryComponent('kubernetes')).toBe('');
  });

  it('returns empty string for nullish names (optional chaining)', () => {
    expect(isTelemetryComponent(undefined)).toBe('');
    expect(isTelemetryComponent(null)).toBe('');
  });

  it('returns empty string for non-string names', () => {
    expect(isTelemetryComponent(123)).toBe('');
    expect(isTelemetryComponent({})).toBe('');
    expect(isTelemetryComponent([])).toBe('');
  });
});
