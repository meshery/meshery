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

/**
 * The function takes in all the context and returns
 *  their respective cluster IDs associated to them
 *
 * @param {Array.<string>} selectedContexts
 * @param {Array.<string>} k8sconfig
 * @returns
 */
export const getK8sClusterIdsFromCtxId = (selectedContexts, k8sconfig) => {
  if (selectedContexts.length === 0){
    return []
  }

  return selectedContexts.map(c => k8sconfig.find(cfg => cfg.contextID === c)?.clusterID) || []
}
