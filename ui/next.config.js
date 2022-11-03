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
  webpack : (config)  => {
    config.resolve.alias = { ...config.resolve.alias,
      "remote-component.config.js" : __dirname + "/remote-component.config.js" };
    return config;
  }
}
