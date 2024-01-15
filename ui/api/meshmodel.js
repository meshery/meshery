import { promisifiedDataFetch } from '../lib/data-fetch';
import {
  MESHMODEL_COMPONENT_ENDPOINT,
  MESHMODEL_ENDPOINT,
  MESHMODEL_RELATIONSHIPS_ENDPOINT,
  SORT,
} from '../constants/endpoints';

const COMPONENTS_ENDPOINT = '/api/meshmodels/components';
const CATEGORIES_ENDPOINT = '/api/meshmodels/categories';

/**
 * @typedef {{
 * paginated: Boolean;
 * pageSize: (Number|"all");
 * page: Number;
 * trim: Boolean
 * }} pageOptions
 */

/** @type {pageOptions} */
const defaultOptions = {
  paginated: false,
  pageSize: 'all',
  page: 0,
  trim: true,
};

/**
 * Fetches the Relationships from the server
 *
 * @returns
 */
export async function fetchRelationships() {
  return await promisifiedDataFetch(`${MESHMODEL_RELATIONSHIPS_ENDPOINT}`);
}

export async function getAllComponents(page = 1, pageSize = 'all') {
  return await promisifiedDataFetch(`${COMPONENTS_ENDPOINT}?page=${page}&pagesize=${pageSize}`);
}

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

export async function getMeshModelsByRegistrants(page = 1, pageSize = 'all', registrant) {
  return await promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}?page=${page}&pagesize=${pageSize}&registrant=${registrant}`,
  );
}

export async function getRelationshipFromModelApi(model, pageSize = 'all', trim = true) {
  return await promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}/${model}/relationships?pagesize=${pageSize}&trim=${trim}`,
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

export async function getComponentsDetailWithPageSize(
  page = 1,
  pageSize = 'all',
  sort = SORT.ASCENDING,
  order = '',
) {
  return await promisifiedDataFetch(
    `api/meshmodels/components?page=${page}&pagesize=${pageSize}&order=${encodeURIComponent(
      order,
    )}&sort=${sort}`,
  );
}

export async function getComponentsDetail(page) {
  return await promisifiedDataFetch(`api/meshmodels/components?page=${page}`);
}

export async function getModelsDetail(page) {
  return await promisifiedDataFetch(`${MESHMODEL_ENDPOINT}?page=${page}`);
}

export async function getRelationshipsDetailWithPageSize(
  page = 1,
  pageSize = 'all',
  sort = SORT.ASCENDING,
  order = '',
) {
  return await promisifiedDataFetch(
    `api/meshmodels/relationships?page=${page}&pagesize=${pageSize}&sort=${sort}&order=${encodeURIComponent(
      order,
    )}`,
  );
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
// export async function queryMeshModel(modelQueryString, paginated = true) {
//   // Note: returns paginated response
//   if (paginated) {
//     return promisifiedDataFetch(`${MESHMODEL_ENDPOINT}/${modelQueryString}?search=true&trim=true`);
//   }

//   // to get full response
//   return promisifiedDataFetch(
//     `${MESHMODEL_ENDPOINT}/${modelQueryString}?search=true&page=1&pagesize=all&trim=true`
//   );
// }

export async function fetchCategories() {
  return promisifiedDataFetch(`${CATEGORIES_ENDPOINT}`);
}

export async function getModelFromCategoryApi(category) {
  return promisifiedDataFetch(`${CATEGORIES_ENDPOINT}/${category}/models?pagesize=all`);
}

/**
 *
 * @param {string} queryString
 * @param {pageOptions} options
 */
export async function searchModels(queryString, options = defaultOptions) {
  return promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}?search=${encodeURI(queryString)}&${optionToQueryConvertor({
      ...defaultOptions,
      ...options,
    })}`,
  );
}

export async function searchComponents(queryString, options = defaultOptions) {
  return promisifiedDataFetch(
    `/api/meshmodels/components?search=${encodeURI(queryString)}&${optionToQueryConvertor(
      options,
    )}`,
  );
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
