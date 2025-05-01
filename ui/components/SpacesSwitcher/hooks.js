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
