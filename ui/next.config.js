const withCSS = require('@zeit/next-css');
const compose = require('next-compose');

module.exports = compose([
  [withCSS, {}],
  { exportPathMap : function () {
    return {
      '/' : { page : '/' },
      '/index' : { page : '/index' },
      '/performance/*' : { page : '/performance/[page]' },
      '/configure' : { page : '/configure' },
      '/playground' : { page : '/playground' },
      '/provider/*' : { page : '/provider' },
      '/settings/*' : { page : '/settings' },
      '/404' : { page : '/_error' },
    };
  }, },
  { webpack : (config, defaultLoaders)  => {
    config.module.rules.push({
      test : /\.css$/,
      use : [
        defaultLoaders.babel,
        {
          loader : require('styles-jsx/webpack').loader,
          options : {
            type : 'scoped'
          }
        }
      ]
    })
    config.resolve.alias = { ...config.resolve.alias,
      "remote-component.config.js" : __dirname + "/remote-component.config.js" };
    return config;
  } }
]);
