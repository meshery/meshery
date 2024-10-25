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
    } else if (isPropertyIncluded('evaluationQuery') || isPropertyIncluded('selector')) {
      return RELATIONSHIPS;
    } else if (isPropertyIncluded('component')) {
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
      const uniqueVersions = _.groupBy(models, (modelDef) => modelDef?.model?.version);
      const arrayOfUniqueVersions = Object.values(uniqueVersions);
      const existingModelDef = acc.find((m) => m.name === name);

      const mergedData = arrayOfUniqueVersions.map((modelsWithSameVersion) => {
        let subVal = {
          relationships: {},
          components: {},
        };
        modelsWithSameVersion.map((model) => {
          subVal.relationships = groupRelationshipsByKind(
            _.union(subVal.relationships, model.relationships),
          );
          subVal.components = _.union(subVal.components, model.components);
        });
        return {
          ...modelsWithSameVersion[0],
          ...subVal,
        };
      });

      if (existingModelDef) {
        existingModelDef.model.version = _.union(
          existingModelDef.model.version,
          mergedData.map((model) => model?.model?.version),
        );
        existingModelDef.versionBasedData = existingModelDef.versionBasedData.concat(mergedData);
      } else {
        const selectedModelDef = models[0];
        acc.push({
          ...selectedModelDef,
          version: mergedData.map((model) => model?.model?.version),
          versionBasedData: mergedData,
        });
      }

      return acc;
    },
    [],
  );

  return result;
};

/**
 * Function takes theme type and returns the theme object
 * @param {object} - theme type
 */
export const reactJsonTheme = (themeType) => ({
  base00: themeType === 'dark' ? '#303030' : '#ffffff',
  base01: '#444c56',
  base02: themeType === 'dark' ? '#586069' : '#abb2bf',
  base03: '#6a737d',
  base04: '#477E96',
  base05: '#9ea7a6',
  base06: '#d8dee9',
  base07: themeType === 'dark' ? '#FFF3C5' : '#002B36',
  base08: '#2a5491',
  base09: '#d19a66',
  base0A: '#EBC017',
  base0B: '#237986',
  base0C: '#56b6c2',
  base0D: '#B1B6B8',
  base0E: '#e1e6cf',
  base0F: '#647881',
});
