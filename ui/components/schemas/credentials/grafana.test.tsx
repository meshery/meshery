import { describe, expect, it } from 'vitest';
import { grafanaSchema } from './grafana';

describe('grafanaSchema', () => {
  it('declares an object type titled "Grafana"', () => {
    expect(grafanaSchema.type).toBe('object');
    expect(grafanaSchema.title).toBe('Grafana');
  });

  it('requires credentialName at the top level', () => {
    expect(grafanaSchema.required).toEqual(['credentialName']);
    expect(grafanaSchema.properties.credentialName.type).toBe('string');
  });

  it('describes secret with grafanaURL and grafanaAPIKey, both required', () => {
    const secret = grafanaSchema.properties.secret;
    expect(secret.type).toBe('object');
    expect(secret.properties.grafanaURL.type).toBe('string');
    expect(secret.properties.grafanaAPIKey.type).toBe('string');
    expect(secret.required).toEqual(['grafanaURL', 'grafanaAPIKey']);
  });
});
