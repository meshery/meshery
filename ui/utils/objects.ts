import _ from 'lodash';

export function isEmptyAtAllDepths(input: unknown): boolean {
  if (_.isArray(input)) {
    // If the input is an array, check if all items are empty at all depths
    return input.every(isEmptyAtAllDepths);
  } else if (_.isObject(input)) {
    // If the input is an object, check if all properties are empty at all depths
    return _.every(input, isEmptyAtAllDepths);
  } else {
    // If the input is not an array or object, check if it's empty
    return _.isEmpty(input);
  }
}

/**
 * Finds the first nested object in a tree-like structure that satisfies a given condition.
 *
 * @param {Object} object - The root object to start the search.
 * @param {Function} condition - A function that takes an object as an argument and returns a boolean indicating if the condition is met.
 * @returns {Object|null} - The first object that satisfies the condition, or null if no matching object is found.
 */
export const findNestedObject = (
  object: unknown,
  condition: (_obj: any) => boolean,
): any | null => {
  const stack: any[] = [object];
  while (stack.length) {
    const currentObject = stack.pop();
    if (condition(currentObject)) {
      return currentObject;
    }
    if (_.isObject(currentObject) || _.isArray(currentObject)) {
      const anyObject = currentObject as any;
      const values = _.values(anyObject).filter((value) => value !== anyObject.models);
      stack.push(...values);
    }
  }
  return null;
};
/**
 * Accept object and removes empty properties from object.
 **/
export const filterEmptyFields = <T extends Record<string, any>>(data: T | null | undefined): T => {
  if (!data) {
    return {} as T;
  }

  return Object.keys(data).reduce((acc, key) => {
    const value = (data as any)[key];
    if (value !== undefined && value !== '') {
      (acc as any)[key] = value;
    }
    return acc;
  }, {} as T);
};
