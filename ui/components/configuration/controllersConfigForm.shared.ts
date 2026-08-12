export const INHERIT = '__inherit__';
export const WATCH_EVENTS = ['ADDED', 'MODIFIED', 'DELETED'] as const;
export const WATCH_MODE_OPTIONS = [
  { value: INHERIT, label: 'Inherit' },
  { value: 'whitelist', label: 'Whitelist (watch only these)' },
  { value: 'blacklist', label: 'Blacklist (default scope minus these)' },
] as const;

export type WatchList = {
  whitelist?: { resource: string; events?: string[] }[];
  blacklist?: string[];
};

// Longest string the control shows (value / placeholder / option) + field chrome.
export const fitWidth = (...parts: Array<string | number | undefined | null>) =>
  `calc(${Math.max(1, ...parts.map((part) => String(part ?? '').length))}ch + 2.75em)`;
