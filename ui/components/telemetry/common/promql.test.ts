// @vitest-environment node

import { spawnSync } from 'node:child_process';
import { describe, expect, it } from 'vitest';
import { resolveExpr } from './promql';
import type { TimeWindow } from './types';

const win: TimeWindow = { start: 1000, end: 4600, step: 15 }; // range 3600s
const nameConstraint = '__name__!="",';

describe('resolveExpr', () => {
  it.each(['"', "'", '`'])(
    'constrains unnamed selectors after known-value upgrades using %s quotes',
    (quote) => {
      expect(resolveExpr(`{pod=${quote}$pod${quote}}`, { pod: '.*' }, win)).toBe(
        `{${nameConstraint}pod=~${quote}.*${quote}}`,
      );
    },
  );

  it('constrains unnamed selectors after known multi-value upgrades', () => {
    expect(resolveExpr('{pod="$pod"}', { pod: 'api|' }, win)).toBe(
      `{${nameConstraint}pod=~"api|"}`,
    );

    expect(resolveExpr('{pod!="$pod"}', { pod: 'api|worker' }, win)).toBe(
      `{${nameConstraint}pod!~"api|worker"}`,
    );
  });

  it('does not add unnecessary constraints after known-value upgrades', () => {
    expect(resolveExpr('m{pod="$pod"}', { pod: '.*' }, win)).toBe('m{pod=~".*"}');

    expect(resolveExpr('{pod="$pod",job="api"}', { pod: '.*' }, win)).toBe('{pod=~".*",job="api"}');

    expect(resolveExpr('{pod="$pod"}', { pod: 'api|worker' }, win)).toBe('{pod=~"api|worker"}');

    expect(resolveExpr('{pod!="$pod"}', { pod: '.*' }, win)).toBe('{pod!~".*"}');
  });
  it.each(['$label', '${label}', '${label:raw}', '[[label]]', '[[label:raw]]'])(
    'substitutes template references in label names: %s',
    (ref) => {
      expect(resolveExpr(`m{${ref}="api"}`, { label: 'job' }, win)).toBe('m{job="api"}');

      expect(resolveExpr(`{${ref}="api",pod="$pod"}`, { label: 'job' }, win)).toBe(
        '{job="api",pod=~".*"}',
      );
    },
  );

  it('substitutes embedded template references in label names', () => {
    expect(resolveExpr('m{prefix_${label}="api"}', { label: 'job' }, win)).toBe(
      'm{prefix_job="api"}',
    );
  });
  it.each(['group_left', 'group_right'])(
    'recognizes %s as a vector matching modifier',
    (modifier) => {
      const prefix = `m + on(job) ${modifier}`;
      expect(resolveExpr(`${prefix} {pod="$pod"}`, {}, win)).toBe(
        `${prefix} {${nameConstraint}pod=~".*"}`,
      );
      expect(resolveExpr(`${prefix}(zone) or{pod="$pod"}`, {}, win)).toBe(
        `${prefix}(zone) or{pod=~".*"}`,
      );
    },
  );

  it.each(['and', 'or', 'unless'])(
    'distinguishes %s as a metric name and as an operator',
    (word) => {
      expect(resolveExpr(`${word}{pod="$pod"}`, {}, win)).toBe(`${word}{pod=~".*"}`);
      expect(resolveExpr(`m ${word}{pod="$pod"}`, {}, win)).toBe(
        `m ${word}{${nameConstraint}pod=~".*"}`,
      );
      expect(resolveExpr(`m ${word} ${word}{pod="$pod"}`, {}, win)).toBe(
        `m ${word} ${word}{pod=~".*"}`,
      );
      expect(resolveExpr(`sum(${word} # comment\n{pod="$pod"})`, {}, win)).toBe(
        `sum(${word} # comment\n{pod=~".*"})`,
      );
    },
  );

  it('keeps operand context across vector matching argument lists', () => {
    expect(resolveExpr('m or on(job) and{pod="$pod"}', {}, win)).toBe(
      'm or on(job) and{pod=~".*"}',
    );
    expect(resolveExpr('m >= bool {pod="$pod"}', {}, win)).toBe(
      `m >= bool {${nameConstraint}pod=~".*"}`,
    );
  });

  it.each([
    ['job=""', {}, false],
    ['job="$job"', { job: '' }, false],
    ['job="$job"', { job: 'api' }, true],
    ['job=~"$job"', { job: '.*' }, false],
    ['job=~"api|"', {}, false],
    ['job=~"api|worker"', {}, true],
    ['job!="api"', {}, false],
    ['job!=""', {}, true],
    ['job!~"api"', {}, false],
    ['job!~"api|"', {}, true],
    [String.raw`job=~"\x61*"`, {}, false],
    ['job=~`(?i)api`', {}, false],
    ['"http.requests"', {}, true],
  ] as const)(
    'checks whether %s excludes empty labels after substitution',
    (other, vars, excludesEmpty) => {
      const input = `{pod="$pod",${other}}`;
      const resolved = other.replace('$job', 'job' in vars ? vars.job : '');
      expect(resolveExpr(input, vars, win)).toBe(
        `{${excludesEmpty ? '' : nameConstraint}pod=~".*",${resolved}}`,
      );
    },
  );

  it('adds one metric-name constraint per unnamed selector and preserves every match-all', () => {
    expect(resolveExpr("{pod='$pod',zone=`[[zone]]`} or {pod=`\${pod}`}", {}, win)).toBe(
      `{${nameConstraint}pod=~'.*',zone=~\`.*\`} or {${nameConstraint}pod=~\`.*\`}`,
    );
  });

  it('counts known regex upgrades using their resolved operator and value', () => {
    expect(resolveExpr('{pod="$pod",job="$job"}', { job: 'api|' }, win)).toBe(
      `{${nameConstraint}pod=~".*",job=~"api|"}`,
    );
    expect(resolveExpr('{pod="$pod",job!="$job"}', { job: 'api|' }, win)).toBe(
      '{pod=~".*",job!~"api|"}',
    );
  });

  it('ignores selectors and constraints written inside strings and comments', () => {
    const expr = `label_join(m, "out", '{pod="$pod"}', "src") # {pod="$pod"}
or {pod="$pod", # job="api" }
job=""}`;
    expect(resolveExpr(expr, {}, win)).toBe(`label_join(m, "out", '{pod=".*"}', "src") # {pod=".*"}
or {${nameConstraint}pod=~".*", # job="api" }
job=""}`);
    expect(resolveExpr('m # {pod="$pod"}', {}, win)).toBe('m # {pod=".*"}');
  });

  it('preserves a named selector with a comment before its opening brace', () => {
    expect(resolveExpr('m # metric\n{pod="$pod"}', {}, win)).toBe('m # metric\n{pod=~".*"}');
  });

  it('recognizes a substituted metric name', () => {
    expect(resolveExpr('${metric}{pod="$pod"}', { metric: 'm' }, win)).toBe('m{pod=~".*"}');
  });

  it('preserves validity and includes series with empty or missing labels in unnamed selectors', () => {
    expect(resolveExpr('{pod="$pod"}', {}, win)).toBe(`{${nameConstraint}pod=~".*"}`);
  });

  it('preserves match-all when another matcher constrains the selector', () => {
    expect(resolveExpr('{pod="$pod",job="api"}', {}, win)).toBe('{pod=~".*",job="api"}');
    expect(resolveExpr('{__name__!="",pod="$pod"}', {}, win)).toBe('{__name__!="",pod=~".*"}');
  });

  it('preserves match-all for named metric selectors', () => {
    expect(resolveExpr('m{pod="$pod"}', {}, win)).toBe('m{pod=~".*"}');
  });

  it.each(["'", '`'])('upgrades unresolved exact matchers using %s quotes', (quote) => {
    expect(resolveExpr(`m{pod=${quote}$pod${quote}}`, {}, win)).toBe(`m{pod=~${quote}.*${quote}}`);

    expect(resolveExpr(`m{pod=${quote}[[pod]]${quote}}`, {}, win)).toBe(
      `m{pod=~${quote}.*${quote}}`,
    );

    expect(resolveExpr(`m{pod!=${quote}$pod${quote}}`, {}, win)).toBe(`m{pod!=${quote}.*${quote}}`);

    expect(resolveExpr(`label_join(m, "out", ${quote}pod="$pod"${quote}, "src")`, {}, win)).toBe(
      `label_join(m, "out", ${quote}pod=".*"${quote}, "src")`,
    );
  });

  it.each(['$pod', '${pod}', '${pod:regex}', '[[pod]]', '[[pod:regex]]'])(
    'upgrades an exact matcher for %s',
    (ref) => {
      expect(resolveExpr(`m{pod="${ref}"}`, {}, win)).toBe('m{pod=~".*"}');
    },
  );

  it.each(['[[pod]]', '[[pod:regex]]'])(
    'preserves exact matching for known legacy variable %s',
    (ref) => {
      expect(resolveExpr(`m{pod="${ref}"}`, { pod: 'frontend' }, win)).toBe('m{pod="frontend"}');
    },
  );

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

  it('keeps delimiters and literal dollars inside escaped strings and comments', () => {
    const input = String.raw`label_join(m{path="a\"}#$5",pod="$pod"}, "out", "pod=\"$pod\"}", "src")
or {pod # fake="}" $pod
= # "}" $
"$pod", __name__=~"api.+"}`;
    expect(resolveExpr(input, {}, win)).toBe(
      String.raw`label_join(m{path="a\"}#$5",pod=~".*"}, "out", "pod=\".*\"}", "src")
or {pod # fake="}" .*
=~ # "}" $
".*", __name__=~"api.+"}`,
    );
  });

  it('consumes delimiter characters in formatted references as part of the reference', () => {
    expect(
      resolveExpr(
        '{[[label:}]]="api",prefix_${label:{raw}suffix="$unknown", cash="$$5"}',
        { label: 'job' },
        win,
      ),
    ).toBe('{job="api",prefix_jobsuffix=~".*", cash="$$5"}');
  });

  it('preserves dashboard values that use replacement-string syntax', () => {
    expect(resolveExpr('left/${job:raw}/right', { job: "$`|$'|$$|$&|$1|$<name>" }, win)).toBe(
      'left/left/|/right|$|.*|$1|$<name>/right',
    );
  });

  it('leaves empty formats untouched and substitutes nonempty formats', () => {
    expect(
      resolveExpr(
        'm{a="${label:}",b="[[label:]]",c="${label:raw}",d="[[label:regex]]"}',
        { label: 'api' },
        win,
      ),
    ).toBe('m{a="${label:}",b="[[label:]]",c="api",d="api"}');
  });

  it('substitutes text in an unclosed selector and still finds a later complete selector', () => {
    expect(resolveExpr('m{pod="$pod" or n{pod="$pod"}', {}, win)).toBe(
      'm{pod=".*" or n{pod=~".*"}',
    );
  });

  it('does not interpret braces inside an unclosed string as selectors', () => {
    expect(resolveExpr("m{pod=\"$pod}, or n{zone='$zone'}", {}, win)).toBe(
      "m{pod=\".*}, or n{zone='.*'}",
    );
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

describe('resolveExpr performance', () => {
  const count = 2000;
  const references = '$abcdefghijklmnop'.repeat(count);
  const unclosedBraced = '${label:'.repeat(50_000);
  const unclosedLegacy = '[[label:'.repeat(50_000);
  const malformedValue = `m{job="${unclosedBraced}"}`;

  it.each([
    ['unclosed selector', `m{${references}`, {}, `m{${'.*'.repeat(count)}`],
    ['malformed label matcher', `m{${references}}`, {}, `m{${'.*'.repeat(count)}}`],
    [
      'valid template label and matcher',
      `m{${references}="api",pod="$pod"}`,
      { abcdefghijklmnop: 'x' },
      `m{${'x'.repeat(count)}="api",pod=~".*"}`,
    ],
    [
      'sequence of incomplete braced references after a valid reference',
      `\${label:raw} ${unclosedBraced}$missing`,
      { label: 'api' },
      `api ${unclosedBraced}.*`,
    ],
    [
      'sequence of incomplete legacy references',
      `${unclosedLegacy}$missing`,
      {},
      `${unclosedLegacy}.*`,
    ],
    [
      'sequence of legacy references closed by a single bracket before a valid reference',
      `${unclosedLegacy}] [[label:raw]]`,
      { label: 'api' },
      `${unclosedLegacy}] api`,
    ],
    [
      'quoted matcher value with incomplete formatted references',
      malformedValue,
      { label: 'api|worker' },
      malformedValue,
    ],
  ] as const)('resolves a long %s within a hard timeout', (_, expr, vars, expected) => {
    const moduleUrl = new URL('./promql.ts', import.meta.url).href;
    const child = spawnSync(
      process.execPath,
      [
        '--experimental-strip-types',
        '--input-type=module',
        '--eval',
        `import { readFileSync } from 'node:fs';
         import { resolveExpr } from ${JSON.stringify(moduleUrl)};
         process.stdout.write(resolveExpr(...JSON.parse(readFileSync(0, 'utf8'))));`,
      ],
      {
        input: JSON.stringify([expr, vars, win]),
        encoding: 'utf8',
        timeout: 3000,
        killSignal: 'SIGKILL',
      },
    );
    expect(child.error).toBeUndefined();
    expect(child.status, child.stderr).toBe(0);
    expect(child.stdout).toBe(expected);
  });
});
