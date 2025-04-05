const { PHASE_DEVELOPMENT_SERVER } = require('next/constants');

module.exports = (phase) => {
  const isDev = phase === PHASE_DEVELOPMENT_SERVER;

  return {
    reactStrictMode: true,
    assetPrefix: isDev ? '/' : '/provider',
    eslint: {
      ignoreDuringBuilds: true,
    },
    output: 'export',
  };
};
