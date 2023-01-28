const { defineConfig } = require('cypress')

module.exports = defineConfig({
  chromeWebSecurity : false,
  projectId : 'gpe84y',
  env : {
    codeCoverage : {
      url : '/api/__coverage__',
    },
  },
  e2e : {
    // We've imported your old cypress plugins here.
    // You may want to clean this up later by importing these.
    setupNodeEvents(on, config) {
      return require('./cypress/plugins/index.js')(on, config)
    },
    baseUrl : 'http://localhost:3001',
    specPattern : 'cypress/e2e/**/*.{js,jsx,ts,tsx}',
  },
})
