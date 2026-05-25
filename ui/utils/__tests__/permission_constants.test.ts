import { describe, expect, it } from 'vitest';
import { keys } from '../permission_constants';

describe('permission_constants.keys', () => {
  it('exports a non-empty record of permission keys', () => {
    expect(typeof keys).toBe('object');
    expect(keys).not.toBeNull();
    expect(Object.keys(keys).length).toBeGreaterThan(0);
  });

  it('every entry has both a `subject` and a UUID-format `action`', () => {
    const uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

    for (const [name, value] of Object.entries(keys)) {
      expect(value, `key=${name}`).toMatchObject({
        subject: expect.any(String),
        action: expect.any(String),
      });
      expect(value.subject.length, `subject for ${name}`).toBeGreaterThan(0);
      expect(value.action, `action for ${name}`).toMatch(uuidPattern);
    }
  });

  it('exposes well-known permission keys with their expected subjects', () => {
    expect(keys.EDIT_ORGANIZATION.subject).toBe('Edit Organization');
    expect(keys.VIEW_ALL_ORGANIZATIONS.subject).toBe('View All Organizations');
    expect(keys.VIEW_PROFILE.subject).toBe('View Profile');
    expect(keys.RESET_DATABASE.subject).toBe('Reset database');
  });

  it('exposes design/filter/workspace lifecycle keys', () => {
    expect(keys.CREATE_NEW_DESIGN).toBeDefined();
    expect(keys.DELETE_A_DESIGN).toBeDefined();
    expect(keys.PUBLISH_DESIGN).toBeDefined();
    expect(keys.UNPUBLISH_DESIGN).toBeDefined();
    expect(keys.DELETE_WORKSPACE).toBeDefined();
    expect(keys.EDIT_WASM_FILTER).toBeDefined();
  });

  it('reuses the same action UUID for the two "download a design" subjects', () => {
    // DOWNLOAD_DESIGN and DOWNLOAD_A_DESIGN intentionally share an action ID
    // in the current key catalog. Pinning this so we notice if the seeds change.
    expect(keys.DOWNLOAD_DESIGN.action).toBe(keys.DOWNLOAD_A_DESIGN.action);
  });

  it('does not share action UUIDs between unrelated keys', () => {
    // Build a map of action -> [names] and assert that every collision is on
    // a known, intentional duplicate pair.
    const knownDuplicateGroups = new Set([
      ['DOWNLOAD_DESIGN', 'DOWNLOAD_A_DESIGN'].sort().join(','),
    ]);

    const byAction: Record<string, string[]> = {};
    for (const [name, value] of Object.entries(keys)) {
      byAction[value.action] = byAction[value.action] || [];
      byAction[value.action].push(name);
    }

    for (const [, names] of Object.entries(byAction)) {
      if (names.length > 1) {
        expect(knownDuplicateGroups).toContain(names.sort().join(','));
      }
    }
  });
});
