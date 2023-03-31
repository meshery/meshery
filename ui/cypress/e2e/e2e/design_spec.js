import { DESIGNER } from "../../support/constants"
import '@4tw/cypress-drag-drop'

describe("Designer Spec", () => {
  beforeEach(() => {
    cy.viewport(1500, 900) 
    cy.login();
    cy.setReleaseTag();
    cy.interceptCapabilities();
    window.localStorage.setItem("tab", 0) // its a bug in designer
    cy.setMode(DESIGNER)
    cy.intercept("/api/provider/extension*").as("extensionFileLoad")
    cy.visit("/extension/meshmap");
    cy.wait("@extensionFileLoad", { timeout: 20000 });
  })

  const argoRolloutDesign = "no-modify-argo-rollout-application-for-cypress";
  const cypressModifiedDesignName = "changed Name with cypress"

  it("Load MeshMap Design with a click", () => {
    cy.get("[data-cy='design-drawer']").click();
    cy.get("#MUIDataTableBodyRow-patterns-0", {timeout: 30000});
    cy.wait(2000);
    cy.intercept("/api/pattern*").as("patternLoad")
    cy.get("#MUIDataTableBodyRow-patterns-0").click(); //convention: MUIDataTableBodyRow + type  + rowIndex
    cy.wait("@patternLoad");
    // cy.get("[data-cy='progress-snackbar']").contains("Rendering your MeshMap...");
    cy.wait(2000);
    cy.get("body").then(body => {
      if (body.find("[aria-describedby='notistack-snackbar'] #notistack-snackbar").length > 0) {
        cy.get("[aria-describedby='notistack-snackbar'] #notistack-snackbar").should("not.contain", "Unable to render")
      }
    })
  });

  it.skip("Rename Design", () => {
    cy.get("#component-drawer-Application").should('be.visible').drag("#cy-canvas-container", {force: true});
    cy.wait(2000); // let it open the rjsf successfully
    cy.get("[data-cy='design-drawer']").click(); // to close the rjsf form by click event
    cy.intercept('/api/pattern').as('patternSave')
    cy.get("#design-name-textfield").focus().clear().type(cypressModifiedDesignName);
    cy.wait("@patternSave").then(() => {
      // move to drawer and check for update
      cy.get("[data-cy='design-drawer']").click();
      cy.get("#MUIDataTableBodyRow-patterns-0 p", {timeout: 30000});
      cy.wait(2500);
      cy.get("#MUIDataTableBodyRow-patterns-0 p").contains(cypressModifiedDesignName);
    })
  })
  
  it.skip("Search a design", () => {
    cy.get("[data-cy='design-drawer']").click();
    cy.get("#MUIDataTableBodyRow-patterns-0", {timeout: 30000})
    cy.wait(2000);
    cy.get("#MUIDataTableBodyRow-patterns-0").click();
    cy.intercept("/api/pattern*").as("patternSearch")
    cy.get('[data-test-id="Search"]').type(cypressModifiedDesignName);
    cy.wait("@patternSearch")
    cy.get("#MUIDataTableBodyRow-patterns-0").should("be.visible").contains(cypressModifiedDesignName);
  })

  it("Validate a design", () => {
    cy.get("[data-cy='design-drawer']").click();
    cy.get("#MUIDataTableBodyRow-patterns-0", {timeout: 30000})
    cy.wait(1500);
    cy.get("#MUIDataTableBodyRow-patterns-0").click();
    cy.wait(1500); // wait for design to load
    cy.get("#verify-design-btn").click();
    cy.get('[data-cy="validate-btn-modal"]').click();
    cy.contains("Validate");
    cy.contains("OK");
  })

  it.skip("Deploy and Undeploy a design", () => {
    cy.get("[data-cy='design-drawer']").click();
    cy.get("#MUIDataTableBodyRow-patterns-0", {timeout: 30000})
    cy.wait(2000);
    cy.get("#MUIDataTableBodyRow-patterns-0").click();
    cy.get('[data-test-id="Search"]').type(argoRolloutDesign);
    cy.intercept("/api/pattern*").as("patternPost")
    cy.wait(1500);
    cy.get("#MUIDataTableBodyRow-patterns-0").should("be.visible").contains(argoRolloutDesign);
    cy.wait(2000);
    cy.get("#MUIDataTableBodyRow-patterns-0").click({force: true}).wait("@patternPost");
    cy.wait(2000);

    // rendering done up until this point
    cy.get("body").then(body => {
      if (body.find("[aria-describedby='notistack-snackbar'] #notistack-snackbar").length > 0) {
        cy.get("[aria-describedby='notistack-snackbar'] #notistack-snackbar").should("not.contain", "Unable to render")
      }
    })

    // modal opens
    cy.get("#deploy-design-btn").click();
    cy.get('[data-cy="deploy-btn-modal"]').click();
    
    cy.intercept("/api/pattern/deploy*").as("patternDeploy")
    cy.get('[data-cy="deploy-btn-confirm"]').click();
    cy.wait("@patternDeploy").then(() => {
      // cy.get("[data-cy='progress-snackbar']").contains("Deploying design");
      cy.get("body").then(body => {
        if (body.find("[aria-describedby='notistack-snackbar'] #notistack-snackbar").length > 0) {
          cy.get("[aria-describedby='notistack-snackbar'] #notistack-snackbar").should("not.contain", "Failed")
        }
      })
    })
    //Undeploy 
    cy.get('#undeploy-design-btn').click();
    cy.get('[data-cy="Undeploy-btn-modal"]').click();
    // modal opens
    cy.intercept("/api/pattern/deploy*").as("patternUndeploy")
    cy.get('[data-cy="deploy-btn-confirm"]').click();
    cy.wait("@patternUndeploy").then(() => {
      // cy.get("[data-cy='progress-snackbar']").contains("Deploying design");
      cy.get("body").then(body => {
        if (body.find("[aria-describedby='notistack-snackbar'] #notistack-snackbar").length > 0) {
          cy.get("[aria-describedby='notistack-snackbar'] #notistack-snackbar").should("not.contain", "Failed")
        }
      })
    })
  });
})
