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

/**
 * Group relationships by kind
 * @param {object} - Relationships arrays
 */
export const groupRelationshipsByKind = (relationships) => {
  const groupedRelationships = {};

  relationships.forEach((relationship) => {
    const { id, kind } = relationship;

    if (!groupedRelationships[kind]) {
      groupedRelationships[kind] = { kind, relationships: [] };
    }

    groupedRelationships[kind].relationships.push({ id, ...relationship });
  });
  const resultArray = Object.values(groupedRelationships);
  return resultArray;
};

/**
 * Function takes models data and merges the duplicate data
 * @param {object} - models data
 */
export const removeDuplicateVersions = (data) => {
  const groupedModels = _.groupBy(data, 'name');

  const result = _.reduce(
    groupedModels,
    (acc, models, name) => {
      const uniqueVersions = _.groupBy(models, 'version');
      const arrayOfUniqueVersions = Object.values(uniqueVersions);

      const existingModel = acc.find((m) => m.name === name);

      const mergedData = arrayOfUniqueVersions.map((modelsWithSameVersion) => {
        let subVal = {
          relationships: {},
          components: {},
        };
        modelsWithSameVersion.map((model) => {
          subVal.relationships = _.union(subVal.relationships, model.relationships);
          subVal.components = _.union(subVal.components, model.components);
        });
        return {
          ...modelsWithSameVersion[0],
          ...subVal,
        };
      });

      if (existingModel) {
        existingModel.version = _.union(
          existingModel.version,
          mergedData.map((model) => model.version),
        );
        existingModel.versionBasedData = existingModel.versionBasedData.concat(mergedData);
      } else {
        const selectedModel = models[0];
        acc.push({
          ...selectedModel,
          version: mergedData.map((model) => model.version),
          versionBasedData: mergedData,
        });
      }

      return acc;
    },
    [],
  );

  return result;
};
