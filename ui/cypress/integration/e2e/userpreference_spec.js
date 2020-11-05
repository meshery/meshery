describe('User Preferences', () => {
  describe('Analytics and Improvement Program', () => {
    beforeEach(() => {
      cy.route2('GET', '/api/user/stats').as('getUserStats')

      cy.selectProviderNone()

      cy.visit('/userpreference')
      cy.get('.MuiFormLabel-root').should('have.text', 'Analytics and Improvement Program')
      cy.wait('@getUserStats')
    })

    it('deactivates "Send Anonymous Usage Statistics"', () => {
      cy.route2('POST', '/api/user/stats').as('postUserStats')

      cy.get('fieldset > div > label:nth-child(1)').click()
      cy.wait('@postUserStats')
      cy.get('.MuiSnackbarContent-message').should('have.text', 'Sending anonymous usage statistics was disabled')
    })

    it('activates "Send Anonymous Usage Statistics"', () => {
      cy.route2('POST', '/api/user/stats').as('postUserStats')

      cy.get('fieldset > div > label:nth-child(1)').click()
      cy.wait('@postUserStats')
      cy.get('.MuiSnackbarContent-message').should('have.text', 'Sending anonymous usage statistics was enabled')
    })

    it('deactivates "Send Anonymous Performance Results"', () => {
      cy.route2('POST', '/api/user/stats').as('postUserStats')

      cy.get('fieldset > div > label:nth-child(2)').click()
      cy.wait('@postUserStats')
      cy.get('.MuiSnackbarContent-message').should('have.text', 'Sending anonymous performance results was disabled')
    })

    it('activates "Send Anonymous Performance Results"', () => {
      cy.route2('POST', '/api/user/stats').as('postUserStats')

      cy.get('fieldset > div > label:nth-child(2)').click()
      cy.wait('@postUserStats')
      cy.get('.MuiSnackbarContent-message').should('have.text', 'Sending anonymous performance results was enabled')
    })
  })
})
