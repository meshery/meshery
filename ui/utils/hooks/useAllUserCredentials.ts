import { useEffect, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { createSelector } from '@reduxjs/toolkit';
import {
  mesheryApi,
  useGetUserCredentialsQuery,
  type GetUserCredentialsApiResponse,
} from '@meshery/schemas/mesheryApi';

// The most GET /api/integrations/credentials serves per page; the handler caps
// any larger request to this.
export const CREDENTIALS_PAGE_SIZE = 25;

type Credential = NonNullable<GetUserCredentialsApiResponse['credentials']>[number];

const pageArg = (page: number) => ({ page, pageSize: CREDENTIALS_PAGE_SIZE });

type PageSelector = ReturnType<typeof mesheryApi.endpoints.getUserCredentials.select>;
type ApiState = Parameters<PageSelector>[0];
const NO_PAGES: ReturnType<PageSelector>[] = [];

/**
 * Every credential of the current user, across all pages of the
 * schemas-generated `getUserCredentials` endpoint.
 *
 * A credential picker has to offer all of them, but the endpoint pages its
 * results, and a request without paging params returns only the 10 newest.
 * The first page reports `totalCount` and the page size the server actually
 * applied, which can be smaller than the size requested, so the number of
 * further pages is derived from that reply rather than from the request.
 *
 * The further pages are held as cache subscriptions and read back from the
 * cache, so a credential saved, updated or deleted elsewhere refetches them
 * through the endpoint's tags exactly as it refetches the first page.
 */
export const useAllUserCredentials = ({ skip = false }: { skip?: boolean } = {}) => {
  const dispatch = useDispatch();
  const first = useGetUserCredentialsQuery(pageArg(0), { skip });

  const totalCount = first.data?.totalCount ?? 0;
  const servedPageSize = first.data?.pageSize || CREDENTIALS_PAGE_SIZE;
  const pageCount = skip || !first.data ? 0 : Math.ceil(totalCount / servedPageSize);

  useEffect(() => {
    if (pageCount <= 1) return undefined;
    const subscriptions = Array.from({ length: pageCount - 1 }, (_, i) =>
      dispatch(mesheryApi.endpoints.getUserCredentials.initiate(pageArg(i + 1))),
    );
    return () => subscriptions.forEach((subscription) => subscription.unsubscribe());
  }, [dispatch, pageCount]);

  const selectRestPages = useMemo(() => {
    if (pageCount <= 1) return () => NO_PAGES;
    return createSelector(
      Array.from({ length: pageCount - 1 }, (_, i) =>
        mesheryApi.endpoints.getUserCredentials.select(pageArg(i + 1)),
      ),
      (...pages) => pages,
    );
  }, [pageCount]);
  const restPages = useSelector((state: ApiState) => selectRestPages(state));

  const credentials = useMemo(() => {
    const byId = new Map<string, Credential>();
    const pages = [first.data, ...restPages.map((page) => page.data)];
    pages.forEach((page) =>
      page?.credentials?.forEach((credential) => {
        // A credential created while the pages load shifts the rest by one, so
        // a row can arrive on two pages; it is listed once.
        byId.set(String(credential.id), credential);
      }),
    );
    return Array.from(byId.values());
  }, [first.data, restPages]);

  return {
    credentials,
    isLoading: first.isLoading || restPages.some((page) => page.isLoading),
    isError: first.isError || restPages.some((page) => page.isError),
  };
};
