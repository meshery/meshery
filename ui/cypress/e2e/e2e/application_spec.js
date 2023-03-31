import {DESIGNER} from "../../support/constants"
import '@4tw/cypress-drag-drop'

describe("Application Spec", () => {
  beforeEach(() => {
    cy.login();
    cy.setReleaseTag();
    cy.interceptCapabilities();
    window.localStorage.setItem("tab", 0)
    cy.setMode(DESIGNER)
    cy.intercept("/api/provider/extension*").as("extensionFileLoad");
    cy.visit("/extension/meshmap");
    cy.wait("@extensionFileLoad", { timeout: 20000 });
    cy.get('[data-cy="application-drawer"]').click();
    cy.get("#MUIDataTableBodyRow-applications-0", {timeout: 30000}).should("be.visible"); // start tests only when the application table is populated
  })

  const bookInfoUrlUploadedName = "k8s.yaml" // sample name of an application
  const bookInfoApp = "bookInfo-istio-cy"; // book-info istio application to search

  it.skip("Render MeshMap Application", () => {
    cy.contains("Applications")
    cy.get("#MUIDataTableBodyRow-applications-0").click({force: true}); //convention: MUIDataTableBodyRow + type  + rowIndex
    cy.wait(2000);
    cy.get("body").then(body => {
        if (body.find("[aria-describedby='notistack-snackbar'] #notistack-snackbar").length > 0) {
          cy.get("[aria-describedby='notistack-snackbar'] #notistack-snackbar").should("not.contain", "Error")
        }
      })
  });

  it.skip("Rename and Saving Application", () => { // renaming application is skipped until the import test is present, because it may change the users app name
    cy.get("#MUIDataTableBodyRow-applications-0", {timeout: 30000}).should("be.visible").click(); // drop the first application
    cy.wait(3000); // let it be dropped on the canvas
    cy.get("#design-name-textfield").type(bookInfoUrlUploadedName);
    cy.intercept('/api/application*').as('applicationSave')
 
    // TODO - Saving Application request intercept need to be fixed.
    cy.wait("@applicationSave").then(() => {
        // move to drawer and check for update
        cy.get("[data-cy='application-drawer']").click();
        cy.wait(5000); // wait for seconds, because the subscritions cannot be tracked for now
        cy.get("#MUIDataTableBodyRow-applications-0 p").contains(bookInfoUrlUploadedName);
      })
  })

  it.skip("Validate an application", () => {
    cy.get("[data-cy='application-drawer']").click();
    cy.get("#MUIDataTableBodyRow-applications-0", {timeout: 30000})
    cy.wait(1500);
    cy.get("#MUIDataTableBodyRow-applications-0").click();
    cy.wait(2000);
    cy.get("#verify-design-btn").click();
    cy.get('[data-cy="validate-btn-modal"]').click();
    cy.contains("OK");
  })

  it.skip("Deploy and Undeploy an Application", () => {
    cy.get('[data-test-id="Search"]').type(bookInfoApp);
    cy.intercept("/api/application*").as("applicationPost")
    cy.wait("@applicationPost")
    cy.wait(1500);
    cy.get("#MUIDataTableBodyRow-applications-0").should("be.visible").contains(bookInfoApp).click();
    cy.wait(2000);
    cy.get("body").then(body => {
      if (body.find("[aria-describedby='notistack-snackbar'] #notistack-snackbar").length > 0) {
        cy.get("[aria-describedby='notistack-snackbar'] #notistack-snackbar").should("not.contain", "Unable to render")
      }
    })
    cy.get("#deploy-design-btn").click();
    cy.get('[data-cy="deploy-btn-modal"]').click();
    // modal opens
    cy.intercept("/api/pattern/deploy*").as("applicationDeploy")
    cy.get('[data-cy="deploy-btn-confirm"]').click();
    cy.wait("@applicationDeploy").then(() => {
      // cy.get("[data-cy='progress-snackbar']").contains(`Deploying application: ${bookInfoApp}`);
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
    cy.intercept("/api/pattern/deploy*").as("applicationUndeploy")
    cy.get('[data-cy="deploy-btn-confirm"]').click();
    // cy.get("[data-cy='progress-snackbar']").contains(`Undeploying application: ${bookInfoApp}`);
    cy.wait("@applicationUndeploy").then(() => {
      cy.get("body").then(body => {
        if (body.find("[aria-describedby='notistack-snackbar'] #notistack-snackbar").length > 0) {
         cy.get("[aria-describedby='notistack-snackbar'] #notistack-snackbar").should("not.contain", "Failed")
        }
      })
    })
  });

  it("Search an Application", () => {
    cy.intercept("/api/application*").as("applicationSearch")
    cy.get('[data-test-id="Search"]').type(bookInfoApp);
    cy.wait("@applicationSearch")
    cy.get("#MUIDataTableBodyRow-applications-0").should("be.visible").contains(bookInfoApp);
  })
})
