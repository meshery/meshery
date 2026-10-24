import { describe, expect, it } from 'vitest';
import { normalizeUserProfileSummary } from '../userProfile';

describe('normalizeUserProfileSummary', () => {
  it('normalizes snake_case profile fields to camelCase', () => {
    expect(
      normalizeUserProfileSummary({
        id: 'user-id',
        email: 'owner@meshery.io',
        user_id: 'user-id',
        avatar_url: 'https://avatars.example.com/u/1',
        first_name: 'Mesh',
        last_name: 'Mate',
      }),
    ).toEqual({
      id: 'user-id',
      email: 'owner@meshery.io',
      userId: 'user-id',
      avatarUrl: 'https://avatars.example.com/u/1',
      firstName: 'Mesh',
      lastName: 'Mate',
    });
  });

  it('preserves camelCase profile fields', () => {
    expect(
      normalizeUserProfileSummary({
        id: 'user-id',
        email: 'owner@meshery.io',
        userId: 'user-id',
        avatarUrl: 'https://avatars.example.com/u/1',
        firstName: 'Mesh',
        lastName: 'Mate',
      }),
    ).toEqual({
      id: 'user-id',
      email: 'owner@meshery.io',
      userId: 'user-id',
      avatarUrl: 'https://avatars.example.com/u/1',
      firstName: 'Mesh',
      lastName: 'Mate',
    });
  });

  it('returns undefined for empty responses', () => {
    expect(normalizeUserProfileSummary()).toBeUndefined();
  });
});
