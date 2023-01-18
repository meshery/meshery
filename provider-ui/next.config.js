module.exports = {
  reactStrictMode: true,
  // call basePath /provider
  basePath: '/provider',
  eslint: {
    // Warning: This allows production builds to successfully complete even if
    // your project has ESLint errors.
    ignoreDuringBuilds: true,
  },
  swcMinify: true,
  compiler: {
    removeConsole: true,
  },
};