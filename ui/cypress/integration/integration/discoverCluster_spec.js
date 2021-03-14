describe("Settings", () => {
    describe("Environment", () => {
      before(() => {
        cy.intercept('GET', '/api/config/sync', { fixture: 'sync.json' }).as('getConfigSync')

        cy.visit("/settings");
        cy.wait('@getConfigSync')
        cy.get('[data-cy="tabEnvironment"]').click();

      });

      it("search Cluster Deployment & ping k8s cluster", () => {
        cy.intercept('GET', '/api/k8sconfig/ping', { fixture: 'clusterVersion.json' }).as('getK8sVersion')
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
