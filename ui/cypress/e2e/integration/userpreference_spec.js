describe('User Preferences', () => {
  describe('Analytics and Improvement Program', () => {
    beforeEach(() => {
      cy.intercept('GET', '/api/user/prefs', { fixture: 'stats.json' }).as('getUserStats');

      cy.visit('/user/preferences');
      cy.get(':nth-child(2) > .MuiFormControl-root > .MuiFormLabel-root').should('have.text', 'Analytics and Improvement Program');
      cy.wait('@getUserStats');
    });

    it('deactivates "Send Anonymous Usage Statistics"', () => {
      cy.intercept('POST', '/api/user/prefs', { fixture: 'stats.json' }).as('postUserStats');

      cy.get('[data-cy="UsageStatsPreference"]').click();
      cy.wait('@postUserStats');
      cy.get('[data-cy="UsageStatsPreference"]').should('not.have.class', 'Mui-checked');
    });

    it('activates "Send Anonymous Performance Results"', () => {
      cy.intercept('POST', '/api/user/prefs', { fixture: 'stats.json' }).as('postUserStats');

      cy.get('[data-cy="PerfResultPreference"]').click();
      cy.wait('@postUserStats');
      cy.get('[data-cy="PerfResultPreference"]').should('not.have.class', 'Mui-checked');
    });
  });
});
