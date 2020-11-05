describe('Visit Meshery', () => {
  it('Visits meshery', () => {
    cy.visit('/')
  })
})

describe('Visit Meshery settings', () => {
  it('Visits meshery settings page', () => {
    cy.visit('/settings')
  })
})