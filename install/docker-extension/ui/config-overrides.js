const webpack = require('webpack');
module.exports = function override(config, env) {
    config.resolve.fallback = {
        url: require.resolve('url'),
        fs: require.resolve('fs'),
        assert: require.resolve('assert'),
        path: require.resolve('path-browserify'),
        http: require.resolve('stream-http'),
        https: require.resolve('https-browserify'),
        os: require.resolve('os-browserify/browser'),
        buffer: require.resolve('buffer'),
        stream: require.resolve('stream-browserify'),
    };
    // Webpack 5 resolves requests from ESM-flagged origins as "fully specified",
    // so the extensionless `process/browser` request emitted by the ProvidePlugin
    // below (and by dependencies such as @react-dnd/invariant) fails to resolve.
    // Pin it to the concrete file.
    config.resolve.alias = {
        ...config.resolve.alias,
        'process/browser': require.resolve('process/browser.js'),
    };
    config.plugins.push(
        new webpack.ProvidePlugin({
            process: 'process/browser',
            Buffer: ['buffer', 'Buffer'],
        }),
    );

    // @docker/docker-mui-theme ships ESM syntax but declares "type": "commonjs"
    // in its package.json. This causes webpack's CJS parser to reject the
    // import/export statements. Override the module type to auto-detect.
    config.module.rules.unshift({
        test: /\.js$/,
        include: /[\\/]node_modules[\\/]@docker[\\/]docker-mui-theme[\\/]/,
        type: 'javascript/auto',
    });

    return config;
}