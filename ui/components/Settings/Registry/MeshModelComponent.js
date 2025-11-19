import React, { useState, useEffect, useCallback } from 'react';
import UploadIcon from '@mui/icons-material/Upload';
import AddIcon from '@mui/icons-material/AddCircleOutline';
import LinkIcon from '@mui/icons-material/Link';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../../constants/navigator';
import {
  MeshModelToolbar,
  MainContainer,
  TreeWrapper,
  DetailsContainer,
  InnerContainer,
  CardStyle,
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
} from '@/rtk-query/meshModel';
import { groupRelationshipsByKind, removeDuplicateVersions } from './helper';
import _ from 'lodash';
import { Button, NoSsr } from '@sistent/sistent';
import { iconSmall } from 'css/icons.styles';
import { useInfiniteScrollRef, useMeshModelComponentRouter } from './hooks';
import ImportModelModal from './ImportModelModal';
import CreateModelModal from './CreateModelModal';
import CreateRelationshipModal from '@/components/RelationshipBuilder/CreateRelationshipModal';

const MeshModelComponent_ = ({
  settingsRouter,
  externalView = null, // External view from modal
  externalSearchText = null, // External search text from modal
  externalSelectedItemUUID = null, // External selected item UUID from modal
}) => {
  const router = useRouter();
  // -> use settingsRouter when not in modal mode (Settings page)
  const { handleChangeSelectedTab } =
    settingsRouter && externalView === null
      ? settingsRouter(router)
      : { handleChangeSelectedTab: null };
  const [resourcesDetail, setResourcesDetail] = useState([]);
  const { searchQuery, selectedPageSize, selectedTab } = useMeshModelComponentRouter();
  const [page, setPage] = useState({
    Models: 0,
    Components: 0,
    Relationships: 0,
    Registrants: 0,
  });

  // Use external search text if provided, otherwise use query from router
  const [searchText, setSearchText] = useState(externalSearchText || searchQuery);
  const [rowsPerPage, setRowsPerPage] = useState(selectedPageSize);
  // Use external view if provided, otherwise use selectedTab or default to 'Models'
  const [view, setView] = useState(externalView || selectedTab || 'Models');
  const [showDetailsData, setShowDetailsData] = useState({
    type: '', // Type of selected data eg. (models, components)
    data: {},
  });
  const [checked, setChecked] = useState(false);
  const [isImportModalOpen, setIsImportModalOpen] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isRelationshipModalOpen, setIsRelationshipModalOpen] = useState(false);
  const [modelFilters, setModelsFilters] = useState({ page: 0 });
  const [registrantFilters, setRegistrantsFilters] = useState({ page: 0 });
  const [componentsFilters, setComponentsFilters] = useState({ page: 0 });
  const [relationshipsFilters, setRelationshipsFilters] = useState({ page: 0 });

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
  }, [relationshipsRes, hasMoreRelationships]);
  /**
   * IntersectionObservers
   */
  const lastModelRef = useInfiniteScrollRef(loadNextModelsPage);
  const lastComponentRef = useInfiniteScrollRef(loadNextComponentsPage);
  const lastRelationshipRef = useInfiniteScrollRef(loadNextRelationshipsPage);
  const lastRegistrantRef = useInfiniteScrollRef(loadNextRegistrantsPage);

  const fetchData = useCallback(async () => {
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
        default:
          break;
      }
      if (response?.data && response.data[view.toLowerCase()]) {
        // When search or "show duplicates" functionality is active:
        // Avoid appending data to the previous dataset.
        // preventing duplicate entries and ensuring the UI reflects the API's response accurately.
        // For instance, during a search, display the data returned by the API instead of appending it to the previous results.
        let newData = [];
        if (response.data[view.toLowerCase()]) {
          newData =
            searchText || view === RELATIONSHIPS
              ? [...response.data[view.toLowerCase()]]
              : [...resourcesDetail, ...response.data[view.toLowerCase()]];
        }

        // Set unique data
        setResourcesDetail(_.uniqWith(newData, _.isEqual));

        // Deeplink may contain higher rowsPerPage val for first time fetch
        // In such case set it to default as 14 after UI renders
        // This ensures the correct pagesize for subsequent API calls triggered on scrolling tree.
        if (rowsPerPage !== 25) {
          setRowsPerPage(25);
        }
      }
    } catch (error) {
      console.error(`Failed to fetch ${view.toLowerCase()}:`, error);
      setResourcesDetail([]); // Set empty array on error
    }
  }, [
    getMeshModelsData,
    getComponentsData,
    getRelationshipsData,
    getRegistrantsData,
    modelFilters,
    registrantFilters,
    view,
    page,
    rowsPerPage,
    searchText,
    resourcesDetail,
    checked,
  ]);

  const getRegistrants = async () => {
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

      for (let registrant of registrants) {
        let hostname = toLower(registrant?.hostname);
        const { data: modelRes } = await getMeshModelsData(
          {
            params: {
              page: page?.Models,
              pagesize: 'all',
              registrant: hostname,
              components: false,
              relationships: false,
            },
          },
          true,
        );
        if (modelRes.models && modelRes.models.length > 0) {
          const updatedRegistrant = {
            ...registrant,
            models: removeDuplicateVersions(modelRes.models) || [],
          };
          tempResourcesDetail.push(updatedRegistrant);
        }
      }
      response = {
        data: {
          registrants: tempResourcesDetail,
        },
      };
    }
    setRowsPerPage(25);
    return response;
  };
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
    setPage({
      Models: 0,
      Components: 0,
      Relationships: 0,
      Registrants: 0,
    });
    setShowDetailsData({
      type: '',
      data: {},
    });
  };
  const modifyData = () => {
    if (!resourcesDetail) return [];

    if (view === MODELS) {
      return removeDuplicateVersions(
        checked ? resourcesDetail.filter((model) => model.duplicates > 0) : resourcesDetail,
      );
    } else if (view === RELATIONSHIPS) {
      return groupRelationshipsByKind(resourcesDetail);
    } else if (view === REGISTRANTS) {
      return resourcesDetail || [];
    } else {
      return resourcesDetail;
    }
  };

  useEffect(() => {
    if (searchText !== null && page[view] > 0) {
      setPage({
        Models: 0,
        Components: 0,
        Relationships: 0,
        Registrants: 0,
      });
    }
  }, [searchText]);

  useEffect(() => {
    fetchData();
  }, [view, page, rowsPerPage, checked, searchText, modelFilters, registrantFilters]);

  // Update view when external view changes (for modal usage)
  useEffect(() => {
    if (externalView && externalView !== view) {
      setView(externalView);
      setResourcesDetail([]);
      setSearchText(externalSearchText || null);
      setModelsFilters({ page: 0 });
      setRegistrantsFilters({ page: 0 });
      setComponentsFilters({ page: 0 });
      setRelationshipsFilters({ page: 0 });
      setPage({
        Models: 0,
        Components: 0,
        Relationships: 0,
        Registrants: 0,
      });
      setShowDetailsData({
        type: '',
        data: {},
      });
    }
  }, [externalView, externalSearchText]);

  useEffect(() => {
    if (externalSearchText !== null && externalSearchText !== searchText) {
      setSearchText(externalSearchText);
    }
  }, [externalSearchText]);

  return (
    <div data-test="workloads">
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
              count={modelsData?.total_count || 0}
              active={view === MODELS}
              onClick={() => handleTabClick(MODELS)}
            />
            <TabCard
              label="Components"
              count={componentsData?.total_count || 0}
              active={view === COMPONENTS}
              onClick={() => handleTabClick(COMPONENTS)}
            />
            <TabCard
              label="Relationships"
              count={relationshipsData?.total_count || 0}
              active={view === RELATIONSHIPS}
              onClick={() => handleTabClick(RELATIONSHIPS)}
            />
            <TabCard
              label="Registrants"
              count={registrantsData?.total_count || 0}
              active={view === REGISTRANTS}
              onClick={() => handleTabClick(REGISTRANTS)}
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
              data={modifyData()}
              view={view}
              setSearchText={setSearchText}
              setPage={setPage}
              checked={checked}
              setChecked={setChecked}
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
            view={view}
            setShowDetailsData={setShowDetailsData}
            showDetailsData={showDetailsData}
          />
        </TreeWrapper>
      </MainContainer>
    </div>
  );
};

const TabBar = ({ openImportModal, openCreateModal, view, openRelationshipModal }) => {
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
        {`(${count})`}
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
