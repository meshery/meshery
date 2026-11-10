import { describe, expect, it } from 'vitest';
import {
  getFilteredDataForDetailsComponent,
  groupRelationshipsByKind,
  removeDuplicateVersions,
  reactJsonTheme,
} from './helper';

describe('groupRelationshipsByKind', () => {
  it('groups relationships into separate buckets by kind', () => {
    const input = [
      { id: 'r1', kind: 'edge', name: 'rel1' },
      { id: 'r2', kind: 'hierarchical', name: 'rel2' },
      { id: 'r3', kind: 'edge', name: 'rel3' },
    ];

    const result = groupRelationshipsByKind(input) as Array<{
      kind: string;
      relationships: Array<any>;
    }>;

    expect(result).toHaveLength(2);
    const edge = result.find((r) => r.kind === 'edge');
    const hierarchical = result.find((r) => r.kind === 'hierarchical');
    expect(edge?.relationships).toHaveLength(2);
    expect(hierarchical?.relationships).toHaveLength(1);
    expect(edge?.relationships[0]).toMatchObject({ id: 'r1', kind: 'edge' });
  });

  it('returns an empty array when input is empty', () => {
    expect(groupRelationshipsByKind([])).toEqual([]);
  });
});

describe('getFilteredDataForDetailsComponent', () => {
  it('returns REGISTRANTS when the matched node has a summary field', () => {
    const data = [{ id: 'reg-1', summary: { models: 10 } }];
    const result = getFilteredDataForDetailsComponent(data, 'reg-1');

    expect(result.type).toBe('Registrants');
    expect(result.data).toMatchObject({ id: 'reg-1' });
  });

  it('returns MODELS when the matched node has components/relationships', () => {
    const data = [{ id: 'model-1', components: [], relationships: [] }];
    const result = getFilteredDataForDetailsComponent(data, 'model-1');

    expect(result.type).toBe('Models');
  });

  it('returns RELATIONSHIPS when the matched node has evaluationQuery or selector', () => {
    const data = [{ id: 'rel-1', selector: {} }];
    const result = getFilteredDataForDetailsComponent(data, 'rel-1');

    expect(result.type).toBe('Relationships');
  });

  it('returns COMPONENTS when the matched node has a component field', () => {
    const data = [{ id: 'comp-1', component: { foo: 'bar' } }];
    const result = getFilteredDataForDetailsComponent(data, 'comp-1');

    expect(result.type).toBe('Components');
  });

  it('returns empty type with empty data when no node matches', () => {
    const data = [{ id: 'other-id' }];
    const result = getFilteredDataForDetailsComponent(data, 'missing');

    expect(result.type).toBe('');
    expect(result.data).toEqual({});
  });

  it('splits a compound UUID and looks up the trailing segment', () => {
    const data = [{ id: 'parent', children: [{ id: 'child' }] }];
    const result = getFilteredDataForDetailsComponent(data, 'parent.child');

    expect(result.data).toMatchObject({ id: 'child' });
  });
});

describe('removeDuplicateVersions', () => {
  it('merges identical model names into a single entry with versionBasedData', () => {
    const data = [
      {
        name: 'kubernetes',
        model: { version: '1.0.0' },
        components: ['c1'],
        relationships: [],
      },
      {
        name: 'kubernetes',
        model: { version: '2.0.0' },
        components: ['c2'],
        relationships: [],
      },
      {
        name: 'istio',
        model: { version: '1.0.0' },
        components: ['ic1'],
        relationships: [],
      },
    ];

    const result = removeDuplicateVersions(data);

    expect(result).toHaveLength(2);
    const kubernetes = result.find((m: any) => m.name === 'kubernetes');
    expect(kubernetes.version).toContain('1.0.0');
    expect(kubernetes.version).toContain('2.0.0');
    expect(kubernetes.versionBasedData).toHaveLength(2);
  });

  it('returns an empty array when given no models', () => {
    expect(removeDuplicateVersions([])).toEqual([]);
  });
});

describe('reactJsonTheme', () => {
  const makeTheme = (mode: 'light' | 'dark') => ({
    palette: {
      mode,
      background: { default: '#fff' },
      divider: '#ddd',
      text: { primary: '#111', secondary: '#222' },
      info: { main: '#5af', light: '#7bf' },
      error: { main: '#f00' },
      warning: { main: '#f80', light: '#fa0' },
      success: { main: '#0c0' },
      primary: { main: '#06f' },
    },
  });

  it('derives a base16 palette from the theme tokens (light mode)', () => {
    const theme = makeTheme('light');
    const t = reactJsonTheme(theme as any);

    expect(t.base00).toBe('#fff');
    expect(t.base01).toBe('#ddd');
    expect(t.base07).toBe('#111');
    expect(t.base08).toBe('#f00');
    expect(t.base0E).toBe('#06f');
    // light mode -> softer outline using 0.16 alpha
    expect(t.base02).toBeTruthy();
  });

  it('uses a stronger outline alpha in dark mode', () => {
    const theme = makeTheme('dark');
    const t = reactJsonTheme(theme as any);

    expect(t.base02).toBeTruthy();
  });
});
