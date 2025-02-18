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
} from '../../constants/navigator';
import DescriptionIcon from '@mui/icons-material/Description';
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
  Modal as SistentModal,
  PrimaryActionButtons,
  ModalFooter,
  useStepper,
  CustomizedStepper,
  ModalBody,
  Box,
  useModal,
  TextField,
  ModalButtonSecondary,
  ModalButtonPrimary,
  Grid,
  Select,
  InputLabel,
  FormControlLabel,
  Checkbox,
  ComponentIcon,
  Typography,
  Button,
  FormControl,
  RadioGroup,
  MenuItem,
  Radio,
  FormLabel,
  NoSsr,
} from '@layer5/sistent';
import BrushIcon from '@mui/icons-material/Brush';
import CategoryIcon from '@mui/icons-material/Category';
import SourceIcon from '@/assets/icons/SourceIcon';
import AppRegistrationIcon from '@mui/icons-material/AppRegistration';
import { UsesSistent } from '../SistentWrapper';
import { RJSFModalWrapper } from '../Modal';
import { useRef } from 'react';
import { updateProgress } from 'lib/store';
import ModelIcon from '@/assets/icons/ModelIcon';
import { iconMedium } from '../../css/icons.styles';
import AddIcon from '@mui/icons-material/AddCircleOutline';
import LanOutlinedIcon from '@mui/icons-material/LanOutlined';

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
  const [view, setView] = useState(OVERVIEW);
  const [convert, setConvert] = useState(false);
  const [importSchema, setImportSchema] = useState({});
  const [importModal, setImportModal] = useState({
    open: false,
  });
  const [generateModal, setGenerateModal] = useState({
    open: false,
  });
  const [showDetailsData, setShowDetailsData] = useState({
    type: '', // Type of selected data eg. (models, components)
    data: {},
  });
  const [animate, setAnimate] = useState(false);
  const [checked, setChecked] = useState(false);
  const [importModelReq] = useImportMeshModelMutation();
  const [uploadMethod, setUploadMethod] = useState('');

  const handleChangeUploadMethod = (e) => {
    setUploadMethod(e.target.value);
  };

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
    setGenerateModal({
      open: false,
    });
    urlModal.openModal({
      title: 'Create Model',
      reactNode: (
        <UrlStepper handleGenerateModal={handleGenerateModal} handleClose={urlModal.closeModal} />
      ),
    });
  };
  const handleCsvStepper = () => {
    setGenerateModal({
      open: false,
    });
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
      default: {
        console.error('Error: Invalid upload type');
        return;
      }
    }
    updateProgress({ showProgress: true });
    await importModelReq({ importBody: requestBody });
    updateProgress({ showProgress: false });
  };

  const handleGenerateModel = () => {
    setGenerateModal({
      open: true,
    });
  };
  const handleGenerateModelClose = () => {
    setGenerateModal({
      open: false,
    });
    setUploadMethod('');
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
      <TabBar
        animate={animate}
        handleUploadImport={handleUploadImport}
        handleGenerateModel={handleGenerateModel}
      />
      {importModal.open && (
        <ImportModal
          importFormSchema={importSchema}
          handleClose={handleUploadImportClose}
          handleImportModel={handleImportModel}
        />
      )}
      {generateModal.open && (
        <GenerateModal
          uploadMethod={uploadMethod}
          handleChange={handleChangeUploadMethod}
          handleClose={handleGenerateModelClose}
          handleUrlStepper={handleUrlStepper}
          handleCsvStepper={handleCsvStepper}
        />
      )}
      <UsesSistent>
        <SistentModal maxWidth="sm" {...urlModal}></SistentModal>
        <SistentModal maxWidth="sm" {...csvModal}></SistentModal>
      </UsesSistent>
      <MainContainer isAnimated={animate}>
        <InnerContainer isAnimated={animate}>
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
        </InnerContainer>
        {convert && (
          <TreeWrapper isAnimated={convert}>
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
        )}
      </MainContainer>
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

const GenerateModal = React.memo((props) => {
  const { handleClose, uploadMethod, handleChange, handleUrlStepper, handleCsvStepper } = props;

  return (
    <>
      <UsesSistent>
        <SistentModal open={true} closeModal={handleClose} maxWidth="sm" title="Generate Model">
          <FormControl style={{ padding: '10px' }}>
            <FormLabel id="upload-method-choices" sx={{ marginBottom: '1rem' }}>
              Upload Method
            </FormLabel>
            <RadioGroup
              aria-labelledby="upload-method-choices"
              name="uploadMethod"
              value={uploadMethod}
              onChange={handleChange}
            >
              <FormControlLabel
                value="url"
                control={<Radio color="primary" />}
                label="URL Import"
              />
              <FormControlLabel
                value="csv"
                control={<Radio color="primary" />}
                label="CSV Import"
              />
            </RadioGroup>
          </FormControl>
          <ModalFooter
            variant="filled"
            helpText="URL Import supports Artifacthub and Github. Csv Import supports bulk generation and import."
          >
            <PrimaryActionButtons
              primaryText="Next"
              secondaryText="Cancel"
              primaryButtonProps={{
                onClick: uploadMethod === 'url' ? handleUrlStepper : handleCsvStepper,
                disabled: !uploadMethod,
              }}
              secondaryButtonProps={{
                onClick: handleClose,
              }}
            />
          </ModalFooter>
        </SistentModal>
      </UsesSistent>
    </>
  );
});

GenerateModal.displayName = 'GenerateModal';
const UrlStepper = React.memo(({ handleGenerateModal, handleClose }) => {
  const [modelSource, setModelSource] = React.useState('');
  const [modelName, setModelName] = React.useState('');
  const [modelDisplayName, setModelDisplayName] = React.useState('');
  const [modelCategory, setModelCategory] = React.useState('');
  const [modelSubcategory, setModelSubcategory] = React.useState('');
  const [modelShape, setModelShape] = React.useState('');
  const [modelUrl, setModelUrl] = React.useState('');
  const [primaryColor, setPrimaryColor] = React.useState('#000000');
  const [secondaryColor, setSecondaryColor] = React.useState('#000000');
  const [logoLightThemePath, setLogoLightThemePath] = React.useState('');
  const [logoDarkThemePath, setLogoDarkThemePath] = React.useState('');
  const [registerModel, setRegisterModel] = React.useState(true);
  const [isAnnotation, setIsAnnotation] = React.useState(true);

  const modelCategories = [
    'Analytics',
    'App Definition and Development',
    'Cloud Native Network',
    'Cloud Native Storage',
    'Database',
    'Machine Learning',
    'Observability and Analysis',
    'Orchestration & Management',
    'Platform',
    'Provisioning',
    'Runtime',
    'Security & Compliance',
    'Serverless',
    'Tools',
    'Uncategorized',
  ];

  const modelSubCategories = [
    'Academic',
    'API Gateway',
    'Application Definition & Image Build',
    'Automation & Configuration',
    'Certified CNFs',
    'Certified Kubernetes - Distribution',
    'Certified Kubernetes - Hosted',
    'Certified Kubernetes - Installer',
    'Chaos Engineering',
    'Cloud Native Network',
    'Cloud Native Storage',
    'Container Registry',
    'Container Runtime',
    'Continuous Integration & Delivery',
    'Continuous Optimization',
    'Coordination & Service Discovery',
    'Database',
    'Debugging and Observability',
    'End User Supporter',
    'Framework',
    'Gold',
    'Hosted Platform',
    'Installable Platform',
    'Key Management',
    'Kubernetes Certified Service Provider',
    'Kubernetes Training Partner',
    'Logging',
    'Metrics',
    'Monitoring',
    'Nonprofit',
    'Packaging, Registries & Application Delivery',
    'PaaS/Container Service',
    'Platinum',
    'Remote Procedure Call',
    'Runtime',
    'Scheduling & Orchestration',
    'Security',
    'Security & Compliance',
    'Service Mesh',
    'Service Proxy',
    'Silver',
    'Specifications',
    'Streaming & Messaging',
    'Toolchain',
    'Tools',
    'Tracing',
  ];
  const handleLogoLightThemeChange = async (event) => {
    const file = event.target.files[0];

    if (file && file.type === 'image/svg+xml') {
      const reader = new FileReader();

      reader.onload = (e) => {
        const svgData = e.target.result;
        setLogoLightThemePath(svgData);
      };
      reader.readAsText(file);
    } else {
      console.error('Please upload a valid SVG file.');
    }
  };

  const handleLogoDarkThemeChange = async (event) => {
    const file = event.target.files[0];

    if (file && file.type === 'image/svg+xml') {
      const reader = new FileReader();

      reader.onload = (e) => {
        const svgData = e.target.result;
        setLogoDarkThemePath(svgData);
      };

      // Read the file as text (since it's an SVG)
      reader.readAsText(file);
    } else {
      console.error('Please upload a valid SVG file.');
    }
  };
  const modelShapes = [
    'rectangle',
    'round-rectangle',
    'bottom-round-rectangle',
    'cut-rectangle',
    'shape',
    'circle',
    'diamond',
    'round-rectang',
    'hexagon',
    'rhomboid',
    'triangle',
    'cilinder',
    'round-triangle',
    'round-pentagon',
    'sheild',
    'vee',
    'cylinder',
    'round-heptagon',
    'concave-hexagon',
    'right-rhomboid',
    'barrel',
    'round-diamond',
  ];
  const handleFinish = () => {
    handleClose();
    handleGenerateModal({
      uploadType: 'URL Import',
      register: registerModel,
      url: modelUrl,
      model: {
        model: modelName,
        modelDisplayName: modelDisplayName,
        registrant: modelSource,
        category: modelCategory,
        subCategory: modelSubcategory,
        shape: modelShape,
        primaryColor: primaryColor,
        secondaryColor: secondaryColor,
        svgColor: logoLightThemePath,
        svgWhite: logoDarkThemePath,
        isAnnotation: isAnnotation,
      },
    });
  };
  const urlStepper = useStepper({
    steps: [
      {
        component: (
          <div>
            <Box display="flex" alignItems="center" mb={2}>
              <Typography>
                Please select the appropriate <strong>Source</strong> based on your URL.
              </Typography>
            </Box>
            <FormControl component="fieldset">
              <RadioGroup
                row
                aria-label="source"
                name="source"
                value={modelSource}
                onChange={(e) => setModelSource(e.target.value.toLowerCase())}
                style={{ gap: '2rem' }}
              >
                {['Github', 'Artifacthub'].map((source, idx) => (
                  <FormControlLabel
                    key={idx}
                    value={source.toLowerCase()}
                    control={<Radio />}
                    label={<>{source}</>}
                  />
                ))}
              </RadioGroup>
            </FormControl>
            <FormControl fullWidth style={{ marginTop: '1rem' }}>
              <TextField
                required
                id="model-url"
                label="Model URL"
                value={modelUrl}
                onChange={(e) => setModelUrl(e.target.value)}
                variant="outlined"
              />
            </FormControl>
          </div>
        ),
        icon: SourceIcon,
        label: 'Source',
        helpText: (
          <>
            <ul>
              <li>
                <strong>GitHub:</strong> Provide a GitHub repository URL. For example,{' '}
                <em>git:://github.com/cert-manager/cert-manager/master/deploy/crds</em>.
              </li>
              <li>
                <br />
                <strong>ArtifactHub:</strong> ArtifactHub package URL. For example,{' '}
                <em>https://artifacthub.io/packages/search?ts_query_web={'{model-name}'}</em>.
              </li>
            </ul>
          </>
        ),
      },
      {
        component: (
          <div>
            <Box display="flex" alignItems="center" mb={2}>
              <Typography>
                Please enter the appropriate<strong> Model name</strong> and{' '}
                <strong> Model Display Name</strong> for your model.
              </Typography>
            </Box>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <FormControl fullWidth>
                  <TextField
                    required
                    id="model-name"
                    label="Model Name"
                    value={modelName}
                    onChange={(e) => setModelName(e.target.value)}
                    variant="outlined"
                  />
                </FormControl>
              </Grid>
              <Grid item xs={6}>
                <FormControl fullWidth>
                  <TextField
                    required
                    id="model-display-name"
                    label="Model Display Name"
                    value={modelDisplayName}
                    onChange={(e) => setModelDisplayName(e.target.value)}
                    variant="outlined"
                  />
                </FormControl>
              </Grid>
            </Grid>
          </div>
        ),
        icon: DescriptionIcon,
        label: 'Model Details',
        helpText: (
          <>
            <ul>
              <li>
                <strong>Model Name:</strong> Should be in lowercase with hyphens. For example,{' '}
                <em>cert-manager</em>.
              </li>
              <br />
              <li>
                <strong>Display Name:</strong> How you want your model to be named. For example,{' '}
                <em>Cert Manager</em>.
              </li>
            </ul>
          </>
        ),
      },
      {
        component: (
          <div>
            <Box display="flex" alignItems="center" mb={2}>
              <Typography>
                Please select the appropriate <strong>Category</strong> and
                <strong>Subcategory</strong> for your model.
                <br />
                <strong>Note:</strong> If you can&apos;t find the appropriate category or
                subcategory, please select <em>Uncategorized</em>.
              </Typography>
            </Box>

            <Grid container spacing={2}>
              <Grid item xs={6}>
                <FormControl fullWidth variant="outlined">
                  <InputLabel id="category-label">Category</InputLabel>
                  <Select
                    labelId="category-label"
                    id="category"
                    value={modelCategory}
                    label="Category"
                    onChange={(e) => setModelCategory(e.target.value)}
                  >
                    {modelCategories.map((category, idx) => (
                      <MenuItem key={idx} value={category}>
                        {category}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={6}>
                <FormControl fullWidth variant="outlined">
                  <InputLabel id="subcategory-label">Subcategory</InputLabel>
                  <Select
                    labelId="subcategory-label"
                    id="subcategory"
                    value={modelSubcategory}
                    label="Subcategory"
                    onChange={(e) => setModelSubcategory(e.target.value)}
                  >
                    {modelSubCategories.map((subCategory, idx) => (
                      <MenuItem key={idx} value={subCategory}>
                        {subCategory}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </div>
        ),
        icon: CategoryIcon,
        label: 'Model Categorization',
        helpText: (
          <>
            <ul>
              <li>
                <strong>Category:</strong> Determines the main grouping.
              </li>
              <li>
                <strong>Subcategory:</strong> Allows for more specific classification.
              </li>
            </ul>
          </>
        ),
      },
      {
        component: (
          <div>
            <Box display="flex" alignItems="center" mb="2rem">
              <Typography>
                Configure logos, colors, and shape for your model.
                <br />
                <strong>Note:</strong> If none of these are provided, default Meshery values will be
                used.
              </Typography>
            </Box>

            <Grid container spacing={1}>
              <Grid item xs={6}>
                <FormControl fullWidth variant="outlined">
                  <Typography>Logo (Dark Theme)</Typography>
                  <input
                    id="logo-dark-theme"
                    type="file"
                    accept=".svg"
                    onChange={handleLogoDarkThemeChange}
                    style={{ marginTop: '1rem' }}
                    label=" "
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6}>
                <FormControl fullWidth>
                  <Typography>Logo (Light Theme)</Typography>
                  <input
                    id="logo-light-theme"
                    type="file"
                    accept=".svg"
                    onChange={handleLogoLightThemeChange}
                    style={{ marginTop: '1rem' }}
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} style={{ marginTop: '2rem' }}>
                <FormControl fullWidth>
                  <Typography>Primary Color</Typography>
                  <input
                    id="primary-color"
                    type="color"
                    value={primaryColor}
                    onChange={(e) => setPrimaryColor(e.target.value)}
                    style={{ width: '100%', marginTop: '1rem' }}
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} style={{ marginTop: '2rem' }}>
                <FormControl fullWidth>
                  <Typography>Secondary Color</Typography>
                  <input
                    id="secondary-color"
                    type="color"
                    value={secondaryColor}
                    onChange={(e) => setSecondaryColor(e.target.value)}
                    style={{ width: '100%', marginTop: '1rem' }}
                  />
                </FormControl>
              </Grid>

              <Grid item xs={12} style={{ marginTop: '1rem' }}>
                <FormControl fullWidth variant="outlined">
                  <InputLabel id="shape-label">Shape</InputLabel>
                  <Select
                    labelId="shape-label"
                    id="shape"
                    value={modelShape}
                    label="Shape"
                    onChange={(e) => setModelShape(e.target.value)}
                    style={{ marginTop: '1rem' }}
                  >
                    {modelShapes.map((shape, idx) => (
                      <MenuItem key={idx} value={shape}>
                        {shape}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
            </Grid>
          </div>
        ),
        icon: BrushIcon,
        label: 'Styling',
        helpText: (
          <>
            <p>
              Configure your model&apos;s logos, primary and secondary colors, and shape. If none of
              these are provided, default Meshery values will be used.
            </p>
            <ul>
              <li>
                <strong>Primary Color:</strong> The main color used in your model&apos;s theme.
              </li>
              <br />
              <li>
                <strong>Secondary Color:</strong> The accent color used in your model&apos;s theme.
              </li>
              <br />
              <li>
                <strong>Shape:</strong> The shape used for visual elements in your model.
              </li>
            </ul>
          </>
        ),
      },
      {
        component: (
          <div>
            <Grid item xs={12}>
              <FormControl component="fieldset">
                <FormControlLabel
                  label="Would you like to register the model now so you can use it immediately after it's generated?"
                  labelPlacement="start"
                  control={
                    <Checkbox
                      checked={registerModel}
                      onChange={(e) => setRegisterModel(e.target.checked)}
                      name="registerModel"
                      color="primary"
                      st
                    />
                  }
                />
              </FormControl>
            </Grid>
            <Grid item xs={12} style={{ marginTop: '1rem' }}>
              <FormControl component="fieldset">
                <FormControlLabel
                  label="Is this model exclusively for visual annotation and not related to infrastructure management?"
                  labelPlacement="start"
                  control={
                    <Checkbox
                      checked={isAnnotation}
                      onChange={(e) => setIsAnnotation(e.target.checked)}
                      name="registerModel"
                      color="primary"
                    />
                  }
                />
              </FormControl>
            </Grid>
          </div>
        ),
        icon: AppRegistrationIcon,
        label: 'Additional Details',
        helpText: (
          <>
            <p>Specify your preferences for model registration and usage:</p>
            <ul>
              <li>
                <strong>Register Model Now</strong>: Choose this option to register the model
                immediately after it&apos;s generated, allowing you to use it right away.
              </li>
              <br />
              <li>
                <strong>Visual Annotation Only</strong>: Select this if the model is exclusively for
                visual annotation purposes and not related to infrastructure management.
              </li>
            </ul>
          </>
        ),
      },
    ],
  });

  const transitionConfig = {
    0: {
      canGoNext: () => modelUrl && modelSource,
      nextButtonText: 'Next',
      nextAction: () => urlStepper.handleNext(),
    },
    1: {
      canGoNext: () => modelDisplayName && modelName,
      nextButtonText: 'Next',
      nextAction: () => urlStepper.handleNext(),
    },
    2: {
      canGoNext: () => true,
      nextButtonText: 'Next',
      nextAction: () => urlStepper.handleNext(),
    },
    3: {
      canGoNext: () => true,
      nextButtonText: 'Next',
      nextAction: () => urlStepper.handleNext(),
    },
    4: {
      canGoNext: () => true,
      nextButtonText: 'Finish',
      nextAction: handleFinish,
    },
  };

  const canGoNext = transitionConfig[urlStepper.activeStep].canGoNext();
  const nextButtonText = transitionConfig[urlStepper.activeStep].nextButtonText;

  return (
    <>
      <ModalBody>
        <CustomizedStepper {...urlStepper}>{urlStepper.activeStepComponent}</CustomizedStepper>
      </ModalBody>
      <ModalFooter
        variant="filled"
        helpText={urlStepper.steps[urlStepper.activeStep]?.helpText || ``}
      >
        <Box style={{ width: '100%', display: 'flex', gap: '1rem', justifyContent: 'end' }}>
          <ModalButtonSecondary onClick={urlStepper.goBack} disabled={!urlStepper.canGoBack}>
            Back
          </ModalButtonSecondary>
          <ModalButtonPrimary
            disabled={!canGoNext}
            onClick={transitionConfig[urlStepper.activeStep].nextAction}
          >
            {nextButtonText}
          </ModalButtonPrimary>
        </Box>
      </ModalFooter>
    </>
  );
});

const CsvStepper = React.memo(({ handleClose, handleGenerateModal }) => {
  const [modelCsvFile, setModelCsvFile] = React.useState(null);
  const [componentCsvFile, setComponentCsvFile] = React.useState(null);
  const [relationshipCsvFile, setRelationshipCsvFile] = React.useState(null);
  const [registerModel, setRegisterModel] = React.useState(true);
  const fileToBase64 = (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = (error) => reject(error);
    });
  };
  const handleModelCsvFileChange = async (event) => {
    const file = event.target.files[0];
    if (file) {
      const base64 = await fileToBase64(file);
      setModelCsvFile(base64);
    }
  };
  const handleComponentCsvFileChange = async (event) => {
    const file = event.target.files[0];
    if (file) {
      const base64 = await fileToBase64(file);
      setComponentCsvFile(base64);
    }
  };

  const handleRelationshipCsvFileChange = async (event) => {
    const file = event.target.files[0];
    if (file) {
      const base64 = await fileToBase64(file);
      setRelationshipCsvFile(base64);
    }
  };

  const handleFinish = () => {
    handleClose();
    handleGenerateModal({
      model_csv: `data:text/csv;base64,${modelCsvFile?.split(',')[1]}`,
      component_csv: `data:text/csv;base64,${componentCsvFile?.split(',')[1]}`,
      relationship_csv: relationshipCsvFile
        ? `data:text/csv;base64,${relationshipCsvFile?.split(',')[1]}`
        : null,
      register: registerModel,
      uploadType: 'CSV Import',
    });
  };

  const csvStepper = useStepper({
    steps: [
      {
        component: (
          <div>
            <Typography>
              Please upload the<strong> model csv</strong>
              <br />
              <strong>Note:</strong> Click on info button to look at csv template.
            </Typography>
            <FormControl fullWidth>
              {modelCsvFile ? (
                <Box>
                  <p>Model CSV file uploaded</p>
                </Box>
              ) : (
                <>
                  <input
                    required
                    id="model-csv-file"
                    type="file"
                    accept=".csv"
                    onChange={handleModelCsvFileChange}
                    style={{ marginTop: '1rem' }}
                  />
                </>
              )}
            </FormControl>
          </div>
        ),
        icon: ModelIcon,
        label: 'Model CSV',
        helpText: (
          <>
            Upload the Model CSV file.
            <a
              href="https://github.com/meshery/meshery/blob/a514f8689260791077bde8171646933cff15dd08/mesheryctl/templates/template-csvs/Models.csv"
              target="_blank"
              rel="noopener noreferrer"
              style={{ color: 'inherit', textDecoration: 'underline' }}
            >
              Click here
            </a>{' '}
            to see the CSV template.
          </>
        ),
      },
      {
        component: (
          <div>
            <Typography>
              Please upload the<strong> component csv</strong>
              <br />
              <strong>Note:</strong> Click on info button to look at csv template.
            </Typography>
            <FormControl fullWidth>
              {componentCsvFile ? (
                <Box>
                  <p>Component CSV file uploaded</p>
                </Box>
              ) : (
                <>
                  <input
                    required
                    id="component-csv-file"
                    type="file"
                    accept=".csv"
                    onChange={handleComponentCsvFileChange}
                    style={{ marginTop: '1rem' }}
                  />
                </>
              )}
            </FormControl>
          </div>
        ),
        icon: ComponentIcon,
        label: 'Component CSV',
        helpText: (
          <>
            Upload the Component CSV file.
            <a
              href="https://github.com/meshery/meshery/blob/a514f8689260791077bde8171646933cff15dd08/mesheryctl/templates/template-csvs/Components.csv"
              target="_blank"
              rel="noopener noreferrer"
              style={{ color: 'inherit', textDecoration: 'underline' }}
            >
              Click here
            </a>{' '}
            to see the CSV template.
          </>
        ),
      },
      {
        component: (
          <div>
            <Typography>
              Please upload the<strong> relationship csv</strong>.
              <br />
              <strong>Note:</strong> Click on info button to look at csv template.
            </Typography>
            <FormControl fullWidth>
              {relationshipCsvFile ? (
                <Box>
                  <p>Relationship CSV file uploaded</p>
                </Box>
              ) : (
                <>
                  <input
                    id="relationship-csv-file"
                    type="file"
                    accept=".csv"
                    onChange={handleRelationshipCsvFileChange}
                    style={{ marginTop: '1rem' }}
                  />
                </>
              )}
            </FormControl>
          </div>
        ),
        icon: LanOutlinedIcon,
        label: 'Relationship CSV',
        helpText: (
          <>
            Upload the relationship CSV file.
            <a
              href="https://github.com/meshery/meshery/blob/a514f8689260791077bde8171646933cff15dd08/mesheryctl/templates/template-csvs/Relationships.csv"
              target="_blank"
              rel="noopener noreferrer"
              style={{ color: 'inherit', textDecoration: 'underline' }}
            >
              Click here
            </a>{' '}
            to see the CSV template.
          </>
        ),
      },
      {
        component: (
          <div>
            <FormControl component="fieldset" marginTop={'1rem'}>
              <FormControlLabel
                labelPlacement="start"
                control={
                  <Checkbox
                    checked={registerModel}
                    onChange={(e) => setRegisterModel(e.target.checked)}
                    name="registerModel"
                    color="primary"
                  />
                }
                label="Would you like to register the model now so you can use it immediately after it's generated?"
              />
            </FormControl>
          </div>
        ),
        icon: AppRegistrationIcon,
        label: 'Register Model',
        helpText: 'Choose whether to register the model.',
      },
    ],
  });

  const transitionConfig = {
    0: {
      canGoNext: () => modelCsvFile !== null,
      nextButtonText: 'Next',
      nextAction: () => csvStepper.handleNext(),
    },
    1: {
      canGoNext: () => componentCsvFile !== null,
      nextButtonText: 'Next',
      nextAction: () => csvStepper.handleNext(),
    },
    2: {
      canGoNext: () => true,
      nextButtonText: 'Next',
      nextAction: () => csvStepper.handleNext(),
    },
    3: {
      canGoNext: () => true,
      nextButtonText: 'Finish',
      nextAction: handleFinish,
    },
  };

  const canGoNext = transitionConfig[csvStepper.activeStep].canGoNext();
  const nextButtonText = transitionConfig[csvStepper.activeStep].nextButtonText;

  return (
    <>
      <ModalBody>
        <CustomizedStepper {...csvStepper}>{csvStepper.activeStepComponent}</CustomizedStepper>
      </ModalBody>
      <ModalFooter
        variant="filled"
        helpText={csvStepper.steps[csvStepper.activeStep]?.helpText || ''}
      >
        <Box style={{ width: '100%', display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
          <ModalButtonSecondary onClick={csvStepper.goBack} disabled={!csvStepper.canGoBack}>
            Back
          </ModalButtonSecondary>
          <ModalButtonPrimary
            disabled={!canGoNext}
            onClick={transitionConfig[csvStepper.activeStep].nextAction}
          >
            {nextButtonText}
          </ModalButtonPrimary>
        </Box>
      </ModalFooter>
    </>
  );
});

CsvStepper.displayName = 'CsvStepper';

UrlStepper.displayName = 'Create';

const TabBar = ({ animate, handleUploadImport, handleGenerateModel }) => {
  return (
    <MeshModelToolbar isAnimated={animate}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          gap: '1rem', // Add some space between buttons
          visibility: `${animate ? 'visible' : 'hidden'}`,
        }}
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

        <Button
          aria-label="Add Pattern"
          variant="contained"
          color="primary"
          size="large"
          onClick={handleGenerateModel}
          style={{ display: 'flex', visibility: `${animate ? 'visible' : 'hidden'}` }}
          disabled={false} //TODO: Need to make key for this component
        >
          <AddIcon style={iconMedium} />
          &nbsp; Generate
        </Button>
      </div>
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
    </MeshModelToolbar>
  );
};

const TabCard = ({ label, count, active, onClick, animate }) => {
  return (
    <CardStyle isAnimated={animate} isSelected={active} elevation={3} onClick={onClick}>
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
