describe('Provider UI', () => {
    beforeEach(() => {
        cy.intercept('GET', '/api/providers', { fixture: 'providers.json' }).as('getProviders')
        cy.visit('/')
        cy.wait('@getProviders')
    })

    it('renders provider component', () => {
        cy
            .get('[data-cy=root]')
            .should('exist')
    })

    it('selects None as provider', () => {
        cy
            .get('[data-cy="select_provider"]').click()
            .get('#split-button-menu').contains('None')
    })

    it('opens and closes modal to learn more about providers', () => {
        cy
            .get('[data-cy="providers-tooltip"]').click()
            .get('[data-cy="providers-modal"]').should('exist')
            .get('[data-cy="providers-modal-button-ok"]').click()
            .get('[data-cy="providers-modal"]').should('not.exist')
    });
})