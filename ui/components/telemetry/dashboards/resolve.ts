import type { Datasource, TemplateVar, TimeWindow } from './types';

// Grafana's "all" sentinel for template variables.
const ALL_VALUE = '$__all';

const escapeRegExp = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

/**
 * buildVarValues resolves each dashboard template variable to the string that
 * should be substituted into a panel query, using the variable's currently
 * selected value(s). "All" / unset variables resolve to a regex-matching `.*`
 * so they work inside `=~"..."` label matchers. Datasource variables are
 * excluded — those are resolved separately to a concrete datasource UID.
 */
export function buildVarValues(templateVars: TemplateVar[] = []): Record<string, string> {
  const values: Record<string, string> = {};
  for (const v of templateVars) {
    if (v.type === 'datasource') continue;
    const current = (v.current ?? []).filter((c) => c && c !== ALL_VALUE);
    if (current.length === 0) {
      // unset or "All" -> match everything
      values[v.name] = '.*';
    } else if (current.length === 1) {
      values[v.name] = current[0];
    } else {
      // multi-value -> regex alternation for `=~` matchers
      values[v.name] = current.join('|');
    }
  }
  return values;
}

/**
 * resolveExpr substitutes Grafana global time macros and dashboard template
 * variables into a PromQL expression so it can be executed directly. Any
 * remaining (unknown) `$var` token falls back to `.*` for a best-effort render.
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

const varName = (ref?: string): string | null => {
  if (!ref) return null;
  const m = ref.match(/^\$\{?(\w+)\}?$/);
  return m ? m[1] : null;
};

const pickDefault = (datasources: Datasource[]): string => {
  if (datasources.length === 0) return '';
  const def = datasources.find((d) => d.isDefault);
  if (def) return def.uid;
  const prom = datasources.find((d) => d.type === 'prometheus');
  return (prom ?? datasources[0]).uid;
};

const matchDatasource = (value: string, datasources: Datasource[]): string => {
  if (!value || value === 'default') return pickDefault(datasources);
  const byUid = datasources.find((d) => d.uid === value);
  if (byUid) return byUid.uid;
  const byName = datasources.find((d) => d.name === value);
  if (byName) return byName.uid;
  return value; // unknown but concrete — let the backend try
};

/**
 * resolveDatasourceUid turns a panel/target datasource reference (which may be a
 * concrete UID, a `$datasource` variable, a name, or empty) into a concrete
 * datasource UID, falling back to the default datasource.
 */
export function resolveDatasourceUid(
  ref: string | undefined,
  templateVars: TemplateVar[] = [],
  datasources: Datasource[] = [],
): string {
  const name = varName(ref);
  if (name) {
    const dsVar = templateVars.find((v) => v.name === name && v.type === 'datasource');
    const current = (dsVar?.current ?? []).filter((c) => c && c !== ALL_VALUE);
    if (current.length > 0) {
      return matchDatasource(current[0], datasources);
    }
    return pickDefault(datasources);
  }
  return matchDatasource(ref ?? '', datasources);
}
