import type { TimeWindow } from './types';

export const escapeRegExp = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

const quotedString = /"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|`[^`]*`/.source;
const comment = /#[^\n]*(?=\n|$)/.source;
const trivia = `(?:\\s|${comment})*`;
const template = /\$\{[A-Za-z_]\w*(?::[^}]+)?\}|\$[A-Za-z_]\w*|\[\[[A-Za-z_]\w*(?::[^\]]+)?\]\]/
  .source;
const wordOperator = /^(?:and|or|unless|atan2|bitand|bitor|bitxor)$/i;
const modifier = /^(?:bool|on|ignoring|group_left|group_right)$/i;
const selectorOrToken = new RegExp(
  `(?:(${template}|[a-zA-Z_:][\\w:]*)(${trivia}))?\\{((?:${quotedString}|${comment}|[^{}"'\x60#])*)\\}` +
    `|((?:on|ignoring|group_left|group_right|fill(?:_left|_right)?)${trivia}\\((?:${quotedString}|${comment}|[^()"'\x60#])*\\))` +
    `|${quotedString}|${comment}|(${template})|[\\w:.]+|\\S`,
  'gi',
);
const matcher = new RegExp(
  `((?:[a-zA-Z_]\\w*|${quotedString})${trivia})(=~|!~|!=|=)(${trivia})(${quotedString})` +
    `|${comment}|(${quotedString})`,
  'g',
);

function excludesEmpty(operator: string, quoted: string): boolean {
  const value = quoted.slice(1, -1);
  if (operator === '=') return value !== '';
  if (operator === '!=') return value === '';

  // Unquote PromQL escapes before testing the fully anchored pattern. JavaScript
  // is only an approximation of RE2. Unrecognized patterns cannot prove a constraint.
  try {
    const pattern =
      quoted[0] === '`'
        ? value
        : value.replace(
            /\\(?:([0-7]{3})|x([\da-fA-F]{2})|u([\da-fA-F]{4})|U([\da-fA-F]{8})|(.))/g,
            (escape, octal, byte, short, long, char) => {
              if (octal) return String.fromCodePoint(parseInt(octal, 8));
              if (byte || short || long)
                return String.fromCodePoint(parseInt(byte || short || long, 16));
              const escapes: Record<string, string> = {
                a: '\x07',
                b: '\b',
                f: '\f',
                n: '\n',
                r: '\r',
                t: '\t',
                v: '\v',
              };
              return escapes[char] ?? (char === '\\' || char === quoted[0] ? char : escape);
            },
          );
    const matchesEmpty = new RegExp(`^(?:${pattern})$`, 'u').test('');
    return operator === '=~' ? !matchesEmpty : matchesEmpty;
  } catch {
    return false;
  }
}

/**
 * resolveExpr substitutes Grafana global time macros and (optionally) dashboard
 * template variables into a PromQL expression so it can be executed directly.
 * Any remaining (unknown) `$var` token falls back to `.*` for a best-effort
 * render. A positive exact matcher uses `=~` when its entire value is an unknown variable.
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

  const variable =
    /\$(?:\{(?!__)([A-Za-z_]\w*)(?::[^}]+)?\}|(?!__)([A-Za-z_]\w*))|\[\[(?!__)([A-Za-z_]\w*)(?::[^\]]+)?\]\]/;
  const known = Object.entries(varValues).map(([name, value]) => {
    const n = escapeRegExp(name);
    return {
      ref: new RegExp(`\\$\\{${n}(?::[^}]+)?\\}|\\[\\[${n}(?::[^\\]]+)?\\]\\]|\\$${n}\\b`, 'g'),
      value,
    };
  });
  const unknown = new RegExp(variable.source, 'g');
  const substitute = (input: string) => {
    for (const { ref, value } of known) input = input.replace(ref, value);
    return input.replace(unknown, '.*');
  };

  // Scan brace-delimited selectors while skipping strings and comments outside them.
  let resolved = '';
  let cursor = 0;
  let expectsOperand = true;
  for (const token of out.matchAll(selectorOrToken)) {
    const [text, metric, spacing, body, modifierWithArguments, variableToken] = token;
    if (body === undefined) {
      if (!text.startsWith('#')) {
        const resolvedToken = variableToken === undefined ? text : substitute(variableToken);
        expectsOperand =
          !!modifierWithArguments ||
          modifier.test(resolvedToken) ||
          /^[,(+*/%^<>=!@-]$/.test(resolvedToken) ||
          (!expectsOperand && wordOperator.test(resolvedToken));
      }
      continue;
    }
    resolved += substitute(out.slice(cursor, token.index));
    cursor = token.index + text.length;

    // The grammar also accepts set operators as metric identifiers in operand position.
    const resolvedMetric = substitute(metric ?? '');
    let constrained =
      /^[a-zA-Z_:][\w:]*$/.test(resolvedMetric) &&
      !modifier.test(resolvedMetric) &&
      !/^atan2$/i.test(resolvedMetric) &&
      (expectsOperand || !wordOperator.test(resolvedMetric));
    expectsOperand = false;
    let firstUpgrade: number | undefined;
    const parts: string[] = [];
    let bodyCursor = 0;
    for (const entry of body.matchAll(matcher)) {
      const [original, label, operator, space, quoted, quotedMetric] = entry;
      parts.push(substitute(body.slice(bodyCursor, entry.index)));
      bodyCursor = entry.index + original.length;
      if (
        label === undefined ||
        operator === undefined ||
        space === undefined ||
        quoted === undefined
      ) {
        // A standalone quoted metric name inside braces also constrains the selector.
        constrained ||= !!quotedMetric && substitute(quotedMetric).length > 2;
        parts.push(substitute(original));
        continue;
      }

      const value = quoted.slice(1, -1);
      const ref = variable.exec(value);
      const name = ref?.[1] ?? ref?.[2] ?? ref?.[3];
      const unresolved =
        operator === '=' && ref?.[0] === value && name && !Object.hasOwn(varValues, name);
      let resolvedOperator = operator;
      if (unresolved) {
        resolvedOperator = '=~';
      } else if (
        (operator === '=' || operator === '!=') &&
        known.some(
          ({ ref, value }) => (value === '.*' || value.includes('|')) && quoted.search(ref) !== -1,
        )
      ) {
        // Match-all and multi-value dashboard variables require a regex matcher.
        resolvedOperator = operator === '=' ? '=~' : '!~';
      }
      const resolvedValue = substitute(quoted);
      constrained ||= excludesEmpty(resolvedOperator, resolvedValue);
      parts.push(substitute(label) + resolvedOperator + substitute(space));
      if (unresolved && firstUpgrade === undefined) firstUpgrade = parts.length;
      parts.push(resolvedValue);
    }
    // A metric-less selector needs one matcher that rejects the empty string.
    if (!constrained && firstUpgrade !== undefined) {
      const quote = parts[firstUpgrade]?.charAt(0);
      if (quote) parts[firstUpgrade] = `${quote}.+${quote}`;
    }
    resolved +=
      resolvedMetric +
      substitute(spacing ?? '') +
      '{' +
      parts.join('') +
      substitute(body.slice(bodyCursor)) +
      '}';
  }
  return resolved + substitute(out.slice(cursor));
}
