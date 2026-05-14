import { describe, expect, it, vi } from 'vitest';

vi.mock('../../lib/trueRandom', () => ({
  trueRandom: () => 0.1234567,
}));

vi.mock('./PatternService/helper', () => ({
  userPromptKeys: ['allOf', 'anyOf', 'oneOf'],
}));

import {
  buildUiSchema,
  createPatternFromConfig,
  getPatternAttributeName,
  recursiveCleanObject,
  recursiveCleanObjectExceptEmptyArray,
} from './helpers';

describe('getPatternAttributeName', () => {
  it('returns the internal pattern attribute name if present', () => {
    expect(getPatternAttributeName({ _internal: { patternAttributeName: 'foo' } })).toBe('foo');
  });

  it('returns NA when no internal data is provided', () => {
    expect(getPatternAttributeName({})).toBe('NA');
    expect(getPatternAttributeName(undefined as any)).toBe('NA');
  });
});

describe('recursiveCleanObject', () => {
  it('removes empty nested objects', () => {
    const obj: any = { a: { b: { c: {} } }, d: 'hello' };
    recursiveCleanObject(obj);
    expect(obj).toEqual({ d: 'hello' });
  });

  it('preserves non-empty descendants', () => {
    const obj: any = { a: { b: 'keep' }, c: {} };
    recursiveCleanObject(obj);
    expect(obj).toEqual({ a: { b: 'keep' } });
  });

  it('ignores primitive falsy values without recursing', () => {
    const obj: any = { a: 0, b: '', c: null, d: 'kept' };
    recursiveCleanObject(obj);
    expect(obj.d).toBe('kept');
  });
});

describe('recursiveCleanObjectExceptEmptyArray', () => {
  it('keeps empty arrays intact', () => {
    const obj: any = { tags: [], nested: { keep: 'me', empty: {} } };
    recursiveCleanObjectExceptEmptyArray(obj);
    expect(obj.tags).toEqual([]);
    expect(obj.nested.empty).toBeUndefined();
  });
});

describe('createPatternFromConfig', () => {
  it('builds a pattern with services from non-empty configs', () => {
    const config = {
      Service1: { settings: { foo: 'bar' } },
      Service2: { settings: {} }, // empty settings will be cleaned
    };

    const pattern = createPatternFromConfig(config, 'default');

    expect(pattern.name).toMatch(/^pattern-/);
    expect(pattern.services.Service1).toEqual({ settings: { foo: 'bar' }, namespace: 'default' });
    // Service2 should be removed since settings was empty and got cleaned
    expect(pattern.services.Service2).toBeUndefined();
  });

  it('removes settings attribute when it is set to literal true', () => {
    const config = { ServiceA: { settings: true } };
    const pattern = createPatternFromConfig(config, 'ns');
    expect(pattern.services.ServiceA).toBeDefined();
    expect(pattern.services.ServiceA.settings).toBeUndefined();
  });

  it('supports partial cleaning that preserves empty arrays', () => {
    const config = { ServiceB: { settings: { tags: [] } } };
    const pattern = createPatternFromConfig(config, 'ns', true);
    expect(pattern.services.ServiceB.settings.tags).toEqual([]);
  });
});

describe('buildUiSchema', () => {
  it('returns ui:order with expected fields', () => {
    const schema = { type: 'object', properties: { name: { type: 'string' } } };
    const uiSchema = buildUiSchema(schema);
    expect(uiSchema['ui:order']).toEqual([
      'metadata',
      'name',
      'namespace',
      'label',
      'annotation',
      '*',
    ]);
  });

  it('marks number percent fields as range widgets', () => {
    const schema = {
      type: 'object',
      properties: { someValuePercent: { type: 'number' } },
    };
    const uiSchema = buildUiSchema(schema);
    expect(uiSchema.someValuePercent['ui:widget']).toBe('range');
  });

  it('marks color fields as color widgets', () => {
    const schema = {
      type: 'object',
      properties: { primaryColor: { type: 'string' } },
    };
    const uiSchema = buildUiSchema(schema);
    expect(uiSchema.primaryColor['ui:widget']).toBe('color');
  });

  it('handles array schemas by descending into items', () => {
    const schema = {
      type: 'array',
      items: {
        type: 'object',
        properties: { secondaryColor: { type: 'string' } },
      },
    };
    const uiSchema = buildUiSchema(schema);
    expect(uiSchema.items['ui:label']).toBe(false);
  });

  it('returns sane uiSchema when given an empty schema', () => {
    const uiSchema = buildUiSchema({});
    expect(uiSchema['ui:order']).toBeDefined();
  });
});
