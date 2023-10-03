import { api } from './index';

const connectionsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getConnectionStatus: builder.query({
        query: queryArg => ({
            url: `integrations/connections/${queryArg.connectionKind}/status`,
            params: { id: queryArg.repoURL }
        })
    }),
    getConnectionDetails: builder.query({
        query: queryArg => ({
            url: `integrations/connections/${queryArg.connectionKind}/details`,
            params: { id: queryArg.repoURL }
        })
    }),
    verifyConnectionURL: builder.mutation({
        query: queryArg => ({
            url: `integrations/connections/${queryArg.connectionKind}/verify`,
            method: "POST",
            params: { id: queryArg.repoURL }
        })
    }),
    connectionMetaData: builder.mutation({
        query: queryArg => ({
            url: `integrations/connections/${queryArg.connectionKind}/metadata`,
            method: "POST",
            body: queryArg.body
        })
    }),
    configureConnection: builder.mutation({
        query: queryArg => ({
            url: `integrations/connections/${queryArg.connectionKind}/configure`,
            method: "POST",
            body: queryArg.body
        })
    }),
  }),
});

export const { useGetConnectionStatusQuery, useLazyGetConnectionDetailsQuery, useVerifyConnectionURLMutation, useConnectionMetaDataMutation, useConfigureConnectionMutation } = connectionsApi;
