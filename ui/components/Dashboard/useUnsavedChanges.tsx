import _ from 'lodash';
import { useRouter } from 'next/router';
import { useCallback, useEffect, useRef, useState } from 'react';

const useUnsavedChanges = ({ isEditMode, dashboardLayout, savedLayout }) => {
  const router = useRouter();
  const [showModal, setShowModal] = useState(false);
  const pendingUrlRef = useRef(null);
  const hasUnsavedChanges = isEditMode && !_.isEqual(dashboardLayout, savedLayout);
  const hasUnsavedRef = useRef(hasUnsavedChanges);
  hasUnsavedRef.current = hasUnsavedChanges;

  useEffect(() => {
    const handleRouteChange = (url) => {
      if (hasUnsavedRef.current) {
        pendingUrlRef.current = url;
        setShowModal(true);
        router.events.emit('routeChangeError');
        // eslint-disable-next-line no-throw-literal
        throw 'Navigation blocked due to unsaved changes';
      }
    };

    router.events.on('routeChangeStart', handleRouteChange);
    return () => {
      router.events.off('routeChangeStart', handleRouteChange);
    };
  }, [router.events]);

  useEffect(() => {
    const handleBeforeUnload = (e) => {
      if (hasUnsavedChanges) {
        e.preventDefault();
        e.returnValue = '';
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, [hasUnsavedChanges]);

  const confirmNavigation = useCallback(() => {
    const destination = pendingUrlRef.current;
    setShowModal(false);
    pendingUrlRef.current = null;
    if (destination) {
      router.push(destination);
    }
  }, [router]);

  const cancelNavigation = useCallback(() => {
    setShowModal(false);
    pendingUrlRef.current = null;
  }, []);

  return {
    showModal,
    hasUnsavedChanges,
    confirmNavigation,
    cancelNavigation,
  };
};

export default useUnsavedChanges;
