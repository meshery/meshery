import { promisifiedDataFetch } from "../lib/data-fetch";
import { MESHMODEL_ENDPOINT, MESHMODEL_RELATIONSHIPS_ENDPOINT } from "../constants/endpoints";

const COMPONENTS_ENDPOINT = "/components";
const CATEGORIES_ENDPOINT = "/api/meshmodels/categories";

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
  paginated : false,
  pageSize : "all",
  page : 0,
  trim : true
};

/**
 * Fetches the Relationships from the server
 *
 * @returns
 */
export async function fetchRelationships() {
  return await promisifiedDataFetch(`${MESHMODEL_RELATIONSHIPS_ENDPOINT}`)
}

export async function getAllComponents() {
  return await promisifiedDataFetch(MESHMODEL_ENDPOINT + COMPONENTS_ENDPOINT);
}

export async function getMeshModels(page = 1, pageSize = "all") {
  return await promisifiedDataFetch(`${MESHMODEL_ENDPOINT}?page=${page}&pagesize=${pageSize}`);
}

export async function getComponentFromModelApi(model, pageSize = "all", trim = true) {
  return await promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}/${model}/components?pagesize=${pageSize}&trim=${trim}`
  );
}

export async function getVersionedComponentFromModel(
  model,
  version,
  pageSize = "all",
  trim = true
) {
  return await promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}/${model}/components?version=${version}&pagesize=${pageSize}&trim=${trim}`
  );
}

export async function getComponentsDetail(page) {
  return await promisifiedDataFetch(
    `api/meshmodels/components?page=${page}`
  );
}

export async function getModelsDetail(page) {
  return await promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}?page=${page}`
  );
}

export async function getRelationshipsDetail(page) {
  return await promisifiedDataFetch(
    `api/meshmodels/relationships?page=${page}`
  );
}

export async function getMeshModelComponent(model, component, version, apiVersion) {
  const versionQueryString = !version ? "" : `?version=${version}`;
  const apiVersionQueryString = !apiVersion
    ? ""
    : !version
      ? `?apiVersion=${apiVersion}`
      : `&apiVersion=${apiVersion}`;

  return promisifiedDataFetch(
    `${MESHMODEL_ENDPOINT}/${model}/components/${component}${versionQueryString}${apiVersionQueryString}`
  );
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
    `${MESHMODEL_ENDPOINT}?search=${encodeURI(queryString)}&${optionToQueryConvertor(
      options
    )}`
  );
}

export async function searchComponents(queryString, options = defaultOptions) {
  return promisifiedDataFetch(
    `/api/meshmodels/components?search=${encodeURI(queryString)}&${optionToQueryConvertor(
      options
    )}`
  );
}
/**
 *
 * @param {pageOptions} options
 */
function optionToQueryConvertor(options) {
  const { paginated, pageSize, page, trim } = options;

  if (paginated) {
    return `trim=${trim}&${pageSize && `pagesize=${pageSize}`}&page=${page || 0}`;
  }

  return `pagesize=all&trim=${trim}`;
}
