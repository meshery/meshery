import { describe, expect, it } from 'vitest';
import { CapabilitiesRegistry } from '../disabledComponents';

describe('CapabilitiesRegistry', () => {
  it('exposes provider capabilities as an array for UI guards', () => {
    const registry = new CapabilitiesRegistry({
      capabilities: [{ feature: 'persist-meshery-patterns' }],
    });

    expect(registry.capabilities).toEqual([{ feature: 'persist-meshery-patterns' }]);
  });

  it('returns an empty capability list when the provider payload is malformed', () => {
    const registry = new CapabilitiesRegistry({
      capabilities: { feature: 'persist-meshery-patterns' },
    });

    expect(registry.capabilities).toEqual([]);
  });
});
