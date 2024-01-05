/// <reference types="cypress" />

describe('Lifecycle Service Mesh', () => {
  beforeEach(() => {
    cy.selectProviderNone();
  });

  const selectServiceMeshType = (adapterLocation) => {
    cy.get('[data-cy="lifecycle-service-mesh-type"]').click();
    cy.get('[role="listbox"]').within(() => {
      cy.contains(adapterLocation).click();
    });
  };

  const mesheryAdapters = [
    { adapterName: 'Istio', adapterPort: '10000', deploy: false },
    { adapterName: 'Consul', adapterPort: '10002', deploy: false },
    // { adapterName: "NGINX Service Mesh", adapterPort: "10010", deploy: true },
  ];

  mesheryAdapters.forEach(({ adapterName, adapterPort, deploy }, index) => {
    const ADAPTER_LOCATION = `localhost:${adapterPort}`;
    it(`User can Configure Existing ${adapterName} adapter through Mesh Adapter URL from Management page`, () => {
      // Settings > Adapters Page
      cy.visit('/settings#service-mesh');

      if (!deploy) {
        // "Mesh Adapter URL" Dropdown
        cy.get('[data-cy="mesh-adapter-url"]').type(`${ADAPTER_LOCATION}{enter}`);
        // "Connect" Button
        cy.get('[data-cy="btnSubmitMeshAdapter"]').click();
        cy.contains('Adapter was configured!', { timeout: 10_000 }).should('exist');
      } else {
        // TODO: Implement test code for Deploying from Settings > Adapters Page's "Available Mesh Adapter" Dropdown
      }

      // Lifecycle > Service Mesh Page
      cy.visit('/management/service-mesh');
      selectServiceMeshType(ADAPTER_LOCATION);
      // "Select Service Mesh Type" Dropdown
      cy.get('[data-cy="lifecycle-service-mesh-type"]').should('contain.text', ADAPTER_LOCATION);
    });

    it(`User can ping ${adapterName} Adapter`, () => {
      cy.visit('/');

      // Lifecycle > Service Mesh Page
      cy.visit('/management');
      selectServiceMeshType(ADAPTER_LOCATION);
      cy.get('[data-cy="lifecycle-service-mesh-type"]', { timeout: 10_000 }).should(
        'contain.text',
        ADAPTER_LOCATION,
      );
      // "Manage Service Mesh" Card's Ping adapter chip
      cy.get('[data-cy="adapter-chip-ping"]').click();
      cy.contains('Adapter pinged!').should('exist');
    });
  });
});
