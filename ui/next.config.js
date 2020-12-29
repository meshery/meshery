const withCSS = require('@zeit/next-css')
const compose = require('next-compose')

module.exports = compose([
  [withCSS, {}],
  {
    exportPathMap: function () {
      return {
        '/': { page: '/' },
        '/index': { page: '/index' },
        '/performance': { page: '/performance' },
        '/configure': { page: '/configure' },
        '/playground': { page: '/playground' },
        '/results': { page: '/results' },
        // '/about': { page: '/about' },
        '/provider/*': { page: '/provider' },
        '/404': { page: '/_error' },
      }
    },
  },
  {
    webpack: (config)  => {
      config.resolve.alias = {
        ...config.resolve.alias, 
        "remote-component.config.js": __dirname + "/remote-component.config.js"
      }
      return config
    }
  }
]);
