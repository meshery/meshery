import { JsonParse } from '../../../utils/utils';

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

export function getStyleOverrides(stringifiedOverrides) {
  if (stringifiedOverrides) {
    return JsonParse(stringifiedOverrides);
  }

  return {};
}
