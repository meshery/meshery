describe('Visit Meshery', () => {
  it('Visits meshery', () => {
    cy.visit('http://localhost:3000')
  })
})

describe('Visit Meshery settings', () => {
  it('Visits meshery settings page', () => {
    cy.visit('http://localhost:3000/settings')
  })
})