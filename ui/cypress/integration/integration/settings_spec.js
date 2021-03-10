describe("Settings", () => {
    describe("Service Meshes", () => {
      before(() => {
        cy.selectProviderNone()

        cy.visit("/settings");

        cy.get('[data-cy="tabServiceMeshes"]').click();
      });
      it("select, submit, and confirm", () => {
        cy.get(".MuiInputBase-input").type("mesherylocal.layer5.io:10000{enter}");

        cy.get("[data-cy=btnSubmitMeshAdapter]").click();

        cy.get("[data-cy=adapterSuccessSnackbar]").should('exist')
      });
    });
  });
