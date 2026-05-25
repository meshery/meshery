import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: { primary: { main: '#0000ff' }, mode: 'light' },
  }),
}));

import {
  HyperLinkDiv,
  calculateGrid,
  getHyperLinkDiv,
  getRefinedJsonSchema,
  getSchema,
  safeDisplayValue,
  safeStringTitle,
  userPromptKeys,
} from './helper';

describe('safeDisplayValue', () => {
  it('returns the original primitive', () => {
    expect(safeDisplayValue('hello')).toBe('hello');
    expect(safeDisplayValue(42)).toBe(42);
    expect(safeDisplayValue(true)).toBe(true);
  });

  it('returns an empty string for null and undefined', () => {
    expect(safeDisplayValue(null)).toBe('');
    expect(safeDisplayValue(undefined)).toBe('');
  });

  it('coerces objects to a string', () => {
    expect(safeDisplayValue({ foo: 'bar' })).toBe('[object Object]');
  });

  it('passes through valid React elements', () => {
    const el = <span>Hi</span>;
    expect(safeDisplayValue(el)).toBe(el);
  });
});

describe('safeStringTitle', () => {
  it('handles primitives and nullish values', () => {
    expect(safeStringTitle(null)).toBe('');
    expect(safeStringTitle('title')).toBe('title');
    expect(safeStringTitle(7)).toBe('7');
    expect(safeStringTitle({})).toBe('[object Object]');
  });
});

describe('userPromptKeys', () => {
  it('lists the schema combinator keywords', () => {
    expect(userPromptKeys).toEqual(['allOf', 'anyOf', 'oneOf']);
  });
});

describe('getRefinedJsonSchema', () => {
  it('strips top-level title/description and sorts properties by type', () => {
    const schema = {
      title: 'My schema',
      description: 'desc',
      type: 'object',
      properties: {
        nested: { type: 'object', properties: {} },
        name: { type: 'string' },
        age: { type: 'integer' },
      },
    };
    const refined = getRefinedJsonSchema(schema, true, vi.fn());
    expect(refined.title).toBe('');
    expect(refined.description).toBe('');
    const keys = Object.keys(refined.properties);
    expect(keys.indexOf('name')).toBeLessThan(keys.indexOf('nested'));
  });

  it('coerces x-kubernetes fields to RJSF-compatible types', () => {
    const schema: any = {
      type: 'object',
      properties: {
        intOrString: { 'x-kubernetes-int-or-string': true },
        preserveUnknown: { 'x-kubernetes-preserve-unknown-fields': true },
      },
    };
    const refined: any = getRefinedJsonSchema(schema, false, vi.fn());
    expect(refined.properties.intOrString.type).toBe('string');
    expect(refined.properties.preserveUnknown.type).toBe('object');
  });

  it('passes schemas through unchanged when hideTitle is false', () => {
    const schema = {
      title: 'keep me',
      type: 'object',
      properties: { name: { type: 'string' } },
    };
    const refined = getRefinedJsonSchema(schema, false, vi.fn());
    expect(refined.title).toBe('keep me');
    expect(refined.description).toBe('');
  });
});

describe('calculateGrid', () => {
  it('returns a wider grid for object types', () => {
    const grid = calculateGrid({
      type: 'object',
      content: { props: { schema: {} } },
    });
    expect(grid.lg).toBe(12);
  });

  it('returns a normal grid for primitive types', () => {
    const grid = calculateGrid({
      type: 'string',
      content: { props: { schema: {} } },
    });
    expect(grid.lg).toBe(6);
  });

  it('honors x-rjsf-grid-area overrides', () => {
    const grid = calculateGrid({
      type: 'string',
      content: { props: { schema: { 'x-rjsf-grid-area': 4 } } },
    });
    expect(grid.lg).toBe(4);
  });
});

describe('HyperLinkDiv / getHyperLinkDiv', () => {
  it('renders markdown style links as anchors', () => {
    render(<HyperLinkDiv text="See [docs](https://example.com)" />);
    const a = document.querySelector('a');
    expect(a?.getAttribute('href')).toBe('https://example.com');
  });

  it('renders plain url text as anchors', () => {
    render(<HyperLinkDiv text="Visit https://meshery.io" />);
    const a = document.querySelector('a');
    expect(a?.getAttribute('href')).toContain('meshery.io');
  });

  it('exposes a legacy helper that returns the same element', () => {
    const node = getHyperLinkDiv('Simple text');
    expect(React.isValidElement(node)).toBe(true);
  });
});

describe('getSchema', () => {
  it('returns undefined for unknown credential types', () => {
    expect(getSchema('unknown')).toBeUndefined();
  });
});
