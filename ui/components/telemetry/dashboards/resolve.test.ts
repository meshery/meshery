import { describe, expect, it } from 'vitest';
import { buildVarValues, resolveDatasourceUid, resolveExpr } from './resolve';
import type { Datasource, TemplateVar, TimeWindow } from './types';

const win: TimeWindow = { start: 1000, end: 4600, step: 15 }; // range 3600s

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

describe('resolveExpr', () => {
  it('substitutes time macros', () => {
    expect(resolveExpr('rate(x[$__rate_interval])', {}, win)).toBe('rate(x[60s])');
    expect(resolveExpr('x[$__interval]', {}, win)).toBe('x[15s]');
    expect(resolveExpr('x[$__range]', {}, win)).toBe('x[3600s]');
    expect(resolveExpr('${__interval_ms}', {}, win)).toBe('15000');
  });

  it('substitutes known variables and falls back to .* for unknowns', () => {
    const vars = { job: '.*', instance: 'a|b' };
    expect(resolveExpr('m{job=~"$job", instance=~"$instance", state="x"}', vars, win)).toBe(
      'm{job=~".*", instance=~"a|b", state="x"}',
    );
    expect(resolveExpr('m{pod=~"${pod}"}', vars, win)).toBe('m{pod=~".*"}');
  });

  it('does not treat $__ macros as unknown variables', () => {
    expect(resolveExpr('rate(x[$__rate_interval_ms])', {}, win)).toBe('rate(x[60000])');
  });

  it('upgrades exact matchers to regex when a variable resolves to match-all or multi-value', () => {
    const vars = { cluster: '.*', namespace: '.*', pod: '.*' };
    // Real kube-prometheus panel: exact matchers must become =~ so .* matches.
    expect(
      resolveExpr(
        'container_memory_rss{namespace="$namespace", pod="$pod", cluster="$cluster", container != ""}',
        vars,
        win,
      ),
    ).toBe('container_memory_rss{namespace=~".*", pod=~".*", cluster=~".*", container != ""}');
    // != upgrades to !~
    expect(resolveExpr('m{ns!="$namespace"}', vars, win)).toBe('m{ns!~".*"}');
    // multi-value (a|b) also needs a regex matcher
    expect(resolveExpr('m{instance="$instance"}', { instance: 'a|b' }, win)).toBe(
      'm{instance=~"a|b"}',
    );
  });

  it('leaves matchers untouched for concrete single values and non-variable literals', () => {
    // concrete value -> keep exact `=`
    expect(resolveExpr('m{job="$job"}', { job: 'api' }, win)).toBe('m{job="api"}');
    // already a regex matcher stays regex
    expect(resolveExpr('m{job=~"$job"}', { job: '.*' }, win)).toBe('m{job=~".*"}');
    // literal matchers with no variable are never rewritten
    expect(resolveExpr('m{container != "", container != "POD"}', { pod: '.*' }, win)).toBe(
      'm{container != "", container != "POD"}',
    );
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
