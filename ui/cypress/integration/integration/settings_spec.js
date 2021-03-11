describe("Settings", () => {
    describe("Service Meshes", () => {
      before(() => {
        cy.intercept('GET', '/api/mesh/adapters').as('getMeshAdapters')

        cy.visit("/settings");
        cy.get('[data-cy="tabServiceMeshes"]').click();
        cy.wait('@getMeshAdapters')
      });

      it("select, submit, and confirm", () => {
        cy.intercept('POST', '/api/mesh/manage', { fixture: 'stats.json' }).as('submitMeshAdapter')

        cy.get(".MuiInputBase-input").type("mesherylocal.layer5.io:10000{enter}");
        cy.get("[data-cy=btnSubmitMeshAdapter]").click();
        cy.wait('@submitMeshAdapter')
        cy.get("[data-cy=adapterSuccessSnackbar]").should('exist')
      });
    });
  });
