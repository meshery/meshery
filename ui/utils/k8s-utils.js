import _ from 'lodash';
import humanizeDuration from 'humanize-duration';

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
