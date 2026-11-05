import { useContext } from 'react';
import { RegistryModalContext } from '../context/RegistryModalContextProvider';

/**
 * Custom hook to interact with the Registry Modal
 * @returns {Object} Registry modal controls and helpers
 */
export const useRegistryModal = () => {
  const context = useContext(RegistryModalContext);

  if (!context) {
    throw new Error('useRegistryModal must be used within a RegistryModalContextProvider');
  }

  /**
   * Opens the registry modal with specific parameters for deep-linking
   * @param {Object} params - Parameters for opening the modal
   * @param {string} params.tab - The tab to open (Models, Components, Relationships, Registrants)
   * @param {string} params.searchText - Search text to apply
   * @param {string} params.selectedItemUUID - UUID of the item to select and highlight
   */
  const openRegistryWithParams = (params = {}) => {
    context.openModalWithParams(params);
  };

  /**
   * Opens the registry modal to a specific model
   * @param {string} modelId - The model ID to navigate to
   * @param {string} modelName - The model name to search for (optional)
   */
  const openToModel = (modelId, modelName = '') => {
    openRegistryWithParams({
      tab: 'Models',
      selectedItemUUID: modelId,
      searchText: modelName,
    });
  };

  /**
   * Opens the registry modal to a specific component
   * @param {string} componentId - The component ID to navigate to
   * @param {string} componentName - The component name to search for (optional)
   */
  const openToComponent = (componentId, componentName = '') => {
    openRegistryWithParams({
      tab: 'Components',
      selectedItemUUID: componentId,
      searchText: componentName,
    });
  };

  /**
   * Opens the registry modal to a specific relationship's group
   * @param {string} relationshipId - The relationship ID to navigate to
   * @param {string} searchText - Optional search text to apply
   */
  const openToRelationship = (relationshipId) => {
    openRegistryWithParams({
      tab: 'Relationships',
      selectedItemUUID: relationshipId,
    });
  };

  /**
   * Opens the registry modal to a specific registrant
   * @param {string} registrantId - The registrant ID to navigate to
   * @param {string} registrantName - The registrant name to search for (optional)
   */
  const openToRegistrant = (registrantId, registrantName = '') => {
    openRegistryWithParams({
      tab: 'Registrants',
      selectedItemUUID: registrantId,
      searchText: registrantName,
    });
  };

  return {
    open: context.open,
    openModal: context.openModal,
    closeModal: context.closeModal,

    openRegistryWithParams,
    openToModel,
    openToComponent,
    openToRelationship,
    openToRegistrant,

    selectedView: context.selectedView,
    searchText: context.searchText,
    selectedItemUUID: context.selectedItemUUID,

    setSelectedView: context.setSelectedView,
    setSearchText: context.setSearchText,
    setSelectedItemUUID: context.setSelectedItemUUID,
  };
};

export default useRegistryModal;
