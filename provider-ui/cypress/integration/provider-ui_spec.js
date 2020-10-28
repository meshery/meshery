describe('Provider UI', () => {
  it('renders provider component', () => {
    cy.server();

    cy.route('GET', '**/api/providers').as('getProviders');

    cy
      .visit('/')
      .url().should('include', 'http://localhost:3001/')
      .get('[data-cy=root]')
      .should('exist')
  })
})
