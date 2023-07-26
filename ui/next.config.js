/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode : false,
  compiler : {
    relay : require("./relay.config"),
  },
  exportPathMap : function () {
    return {
      '/404' : { page : '/404' },
      '/configuration/applications' : { page : '/configuration/applications' },
      '/configuration/filters' : { page : '/configuration/filters' },
      '/configuration/designs' : { page : '/configuration/designs' },
      '/configuration/designs/configurator' : { page : '/configuration/designs/configurator' },
      '/extension/[...component]' : { page : '/extension/[...component]' },
      '/extensions' : { page : '/extensions' },
      '/' : { page : '/', query : { __nextDefaultLocale : 'en' } },
      '/management/[mesh]' : { page : '/management/[mesh]' },
      '/management' : { page : '/management' },
      '/performance' : { page : '/performance' },
      '/performance/profiles' : { page : '/performance/profiles' },
      '/performance/results' : { page : '/performance/results' },
      '/settings' : { page : '/settings' },
      '/smi_results' : { page : '/smi_results' },
      '/system/connections' : { page : '/system/connections' },
      '/user/preferences' : { page : '/user/preferences' }
    };
  },
  //  exportPathMap: function (pathMap) {
  //    console.log(pathMap)
  //    if (process.env.PLAYGROUND === "true") {
  //      return {
  //        '/': { page: '/' },
  //        '/extension/*': { page: "/extension/[component]" }
  //      }
  //    } else {
  //      return pathMap
  //    }
  //  },
  webpack : (config) => {
    config.resolve.alias = {
      ...config.resolve.alias,
      "remote-component.config.js" : __dirname + "/remote-component.config.js"
    };
    config.output.webassemblyModuleFilename = 'static/[modulehash].wasm',
    config.experiments = { asyncWebAssembly : true, layers : true }
    return config;
  }
}

module.exports = nextConfig
