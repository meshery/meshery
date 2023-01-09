const compose = require('next-compose');

// module.exports = {
//   exportPathMap : function (pathMap) {
//     console.log(pathMap)
//     if (process.env.PLAYGROUND === "true") {
//       return {
//         '/' : { page : '/' },
//         '/extension/*' : { page : "/extension/[component]" }
//       }
//     } else {
//       return pathMap
//     }
//   },
//   webpack : (config)  => {
//     config.resolve.alias = { ...config.resolve.alias,
//       "remote-component.config.js" : __dirname + "/remote-component.config.js" };
//     return config;
//   }
// }

module.exports = compose([
  { exportPathMap : function () {
    return {
      '/404' : { page : '/404' },
      '/configuration/applications' : { page : '/configuration/applications' },
      '/configuration/filters' : { page : '/configuration/filters' },
      '/configuration/patterns' : { page : '/configuration/patterns' },
      '/extension/[...component]' : { page : '/extension/[...component]' },
      '/extensions' : { page : '/extensions' },
      '/' : { page : '/' },
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
  }, },
  { webpack : (config)  => {
    config.resolve.alias = { ...config.resolve.alias,
      "remote-component.config.js" : __dirname + "/remote-component.config.js" };
    return config;
  } }
]);