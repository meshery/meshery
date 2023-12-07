describe('Visit Meshery', () => {
  it('Visits meshery', () => {
    cy.selectProviderNone();
    cy.visit('/');
  });
});

describe('Visit Meshery settings', () => {
  it('Visits meshery settings page', () => {
    cy.selectProviderNone();
    cy.visit('/settings');
  });
});
