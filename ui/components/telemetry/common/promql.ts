import type { TimeWindow } from './types';

export const escapeRegExp = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

/**
 * resolveExpr substitutes Grafana global time macros and (optionally) dashboard
 * template variables into a PromQL expression so it can be executed directly.
 * Any remaining (unknown) `$var` token falls back to `.*` for a best-effort
 * render. With an empty varValues map it performs macro expansion only — which
 * is the common case for hand-written Prometheus panel queries.
 */
export function resolveExpr(
  expr: string,
  varValues: Record<string, string>,
  window: TimeWindow,
): string {
  if (!expr) return expr;
  const step = Math.max(1, window.step);
  const rangeSec = Math.max(1, window.end - window.start);
  const rateInterval = Math.max(step * 4, step + 15);

  let out = expr;

  // Global time macros (must run before generic variable substitution).
  const macros: Record<string, string> = {
    $__rate_interval_ms: String(rateInterval * 1000),
    $__rate_interval: `${rateInterval}s`,
    $__interval_ms: String(step * 1000),
    $__interval: `${step}s`,
    $__range_ms: String(rangeSec * 1000),
    $__range_s: String(rangeSec),
    $__range: `${rangeSec}s`,
  };
  for (const [macro, value] of Object.entries(macros)) {
    out = out
      .split(`${macro}`)
      .join(value)
      .split(`\${${macro.slice(1)}}`)
      .join(value);
  }

  // A variable that resolves to a regex — the match-all `.*` (unset / "All") or a
  // multi-value alternation `a|b` — is only meaningful inside a regex matcher.
  // Dashboards frequently use an exact matcher (e.g. `cluster="$cluster"`), which
  // would become `cluster=".*"` and match the literal string `.*` (returning no
  // data). Upgrade `=`/`!=` to `=~`/`!~` for matchers referencing such variables
  // before substituting, so the query returns data as Grafana would.
  for (const [name, value] of Object.entries(varValues)) {
    if (value !== '.*' && !value.includes('|')) continue; // concrete single value: leave `=` as-is
    const n = escapeRegExp(name);
    const ref = `(?:\\$\\{${n}(?::[^}]+)?\\}|\\[\\[${n}(?::[^\\]]+)?\\]\\]|\\$${n}\\b)`;
    const exactMatcher = new RegExp(`([a-zA-Z_]\\w*)\\s*(!?=)(?!~)\\s*("[^"]*${ref}[^"]*")`, 'g');
    out = out.replace(
      exactMatcher,
      (_m, label, op, quoted) => `${label}${op === '=' ? '=~' : '!~'}${quoted}`,
    );
  }

  // Known template variables: ${name}, $name, [[name]] / [[name:format]].
  for (const [name, value] of Object.entries(varValues)) {
    const n = escapeRegExp(name);
    const re = new RegExp(`\\$\\{${n}(?::[^}]+)?\\}|\\[\\[${n}(?::[^\\]]+)?\\]\\]|\\$${n}\\b`, 'g');
    out = out.replace(re, value);
  }

  // Any leftover non-macro variable -> match-all, so the query still runs.
  out = out.replace(/\$\{(?!__)[A-Za-z_]\w*(?::[^}]+)?\}|\$(?!__)[A-Za-z_]\w*/g, '.*');

  return out;
}
