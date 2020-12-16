describe('General Preferences', () => {
  describe('Analytics and Improvement Program', () => {
    beforeEach(() => {
      cy.route2('GET', '/api/user/stats', { fixture: 'stats.json' }).as('getUserStats')

      cy.visit('/userpreference')
      cy.get('.MuiFormLabel-root').should('have.text', 'Analytics and Improvement Program')
      cy.wait('@getUserStats')
    })

    it('deactivates "Send Anonymous Usage Statistics"', () => {
      cy.route2('POST', '/api/user/stats', { fixture: 'stats.json' }).as('postUserStats')

      cy.get('[data-cy="UsageStatsPreference"]').click()
      cy.wait('@postUserStats')
      cy.get('[data-cy="UsageStatsPreference"]').should('not.have.class', 'Mui-checked')
    })

    it('deactivates "Send Anonymous Performance Results"', () => {
      cy.route2('POST', '/api/user/stats', { fixture: 'stats.json' }).as('postUserStats')

      cy.get('[data-cy="PerfResultPreference"]').click()
      cy.wait('@postUserStats')
      cy.get('[data-cy="PerfResultPreference"]').should('not.have.class', 'Mui-checked')
    })
  })
})


describe('Remote Provider Preferences', () => {
  describe('MeshMap User Settings', () => {
    beforeEach(() => {
      cy.route2('POST', '/api/user/stats', { fixture: 'stats.json' }).as('postUserStats')

      cy.visit('/userpreference')
      cy.get('.MuiTab-wrapper').contains('Remote Provider').click()
      cy.wait('@getUserStats')
    })

    it('deactivates "Start MeshMap on Zoom"', () => {
      cy.route2('POST', '/api/user/stats', { fixture: 'stats.json' }).as('postUserStats')

      cy.get('[data-cy="StartOnZoomPreference"]').click()
      cy.wait('@postUserStats')
      cy.get('[data-cy="StartOnZoomPreference"]').should('not.have.class', 'Mui-checked')
    })
  })
})
