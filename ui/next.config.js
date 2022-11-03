// const withCSS = require('@zeit/next-css');
// const compose = require('next-compose');

module.exports = {
  exportPathMap : function (pathMap) {
    if (process.env.PLAYGROUND === "true") {
      return {
        '/' : { page : '/' },
        '/extension/*' : { page : "/extension/[component]" }
      }
    } else {
      return pathMap
    }
  },
}
