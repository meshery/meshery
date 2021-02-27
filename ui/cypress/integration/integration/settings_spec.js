describe("Settings", () => {
    describe("Service Meshe", () => {
      beforeEach(() => {
        cy.selectProviderNone()
  
        cy.visit("/settings");
  
        cy.get('[data-cy="tabServiceMeshes"]').click();
      });
      it("Select and Submit", () => {
        cy.get(".MuiInputBase-input").type("mesherylocal.layer5.io:10000{enter}");
  
        cy.get("[data-cy=btnSubmitMeshAdapter]").click();      
      });
    });
  });