describe('Test if UI components are displayed on Index Page', () => {
    it('Left Navigation Panel', () => {
      cy.visit('/');
      cy.get('[data-test=navigation]').should('be.visible');
    });

    it('Settings button', () => {
      cy.visit('/');
      cy.get('[data-test=settings-button]', { timeout: 10_000 }).should('be.visible');
    });

    it('Notification button', () => {
      cy.visit('/');
      cy.get('[data-test=notification-button]').should('be.visible');
    });

    it('Profile button', () => {
      cy.visit('/');
      cy.get('[data-test=profile-button]').should('be.visible');
    });

    it('Service Mesh Section', () => {
      cy.visit('/');
      cy.get('[data-test=service-mesh]').should('be.visible');
    });

});
