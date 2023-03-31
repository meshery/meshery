/// <reference types="cypress" /> // Add and extra '/' to load the cypress auto-complete feature following npm i cypress
import { VISUALIZER } from "../../../support/constants"

// DO NOT INTERCEPT CAPABILITIES WITH FIXTURES 
describe("Visualizer Spec", () => {
  beforeEach(() => {
    cy.login();
    cy.setReleaseTag();
    window.localStorage.setItem("tab", 0);
    const token = Cypress.env('token');
    cy.request(`/api/user/token?token=${token}`);
  })

  it.skip("Visit MeshMap Visualizer, check gql status", () => {
    cy.setMode(VISUALIZER)
    cy.visit("/extension/meshmap")
    cy.intercept("/api/provider/extension*").as("extensionFileLoad")
    cy.wait("@extensionFileLoad");
    cy.get("body").then(body => {
      if (body.find('[data-cy="modal-close-btn"]').length > 0) {
        cy.get('[data-cy="modal-close-btn"]').click();
      }
    })
    cy.contains("MeshMap")
    cy.wait(1500)
    cy.get("body").then(body => {
      if (body.find(`[data-cy="modal-close-btn"]`).length > 0) {
        // ensures the message of graphQL error, if any
        cy.get('[data-cy="modal-close-btn"]').click();
      }
    })
    //tabs
    cy.contains("Details")
    cy.contains("Metrics")
    cy.contains("Actions")
    cy.get("body").should("not.contain", "Uh-oh!")
  });
});
