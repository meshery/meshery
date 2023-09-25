describe('Settings', () => {
  describe('Service Meshes', () => {
    beforeEach(() => {
      cy.visit('settings', { timeout: 45_000, headers: { 'Accept-Encoding': 'gzip, deflate' } });
      cy.intercept('GET', '/api/system/adapters', { fixture: 'getMeshAdapters.json' }).as(
        'getMeshAdapters',
      );
      cy.intercept('POST', '/api/system/adapter/manage', { fixture: 'postMeshManage.json' }).as(
        'postMeshManage',
      );
      cy.get('[data-cy="tabServiceMeshes"]').click({ force: true });
    });

    it('Adapter Connection Status', () => {
      cy.wait('@getMeshAdapters');
      cy.get("[data-cy='mesh-adapter-connections']").should('be.visible');
    });

    it('select, submit, and confirm', () => {
      cy.get("[data-cy='mesh-adapter-url']").type('localhost:10000{enter}');
      cy.get('[data-cy=btnSubmitMeshAdapter]').should('be.visible').click();
      cy.wait('@postMeshManage');
      cy.wait('@getMeshAdapters');
      cy.get('[data-cy=adapterSuccessSnackbar]').should('exist');
    });
  });
});
