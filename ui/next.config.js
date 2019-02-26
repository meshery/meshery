module.exports = {
    exportPathMap: function () {
      return {
        '/': { page: '/performance' },
        '/index': { page: '/performance' },
        '/performance': { page: '/performance' },
        '/configure': { page: '/configure' },
        '/play': { page: '/play' },
        // '/post': { page: '/post' },
        '/404': { page: '/_error' },
      }
    }
  }