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

  it('produces empty string fields when token / certs are missing', () => {
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
    expect(creds.secret.auth.clusterToken).toBe('');
    expect(creds.secret.auth.clusterClientCertificateData).toBe('');
    expect(creds.secret.auth.clusterCertificateAuthorityData).toBe('');
    expect(creds.secret.auth.clusterClientKeyData).toBe('');
  });

  it('safely handles missing auth and cluster sections without throwing', () => {
    const noAuth = {
      name: 'no-auth-cluster',
      cluster: { name: 'c2', cluster: { server: 'https://c2' } },
    };
    const credsNoAuth = extractKubernetesCredentials(noAuth);
    expect(credsNoAuth.credentialName).toBe('no-auth-cluster');
    expect(credsNoAuth.secret.clusterName).toBe('c2');
    expect(credsNoAuth.secret.auth.clusterUserName).toBe('');
    expect(credsNoAuth.secret.auth.clusterToken).toBe('');

    const noCluster = {
      name: 'no-cluster',
      auth: { name: 'u2' },
    };
    const credsNoCluster = extractKubernetesCredentials(noCluster);
    expect(credsNoCluster.credentialName).toBe('no-cluster');
    expect(credsNoCluster.secret.clusterName).toBe('');
    expect(credsNoCluster.secret.clusterServerURL).toBe('');
    expect(credsNoCluster.secret.auth.clusterUserName).toBe('u2');
  });

  it('safely handles nullish and non-object inputs without throwing', () => {
    expect(() => extractKubernetesCredentials(null)).not.toThrow();
    expect(extractKubernetesCredentials(null)).toEqual({
      credentialName: '',
      secret: {
        clusterName: '',
        clusterServerURL: '',
        auth: {
          clusterUserName: '',
          clusterToken: '',
          clusterClientCertificateData: '',
          clusterCertificateAuthorityData: '',
          clusterClientKeyData: '',
        },
      },
    });

    expect(() => extractKubernetesCredentials(undefined)).not.toThrow();
    expect(extractKubernetesCredentials(undefined)).toEqual({
      credentialName: '',
      secret: {
        clusterName: '',
        clusterServerURL: '',
        auth: {
          clusterUserName: '',
          clusterToken: '',
          clusterClientCertificateData: '',
          clusterCertificateAuthorityData: '',
          clusterClientKeyData: '',
        },
      },
    });
    expect(() => extractKubernetesCredentials('not an object')).not.toThrow();
  });
});
