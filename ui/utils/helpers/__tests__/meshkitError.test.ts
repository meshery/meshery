import { describe, expect, it } from 'vitest';
import { formatApiError, formatMeshkitErrorMarkdown, hasMeshkitError } from '../meshkitError';

// ---------------------------------------------------------------------------
// Unit tests for the MeshKit error formatter that backs `notifyApiError`.
//
// The schemas v1.2.2 baseQuery wrapper attaches a structured `meshkit`
// envelope on every non-2xx JSON response that matches the MeshKit shape.
// `formatApiError` turns that envelope into a markdown string the snackbar
// renders via `BasicMarkdown`. These tests guard:
//
//   1. Pure-string inputs and unknown shapes degrade safely.
//   2. A populated MeshKit envelope produces title + remediation list + code.
//   3. A minimal MeshKit envelope renders as just the bold title — no broken
//      "Try:" header or trailing whitespace.
// ---------------------------------------------------------------------------

describe('hasMeshkitError', () => {
  it('returns true for an RTK error carrying a meshkit object with a message', () => {
    expect(hasMeshkitError({ status: 500, data: {}, meshkit: { message: 'boom' } })).toBe(true);
  });

  it('returns false when meshkit is missing or malformed', () => {
    expect(hasMeshkitError({ status: 500, data: 'plain' })).toBe(false);
    expect(hasMeshkitError({ meshkit: {} })).toBe(false);
    expect(hasMeshkitError(null)).toBe(false);
    expect(hasMeshkitError(undefined)).toBe(false);
    expect(hasMeshkitError('string error')).toBe(false);
  });
});

describe('formatMeshkitErrorMarkdown', () => {
  it('renders the message as a bold title', () => {
    const md = formatMeshkitErrorMarkdown({ message: 'Workspace name already taken' });
    expect(md).toBe('**Workspace name already taken**');
  });

  it('includes a Try: section when remediations are present', () => {
    const md = formatMeshkitErrorMarkdown({
      message: 'Failed to publish design',
      suggestedRemediation: ['Re-run with a unique name', 'Check catalog permissions'],
    });
    expect(md).toContain('**Failed to publish design**');
    expect(md).toContain('*Try:*');
    expect(md).toContain('- Re-run with a unique name');
    expect(md).toContain('- Check catalog permissions');
  });

  it('emits a muted code reference when meshkit.code is set', () => {
    const md = formatMeshkitErrorMarkdown({
      message: 'Connection refused',
      code: 'meshery-server-1033',
    });
    expect(md).toContain('`meshery-server-1033`');
  });

  it('skips empty remediation entries', () => {
    const md = formatMeshkitErrorMarkdown({
      message: 'Bad request',
      suggestedRemediation: ['', '   ', 'Fix the body'],
    });
    expect(md).toContain('- Fix the body');
    expect(md).not.toMatch(/^- *$/m);
  });

  it('falls back to the supplied title when meshkit.message is empty', () => {
    const md = formatMeshkitErrorMarkdown({ message: '' }, 'Failed to create workspace');
    expect(md).toBe('**Failed to create workspace**');
  });
});

describe('formatApiError', () => {
  it('renders MeshKit metadata when the envelope is present', () => {
    const error = {
      status: 409,
      data: { error: 'duplicate' },
      meshkit: {
        message: 'Workspace "demo" already exists',
        code: 'meshery-server-1014',
        suggestedRemediation: ['Pick a different name'],
      },
    };
    const result = formatApiError(error, 'Failed to create workspace');
    expect(result.meshkit).toBeDefined();
    expect(result.message).toContain('**Workspace "demo" already exists**');
    expect(result.message).toContain('- Pick a different name');
    expect(result.message).toContain('`meshery-server-1014`');
  });

  it('falls back to FetchBaseQueryError.data string', () => {
    const result = formatApiError(
      { status: 500, data: 'internal server error' },
      'Failed to create workspace',
    );
    expect(result.meshkit).toBeUndefined();
    expect(result.message).toBe('internal server error');
  });

  it('falls back to nested data.error', () => {
    const result = formatApiError({ status: 400, data: { error: 'bad input' } });
    expect(result.message).toBe('bad input');
  });

  it('falls back to a SerializedError.message', () => {
    const result = formatApiError(new Error('network down'), 'Request failed');
    expect(result.message).toBe('network down');
  });

  it('uses fallbackTitle when no payload information is available', () => {
    const result = formatApiError({}, 'Failed to load workspaces');
    expect(result.message).toBe('Failed to load workspaces');
  });

  it('uses a generic fallback when nothing is supplied', () => {
    const result = formatApiError(undefined);
    expect(result.message).toBe('An unexpected error occurred');
  });
});
