const { PHASE_DEVELOPMENT_SERVER } = require('next/constants')

module.exports = (phase) => {
  if (phase === PHASE_DEVELOPMENT_SERVER) {
    return {
      'assetPrefix' : '/',
      reactStrictMode : true,
    }
  }

  return {
    'assetPrefix' : '/provider',
    eslint : {
      // THis will not be needed since we are already using lint check in build
      ignoreDuringBuilds : true,
    },
  }
}
