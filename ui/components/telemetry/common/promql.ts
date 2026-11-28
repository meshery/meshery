import type { TimeWindow } from './types';

export const escapeRegExp = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

const quotedString = /"(?:\\.|[^"\\])*"|'(?:\\.|[^'\\])*'|`[^`]*`/.source;
const identifier = '[A-Za-z_]\\w*';
const referencePrefix = (name: string) =>
  `\\$\\{(?<braced>${name})(?=[:}])|\\[\\[(?<legacy>${name})(?=:|\\]\\])|\\$(?<bare>${name})\\b`;
const unknownReference = new RegExp(referencePrefix(`(?!__)${identifier}`), 'g');
const metricName = /^[a-zA-Z_:][\w:]*$/;
const wordOperator = /^(?:and|or|unless|atan2|bitand|bitor|bitxor)$/i;
const modifier = /^(?:bool|on|ignoring|group_left|group_right)$/i;
const argumentModifier = /^(?:on|ignoring|group_left|group_right|fill(?:_left|_right)?)$/i;
// Consume complete references and strings before punctuation. An unmatched
// quote owns the rest of the input, and a literal dollar sign falls through to \S.
const tokenPrefix = `${referencePrefix(identifier)}|(?<quoted>${quotedString})|["'\x60][\\s\\S]*|#[^\\n]*`;
const expressionToken = new RegExp(`${tokenPrefix}|[\\w:.]+|\\S`, 'g');
const matcherToken = new RegExp(`${tokenPrefix}|\\w+|=~|!~|!=|\\S`, 'g');

interface Token {
  text: string;
  start: number;
  end: number;
  kind: 'template' | 'quoted' | 'comment' | 'plain';
  name?: string;
}

function* scanTokens(
  input: string,
  pattern: RegExp,
  keepUnmatched = false,
): Generator<Token, void> {
  const re = new RegExp(pattern);
  // Cache the next closing character for each format syntax. Many incomplete
  // prefixes can share the same terminator or the same unterminated suffix.
  const nextClose = { '}': -1, ']': -1 };
  for (let match = re.exec(input); match; match = re.exec(input)) {
    const { braced, legacy, bare, quoted } = match.groups ?? {};
    const start = match.index;
    let end = re.lastIndex;
    const closer = braced !== undefined ? '}' : legacy !== undefined ? ']' : undefined;
    if (closer) {
      const terminator = closer === '}' ? '}' : ']]';
      if (input[end] === ':') {
        const formatStart = end + 1;
        if (nextClose[closer] < formatStart) {
          const index = input.indexOf(closer, formatStart);
          nextClose[closer] = index === -1 ? input.length : index;
        }
        end = nextClose[closer];
        // Legacy formats end at the first ']', which must be part of ']]'.
        if (end === formatStart || !input.startsWith(terminator, end)) {
          re.lastIndex = start + 1;
          if (keepUnmatched) {
            yield { text: input.slice(start, start + 1), start, end: start + 1, kind: 'plain' };
          }
          continue;
        }
      }
      end += terminator.length;
      re.lastIndex = end;
    }
    const text = input.slice(start, end);
    const name = braced ?? legacy ?? bare;
    if (name !== undefined) {
      yield { text, start, end, kind: 'template', name };
    } else {
      yield { text, start, end, kind: quoted ? 'quoted' : text[0] === '#' ? 'comment' : 'plain' };
    }
  }
}

function replaceReferences(input: string, pattern: RegExp, value: string): string {
  let result = '';
  let cursor = 0;
  for (const { text, start, end } of scanTokens(input, pattern)) {
    // Keep String.replace's replacement-string semantics for dashboard values.
    const replacement = value.replace(/\$([$&\x60'])/g, (_, key) => {
      switch (key) {
        case '&':
          return text;
        case '`':
          return input.slice(0, start);
        case "'":
          return input.slice(end);
        default:
          return '$';
      }
    });
    result += input.slice(cursor, start) + replacement;
    cursor = end;
  }
  return result + input.slice(cursor);
}

function skipComments(tokens: Token[], index: number): number {
  while (tokens[index]?.kind === 'comment') index++;
  return index;
}

function closingToken(tokens: Token[], index: number, open: string, close: string): number {
  // Selectors and modifier argument lists do not nest. Stop at another opening
  // delimiter so an unmatched opener cannot repeatedly scan the remaining input.
  for (let i = index + 1; i < tokens.length; i++) {
    const text = tokens[i]?.text;
    if (text === open) break;
    if (text === close) return i;
  }
  return -1;
}

function* scanMatchers(body: string) {
  const tokens = Array.from(scanTokens(body, matcherToken, true));
  for (let i = 0; i < tokens.length; i++) {
    const token = tokens[i];
    if (!token) break;
    let labelEnd = token.end;
    const quotedLabel = token.kind === 'quoted';
    const bareLabel = token.kind === 'template' || /^[a-zA-Z_]/.test(token.text);
    if (bareLabel) {
      // A label can contain adjacent words and template references. Consume
      // each whole token once, even if no operator follows the label.
      for (let next = tokens[i + 1]; next?.start === labelEnd; next = tokens[i + 1]) {
        if (next.kind !== 'template' && !/^\w+$/.test(next.text)) break;
        labelEnd = next.end;
        i++;
      }
    }
    if (bareLabel || quotedLabel) {
      const operatorIndex = skipComments(tokens, i + 1);
      const operator = tokens[operatorIndex];
      if (operator && /^(?:=~|!~|!=|=)$/.test(operator.text)) {
        const valueIndex = skipComments(tokens, operatorIndex + 1);
        const value = tokens[valueIndex];
        if (value?.kind === 'quoted') {
          yield {
            start: token.start,
            end: value.end,
            label: body.slice(token.start, operator.start),
            operator: operator.text,
            space: body.slice(operator.end, value.start),
            quoted: value.text,
          };
          i = valueIndex;
          continue;
        }
      }
    }
    if (quotedLabel || token.kind === 'comment') {
      yield { start: token.start, end: token.end, quotedMetric: quotedLabel };
    }
  }
}

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

  const known = Object.entries(varValues).map(([name, value]) => ({
    ref: new RegExp(referencePrefix(escapeRegExp(name)), 'g'),
    value,
  }));
  const substitute = (input: string) => {
    for (const { ref, value } of known) input = replaceReferences(input, ref, value);
    return replaceReferences(input, unknownReference, '.*');
  };

  // Scan brace-delimited selectors while skipping strings and comments outside them.
  let resolved = '';
  let cursor = 0;
  let expectsOperand = true;
  const tokens = Array.from(scanTokens(out, expressionToken, true));
  for (let i = 0; i < tokens.length; i++) {
    const token = tokens[i];
    if (!token) break;
    if (token.kind === 'comment') continue;
    const { text, start } = token;
    const next = skipComments(tokens, i + 1);
    const hasMetric = token.kind === 'template' || metricName.test(text);
    const openIndex = hasMetric ? next : i;
    const open = tokens[openIndex];
    const closeIndex = open?.text === '{' ? closingToken(tokens, openIndex, '{', '}') : -1;
    const close = tokens[closeIndex];
    if (!open || !close) {
      if (argumentModifier.test(text) && tokens[next]?.text === '(') {
        const argsEnd = closingToken(tokens, next, '(', ')');
        if (argsEnd !== -1) {
          expectsOperand = true;
          i = argsEnd;
          continue;
        }
      }
      const resolvedToken = token.kind === 'template' ? substitute(text) : text;
      expectsOperand =
        modifier.test(resolvedToken) ||
        /^[,(+*/%^<>=!@-]$/.test(resolvedToken) ||
        (!expectsOperand && wordOperator.test(resolvedToken));
      continue;
    }

    resolved += substitute(out.slice(cursor, start));
    cursor = close.end;
    i = closeIndex;
    const body = out.slice(open.end, close.start);

    // The grammar also accepts set operators as metric identifiers in operand position.
    const resolvedMetric = hasMetric ? substitute(text) : '';
    let constrained =
      metricName.test(resolvedMetric) &&
      !modifier.test(resolvedMetric) &&
      !/^atan2$/i.test(resolvedMetric) &&
      (expectsOperand || !wordOperator.test(resolvedMetric));
    expectsOperand = false;
    let upgradedMatcher = false;
    const parts: string[] = [];
    let bodyCursor = 0;
    for (const entry of scanMatchers(body)) {
      parts.push(substitute(body.slice(bodyCursor, entry.start)));
      bodyCursor = entry.end;
      if (entry.operator === undefined) {
        // A standalone quoted metric name inside braces also constrains the selector.
        const resolvedToken = substitute(body.slice(entry.start, entry.end));
        constrained ||= entry.quotedMetric && resolvedToken.length > 2;
        parts.push(resolvedToken);
        continue;
      }

      const { label, operator, space, quoted } = entry;
      const value = quoted.slice(1, -1);
      const ref = operator === '=' ? scanTokens(value, unknownReference).next().value : undefined;
      const unresolved =
        ref?.name !== undefined &&
        ref.start === 0 &&
        ref.end === value.length &&
        !Object.hasOwn(varValues, ref.name);

      let resolvedOperator = operator;

      if (unresolved) {
        resolvedOperator = '=~';
      } else if (
        (operator === '=' || operator === '!=') &&
        known.some(
          ({ ref, value }) =>
            (value === '.*' || value.includes('|')) && !scanTokens(quoted, ref).next().done,
        )
      ) {
        // Match-all and multi-value dashboard variables require a regex matcher.
        resolvedOperator = operator === '=' ? '=~' : '!~';
      }

      upgradedMatcher ||= resolvedOperator !== operator;

      const resolvedValue = substitute(quoted);
      constrained ||= excludesEmpty(resolvedOperator, resolvedValue);
      parts.push(substitute(label) + resolvedOperator + substitute(space) + resolvedValue);
    }

    // PromQL requires an unnamed selector to reject the empty string somewhere.
    // Constrain the metric name so =~".*" can still match missing labels. With no
    // other selective matcher this can query all named series.
    const nameConstraint = !constrained && upgradedMatcher ? '__name__!="",' : '';

    resolved +=
      resolvedMetric +
      substitute(out.slice(hasMetric ? token.end : start, open.start)) +
      '{' +
      nameConstraint +
      parts.join('') +
      substitute(body.slice(bodyCursor)) +
      '}';
  }
  return resolved + substitute(out.slice(cursor));
}
