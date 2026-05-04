/** @type {import('next').NextConfig} */
const isPlaygroundBuild = process.env.PLAYGROUND === 'true';

const nextConfig = {
  reactStrictMode: true,

  env: {
    NEXT_PUBLIC_PLAYGROUND_BUILD: isPlaygroundBuild ? 'true' : 'false',
  },

  // Static export (replaces removed `next export` in Next.js 15)
  output: 'export',

  typescript: {
    ignoreBuildErrors: false,
  },

  // Transpile packages that have CSS imports or ship ESM-only (needed for
  // SSR during static export because remote-component.config.js require()s them).
  transpilePackages: [
    '@uiw/react-md-editor',
    '@uiw/react-markdown-preview',
    '@uiw/react-codemirror',
    'billboard.js',
    '@reduxjs/toolkit',
    'react-redux',
    'uuid',
    'xstate',
    '@xstate/react',
    'graphql-ws',
    'yjs',
  ],

  // SWC Compiler Configuration (SWC is the default in Next.js 15)
  compiler: {
    relay: require('./relay.config'),
  },

  // Turbopack configuration (default bundler in Next.js 16+).
  turbopack: {
    resolveAlias: {
      'remote-component.config.js': './remote-component.config.js',
    },
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
      '/extension/[...component]': { page: '/extension/[...component]' },
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
};

module.exports = nextConfig;
