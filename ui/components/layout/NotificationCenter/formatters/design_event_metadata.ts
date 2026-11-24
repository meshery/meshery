import { MESHERY_DOCS_URL } from '@/constants/endpoints';

// Design-related events carry the design's identity and a bounded summary,
// never the design itself (#22764). Rows stored earlier, and rows Meshery
// Cloud rewrote to that shape (layer5io/meshery-cloud#6387), may still spell
// the identity in snake_case - beside or instead of the camelCase keys - and
// may lack fields that could not be derived. Read event metadata through these
// helpers rather than indexing the keys directly.

type EventMetadata = Record<string, unknown> | null | undefined;

const isPresent = (value: unknown) => value !== undefined && value !== null && value !== '';

// firstPresent returns the first value that is set, so a camelCase key wins
// over its snake_case fallback only when it actually carries a value.
const firstPresent = <T = unknown>(...values: unknown[]): T | undefined =>
  values.find(isPresent) as T | undefined;

export type DesignEventIdentity = {
  designId?: string;
  designName?: string;
  designVersion?: string;
  designLink?: string;
  viewLink?: string;
};

export const getDesignEventIdentity = (metadata: EventMetadata): DesignEventIdentity => ({
  designId: firstPresent<string>(metadata?.designId, metadata?.design_id),
  designName: firstPresent<string>(metadata?.designName, metadata?.design_name),
  designVersion: firstPresent<string>(metadata?.designVersion),
  designLink: firstPresent<string>(metadata?.designLink),
  viewLink: firstPresent<string>(metadata?.viewLink, metadata?.view_link),
});

// Every key getDesignEventIdentity reads, for views that render the identity
// another way and must not also list it key by key.
export const DESIGN_EVENT_IDENTITY_KEYS = [
  'designId',
  'design_id',
  'designName',
  'design_name',
  'designVersion',
  'designLink',
  'viewLink',
  'view_link',
];

export type DesignEventError = {
  code?: string;
  componentId?: string;
};

export type DesignEventErrorSummary = {
  status?: string;
  errorCount: number;
  errors: DesignEventError[];
  truncated: boolean;
};

// getDesignEventErrorSummary reads the status, errorCount, errors[] and
// truncated fields of a pattern validate, deploy, undeploy or dry-run event.
// It returns null for an event that carries none of them (a pre-#22764 event,
// or a rewritten dry-run whose stored response was empty). errors[] holds at
// most 20 entries; errorCount is the full count. A rewritten deploy error may
// lack componentId, and a rewritten validate error's code is the JSON-schema
// validator keyword (`required`, `type`, ...) rather than a MeshKit code.
export const getDesignEventErrorSummary = (
  metadata: EventMetadata,
): DesignEventErrorSummary | null => {
  const hasStatus = isPresent(metadata?.status);
  const hasErrors = Array.isArray(metadata?.errors);
  if (!hasStatus && !hasErrors) {
    return null;
  }

  const errors = (hasErrors ? (metadata?.errors as unknown[]) : [])
    .filter((error) => error && typeof error === 'object')
    .map((error: Record<string, unknown>) => ({
      code: isPresent(error.code) ? String(error.code) : undefined,
      componentId: isPresent(error.componentId) ? String(error.componentId) : undefined,
    }));
  const errorCount = isPresent(metadata?.errorCount) ? Number(metadata?.errorCount) : NaN;

  return {
    status: hasStatus ? String(metadata?.status) : undefined,
    errorCount: Number.isFinite(errorCount) ? errorCount : errors.length,
    errors,
    truncated: metadata?.truncated === true,
  };
};

const MESHKIT_ERROR_CODE = /^(meshery-[a-z]+-\d+|meshkit-\d+)$/i;

// isMeshkitErrorCode tells a MeshKit error code (meshery-server-1480, or
// meshkit-11001 for a MeshKit library error), which has an entry in the error
// code reference, from any other code an event error may carry: a JSON-schema
// keyword or a Kubernetes cause type.
export const isMeshkitErrorCode = (code: unknown): code is string =>
  typeof code === 'string' && MESHKIT_ERROR_CODE.test(code);

export const meshkitErrorDocLink = (code: string) =>
  `${MESHERY_DOCS_URL}/reference/error-codes#${code}`;
