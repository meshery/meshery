import React, { useState, useEffect, useCallback } from 'react';
import UploadIcon from '@mui/icons-material/Upload';
import AddIcon from '@mui/icons-material/AddCircleOutline';
import LinkIcon from '@mui/icons-material/Link';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../../constants/navigator';
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
  modelsCount: initialModelsCount,
  componentsCount: initialComponentsCount,
  relationshipsCount: initialRelationshipsCount,
  registrantCount: initialRegistrantCount,
  settingsRouter,
}) => {
  const router = useRouter();
  const { handleChangeSelectedTab, selectedTab } = settingsRouter(router);
  const [resourcesDetail, setResourcesDetail] = useState([]);
  const { searchQuery, selectedPageSize } = useMeshModelComponentRouter();
  const [page, setPage] = useState({
    Models: 0,
    Components: 0,
    Relationships: 0,
    Registrants: 0,
  });

  const [counts, setCounts] = useState({
    models: initialModelsCount,
    components: initialComponentsCount,
    relationships: initialRelationshipsCount,
    registrants: initialRegistrantCount,
  });

  const [searchText, setSearchText] = useState(searchQuery);
  const [rowsPerPage, setRowsPerPage] = useState(selectedPageSize);
  const [view, setView] = useState(selectedTab ?? 'Models');
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
          if (response.data && response.data.total_count !== undefined) {
            setCounts((prevCounts) => ({
              ...prevCounts,
              models: response.data.total_count,
            }));
          }
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
          if (response.data && response.data.total_count !== undefined) {
            setCounts((prevCounts) => ({
              ...prevCounts,
              components: response.data.total_count,
            }));
          }
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
          if (response.data && response.data.total_count !== undefined) {
            setCounts((prevCounts) => ({
              ...prevCounts,
              relationships: response.data.total_count,
            }));
          }
          break;
        case REGISTRANTS:
          response = await getRegistrants();
          if (response?.data?.registrants) {
            setCounts((prevCounts) => ({
              ...prevCounts,
              registrants: response.data.registrants.length,
            }));
          }
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

  useEffect(() => {
    const fetchAllCounts = async () => {
      try {
        const modelsResponse = await getMeshModelsData(
          {
            params: {
              page: 0,
              pagesize: 1,
              components: false,
              relationships: false,
            },
          },
          true,
        );
        if (modelsResponse.data && modelsResponse.data.total_count !== undefined) {
          setCounts((prevCounts) => ({
            ...prevCounts,
            models: modelsResponse.data.total_count,
          }));
        }

        const componentsResponse = await getComponentsData(
          {
            params: {
              page: 0,
              pagesize: 1,
              trim: true,
            },
          },
          true,
        );
        if (componentsResponse.data && componentsResponse.data.total_count !== undefined) {
          setCounts((prevCounts) => ({
            ...prevCounts,
            components: componentsResponse.data.total_count,
          }));
        }

        const relationshipsResponse = await getRelationshipsData(
          {
            params: {
              page: 0,
              pagesize: 1, 
            },
          },
          true,
        );
        if (relationshipsResponse.data && relationshipsResponse.data.total_count !== undefined) {
          setCounts((prevCounts) => ({
            ...prevCounts,
            relationships: relationshipsResponse.data.total_count,
          }));
        }

        const registrantsResponse = await getRegistrantsData(
          {
            params: {
              page: 0,
              pagesize: 1,
            },
          },
          true,
        );
        if (registrantsResponse.data && registrantsResponse.data.total_count !== undefined) {
          setCounts((prevCounts) => ({
            ...prevCounts,
            registrants: registrantsResponse.data.total_count,
          }));
        }
      } catch (error) {
        console.error('Failed to fetch counts:', error);
      }
    };

    fetchAllCounts();
  }, []);

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
    handleChangeSelectedTab(selectedView);
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
    setView(selectedView);
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

  return (
    <div data-test="workloads">
      <TabBar
        openImportModal={() => setIsImportModalOpen(true)}
        openCreateModal={() => setIsCreateModalOpen(true)}
        openRelationshipModal={() => setIsRelationshipModalOpen(true)}
        view={view}
      />

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
        <InnerContainer>
          <TabCard
            label="Models"
            count={counts.models}
            active={view === MODELS}
            onClick={() => handleTabClick(MODELS)}
          />
          <TabCard
            label="Components"
            count={counts.components}
            active={view === COMPONENTS}
            onClick={() => handleTabClick(COMPONENTS)}
          />
          <TabCard
            label="Relationships"
            count={counts.relationships}
            active={view === RELATIONSHIPS}
            onClick={() => handleTabClick(RELATIONSHIPS)}
          />
          <TabCard
            label="Registrants"
            count={counts.registrants}
            active={view === REGISTRANTS}
            onClick={() => handleTabClick(REGISTRANTS)}
          />
        </InnerContainer>

        <TreeWrapper>
          <DetailsContainer
            isEmpty={!resourcesDetail.length}
            style={{
              padding: '0.6rem',
              overflow: 'hidden',
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
          Create
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
          Import
        </Button>
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
