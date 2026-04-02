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