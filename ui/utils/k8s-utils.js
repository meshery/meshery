import _ from 'lodash';
import humanizeDuration from 'humanize-duration';
import React from 'react';
// import { matchPath, useHistory } from 'react-router';
// import helpers from '../helpers';
// import { filterGeneric, filterResource } from '../redux/filterSlice';
// import { useTypedSelector } from '../redux/reducers/reducers';
// import store from '../redux/stores/store';
// import { ApiError } from './k8s/apiProxy';
// import { KubeMetrics, KubeObjectInterface, Workload } from './k8s/cluster';
// import { KubeEvent } from './k8s/event';
// import Node from './k8s/node';
// import { parseCpu, parseRam, unparseCpu, unparseRam } from './units';

// Exported to keep compatibility for plugins that may have used them.
// export { filterGeneric, filterResource };

const humanize = humanizeDuration.humanizer();
humanize.languages['en-mini'] = {
  y: () => 'y',
  mo: () => 'mo',
  w: () => 'w',
  d: () => 'd',
  h: () => 'h',
  m: () => 'm',
  s: () => 's',
  ms: () => 'ms',
};

export const CLUSTER_ACTION_GRACE_PERIOD = 5000; // ms

// export type DateParam = string | number | Date;

// export type DateFormatOptions = 'brief' | 'mini';

// export interface TimeAgoOptions {
//   format?: DateFormatOptions;
// }

/**
 * Show the time passed since the given date, in the desired format.
 *
 * @param date - The date since which to calculate the duration.
 * @param options - `format` takes "brief" or "mini". "brief" rounds the date and uses the largest suitable unit (e.g. "4 weeks"). "mini" uses something like "4w" (for 4 weeks).
 * @returns The formatted date.
 */
export function timeAgo(date, options = {}) {
  const fromDate = new Date(date);
  let now = new Date();

  if (process.env.UNDER_TEST === 'true') {
    // For testing, we consider the current moment to be 3 months from the dates we are testing.
    const days = 24 * 3600 * 1000; // in ms
    now = new Date(fromDate.getTime() + 90 * days);
  }

  return formatDuration(now.getTime() - fromDate.getTime(), options);
}

/** Format a duration in milliseconds to a human-readable string.
 *
 * @param duration - The duration in milliseconds.
 * @param options - `format` takes "brief" or "mini". "brief" rounds the date and uses the largest suitable unit (e.g. "4 weeks"). "mini" uses something like "4w" (for 4 weeks).
 * @returns The formatted duration.
 * */
export function formatDuration(duration, options) {
  const { format = 'brief' } = options;

  if (format === 'brief') {
    return humanize(duration, {
      fallbacks: ['en'],
      round: true,
      largest: 1,
    });
  }

  return humanize(duration, {
    language: 'en-mini',
    spacer: '',
    fallbacks: ['en'],
    round: true,
    largest: 1,
  });
}

export function localeDate(date) {
  const options = { timeZoneName: 'short' };
  let locale = undefined;

  // Force the same conditions under test, so snapshots are the same.
  if (process.env.UNDER_TEST === 'true') {
    options.timeZone = 'UTC';
    options.hour12 = true;
    locale = 'en-US';
    return new Date(date).toISOString();
  } else {
    options.timeZone = store.getState().config.settings.timezone;
  }

  return new Date(date).toLocaleString(locale, options);
}

export function getPercentStr(value, total) {
  if (total === 0) {
    return null;
  }
  const percentage = (value / total) * 100;
  const decimals = percentage % 10 > 0 ? 1 : 0;
  return `${percentage.toFixed(decimals)} %`;
}

export function getReadyReplicas(item) {
  return item.status.readyReplicas || item.status.numberReady || 0;
}

export function getTotalReplicas(item) {
  return item.spec.replicas || item.status.currentNumberScheduled || 0;
}

// // /**
// //  * @returns A filter function that can be used to filter a list of items.
// //  *
// //  * @param matchCriteria - The JSONPath criteria to match.
// //  */
// // export function useFilterFunc<
// //   T extends { [key: string]: any } | KubeObjectInterface | KubeEvent =
// //     | KubeObjectInterface
// //     | KubeEvent
// // >(matchCriteria?: string[]) {
// //   const filter = useTypedSelector(state => state.filter);

// //   return (item: T) => {
// //     if (!!item.metadata) {
// //       return filterResource(item as KubeObjectInterface | KubeEvent, filter, matchCriteria);
// //     }
// //     return filterGeneric<T>(item, filter, matchCriteria);
// //   };
// // }

// export function getClusterPrefixedPath(path?: string | null) {
//   const baseClusterPath = '/c/:cluster';
//   if (!path) {
//     return baseClusterPath;
//   }
//   return baseClusterPath + (path[0] === '/' ? '' : '/') + path;
// }

// export function getCluster(): string | null {
//   const prefix = helpers.getBaseUrl();
//   const urlPath = helpers.isElectron()
//     ? window.location.hash.substr(1)
//     : window.location.pathname.slice(prefix.length);

//   const clusterURLMatch = matchPath<{ cluster?: string }>(urlPath, {
//     path: getClusterPrefixedPath(),
//   });
//   return (!!clusterURLMatch && clusterURLMatch.params.cluster) || null;
// }

// export function useErrorState(dependentSetter?: (...args: any) => void) {
//   const [error, setError] = React.useState<ApiError | null>(null);

//   React.useEffect(
//     () => {
//       if (!!error && !!dependentSetter) {
//         dependentSetter(null);
//       }
//     },
//     // eslint-disable-next-line
//     [error]
//   );

//   // Adding "as any" here because it was getting difficult to validate the setter type.
//   return [error, setError as any];
// }

// type URLStateParams<T> = {
//   /** The defaultValue for the URL state. */
//   defaultValue: T;
//   /** Whether to hide the parameter when the value is the default one (true by default). */
//   hideDefault?: boolean;
//   /** The prefix of the URL key to use for this state (a prefix 'my' with a key name 'key' will be used in the URL as 'my.key'). */
//   prefix?: string;
// };
// export function useURLState(
//   key: string,
//   defaultValue: number
// ): [number, React.Dispatch<React.SetStateAction<number>>];
// export function useURLState(
//   key: string,
//   valueOrParams: number | URLStateParams<number>
// ): [number, React.Dispatch<React.SetStateAction<number>>];
// /**
//  * A hook to manage a state variable that is also stored in the URL.
//  *
//  * @param key The name of the key in the URL. If empty, then the hook behaves like useState.
//  * @param paramsOrDefault The default value of the state variable, or the params object.
//  *
//  */
// export function useURLState<T extends string | number | undefined = string>(
//   key: string,
//   paramsOrDefault: T | URLStateParams<T>
// ): [T, React.Dispatch<React.SetStateAction<T>>] {
//   const params: URLStateParams<T> =
//     typeof paramsOrDefault === 'object' ? paramsOrDefault : { defaultValue: paramsOrDefault };
//   const { defaultValue, hideDefault = true, prefix = '' } = params;
//   const history = useHistory();
//   // Don't even use the prefix if the key is empty
//   const fullKey = !key ? '' : !!prefix ? prefix + '.' + key : key;

//   function getURLValue() {
//     // An empty key means that we don't want to use the state from the URL.
//     if (fullKey === '') {
//       return null;
//     }

//     const urlParams = new URLSearchParams(history.location.search);
//     const urlValue = urlParams.get(fullKey);
//     if (urlValue === null) {
//       return null;
//     }
//     let newValue: string | number = urlValue;
//     if (typeof defaultValue === 'number') {
//       newValue = Number(urlValue);
//       if (Number.isNaN(newValue)) {
//         return null;
//       }
//     }

//     return newValue;
//   }

//   const initialValue = React.useMemo(() => {
//     const newValue = getURLValue();
//     if (newValue === null) {
//       return defaultValue;
//     }
//     return newValue;
//   }, []);
//   const [value, setValue] = React.useState<T>(initialValue as T);

//   React.useEffect(
//     () => {
//       const newValue = getURLValue();
//       if (newValue === null) {
//         if (defaultValue !== undefined && defaultValue !== value) {
//           setValue(defaultValue);
//         }
//       } else if (newValue !== value) {
//         setValue(newValue as T);
//       }
//     },
//     // eslint-disable-next-line
//     [history]
//   );

//   React.useEffect(() => {
//     // An empty key means that we don't want to use the state from the URL.
//     if (fullKey === '') {
//       return;
//     }

//     const urlCurrentValue = getURLValue();

//     if (urlCurrentValue === value) {
//       return;
//     }

//     const urlParams = new URLSearchParams(history.location.search);
//     let shouldUpdateURL = false;

//     if ((value === null || value === defaultValue) && hideDefault) {
//       if (urlParams.has(fullKey)) {
//         urlParams.delete(fullKey);
//         shouldUpdateURL = true;
//       }
//     } else if (value !== undefined) {
//       const urlValue = value as NonNullable<T>;

//       urlParams.set(fullKey, urlValue.toString());
//       shouldUpdateURL = true;
//     }

//     if (shouldUpdateURL) {
//       history.replace({ search: urlParams.toString() });
//     }
//   }, [value]);

//   return [value, setValue] as [T, React.Dispatch<React.SetStateAction<T>>];
// }

// compareUnits compares two units and returns true if they are equal
export function compareUnits(quantity1, quantity2) {
  // strip whitespace and convert to lowercase
  const qty1 = quantity1.replace(/\s/g, '').toLowerCase();
  const qty2 = quantity2.replace(/\s/g, '').toLowerCase();

  // compare numbers
  return parseInt(qty1) === parseInt(qty2);
}

export function normalizeUnit(resourceType, quantity) {
  let type = resourceType;

  if (type.includes('.')) {
    type = type.split('.')[1];
  }

  let normalizedQuantity = '';
  let bytes = 0;
  switch (type) {
    case 'cpu':
      normalizedQuantity = quantity?.endsWith('m')
        ? `${Number(quantity.substring(0, quantity.length - 1)) / 1000}`
        : `${quantity}`;
      if (normalizedQuantity === '1') {
        normalizedQuantity = normalizedQuantity + ' ' + 'core';
      } else {
        normalizedQuantity = normalizedQuantity + ' ' + 'cores';
      }
      break;

    case 'memory':
      /**
       * Decimal: m | n | "" | k | M | G | T | P | E
       * Binary: Ki | Mi | Gi | Ti | Pi | Ei
       * Refer https://github.com/kubernetes-client/csharp/blob/840a90e24ef922adee0729e43859cf6b43567594/src/KubernetesClient.Models/ResourceQuantity.cs#L211
       */
      console.log('debug:', quantity, parseInt(quantity), quantity.endsWith('m'));
      bytes = parseInt(quantity);
      if (quantity.endsWith('Ki')) {
        bytes *= 1024;
      } else if (quantity.endsWith('Mi')) {
        bytes *= 1024 * 1024;
      } else if (quantity.endsWith('Gi')) {
        bytes *= 1024 * 1024 * 1024;
      } else if (quantity.endsWith('Ti')) {
        bytes *= 1024 * 1024 * 1024 * 1024;
      } else if (quantity.endsWith('Ei')) {
        bytes *= 1024 * 1024 * 1024 * 1024 * 1024;
      } else if (quantity.endsWith('m')) {
        bytes /= 1000;
      } else if (quantity.endsWith('u')) {
        bytes /= 1000 * 1000;
      } else if (quantity.endsWith('n')) {
        bytes /= 1000 * 1000 * 1000;
      } else if (quantity.endsWith('k')) {
        bytes *= 1000;
      } else if (quantity.endsWith('M')) {
        bytes *= 1000 * 1000;
      } else if (quantity.endsWith('G')) {
        bytes *= 1000 * 1000 * 1000;
      } else if (quantity.endsWith('T')) {
        bytes *= 1000 * 1000 * 1000 * 1000;
      } else if (quantity.endsWith('E')) {
        bytes *= 1000 * 1000 * 1000 * 1000 * 1000;
      }

      if (bytes === 0) {
        normalizedQuantity = '0 Bytes';
      } else {
        const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
        const k = 1000;
        const dm = 2;
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        console.debug('debug bytes:', bytes, i, sizes);
        normalizedQuantity = parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
      }
      break;

    default:
      normalizedQuantity = quantity;
      break;
  }
  return normalizedQuantity;
}

/** Creates a unique ID, with the given prefix.
 * If UNDER_TEST is set to true, it will return the same ID every time, so snapshots do not get invalidated.
 */
export function useId(prefix = '') {
  const [id] =
    (React.useState < string) |
    (undefined >
      (process.env.UNDER_TEST === 'true'
        ? prefix + 'id'
        : prefix + Math.random().toString(16).slice(2)));

  return id;
}

// // Make units available from here
// export * as auth from './auth';
// export * as units from './units';

const RAM_TYPES = ['Bi', 'Ki', 'Mi', 'Gi', 'Ti', 'Pi', 'Ei'];
const UNITS = ['B', 'K', 'M', 'G', 'T', 'P', 'E'];

export const TO_GB = 1024 * 1024 * 1024;
export const TO_ONE_M_CPU = 1000000;
export const TO_ONE_CPU = 1000000000;

export function parseDiskSpace(value) {
  return parseUnitsOfBytes(value);
}

export function parseRam(value) {
  return parseUnitsOfBytes(value);
}

function parseUnitsOfBytes(value) {
  if (!value) return 0;

  const groups = value.match(/(\d+)([BKMGTPEe])?(i)?(\d+)?/) || [];
  const number = parseInt(groups[1], 10);

  // number ex. 1000
  if (groups[2] === undefined) {
    return number;
  }

  // number with exponent ex. 1e3
  if (groups[4] !== undefined) {
    return number * 10 ** parseInt(groups[4], 10);
  }

  const unitIndex = _.indexOf(UNITS, groups[2]);

  // Unit + i ex. 1Ki
  if (groups[3] !== undefined) {
    return number * 1024 ** unitIndex;
  }

  // Unit ex. 1K
  return number * 1000 ** unitIndex;
}

export function unparseRam(value) {
  let i = 0;
  while (value >= 1024 && i < RAM_TYPES.length - 1) {
    i++;
    value /= 1024; // eslint-disable-line no-param-reassign
  }

  return {
    value: _.round(value, 1),
    unit: RAM_TYPES[i],
  };
}

export function parseCpu(value) {
  if (!value) return 0;

  const number = parseInt(value, 10);
  if (value.endsWith('n')) return number;
  if (value.endsWith('u')) return number * 1000;
  if (value.endsWith('m')) return number * 1000 * 1000;
  return number * 1000 * 1000 * 1000;
}

export function unparseCpu(value) {
  const result = parseFloat(value);

  return {
    value: _.round(result / 1000000, 2),
    unit: 'm',
  };
}

export const resourceParsers = {
  cpu: parseCpu,
  memory: parseRam,
};

export function getResourceStr(value, resourceType) {
  const resourceFormatters = {
    cpu: unparseCpu,
    memory: unparseRam,
  };

  const valueInfo = resourceFormatters[resourceType](value);
  return `${valueInfo.value}${valueInfo.unit}`;
}
