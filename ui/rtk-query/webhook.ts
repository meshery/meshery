import { api, mesheryApiPath } from './index';
import { shouldOverrideExisting } from './utils';

const webhookApi = api.injectEndpoints({
  overrideExisting: shouldOverrideExisting,
  endpoints: (builder) => ({
    supportWebHook: builder.mutation({
      query: (queryArg) => ({
        url: mesheryApiPath(`extensions/api/webhook/${queryArg.type}`),
        method: 'POST',
        body: queryArg.body,
      }),
    }),
  }),
});

export const { useSupportWebHookMutation } = webhookApi;
