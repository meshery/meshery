Cypress.Commands.add('selectProviderNone', () => {
  cy.request({
    method: 'GET',
    url: 'http://localhost:9081/api/provider?provider=None'
  });
});