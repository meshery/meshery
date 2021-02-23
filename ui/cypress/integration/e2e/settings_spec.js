describe('Settings', () => {
  describe('Environment', () => {
    beforeEach(() => {
      cy.intercept('GET', '/api/promGrafana/scan').as('getScan')

      cy.selectProviderNone()

      cy.visit('/settings')
      cy.wait('@getScan')
    })

    it('click on Discover Cluster and send a ping to the cluster', () => {
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
    })
  })

  describe('Service Meshes', () => {
    beforeEach(() => {
      cy.intercept('GET', '/api/promGrafana/scan').as('getScan')
      cy.intercept('GET', '/api/mesh/adapters').as('getMeshAdapters')

      cy.selectProviderNone()
      cy.visit('/settings')
      cy.wait('@getScan')

      cy.get('[data-cy="tabServiceMeshes"]').click()
      cy.wait('@getMeshAdapters')
    })
    it('ping and submit Consul', () => {
      cy.intercept('GET', '/api/mesh/adapter/ping').as('getAdapterPing')
      cy.intercept('POST', '/api/mesh/manage').as('postMeshManage')
      cy.intercept('GET', '/api/mesh/adapters').as('getMeshAdapters')

      cy.get('[data-cy=chipAdapterLocation]')
        .contains('.MuiChip-label', 'mesherylocal.layer5.io:10002')
        .click()

      cy.wait('@getAdapterPing')

      cy.get('.MuiInputBase-input')
        .type('mesherylocal.layer5.io:10002{enter}')

      cy.get('[data-cy=btnSubmitMeshAdapter]')
        .click()

      cy.wait('@postMeshManage')
    })
  })
})
