/** @type {import('next').NextConfig} */

// Bundle analyzer for build optimization monitoring
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});

const nextConfig = {
  // Note: Strict mode disabled temporarily due to initialization issues with some libraries
  // Re-enable once library compatibility is verified
  reactStrictMode: false,

  // Static export configuration (replaces `next export`)
  output: 'export',
  distDir: '.next',

  // Trailing slash for static hosting compatibility
  trailingSlash: true,

  eslint: {
    ignoreDuringBuilds: true,
  },

  // Transpile packages that have CSS imports or ESM issues
  transpilePackages: [
    '@uiw/react-md-editor',
    '@uiw/react-markdown-preview',
    '@uiw/react-codemirror',
    'billboard.js',
  ],

  typescript: {
    // Allow production builds to complete even with type errors
    ignoreBuildErrors: true,
  },

  // SWC Compiler Configuration (SWC is now the default in Next.js 15)
  compiler: {
    relay: require('./relay.config'),
    // Remove console.log in production (except errors/warnings)
    removeConsole: process.env.NODE_ENV === 'production' ? {
      exclude: ['error', 'warn'],
    } : false,
    // Enable React optimizations - removes data-testid in production
    reactRemoveProperties: process.env.NODE_ENV === 'production',
  },

  // Modularize imports for tree-shaking (now stable in Next.js 15)
  modularizeImports: {
    // Transform MUI icons for smaller bundles
    '@mui/icons-material': {
      transform: '@mui/icons-material/{{member}}',
    },
    // Note: lodash modularization disabled due to compatibility issues
    // with some lodash usage patterns in the codebase
  },

  // Experimental features
  experimental: {
    // Turbopack for faster dev builds (optional - uncomment to enable)
    // turbo: {},
  },

  // Images configuration for static export
  images: {
    unoptimized: true,
  },

  webpack: (config, { isServer, dev }) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      'remote-component.config.js': __dirname + '/remote-component.config.js',
    };

    // Handle markdown imports (replacement for next-remove-imports)
    config.module.rules.push({
      test: /\.md$/,
      type: 'asset/source',
    });

    // Production optimizations
    if (!dev && !isServer) {
      // Split chunks optimization for better caching
      config.optimization = {
        ...config.optimization,
        splitChunks: {
          ...config.optimization?.splitChunks,
          cacheGroups: {
            ...config.optimization?.splitChunks?.cacheGroups,
            // Separate vendor chunks for better caching
            vendor: {
              test: /[\\/]node_modules[\\/]/,
              name: 'vendors',
              chunks: 'all',
              priority: 10,
            },
            // Separate MUI into its own chunk
            mui: {
              test: /[\\/]node_modules[\\/]@mui[\\/]/,
              name: 'mui',
              chunks: 'all',
              priority: 20,
            },
            // Separate Sistent into its own chunk
            sistent: {
              test: /[\\/]node_modules[\\/]@sistent[\\/]/,
              name: 'sistent',
              chunks: 'all',
              priority: 20,
            },
            // Separate Redux/RTK into its own chunk
            redux: {
              test: /[\\/]node_modules[\\/](@reduxjs|react-redux|redux)[\\/]/,
              name: 'redux',
              chunks: 'all',
              priority: 20,
            },
          },
        },
      };
    }

    return config;
  },
};

module.exports = withBundleAnalyzer(nextConfig);
