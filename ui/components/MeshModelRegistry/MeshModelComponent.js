import React, { useState, useEffect, useCallback, useMemo } from 'react';
import UploadIcon from '@mui/icons-material/Upload';
import DoNotDisturbOnIcon from '@mui/icons-material/DoNotDisturbOn';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import {
  MeshModelToolbar,
  MainContainer,
  InnerContainer,
  CardStyle,
  TreeWrapper,
  DetailsContainer,
} from '@/assets/styles/general/tool.styles';
import MesheryTreeView from './MesheryTreeView';
import MeshModelDetails from './MeshModelDetails';
import { toLower } from 'lodash';
import { DisableButton } from './MeshModel.style';
import { useRouter } from 'next/router';
import { Provider } from 'react-redux';
import { store } from '../../store';
import {
  useLazyGetMeshModelsQuery,
  useLazyGetComponentsQuery,
  useLazyGetRelationshipsQuery,
  useLazyGetRegistrantsQuery,
} from '@/rtk-query/meshModel';
import { groupRelationshipsByKind, removeDuplicateVersions } from './helper';
import _ from 'lodash';
import { Button, NoSsr } from '@layer5/sistent';
import { iconSmall } from '../../css/icons.styles';
import AddIcon from '@mui/icons-material/AddCircleOutline';
import { useInfiniteScrollRef, useMeshModelComponentRouter } from './hooks';
import ImportModelModal from './ImportModelModal';
import CreateModelModal from './CreateModelModal';

const MeshModelComponent_ = ({
  modelsCount,
  componentsCount,
  relationshipsCount,
  registrantCount,
  settingsRouter,
}) => {
  const router = useRouter();
  const { handleChangeSelectedTab } = settingsRouter(router);
  const { searchQuery, selectedPageSize } = useMeshModelComponentRouter();
  const [checked, setChecked] = useState(false);

  const [filters, setFilters] = useState({
    models: { page: 0 },
    registrants: { page: 0 },
    components: { page: 0 },
    relationships: { page: 0 },
  });

  const [viewState, setViewState] = useState({
    view: MODELS,
    searchText: searchQuery,
    rowsPerPage: selectedPageSize,
    resourcesDetail: [],
    showDetailsData: { type: '', data: {} },
  });

  const [isImportModalOpen, setIsImportModalOpen] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  /**
   * RTK Lazy Queries
   */
  const [getMeshModelsData, modelsRes] = useLazyGetMeshModelsQuery();
  const [getComponentsData, componentsRes] = useLazyGetComponentsQuery();
  const [getRelationshipsData, relationshipsRes] = useLazyGetRelationshipsQuery();
  const [getRegistrantsData, registrantsRes] = useLazyGetRegistrantsQuery();

  const modelsData = modelsRes.data;
  const registrantsData = registrantsRes.data;
  const componentsData = componentsRes.data;
  const relationshipsData = relationshipsRes.data;

  const hasMoreModels = modelsData?.total_count > modelsData?.page_size * modelsData?.page;
  const hasMoreRegistrants =
    registrantsData?.total_count > registrantsData?.page_size * registrantsData?.page;
  const hasMoreComponents =
    componentsData?.total_count > componentsData?.page_size * componentsData?.page;
  const hasMoreRelationships =
    componentsData?.total_count > relationshipsData?.page_size * relationshipsData?.page;

  const loadNextModelsPage = useCallback(() => {
    if (modelsRes.isLoading || modelsRes.isFetching || !hasMoreModels) {
      return;
    }
    setFilters((prev) => ({ ...prev, models: { page: prev.models.page + 1 } }));
  }, [modelsRes, hasMoreModels]);

  const loadNextRegistrantsPage = useCallback(() => {
    if (registrantsRes.isLoading || registrantsRes.isFetching || !hasMoreRegistrants) {
      return;
    }
    setFilters((prev) => ({ ...prev, registrants: { page: prev.registrants.page + 1 } }));
  }, [registrantsRes, hasMoreRegistrants]);

  const loadNextComponentsPage = useCallback(() => {
    if (componentsRes.isLoading || componentsRes.isFetching || !hasMoreComponents) {
      return;
    }
    setFilters((prev) => ({ ...prev, components: { page: prev.components.page + 1 } }));
  }, [componentsRes, hasMoreComponents]);

  const loadNextRelationshipsPage = useCallback(() => {
    if (relationshipsRes.isLoading || relationshipsRes.isFetching || !hasMoreRelationships) {
      return;
    }
  }, [relationshipsRes, hasMoreRelationships]);

  /**
   * IntersectionObservers
   */
  const lastModelRef = useInfiniteScrollRef(loadNextModelsPage);
  const lastComponentRef = useInfiniteScrollRef(loadNextComponentsPage);
  const lastRelationshipRef = useInfiniteScrollRef(loadNextRelationshipsPage);
  const lastRegistrantRef = useInfiniteScrollRef(loadNextRegistrantsPage);

  const modifiedData = useMemo(() => {
    if (viewState.view === MODELS) {
      return removeDuplicateVersions(
        checked
          ? viewState.resourcesDetail.filter((model) => model.duplicates > 0)
          : viewState.resourcesDetail,
      );
    }
    if (viewState.view === RELATIONSHIPS) {
      return groupRelationshipsByKind(viewState.resourcesDetail);
    }
    return viewState.resourcesDetail;
  }, [viewState.view, viewState.resourcesDetail, checked]);

  const debouncedFetchData = useCallback(
    _.debounce(async () => {
      try {
        const response = await fetchDataForView(viewState.view, {
          page: viewState.searchText ? 0 : filters[viewState.view.toLowerCase()].page,
          pagesize: viewState.searchText
            ? 'all'
            : viewState.view === RELATIONSHIPS
              ? 'all'
              : viewState.rowsPerPage,
          search: viewState.searchText || '',
        });

        if (response?.data) {
          setViewState((prev) => ({
            ...prev,
            resourcesDetail: handleNewData(
              response.data,
              prev.resourcesDetail,
              viewState.searchText,
            ),
          }));
        }
      } catch (error) {
        console.error('Fetch failed:', error);
      }
    }, 300),
    [viewState.view, filters, viewState.searchText],
  );

  const fetchDataForView = async (view, params) => {
    const searchParam = params.search || '';
    const pageParam = typeof params.page === 'number' ? params.page : 0;

    const baseParams = {
      ...params,
      search: searchParam,
      page: pageParam,
      pagesize: searchParam ? 'all' : params.pagesize || 25,
    };

    switch (view) {
      case MODELS:
        return await getMeshModelsData(
          {
            params: {
              ...baseParams,
              trim: false,
              components: false,
              relationships: false,
            },
          },
          true,
        );
      case COMPONENTS:
        return await getComponentsData({ params }, true);
      case RELATIONSHIPS:
        return await getRelationshipsData({ params }, true);
      case REGISTRANTS:
        return await getRegistrants();
      default:
        return null;
    }
  };

  const getRegistrants = async () => {
    let registrantResponse = await getRegistrantsData({
      params: {
        page: viewState.searchText ? 0 : filters.registrants.page,
        pagesize: viewState.searchText ? 'all' : 25,
        search: viewState.searchText || '',
      },
    });

    if (!registrantResponse.data?.registrants) {
      return registrantResponse;
    }

    const registrants = registrantResponse.data.registrants;
    const tempResourcesDetail = [];

    for (let registrant of registrants) {
      let hostname = toLower(registrant?.hostname);
      const modelRes = await getMeshModelsData({
        params: {
          page: filters.models.page,
          pagesize: 'all',
          registrant: hostname,
          components: false,
          relationships: false,
        },
      });

      if (modelRes.data?.models?.length > 0) {
        tempResourcesDetail.push({
          ...registrant,
          models: removeDuplicateVersions(modelRes.data.models) || [],
        });
      }
    }

    return {
      data: {
        registrants: tempResourcesDetail,
      },
    };
  };

  const handleTabClick = useCallback(
    (selectedView) => {
      handleChangeSelectedTab(selectedView);
      setViewState((prev) => ({
        ...prev,
        view: selectedView,
        searchText: '',
        resourcesDetail: [],
        showDetailsData: { type: '', data: {} },
      }));
      setFilters({
        models: { page: 0 },
        registrants: { page: 0 },
        components: { page: 0 },
        relationships: { page: 0 },
      });
    },
    [handleChangeSelectedTab],
  );

  const handleNewData = (newData, existingData, searchText) => {
    // Get the appropriate array from response data based on view type
    const dataArray = newData?.[viewState.view.toLowerCase()] || [];

    if (searchText) {
      return _.uniqWith(dataArray, _.isEqual);
    }

    // Ensure both arrays exist before combining
    const existingArray = Array.isArray(existingData) ? existingData : [];
    return _.uniqWith([...existingArray, ...dataArray], _.isEqual);
  };

  useEffect(() => {
    // Only reset pagination if searchText is being set to empty string
    if (viewState.searchText !== null && viewState.searchText !== undefined) {
      setFilters({
        models: { page: 0 },
        registrants: { page: 0 },
        components: { page: 0 },
        relationships: { page: 0 },
      });
    }
  }, [viewState.searchText]);

  useEffect(() => {
    debouncedFetchData();
  }, [viewState.view, filters, viewState.rowsPerPage, checked, viewState.searchText]);

  return (
    <div data-test="workloads">
      <TabBar
        openImportModal={() => setIsImportModalOpen(true)}
        openCreateModal={() => setIsCreateModalOpen(true)}
      />

      <ImportModelModal
        isImportModalOpen={isImportModalOpen}
        setIsImportModalOpen={setIsImportModalOpen}
      />
      <CreateModelModal
        isCreateModalOpen={isCreateModalOpen}
        setIsCreateModalOpen={setIsCreateModalOpen}
      />

      <MainContainer>
        <InnerContainer>
          <TabCard
            label="Models"
            count={modelsCount}
            active={viewState.view === MODELS}
            onClick={() => handleTabClick(MODELS)}
          />
          <TabCard
            label="Components"
            count={componentsCount}
            active={viewState.view === COMPONENTS}
            onClick={() => handleTabClick(COMPONENTS)}
          />
          <TabCard
            label="Relationships"
            count={relationshipsCount}
            active={viewState.view === RELATIONSHIPS}
            onClick={() => handleTabClick(RELATIONSHIPS)}
          />
          <TabCard
            label="Registrants"
            count={registrantCount}
            active={viewState.view === REGISTRANTS}
            onClick={() => handleTabClick(REGISTRANTS)}
          />
        </InnerContainer>

        <TreeWrapper>
          <DetailsContainer
            isEmpty={!viewState.resourcesDetail.length}
            style={{
              padding: '0.6rem',
              overflow: 'hidden',
            }}
          >
            <MesheryTreeView
              data={modifiedData}
              view={viewState.view}
              setSearchText={(text) => setViewState((prev) => ({ ...prev, searchText: text }))}
              setPage={(page) =>
                setFilters((prev) => ({ ...prev, [viewState.view.toLowerCase()]: { page } }))
              }
              checked={checked}
              setChecked={setChecked}
              searchText={viewState.searchText}
              setShowDetailsData={(data) =>
                setViewState((prev) => ({ ...prev, showDetailsData: data }))
              }
              showDetailsData={viewState.showDetailsData}
              setResourcesDetail={(data) =>
                setViewState((prev) => ({ ...prev, resourcesDetail: data }))
              }
              lastItemRef={{
                [MODELS]: lastModelRef,
                [REGISTRANTS]: lastRegistrantRef,
                [COMPONENTS]: lastComponentRef,
                [RELATIONSHIPS]: lastRelationshipRef,
              }}
              isFetching={{
                [MODELS]: modelsRes.isFetching,
                [REGISTRANTS]: registrantsRes.isFetching,
                [COMPONENTS]: componentsRes.isFetching,
                [RELATIONSHIPS]: relationshipsRes.isFetching,
              }}
              isLoading={{
                [MODELS]: modelsRes.isLoading,
                [REGISTRANTS]: registrantsRes.isLoading,
                [COMPONENTS]: componentsRes.isLoading,
                [RELATIONSHIPS]: relationshipsRes.isLoading,
              }}
            />
          </DetailsContainer>
          <MeshModelDetails
            view={viewState.view}
            setShowDetailsData={(data) =>
              setViewState((prev) => ({ ...prev, showDetailsData: data }))
            }
            showDetailsData={viewState.showDetailsData}
          />
        </TreeWrapper>
      </MainContainer>
    </div>
  );
};

const TabBar = ({ openImportModal, openCreateModal }) => {
  return (
    <MeshModelToolbar>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          gap: '1rem', // Add some space between buttons
        }}
      >
        <Button
          aria-label="Add Pattern"
          variant="contained"
          color="primary"
          onClick={openCreateModal}
          style={{ display: 'flex' }}
          disabled={false} //TODO: Need to make key for this component
          startIcon={<AddIcon style={iconSmall} />}
        >
          Create
        </Button>
        <Button
          aria-label="Add Pattern"
          variant="contained"
          color="primary"
          onClick={openImportModal}
          style={{ display: 'flex' }}
          disabled={false} //TODO: Need to make key for this component
          startIcon={<UploadIcon />}
        >
          Import
        </Button>
      </div>
      <DisableButton disabled variant="contained" startIcon={<DoNotDisturbOnIcon />}>
        Ignore
      </DisableButton>
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
        {`(${count})`}
      </span>
      {label}
    </CardStyle>
  );
};

const MeshModelComponent = (props) => {
  return (
    <NoSsr>
      <Provider store={store}>
        <MeshModelComponent_ {...props} />
      </Provider>
    </NoSsr>
  );
};
export default MeshModelComponent;
