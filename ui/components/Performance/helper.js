/**
 * generateTestName takes in test name and service mesh name
 * and generates a random name (if test name is an empty string or is falsy) or
 * will return the given name
 *
 * @param {string} name
 * @param {string} meshName
 * @returns {string}
 */
export const generateTestName = (name, meshName) => {
  if (!name || name.trim() === '') {
    const mesh = meshName === '' || meshName === 'None' ? 'No mesh' : meshName;
    return `${mesh}_${new Date().getTime()}`;
  }

  return name;
};

export function generateUUID() {
  const { v4: uuid } = require('uuid');
  return uuid();
}
