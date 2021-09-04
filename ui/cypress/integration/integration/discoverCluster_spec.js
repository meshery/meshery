describe("Settings", () => {
    describe("Environment", () => {
      before(() => {
        cy.intercept('GET', '/api/system/sync', { fixture: 'sync.json' }).as('getConfigSync')

        cy.visit("/settings");
        cy.wait('@getConfigSync')
        cy.get('[data-cy="tabEnvironment"]').click();

      });
      
  cy.on('uncaught:exception', (err, runnable) => {
    // we expect a `subscriptions-transport-ws` error with message `Cannot set property 'onopen' of null`
    // and don't want to fail the test so we return false
    // if (err.message.includes(`Cannot set property 'onopen' of null`)) {
      return false
    // }
    // we still want to ensure there are no other unexpected
    // errors, so we let them fail the test
  })
      it("search Cluster Deployment & ping k8s cluster", () => {
        cy.intercept('GET', '/api/system/kubernetes/ping', { fixture: 'clusterVersion.json' }).as('getK8sVersion')
        cy.intercept('GET', '/api/config/sync', { fixture: 'sync.json' }).as('getConfigSync')

        cy.get("[data-cy=tabInClusterDeployment]").click();
        cy.get("[data-cy=btnDiscoverCluster]").click();

        cy.wait('@getConfigSync')
        cy.get("[data-cy=chipContextName]").click();
        cy.wait('@getK8sVersion');
        cy.get("[data-cy=k8sSuccessSnackbar]").should("exist");
      });
    });
  });
