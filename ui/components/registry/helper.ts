import { MODELS, REGISTRANTS, COMPONENTS, RELATIONSHIPS } from '@/constants/navigator';
import _ from 'lodash';
import { findNestedObject } from '@/utils/objects';
import { alpha, type Theme } from '@/theme';

/**
 * Retrieves filtered data for the details component based on the selected item ID.
 *
 * @param {Array} data - An array of data representing the tree.
 * @param {string} selectedItemUUID - Node ID of the selected element in the tree.
 * @returns {Object} - An object containing the selected component, model, and relationship, with type and data properties.
 */
export const getFilteredDataForDetailsComponent = (data: Array<any>, selectedItemUUID: string) => {
  const selectedIdArr = selectedItemUUID.split('.');
  const resultObject = findNestedObject(
    data,
    (obj) => _.get(obj, 'id') === selectedIdArr[selectedIdArr.length - 1],
  );

  const propertiesArr = resultObject ? Object.keys(resultObject).map(_.toLower) : [];
  const isPropertyIncluded = (property: string) => propertiesArr.includes(_.toLower(property));

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
export const groupRelationshipsByKind = (relationships: Array<any>) => {
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
export const removeDuplicateVersions = (data: Array<any>) => {
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
 * To style JSON viewing in react-json-tree.
 * Refer to base-16 theme styling guidelines for more info.
 * https://github.com/chriskempson/base16/blob/main/styling.md
 *
 * Accepts a Sistent {@link Theme} so the base16 ramp can be derived from
 * `theme.palette.*` tokens (no inline hex literals). Callers in components
 * pass the result of `useTheme()` directly.
 *
 * @property {string} base00 - BACKGROUND_COLOR
 * @property {string} base02 - OBJECT_OUTLINE_COLOR
 * @property {string} base04 - OBJECT_DETAILS_COLOR
 * @property {string} base07 - OBJECT_KEY_COLOR
 * @property {string} base09 - ITEM_STRING_COLOR, DATE_COLOR, STRING_COLOR
 * @property {string} base0A - SYMBOL_COLOR, FUNCTION_COLOR, UNDEFINED_COLOR, NULL_COLOR
 * @property {string} base0D - ITEM_STRING_EXPANDED_COLOR, ARROW_COLOR
 * @property {string} base0E - BOOLEAN_COLOR, NUMBER_COLOR
 */
export const reactJsonTheme = (theme: Theme) => {
  const isDark = theme.palette.mode === 'dark';
  // Object-outline contrast needs more punch in dark mode where the
  // background is already dim; in light mode a softer divider tone reads
  // better against the white canvas.
  const outline = isDark
    ? alpha(theme.palette.text.primary, 0.32)
    : alpha(theme.palette.text.primary, 0.16);
  return {
    base00: theme.palette.background.default,
    base01: theme.palette.divider,
    base02: outline,
    base03: theme.palette.text.secondary,
    base04: theme.palette.info.main,
    base05: theme.palette.text.secondary,
    base06: alpha(theme.palette.text.primary, 0.24),
    base07: theme.palette.text.primary,
    base08: theme.palette.error.main,
    base09: theme.palette.warning.main,
    base0A: theme.palette.warning.light,
    base0B: theme.palette.success.main,
    base0C: theme.palette.info.light,
    base0D: theme.palette.info.main,
    base0E: theme.palette.primary.main,
    base0F: theme.palette.text.secondary,
  };
};
