describe('User Preferences', () => {
  describe('Extensions | Analytics and Improvement Program', () => {
    beforeEach(() => {
      cy.intercept('GET', '/api/user/prefs*').as('getUserStats');

      cy.selectProviderNone();

      cy.visit('/user/preferences');
      cy.get(':nth-child(1) > .MuiFormControl-root > .MuiFormLabel-root').should('have.text', 'Extensions');
      cy.get(':nth-child(2) > .MuiFormControl-root > .MuiFormLabel-root').should('have.text', 'Analytics and Improvement Program');
      cy.wait('@getUserStats');
    });

    context('Extensions', () => {
      it('activates Meshery Catalog Content', () => {
        cy.intercept('POST', '/api/user/prefs').as('postUserStatsToActivateExtensions');
        cy.get('[data-cy="CatalogContentPreference"]').click();
        cy.wait('@postUserStatsToActivateExtensions');
      });
      
      it('deactivates Meshery Catalog Content', () => {
        cy.intercept('POST', '/api/user/prefs').as('postUserStatsForExtensions');
        cy.get('[data-cy="CatalogContentPreference"]').click();
        cy.wait('@postUserStatsForExtensions');
      });
      
    })

    context('Analytics and Improvement Program', () => {
      it('deactivates "Send Anonymous Usage Statistics"', () => {
        cy.intercept('POST', '/api/user/prefs?contexts=all').as('postUserStats');
  
        cy.get('[data-cy="UsageStatsPreference"]').click();
        cy.wait('@postUserStats');
        cy.wait(2000)
        cy.get('[data-cy="UsageStatsPreference"]').should('not.have.class', 'Mui-checked');
      });
  
      it('activates "Send Anonymous Usage Statistics"', () => {
        cy.intercept('POST', '/api/user/prefs?contexts=all').as('postUserStats');
  
        cy.get('[data-cy="UsageStatsPreference"]').click();
        cy.wait('@postUserStats');
        cy.wait(2000)
        cy.get('[data-cy="UsageStatsPreference"]').should('have.class', 'Mui-checked');
      });
  
      it('deactivates "Send Anonymous Performance Results"', () => {
        cy.intercept('POST', '/api/user/prefs?contexts=all').as('postUserStats');
  
        cy.get('[data-cy="PerfResultPreference"]').click();
        cy.wait('@postUserStats');
        cy.wait(2000)
        cy.get('[data-cy="PerfResultPreference"]').should('not.have.class', 'Mui-checked');
      });
  
      it('activates "Send Anonymous Performance Results"', () => {
        cy.intercept('POST', '/api/user/prefs?contexts=all').as('postUserStats');

        cy.get('[data-cy="PerfResultPreference"]').click();
        cy.wait('@postUserStats');
        cy.wait(2000)
        cy.get('[data-cy="PerfResultPreference"]').should('have.class', 'Mui-checked');
      });
    })
  });
});
