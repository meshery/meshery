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
        '/play': { page: '/play' },
        '/results': { page: '/results' },
        // '/about': { page: '/about' },
        '/404': { page: '/_error' },
      }
    }
  }]);