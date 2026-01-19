/**
 * PascalCaseToKebab takes in a string in pascal case and returns
 * a kebab cased string
 * @param str - String in pascal case or camel case
 * @returns Kebab-cased string
 */
export default function PascalCaseToKebab(str: string): string {
  return pascalCaseToCamelCase(str).replace(/[A-Z]/g, '-$&').toLowerCase();
}

/**
 * pascalCaseToCamelCase takes in a string in pascal case and
 * returns it in camelcase format
 * @param str - String that needs to be transformed
 */
function pascalCaseToCamelCase(str: string): string {
  return str.charAt(0).toLowerCase() + str.slice(1);
}
