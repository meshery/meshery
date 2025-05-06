import { APP_MODE, RESOURCE_TYPE } from '@/utils/Enum';
import { isInOperatorMode, JsonParse } from '@/utils/utils';
import _ from 'lodash';
import { useEffect, useRef } from 'react';

const useInfiniteScroll = ({ isLoading, hasMore, onLoadMore }) => {
  const loadingRef = useRef(null);
  const observerRef = useRef(null);

  useEffect(() => {
    if (isLoading || !hasMore) return;

    if (observerRef.current && loadingRef.current) {
      observerRef.current.unobserve(loadingRef.current);
    }

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !isLoading) {
          onLoadMore();
        }
      },
      { threshold: 0.1 },
    );

    observerRef.current = observer;

    if (loadingRef.current) {
      observer.observe(loadingRef.current);
    }

    return () => {
      observer.disconnect();
    };
  }, [isLoading, hasMore, onLoadMore]);

  return { loadingRef };
};

export default useInfiniteScroll;

export const getDesignPath = (id) => {
  const currentRoute = new URL(window.location.href);
  const currentURI = currentRoute.origin + currentRoute.pathname;

  const newParams = new URLSearchParams({
    mode: APP_MODE.DESIGN,
    ...(id ? { design: id } : {}),
  });
  const newURI = currentURI + '?' + newParams.toString();
  return newURI;
};
export const viewPath = ({ id, name }) => {
  const currentRoute = new URL(window.location.href);
  const currentURI = currentRoute.origin + currentRoute.pathname;
  const newParams = new URLSearchParams({
    mode: APP_MODE.OPERATOR,
    type: RESOURCE_TYPE.VIEW,
    ...(id ? { id } : {}),
    ...(name ? { name } : {}),
  });
  const newURI = currentURI + '?' + newParams.toString();
  return newURI;
};

export const catalogPath = ({ id, name }) => {
  const currentRoute = new URL(window.location.href);
  const currentURI = currentRoute.origin + currentRoute.pathname;
  const newParams = new URLSearchParams({
    mode: APP_MODE.DESIGN,
    type: RESOURCE_TYPE.CATALOG,
    ...(id ? { id } : {}),
    ...(name ? { name } : {}),
  });
  const newURI = currentURI + '?' + newParams.toString();
  return newURI;
};

export const getShareableResourceRoute = (type, id, name) => {
  if (type === RESOURCE_TYPE.DESIGN) {
    return getDesignPath(id);
  }

  if (type === RESOURCE_TYPE.VIEW) {
    return viewPath({ id, name });
  }

  if (type === RESOURCE_TYPE.CATALOG) {
    return catalogPath({ id, name });
  }

  throw new Error(`Unknown resource type ${type}`);
};

/**
 * Get model names based on their display names
 * @param {object} - Models data
 * @param {array} - Array of model display names
 * @return {array} - Array of unique model names
 */
export const getModelNamesBasedOnDisplayNames = (meshModels, displayNames) => {
  const compatibilityStore = _.uniqBy(meshModels, (model) => _.toLower(model.displayName))
    ?.filter((model) =>
      displayNames.some((comp) => _.toLower(comp) === _.toLower(model.displayName)),
    )
    ?.map((model) => model.name);
  return compatibilityStore;
};

export const handleUpdatePatternVisibility = async ({
  value,
  updatePatterns,
  selectedResource,
}) => {
  const res = await updatePatterns({
    updateBody: {
      id: selectedResource?.id,
      name: selectedResource.name,
      catalog_data: selectedResource.catalog_data,
      design_file: JsonParse(selectedResource.pattern_file),
      visibility: value,
    },
  });
  return {
    error: res?.error?.error,
  };
};

export const handleUpdateViewVisibility = async ({ value, updateView, selectedResource }) => {
  const res = await updateView({
    id: selectedResource?.id,
    body: {
      visibility: value,
    },
  });
  return {
    error: res.error?.error,
  };
};

export const getDefaultFilterType = () => {
  if (isInOperatorMode()) {
    return RESOURCE_TYPE.VIEW;
  }
  return RESOURCE_TYPE.DESIGN;
};
