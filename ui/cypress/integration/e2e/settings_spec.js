describe('Settings', () => {
  describe('Environment', () => {
    beforeEach(() => {
      cy.intercept('GET', '/api/promGrafana/scan').as('getScan')

      cy.selectProviderNone()

      cy.visit('/settings')
      cy.wait('@getScan')
    })

    it('clicks on Discover Cluster and send a ping to the cluster', () => {
      cy.intercept('GET', '/api/config/sync').as('getConfigSync')
      cy.intercept('GET', '/api/k8sconfig/ping').as('getK8sConfigPing')

      cy.get('[data-cy="tabInClusterDeployment"]').click()
      cy.get('[data-cy="btnDiscoverCluster"]').click()

      cy.wait('@getConfigSync')

      cy
        .get('[data-cy="chipContextName"]')
        .should('have.text', 'kind-kind-cluster')
        .click()

      cy.wait('@getK8sConfigPing')

      cy
        .get('[data-cy="itemListContextName"] > .MuiListItemText-secondary')
        .should('have.text', 'kind-kind-cluster')
    });
  });
});
