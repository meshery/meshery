import { api } from './index';

const webhookApi = api.injectEndpoints({
  endpoints: (builder) => ({
    supportWebHook: builder.mutation({
      query: (queryArg) => ({
        url: `extensions/api/webhook/${queryArg.type}`,
        method: 'POST',
        body: queryArg.body,
      }),
    }),
  }),
});

export const { useSupportWebHookMutation } = webhookApi;
