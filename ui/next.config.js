/** @type {import('next').NextConfig} */

const isPlaygroundBuild = process.env.PLAYGROUND === 'true';

const nextConfig = {
  reactStrictMode: true,

  env: {
    NEXT_PUBLIC_PLAYGROUND_BUILD: isPlaygroundBuild ? 'true' : 'false',
  },

  // Static export (replaces removed `next export` in Next.js 15)
  output: 'export',

  // Transpile packages that have CSS imports or ESM issues
  transpilePackages: [
    '@uiw/react-md-editor',
    '@uiw/react-markdown-preview',
    '@uiw/react-codemirror',
    'billboard.js',
  ],

  // SWC Compiler Configuration (SWC is the default in Next.js 15)
  compiler: {
    relay: require('./relay.config'),
  },

  // Turbopack configuration (default bundler in Next.js 16+).
  // An explicit config block is required when a `webpack` config is also present
  // so that Next.js 16 does not abort the build with a configuration-mismatch error.
  // resolveAlias mirrors the webpack alias for remote-component.config.js so that
  // @paciolan/remote-component can find the config in dev mode (Turbopack).
  turbopack: {
    resolveAlias: {
      'remote-component.config.js': './remote-component.config.js',
    },
  },

  // Required for static export
  images: {
    unoptimized: true,
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

    // `import.meta.webpackHot.accept()` into CommonJS MUI files under @mui/material.
    // Treat all JS files in this package as generic JS modules so `import.meta`
    // is accepted and dev server can compile.
    config.module.rules.push({
      test: /node_modules[\\/]@mui[\\/]material[\\/].*\.js$/,
      type: 'javascript/auto',
      resolve: {
        fullySpecified: false,
      },
    });

    // Similarly, billboard.js is bundled as CommonJS and is being transformed
    // by React Fast Refresh. Treat its dist bundle as generic JS as well so
    // injected `import.meta.webpackHot.accept()` is accepted.
    config.module.rules.push({
      test: /node_modules[\\/]billboard\.js[\\/]dist[\\/].*\.js$/,
      type: 'javascript/auto',
      resolve: {
        fullySpecified: false,
      },
    });

    return config;
  },
};

module.exports = nextConfig;
