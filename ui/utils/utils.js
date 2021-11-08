/**
 * captitalizestring capitalises the string presented to it
 * "hello world!" -> "Hello World!"
 * "hello_world!" -> "Hello World"
 *
 * @param {string} str The uncapitalised string
 *
 * @returns {string} the capitalised string
 */
export function capitalizeString(str) {
  return str.split(/[\s|_]/) // separating regex for "_" or " "
    .map(smallStr => smallStr.charAt(0).toUpperCase() + smallStr.substring(1).toLowerCase()) // transforming string
    .join(" ") // join by whitespace
}
