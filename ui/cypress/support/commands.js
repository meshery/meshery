import 'cypress-file-upload';

Cypress.Commands.add('selectProviderNone', () => {
  cy.request({
    method: 'GET',
    url: 'http://localhost:3000/api/provider?provider=None'
  });
});
