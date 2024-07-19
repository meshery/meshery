import { useSelectorRtk, useDispatchRtk, useStoreRtk } from '@/store/hooks';
import { getCurrentOrg, getCurrentWorkspace, selectCurrentOrg } from '@/store/slices/globalContext';
import { CURRENT_ORG_KEY, CURRENT_WORKSPACE_KEY } from '@/utils/Enum';
import {
  coreModule,
  reactHooksModule,
  buildCreateApi,
  fetchBaseQuery,
} from '@reduxjs/toolkit/query/react';

const createApi = buildCreateApi(
  coreModule(),
  reactHooksModule({
    useSelector: useSelectorRtk,
    useDispatch: useDispatchRtk,
    useStore: useStoreRtk,
  }),
);
export const api = createApi({
  reducerPath: 'mesheryApi',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api/',
    credentials: 'include',
    prepareHeaders: (headers, { getState }) => {
      const state = getState();
      const currentOrg = getCurrentOrg(state);
      const currentWorkspace = getCurrentWorkspace(state);
      headers.set(CURRENT_ORG_KEY, currentOrg?.id);
      headers.set(CURRENT_WORKSPACE_KEY, currentWorkspace?.id);
      return headers;
    },
  }),
  endpoints: () => ({}),
});
