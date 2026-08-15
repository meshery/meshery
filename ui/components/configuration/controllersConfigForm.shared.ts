import type { UpdateControllersDefaultConfigApiArg } from '@meshery/schemas/mesheryApi';

export const INHERIT = '__inherit__';
export const WATCH_EVENTS = ['ADDED', 'MODIFIED', 'DELETED'] as const;
export const WATCH_MODE_OPTIONS = [
  { value: INHERIT, label: 'Inherit' },
  { value: 'whitelist', label: 'Whitelist (watch only these)' },
  { value: 'blacklist', label: 'Blacklist (default scope minus these)' },
] as const;

export type WatchList = NonNullable<
  NonNullable<UpdateControllersDefaultConfigApiArg['body']['meshsync']>['watchList']
>;

const longestLine = (value: string) => Math.max(1, ...value.split('\n').map((line) => line.length));

// Longest displayed line + field chrome, never wider than the container.
export const fitWidth = (...parts: Array<string | number | undefined | null>) =>
  `min(calc(${Math.max(1, ...parts.map((part) => longestLine(String(part ?? ''))))}ch + 2.75em), 100%)`;
