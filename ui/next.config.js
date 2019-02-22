module.exports = {
    exportPathMap: function () {
      return {
        '/': { page: '/' },
        '/index': { page: '/index' },
        '/about': { page: '/about' },
        // '/post': { page: '/post' },
        '/404': { page: '/_error' },
      }
    }
  }