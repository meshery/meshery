describe("Settings", () => {
  describe("Service Meshes", () => {
    before(() => {
      cy.visit("/settings");
      cy.intercept("GET", "/api/system/adapters", { fixture: "getMeshAdapters.json" }).as("getMeshAdapters");
      cy.intercept("POST", "/api/system/adapter/manage", { fixture: "postMeshManage.json" }).as("postMeshManage");
      cy.get('[data-cy="tabServiceMeshes"]').click();
    });

    it("select, submit, and confirm", () => {
      cy.get("[data-cy='mesh-adapter-url']").type("localhost:10000{enter}");
      cy.get("[data-cy=btnSubmitMeshAdapter]").click();
      cy.wait("@postMeshManage");
      cy.wait("@getMeshAdapters");
      cy.get("[data-cy=adapterSuccessSnackbar]").should("exist");
    });
  });
});
