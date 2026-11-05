import { describe, expect, it } from 'vitest';
import { buildVarValues, resolveDatasourceUid } from './resolve';
import type { Datasource, TemplateVar } from './types';

describe('buildVarValues', () => {
  it('resolves single, multi, all/unset, and skips datasource vars', () => {
    const vars: TemplateVar[] = [
      { name: 'job', type: 'query', current: ['api'] },
      { name: 'instance', type: 'query', current: ['a', 'b'], multi: true },
      { name: 'env', type: 'query', current: ['$__all'], includeAll: true },
      { name: 'empty', type: 'query', current: [] },
      { name: 'datasource', type: 'datasource', current: ['uid-1'] },
    ];
    expect(buildVarValues(vars)).toEqual({
      job: 'api',
      instance: 'a|b',
      env: '.*',
      empty: '.*',
    });
  });
});

describe('resolveDatasourceUid', () => {
  const datasources: Datasource[] = [
    { uid: 'prom-uid', name: 'Prometheus', type: 'prometheus', isDefault: true },
    { uid: 'loki-uid', name: 'Loki', type: 'loki' },
  ];

  it('resolves a $datasource variable to its current value (by name or uid)', () => {
    const vars: TemplateVar[] = [
      { name: 'datasource', type: 'datasource', current: ['Prometheus'] },
    ];
    expect(resolveDatasourceUid('$datasource', vars, datasources)).toBe('prom-uid');
    expect(resolveDatasourceUid('${datasource}', vars, datasources)).toBe('prom-uid');
  });

  it('falls back to the default datasource for unset var, empty, or "default"', () => {
    expect(resolveDatasourceUid('$datasource', [], datasources)).toBe('prom-uid');
    expect(resolveDatasourceUid('', [], datasources)).toBe('prom-uid');
    expect(resolveDatasourceUid('default', [], datasources)).toBe('prom-uid');
  });

  it('passes through a concrete uid', () => {
    expect(resolveDatasourceUid('loki-uid', [], datasources)).toBe('loki-uid');
  });
});
