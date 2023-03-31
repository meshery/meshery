import { DESIGNER, extension, MESHMAP_PATH } from "../support/constants"


export function waitFor(str) {
  return "@" + str;
}

export function id(str) {
  return "#" + str
}

export const beforeEachCallback = () => {
  cy.viewport(1500, 900);
  cy.login();
  cy.setReleaseTag();
  cy.interceptCapabilities();
  cy.setMode(DESIGNER);
  cy.intercept(extension.path).as(extension.alias);
  cy.visit(MESHMAP_PATH)
  cy.wait(waitFor(extension.alias));
}