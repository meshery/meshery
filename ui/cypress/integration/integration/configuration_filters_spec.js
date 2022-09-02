import { getConfigurationGridItemName } from '../../actionHelpers/service-mesh-configuration-management'

describe('Configuration Management', () => {
  beforeEach(() => {
    cy.selectProviderNone();
    // Prepare Stub Interception for Initial Filters
    cy.intercept('GET', '/api/filter**', { fixture : 'configuration/filter-stubs/initial-filters.json' }).as('getInitialFilters');
  })

  describe('Filters', () => {

    context('UI', () => {
      beforeEach(() => {
        cy.visit('/configuration/filters');
      })

      it('displays correct title', () => {
        const expectedPageTitle = 'FiltersBETA';
        cy.get('[data-cy="headerPageTitle"]').should('have.text', expectedPageTitle);
      })

      it('displays correct table title', () => {
        const expectedTableTitle = 'Filters';
        cy.get('[data-cy="table-view"]').click();
        cy.get('[data-cy="filters-grid"] [class*="tableHeader"]').should('have.text', expectedTableTitle)
      })
    })

    context('Upload File', () => {
      beforeEach(() => {
        cy.intercept('POST', '/api/filter**', {}).as('postFilter');
        cy.visit('/configuration/filters');
        cy.wait('@getInitialFilters')
      })

      it('send correct POST /api/filter request', () => {
        const testFilePath = 'configuration/wasm-filters-v0.1.2/http-auth-envoy.yaml'
        const expectedUploadedFilterName = 'Test Filter';

        // Prepare Stub Interception for Post Upload Filters
        cy.intercept('GET', '/api/filter**', { fixture : 'configuration/filter-stubs/post-upload-filters.json' }).as('getPostUploadFilters');

        // Load test file fixture data
        cy.fixture(testFilePath).then((expectedContent) => {
          // Custom command 'attachFile' is provided by https://www.npmjs.com/package/cypress-file-upload#html5-file-input
          // It internally calls https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/dispatchEvent
          // with a Custom 'change' input event.
          //import button shifted to Filters main UI under service_mesh_config_management_spec.js
          cy.get('[data-cy="import-button"]').click();
          cy.get('[data-cy="file-upload-button"]').attachFile(testFilePath);

          cy.wait('@postFilter').then(interception => {

            cy.wrap(interception.request).then(req => {
              const body = JSON.parse(req.body);
              expect(body).to.have.nested.property('filter_data.filter_file');
              expect(body.filter_data.filter_file).to.eq(expectedContent);
              expect(body).to.have.property('save');
              expect(body.save).to.eq(true);
            })
          });

        })

        cy.wait('@getPostUploadFilters')
        cy.get('[data-cy="table-view"]').click();
        getConfigurationGridItemName(1).should('have.text', expectedUploadedFilterName)
      })
    })
  });
});
