
const getConfigurationGridItem = (itemNumber) => {
  const index = itemNumber - 1
  return cy.get(`[data-cy="config-row-${index}"]`)
}

export const getConfigurationGridItemName = (itemNumber) => {
  const itemNameColIndex = '0'
  return getConfigurationGridItem(itemNumber).find(`[data-colindex="${itemNameColIndex}"]`)
}