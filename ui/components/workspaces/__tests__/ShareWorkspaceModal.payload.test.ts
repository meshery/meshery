/**
 * Wire-contract regression test for the resource share payload.
 *
 * `ShareWorkspaceModal` hands Sistent's `ShareModal` the
 * `useCreateAndRevokeResourceAccessRecordMutation` mutator, and that mutation
 * forwards Sistent's `ResourceAccessMappingPayload` verbatim as the POST body
 * to `extensions/api/resource/{type}/share/{id}` (see `rtk-query/resource`).
 * The Cloud API reads camelCase keys and silently drops anything else, so a
 * snake_case payload grants nothing while still returning success - which is
 * exactly how sharing regressed between Sistent 0.20.x and 0.22.0
 * (layer5io/sistent#1786).
 *
 * Meshery owns this boundary even though the payload is built upstream, so
 * assert the emitted keys rather than trusting the dependency's version.
 */
import { buildGrantAccessPayload, buildRevokeAccessPayload } from '@sistent/sistent';
import { describe, expect, it } from 'vitest';

const user = (id: string) => ({ id, email: `${id}@example.com` });

const collectKeys = (value: unknown): string[] => {
  if (Array.isArray(value)) {
    return value.flatMap(collectKeys);
  }
  if (value !== null && typeof value === 'object') {
    return Object.entries(value as Record<string, unknown>).flatMap(([key, nested]) => [
      key,
      ...collectKeys(nested),
    ]);
  }
  return [];
};

describe('resource share payload wire contract', () => {
  it('emits camelCase top-level keys for a grant', () => {
    const payload = buildGrantAccessPayload([user('alice')]);

    expect(Object.keys(payload).sort()).toEqual(['grantAccess', 'notifyUsers', 'revokeAccess']);
  });

  it('emits camelCase top-level keys for a revoke', () => {
    const payload = buildRevokeAccessPayload([user('bob')]);

    expect(Object.keys(payload).sort()).toEqual(['grantAccess', 'notifyUsers', 'revokeAccess']);
  });

  it('places granted users under grantAccess and revoked users under revokeAccess', () => {
    expect(buildGrantAccessPayload([user('alice')]).grantAccess).toHaveLength(1);
    expect(buildGrantAccessPayload([user('alice')]).revokeAccess).toHaveLength(0);

    expect(buildRevokeAccessPayload([user('bob')]).revokeAccess).toHaveLength(1);
    expect(buildRevokeAccessPayload([user('bob')]).grantAccess).toHaveLength(0);
  });

  it('emits camelCase actor keys for a granted user', () => {
    const payload = buildGrantAccessPayload([user('alice')]);

    expect(payload.grantAccess).toHaveLength(1);
    expect(Object.keys(payload.grantAccess[0]).sort()).toEqual(['actorId', 'actorType']);
  });

  it('emits camelCase actor keys for a revoked user', () => {
    const payload = buildRevokeAccessPayload([user('bob')]);

    expect(payload.revokeAccess).toHaveLength(1);
    expect(Object.keys(payload.revokeAccess[0]).sort()).toEqual(['actorId', 'actorType']);
  });

  it('never emits a snake_case key at any nesting depth', () => {
    const payloads = [
      buildGrantAccessPayload([user('alice')]),
      buildRevokeAccessPayload([user('bob')]),
    ];

    for (const payload of payloads) {
      const keys = collectKeys(payload);

      expect(keys.length).toBeGreaterThan(Object.keys(payload).length);
      for (const key of keys) {
        expect(key).not.toMatch(/_/);
      }
    }
  });
});
