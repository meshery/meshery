import type { MeshkitError, MeshkitFetchBaseQueryError } from '@meshery/schemas/api';

/**
 * Result of formatting an RTK Query error for the user-facing notification
 * layer. The notification system renders `message` through `BasicMarkdown`
 * (see `ThemeResponsiveSnackbar` in `themes/App.styles.tsx`), which means the
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
    return {
      message: formatMeshkitErrorMarkdown(error.meshkit, fallbackTitle),
      meshkit: error.meshkit,
    };
  }

  const fallback = extractFallbackMessage(error) ?? fallbackTitle ?? 'An unexpected error occurred';
  return { message: fallback };
};
