import { describe, expect, it } from 'vitest';
import { getFilterConfig } from './filterConfig';

const withResource = (resource: unknown) => ({ filterResource: JSON.stringify(resource) });

describe('getFilterConfig', () => {
  it('reads the configuration a save stores under filterResource.configuration.config', () => {
    expect(
      getFilterConfig(
        withResource({ component: { kind: 'WASMFilter' }, configuration: { config: 'a: 1' } }),
      ),
    ).toBe('a: 1');
  });

  it('reads filters saved by older releases under settings.config', () => {
    expect(getFilterConfig(withResource({ settings: { config: 'legacy: true' } }))).toBe(
      'legacy: true',
    );
  });

  it('prefers configuration.config when both shapes are present', () => {
    expect(
      getFilterConfig(
        withResource({ configuration: { config: 'new' }, settings: { config: 'old' } }),
      ),
    ).toBe('new');
  });

  it('ignores the snake_case filter_resource key, which the API does not send', () => {
    expect(
      getFilterConfig({
        filter_resource: JSON.stringify({ configuration: { config: 'x' } }),
      } as any),
    ).toBe('');
  });

  it.each([
    ['a missing filter', undefined],
    ['an empty filterResource', { filterResource: '' }],
    ['invalid JSON', { filterResource: '{invalid json' }],
    ['a resource without a config', withResource({ component: { kind: 'WASMFilter' } })],
    ['a non-string config', withResource({ configuration: { config: 42 } })],
  ])('returns an empty string for %s', (_label, filter) => {
    expect(getFilterConfig(filter as any)).toBe('');
  });
});
