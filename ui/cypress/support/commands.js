import 'cypress-file-upload';

Cypress.Commands.add('selectProviderNone', () => {
  cy.request({
    method: 'GET',
    url: 'api/provider?provider=None'
  });
});