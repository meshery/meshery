import { withStyles } from '@material-ui/core';
import React, { useState, useEffect, useCallback } from 'react';
import { Paper, Button } from '@material-ui/core';
import UploadIcon from '@mui/icons-material/Upload';
import DoNotDisturbOnIcon from '@mui/icons-material/DoNotDisturbOn';
import { getUnit8ArrayDecodedFile } from '../../utils/utils';
import {
  OVERVIEW,
  MODELS,
  COMPONENTS,
  RELATIONSHIPS,
  REGISTRANTS,
  GRAFANA,
  PROMETHEUS,
} from '../../constants/navigator';
// import { SORT } from '../../constants/endpoints';
import useStyles from '../../assets/styles/general/tool.styles';
import MesheryTreeView from './MesheryTreeView';
import MeshModelDetails from './MeshModelDetails';
import { toLower } from 'lodash';
import { DisableButton } from './MeshModel.style';
import { useRouter } from 'next/router';
import { Provider } from 'react-redux';
import { store } from '../../store';
import { ErrorBoundary } from '../General/ErrorBoundary';
import {
  useLazyGetMeshModelsQuery,
  useLazyGetComponentsQuery,
  useLazyGetRelationshipsQuery,
  useLazyGetRegistrantsQuery,
  useImportMeshModelMutation,
} from '@/rtk-query/meshModel';
import NoSsr from '@material-ui/core/NoSsr';
import { groupRelationshipsByKind, removeDuplicateVersions } from './helper';
import _ from 'lodash';
import { Modal as SistentModal } from '@layer5/sistent';
import { UsesSistent } from '../SistentWrapper';
import { RJSFModalWrapper } from '../Modal';
import { useRef } from 'react';
import { updateProgress } from 'lib/store';

const meshmodelStyles = (theme) => ({
  wrapperClss: {
    flexGrow: 1,
    maxWidth: '100%',
    height: 'auto',
  },
  tab: {
    minWidth: 40,
    paddingLeft: 0,
    paddingRight: 0,
    '&.Mui-selected': {
      color: theme.palette.secondary.focused,
    },
  },
  tabs: {
    '& .MuiTabs-indicator': {
      backgroundColor: theme.palette.secondary.focused,
    },
  },
  dashboardSection: {
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
    overflowY: 'scroll',
  },
  duplicatesModelStyle: {
    backgroundColor: theme.palette.secondary.focused,
  },
});

const useMeshModelComponentRouter = () => {
  const router = useRouter();
  const { query } = router;

  const searchQuery = query.searchText || null;
  const selectedTab = query.tab === GRAFANA || query.tab === PROMETHEUS ? OVERVIEW : query.tab;
  const selectedPageSize = query.pagesize || 25;

  return { searchQuery, selectedTab, selectedPageSize };
};

const useInfiniteScrollRef = (callback) => {
  const observerRef = useRef(null);
  const triggerRef = useRef(null);

  useEffect(() => {
    // setTimeout gives the browser time to finish rendering the DOM elements before executing the callback function.
    const timeoutId = setTimeout(() => {
      if (!triggerRef.current) {
        return () => observerRef.current && observerRef.current.disconnect();
      }

      observerRef.current = new IntersectionObserver(
        (entries) => {
          entries.forEach((entry) => {
            if (entry.isIntersecting) {
              callback();
            }
          });
        },
        { threshold: 0.01 },
      );
      observerRef.current.observe(triggerRef.current);
    }, 0);

    return () => {
      observerRef.current && observerRef.current.disconnect();
      clearTimeout(timeoutId);
    };
  }, [callback, triggerRef.current]);

  return triggerRef;
};

const MeshModelComponent_ = ({
  modelsCount,
  componentsCount,
  relationshipsCount,
  registrantCount,
  settingsRouter,
}) => {
  const router = useRouter();
  const { handleChangeSelectedTab } = settingsRouter(router);
  const [resourcesDetail, setResourcesDetail] = useState([]);
  const { selectedTab, searchQuery, selectedPageSize } = useMeshModelComponentRouter();
  const [page, setPage] = useState({
    Models: 0,
    Components: 0,
    Relationships: 0,
    Registrants: 0,
  });
  const [searchText, setSearchText] = useState(searchQuery);
  const [rowsPerPage, setRowsPerPage] = useState(selectedPageSize);
  const StyleClass = useStyles();
  const [view, setView] = useState(OVERVIEW);
  const [convert, setConvert] = useState(false);
  const [importSchema, setImportSchema] = useState({});
  const [importModal, setImportModal] = useState({
    open: false,
  });
  const [showDetailsData, setShowDetailsData] = useState({
    type: '', // Type of selected data eg. (models, components)
    data: {},
  });
  const [animate, setAnimate] = useState(false);
  const [checked, setChecked] = useState(false);
  const [importModelReq] = useImportMeshModelMutation();

  const handleUploadImport = () => {
    setImportModal({
      open: true,
    });
  };

  const handleUploadImportClose = () => {
    setImportModal({
      open: false,
    });
  };

  const handleImportModel = async (data) => {
    const { uploadType, url, file, model } = data;
    let requestBody = null;

    const fileElement = document.getElementById('root_file');

    switch (uploadType) {
      case 'File Upload': {
        const fileName = fileElement.files[0].name;
        const fileData = getUnit8ArrayDecodedFile(file);
        if (fileData) {
          requestBody = {
            importBody: {
              model_file: fileData,
              url: '',
              filename: fileName,
            },
            uploadType: 'file',
            register: true,
          };
        } else {
          console.error('Error: File data is empty or invalid');
          return;
        }
        break;
      }
      case 'URL Import': {
        if (url) {
          requestBody = {
            importBody: {
              url: url,
              model: model,
            },
            uploadType: 'url',
            register: true,
          };
        } else {
          console.error('Error: URL is empty');
          return;
        }
        break;
      }
      default: {
        console.error('Error: Invalid upload type');
        return;
      }
    }
    updateProgress({ showProgress: true });
    await importModelReq({ importBody: requestBody });
    updateProgress({ showProgress: false });
  };
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

      if (response.data && response.data[view.toLowerCase()]) {
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
        const updatedRegistrant = {
          ...registrant,
          models: removeDuplicateVersions(modelRes.models) || [],
        };
        tempResourcesDetail.push(updatedRegistrant);
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
    if (!animate) {
      setAnimate(true);
      setConvert(true);
    }
  };

  const modifyData = () => {
    if (view === MODELS) {
      return removeDuplicateVersions(
        checked ? resourcesDetail.filter((model) => model.duplicates > 0) : resourcesDetail,
      );
    } else if (view === RELATIONSHIPS) {
      return groupRelationshipsByKind(resourcesDetail);
    } else {
      return resourcesDetail;
    }
  };
  useEffect(() => {
    if (selectedTab && selectedTab !== OVERVIEW) {
      setAnimate(true);
      setConvert(true);
      setView(selectedTab);
    }
  }, [selectedTab]);

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
  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch('/api/schema/resource/model', {
          method: 'GET',
          credentials: 'include',
        });
        const result = await response.json();
        setImportSchema(result);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <div data-test="workloads">
      <TabBar animate={animate} handleUploadImport={handleUploadImport} />
      {importModal.open && (
        <ImportModal
          importFormSchema={importSchema}
          handleClose={handleUploadImportClose}
          handleImportModel={handleImportModel}
        />
      )}
      <div
        className={`${StyleClass.mainContainer} ${animate ? StyleClass.mainContainerAnimate : ''}`}
      >
        <div
          className={`${StyleClass.innerContainer} ${
            animate ? StyleClass.innerContainerAnimate : ''
          }`}
        >
          <TabCard
            label="Models"
            count={modelsCount}
            active={view === MODELS && animate}
            animate={animate}
            onClick={() => handleTabClick(MODELS)}
          />
          <TabCard
            label="Components"
            count={componentsCount}
            active={view === COMPONENTS && animate}
            animate={animate}
            onClick={() => handleTabClick(COMPONENTS)}
          />
          <TabCard
            label="Relationships"
            count={relationshipsCount}
            active={view === RELATIONSHIPS && animate}
            animate={animate}
            onClick={() => handleTabClick(RELATIONSHIPS)}
          />
          <TabCard
            label="Registrants"
            count={registrantCount}
            active={view === REGISTRANTS && animate}
            animate={animate}
            onClick={() => handleTabClick(REGISTRANTS)}
          />
        </div>
        {convert && (
          <div
            className={`${StyleClass.treeWrapper} ${convert ? StyleClass.treeWrapperAnimate : ''}`}
          >
            <Paper
              className={StyleClass.detailsContainer}
              style={{
                display: 'flex',
                alignItems: resourcesDetail.length === 0 ? 'center' : '',
                justifyContent: resourcesDetail.length === 0 ? 'center' : '',
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
              />
            </Paper>
            <MeshModelDetails
              view={view}
              setShowDetailsData={setShowDetailsData}
              showDetailsData={showDetailsData}
            />
          </div>
        )}
      </div>
    </div>
  );
};
const ImportModal = React.memo((props) => {
  const { importFormSchema, handleClose, handleImportModel } = props;

  return (
    <>
      <UsesSistent>
        <SistentModal open={true} closeModal={handleClose} maxWidth="sm" title="Import Model">
          <RJSFModalWrapper
            schema={importFormSchema.rjsfSchema}
            uiSchema={{
              ...importFormSchema.uiSchema,
            }}
            handleSubmit={handleImportModel}
            submitBtnText="Import"
            handleClose={handleClose}
          />
        </SistentModal>
      </UsesSistent>
    </>
  );
});

ImportModal.displayName = 'ImportModal';

const TabBar = ({ animate, handleUploadImport }) => {
  const StyleClass = useStyles();

  return (
    <div
      className={`${StyleClass.meshModelToolbar} ${animate ? StyleClass.toolWrapperAnimate : ''}`}
    >
      <Button
        aria-label="Add Pattern"
        variant="contained"
        color="primary"
        size="large"
        onClick={handleUploadImport}
        style={{ display: 'flex', visibility: `${animate ? 'visible' : 'hidden'}` }}
        disabled={false} //TODO: Need to make key for this component
      >
        <UploadIcon />
        Import
      </Button>

      <DisableButton
        disabled
        variant="contained"
        size="large"
        style={{
          visibility: `${animate ? 'visible' : 'hidden'}`,
        }}
        startIcon={<DoNotDisturbOnIcon />}
      >
        Ignore
      </DisableButton>
    </div>
  );
};

const TabCard = ({ label, count, active, onClick, animate }) => {
  const StyleClass = useStyles();
  return (
    <Paper
      elevation={3}
      className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''} ${
        active ? StyleClass.activeTab : ''
      }`}
      onClick={onClick}
    >
      <span
        style={{
          fontWeight: `${animate ? 'normal' : 'bold'}`,
          fontSize: `${animate ? '1rem' : '3rem'}`,
          marginLeft: `${animate && '4px'}`,
        }}
      >
        {animate ? `(${count})` : `${count}`}
      </span>
      {label}
    </Paper>
  );
};

const MeshModelComponent = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary
        FallbackComponent={() => null}
        onError={(e) => console.error('Error in NotificationCenter', e)}
      >
        <Provider store={store}>
          <MeshModelComponent_ {...props} />
        </Provider>
      </ErrorBoundary>
    </NoSsr>
  );
};
export default withStyles(meshmodelStyles)(MeshModelComponent);
