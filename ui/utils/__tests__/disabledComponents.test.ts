import { describe, expect, it } from 'vitest';
import { ProviderUiAccessControl } from '../disabledComponents';

describe('ProviderUiAccessControl', () => {
  it('exposes provider capabilities as an array for UI guards', () => {
    const registry = new ProviderUiAccessControl({
      capabilities: [{ feature: 'persist-meshery-patterns' }],
    });

    expect(registry.providerCapabilities).toEqual([{ feature: 'persist-meshery-patterns' }]);
  });

  it('returns an empty capability list when the provider payload is malformed', () => {
    const registry = new ProviderUiAccessControl({
      capabilities: { feature: 'persist-meshery-patterns' },
    });

    expect(registry.providerCapabilities).toEqual([]);
  });
});
