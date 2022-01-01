/// <reference types="cypress" />

describe("Login", () => {
  before(()=>{
    const token = Cypress.env('token')
    cy.setCookie("meshery-provider", "Meshery")
    cy.setCookie("token", token)
    window.localStorage.setItem("mode", "designer")
    window.localStorage.setItem("tab", 0)
  })

  it("Visit MeshMap Designer", () => {
    cy.visit("/")
    cy.get('[data-cy="MeshMap"]').click();
  })

})