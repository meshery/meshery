// ***********************************************************
// This example support/index.js is processed and
// loaded automatically before your test files.
//
// This is a great place to put global configuration and
// behavior that modifies Cypress.
//
// You can change the location of this file or turn off
// automatically serving support files with the
// 'supportFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

// Import commands.js using ES2015 syntax:
import './commands';
// import '@cypress/code-coverage/support';

// Alternatively you can use CommonJS syntax:
// require('./commands')

Cypress.on('uncaught:exception', (err, runnable) => {
  // While we don't expect to have uncaught exceptions nor unhandled promise rejections, we'd like our end-to-end tests to fail when a user-facing aspect
  // of Meshery UI / Meshery Server has been affected in order to identify root cause of defects besides the fact that there's some error being printed to
  // browser dev tool's console so we've decided to disable the default behavior of cypress to fail on such circumstances as part of issue #8640
  return false;
});
