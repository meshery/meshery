import { MODELS, REGISTRANTS } from '@/constants/navigator';

/**
 * Retrieves filtered data for the details component based on the selected item ID.
 *
 * @param {Array} data - An array of data representing the tree.
 * @param {Array} selectedIdArr - An array containing the IDs of the selected item. As the length increases, it represents sub-child levels. For example, selectedIdArr[0] is the main parent, selectedIdArr[1] is a sub-child of it, and so on.
 * @returns {Object} - An object containing the selected component, model, and relationship.
 */
export const getFilteredDataForDetailsComponent = (data, selectedItemUUID, view) => {
  let selectedIdArr = selectedItemUUID.split('.');
  const selectedParent = data.find((parentItem) => parentItem.id === selectedIdArr[0]);
  let selectedComponent = [];
  let selectedModel = {};
  let selectedRelationship = [];

  selectedModel = view === REGISTRANTS ? selectedParent?.models : selectedParent;
  const componentKey = view === MODELS ? 'components' : 'relationships';
  selectedComponent = selectedModel?.[componentKey]?.find(
    (component) => component && `${selectedModel.id}.1.${component.id}` === selectedItemUUID,
  );
  selectedRelationship = selectedModel?.relationships?.find(
    (relationship) =>
      relationship && `${selectedModel.id}.2.${relationship.id}` === selectedItemUUID,
  );

  return { selectedComponent, selectedModel, selectedRelationship };
};
