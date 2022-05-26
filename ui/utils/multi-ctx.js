/**
 * A function to be used by the requests sent for the
 * operations based on multi-context support
 *
 * @param {string} url The request URL
 * @param {Array.<string>} ctx The context Array
 * @returns {string} The final query-parametrised URL
 */
export function ctxUrl(url, ctx) {
  if (ctx?.length) {
    const contextQuery = ctx.map(context => `contexts=${context}`).join("&")
    return `${url}?${contextQuery}`
  }
  return url;
}
