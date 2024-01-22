import { MODELS, REGISTRANTS } from '@/constants/navigator';
import _ from 'lodash';

/**
 * Retrieves filtered data for the details component based on the selected item ID.
 *
 * @param {Array} data - An array of data representing the tree.
 * @param {Array} selectedItemUUID - Node ID of selected element in tree.
 * @param {string} view - Selected tab in registry.
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

export const removeDuplicateVersions = (data) => {
  const groupedModels = _.groupBy(data, 'name');

  const result = _.reduce(
    groupedModels,
    (acc, models, name) => {
      const uniqueVersions = _.uniqBy(models, 'version');

      const versionDataArray = uniqueVersions.map((model) => {
        return {
          version: model.version,
          components: _.uniq(model.components),
          relationships: _.uniq(model.relationships),
          ...model,
        };
      });

      const existingModel = acc.find((m) => m.name === name);

      if (existingModel) {
        existingModel.version = _.union(
          existingModel.version,
          uniqueVersions.map((model) => model.version),
        );
        existingModel.versionBasedData = existingModel.versionBasedData.concat(versionDataArray);
      } else {
        const selectedModel = models[0];
        acc.push({
          ...selectedModel,
          version: uniqueVersions.map((model) => model.version),
          versionBasedData: versionDataArray,
        });
      }

      return acc;
    },
    [],
  );

  return result;
};
