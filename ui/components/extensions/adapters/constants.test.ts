import { describe, expect, it } from 'vitest';
import { ADAPTER_STATUS, adaptersList } from './constants';

describe('adapters constants', () => {
  it('exposes ENABLED and DISABLED statuses', () => {
    expect(ADAPTER_STATUS.ENABLED).toBe('ENABLED');
    expect(ADAPTER_STATUS.DISABLED).toBe('DISABLED');
  });

  it('includes the canonical adapter list with required fields', () => {
    expect(Object.keys(adaptersList).length).toBeGreaterThan(0);

    Object.entries(adaptersList).forEach(([key, adapter]) => {
      expect(adapter.name).toBeTruthy();
      expect(adapter.label).toMatch(/^meshery-/);
      expect(adapter.imageSrc).toMatch(/^\/static\/img\//);
      expect(typeof adapter.defaultPort).toBe('number');
      expect(adapter.enabled).toBe(false);
      expect(adapter.url).toBe('');
      expect(typeof key).toBe('string');
    });
  });

  it('describes each adapter with the standard meshery wording', () => {
    expect(adaptersList.ISTIO.description).toContain('Istio');
    expect(adaptersList.LINKERD.description).toContain('Linkerd');
  });

  it('defines unique default ports per adapter', () => {
    const ports = Object.values(adaptersList).map((a) => a.defaultPort);
    const uniquePorts = new Set(ports);
    expect(uniquePorts.size).toBe(ports.length);
  });
});
