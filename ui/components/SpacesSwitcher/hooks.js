//@ts-check
import { APP_MODE, RESOURCE_TYPE } from '@/utils/Enum';
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
