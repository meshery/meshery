import { useEffect, useRef } from 'react';
import { useRouter } from 'next/router';

export interface UseImportDesignQueryActionOptions {
  canImport: boolean;
  openImportModal: () => void;
}

/**
 * Handles the query-driven `?action=import` shortcut on the designs catalog page.
 * When present and permitted, opens the import modal and shallowly strips the parameter
 * so browser reload and history navigation do not re-trigger it.
 */
export function useImportDesignQueryAction({
  canImport,
  openImportModal,
}: UseImportDesignQueryActionOptions): void {
  const router = useRouter();
  const modalOpenedRef = useRef(false);
  const lastReplacedQueryRef = useRef<Record<string, unknown> | null>(null);

  useEffect(() => {
    if (!router.isReady) {
      return;
    }

    const action = router.query.action;

    if (action !== 'import') {
      modalOpenedRef.current = false;
      lastReplacedQueryRef.current = null;
      return;
    }

    if (!modalOpenedRef.current) {
      modalOpenedRef.current = true;
      if (canImport) {
        openImportModal();
      }
    }

    if (lastReplacedQueryRef.current === router.query) {
      return;
    }
    lastReplacedQueryRef.current = router.query;

    const { action: _action, ...restQuery } = router.query;
    router.replace(
      {
        pathname: router.pathname,
        query: restQuery,
      },
      undefined,
      { shallow: true },
    );
  }, [canImport, openImportModal, router, router.isReady, router.pathname, router.query]);
}

export default useImportDesignQueryAction;
