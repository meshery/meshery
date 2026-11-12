import { useEffect, useRef } from 'react';
import { useDispatch } from 'react-redux';
import { useGetUserPrefQuery, useUpdateUserPrefMutation } from '@/rtk-query/user';
import { toggleCatalogContent } from '@/store/slices/mesheryUi';

export const CATALOG_CONTENT_PREFERENCE_KEY = 'catalog-content-preference';

export const useSyncCatalogContentPreference = () => {
  const dispatch = useDispatch();
  const { data } = useGetUserPrefQuery();
  const [, { isLoading: isSaving }] = useUpdateUserPrefMutation({
    fixedCacheKey: CATALOG_CONTENT_PREFERENCE_KEY,
  });
  const savingRef = useRef(isSaving);
  const catalogContent = data?.usersExtensionPreferences?.catalogContent;

  useEffect(() => {
    savingRef.current = isSaving;
  }, [isSaving]);

  // Only a changed saved value should sync. Finishing a save must not replay
  // cached data over the newer optimistic choice or its explicit rollback.
  useEffect(() => {
    if (typeof catalogContent === 'boolean' && !savingRef.current) {
      dispatch(toggleCatalogContent({ catalogVisibility: catalogContent }));
    }
  }, [catalogContent, dispatch]);
};
