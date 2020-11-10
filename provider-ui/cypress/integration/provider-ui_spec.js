describe('Provider UI', () => {
  beforeEach(() => {
    cy.route2('GET', '/api/providers', { fixture: 'providers.json' }).as('getProviders')
    cy.visit('/')
    cy.wait('@getProviders')
  })

  it('renders provider component', () => {
    cy
      .get('[data-cy=root]')
      .should('exist')
  })

  it('selects None as provider', () => {
    cy
      .get('[data-cy="select_provider"]').click()
      .get('#split-button-menu').contains('None').click()
      .get('[data-cy="select_provider"]').click()
  })
})
