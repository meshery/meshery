import { describe, expect, it } from 'vitest';
import { eventErrorCode, eventErrorReferenceUrl, normalizeEventError } from './event_error';

// The acceptance sample of meshery/schemas#1334: two causes that share one
// remediation and each have one of their own, with action policies.
const sharedReconnect = {
  id: 'reconnect',
  description: 'Reconnect the connection once the cluster is reachable.',
  action: 'connection.reconnect',
  actionPolicy: 'requiresApproval',
};
const eventError = {
  code: 'meshery-server-1304',
  shortDescription: 'Error communicating with KubeAPI at https://prod:6443.',
  longDescription: 'dial tcp 10.0.0.1:6443: connect: connection refused',
  probableCauses: [
    {
      id: 'cluster-unreachable',
      description: "The cluster's API server is not reachable from Meshery Server.",
      remediations: [
        {
          id: 'check-network',
          description: "Check that Meshery Server can reach the cluster's API server address.",
          actionPolicy: 'manual',
        },
        sharedReconnect,
      ],
    },
    {
      id: 'credentials-expired',
      description: "The connection's credentials for the cluster have expired.",
      remediations: [
        {
          id: 'refresh-credentials',
          description: "Refresh the connection's credentials from its credential store.",
          action: 'connection.refresh-credentials',
          actionPolicy: 'automatic',
        },
        sharedReconnect,
      ],
    },
  ],
  matchedCauses: [{ causeId: 'credentials-expired', confidence: 0.8 }],
  referenceUrl: 'https://docs.meshery.io/reference/error-codes/component-meshery-server/x/',
};

describe('normalizeEventError', () => {
  it('leads with the matched cause and keeps each cause with its own remediations', () => {
    const error = normalizeEventError(eventError)!;
    expect(error.code).toBe('meshery-server-1304');
    expect(error.causes.map((cause) => [cause.id, cause.matched, cause.confidence])).toEqual([
      ['credentials-expired', true, 0.8],
      ['cluster-unreachable', false, undefined],
    ]);
    expect(error.causes[0].remediations.map((r) => [r.id, r.actionPolicy])).toEqual([
      ['refresh-credentials', 'automatic'],
      ['reconnect', 'requiresApproval'],
    ]);
    // The shared remediation is listed under every cause it answers.
    expect(error.causes[1].remediations.map((r) => r.id)).toEqual(['check-network', 'reconnect']);
    expect(error.referenceUrl).toBe(eventError.referenceUrl);
    expect(error.longDescription).toEqual([eventError.longDescription]);
  });

  it('derives the flat lists exactly as the contract does', () => {
    const error = normalizeEventError(eventError)!;
    expect(error.probableCause).toEqual([
      "The cluster's API server is not reachable from Meshery Server.",
      "The connection's credentials for the cluster have expired.",
    ]);
    expect(error.suggestedRemediation).toEqual([
      "Check that Meshery Server can reach the cluster's API server address.",
      'Reconnect the connection once the cluster is reachable.',
      "Refresh the connection's credentials from its credential store.",
    ]);
  });

  it('treats an unknown or omitted action policy as manual', () => {
    const error = normalizeEventError({
      code: 'meshery-server-1',
      probableCauses: [
        {
          id: 'a',
          description: 'A',
          remediations: [
            { id: 'r1', description: 'R1' },
            { id: 'r2', description: 'R2', action: 'x.y', actionPolicy: 'someday' },
          ],
        },
      ],
    })!;
    expect(error.causes[0].remediations.map((r) => r.actionPolicy)).toEqual(['manual', 'manual']);
  });

  it('ignores a matched cause the error does not declare', () => {
    const error = normalizeEventError({ ...eventError, matchedCauses: [{ causeId: 'unknown' }] })!;
    expect(error.causes.map((cause) => cause.matched)).toEqual([false, false]);
  });

  it('reads an EventError that carries only the flat lists as unpaired', () => {
    const error = normalizeEventError({
      code: 'meshery-server-1',
      probableCause: ['Cause A'],
      suggestedRemediation: ['Fix B', 'Fix C'],
    })!;
    expect(error.causes).toEqual([]);
    expect(error.probableCause).toEqual(['Cause A']);
    expect(error.suggestedRemediation).toEqual(['Fix B', 'Fix C']);
  });

  it('reads a legacy MeshKit error stored with Go field names', () => {
    const error = normalizeEventError({
      Code: 'meshery-server-1033',
      Severity: 2,
      ShortDescription: ['Unable to do X'],
      LongDescription: ['First line', '', 'Second line'],
      ProbableCause: ['Cause A', 'None'],
      SuggestedRemediation: ['Fix B'],
    })!;
    expect(error).toEqual({
      code: 'meshery-server-1033',
      shortDescription: 'Unable to do X',
      longDescription: ['First line', 'Second line'],
      causes: [],
      probableCause: ['Cause A'],
      suggestedRemediation: ['Fix B'],
    });
  });

  it('reads a design event operation error', () => {
    const error = normalizeEventError({ code: 'meshery-server-1480', message: 'deploy failed' })!;
    expect(error.code).toBe('meshery-server-1480');
    expect(error.longDescription).toEqual(['deploy failed']);
  });

  it('returns null for a value that is not an error object', () => {
    expect(normalizeEventError(undefined)).toBeNull();
    expect(normalizeEventError('boom')).toBeNull();
    expect(normalizeEventError([{ Code: 'A' }])).toBeNull();
  });
});

describe('eventErrorCode and eventErrorReferenceUrl', () => {
  it('reads the code of either shape and the reference URL of an EventError', () => {
    expect(eventErrorCode(eventError)).toBe('meshery-server-1304');
    expect(eventErrorCode({ Code: 'meshery-server-1033' })).toBe('meshery-server-1033');
    expect(eventErrorCode(null)).toBeUndefined();
    expect(eventErrorReferenceUrl(eventError)).toBe(eventError.referenceUrl);
    expect(
      eventErrorReferenceUrl({ Code: 'x', referenceUrl: 'https://example.com' }),
    ).toBeUndefined();
  });
});
