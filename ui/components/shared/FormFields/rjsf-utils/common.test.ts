import { describe, expect, it } from 'vitest';
import { selectCompSchema } from './common';

describe('selectCompSchema', () => {
  it('creates a single-select schema by default', () => {
    const schema = selectCompSchema(['A', 'B'], 'desc', 'Title', 'mySel');

    expect(schema).toEqual({
      properties: {
        mySel: {
          description: 'desc',
          enum: ['A', 'B'],
          title: 'Title',
          uniqueItems: true,
          'x-rjsf-grid-area': 12,
        },
      },
      required: ['mySel'],
      type: 'object',
    });
  });

  it('creates a multi-select schema when multiSelect is true', () => {
    const schema = selectCompSchema(['A', 'B', 'C'], 'desc', 'Title', 'multiSel', true);

    expect(schema.properties.multiSel).toEqual({
      description: 'desc',
      items: { enum: ['A', 'B', 'C'], type: 'string' },
      minItems: 1,
      type: 'array',
      title: 'Title',
      uniqueItems: true,
      'x-rjsf-grid-area': 12,
    });
    expect(schema.required).toEqual(['multiSel']);
    expect(schema.type).toBe('object');
  });

  it('handles an empty enum list', () => {
    const schema = selectCompSchema([], 'desc', 'Title', 'empty');
    expect(schema.properties.empty.enum).toEqual([]);
  });

  it('handles long enum lists for multi-select', () => {
    const enums = ['x', 'y', 'z', 'w'];
    const schema = selectCompSchema(enums, 'd', 't', 'list', true);
    expect(schema.properties.list.items?.enum).toEqual(enums);
    expect(schema.properties.list.minItems).toBe(1);
  });
});
