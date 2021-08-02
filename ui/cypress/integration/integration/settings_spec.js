describe("Settings", () => {
  describe("Service Meshes", () => {
    before(() => {
      cy.visit("/settings");
      cy.get('[data-cy="tabServiceMeshes"]').click();
    });

    it("select, submit, and confirm", () => {
      cy.intercept("GET", "/api/mesh/adapters", { fixture: "getMeshAdapters.json" }).as("getMeshAdapters");
      cy.intercept("POST", "/api/mesh/manage", { fixture: "postMeshManage.json" }).as("postMeshManage");

      cy.get(".MuiInputBase-input").type("mesherylocal.layer5.io:10000{enter}");

      cy.get("[data-cy=btnSubmitMeshAdapter]").click();
      cy.wait("@postMeshManage");
      cy.wait("@getMeshAdapters");
      cy.get("[data-cy=adapterSuccessSnackbar]").should("exist");
    });
  });
});
