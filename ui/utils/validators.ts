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
