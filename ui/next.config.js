/** @type {import('next').NextConfig} */

const nextConfig = {
  reactStrictMode: false,

  // Static export (replaces removed `next export` in Next.js 15)
  output: 'export',
  trailingSlash: true,

  eslint: {
    ignoreDuringBuilds: true,
  },

  typescript: {
    ignoreBuildErrors: true,
  },

  // Transpile packages that have CSS imports or ESM issues
  transpilePackages: [
    '@uiw/react-md-editor',
    '@uiw/react-markdown-preview',
    '@uiw/react-codemirror',
    'billboard.js',
  ],

  compiler: {
    relay: require('./relay.config'),
  },

  // Required for static export
  images: {
    unoptimized: true,
  },

  exportPathMap: function () {
    return {
      '/404': { page: '/404' },
      '/configuration/filters': { page: '/configuration/filters' },
      '/configuration/catalog': { page: '/configuration/catalog' },
      '/configuration/designs': { page: '/configuration/designs' },
      '/configuration/designs/configurator': { page: '/configuration/designs/configurator' },
      // /extension/[...component] handled via getStaticPaths
      '/extensions': { page: '/extensions' },
      '/': { page: '/' },
      '/management/adapter': { page: '/management/adapter' },
      '/management/environments': { page: '/management/environments' },
      '/management/connections': { page: '/management/connections' },
      '/management/workspaces': { page: '/management/workspaces' },
      '/performance': { page: '/performance' },
      '/performance/profiles': { page: '/performance/profiles' },
      '/settings': { page: '/settings' },
      '/user/preferences': { page: '/user/preferences' },
    };
  },

  webpack: (config) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      'remote-component.config.js': __dirname + '/remote-component.config.js',
    };

    // Handle markdown imports (replaces next-remove-imports)
    config.module.rules.push({
      test: /\.md$/,
      type: 'asset/source',
    });

    return config;
  },
};

module.exports = nextConfig;
