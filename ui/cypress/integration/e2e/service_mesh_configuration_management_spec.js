import { getConfigurationGridItemName } from '../../actionHelpers/service-mesh-configuration-management'

// Shared Expect
const filtersAndApplicationsExpect = (body, itemType) => expect(body).to.have.nested.property(`${itemType}_data.${itemType}_file`)

// Test Data
const configurationTests = [
  {
    itemType : "filter",
    testFilePath : "configuration/wasm-filters-v0.1.2/http-auth-envoy.yaml",
    expectedUploadConfigItemName : "Test Filter",
    expectUploadRequestBody : filtersAndApplicationsExpect
  },
  // below deploy pattern test returns this error (500 Internal Server Error):
  // AssertionError: expected 'rpc error: code = Unknown desc = no matches for kind "EnvoyFilter"
  // in version "networking.istio.io/v1alpha3"\n' to equal ''
  // {
  //   itemType : "pattern",
  //   testFilePath : "configuration/service-mesh-patterns-samples/IstioFilterPattern-commit-b23b219.yaml",
  //   expectedUploadConfigItemName : "IstioFilterPattern-commit-b23b219.yaml",
  //   expectUploadRequestBody : filtersAndApplicationsExpect
  // },
  // below deploy application test returns this error (500 Internal Server Error) because
  // incorrect request payload to POST /api/application/deploy endpoint instead of actual yaml file content
  // {
  //   itemType : "application",
  //   testFilePath : "configuration/service-mesh-patterns-samples/bookInfoPattern-commit-5eb2369.yaml",
  //   expectedUploadConfigItemName : "bookInfoPattern-commit-5eb2369.yaml",
  //   expectUploadRequestBody : filtersAndApplicationsExpect
  // },
];

// Test Template
const configurationTestTemplate = (itemType, testFilePath, expectedUploadConfigItemName, expectUploadRequestBody) => {
  describe(`${itemType}s`, () => {
    beforeEach(() => {
      // For E2E tests either 'None' provider needs to be selected,
      // authentication through a specific provider must be performed,
      // OR we must ensure CI environment has proper access token properly setup and enabled for current browser session.
      cy.selectProviderNone();
      // Interception for Get Filters to spy/wait/assert on actual server requests/responses
      cy.intercept("GET", `/api/${itemType}**`).as("getConfigItems");

      // Interception for Post Filter to spy/wait/assert on actual server requests/responses
      cy.intercept("POST", `/api/${itemType}**`).as("uploadConfigItem");

      // Interception for Post Filter Deploy to spy/wait/assert on actual server requests/responses
      cy.intercept("POST", `/api/${itemType}/deploy**`).as("deployConfigItem");

      // Visit current page under testing
      cy.visit(`/configuration/${itemType}s`);
      cy.wait("@getConfigItems");
    });
    it(`Deploys ${itemType}`, () => {
      // Load test file fixture data
      cy.fixture(testFilePath).then((expectedContent) => {
        // Custom command 'attachFile' is provided by https://www.npmjs.com/package/cypress-file-upload#html5-file-input
        // It internally calls https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/dispatchEvent
        // with a Custom 'change' input event.
        cy.get('[data-cy="import-button"]').click();
        cy.get('[data-cy="file-upload-button"]').attachFile(testFilePath);
        cy.wait("@uploadConfigItem").then((interception) => {
          cy.wrap(interception.request).then((req) => {
            const body = JSON.parse(req.body);
            // expect(body).to.have.nested.property(`${itemType}_data.${itemType}_file`);
            expectUploadRequestBody(body, itemType)
            expect(body[`${itemType}_data`][`${itemType}_file`]).to.eq(expectedContent);
            expect(body).to.have.property("save");
            expect(body.save).to.eq(true);
          });
        });

        cy.wait("@getConfigItems").then((interception) => {
          cy.wrap(interception.response).then((res) => {
            expect(res.statusCode).to.eq(200);
            const body = res.body;
            expect(body).to.have.property(`${itemType}s`);
            expect(body[`${itemType}s`][0][`${itemType}_file`]).to.eq(expectedContent);
          });
        });

        getConfigurationGridItemName(1).should("have.text", expectedUploadConfigItemName);

        cy.get('[data-cy="config-row-0"] [data-cy="deploy-button"]').click();

        cy.wait("@deployConfigItem").then((interception) => {
          cy.wrap(interception.request).then((req) => {
            const body = req.body;
            expect(body).to.eq(expectedContent);
          });
          cy.wrap(interception.response).then((res) => {
            expect(res.body).to.eq('')
            expect(res.statusCode).to.eq(200);
          });
        });
      });
    });
  });
};

describe("Configuration Management", () => {
  configurationTests.forEach(({ itemType, testFilePath, expectedUploadConfigItemName, expectUploadRequestBody }) =>
    configurationTestTemplate(itemType, testFilePath, expectedUploadConfigItemName, expectUploadRequestBody)
  );
});

