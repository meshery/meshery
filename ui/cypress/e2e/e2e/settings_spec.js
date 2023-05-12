describe('Settings', () => {
  describe('Service Meshes', () => {
    beforeEach(() => {
      cy.selectProviderNone();
      cy.visit('/settings');
      cy.intercept('GET', '/api/system/adapters').as('getMeshAdapters');

      cy.get('[data-cy="tabServiceMeshes"]').click();
      cy.request('GET', '/api/system/adapters').then((response) => {
        cy.log(response.body);
      });
      cy.wait('@getMeshAdapters')
      cy.reload()
    })

    it('ping and submit Consul', () => {
      cy.intercept('GET', '/api/system/adapters*').as('getAdapterPing');
      cy.intercept('POST', '/api/system/adapter/manage').as('postMeshManage');
      cy.intercept('GET', '/api/system/adapters').as('getMeshAdapters');

      cy.get('[data-cy=chipAdapterLocation]')
        .contains('.MuiChip-label', 'localhost:10002')
        .click();

      cy.wait('@getAdapterPing');

      cy.get('.MuiInputBase-input')
        .type('localhost:10002{enter}');

      cy.get('[data-cy=btnSubmitMeshAdapter]')
        .click();

      cy.wait('@postMeshManage');
      cy.contains("Adapter was configured!")
    });

    it('ping and Submit Istio', () => {
      cy.intercept('GET', '/api/system/adapters*').as('getAdapterPing');
      cy.intercept('POST', '/api/system/adapter/manage').as('postMeshManage');
      cy.intercept('GET', '/api/system/adapters').as('getMeshAdapters');

      cy.get('[data-cy=chipAdapterLocation]')
        .contains('.MuiChip-label', 'localhost:10000')
        .click();

      cy.wait('@getAdapterPing');

      cy.get('.MuiInputBase-input')
        .type('localhost:10000{enter}');

      cy.get('[data-cy=btnSubmitMeshAdapter]')
        .click();

      cy.wait('@postMeshManage');
      cy.contains("Adapter was configured!")
    });
  });
});
