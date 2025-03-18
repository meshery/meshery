import React, { useState, useEffect, useCallback } from 'react';
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
  REGISTRY,
} from '../../constants/navigator';
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
  useImportMeshModelMutation,
} from '@/rtk-query/meshModel';
import { groupRelationshipsByKind, removeDuplicateVersions } from './helper';
import _ from 'lodash';
import {
  Modal,
  useModal,
  FormControlLabel,
  Button,
  FormControl,
  RadioGroup,
  Radio,
  NoSsr,
  importModelUiSchema,
  importModelSchema,
  Typography,
} from '@layer5/sistent';
import { RJSFModalWrapper } from '../Modal';
import { useRef } from 'react';
import { updateProgress } from 'lib/store';
import { iconSmall } from '../../css/icons.styles';
import AddIcon from '@mui/icons-material/AddCircleOutline';
import CsvStepper, { StyledDocsRedirectLink } from './Stepper/CSVStepper';
import UrlStepper from './Stepper/UrlStepper';
import { MESHERY_DOCS_URL } from '@/constants/endpoints';

const useMeshModelComponentRouter = () => {
  const router = useRouter();
  const { query } = router;

  if (query.settingsCategory === REGISTRY && !query.tab) {
    router.push({
      pathname: router.pathname,
      query: {
        ...query,
        tab: MODELS,
      },
    });
  }
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
  const { searchQuery, selectedPageSize } = useMeshModelComponentRouter();
  const [page, setPage] = useState({
    Models: 0,
    Components: 0,
    Relationships: 0,
    Registrants: 0,
  });
  const [searchText, setSearchText] = useState(searchQuery);
  const [rowsPerPage, setRowsPerPage] = useState(selectedPageSize);
  const [view, setView] = useState(MODELS);
  const [importModal, setImportModal] = useState({
    open: false,
  });
  const [showDetailsData, setShowDetailsData] = useState({
    type: '', // Type of selected data eg. (models, components)
    data: {},
  });
  const [checked, setChecked] = useState(false);
  const [importModelReq] = useImportMeshModelMutation();

  const handleGenerateModal = async (data) => {
    const { uploadType, url, model, component_csv, model_csv, relationship_csv, register } = data;
    let requestBody = null;
    switch (uploadType) {
      case 'CSV Import': {
        requestBody = {
          importBody: {
            model_csv: model_csv,
            component_csv: component_csv,
            relationship_csv: relationship_csv,
          },
          uploadType: 'csv',
          register: register,
        };
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
            register: register,
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

  const urlModal = useModal({});

  const csvModal = useModal({});

  const handleUrlStepper = () => {
    urlModal.openModal({
      title: 'Create Model',
      reactNode: (
        <UrlStepper handleGenerateModal={handleGenerateModal} handleClose={urlModal.closeModal} />
      ),
    });
  };
  const handleCsvStepper = () => {
    csvModal.openModal({
      title: 'Import CSV ',
      reactNode: (
        <CsvStepper handleGenerateModal={handleGenerateModal} handleClose={csvModal.closeModal} />
      ),
    });
  };

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
    const { uploadType, url, file } = data;
    let requestBody = null;

    const fileElement = document.getElementById('root_file');

    switch (uploadType) {
      case 'File Import': {
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
            },
            uploadType: 'urlImport',
            register: true,
          };
        } else {
          console.error('Error: URL is empty');
          return;
        }
        break;
      }
      case 'CSV Import': {
        handleCsvStepper();
        return;
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

  return (
    <div data-test="workloads">
      <TabBar handleUploadImport={handleUploadImport} handleUrlStepper={handleUrlStepper} />
      {importModal.open && (
        <ImportModal handleClose={handleUploadImportClose} handleImportModel={handleImportModel} />
      )}
      <>
        <Modal maxWidth="sm" {...urlModal}></Modal>
        <Modal maxWidth="sm" {...csvModal}></Modal>
      </>
      <MainContainer>
        <InnerContainer>
          <TabCard
            label="Models"
            count={modelsCount}
            active={view === MODELS}
            onClick={() => handleTabClick(MODELS)}
          />
          <TabCard
            label="Components"
            count={componentsCount}
            active={view === COMPONENTS}
            onClick={() => handleTabClick(COMPONENTS)}
          />
          <TabCard
            label="Relationships"
            count={relationshipsCount}
            active={view === RELATIONSHIPS}
            onClick={() => handleTabClick(RELATIONSHIPS)}
          />
          <TabCard
            label="Registrants"
            count={registrantCount}
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

const ImportModal = React.memo((props) => {
  const [importModalDescription, setImportModalDescription] = useState('');

  const CustomRadioWidget = (props) => {
    const { options, value, onChange, label, schema } = props;
    const { enumOptions } = options;

    setImportModalDescription(schema.description);

    return (
      <FormControl component="fieldset">
        <RadioGroup
          value={value}
          onChange={(e) => onChange(e.target.value)}
          style={{ marginTop: '-1.7rem', marginLeft: '-1rem' }}
        >
          <Typography fontWeight={'bold'} fontSize={'1rem'}>
            {label}
          </Typography>

          {enumOptions.map((option, index) => (
            <FormControlLabel
              key={option.value}
              value={option.value}
              control={<Radio />}
              label={
                <div>
                  <Typography variant="subtitle1">{option.label}</Typography>
                  <Typography variant="body2" color="textSecondary" textTransform={'none'}>
                    {schema.enumDescriptions[index]}
                  </Typography>
                </div>
              }
            />
          ))}
        </RadioGroup>
      </FormControl>
    );
  };

  const widgets = {
    RadioWidget: CustomRadioWidget,
  };
  const { handleClose, handleImportModel } = props;

  return (
    <Modal open={true} closeModal={handleClose} maxWidth="sm" title="Import Model">
      <RJSFModalWrapper
        schema={importModelSchema}
        uiSchema={importModelUiSchema}
        handleSubmit={handleImportModel}
        submitBtnText="Import"
        handleClose={handleClose}
        widgets={widgets}
        helpText={
          <p>
            {importModalDescription} <br />
            Learn more about importing Models in our{' '}
            <StyledDocsRedirectLink
              href={`${MESHERY_DOCS_URL}/guides/configuration-management/importing-models`}
              target="_blank"
              rel="noopener noreferrer"
            >
              documentation
            </StyledDocsRedirectLink>
            .
          </p>
        }
      />
    </Modal>
  );
});

ImportModal.displayName = 'ImportModal';

const TabBar = ({ handleUploadImport, handleUrlStepper }) => {
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
          onClick={handleUrlStepper}
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
          onClick={handleUploadImport}
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
