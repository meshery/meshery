import { useSelectorRtk, useDispatchRtk, useStoreRtk } from '@/store/hooks';
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
  baseQuery: fetchBaseQuery({ baseUrl: '/api/' }),
  endpoints: () => ({}),
});
