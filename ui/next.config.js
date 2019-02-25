module.exports = {
    exportPathMap: function () {
      return {
        '/': { page: '/k8s-config' },
        '/index': { page: '/k8s-config' },
        '/k8s-config': { page: '/k8s-config' },
        '/load-test': { page: '/load-test' },
        '/about': { page: '/about' },
        // '/post': { page: '/post' },
        '/404': { page: '/_error' },
      }
    }
  }