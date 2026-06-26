import type { Datasource, TemplateVar } from './types';

// resolveExpr (PromQL macro + variable substitution) is product-agnostic and
// lives in ../common/promql; re-exported here so existing imports keep working.
export { resolveExpr } from '../common/promql';

// Grafana's "all" sentinel for template variables.
const ALL_VALUE = '$__all';

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
