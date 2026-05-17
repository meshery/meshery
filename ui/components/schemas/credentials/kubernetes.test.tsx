import { describe, expect, it } from 'vitest';
import { kubernetesSchema } from './kubernetes';

describe('kubernetesSchema', () => {
  it('declares an object type titled "Kubernetes"', () => {
    expect(kubernetesSchema.type).toBe('object');
    expect(kubernetesSchema.title).toBe('Kubernetes');
  });

  it('requires credentialName at the top level', () => {
    expect(kubernetesSchema.required).toEqual(['credentialName']);
  });

  it('describes secret with clusterName and clusterServerURL', () => {
    const secret = kubernetesSchema.properties.secret;
    expect(secret.properties.clusterName.type).toBe('string');
    expect(secret.properties.clusterServerURL.type).toBe('string');
    expect(secret.required).toEqual(['clusterName', 'clusterServerURL']);
  });

  it('declares full set of cluster auth fields as required', () => {
    const auth = kubernetesSchema.properties.secret.properties.auth;
    expect(auth.required).toEqual([
      'clusterUserName',
      'clusterToken',
      'clusterClientCertificateData',
      'clusterClientKeyData',
      'clusterCertificateAuthorityData',
    ]);
    Object.keys(auth.properties).forEach((field) => {
      expect(auth.properties[field].type).toBe('string');
    });
  });
});
