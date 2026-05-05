import _ from 'lodash';
import { useRouter } from 'next/router';
import { useCallback, useEffect, useRef, useState } from 'react';

interface UseUnsavedChangesProps {
  isEditMode: boolean;
  dashboardLayout: unknown;
  savedLayout: unknown;
}

const useUnsavedChanges = ({
  isEditMode,
  dashboardLayout,
  savedLayout,
}: UseUnsavedChangesProps) => {
  const router = useRouter();
  const [showModal, setShowModal] = useState(false);
  const pendingUrlRef = useRef<string | null>(null);
  const allowNavigationRef = useRef(false);
  const hasUnsavedChanges = isEditMode && !_.isEqual(dashboardLayout, savedLayout);
  const hasUnsavedRef = useRef(hasUnsavedChanges);
  hasUnsavedRef.current = hasUnsavedChanges;

  useEffect(() => {
    const handleRouteChange = (url: string) => {
      if (allowNavigationRef.current) {
        allowNavigationRef.current = false;
        return;
      }

      if (hasUnsavedRef.current) {
        pendingUrlRef.current = url;
        setShowModal(true);
        router.events.emit('routeChangeError');
        throw new Error('Navigation blocked due to unsaved changes');
      }
    };

    router.events.on('routeChangeStart', handleRouteChange);
    return () => {
      router.events.off('routeChangeStart', handleRouteChange);
    };
  }, [router.events]);

  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
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
      hasUnsavedRef.current = false;
      allowNavigationRef.current = true;
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
