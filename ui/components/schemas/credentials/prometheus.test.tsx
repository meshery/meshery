import { describe, expect, it } from 'vitest';
import { prometheusSchema } from './prometheus';

describe('prometheusSchema', () => {
  it('declares an object type titled "Prometheus"', () => {
    expect(prometheusSchema.type).toBe('object');
    expect(prometheusSchema.title).toBe('Prometheus');
  });

  it('requires credentialName at the top level', () => {
    expect(prometheusSchema.required).toEqual(['credentialName']);
    expect(prometheusSchema.properties.credentialName.type).toBe('string');
  });

  it('describes a secret block containing prometheusURL', () => {
    const secret = prometheusSchema.properties.secret;
    expect(secret.type).toBe('object');
    expect(secret.properties.prometheusURL.type).toBe('string');
    expect(secret.required).toContain('prometheusURL');
  });
});
