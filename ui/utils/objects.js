import _ from 'lodash';

export function isEmptyAtAllDepths(input) {
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
export const findNestedObject = (object, condition) => {
  const stack = [object];
  while (stack.length) {
    const currentObject = stack.pop();
    if (condition(currentObject)) {
      return currentObject;
    }
    if (_.isObject(currentObject)) {
      // Push each property of the currentObject onto the stack for further exploration
      stack.push(..._.values(currentObject));
    }
  }
  return null;
};
