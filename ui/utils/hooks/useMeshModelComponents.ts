import _ from 'lodash';

export const WILDCARD_V = 'All Versions';

import {
  useGetModelCategoriesQuery,
  useLazyGetModelFromCategoryQuery,
  useLazyGetComponentsFromModalQuery,
} from '../../rtk-query/meshModel';
import { compose } from 'lodash/fp';
import { useEffect, useState } from 'react';
import getMostRecentVersion, {
  versionSortComparatorFn,
  sortAndGroupVersionsInModel,
} from '../versionSort';

const handleError = (e) => {
  console.error('MeshModel error occurred', e);
};

function componentToLatestApiVersion(components) {
  const componentToAPiVersionMap = {}; // this is for storing all the apiVersions of similar components in order to get the most recent at the end

  [...components].forEach((componentDef) => {
    if (componentToAPiVersionMap?.[componentDef.component.kind]) {
      componentToAPiVersionMap[componentDef.component.kind] = [
        ...componentToAPiVersionMap[componentDef.component.kind],
        componentDef.component.version,
      ];
    } else {
      componentToAPiVersionMap[componentDef.component.kind] = [componentDef.component.version];
    }
  });

  Object.keys(componentToAPiVersionMap).forEach((key) => {
    componentToAPiVersionMap[key] = getMostRecentVersion(componentToAPiVersionMap[key]);
  });

  return componentToAPiVersionMap;
}

function removeDuplicateMeshModelComponents(componentDefs) {
  const componentClone = [...componentDefs];
  const cmpToApiVersion = componentToLatestApiVersion(componentClone);
  // component kind set keeps track of redudant components
  const componentKindUniqueSet = new Set();
  return componentClone
    .filter((componentDef) => {
      const kind = componentDef.component;
      // filter unique components
      // already found in the unique set, means that it is already filtered
      if (componentKindUniqueSet.has(kind)) {
        return false;
      }

      componentKindUniqueSet.add(kind);
      return true;
    })
    .map((componentDef) => ({
      // on all unique components, set the apiVersion to latest one
      ...componentDef,
      apiVersion: cmpToApiVersion[componentDef.component.kind] || componentDef.component.version, // fallback in case of mishap
    }));
}

function sortMeshModelComponents(componentDefs) {
  return [...componentDefs].sort((a, b) => a.component.kind.localeCompare(b.component.kind));
}

export const removeDuplicatesAndSortByAlphabet = _.flowRight(
  sortMeshModelComponents,
  removeDuplicateMeshModelComponents,
);

// processing includes sorting and deduplicating components
function getProcessedMeshModelResponseData(meshModelResponse) {
  return [...meshModelResponse]
    .map((meshmodel) => ({
      ...meshmodel,
      components: removeDuplicatesAndSortByAlphabet(meshmodel.components),
    }))
    .sort((modelA, modelB) => versionSortComparatorFn(modelA.version, modelB.version))
    .reverse(); // sort the versions in reverse order
}

function deduplicatedListOfComponentsFromAllVersions(componentDefs) {
  const uniqueComponents = [
    ...new Set(componentDefs.map((componentDef) => componentDef.component.kind)),
  ];
  return uniqueComponents.map((compKind) =>
    componentDefs.find((componentDef) => componentDef.component.kind === compKind),
  );
}

function groupComponentsByVersion(componentDefs) {
  const versions = [
    ...new Set(componentDefs?.map((componentDef) => componentDef.model.model.version) || []),
  ];

  if (versions.length > 1) {
    return [
      {
        version: WILDCARD_V,
        components: deduplicatedListOfComponentsFromAllVersions(componentDefs),
      },
      ...versions.map((version) => ({
        version: version,
        components: componentDefs.filter(
          (componentDef) => componentDef.model.model.version === version,
        ),
      })),
    ];
  }

  // don't attach the wildcards
  return versions.map((version) => ({
    version: version,
    components: componentDefs.filter(
      (componentDef) => componentDef.model.model.version === version,
    ),
  }));
}
const getProcessedComponentsData = compose(
  getProcessedMeshModelResponseData,
  groupComponentsByVersion,
);

function convertToArray(item) {
  if (Array.isArray(item)) return item;

  return [item];
}

export function useMeshModelComponents() {
  const { data: categoriesData } = useGetModelCategoriesQuery();
  const [triggerGetModelFromCategory] = useLazyGetModelFromCategoryQuery();
  const [triggerGetComponentsFromModal] = useLazyGetComponentsFromModalQuery();

  const [meshmodelComponents, setMeshModelComponents] = useState({});
  const [categories, setCategories] = useState([]);
  const [models, setModels] = useState({});

  useEffect(() => {
    if (categoriesData?.categories) {
      setCategories(
        [...categoriesData.categories].sort((catA, catB) => catA.name.localeCompare(catB.name)),
      );
    }
  }, [categoriesData]);

  async function getModelFromCategory(category) {
    // already fetched the models from catgory and stored
    if (models[category]) {
      return;
    }

    triggerGetModelFromCategory({ category, params: { pagesize: 'all' } })
      .unwrap()
      .then((response) => {
        if (response?.models) {
          setModels((prevModels) => ({
            ...prevModels,
            [category]: sortAndGroupVersionsInModel(response.models),
          }));
        }
      })
      .catch(handleError);
  }

  async function getComponentsFromModel(modelName, version) {
    if (!version) {
      if (!meshmodelComponents[modelName]) {
        try {
          const response = await triggerGetComponentsFromModal({
            model: modelName,
            params: { pagesize: 'all', trim: true },
          }).unwrap();

          if (response?.components) {
            setMeshModelComponents((prev) => ({
              ...prev,
              [modelName]: getProcessedComponentsData(response.components),
            }));
          }
        } catch (e) {
          handleError(e);
        }
      }
      return;
    }

    if (
      !meshmodelComponents[modelName] ||
      !convertToArray(meshmodelComponents[modelName])?.find(
        (item) => item.version === version,
      )
    ) {
      try {
        const response = await triggerGetComponentsFromModal({
          model: modelName,
          params: { version, pagesize: 'all', trim: true },
        }).unwrap();

        if (response?.components) {
          setMeshModelComponents((prev) => ({
            ...prev,
            [modelName]: getProcessedComponentsData(response.components),
          }));
        }
      } catch (e) {
        handleError(e);
      }
    }
  }

  return {
    models,
    meshmodelComponents,
    getModelFromCategory,
    getComponentsFromModel,
    categories,
  };
}
