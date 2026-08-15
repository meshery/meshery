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

const chars = (...parts: Array<string | number | undefined | null>) =>
  Math.max(1, ...parts.map((part) => longestLine(String(part ?? ''))));

// Definite width (no %): a percentage inside a wrap row makes every field
// share the line equally and leaves a hole between the compact inputs.
export const fitWidth = (...parts: Array<string | number | undefined | null>) =>
  `calc(${chars(...parts)}ch + 2.75em)`;

// Number inputs keep stepper buttons; those eat the default chrome and wrap "Inherit".
export const fitNumberWidth = (...parts: Array<string | number | undefined | null>) =>
  `calc(${chars(...parts)}ch + 4.75em)`;
