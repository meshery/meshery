import { describe, expect, it } from 'vitest';
import { resolveCredentialAuthSecret, resolveCredentialPayload } from '../credentialSecret';

// The shapes below are the ones read-only production inspection found on
// layer5io/meshery-cloud issue #5918: ~36.9k Kubernetes credentials shaped
// {auth, cluster}, ~154 shaped {secret: '<string>'}, and 8 double-nested
// {credentialName, secret: {...}} rows written by Meshery UI's credential form.
// The canonical shape is what meshery/schemas declares and what Cloud is moving
// to. All four have to keep resolving.

describe('resolveCredentialPayload', () => {
  it('returns the canonical secret object as the payload', () => {
    const secret = { prometheusURL: 'https://prom.example' };
    expect(resolveCredentialPayload(secret)).toEqual(secret);
  });

  it('returns the kubernetes secret object as the payload', () => {
    const secret = {
      auth: { clusterToken: 'tok' },
      cluster: { server: 'https://k8s.example' },
    };
    expect(resolveCredentialPayload(secret)).toEqual(secret);
  });

  it('unwraps the legacy double-nested payload', () => {
    expect(
      resolveCredentialPayload({
        credentialName: 'kube-cred',
        secret: { auth: { clusterToken: 'tok' } },
      }),
    ).toEqual({ auth: { clusterToken: 'tok' } });
  });

  it('unwraps a legacy double-nested payload keyed by name', () => {
    expect(
      resolveCredentialPayload({ name: 'kube-cred', secret: { clusterName: 'cluster-a' } }),
    ).toEqual({ clusterName: 'cluster-a' });
  });

  it('unwraps the legacy string-valued secret', () => {
    expect(resolveCredentialPayload({ secret: 'super-secret' })).toBe('super-secret');
  });

  it('leaves a canonical payload that carries its own secret field alone', () => {
    const secret = { grafanaURL: 'https://grafana.example', secret: { nested: 'value' } };
    expect(resolveCredentialPayload(secret)).toEqual(secret);
  });

  it('ignores an inherited secret property when classifying the wrapper', () => {
    // Every own key is a wrapper key, so the own-key check alone cannot decide
    // this: only the `secret` probe can. With `'secret' in obj` the object would
    // be treated as a wrapper and unwrapped to the *prototype's* value, even
    // though it owns no `secret` at all.
    const payload = Object.create({ secret: { grafanaAPIKey: 'inherited' } });
    payload.name = 'grafana-cred';

    expect(resolveCredentialPayload(payload)).toBe(payload);
    expect(resolveCredentialAuthSecret(payload)).toBeUndefined();
  });

  it('passes through values that are not objects', () => {
    expect(resolveCredentialPayload(undefined)).toBeUndefined();
    expect(resolveCredentialPayload(null)).toBeNull();
    expect(resolveCredentialPayload('raw')).toBe('raw');
  });
});

describe('resolveCredentialAuthSecret', () => {
  it('reads the legacy string-valued secret', () => {
    expect(resolveCredentialAuthSecret({ secret: 'super-secret' })).toBe('super-secret');
    expect(
      resolveCredentialAuthSecret({ credentialName: 'prom-token', secret: 'super-secret' }),
    ).toBe('super-secret');
  });

  it('reads the canonical grafana API key', () => {
    expect(
      resolveCredentialAuthSecret({
        grafanaURL: 'https://grafana.example',
        grafanaAPIKey: 'canonical-key',
      }),
    ).toBe('canonical-key');
  });

  it('reads the grafana API key out of the legacy double-nested shape', () => {
    expect(
      resolveCredentialAuthSecret({
        credentialName: 'grafana-cred',
        secret: { grafanaURL: 'https://grafana.example', grafanaAPIKey: 'nested-key' },
      }),
    ).toBe('nested-key');
  });

  it('prefers the canonical field over a sibling legacy secret string', () => {
    expect(
      resolveCredentialAuthSecret({ grafanaAPIKey: 'canonical-key', secret: 'legacy-key' }),
    ).toBe('canonical-key');
  });

  it('reports no auth material for an anonymous canonical credential', () => {
    expect(resolveCredentialAuthSecret({ prometheusURL: 'https://prom.example' })).toBeUndefined();
    expect(
      resolveCredentialAuthSecret({
        credentialName: 'prom-cred',
        secret: { prometheusURL: 'https://prom.example' },
      }),
    ).toBeUndefined();
  });

  it('reports no auth material for a kubernetes credential', () => {
    expect(
      resolveCredentialAuthSecret({
        auth: { clusterToken: 'tok' },
        cluster: { server: 'https://k8s.example' },
      }),
    ).toBeUndefined();
  });

  it('reports no auth material for absent or non-object secrets', () => {
    expect(resolveCredentialAuthSecret(undefined)).toBeUndefined();
    expect(resolveCredentialAuthSecret(null)).toBeUndefined();
    expect(resolveCredentialAuthSecret({})).toBeUndefined();
  });
});
