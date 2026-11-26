import { describe, expect, it } from 'vitest';
import { resolveExpr } from './promql';
import type { TimeWindow } from './types';

const win: TimeWindow = { start: 1000, end: 4600, step: 15 }; // range 3600s

describe('resolveExpr', () => {
  it.each(['$pod', '${pod}', '${pod:regex}'])('upgrades an exact matcher for %s', (ref) => {
    expect(resolveExpr(`m{pod="${ref}"}`, {}, win)).toBe('m{pod=~".*"}');
  });

  it('keeps exact matching when an unresolved variable is only part of the value', () => {
    expect(resolveExpr('m{app="app.v1-$pod",node="${pod}[0]",zone="$pod$region"}', {}, win)).toBe(
      'm{app="app.v1-.*",node=".*[0]",zone=".*.*"}',
    );
  });

  it('preserves exact matching for a known value containing a variable reference', () => {
    expect(resolveExpr('m{app="$app",pod="$pod"}', { app: '$pod' }, win)).toBe(
      'm{app=".*",pod=~".*"}',
    );
  });

  it('preserves exact literals, negative matchers, and existing regex syntax', () => {
    expect(resolveExpr('m{a=".*",b!="$pod",c!~"v.+$pod",d=~"v.+$pod"}', {}, win)).toBe(
      'm{a=".*",b!=".*",c!~"v.+.*",d=~"v.+.*"}',
    );
  });

  it('only upgrades actual matchers in expressions containing matcher-shaped strings', () => {
    expect(resolveExpr(`label_join(m{pod = "$pod"}, "out", 'pod="$pod"', "src")`, {}, win)).toBe(
      `label_join(m{pod =~ ".*"}, "out", 'pod=".*"', "src")`,
    );
  });

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
    expect(resolveExpr('m{pod="$__custom"}', {}, win)).toBe('m{pod="$__custom"}');
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
