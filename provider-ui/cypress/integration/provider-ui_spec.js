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

  it('selects a provider', () => {
    cy
      .get('[aria-label="Select Provider"]').click()
      .get('#split-button-menu').contains('None').click()
      .get('[aria-label="Select Provider"]').click()
  })
})
