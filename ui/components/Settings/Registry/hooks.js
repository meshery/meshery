import { useEffect } from 'react';
import { OVERVIEW, MODELS, GRAFANA, PROMETHEUS, REGISTRY } from '../../../constants/navigator';
import { useRouter } from 'next/router';
import { useRef } from 'react';

export const useMeshModelComponentRouter = () => {
  const router = useRouter();
  const { query } = router;

  if (query.settingsCategory === REGISTRY && !query.tab) {
    router.push({
      pathname: router.pathname,
      query: {
        ...query,
        tab: MODELS,
      },
    });
  }
  const searchQuery = query.searchText || null;
  const selectedTab = query.tab === GRAFANA || query.tab === PROMETHEUS ? OVERVIEW : query.tab;
  const selectedPageSize = query.pagesize || 25;

  return { searchQuery, selectedTab, selectedPageSize };
};

export const useInfiniteScrollRef = (callback) => {
  const observerRef = useRef(null);
  const triggerRef = useRef(null);

  useEffect(() => {
    // setTimeout gives the browser time to finish rendering the DOM elements before executing the callback function.
    const timeoutId = setTimeout(() => {
      if (!triggerRef.current) {
        return () => observerRef.current && observerRef.current.disconnect();
      }

      observerRef.current = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              callback();
            }
          });
        },
        { threshold: 0.01 },
      );
      observerRef.current.observe(triggerRef.current);
    }, 0);

    return () => {
      observerRef.current && observerRef.current.disconnect();
      clearTimeout(timeoutId);
    };
  }, [callback, triggerRef.current]);

  return triggerRef;
};

export const useRegistryRouter = () => {
  const router = useRouter();
  const { query, push: pushRoute, route } = router;

  const settingsCategory = query.settingsCategory;
  const tab = query.tab;
  const selectedItemUUID = query.selectedItemUUID || '';
  const searchText = query.searchText || null;
  let filters = {
    searchText: searchText,
  };

  const handleUpdateSelectedRoute = (nodeIds, filters) => {
    const id = nodeIds[0];
    const validSettingsCategory = settingsCategory || 'Registry';
    const validTab = tab || 'Models';

    const queryString = new URLSearchParams({
      settingsCategory: validSettingsCategory,
      tab: validTab,
      selectedItemUUID: id,
      ...filters,
    }).toString();

    let targetRoute = route;
    if (typeof window !== 'undefined') {
      targetRoute = window.location.pathname;
    }

    pushRoute(`${targetRoute}?${queryString}`, undefined, { shallow: true });
  };

  return {
    settingsCategory,
    tab,
    handleUpdateSelectedRoute,
    selectedItemUUID,
    filters,
  };
};
