import { api, mesheryApiPath } from './index';

const webhookApi = api.injectEndpoints({
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
