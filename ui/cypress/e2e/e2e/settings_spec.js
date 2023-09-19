describe('Settings', () => {
  describe('Service Meshes', () => {
    beforeEach(() => {
      cy.intercept('GET', '/api/system/adapters').as('getMeshAdapters');

      cy.selectProviderNone();
      cy.visit('/settings');

      cy.get('[data-cy="tabServiceMeshes"]').click();
      cy.wait('@getMeshAdapters');
    });

    it('ping and submit Consul', () => {
      cy.intercept('GET', '/api/system/adapters*').as('getAdapterPing');
      cy.intercept('POST', '/api/system/adapter/manage').as('postMeshManage');
      cy.intercept('GET', '/api/system/adapters').as('getMeshAdapters');

      cy.get('[data-cy=chipAdapterLocation]', { timeout: 30_000 })
        .contains('.MuiChip-label', 'localhost:10002')
        .click();

      cy.wait('@getAdapterPing');

      cy.get("[data-cy='mesh-adapter-url']").type('localhost:10002{enter}');

      cy.get('[data-cy=btnSubmitMeshAdapter]').click();

      cy.wait('@postMeshManage');
      cy.contains('Adapter was configured!');
    });

    it('ping and Submit Istio', () => {
      cy.intercept('GET', '/api/system/adapters*').as('getAdapterPing');
      cy.intercept('POST', '/api/system/adapter/manage').as('postMeshManage');
      cy.intercept('GET', '/api/system/adapters').as('getMeshAdapters');

      cy.get('[data-cy=chipAdapterLocation]', { timeout: 30_000 })
        .contains('.MuiChip-label', 'localhost:10000')
        .click();

      cy.wait('@getAdapterPing');

      cy.get("[data-cy='mesh-adapter-url']").type('localhost:10000{enter}');

      cy.get('[data-cy=btnSubmitMeshAdapter]').click();

      cy.wait('@postMeshManage');
      cy.contains('Adapter was configured!');
    });
  });
});
