import _ from "lodash";

export const WILDCARD_V = "*(All Versions)"

import {
  fetchCategories,
  getComponentFromModelApi,
  getModelFromCategoryApi,
  getVersionedComponentFromModel
} from "../../api/meshmodel";
import { compose } from "lodash/fp";
import { useEffect, useState } from "react";
import getMostRecentVersion, { versionSortComparatorFn, sortAndGroupVersionsInModel } from "../versionSort";

const handleError = e => {
  console.error("MeshModel axios error ocurred", e);
};

function componentToLatestApiVersion(components) {
  const componentToAPiVersionMap = {}; // this is for storing all the apiVersions of similar components in order to get the most recent at the end

  [...components].forEach(component => {
    if (componentToAPiVersionMap?.[component.kind]) {
      componentToAPiVersionMap[component.kind] = [
        ...componentToAPiVersionMap[component.kind],
        component.apiVersion
      ];
    } else {
      componentToAPiVersionMap[component.kind] = [component.apiVersion];
    }
  });

  Object.keys(componentToAPiVersionMap).forEach(key => {
    componentToAPiVersionMap[key] = getMostRecentVersion(componentToAPiVersionMap[key]);
  });

  return componentToAPiVersionMap;
}

function removeDuplicateMeshModelComponents(components) {
  const componentClone = [...components];
  const cmpToApiVersion = componentToLatestApiVersion(componentClone);

  // component kind set keeps track of redudan components
  const componentKindUniqueSet = new Set();

  return componentClone
    .filter(({ kind }) => {
      // filter unique components
      // already found in the unique set, means that it is already filtered
      if (componentKindUniqueSet.has(kind)) {
        return false;
      }

      componentKindUniqueSet.add(kind);
      return true;
    })
    .map(component => ({
      // on all unique components, set the apiVersion to latest one
      ...component,
      apiVersion : cmpToApiVersion[component.kind] || component.apiVersion // fallback in case of mishap
    }));
}

function sortMeshModelComponents(components) {
  return [...components].sort((a, b) => a.kind.localeCompare(b.kind));
}

export const removeDuplicatesAndSortByAlphabet = _.flowRight(
  sortMeshModelComponents,
  removeDuplicateMeshModelComponents
);

// processing includes sorting and deduplicating components
function getProcessedMeshModelResponseData(meshModelResponse) {
  return [...meshModelResponse]
    .map(meshmodel => ({
      ...meshmodel,
      components : removeDuplicatesAndSortByAlphabet(meshmodel.components)
    }))
    .sort((modelA, modelB) => versionSortComparatorFn(modelA.version, modelB.version))
    .reverse(); // sort the versions in reverse order
}

function deduplicatedListOfComponentsFromAllVersions(components) {
  const uniqueComponents = [...new Set(components.map(comp => comp.kind))];
  return uniqueComponents.map(compKind => components.find(c => c.kind === compKind));
}

function groupComponentsByVersion(components) {
  const versions = [...new Set(components.map(component => component.model.version))];

  if (versions.length > 1) {
    return [
      {
        version : WILDCARD_V,
        components : deduplicatedListOfComponentsFromAllVersions(components)
      },
      ...versions.map(version => ({
        version : version,
        components : components.filter(component => component.model.version === version)
      }))
    ]
  }

  // don't attach the wildcards
  return versions.map(version => ({
    version : version,
    components : components.filter(component => component.model.version === version)
  }));
}

const getProcessedComponentsData = compose(
  getProcessedMeshModelResponseData,
  groupComponentsByVersion
);

function convertToArray(item) {
  if (Array.isArray(item)) return item;

  return [item];
}

export function useMeshModelComponents() {
  const [meshmodelComponents, setMeshModelComponents] = useState({});
  const [categories, setCategories] = useState([]);
  const [models, setModels] = useState({});

  useEffect(() => {
    fetchCategories()
      .then((categoryJson) => {
        setCategories(categoryJson.categories.sort((catA, catB) => catA.name.localeCompare(catB.name)));
      })
      .catch(handleError);
  }, []);

  async function getModelFromCategory(category) {
    // already fetched the models from catgory and stored
    if (models[category]) {
      return;
    }

    getModelFromCategoryApi(category)
      .then((response) => {
        setModels(
          Object.assign(
            { ...models },
            {
              [category] : sortAndGroupVersionsInModel(response.models)
            }
          )
        );
      })
      .catch(handleError);
  }

  async function getComponentsFromModel(modelName, version) {
    if (!version) {
      if (!meshmodelComponents[modelName]) {
        const modelData = (await getComponentFromModelApi(modelName));
        console.log("modeled data....", modelData)
        setMeshModelComponents(
          Object.assign(
            { ...meshmodelComponents },
            {
              [modelName] : getProcessedComponentsData(modelData.components)
            }
          )
        );
      }
      return;
    }

    if (
      !meshmodelComponents[modelName] ||
      !convertToArray(meshmodelComponents[modelName])?.find(
        model => model.version === version
      )
    ) {
      const modelData = (await getVersionedComponentFromModel(modelName, version));
      setMeshModelComponents(
        Object.assign(
          { ...meshmodelComponents },
          {
            [modelName] : getProcessedComponentsData(modelData.components)
          }
        )
      );
    }
  }

  return {
    models,
    meshmodelComponents,
    getModelFromCategory,
    getComponentsFromModel,
    categories
  };
}
