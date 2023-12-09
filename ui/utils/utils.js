import { trueRandom } from '../lib/trueRandom';
import jsYaml from 'js-yaml';
import { findWorkloadByName } from './workloadFilter';
import { EVENT_TYPES } from './Enum';
import _ from 'lodash';
import { getWebAdress } from './webApis';
import { APPLICATION, DESIGN, FILTER } from '../constants/navigator';
import { Tooltip } from '@mui/material';

/**
 * Check if an object is empty
 *
 * @param {Object} obj
 *
 * @returns {Boolean} if obj is empty
 */
export function isEmptyObj(obj) {
  return (
    !obj ||
    (obj && Object.keys(obj).length === 0 && Object.getPrototypeOf(obj) === Object.prototype)
  );
}

/**
 * Check if array is empty
 *
 * @param {Array} arr
 * @returns {Boolean} if arr is empty
 */
export function isEmptyArr(arr) {
  return arr && arr.length === 0;
}

/**
 * Check if two arrays are equal
 *
 * @param {Array} arr1
 * @param {Array} arr2
 * @param {Boolean} orderMatters
 * @returns
 */
export function isEqualArr(arr1, arr2, orderMatters = true) {
  if (arr1 === arr2) return true;
  if (arr1 == null || arr2 == null) return false;
  if (arr1.length !== arr2.length) return false;

  if (!orderMatters) {
    arr1.sort();
    arr2.sort();
  }

  for (var i = 0; i < arr1.length; ++i) {
    if (arr1[i] !== arr2[i]) return false;
  }
  return true;
}

/**
 * ScrollToTop scrolls the window to top
 *
 * @param {(
 * "auto"
 * |"smooth"
 * |"inherit"
 * |"initial"
 * |"revert"
 * |"unset"
 * )} behavior : scroll-behaviour, see https://developer.mozilla.org/en-US/docs/Web/CSS/scroll-behavior
 */
export function scrollToTop(behavior = 'smooth') {
  setTimeout(() => {
    window.scrollTo({
      top: 0,
      left: 0,
      behavior,
    });
  }, 0);
}

/**
 * Generates random Pattern Name with the prefix meshery_
 */
export function randomPatternNameGenerator() {
  return 'meshery_' + Math.floor(trueRandom() * 100);
}

/**
 * Returns number of components in Pattern/Application/Filters file
 */
export function getComponentsinFile(file) {
  if (file) {
    try {
      let keys = Object.keys(jsYaml.load(file).services);
      return keys.length;
    } catch (e) {
      if (e.reason?.includes('expected a single document')) {
        return file.split('---').length;
      }
    }
  }
  return 0;
}

export function generateValidatePayload(pattern_file, workloadTraitSet) {
  let pattern = jsYaml.loadAll(pattern_file);
  const services = pattern[0]?.services;
  if (!services) {
    return { err: 'Services not found in the design' };
  }

  const validationPayloads = {};

  for (const serviceId in services) {
    let valueType;

    let { workload } = findWorkloadByName(services[serviceId].type, workloadTraitSet);

    if (!(workload && workload?.oam_ref_schema)) {
      continue;
    }
    const schema = workload.oam_ref_schema;
    const value = services[serviceId]?.settings;
    if (!value) {
      continue;
    }
    valueType = 'JSON';
    const validationPayload = {
      schema,
      value: JSON.stringify(value),
      valueType,
    };
    validationPayloads[serviceId] = validationPayload;
  }

  return validationPayloads;
}

export function updateURLs(urlsSet, newUrls, eventType) {
  switch (eventType) {
    case EVENT_TYPES.DELETED:
      newUrls.forEach((url) => {
        urlsSet.delete(url);
      });
      break;
    case EVENT_TYPES.ADDED:
    case EVENT_TYPES.MODIFIED:
      newUrls.forEach((url) => {
        urlsSet.add(url);
      });
  }
}

/**
 * Gets the raw b64 file and convert it to Binary
 *
 * @param {string} file
 * @returns
 */
export function getDecodedFile(file) {
  // Extract base64-encoded content
  var encodedContent = file.split(';base64,')[1];

  // Decode base64 content
  return atob(encodedContent);
}

/**
 * Gets the raw b64 file and convert it to uint8Array
 *
 * @param {string} file
 * @returns {array} - return array of uint8Array
 */
export const getUnit8ArrayDecodedFile = (dataUrl) => {
  // Extract base64 content
  const [, base64Content] = dataUrl.split(';base64,');

  // Decode base64 content
  const decodedContent = atob(base64Content);

  // Convert decoded content to Uint8Array directly
  const uint8Array = Uint8Array.from(decodedContent, (char) => char.charCodeAt(0));

  return Array.from(uint8Array);
};

/**
 * Change the value of a property in RJSF schema
 *
 * @param {string} schema - RJSF schema
 * @param {string} propertyPath - path of the property to be modified
 * @param {any} newValue - new value to be set
 * @returns {object} - modified schema
 */
export const modifyRJSFSchema = (schema, propertyPath, newValue) => {
  const clonedSchema = _.cloneDeep(schema);
  _.set(clonedSchema, propertyPath, newValue);
  return clonedSchema;
};

/**
 * get sharable link with same and host and protocol, here until meshery cloud interception
 * @param {Object.<string,string>} sharedResource
 * @returns {string}
 */
export function getSharableCommonHostAndprotocolLink(sharedResource) {
  const webAddr = getWebAdress() + '/extension/meshmap';
  if (sharedResource?.application_file) {
    return `${webAddr}?${APPLICATION}=${sharedResource.id}`;
  }
  if (sharedResource?.pattern_file) {
    return `${webAddr}?${DESIGN}=${sharedResource.id}`;
  }
  if (sharedResource?.filter_resource) {
    return `${webAddr}?${FILTER}=${sharedResource.id}`;
  }

  return '';
}

/**
 * Retrieves the value of a specified column from a row of data.
 *
 * @param {Array} rowData - The array representing the row of data.
 * @param {string} columnName - The name of the column whose value you want to retrieve.
 * @param {Array} columns - An array of column definitions.
 * @returns {*} The value of the specified column in the row, or undefined if not found.
 */

export const getColumnValue = (rowData, columnName, columns) => {
  const columnIndex = columns.findIndex((column) => column.name === columnName);
  return rowData[columnIndex];
};

/**
 * Filter the columns to show in visibility switch.
 *
 * @param {string} columns - Full list of columns name.
 *
 */

export const getVisibilityColums = (columns) => {
  return columns.filter((col) => col?.options?.display !== false);
};

export function JsonParse(item, safe = true) {
  if (typeof item === 'string') {
    try {
      return JSON.parse(item || '{}');
    } catch (e) {
      if (safe) {
        return {};
      }
      throw e;
    }
  }

  return item;
}

export const ConditionalTooltip = ({ value, maxLength, ...restProps }) => {
  return value?.length > maxLength ? (
    <Tooltip title={value} arrow placement="top">
      <div
        style={{
          maxWidth: '15rem',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
        }}
        {...restProps}
      >
        {`${value.slice(0, maxLength)}...`}
      </div>
    </Tooltip>
  ) : (
    <div
      style={{
        maxWidth: '15rem',
        overflow: 'hidden',
        textOverflow: 'ellipsis',
        whiteSpace: 'nowrap',
      }}
      {...restProps}
    >
      {value}
    </div>
  );
};

/**
 * Handle scroll event for infinite scrolling
 *
 * @param {string} scrollingView - The view identifier for which scrolling is handled
 * @param {function} setPage - The function to set the page state
 * @param {object} scrollRef - Reference to the scroll element
 * @param {number} buffer - The buffer value for infinite scrolling
 * @returns {function} - Scroll event handler function
 */
export const createScrollHandler = (scrollingView, setPage, scrollRef, buffer) => (event) => {
  const div = event.target;
  if (div.scrollTop >= div.scrollHeight - div.clientHeight - buffer) {
    setPage((prevPage) => ({
      ...prevPage,
      [scrollingView]: prevPage[scrollingView] + 1,
    }));
  }

  scrollRef.current = div.scrollTop;
};

/**
 *
 * Add underscore to camal case variable name.
 *
 * @param {string} value - An array of column definitions.
 *
 */

export const camelcaseToSnakecase = (value) => {
  return value?.replace(/[A-Z]/g, (match) => `_${match.toLowerCase()}`);
};

export const formatToTitleCase = (value) => {
  if (typeof value === 'string') {
    return value.substring(0, 1).toUpperCase().concat('', value.substring(1).toLowerCase());
  }
  return '';
};
