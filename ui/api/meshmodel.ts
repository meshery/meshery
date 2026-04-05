import { promisifiedDataFetch } from '../lib/data-fetch';
import { MESHMODEL_COMPONENT_ENDPOINT, MESHMODEL_ENDPOINT } from '../constants/endpoints';

const COMPONENTS_ENDPOINT = '/api/meshmodels/components';
const CATEGORIES_ENDPOINT = '/api/meshmodels/categories';

/**
 * @typedef {{
 * pageSize: (Number|"all");
 * page: Number;
 * trim: Boolean
 * }} pageOptions
 */

/** @type {pageOptions} */
const defaultOptions = {
  pageSize: 'all',
  page: 0,
  trim: true,
};

// TODO: Migrate consumers to RTK Query hooks from rtk-query/meshModel.ts
// Each function below has an RTK Query equivalent. Consumer migration tracked separately.

export async function getMeshModels(page = 1, pageSize = 'all', options = defaultOptions) {
  return await promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}?page=${page}&pagesize=${pageSize}&${optionToQueryConvertor({
      ...defaultOptions,
      ...options,
    })}`,
  );
}

export async function getComponentFromModelApi(model, pageSize = 'all', trim = true) {
  return await promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}/${model}/components?pagesize=${pageSize}&trim=${trim}`,
  );
}

export async function getDuplicateModels(model, version) {
  return await promisifiedDataFetch(`${MESHMODEL_ENDPOINT}/${model}?version=${version}      `);
}

export async function getDuplicateComponents(componentKind, apiVersion, modelName) {
  return await promisifiedDataFetch(
    `${COMPONENTS_ENDPOINT}/${componentKind}?apiVersion=${apiVersion}&?model=${modelName}`,
  );
}

export async function getMeshModelRegistrants(page = 1, pageSize = 'all') {
  return await promisifiedDataFetch(
    `/api/meshmodels/registrants?page=${page}&pageSize=${pageSize}`,
  );
}

export async function getVersionedComponentFromModel(
  model,
  version,
  pageSize = 'all',
  trim = true,
) {
  return await promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}/${model}/components?version=${version}&pagesize=${pageSize}&trim=${trim}`,
  );
}

export async function getComponentsDetail(page) {
  return await promisifiedDataFetch(`api/meshmodels/components?page=${page}`);
}

export async function getRelationshipsDetail(page) {
  return await promisifiedDataFetch(`api/meshmodels/relationships?page=${page}`);
}

export async function getMeshModelComponent(model, component, version, apiVersion) {
  const versionQueryString = !version ? '' : `?version=${version}`;
  const apiVersionQueryString = !apiVersion
    ? ''
    : !version
      ? `?apiVersion=${apiVersion}`
      : `&apiVersion=${apiVersion}`;

  return promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}/${model}/components/${component}${versionQueryString}${apiVersionQueryString}`,
  );
}

export async function getMeshModelComponentByName(component) {
  return promisifiedDataFetch(`${MESHMODEL_COMPONENT_ENDPOINT}/components/${component}`);
}

export async function fetchCategories() {
  return promisifiedDataFetch(`${CATEGORIES_ENDPOINT}`);
}

export async function getModelFromCategoryApi(category) {
  return promisifiedDataFetch(`${CATEGORIES_ENDPOINT}/${category}/models?pagesize=all`);
}

export async function getModelByName(modelName, options = defaultOptions) {
  return promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}/${modelName}?${optionToQueryConvertor(options)}`,
  );
}

/**
 *
 * @param {pageOptions} options
 */
function optionToQueryConvertor(options) {
  const uri = new URLSearchParams();
  const { pageSize, page, trim, components, relationships } = options;

  if (trim) {
    uri.append('trim', `${trim}`);
  }

  if (pageSize) {
    uri.append('pagesize', `${pageSize}`);
  }

  if (page) {
    uri.append('page', `${page}`);
  }

  if (components) {
    uri.append('components', `${components}`);
  }
  if (relationships) {
    uri.append('relationships', `${relationships}`);
  }
  return uri.toString();
}
