import { describe, expect, it } from 'vitest';
import { extractKubernetesCredentials, isKubernetesConnected } from '../kubernetesHelpers';

describe('isKubernetesConnected', () => {
  it('returns true only when both cluster is configured AND ping succeeded', () => {
    expect(isKubernetesConnected(true, true)).toBe(true);
  });

  it('returns false when cluster is configured but ping failed', () => {
    expect(isKubernetesConnected(true, false)).toBe(false);
  });

  it('returns false when cluster is not configured', () => {
    expect(isKubernetesConnected(false, true)).toBe(false);
    expect(isKubernetesConnected(false, false)).toBe(false);
  });

  it('treats nullish inputs as falsy', () => {
    expect(isKubernetesConnected(null, true)).toBe(false);
    expect(isKubernetesConnected(undefined, true)).toBe(false);
    expect(isKubernetesConnected(true, null)).toBe(false);
    expect(isKubernetesConnected(true, undefined)).toBe(false);
  });
});

describe('extractKubernetesCredentials', () => {
  it('maps a server kubeconfig payload into the UI credential shape', () => {
    const payload = {
      name: 'demo-cluster-cred',
      cluster: {
        name: 'demo-cluster',
        cluster: {
          server: 'https://demo.example.com',
          'certificate-authority-data': 'CA_DATA',
        },
      },
      auth: {
        name: 'demo-user',
        user: {
          token: 'tok-abc',
          'client-certificate-data': 'CERT_DATA',
          'client-key-data': 'KEY_DATA',
        },
      },
    };

    expect(extractKubernetesCredentials(payload)).toEqual({
      credentialName: 'demo-cluster-cred',
      secret: {
        clusterName: 'demo-cluster',
        clusterServerURL: 'https://demo.example.com',
        auth: {
          clusterUserName: 'demo-user',
          clusterToken: 'tok-abc',
          clusterClientCertificateData: 'CERT_DATA',
          clusterCertificateAuthorityData: 'CA_DATA',
          clusterClientKeyData: 'KEY_DATA',
        },
      },
    });
  });

  it('produces undefined fields when token / certs are missing', () => {
    const payload = {
      name: 'no-secret',
      cluster: { name: 'c1', cluster: { server: 'https://c1' } },
      auth: { name: 'u1', user: {} },
    };

    const creds = extractKubernetesCredentials(payload);
    expect(creds.credentialName).toBe('no-secret');
    expect(creds.secret.clusterName).toBe('c1');
    expect(creds.secret.clusterServerURL).toBe('https://c1');
    expect(creds.secret.auth.clusterUserName).toBe('u1');
    expect(creds.secret.auth.clusterToken).toBeUndefined();
    expect(creds.secret.auth.clusterClientCertificateData).toBeUndefined();
    expect(creds.secret.auth.clusterCertificateAuthorityData).toBeUndefined();
    expect(creds.secret.auth.clusterClientKeyData).toBeUndefined();
  });
});
