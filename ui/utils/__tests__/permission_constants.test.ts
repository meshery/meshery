import { describe, expect, it } from 'vitest';
import { keys } from '../permission_constants';
import CAN, { ability } from '../can';

describe('permission_constants dynamic bridge and CASL capability tests', () => {
  it('should resolve legacy keys correctly with action and subject', () => {
    // legacy key AccountManagementViewProfile
    expect(keys.VIEW_PROFILE).toBeDefined();
    expect(keys.VIEW_PROFILE.action).toBe('fa7de118-2d08-4b07-b9d7-3e0baead6d04');
    expect(keys.VIEW_PROFILE.subject).toBe('View Profile');

    // legacy key AccountManagementEditAccount
    expect(keys.EDIT_ACCOUNT).toBeDefined();
    expect(keys.EDIT_ACCOUNT.action).toBe('f5e6bb39-c89a-4172-86f9-14a4a59792c1');
    expect(keys.EDIT_ACCOUNT.subject).toBe('Edit Account');
  });

  it('should resolve new PascalCase keys directly mapped from schemas', () => {
    expect(keys.AccountManagementViewProfile).toBeDefined();
    expect(keys.AccountManagementViewProfile.action).toBe('fa7de118-2d08-4b07-b9d7-3e0baead6d04');
    expect(keys.AccountManagementViewProfile.subject).toBe('View Profile');
  });

  it('every exported permission key has a non-empty subject and UUID action', () => {
    const uuidPattern = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    for (const [name, value] of Object.entries(keys)) {
      expect(value, `key=${name}`).toMatchObject({
        subject: expect.any(String),
        action: expect.any(String),
      });
      expect((value as { subject: string }).subject.length, `subject for ${name}`).toBeGreaterThan(
        0,
      );
      expect((value as { action: string }).action, `action for ${name}`).toMatch(uuidPattern);
    }
  });

  it('should verify CASL CAN capability with active abilities', () => {
    // Initialize ability with specific permissions
    ability.update([
      {
        action: 'fa7de118-2d08-4b07-b9d7-3e0baead6d04',
        subject: 'view profile',
      },
    ]);

    // Check CAN check with legacy key
    expect(CAN(keys.VIEW_PROFILE.action, keys.VIEW_PROFILE.subject)).toBe(true);

    // Check CAN check with PascalCase key
    expect(
      CAN(keys.AccountManagementViewProfile.action, keys.AccountManagementViewProfile.subject),
    ).toBe(true);

    // Check CAN check with a key the user doesn't have
    expect(CAN(keys.EDIT_ACCOUNT.action, keys.EDIT_ACCOUNT.subject)).toBe(false);
  });
});
