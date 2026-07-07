import { useCan } from '@/utils/hooks/useCan';

/**
 * A thin hook-based helper for computing a `disabled` prop from a permission
 * descriptor. This exists because `useCan()` is a React hook and cannot be
 * called conditionally or inside `.map()` callbacks.
 *
 * @param permission - Object with `action` and `subject`, or `undefined`.
 * @returns `true` when the permission exists and the user is NOT allowed.
 */
export function usePermissionDisabled(
  permission: { action: string; subject: string } | undefined,
): boolean {
  // Hooks must always be called, so we pass empty strings when there is no
  // permission and ignore the result.
  const allowed = useCan(permission?.action ?? '', permission?.subject ?? '');

  if (!permission) return false;
  return !allowed;
}
