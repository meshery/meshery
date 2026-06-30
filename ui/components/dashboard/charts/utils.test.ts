import { describe, expect, it } from 'vitest';
import { getLegendTemplate } from './utils';

describe('getLegendTemplate', () => {
  it('embeds the matching node count from the data list', () => {
    const html = getLegendTemplate('Ready', '#0f0', [
      { id: 'Ready', value: 7 },
      { id: 'Not Ready', value: 1 },
    ]);
    expect(html).toContain('>7<');
    expect(html).toContain('Ready');
    expect(html).toContain('color:#0f0');
  });

  it('falls back to 0 when the matching id is missing', () => {
    const html = getLegendTemplate('Unknown', '#abc', [{ id: 'Ready', value: 4 }]);
    expect(html).toContain('>0<');
    expect(html).toContain('Unknown');
    expect(html).toContain('color:#abc');
  });

  it('falls back to 0 when the matching value is zero (|| short-circuit)', () => {
    const html = getLegendTemplate('Ready', '#fff', [{ id: 'Ready', value: 0 }]);
    expect(html).toContain('>0<');
  });

  it('returns a string snippet wrapped in a flex column container div', () => {
    const html = getLegendTemplate('Ready', '#000', []);
    expect(html).toMatch(/^<div[\s\S]+<\/div>$/);
    expect(html).toContain('flex-direction:column');
  });
});
