import type { MeshkitError, MeshkitFetchBaseQueryError } from '@meshery/schemas/api';

/**
 * Result of formatting an RTK Query error for the user-facing notification
 * layer. The notification system renders `message` through `BasicMarkdown`
 * (see `ThemeResponsiveSnackbar` in `theme/snackbar.tsx`), which means the
 * formatter can use markdown — bold, lists, italics — to surface multiple
 * lines of MeshKit metadata in a single toast.
 */
export interface FormattedApiError {
  /** Markdown string ready for `BasicMarkdown`. */
  message: string;
  /** The MeshKit metadata block, if the response carried one. */
  meshkit?: MeshkitError;
}

/**
 * Type-guard that narrows an unknown RTK Query error to one carrying a
 * MeshKit envelope. The schemas package's `transformErrorResponse` wrapper
 * (v1.2.2+) sets `error.meshkit` on every non-2xx JSON response that matches
 * the MeshKit shape, so this guard can be used directly on the `error`
 * returned by mutation/query hooks without manual casting.
 */
export const hasMeshkitError = (
  error: unknown,
): error is MeshkitFetchBaseQueryError & { meshkit: MeshkitError } => {
  if (!error || typeof error !== 'object') return false;
  const candidate = (error as { meshkit?: unknown }).meshkit;
  return (
    !!candidate &&
    typeof candidate === 'object' &&
    typeof (candidate as MeshkitError).message === 'string'
  );
};

/**
 * Render a MeshKit error to a markdown string suitable for the snackbar.
 *
 * Layout:
 *   **<message>**            ← bold title
 *   *Try:*
 *   - <remediation 1>        ← bullet list
 *   - <remediation 2>
 *   `<code>`                 ← muted reference for support tickets
 *
 * Sections are emitted only when their source data is non-empty, so a
 * minimally-populated MeshKit response (just `message`) renders as a single
 * bold line and degrades gracefully.
 */
export const formatMeshkitErrorMarkdown = (
  meshkit: MeshkitError,
  fallbackMessage?: string,
): string => {
  const lines: string[] = [];
  const title = meshkit.message?.trim() || fallbackMessage?.trim() || 'Request failed';
  lines.push(`**${title}**`);

  const remediations = (meshkit.suggestedRemediation ?? []).filter(
    (line) => typeof line === 'string' && line.trim().length > 0,
  );
  if (remediations.length > 0) {
    lines.push('');
    lines.push('*Try:*');
    for (const remediation of remediations) {
      lines.push(`- ${remediation.trim()}`);
    }
  }

  if (meshkit.code) {
    lines.push('');
    lines.push(`\`${meshkit.code}\``);
  }

  return lines.join('\n');
};

type MeshkitDetailKey = 'probableCause' | 'suggestedRemediation' | 'longDescription';

/**
 * Wire aliases for each MeshKit detail array, most-preferred spelling first.
 *
 * The Meshery server emits camelCase - see the `errorResponse` struct in
 * `server/models/httputil/httputil.go` and the identifier-naming contract in
 * AGENTS.md - but the `withMeshkitError` baseQuery wrapper in
 * `@meshery/schemas` reads the snake_case spellings (`probable_cause`,
 * `suggested_remediation`, `long_description`) instead. The result is an
 * `error.meshkit` carrying only message/code/severity, so the "*Try:*" section
 * below never renders even though the server sent remediations.
 *
 * Tracked upstream as meshery/schemas#1081
 * (https://github.com/meshery/schemas/issues/1081). The defect is present in
 * the pinned 1.3.x line and in 1.4.0. Once it is fixed upstream and this
 * repo's dependency is bumped past it, the raw-body fallback in
 * `resolveMeshkitError` becomes redundant and can be deleted.
 */
const MESHKIT_DETAIL_ALIASES: Record<MeshkitDetailKey, readonly string[]> = {
  probableCause: ['probableCause', 'probable_cause'],
  suggestedRemediation: ['suggestedRemediation', 'suggested_remediation'],
  longDescription: ['longDescription', 'long_description'],
};

/** Narrow an unknown value to a non-empty array of non-blank strings. */
const toStringArray = (value: unknown): string[] | undefined => {
  if (!Array.isArray(value)) return undefined;
  const entries = value.filter(
    (entry): entry is string => typeof entry === 'string' && entry.trim().length > 0,
  );
  return entries.length > 0 ? entries : undefined;
};

/**
 * Merge the MeshKit envelope attached by the schemas baseQuery with the detail
 * arrays still present on the raw response body at `error.data`. Values already
 * on `meshkit` win; within the body camelCase wins over snake_case, so a future
 * producer emitting either spelling keeps working.
 */
const resolveMeshkitError = (error: { meshkit: MeshkitError; data?: unknown }): MeshkitError => {
  const body =
    error.data && typeof error.data === 'object' && !Array.isArray(error.data)
      ? (error.data as Record<string, unknown>)
      : undefined;
  if (!body) return error.meshkit;

  const resolved: MeshkitError = { ...error.meshkit };
  for (const key of Object.keys(MESHKIT_DETAIL_ALIASES) as MeshkitDetailKey[]) {
    if (toStringArray(resolved[key])) continue;
    for (const alias of MESHKIT_DETAIL_ALIASES[key]) {
      const fromBody = toStringArray(body[alias]);
      if (fromBody) {
        resolved[key] = fromBody;
        break;
      }
    }
  }
  return resolved;
};

/**
 * Extract a single best-effort string from a non-MeshKit RTK Query error or
 * any thrown value. Mirrors the legacy `error?.data` / `error?.message`
 * pattern that scattered through the codebase before the MeshKit envelope
 * landed, so callers that fall through this branch keep their previous UX.
 */
const extractFallbackMessage = (error: unknown): string | undefined => {
  if (typeof error === 'string') return error;
  if (!error || typeof error !== 'object') return undefined;

  const obj = error as Record<string, unknown>;

  // RTK FetchBaseQueryError shape: { status, data }
  if (typeof obj.data === 'string' && obj.data.length > 0) return obj.data;
  if (obj.data && typeof obj.data === 'object') {
    const nested = obj.data as Record<string, unknown>;
    if (typeof nested.error === 'string') return nested.error;
    if (typeof nested.message === 'string') return nested.message;
  }

  // SerializedError / generic Error shape: { message }
  if (typeof obj.message === 'string') return obj.message;
  if (typeof obj.error === 'string') return obj.error;

  return undefined;
};

/**
 * Format any RTK Query error (with or without MeshKit metadata) into a
 * markdown string and the originating MeshKit block. When `meshkit` is
 * absent the result is a single-line message — identical to the pre-v1.2.2
 * UX — so this function is safe to call unconditionally on every error.
 *
 * @param error          The error returned by `useMutation()`/`useQuery()` or
 *                       caught from `.unwrap()`. Typed as `unknown` because
 *                       call sites disagree on the wire type.
 * @param fallbackTitle  Human-friendly description of the failed operation
 *                       (e.g. "Failed to create workspace") used when the
 *                       wire payload doesn't carry a `message`.
 */
export const formatApiError = (error: unknown, fallbackTitle?: string): FormattedApiError => {
  if (hasMeshkitError(error)) {
    const meshkit = resolveMeshkitError(error);
    return {
      message: formatMeshkitErrorMarkdown(meshkit, fallbackTitle),
      meshkit,
    };
  }

  const fallback = extractFallbackMessage(error) ?? fallbackTitle ?? 'An unexpected error occurred';
  return { message: fallback };
};
