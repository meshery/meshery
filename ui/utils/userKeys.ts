import type { GetUserKeysApiResponse } from '@meshery/schemas/mesheryApi';

export type UserKey = GetUserKeysApiResponse['keys'][number];

const USER_KEYS_STORAGE_KEY = 'keys';

const isNonBlankString = (value: unknown): value is string =>
  typeof value === 'string' && value.trim() !== '';

// `function` is required too: _app.tsx maps it into the CASL subject.
const isUserKey = (value: unknown): value is UserKey =>
  typeof value === 'object' &&
  value !== null &&
  isNonBlankString((value as UserKey).id) &&
  isNonBlankString((value as UserKey).function);

/**
 * Reads the user's permission keys cached for this browser session.
 *
 * @returns The cached keys, or `null` when nothing valid is cached and the
 *   caller should fetch fresh keys. A corrupted entry is removed.
 */
export function loadCachedUserKeys(): UserKey[] | null {
  if (typeof window === 'undefined') return null;
  const raw = window.sessionStorage.getItem(USER_KEYS_STORAGE_KEY);
  if (raw === null) return null;
  try {
    const parsed = JSON.parse(raw);
    if (Array.isArray(parsed) && parsed.every(isUserKey)) return parsed;
  } catch {
    // Fall through: malformed JSON is treated the same as a wrong-type value.
  }
  window.sessionStorage.removeItem(USER_KEYS_STORAGE_KEY);
  return null;
}
