/// <reference types="cypress" />
import { aliasQuery } from "../../actionHelpers/graphql-test-utils";

function doGeneralIstioLifecycleChecks() {
  cy.contains("Lifecycle");
  cy.contains("Manage Service Mesh Lifecycle", { timeout: 10_000 });
  cy.contains("Namespace");
  cy.contains("Manage Sample Application Lifecycle");
  cy.contains("Apply Service Mesh Configuration");
  cy.contains("Validate Service Mesh Configuration");
  cy.contains("Apply Custom Configuration");
}

function verifyAllIstioAddons() {
  cy.contains("Grafana");
  cy.contains("Jaeger");
  cy.contains("Kiali");
  cy.contains("Prometheus");
  cy.contains("Zipkin");
}

function doGeneralConsulLifecycleChecks() {
  cy.contains("Lifecycle");
  cy.contains("Manage Service Mesh Lifecycle");
  cy.contains("Namespace");
  cy.contains("Manage Sample Application Lifecycle");
  cy.contains("Apply Custom Configuration");
  cy.get('body').should("not.contain", "Apply Service Mesh Configuration");
}

describe("Test Lifecycle Paths With None Provider", () => {
  beforeEach(() => {
    cy.selectProviderNone();
  })

  const ISTIO_URL = "/management?adapter=localhost:10000";
  const CONSUL_URL = "/management?adapter=localhost:10002";


  it("visit Management page", () => {
    cy.visit("/management");
    cy.contains("Lifecycle");
  })

  // assumption: istio adapter is live
  it("Visit Istio by link", () => {
    cy.visit(ISTIO_URL);
    doGeneralIstioLifecycleChecks(cy);
    verifyAllIstioAddons(cy);
  })

  it.only("Visit Istio by button Click", () => {
    cy.intercept('POST', '/api/system/graphql/query', (req) => {
      // Queries
      aliasQuery(req, "AddonsStatus")
    })
    cy.visit("/")
    cy.get('[data-cy="lifecycle"]').click();
    cy.wait('@gqlAddonsStatusQuery', { timeout: 40_000 })
    cy.get('[data-cy="istio"]').click();
    doGeneralIstioLifecycleChecks(cy);
    verifyAllIstioAddons(cy);
  })

  it("Visit Consul by link", () => {
    cy.visit(CONSUL_URL);
    doGeneralConsulLifecycleChecks(cy);
  })

  it.only("Visit Consul by button Click", () => {
    cy.intercept('POST', '/api/system/graphql/query', (req) => {
      // Queries
      aliasQuery(req, "AddonsStatus")
    })
    cy.visit("/")
    cy.get('[data-cy="lifecycle"]').click();
    cy.wait('@gqlAddonsStatusQuery', { timeout: 30_000 })
    cy.get('[data-cy="consul"]').click();
    doGeneralConsulLifecycleChecks(cy);
  })

});
