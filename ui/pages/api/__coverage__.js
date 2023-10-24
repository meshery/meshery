// Needed for fetching server-side code coverage. See: https://github.com/bahmutov/next-and-cypress-example

if (process.env.NODE_ENV === 'development') {
  module.exports = require('@cypress/code-coverage/middleware/nextjs');
}
