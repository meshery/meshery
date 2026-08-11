/* eslint-disable max-lines -- MeshModelComponent is currently large; prefer refactor/splitting when feasible. */
import React, { useState, useEffect, useLayoutEffect, useCallback, useMemo, useRef } from 'react';
import {
  MODELS,
  COMPONENTS,
  RELATIONSHIPS,
  REGISTRANTS,
  CONNECTIONS,
} from '../../constants/navigator';
import {
  MeshModelToolbar,
  MainContainer,
  TreeWrapper,
  DetailsContainer,
  InnerContainer,
  CardStyle,
  WorkloadsContainer,
} from '@/assets/styles/general/tool.styles';
import MesheryTreeView from './MesheryTreeView';
import MeshModelDetails from './MeshModelDetails';
import { toLower } from 'lodash';
import { useRouter } from 'next/router';
import {
  useLazyGetMeshModelsQuery,
  useLazyGetComponentsQuery,
  useLazyGetRelationshipsQuery,
  useLazyGetRegistrantsQuery,
  useLazyGetConnectionDefinitionsQuery,
  useGetMeshModelsQuery,
  useGetComponentsQuery,
  useGetRelationshipsQuery,
  useGetRegistrantsQuery,
} from '@/rtk-query/meshModel';
import { groupRelationshipsByKind, removeDuplicateVersions } from './helper';
import {
  Button,
  NoSsr,
  AddCircleIcon as AddIcon,
  ExternalLinkIcon as LinkIcon,
  FileUploadIcon as UploadIcon,
  useMediaQuery,
} from '@sistent/sistent';
import { useTheme } from '@/theme';
import { iconSmall } from 'css/icons.styles';
import { useInfiniteScrollRef, useMeshModelComponentRouter } from './hooks';
import ImportModelModal from './ImportModelModal';
import CreateModelModal from './CreateModelModal';
import CreateRelationshipModal from './CreateRelationshipModal';
import MeshModelMobileDetails from './MeshModelMobileDetails';

/**
 * Stable per-row identity used when merging fetched pages into `resourcesDetail`.
 * Prefers the row's own `id` — present on models, components, relationships, and
 * most registrants/connections, and already how MesheryTreeView looks rows up —
 * over a full deep-equality scan. Deep equality is O(n²) on large lists and, for
 * Models specifically, is the wrong notion of "duplicate": two rows can
 * legitimately share every visible field except id (same name+version, different
 * registrant), which is exactly what the Duplicates toggle surfaces, so we must
 * not collapse them. Rows that genuinely lack an id (some connection
 * definitions — see ConnectionDefinitionTree's `connection.id || index`
 * fallback) fall back to a content key so exact repeats across overlapping
 * pages still collapse to one row.
 */
const getRowIdentity = (item: any): string => {
  if (item?.id !== undefined && item?.id !== null) {
    return `id:${item.id}`;
  }
  return `json:${JSON.stringify(item)}`;
};

const dedupeRows = (items: any[]): any[] => {
  const seen = new Set<string>();
  const result: any[] = [];
  for (const item of items) {
    const key = getRowIdentity(item);
    if (!seen.has(key)) {
      seen.add(key);
      result.push(item);
    }
  }
  return result;
};

type MeshModelComponentProps = {
  settingsRouter?: (_router: any) => { handleChangeSelectedTab?: (_tab: string) => void };
  externalView?: string | null;
  externalSearchText?: string | null;
  externalSelectedItemUUID?: string | null;
};

const MeshModelComponent_ = ({
  settingsRouter,
  externalView = null, // External view from modal
  externalSearchText = null, // External search text from modal
  externalSelectedItemUUID = null, // External selected item UUID from modal
}: MeshModelComponentProps) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const router = useRouter();
  // -> use settingsRouter when not in modal mode (Settings page)
  const { handleChangeSelectedTab } =
    settingsRouter && externalView === null
      ? settingsRouter(router)
      : { handleChangeSelectedTab: null };
  const [resourcesDetail, setResourcesDetail] = useState<any[]>([]);
  const { searchQuery, selectedPageSize, selectedTab } = useMeshModelComponentRouter();
  const [page, setPage] = useState<{ [key: string]: number }>({
    Models: 0,
    Components: 0,
    Relationships: 0,
    Registrants: 0,
    Connections: 0,
  });

  // Use external search text if provided, otherwise use query from router
  const [searchText, setSearchText] = useState<string | null>(externalSearchText || searchQuery);
  const [rowsPerPage, setRowsPerPage] = useState(selectedPageSize);
  // Use external view if provided, otherwise use selectedTab or default to 'Models'
  const [view, setView] = useState<string>(externalView || selectedTab || 'Models');
  const [showDetailsData, setShowDetailsData] = useState<{ type: string; data: any }>({
    type: '', // Type of selected data eg. (models, components)
    data: {},
  });
  const [checked, setChecked] = useState<boolean>(false);
  const [isImportModalOpen, setIsImportModalOpen] = useState<boolean>(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState<boolean>(false);
  const [isRelationshipModalOpen, setIsRelationshipModalOpen] = useState<boolean>(false);
  const [modelFilters, setModelsFilters] = useState<{ page: number }>({ page: 0 });
  const [registrantFilters, setRegistrantsFilters] = useState<{ page: number }>({ page: 0 });
  const [componentsFilters, setComponentsFilters] = useState<{ page: number }>({ page: 0 });
  const [relationshipsFilters, setRelationshipsFilters] = useState<{ page: number }>({ page: 0 });
  const [connectionsFilters, setConnectionsFilters] = useState<{ page: number }>({ page: 0 });

  // Monotonically increasing id used by fetchData to detect and discard a
  // response that resolves after a newer fetchData call has already started
  // (see fetchData below).
  const fetchRequestIdRef = useRef(0);
  // Tracks whether the one-shot "fetch every model" backfill (triggered when
  // the Duplicates toggle turns on) has already run for the current Models
  // view, so we don't repeat it on every render (see the effect below).
  const hasFetchedAllModelsRef = useRef(false);

  /**
   * RTK Lazy Queries
   */
  const [getMeshModelsData, modelsRes] = useLazyGetMeshModelsQuery();
  const [getComponentsData, componentsRes] = useLazyGetComponentsQuery();
  const [getRelationshipsData, relationshipsRes] = useLazyGetRelationshipsQuery();
  const [getRegistrantsData, registrantsRes] = useLazyGetRegistrantsQuery();
  const [getConnectionDefinitionsData, connectionsRes] = useLazyGetConnectionDefinitionsQuery();
  // Separate instance from getMeshModelsData above: getRegistrants() below fetches
  // models per-registrant, and sharing the Models-tab query would overwrite
  // modelsRes.data with whichever registrant queried last, corrupting the
  // Models tab's count badge.
  const [getRegistrantModelsData] = useLazyGetMeshModelsQuery();

  /**
   * RTK Queries for counts
   */
  const { data: modelsCountData } = useGetMeshModelsQuery({
    params: { page: 0, pagesize: 1, components: false, relationships: false },
  });
  const { data: componentsCountData } = useGetComponentsQuery({
    params: { page: 0, pagesize: 1, trim: true },
  });
  const { data: relationshipsCountData } = useGetRelationshipsQuery({
    params: { page: 0, pagesize: 1 },
  });
  const { data: registrantsCountData } = useGetRegistrantsQuery({
    params: { page: 0, pagesize: 1 },
  });

  const modelsData = modelsRes.data;
  const registrantsData = registrantsRes.data;
  const componentsData = componentsRes.data;
  const relationshipsData = relationshipsRes.data;
  const connectionsData = connectionsRes.data;

  const hasMoreModels = modelsData?.totalCount > modelsData?.pageSize * modelsData?.page;
  const hasMoreRegistrants =
    registrantsData?.totalCount > registrantsData?.pageSize * registrantsData?.page;
  const hasMoreComponents =
    componentsData?.totalCount > componentsData?.pageSize * componentsData?.page;
  const hasMoreRelationships =
    relationshipsData?.totalCount > relationshipsData?.pageSize * relationshipsData?.page;
  const hasMoreConnections =
    connectionsData?.totalCount > connectionsData?.pageSize * connectionsData?.page;

  const loadNextModelsPage = useCallback(() => {
    if (modelsRes.isLoading || modelsRes.isFetching || !hasMoreModels) {
      return;
    }
    setModelsFilters((prev) => ({ ...prev, page: prev.page + 1 }));
  }, [modelsRes, hasMoreModels]);

  const loadNextRegistrantsPage = useCallback(() => {
    if (registrantsRes.isLoading || registrantsRes.isFetching || !hasMoreRegistrants) {
      return;
    }
    setRegistrantsFilters((prev) => ({ ...prev, page: prev.page + 1 }));
  }, [registrantsRes, hasMoreRegistrants]);

  const loadNextComponentsPage = useCallback(() => {
    if (componentsRes.isLoading || componentsRes.isFetching || !hasMoreComponents) {
      return;
    }
    setComponentsFilters((prev) => ({ ...prev, page: prev.page + 1 }));
  }, [componentsRes, hasMoreComponents]);

  const loadNextRelationshipsPage = useCallback(() => {
    if (relationshipsRes.isLoading || relationshipsRes.isFetching || !hasMoreRelationships) {
      return;
    }
    setRelationshipsFilters((prev) => ({ ...prev, page: prev.page + 1 }));
  }, [relationshipsRes, hasMoreRelationships]);

  const loadNextConnectionsPage = useCallback(() => {
    if (connectionsRes.isLoading || connectionsRes.isFetching || !hasMoreConnections) {
      return;
    }
    setConnectionsFilters((prev) => ({ ...prev, page: prev.page + 1 }));
  }, [connectionsRes, hasMoreConnections]);

  /**
   * IntersectionObservers
   */
  const lastModelRef = useInfiniteScrollRef(loadNextModelsPage);
  const lastComponentRef = useInfiniteScrollRef(loadNextComponentsPage);
  const lastRelationshipRef = useInfiniteScrollRef(loadNextRelationshipsPage);
  const lastRegistrantRef = useInfiniteScrollRef(loadNextRegistrantsPage);
  const lastConnectionRef = useInfiniteScrollRef(loadNextConnectionsPage);

  const getRegistrants = useCallback(async () => {
    let registrantResponse;
    let response;
    registrantResponse = await getRegistrantsData(
      {
        params: {
          page: searchText ? 0 : registrantFilters.page,
          pagesize: searchText ? 'all' : 25,
          search: searchText || '',
        },
      },
      true,
    );
    if (registrantResponse.data && registrantResponse.data.registrants) {
      const registrants = registrantResponse.data.registrants;
      const tempResourcesDetail = [];

      // Fetch each registrant's models in bounded-size batches: fully sequential
      // is slow for large registrant lists, but firing every request at once via
      // Promise.all(registrants.map(...)) can fan out into dozens of concurrent
      // requests and overwhelm the browser/API. Batching keeps a fixed cap on
      // in-flight requests while still parallelizing within each batch.
      const REGISTRANT_FETCH_BATCH_SIZE = 5;
      for (let i = 0; i < registrants.length; i += REGISTRANT_FETCH_BATCH_SIZE) {
        const batch = registrants.slice(i, i + REGISTRANT_FETCH_BATCH_SIZE);
        const batchResults = await Promise.all(
          batch.map(async (registrant) => {
            const hostname = registrant?.hostname ? toLower(registrant.hostname) : '';
            if (!hostname) {
              return null;
            }
            const { data: modelRes } = await getRegistrantModelsData(
              {
                params: {
                  page: 0,
                  pagesize: 'all',
                  registrant: hostname,
                  components: false,
                  relationships: false,
                },
              },
              true,
            );
            if (modelRes?.models && modelRes.models.length > 0) {
              return {
                ...registrant,
                models: removeDuplicateVersions(modelRes.models) || [],
              };
            }
            return null;
          }),
        );
        tempResourcesDetail.push(...batchResults.filter(Boolean));
      }
      response = {
        data: {
          registrants: tempResourcesDetail,
        },
      };
    }
    setRowsPerPage(25);
    return response;
  }, [getRegistrantsData, getRegistrantModelsData, searchText, registrantFilters.page]);

  const fetchData = useCallback(async () => {
    // Captured before any await: if a newer fetchData call starts (e.g. view,
    // search, or a page change) before this one's response comes back, that
    // newer call bumps the ref and this stale response is dropped below
    // instead of being committed over the correct, more recent one. This also
    // covers the case where the search-reset useLayoutEffect and the fetchData
    // useEffect land on the same render: a fetch kicked off with a stale page
    // value doesn't need the effects to run in a particular order to be safe.
    const requestId = ++fetchRequestIdRef.current;
    try {
      let response;
      switch (view) {
        case MODELS:
          response = await getMeshModelsData(
            {
              params: {
                page: searchText ? 0 : modelFilters.page,
                pagesize: searchText ? 'all' : 25,
                components: false,
                relationships: false,
                search: searchText || '',
              },
            },
            true, // arg to use cache as default
          );
          break;
        case COMPONENTS:
          response = await getComponentsData(
            {
              params: {
                page: searchText ? 0 : componentsFilters.page,
                pagesize: searchText ? 'all' : rowsPerPage,
                search: searchText || '',
                trim: true,
              },
            },
            true,
          );
          break;
        case RELATIONSHIPS:
          response = await getRelationshipsData(
            {
              params: {
                page: searchText ? 0 : relationshipsFilters.page,
                pagesize: 'all',
                search: searchText || '',
              },
            },
            true,
          );
          break;
        case REGISTRANTS:
          response = await getRegistrants();
          break;
        case CONNECTIONS: {
          const res = await getConnectionDefinitionsData(
            {
              params: {
                page: searchText ? 0 : connectionsFilters.page,
                pagesize: searchText ? 'all' : rowsPerPage,
                search: searchText || '',
              },
            },
            true,
          );
          // The endpoint returns the page under `connectionDefinitions`; normalize
          // it to the view key (`connections`) the generic handler below expects.
          response = res?.data
            ? { ...res, data: { ...res.data, connections: res.data.connectionDefinitions || [] } }
            : res;
          break;
        }
        default:
          break;
      }
      if (requestId !== fetchRequestIdRef.current) {
        // A newer fetchData call has started since this one began; this
        // response is stale, so don't let it overwrite more recent state.
        return;
      }

      if (response?.data && response.data[view.toLowerCase()]) {
        // Use functional state update to avoid depending on resourcesDetail in the
        // useCallback dependency array (which caused a stale-closure re-fetch loop).
        // Replace vs append is determined by whether this is the first page of the
        // current view, so infinite scroll pagination works correctly in all cases.
        // Relationships always fetch all pages at once (pagesize: 'all'), so they
        // always replace.
        setResourcesDetail((prev) => {
          const fresh = response.data[view.toLowerCase()] ?? [];
          const isFirstPage =
            (view === MODELS && modelFilters.page === 0) ||
            (view === COMPONENTS && componentsFilters.page === 0) ||
            (view === REGISTRANTS && registrantFilters.page === 0) ||
            (view === RELATIONSHIPS && relationshipsFilters.page === 0) ||
            (view === CONNECTIONS && connectionsFilters.page === 0);
          const shouldReplace = searchText || isFirstPage || view === RELATIONSHIPS;
          const newData = shouldReplace ? [...fresh] : [...prev, ...fresh];
          // id-based, O(n) dedupe — see getRowIdentity for why this replaces
          // the deep-equality _.uniqWith scan.
          return dedupeRows(newData);
        });

        // Deeplink may contain higher rowsPerPage val for first time fetch
        // In such case set it to default as 14 after UI renders
        // This ensures the correct pagesize for subsequent API calls triggered on scrolling tree.
        if (rowsPerPage !== 25) {
          setRowsPerPage(25);
        }
      }
    } catch (error) {
      if (requestId !== fetchRequestIdRef.current) {
        // A newer fetchData call has already superseded this one; don't clobber
        // its (possibly successful) result with this stale failure.
        return;
      }
      console.error(`Failed to fetch ${view.toLowerCase()}:`, error);
      setResourcesDetail([]); // Set empty array on error
    }
  }, [
    getMeshModelsData,
    getComponentsData,
    getRelationshipsData,
    getRegistrants,
    getConnectionDefinitionsData,
    modelFilters,
    registrantFilters,
    componentsFilters,
    relationshipsFilters,
    connectionsFilters,
    view,
    rowsPerPage,
    searchText,
  ]);

  // NOTE: The "Duplicates" toggle is a client-side display filter over
  // `resourcesDetail` — toggling it must not clear state or re-trigger the
  // paginated fetchData flow above (that's exactly the infinite re-fetch loop
  // this PR fixes). The one deliberate exception is the backfill effect
  // immediately below: modifyData's duplicate count can only be correct once
  // every model page is loaded, so the first time the toggle turns on we do a
  // single one-shot full-set fetch, independent of fetchData's pagination.

  // Server-side `duplicates` is page-scoped, and modifyData's own recomputation
  // below only sees whatever has been scrolled into `resourcesDetail` so far —
  // so a model's duplicate copy on a not-yet-loaded page is silently missed.
  // Search already asks fetchData for `pagesize: 'all'`, so it's unaffected;
  // this only backfills the plain paginated case.
  useEffect(() => {
    if (view !== MODELS || searchText) {
      hasFetchedAllModelsRef.current = false;
      return;
    }
    if (!checked || hasFetchedAllModelsRef.current) {
      return;
    }
    let cancelled = false;
    (async () => {
      try {
        const response = await getMeshModelsData(
          {
            params: {
              page: 0,
              pagesize: 'all',
              components: false,
              relationships: false,
              search: '',
            },
          },
          true,
        );
        if (cancelled) return;
        const fullSet = response?.data?.models;
        if (fullSet) {
          hasFetchedAllModelsRef.current = true;
          setResourcesDetail(dedupeRows(fullSet));
        }
      } catch (error) {
        console.error('Failed to fetch full model set for duplicate detection:', error);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [checked, view, searchText, getMeshModelsData]);

  const handleTabClick = (selectedView) => {
    // -> use settingsRouter when not in modal mode (Settings page)
    if (handleChangeSelectedTab && externalView === null) {
      handleChangeSelectedTab(selectedView);
    }
    setView(selectedView);
    if (view !== selectedView) {
      setSearchText(null);
      setResourcesDetail([]);
    }
    setModelsFilters({ page: 0 });
    setRegistrantsFilters({ page: 0 });
    setComponentsFilters({ page: 0 });
    setRelationshipsFilters({ page: 0 });
    setConnectionsFilters({ page: 0 });
    setPage({ Models: 0, Components: 0, Relationships: 0, Registrants: 0, Connections: 0 });
    setShowDetailsData({
      type: '',
      data: {},
    });
  };

  const modifyData = () => {
    if (!resourcesDetail) return [];

    if (view === MODELS) {
      if (!checked) {
        return removeDuplicateVersions(resourcesDetail);
      }
      // Server's `duplicates` field is page-scoped; recompute across all loaded models.
      const keyOf = (m) => `${m?.name}@${m?.model?.version}`;
      const countByKey: Record<string, number> = {};
      resourcesDetail.forEach((m) => {
        countByKey[keyOf(m)] = (countByKey[keyOf(m)] || 0) + 1;
      });
      const dupeRows = resourcesDetail.filter((m) => countByKey[keyOf(m)] > 1);
      return removeDuplicateVersions(dupeRows);
    } else if (view === RELATIONSHIPS) {
      return groupRelationshipsByKind(resourcesDetail);
    } else if (view === REGISTRANTS) {
      return resourcesDetail || [];
    } else {
      return resourcesDetail;
    }
  };

  // Memoize so MesheryTreeView receives the same array reference when nothing
  // has changed. This is what makes the O(1) referential guard in
  // MesheryTreeView.tsx safe (prevState.data === data).
  const treeData = useMemo(modifyData, [resourcesDetail, view, checked]);

  // Reset every view's page state whenever a search opens or closes. fetchData's
  // replace-vs-append decision (isFirstPage) trusts each view's *Filters.page to
  // be 0 at that point — otherwise, after infinite-scrolling a view and then
  // opening/clearing search, a stale non-zero page gets requested into a
  // freshly emptied list, silently dropping earlier results.
  // Uses useLayoutEffect (not useEffect) so this reset is flushed synchronously
  // before the fetchData effect below runs, guaranteeing fetchData never sees
  // a stale non-zero *Filters.page value on the same render pass.
  useLayoutEffect(() => {
    const isStale =
      page[view] > 0 ||
      modelFilters.page > 0 ||
      componentsFilters.page > 0 ||
      registrantFilters.page > 0 ||
      relationshipsFilters.page > 0 ||
      connectionsFilters.page > 0;
    if (isStale) {
      setModelsFilters({ page: 0 });
      setRegistrantsFilters({ page: 0 });
      setComponentsFilters({ page: 0 });
      setRelationshipsFilters({ page: 0 });
      setConnectionsFilters({ page: 0 });
      setPage({ Models: 0, Components: 0, Relationships: 0, Registrants: 0, Connections: 0 });
    }
  }, [searchText]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Sync view state with externalView or selectedTab (for modal or route usage)
  useEffect(() => {
    const newView =
      externalView ?? (typeof selectedTab === 'string' ? selectedTab : selectedTab?.[0]);
    if (newView && newView !== view) {
      setView(newView);
      setResourcesDetail([]);
      setSearchText(externalSearchText || null);
      setModelsFilters({ page: 0 });
      setRegistrantsFilters({ page: 0 });
      setComponentsFilters({ page: 0 });
      setRelationshipsFilters({ page: 0 });
      setConnectionsFilters({ page: 0 });
      setPage({ Models: 0, Components: 0, Relationships: 0, Registrants: 0, Connections: 0 });
      setShowDetailsData({
        type: '',
        data: {},
      });
    }
  }, [externalView, selectedTab, externalSearchText]);

  useEffect(() => {
    if (externalSearchText !== null && externalSearchText !== searchText) {
      setSearchText(externalSearchText);
    }
  }, [externalSearchText]);

  return (
    <WorkloadsContainer data-test="workloads">
      <ImportModelModal
        isImportModalOpen={isImportModalOpen}
        setIsImportModalOpen={setIsImportModalOpen}
      />
      <CreateModelModal
        isCreateModalOpen={isCreateModalOpen}
        setIsCreateModalOpen={setIsCreateModalOpen}
      />
      <CreateRelationshipModal
        isRelationshipModalOpen={isRelationshipModalOpen}
        setIsRelationshipModalOpen={setIsRelationshipModalOpen}
      />

      <MainContainer>
        {(view === MODELS || view === RELATIONSHIPS) && (
          <TabBar
            openImportModal={() => setIsImportModalOpen(true)}
            openCreateModal={() => setIsCreateModalOpen(true)}
            openRelationshipModal={() => setIsRelationshipModalOpen(true)}
            view={view}
          />
        )}
        {externalView === null && (
          <InnerContainer>
            <TabCard
              label="Models"
              count={modelsData?.totalCount ?? modelsCountData?.totalCount ?? 0}
              active={view === MODELS}
              onClick={() => handleTabClick(MODELS)}
            />
            <TabCard
              label="Components"
              count={componentsData?.totalCount ?? componentsCountData?.totalCount ?? 0}
              active={view === COMPONENTS}
              onClick={() => handleTabClick(COMPONENTS)}
            />
            <TabCard
              label="Relationships"
              count={relationshipsData?.totalCount ?? relationshipsCountData?.totalCount ?? 0}
              active={view === RELATIONSHIPS}
              onClick={() => handleTabClick(RELATIONSHIPS)}
            />
            <TabCard
              label="Registrants"
              count={registrantsData?.totalCount ?? registrantsCountData?.totalCount ?? 0}
              active={view === REGISTRANTS}
              onClick={() => handleTabClick(REGISTRANTS)}
            />
            <TabCard
              label="Connections"
              count={connectionsData?.totalCount || 0}
              active={view === CONNECTIONS}
              onClick={() => handleTabClick(CONNECTIONS)}
            />
          </InnerContainer>
        )}

        <TreeWrapper>
          <DetailsContainer
            isEmpty={!resourcesDetail.length}
            style={{
              padding: '0.6rem',
              overflow: 'auto',
              height: '100%',
              flex: 1,
              display: 'flex',
              flexDirection: 'column',
            }}
          >
            <MesheryTreeView
              data={treeData}
              view={view}
              setSearchText={setSearchText}
              setPage={setPage}
              checked={checked}
              setChecked={setChecked}
              // Whether ANY (unfiltered) models are loaded. Used to decide
              // whether the Duplicates switch should be usable — must NOT be
              // derived from treeData/data, since that's already filtered by
              // `checked` and would make the switch disable itself the
              // moment it finds zero duplicates, permanently locking it on.
              hasResourcesLoaded={view === MODELS ? resourcesDetail.length > 0 : true}
              searchText={searchText}
              setShowDetailsData={setShowDetailsData}
              showDetailsData={showDetailsData}
              setResourcesDetail={setResourcesDetail}
              setModelsFilters={setModelsFilters}
              externalSelectedItemUUID={externalSelectedItemUUID} // Pass external UUID
              isModalMode={externalView !== null} // Modal mode // external view
              lastItemRef={{
                [MODELS]: lastModelRef,
                [REGISTRANTS]: lastRegistrantRef,
                [COMPONENTS]: lastComponentRef,
                [RELATIONSHIPS]: lastRelationshipRef,
                [CONNECTIONS]: lastConnectionRef,
              }}
              isFetching={{
                [MODELS]: modelsRes.isFetching,
                [REGISTRANTS]: registrantsRes.isFetching,
                [COMPONENTS]: componentsRes.isFetching,
                [RELATIONSHIPS]: relationshipsRes.isFetching,
                [CONNECTIONS]: connectionsRes.isFetching,
              }}
              isLoading={{
                [MODELS]: modelsRes.isLoading,
                [REGISTRANTS]: registrantsRes.isLoading,
                [COMPONENTS]: componentsRes.isLoading,
                [RELATIONSHIPS]: relationshipsRes.isLoading,
                [CONNECTIONS]: connectionsRes.isLoading,
              }}
            />
          </DetailsContainer>
          {isMobile ? (
            <MeshModelMobileDetails
              view={view}
              showDetailsData={showDetailsData}
              setShowDetailsData={setShowDetailsData}
            />
          ) : (
            <MeshModelDetails view={view} showDetailsData={showDetailsData} />
          )}
        </TreeWrapper>
      </MainContainer>
    </WorkloadsContainer>
  );
};

const TabBar = ({ openImportModal, openCreateModal, view, openRelationshipModal }) => {
  return (
    <MeshModelToolbar>
      <div
        style={{
          display: 'flex',
          justifyContent: 'flex-start',
          alignItems: 'center',
          gap: '0.75rem',
          flexWrap: 'wrap',
        }}
      >
        {view === MODELS && (
          <>
            <Button
              aria-label="Create Model"
              variant="contained"
              color="primary"
              onClick={openCreateModal}
              style={{ display: 'flex' }}
              disabled={false} //TODO: Need to make key for this component
              startIcon={<AddIcon style={iconSmall} />}
              data-testid="TabBar-Button-CreateModel"
            >
              Create Model
            </Button>
            <Button
              aria-label="Import Model"
              variant="contained"
              color="primary"
              onClick={openImportModal}
              style={{ display: 'flex' }}
              disabled={false} //TODO: Need to make key for this component
              startIcon={<UploadIcon />}
              data-testid="TabBar-Button-ImportModel"
            >
              Import Model
            </Button>
          </>
        )}

        {view === RELATIONSHIPS && (
          <Button
            aria-label="Create Relationship"
            variant="contained"
            color="primary"
            onClick={openRelationshipModal}
            style={{ display: 'flex' }}
            disabled={false}
            startIcon={<LinkIcon />}
            data-testid="TabBar-Button-CreateRelationship"
          >
            Create Relationship
          </Button>
        )}
      </div>
      {/*
      This builk operation is not yet supported
      <DisableButton disabled variant="contained" startIcon={<DoNotDisturbOnIcon />}>
        Ignore
      </DisableButton> */}
    </MeshModelToolbar>
  );
};

const TabCard = ({ label, count, active, onClick }) => {
  return (
    <CardStyle isSelected={active} elevation={3} onClick={onClick}>
      <span
        style={{
          fontSize: '1rem',
          marginLeft: '4px',
        }}
      >
        {`(${count?.toLocaleString() || 0})`}
      </span>
      {label}
    </CardStyle>
  );
};

const MeshModelComponent = (props) => {
  return (
    <NoSsr>
      <MeshModelComponent_ {...props} />
    </NoSsr>
  );
};

export default MeshModelComponent;
