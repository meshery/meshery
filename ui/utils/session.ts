// Session-scoped persistence for the active organization and the CASL
// permission keys. Both are written from authenticated server responses, so a
// corrupted entry (interrupted write, tampering, a browser extension) must
// never take the app bootstrap down: parse defensively and clear the bad entry
// so the next read starts clean.
interface StoredOrganization {
  id: string;
  [key: string]: unknown;
}

interface StoredPermissionKey {
  id: string;
  function?: string;
}

export function loadStoredOrganization(): StoredOrganization | null {
  const raw = sessionStorage.getItem('currentOrg');
  if (!raw || raw === 'undefined' || raw === 'null') {
    if (raw === 'undefined' || raw === 'null') {
      sessionStorage.removeItem('currentOrg');
    }
    return null;
  }

  try {
    const parsed = JSON.parse(raw);
    if (
      parsed &&
      typeof parsed === 'object' &&
      typeof parsed.id === 'string' &&
      parsed.id.trim() !== ''
    ) {
      return parsed;
    }
  } catch {
    // Fall through to clearing the corrupted value.
  }

  sessionStorage.removeItem('currentOrg');
  return null;
}

export function loadStoredKeys(): StoredPermissionKey[] | null {
  const raw = sessionStorage.getItem('keys');
  if (!raw || raw === 'undefined' || raw === 'null') {
    if (raw === 'undefined' || raw === 'null') {
      sessionStorage.removeItem('keys');
    }
    return null;
  }

  try {
    const parsed = JSON.parse(raw);
    if (
      Array.isArray(parsed) &&
      parsed.every((key) => key && typeof key === 'object' && typeof key.id === 'string')
    ) {
      return parsed;
    }
  } catch {
    // Fall through to clearing the corrupted value.
  }

  sessionStorage.removeItem('keys');
  return null;
}
