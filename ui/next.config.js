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
  { webpack : (config)  => {
    config.resolve.alias = { ...config.resolve.alias,
      "remote-component.config.js" : __dirname + "/remote-component.config.js" };
    return config;
  } }
]);
