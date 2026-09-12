import { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';
import { getPath } from '../../lib/path';

export function usePageTitle(title: string): void {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title }));
  }, [dispatch, title]);
}
