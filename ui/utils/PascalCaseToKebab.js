/**
 * PascalCaseToKebab takes in a string in pascal case and returns
 * a kebab cased string
 * @param {string} str string in pascal case or camel case
 * @returns
 */
export default function PascalCaseToKebab(str) {
  return pascalCaseToCamelCase(str).replace(/[A-Z]/g, "-$&").toLowerCase();
}

/**
 * pascalCaseToCamelCase takes in a string in pascal case and
 * returns it in camelcase format
 * @param {string} str string that needs to be transformed
 */
function pascalCaseToCamelCase(str) {
  return str.charAt(0).toLowerCase() + str.slice(1);
}