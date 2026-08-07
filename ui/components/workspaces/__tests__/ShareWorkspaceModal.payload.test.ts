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

  it('never emits a snake_case key', () => {
    const payloads = [
      buildGrantAccessPayload([user('alice')]),
      buildRevokeAccessPayload([user('bob')]),
    ];

    for (const payload of payloads) {
      for (const key of Object.keys(payload)) {
        expect(key).not.toMatch(/_/);
      }
    }
  });
});
