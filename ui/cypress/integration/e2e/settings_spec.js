describe('Settings', () => {
  describe('Environment', () => {
    beforeEach(() => {
      cy.selectProviderNone();

      cy.visit('/settings');
    });

    it('click on Discover Cluster and send a ping to the cluster', () => {

      cy.intercept('GET', '/api/system/sync').as('getConfigSync');
      cy.intercept('GET', '/api/system/kubernetes/ping').as('getK8sConfigPing');

      cy.get('[data-cy="tabInClusterDeployment"]').click();
      cy.get('[data-cy="btnDiscoverCluster"]').click();

      cy.wait('@getConfigSync');

      // TODO: re-enable once we have chips for pinging the cluster
      // cy
      //   .get('[data-cy="chipContextName"]')
      //   .should('have.text', 'kind-kind-cluster')
      //   .click();

      // cy.wait('@getK8sConfigPing');

      cy
        .get('[data-cy="itemListContextName"] > .MuiListItemText-secondary')
        .should('have.text', 'kind-kind-cluster');
    });
  });

  describe('Service Meshes', () => {
    beforeEach(() => {
      cy.intercept('GET', '/api/system/adapters').as('getMeshAdapters');

      cy.selectProviderNone();
      cy.visit('/settings');

      cy.get('[data-cy="tabServiceMeshes"]').click();
      cy.wait('@getMeshAdapters');
    })
    it.skip('ping and submit Consul', () => {
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
    });
  });
});
