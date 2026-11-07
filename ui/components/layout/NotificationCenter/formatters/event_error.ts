import { flattenProbableCauses, effectiveActionPolicy } from '@meshery/schemas';
import type { v1beta3 } from '@meshery/schemas';

// An error event's metadata.error comes in three shapes
// (docs/error-definition-contract.md in meshery/schemas):
//
//   - the schemas EventError: `code`, rendered `shortDescription` and
//     `longDescription`, `probableCauses` nested with their remediations,
//     `matchedCauses`, the flat lists derived from the nesting, `referenceUrl`;
//   - the legacy MeshKit error as Go serialized it, on rows recorded before
//     EventError: `Code`, `LongDescription`, `ProbableCause` and
//     `SuggestedRemediation` as lists of lines, never paired by position;
//   - a design event's DesignEventOperationError: `code` and `message`.
//
// normalizeEventError reads any of them into one view model, so the
// notification center renders them alike and falls back to the flat lists
// whenever the nesting is absent.

type EventError = v1beta3.EventError;
type ProbableCause = v1beta3.ProbableCause;
type Remediation = v1beta3.Remediation;
type ActionPolicy = v1beta3.ActionPolicy;

export type ErrorRemediation = {
  id: string;
  description: string;
  action?: string;
  // The policy a client acts on: manual when omitted or unknown, so nothing
  // is ever run by mistake. The notification center runs no action; it shows
  // the policy so a person knows what a client may do.
  actionPolicy: ActionPolicy;
};

export type ErrorCause = {
  id: string;
  description: string;
  remediations: ErrorRemediation[];
  // Whether the producer named this cause as likely for this occurrence, and
  // how sure it was (0 to 1) when it said.
  matched: boolean;
  confidence?: number;
};

export type NormalizedEventError = {
  code?: string;
  shortDescription?: string;
  // One entry per line or paragraph.
  longDescription: string[];
  // Nested causes, the matched ones first in the producer's order, then the
  // rest in the definition's order. Empty for an unpaired error.
  causes: ErrorCause[];
  // The flat lists, for an unpaired error (or one whose producer sent only
  // them). Never paired by position.
  probableCause: string[];
  suggestedRemediation: string[];
  referenceUrl?: string;
};

const isRecord = (value: unknown): value is Record<string, unknown> =>
  !!value && typeof value === 'object' && !Array.isArray(value);

const nonEmptyString = (value: unknown): string | undefined =>
  typeof value === 'string' && value.trim() !== '' ? value : undefined;

// MeshKit's legacy lists hold "None" or "" for an absent value.
const lines = (value: unknown): string[] => {
  const list = Array.isArray(value) ? value : value === undefined || value === null ? [] : [value];
  return list
    .filter((line): line is string => typeof line === 'string')
    .map((line) => line.trim())
    .filter((line) => line !== '' && line !== 'None');
};

const toRemediation = (remediation: Remediation): ErrorRemediation => ({
  id: remediation.id,
  description: remediation.description,
  action: nonEmptyString(remediation.action),
  actionPolicy: effectiveActionPolicy(remediation),
});

const isProbableCause = (value: unknown): value is ProbableCause =>
  isRecord(value) &&
  typeof value.id === 'string' &&
  typeof value.description === 'string' &&
  Array.isArray(value.remediations);

const orderCauses = (
  causes: ProbableCause[],
  matched: v1beta3.EventMatchedCause[],
): ErrorCause[] => {
  const byId = new Map(causes.map((cause) => [cause.id, cause]));
  const seen = new Set<string>();
  const ordered: ErrorCause[] = [];
  for (const match of matched) {
    const cause = byId.get(match.causeId);
    if (!cause || seen.has(cause.id)) continue;
    seen.add(cause.id);
    ordered.push({
      id: cause.id,
      description: cause.description,
      remediations: cause.remediations.map(toRemediation),
      matched: true,
      confidence: typeof match.confidence === 'number' ? match.confidence : undefined,
    });
  }
  for (const cause of causes) {
    if (seen.has(cause.id)) continue;
    ordered.push({
      id: cause.id,
      description: cause.description,
      remediations: cause.remediations.map(toRemediation),
      matched: false,
    });
  }
  return ordered;
};

const normalizeSchemaEventError = (error: Partial<EventError> & Record<string, unknown>) => {
  const nested = Array.isArray(error.probableCauses)
    ? error.probableCauses.filter(isProbableCause)
    : [];
  const matched = Array.isArray(error.matchedCauses)
    ? error.matchedCauses.filter(
        (match): match is v1beta3.EventMatchedCause =>
          isRecord(match) && typeof match.causeId === 'string',
      )
    : [];
  const derived = nested.length > 0 ? flattenProbableCauses(nested) : undefined;
  return {
    code: nonEmptyString(error.code),
    shortDescription: nonEmptyString(error.shortDescription),
    // A DesignEventOperationError carries its text as `message`.
    longDescription: lines(error.longDescription ?? error.message),
    causes: orderCauses(nested, matched),
    probableCause: derived ? derived.probableCause : lines(error.probableCause),
    suggestedRemediation: derived
      ? derived.suggestedRemediation
      : lines(error.suggestedRemediation),
    referenceUrl: nonEmptyString(error.referenceUrl),
  };
};

const normalizeLegacyMeshkitError = (error: Record<string, unknown>): NormalizedEventError => ({
  code: nonEmptyString(error.Code),
  shortDescription: lines(error.ShortDescription).join('. ') || undefined,
  longDescription: lines(error.LongDescription),
  causes: [],
  probableCause: lines(error.ProbableCause),
  suggestedRemediation: lines(error.SuggestedRemediation),
});

/**
 * Reads an event's metadata.error, in any shape it is stored in, into one
 * view model. Returns null when the value is not an error object.
 */
export const normalizeEventError = (value: unknown): NormalizedEventError | null => {
  if (!isRecord(value)) return null;
  // `code` marks the schemas shapes; rows recorded before EventError carry
  // MeshKit's Go field names instead.
  if (typeof value.code === 'string') return normalizeSchemaEventError(value);
  return normalizeLegacyMeshkitError(value);
};

/** The error code an event's metadata.error names, in any shape. */
export const eventErrorCode = (value: unknown): string | undefined => {
  if (!isRecord(value)) return undefined;
  return nonEmptyString(value.code) ?? nonEmptyString(value.Code);
};

/** The error's own reference page, when the producer linked one. */
export const eventErrorReferenceUrl = (value: unknown): string | undefined =>
  isRecord(value) && typeof value.code === 'string'
    ? nonEmptyString(value.referenceUrl)
    : undefined;
