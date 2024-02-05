import { MODELS, REGISTRANTS, COMPONENTS, RELATIONSHIPS } from '@/constants/navigator';
import _ from 'lodash';
import { findNestedObject } from '@/utils/objects';

/**
 * Retrieves filtered data for the details component based on the selected item ID.
 *
 * @param {Array} data - An array of data representing the tree.
 * @param {string} selectedItemUUID - Node ID of the selected element in the tree.
 * @returns {Object} - An object containing the selected component, model, and relationship, with type and data properties.
 */
export const getFilteredDataForDetailsComponent = (data, selectedItemUUID) => {
  const selectedIdArr = selectedItemUUID.split('.');
  const resultObject = findNestedObject(
    data,
    (obj) => _.get(obj, 'id') === selectedIdArr[selectedIdArr.length - 1],
  );

  const propertiesArr = resultObject ? Object.keys(resultObject).map(_.toLower) : [];
  const isPropertyIncluded = (property) => propertiesArr.includes(_.toLower(property));

  const isDepthCheck = () => {
    if (isPropertyIncluded('summary')) {
      return REGISTRANTS;
    } else if (isPropertyIncluded(COMPONENTS) || isPropertyIncluded(RELATIONSHIPS)) {
      return MODELS;
    } else if (isPropertyIncluded('evaluationQuery')) {
      return RELATIONSHIPS;
    } else if (isPropertyIncluded('schema')) {
      return COMPONENTS;
    }
  };

  const selectedType = isDepthCheck();

  return {
    type: selectedType || '',
    data: resultObject || {},
  };
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
