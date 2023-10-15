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
