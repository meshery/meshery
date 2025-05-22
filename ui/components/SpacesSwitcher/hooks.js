import { APP_MODE, RESOURCE_TYPE } from '@/utils/Enum';
import { isInOperatorMode, JsonParse } from '@/utils/utils';
import _ from 'lodash';
import { useEffect, useRef } from 'react';
import { useDeletePatternFileMutation } from '@/rtk-query/design';
import { useDeleteViewMutation } from '@/rtk-query/view';
import { useNotification } from '@/utils/hooks/useNotification';
import { DesignIcon, PROMPT_VARIANTS, useTheme, ViewIcon } from '@layer5/sistent';
import { EVENT_TYPES } from 'lib/event-types';
import { iconMedium } from 'css/icons.styles';
import { updateProgress } from '@/store/slices/mesheryUi';
import downloadContent, { downloadFileFromContent } from '@/utils/fileDownloader';

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

export const useGetIconBasedOnMode = ({ mode, designStyles, viewStyles }) => {
  const theme = useTheme();
  if (mode === RESOURCE_TYPE.DESIGN) {
    return <DesignIcon {...designStyles} />;
  } else if (mode == RESOURCE_TYPE.VIEW) {
    return <ViewIcon {...viewStyles} fill={theme.palette.icon.brand} {...iconMedium} />;
  }
};

export const useContentDelete = (modalRef) => {
  const [deleteView] = useDeleteViewMutation();
  const [deletePatternFile] = useDeletePatternFileMutation();
  const { notify } = useNotification();

  const handleDelete = async (items, type = RESOURCE_TYPE.DESIGN) => {
    const isDesign = type === RESOURCE_TYPE.DESIGN;
    const itemType = isDesign ? 'Design' : 'View';
    const deleteMutation = isDesign ? deletePatternFile : deleteView;

    const response = await modalRef.current.show({
      title: `Delete catalog item?`,
      subtitle: `Are you sure you want to delete ${
        items.length > 1
          ? `${items.length} ${itemType.toLowerCase()}s`
          : `the "${items[0].name}" ${itemType.toLowerCase()}`
      }?`,
      primaryOption: 'DELETE',
      variant: PROMPT_VARIANTS.DANGER,
    });

    if (response === 'DELETE') {
      items.forEach((item) => {
        const { name, id } = item;
        deleteMutation({ id })
          .unwrap()
          .then(() => {
            notify({
              message: `"${name}" ${itemType} deleted`,
              event_type: EVENT_TYPES.SUCCESS,
            });
          })
          .catch(() => {
            notify({
              message: `Unable to delete "${name}" ${itemType}`,
              event_type: EVENT_TYPES.ERROR,
            });
          });
      });
    }
  };

  return { handleDelete };
};

export const useContentDownload = () => {
  const { notify } = useNotification();
  const handleDesignDownload = (e, designs, source_type, params) => {
    e.stopPropagation();

    try {
      designs = Array.isArray(designs) ? designs : [designs];
      designs.forEach((design) => {
        updateProgress({ showProgress: true });
        let id = design.id;
        let name = design.name;
        downloadContent({ id, name, type: 'pattern', source_type, params });
        updateProgress({ showProgress: false });
        notify({ message: `"${name}" design downloaded`, event_type: EVENT_TYPES.INFO });
      });
    } catch (e) {
      console.error(e);
    }
  };

  const handleViewDownload = (views) => {
    try {
      views = Array.isArray(views) ? views : [views];
      views.forEach((view) => {
        updateProgress({ showProgress: true });
        let name = view.name;
        downloadFileFromContent(JSON.stringify(view), `${name}.json`, 'application/json');
        updateProgress({ showProgress: false });
        notify({ message: `"${name}" view downloaded`, event_type: EVENT_TYPES.INFO });
      });
    } catch (e) {
      console.error(e);
    }
  };

  return { handleDesignDownload, handleViewDownload };
};
