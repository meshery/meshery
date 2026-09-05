export const isFieldEmpty = (value: unknown): boolean => {
  if (value === null || value === undefined) {
    return true;
  }
  if (typeof value === 'string') {
    return value.trim() === '';
  }
  if (Array.isArray(value)) {
    return value.length === 0;
  }
  return false;
};

export const isArrayEmpty = <T>(arr: ReadonlyArray<T> | null | undefined): boolean =>
  !arr || arr.length === 0;

export const isValidJSON = (str: unknown): boolean => {
  if (typeof str !== 'string' || str.trim() === '') {
    return false;
  }
  try {
    JSON.parse(str);
    return true;
  } catch {
    return false;
  }
};

export const normalizeSearchTerm = (term: string): string => term.toLowerCase().trim();

export const matchesSearch = (text: string, searchTerm: string): boolean => {
  const needle = normalizeSearchTerm(searchTerm);
  if (needle === '') {
    return true;
  }
  return text.toLowerCase().includes(needle);
};

/**
 * Validates a duration string consisting of a positive integer followed by a unit suffix.
 *
 * @param duration - The duration string to validate (e.g. '30s', '5m', '1h').
 * @param validUnits - Array of allowed unit characters (default: ['h', 'm', 's']).
 * @returns true if the duration is a non-empty string with a positive integer and valid unit suffix.
 */
export const isValidDuration = (
  duration: unknown,
  validUnits: string[] = ['h', 'm', 's'],
): boolean => {
  if (typeof duration !== 'string') {
    return false;
  }
  const trimmed = duration.trim();
  if (trimmed.length < 2) {
    return false;
  }
  const unit = trimmed.slice(-1).toLowerCase();
  const normalizedValidUnits = validUnits.map((u) => u.toLowerCase());
  if (!normalizedValidUnits.includes(unit)) {
    return false;
  }
  const numStr = trimmed.slice(0, -1);
  if (!/^\d+$/.test(numStr)) {
    return false;
  }
  const tNum = parseInt(numStr, 10);
  return !isNaN(tNum) && tNum > 0;
};
