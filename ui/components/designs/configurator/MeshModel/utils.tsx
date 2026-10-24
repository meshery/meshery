import { JsonParse } from '../../../../utils/utils';

export function getComponentMetadata(metadata) {
  const cleanedMeta = { ...metadata };

  // for the older designs where metadata is stored in css styles, replace the svgWhite with it
  if (cleanedMeta['background-image']) {
    cleanedMeta.svgWhite = cleanedMeta['background-image'];
  }

  return Object.assign(
    {
      // default assignments
      shape: 'circle', // shape should be small case
      primaryColor: '#00B39F',
      secondaryColor: '#b2e8e2',
      svgColor: 'static/img/component-svg/meshery.svg',
      svgWhite: 'static/img/component-svg/meshery.svg',
    },
    cleanedMeta,
  );
}

/**
 * Builds the React key for a component form in the Design Configurator list.
 * The model name is part of the key so that switching models never reuses a form
 * (and the component id it holds) that was opened for another model's component.
 * @param {string} modelName name of the model whose components are listed
 * @param {Types.ComponentDefinition} componentDefinition component shown by the form
 * @param {number} index position of the component in the list
 * @returns {string}
 */
export function getComponentFormKey(modelName, componentDefinition, index) {
  return `${modelName}-${componentDefinition?.component?.kind}-${index}`;
}

export function getStyleOverrides(stringifiedOverrides) {
  if (stringifiedOverrides) {
    return JsonParse(stringifiedOverrides);
  }

  return {};
}
