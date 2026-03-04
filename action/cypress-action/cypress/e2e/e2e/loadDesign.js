/// <reference types="cypress" />
/// <reference types="../../support" />

import { TIME, canvasContainer } from "../../support/constants";
import {
    beforeEachCallbackForCustomUrl,
    doSnapshotSetup,
    waitFor,
} from "../../support/helpers";

const InfraShot = (theme) => {
    return describe(`Infra Shot Automated Runner ${theme} Mode`, () => {
        beforeEach(() =>
            beforeEachCallbackForCustomUrl(
                `/extension/meshmap?mode=design&design=${getDesignId()}&render=full`,
                theme,
            ),
        );

        it(`take light mode infra shot`, () => {
            const designId = getDesignId();
            waitForDesignRender();
            cy.window().then((window) => {
                cy.wait(TIME.XLARGE *4);
                captureSnapshot({
                    window,
                    designId: designId,
                    theme,
                });
            });
        });
    });
};

const getDesignId = () => {
    return Cypress.env("applicationId").replace(/['"]+/g, "");
};

const waitForDesignRender = () => {
    waitFor(canvasContainer.query, { timeout: 60_000 });
    cy.wait(TIME.X4LARGE * 2);
};

const snapshotPath = (designId, theme) => {
    return `snapshot-${theme}`;
};

const captureSnapshot = ({ window, designId, theme }) => {
    console.log("Taking snapshot", designId, theme);
    //removeWidgets();
    cy.window()
        .its("cyto", { timeout: 10000 })
        .should("exist")
        .then((cytoscape) => {
            cytoscape.fit();
            cytoscape.center();
        });
    const path = snapshotPath(designId, theme);
    cy.wait(2000);

    cy.get("main", { timeout: 10 * 1000 })
        .should("exist")
        .screenshot(path, {
            //blackout: [".hide-from-snapshot"], // hides elements before screenshot
            scale: true,
        });
    console.log(`Snapshot taken at ${path}`);
};

const safeHide = (selector) => {
    cy.get("body").then(($body) => {
        if ($body.find(selector).length) {
            cy.get(selector).invoke("css", "visibility", "hidden");
        }
    });
};

const removeWidgets = () => {
    safeHide("#action-toolbar");
    safeHide("#kanvas-bottom-dock");
    safeHide("#left-navigation-bar");
    safeHide("#top-navigation-bar");
    safeHide(".hide-from-snapshot");
};

["light", "dark"].forEach((theme) => {
    InfraShot(theme);
});
