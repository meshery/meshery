import { DESIGNER } from "../../support/constants"
// import '@4tw/cypress-drag-drop' // Filters drag and drop on canvas is disabled for now.

describe("Filter Spec", () => {
    beforeEach(() => {
        cy.login();
        cy.setReleaseTag();
        cy.interceptCapabilities();
        window.localStorage.setItem("tab", 0) // its a bug in designer
        cy.setMode(DESIGNER)
        cy.visit("/extension/meshmap");
        cy.intercept("/api/provider/extension*").as("extensionFileLoad")
        cy.wait("@extensionFileLoad", { timeout: 20000 });
        cy.get('[data-cy="filter-drawer"]').click();
        cy.get("#MUIDataTableBodyRow-filters-0", {timeout: 30000}).should("be.visible"); // start tests only when the filter table is populated
      })

      const filterName = "Basic Auth for Istio" // sample name of an filter

      it.skip("Render Filters", () => { // filters are not rendered on canvas for now
        cy.contains("Filters")
        cy.get("#MUIDataTableBodyRow-filters-0").should("be.visible"); //convention: MUIDataTableBodyRow + type  + rowIndex
        cy.get("body").then(body => {
        if (body.find("[aria-describedby='notistack-snackbar'] #notistack-snackbar").length > 0) {
          cy.get("[aria-describedby='notistack-snackbar'] #notistack-snackbar").should("not.contain", "Error")
        }
      })
      });

      // Filters are not clickable for now. So they are skipped.
      it.skip("Rename Filter", () => {
        cy.get("#design-name-textfield").type(filterName);
        cy.intercept('/api/filter').as('filterSave')
        cy.wait("@filterSave").then(() => {
          // move to drawer and check for update
          cy.get("[data-cy='filter-drawer']").click();
          cy.wait(5000); // wait for seconds, because the subscritions cannot be tracked for now
          cy.get("#MUIDataTableBodyRow-filters-0 p").contains(filterName);
        })
      })

      // No Filters are being rendered or seeded so no filters will be found searched.
      it.skip("Search a filter", () => {
        cy.get("[data-cy='filter-drawer']").click();
        cy.get('[data-test-id="Search"]').type(filterName);
        cy.intercept("/api/filter*").as("filterSearch")
        cy.wait("@filterSearch")
        cy.get("#MUIDataTableBodyRow-filters-0").should("be.visible").contains(filterName);
      });
})