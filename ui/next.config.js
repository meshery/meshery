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
        '/performance/profiles': { page: '/performance/profiles'},
        // '/about': { page: '/about' },
        '/404': { page: '/_error' },
      }
    }
  }]);