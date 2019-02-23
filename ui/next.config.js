module.exports = {
    exportPathMap: function () {
      return {
        '/': { page: '/' },
        '/index': { page: '/index' },
        '/load-test': { page: '/load-test' },
        '/about': { page: '/about' },
        // '/post': { page: '/post' },
        '/404': { page: '/_error' },
      }
    }
  }