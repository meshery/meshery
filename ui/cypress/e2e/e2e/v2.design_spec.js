// type definitions for Cypress object "cy"
/// <reference types="cypress" />
/// <reference types="../../support" /> 

import { canvasLayer0, designEndpoint, TIME } from "../../support/constants";
import { beforeEachCallback, waitFor } from "../../support/helpers";

describe("canvas test", () => {
  beforeEach(beforeEachCallback)

  it("add element on canvas with double click", () => {
    // double click on the center of the canvas to popup the dialog
    cy.get(canvasLayer0.query).as(canvasLayer0.alias)
    cy.get(waitFor(canvasLayer0.alias)).dblclick('center', { force: true, multiple: true, })

    // the popup wizard should exist
    cy.get('[data-cy="canvas-component-menu"]').should("exist");

    // get the first category from list and click on it
    cy.get('[data-cy="canvas-component-menu"]')
      .find('[data-cy="category-menu"]')
      .children().first().click()

    cy.intercept({ url: designEndpoint.path, method: "POST" }).as(designEndpoint.alias) // create an interception alias for POST Method
    cy.get('[data-cy="canvas-component-menu"]')
      .find('[data-cy="model-table"]')
      .children().should("have.length.above", 0) // check for model length > 0
      .first().click() // select the first Model
      .wait(TIME.MEDIUM) // wait for ajax
      .find('[data-cy="component-table-body"]')
      .children().first() // get the first model selected
      .children().eq(1) // get the node svg selected
      .click() // click to add the node on canvas

    let designId = null;
    cy.wait(designEndpoint.wait)
      // .its("request").should("include", "cytoscape_json").should("include", "name").its("response").should("include", "id")
      .then(interceptionObj => {
        console.log("result of endpoint", interceptionObj)
        const reqBody = interceptionObj.request.body;
        expect(reqBody).to.have.ownProperty("cytoscape_json")
        expect(reqBody).to.have.ownProperty("name")
        expect(reqBody).to.have.include({ save: true })

        const response = interceptionObj.response
        expect(response).to.include({ statusCode: 200 })

        expect(response.body).to.be.an("array") // response body structure should be array
        designId = response.body[0]?.id

        // do cleanup by deleting the design
        cy.deleteDesign(designId)

      }) // ensuring on the click event the design endpoint is called

    cy.window().its("cyto").should("exist") // there should be a global cytoscape property exposed
    cy.window().then(window => {
      const cyto = window.cyto;
      cy.wrap(cyto.nodes().length).should("equal", 1)
    })
  })
})
