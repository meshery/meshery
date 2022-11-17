// eslint-disable-next-line no-undef no-unused-vars
const compose = require('next-compose');
const withCSS = require('@zeit/next-css');

module.exports = compose([
  [withCSS, {}],
  {
    images: {
      loader: 'imgix',
      path: 'https://example.com/myaccount/',
    },
  },
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
  } },
  
]);
