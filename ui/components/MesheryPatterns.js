/* eslint-disable react/display-name */

import {
  CustomColumnVisibilityControl,
  CustomTooltip,
  OutlinedPatternIcon,
  SearchBar,
  UniversalFilter,
  importDesignSchema,
  importDesignUiSchema,
  publishCatalogItemSchema,
  publishCatalogItemUiSchema,
  Box,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  IconButton,
  ResponsiveDataTable,
  Typography,
  styled,
  PROMPT_VARIANTS,
} from '@layer5/sistent';
import { NoSsr } from '@layer5/sistent';
import CloseIcon from '@mui/icons-material/Close';
import DeleteIcon from '@mui/icons-material/Delete';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import GetAppIcon from '@mui/icons-material/GetApp';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import SaveIcon from '@mui/icons-material/Save';
import CustomToolbarSelect from './MesheryPatterns/CustomToolbarSelect';
import AddIcon from '@mui/icons-material/AddCircleOutline';
import React, { useEffect, useRef, useState } from 'react';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import Moment from 'react-moment';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { toggleCatalogContent, updateProgress } from '../lib/store';
import { encodeDesignFile, getUnit8ArrayDecodedFile, parseDesignFile } from '../utils/utils';
import ViewSwitch from './ViewSwitch';
import MesheryPatternGrid from './MesheryPatterns/MesheryPatternGridView';
import UndeployIcon from '../public/static/img/UndeployIcon';
import DoneAllIcon from '@mui/icons-material/DoneAll';
import PublicIcon from '@mui/icons-material/Public';
import PublishIcon from '@mui/icons-material/Publish';
import _PromptComponent from './PromptComponent';
import LoadingScreen from './LoadingComponents/LoadingComponent';
import { FILE_OPS, MesheryPatternsCatalog, VISIBILITY } from '../utils/Enum';
import CloneIcon from '../public/static/img/CloneIcon';
import { useRouter } from 'next/router';
import { RJSFModalWrapper } from './Modal';
import downloadContent from '../utils/fileDownloader';
import ConfigurationSubscription from './graphql/subscriptions/ConfigurationSubscription';
import Pattern from '../public/static/img/drawer-icons/pattern_svg.js';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import _ from 'lodash';
import { getMeshModels } from '../api/meshmodel';
import { modifyRJSFSchema } from '../utils/utils';
import { Edit as EditIcon } from '@mui/icons-material';
import { updateVisibleColumns } from '../utils/responsive-column';
import { useWindowDimensions } from '../utils/dimension';
import InfoModal from './Modals/Information/InfoModal';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import { DefaultTableCell, SortableTableCell } from './connections/common/index.js';
import DefaultError from './General/error-404/index';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import ExportModal from './ExportModal';
import { useModal, Modal as SistentModal, ModalBody } from '@layer5/sistent';
import PatternIcon from '@/assets/icons/Pattern';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import { useActorRef } from '@xstate/react';
import { designValidationMachine } from 'machines/validator/designValidator';
import { UnDeployStepper, DeployStepper } from './DesignLifeCycle/DeployStepper';
import { DryRunDesign } from './DesignLifeCycle/DryRun';
import { DEPLOYMENT_TYPE } from './DesignLifeCycle/common';
import {
  useClonePatternMutation,
  useDeletePatternFileMutation,
  useDeletePatternMutation,
  useDeployPatternMutation,
  useGetPatternsQuery,
  useImportPatternMutation,
  usePublishPatternMutation,
  useUndeployPatternMutation,
  useUnpublishPatternMutation,
  useUpdatePatternFileMutation,
  useUploadPatternFileMutation,
} from '@/rtk-query/design';
import CheckIcon from '@/assets/icons/CheckIcon';
import { ValidateDesign } from './DesignLifeCycle/ValidateDesign';
import PatternConfigureIcon from '@/assets/icons/PatternConfigure';
// import { useGetUserPrefQuery } from '@/rtk-query/user';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import TooltipButton from '@/utils/TooltipButton';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import yaml from 'js-yaml';
import ActionPopover from './MesheryPatterns/ActionPopover';

const genericClickHandler = (ev, fn) => {
  ev.stopPropagation();
  fn(ev);
};

const ViewSwitchButton = styled(Box)(() => ({
  justifySelf: 'flex-end',
}));

const CreateButton = styled(Box)(() => ({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  whiteSpace: 'nowrap',
}));

const AddIconStyled = styled(AddIcon)(() => ({
  paddingRight: '.35rem',
}));

const SearchWrapper = styled(Box)(() => ({
  justifySelf: 'flex-end',
  marginLeft: 'auto',
  paddingLeft: '1rem',
  display: 'flex',
  '@media (max-width: 965px)': {
    width: 'max-content',
  },
}));

const BtnText = styled('span')(() => ({
  display: 'block',
  '@media (max-width: 765px)': {
    display: 'none',
  },
}));

const EllipsisButtonWrapper = styled('div')(({ theme }) => ({
  display: 'none',
  [theme.breakpoints.down('lg')]: {
    display: 'block',
  },
}));

const ActionWrapper = styled('div')(({ theme }) => ({
  display: 'block',
  [theme.breakpoints.down('lg')]: {
    display: 'none',
  },
}));

const YamlDialogTitle = styled(DialogTitle)(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'end',
}));

const YamlDialogTitleText = styled(Typography)(() => ({
  flexGrow: 1,
}));

function TooltipIcon({ children, onClick, title, placement, disabled }) {
  return (
    <>
      <CustomTooltip title={title} placement={placement} interactive>
        <div>
          <IconButton disabled={disabled} onClick={onClick}>
            {children}
          </IconButton>
        </div>
      </CustomTooltip>
    </>
  );
}

function YAMLEditor({ pattern, onClose, onSubmit, isReadOnly = false }) {
  const [yaml, setYaml] = useState(pattern.pattern_file);
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const FullScreenCodeMirrorWrapper = styled('div')(() => ({
    height: '100%',
    '& .CodeMirror': {
      minHeight: '300px',
      height: fullScreen ? '80vh' : '30vh',
    },
  }));

  return (
    <Dialog
      onClose={onClose}
      aria-labelledby="pattern-dialog-title"
      open
      maxWidth="md"
      fullScreen={fullScreen}
      fullWidth={!fullScreen}
    >
      <YamlDialogTitle
        style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'end' }}
        disableTypography
        id="pattern-dialog-title"
      >
        <div>
          <YamlDialogTitleText variant="h6">{pattern.name}</YamlDialogTitleText>
        </div>
        <div>
          <CustomTooltip
            placement="top"
            title={fullScreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
            onClick={toggleFullScreen}
          >
            {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
          </CustomTooltip>
          <CustomTooltip placement="top" title="Exit" onClick={onClose}>
            <CloseIcon />
          </CustomTooltip>
        </div>
      </YamlDialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <FullScreenCodeMirrorWrapper>
          <CodeMirror
            value={pattern.pattern_file}
            options={{
              theme: 'material',
              lineNumbers: true,
              lineWrapping: true,
              gutters: ['CodeMirror-lint-markers'],
              // @ts-ignore
              lint: true,
              mode: 'text/x-yaml',
              readOnly: isReadOnly,
            }}
            onChange={(_, data, val) => setYaml(val)}
          />
        </FullScreenCodeMirrorWrapper>
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        {isReadOnly ? null : (
          <>
            <CustomTooltip title="Update Pattern">
              <IconButton
                aria-label="Update"
                disabled={!CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject)}
                onClick={() =>
                  onSubmit({
                    data: yaml,
                    id: pattern.id,
                    name: pattern.name,
                    type: FILE_OPS.UPDATE,
                    catalog_data: pattern.catalog_data,
                  })
                }
              >
                <SaveIcon />
              </IconButton>
            </CustomTooltip>
            <CustomTooltip title="Delete Pattern">
              <IconButton
                aria-label="Delete"
                disabled={!CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject)}
                onClick={() =>
                  onSubmit({
                    data: yaml,
                    id: pattern.id,
                    name: pattern.name,
                    type: FILE_OPS.DELETE,
                    catalog_data: pattern.catalog_data,
                  })
                }
              >
                <DeleteIcon />
              </IconButton>
            </CustomTooltip>
          </>
        )}
      </DialogActions>
    </Dialog>
  );
}

function resetSelectedPattern() {
  return { show: false, pattern: null };
}

function MesheryPatterns({
  updateProgress,
  user,
  selectedK8sContexts,
  catalogVisibility,
  disableCreateImportDesignButton = false,
  disableUniversalFilter = false,
  hideVisibility = false,
  initialFilters = { visibility: 'All' },
  pageTitle = 'Designs',
  arePatternsReadOnly = false,
}) {
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('updated_at desc');
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const modalRef = useRef();
  const [patterns, setPatterns] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [selectedPattern, setSelectedPattern] = useState(resetSelectedPattern());
  const router = useRouter();
  const [meshModels, setMeshModels] = useState([]);
  const [selectedFilters, setSelectedFilters] = useState(initialFilters);
  const [canPublishPattern, setCanPublishPattern] = useState(false);
  const [publishSchema, setPublishSchema] = useState({});
  const [infoModal, setInfoModal] = useState({
    open: false,
    ownerID: '',
    selectedResource: {},
  });
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [viewType, setViewType] = useState('grid');
  const { notify } = useNotification();
  const [visibilityFilter, setVisibilityFilter] = useState(null);

  const [deployPatternMutation] = useDeployPatternMutation();
  const [undeployPatternMutation] = useUndeployPatternMutation();
  const {
    data: patternsData,
    isLoading: ispatternsLoading,
    refetch: getPatterns,
  } = useGetPatternsQuery({
    page: page,
    pagesize: pageSize,
    search: search,
    order: sortOrder,
    visibility: visibilityFilter ? JSON.stringify([visibilityFilter]) : '',
    populate: 'pattern_file',
  });
  const [clonePattern] = useClonePatternMutation();
  const [publishCatalog] = usePublishPatternMutation();
  const [unpublishCatalog] = useUnpublishPatternMutation();
  const [deletePattern] = useDeletePatternMutation();
  const [importPattern] = useImportPatternMutation();
  const [updatePattern] = useUpdatePatternFileMutation();
  const [uploadPatternFile] = useUploadPatternFileMutation();
  const [deletePatternFile] = useDeletePatternFileMutation();

  useEffect(() => {
    if (patternsData) {
      const filteredPatterns = patternsData.patterns.filter((content) => {
        if (visibilityFilter === null || content.visibility === visibilityFilter) {
          return true;
        }
        return false;
      });
      setCount(patternsData.total_count || 0);
      handleSetPatterns(filteredPatterns);
      setVisibilityFilter(visibilityFilter);
      setPatterns(patternsData.patterns || []);
    }
  }, [patternsData]);

  const [importModal, setImportModal] = useState({
    open: false,
  });
  const [publishModal, setPublishModal] = useState({
    open: false,
    pattern: {},
    name: '',
  });

  const [downloadModal, setDownloadModal] = useState({
    open: false,
    content: null,
  });

  const designValidationActorRef = useActorRef(designValidationMachine);

  const designLifecycleModal = useModal({
    headerIcon: <PatternIcon fill="#fff" height={'2rem'} width={'2rem'} />,
  });
  const sistentInfoModal = useModal({
    headerIcon: OutlinedPatternIcon,
  });
  const handleDeploy = async ({ design, selectedK8sContexts }) => {
    updateProgress({ showProgress: true });
    await deployPatternMutation({
      pattern_file: encodeDesignFile(design),
      pattern_id: design.id,
      selectedK8sContexts,
    });
    updateProgress({ showProgress: false });
  };

  const handleUndeploy = async ({ design, selectedK8sContexts }) => {
    updateProgress({ showProgress: true });
    await undeployPatternMutation({
      pattern_file: encodeDesignFile(design),
      pattern_id: design.id,
      selectedK8sContexts,
    });
    updateProgress({ showProgress: false });
  };

  const handleDownloadDialogClose = () => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: false,
      content: null,
    }));
  };

  const handleDesignDownloadModal = (e, pattern) => {
    e.stopPropagation();
    setDownloadModal((prevState) => ({
      ...prevState,
      open: true,
      content: pattern,
    }));
  };

  // const [loading, stillLoading] = useState(true);
  const { width } = useWindowDimensions();

  const catalogVisibilityRef = useRef(false);
  const catalogContentRef = useRef();
  const disposeConfSubscriptionRef = useRef(null);

  const ACTION_TYPES = {
    FETCH_PATTERNS: {
      name: 'FETCH_PATTERNS',
      error_msg: 'Failed to fetch designs',
    },
    UPDATE_PATTERN: {
      name: 'UPDATE_PATTERN',
      error_msg: 'Failed to update design file',
    },
    DELETE_PATTERN: {
      name: 'DELETE_PATTERN',
      error_msg: 'Failed to delete design file',
    },
    DEPLOY_PATTERN: {
      name: 'DEPLOY_PATTERN',
      error_msg: 'Failed to deploy design file',
    },
    UNDEPLOY_PATTERN: {
      name: 'UNDEPLOY_PATTERN',
      error_msg: 'Failed to undeploy design file',
    },
    UPLOAD_PATTERN: {
      name: 'UPLOAD_PATTERN',
      error_msg: 'Failed to upload design file',
    },
    CLONE_PATTERN: {
      name: 'CLONE_PATTERN',
      error_msg: 'Failed to clone design file',
    },
    PUBLISH_CATALOG: {
      name: 'PUBLISH_CATALOG',
      error_msg: 'Failed to publish catalog',
    },
    UNPUBLISH_CATALOG: {
      name: 'UNPUBLISH_CATALOG',
      error_msg: 'Failed to unpublish catalog',
    },
    SCHEMA_FETCH: {
      name: 'SCHEMA_FETCH',
      error_msg: 'failed to fetch import schema',
    },
  };

  /**
   * Checking whether users are signed in under a provider that doesn't have
   * publish pattern capability and setting the canPublishPattern state accordingly
   */
  const { data: capabilitiesData } = useGetProviderCapabilitiesQuery();

  useEffect(() => {
    if (capabilitiesData) {
      const capabilitiesRegistry = capabilitiesData;
      const patternsCatalogueCapability = capabilitiesRegistry?.capabilities.filter(
        (val) => val.feature === MesheryPatternsCatalog,
      );
      if (patternsCatalogueCapability?.length > 0) setCanPublishPattern(true);
    }
  }, []);

  const searchTimeout = useRef(null);
  /**
   * fetch patterns when the page loads
   */
  // @ts-ignore
  useEffect(() => {
    document.body.style.overflowX = 'hidden';
    const visibilityFilter =
      selectedFilters.visibility === 'All' ? null : selectedFilters.visibility;
    setVisibilityFilter(visibilityFilter);
    return () => (document.body.style.overflowX = 'auto');
  }, [visibilityFilter]);

  useEffect(() => {
    if (viewType === 'grid') {
      setSearch('');
    }
  }, [viewType]);

  const initPatternsSubscription = (
    pageNo = page.toString(),
    pagesize = pageSize.toString(),
    searchText = search,
    order = sortOrder,
  ) => {
    if (disposeConfSubscriptionRef.current) {
      disposeConfSubscriptionRef.current.dispose();
    }
    const configurationSubscription = ConfigurationSubscription(
      () => {
        // stillLoading(false);
        /**
         * We are not using pattern subscription and this code is commented to prevent
         * unnecessary state updates
         */
        // setPage(result.configuration?.patterns?.page || 0);
        // setPageSize(result.configuration?.patterns?.page_size || 10);
        // setCount(result.configuration?.patterns?.total_count || 0);
        // handleSetPatterns(result.configuration?.patterns?.patterns);
      },
      {
        applicationSelector: {
          pageSize: pagesize,
          page: pageNo,
          search: searchText,
          order: order,
        },
        patternSelector: {
          pageSize: pagesize,
          page: pageNo,
          search: searchText,
          order: order,
        },
        filterSelector: {
          pageSize: pagesize,
          page: pageNo,
          search: searchText,
          order: order,
        },
      },
    );
    disposeConfSubscriptionRef.current = configurationSubscription;
  };

  // const handleCatalogVisibility = () => {
  //   handleCatalogPreference(!catalogVisibilityRef.current);
  //   catalogVisibilityRef.current = !catalogVisibility;
  //   toggleCatalogContent({ catalogVisibility: !catalogVisibility });
  // };

  useEffect(async () => {
    try {
      const { models } = await getMeshModels();
      const modelNames = _.uniqBy(
        models?.map((model) => {
          if (model.displayName && model.displayName !== '') {
            return model.displayName;
          }
        }),
        _.toLower,
      );
      modelNames.sort();

      // Modify the schema using the utility function
      const modifiedSchema = modifyRJSFSchema(
        publishCatalogItemSchema,
        'properties.compatibility.items.enum',
        modelNames,
      );

      setPublishSchema({ rjsfSchema: modifiedSchema, uiSchema: publishCatalogItemUiSchema });
      setMeshModels(models);
    } catch (err) {
      console.error(err);
      handleError(ACTION_TYPES.SCHEMA_FETCH);
    }

    catalogVisibilityRef.current = catalogVisibility;

    /*
                                       Below is a graphql query that fetches the catalog patterns that is published so
                                       when catalogVisibility is true, we fetch the catalog patterns and set it to the patterns state
                                       which show the catalog patterns only in the UI at the top of the list always whether we filter for public or private patterns.
                                       Meshery's REST API already fetches catalog items with `published` visibility, hence this function is commented out.
                                      */
    // const fetchCatalogPatterns = fetchCatalogPattern({
    //   selector: {
    //     search: '',
    //     order: '',
    //     page: 0,
    //     pagesize: 0,
    //   },
    // }).subscribe({
    //   next: (result) => {
    //     catalogContentRef.current = result?.catalogPatterns;
    //     initPatternsSubscription();
    //   },
    //   error: (err) => console.log('There was an error fetching Catalog Filter: ', err),
    // });

    // return () => {
    //   fetchCatalogPatterns.unsubscribe();
    //   disposeConfSubscriptionRef.current?.dispose();
    // };
  }, []);

  // useEffect(() => {
  //   handleSetPatterns(patterns);
  // }, [catalogVisibility]);

  const handleSetPatterns = (patterns) => {
    if (catalogVisibilityRef.current && catalogContentRef.current?.length > 0) {
      setPatterns([
        ...(catalogContentRef.current || []),
        ...(patterns?.filter((content) => content.visibility !== VISIBILITY.PUBLISHED) || []),
      ]);
      return;
    }

    setPatterns(patterns?.filter((content) => content.visibility !== VISIBILITY.PUBLISHED) || []);
  };

  const openDeployModal = (e, pattern_file, name) => {
    const design = parseDesignFile(pattern_file);
    e.stopPropagation();
    designLifecycleModal.openModal({
      title: `Deploy design "${name}"`,
      headerIcon: <DoneAllIcon fill="#fff" height={'2rem'} width={'2rem'} />,
      reactNode: (
        <DeployStepper
          handleClose={designLifecycleModal.closeModal}
          validationMachine={designValidationActorRef}
          design={design}
          handleDeploy={handleDeploy}
          deployment_type={DEPLOYMENT_TYPE.DEPLOY}
          selectedK8sContexts={selectedK8sContexts}
        />
      ),
    });
  };

  const openUndeployModal = (e, pattern_file, name) => {
    e.stopPropagation();
    const design = parseDesignFile(pattern_file);
    designLifecycleModal.openModal({
      title: `Undeploy design "${name}"`,
      headerIcon: <UndeployIcon fill="#fff" height={'2rem'} width={'2rem'} />,
      reactNode: (
        <UnDeployStepper
          handleClose={designLifecycleModal.closeModal}
          validationMachine={designValidationActorRef}
          design={design}
          handleUndeploy={handleUndeploy}
          deployment_type={DEPLOYMENT_TYPE.UNDEPLOY}
          selectedK8sContexts={selectedK8sContexts}
        />
      ),
    });
  };

  const openDryRunModal = (e, pattern_file, name) => {
    e.stopPropagation();

    const design = parseDesignFile(pattern_file);
    designLifecycleModal.openModal({
      title: `Dryrun design "${name}"`,
      headerIcon: <DryRunIcon fill="#fff" height={'2rem'} width={'2rem'} />,
      reactNode: (
        <ModalBody style={{ minWidth: '30rem', width: 'auto' }}>
          <DryRunDesign
            handleClose={designLifecycleModal.closeModal}
            validationMachine={designValidationActorRef}
            design={design}
            deployment_type={DEPLOYMENT_TYPE.DEPLOY}
            selectedK8sContexts={selectedK8sContexts}
          />
        </ModalBody>
      ),
    });
  };

  const openValidateModal = (e, pattern_file, name) => {
    e.stopPropagation();

    const design = parseDesignFile(pattern_file);
    designLifecycleModal.openModal({
      title: `Validate design "${name}"`,
      headerIcon: <CheckIcon fill="#fff" height={'2rem'} width={'2rem'} />,
      reactNode: (
        <ModalBody style={{ minWidth: '30rem', width: 'auto' }}>
          <ValidateDesign
            handleClose={designLifecycleModal.closeModal}
            validationMachine={designValidationActorRef}
            design={design}
            deployment_type={DEPLOYMENT_TYPE.DEPLOY}
            selectedK8sContexts={selectedK8sContexts}
          />
        </ModalBody>
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

  const handleInfoModalClose = () => {
    sistentInfoModal.closeModal();
    setInfoModal({
      open: false,
    });
  };

  const handleInfoModal = (pattern) => {
    sistentInfoModal.openModal({
      title: pattern.name,
    });

    setInfoModal({
      open: true,
      ownerID: pattern.user_id,
      selectedResource: pattern,
    });
  };

  // const handlePublishModal = (ev, pattern) => {
  //   if (canPublishPattern) {
  //     ev.stopPropagation();
  //     setPublishModal({
  //       open: true,
  //       pattern: pattern,
  //       name: '',
  //     });
  //   }
  // };

  const handleUnpublishModal = (ev, pattern) => {
    return async () => {
      let response = await modalRef.current.show({
        title: `Unpublish Catalog item?`,
        subtitle: `Are you sure you want to unpublish ${pattern?.name}?`,
        variant: PROMPT_VARIANTS.DANGER,
        primaryOption: 'UNPUBLISH',
        showInfoIcon:
          "Unpublishing a catolog item removes the item from the public-facing catalog (a public website accessible to anonymous visitors at meshery.io/catalog). The catalog item's visibility will change to either public (or private with a subscription). The ability to for other users to continue to access, edit, clone and collaborate on your content depends upon the assigned visibility level (public or private). Prior collaborators (users with whom you have shared your catalog item) will retain access. However, you can always republish it whenever you want. Remember: unpublished catalog items can still be available to other users if that item is set to public visibility. For detailed information, please refer to the [documentation](https://docs.meshery.io/concepts/designs).",
      });
      if (response === 'UNPUBLISH') {
        updateProgress({ showProgress: true });
        unpublishCatalog({
          unpublishBody: JSON.stringify({ id: pattern?.id }),
        })
          .unwrap()
          .then(() => {
            updateProgress({ showProgress: false });
            notify({
              message: `Design Unpublished`,
              event_type: EVENT_TYPES.SUCCESS,
            });
          })
          .catch(() => {
            updateProgress({ showProgress: false });
            handleError(ACTION_TYPES.UNPUBLISH_CATALOG);
          });
      }
    };
  };

  const handlePublishModalClose = () => {
    setPublishModal({
      open: false,
      pattern: {},
      name: '',
    });
  };

  const handlePublish = (formData) => {
    const compatibilityStore = _.uniqBy(meshModels, (model) => _.toLower(model.displayName))
      ?.filter((model) =>
        formData?.compatibility?.some((comp) => _.toLower(comp) === _.toLower(model.displayName)),
      )
      ?.map((model) => model.name);

    const payload = {
      id: infoModal.selectedResource?.id,
      catalog_data: {
        ...formData,
        compatibility: compatibilityStore,
        type: _.toLower(formData?.type),
      },
    };
    updateProgress({ showProgress: true });
    publishCatalog({
      publishBody: JSON.stringify(payload),
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        if (user.role_names.includes('admin')) {
          notify({
            message: `${publishModal?.name} Design Published`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        } else {
          notify({
            message:
              'Design queued for publishing into Meshery Catalog. Maintainers notified for review',
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.PUBLISH_CATALOG);
      });
  };

  function handleClone(patternID, name) {
    updateProgress({ showProgress: true });
    clonePattern({
      body: JSON.stringify({ name: name + ' (Copy)' }),
      patternID: patternID,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `"${name}" Design cloned`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `Failed to clone "${name}" Design`,
          event_type: EVENT_TYPES.ERROR,
        });
      });
  }

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });

    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
    });
  };

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  async function handleSubmit({ data, id, name, type, metadata, catalog_data }) {
    updateProgress({ showProgress: true });
    if (type === FILE_OPS.DELETE) {
      const response = await showModal(1, name);
      if (response == 'No') {
        updateProgress({ showProgress: false });
        return;
      }
      deletePatternFile({
        id: id,
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
          notify({ message: `"${name}" Design deleted`, event_type: EVENT_TYPES.SUCCESS });
          getPatterns();
          resetSelectedRowData()();
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.DELETE_PATTERN);
        });
    }

    if (type === FILE_OPS.UPDATE) {
      const design = yaml.load(data);

      updatePattern({
        updateBody: JSON.stringify({
          id,
          name: data.name,
          design_file: design,
          catalog_data,
        }),
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
          notify({ message: `"${name}" Design updated`, event_type: EVENT_TYPES.SUCCESS });
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.UPDATE_PATTERN);
        });
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      let body;
      if (type === FILE_OPS.FILE_UPLOAD) {
        body = JSON.stringify({
          pattern_data: {
            name: metadata?.name || name,
            pattern_file: getUnit8ArrayDecodedFile(data),
            catalog_data,
          },
          save: true,
        });
      }
      if (type === FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({
          url: data,
          save: true,
          name: metadata?.name || name,
          catalog_data,
        });
      }
      uploadPatternFile({
        uploadBody: body,
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.UPLOAD_PATTERN);
        });
    }
  }

  const handleDownload = (e, design, source_type, params) => {
    e.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      let id = design.id;
      let name = design.name;
      downloadContent({ id, name, type: 'pattern', source_type, params });
      updateProgress({ showProgress: false });
      notify({ message: `"${name}" design downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (e) {
      console.error(e);
    }
  };

  const userCanEdit = (pattern) => {
    return (
      CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject) && user?.user_id == pattern?.user_id
    );
  };

  const handleOpenInConfigurator = (id) => {
    router.push('/configuration/designs/configurator?design_id=' + id);
  };

  let colViews = [
    ['name', 'xs'],
    ['created_at', 'm'],
    ['updated_at', 'm'],
    ['visibility', 's'],
    ['Actions', 'xs'],
  ];

  const columns = [
    {
      name: 'name',
      label: 'Name',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
      },
    },
    {
      name: 'created_at',
      label: 'Created At',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'updated_at',
      label: 'Update At',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'visibility',
      label: 'Visibility',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
      },
    },
    {
      name: 'Actions',
      label: 'Actions',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          const rowData = patterns[tableMeta.rowIndex];
          const visibility = patterns[tableMeta.rowIndex]?.visibility;
          const actions = [
            {
              label: 'Edit',
              icon: <EditIcon fill="currentColor" />,
              onClick: (e) => {
                e.stopPropagation();
                handleOpenInConfigurator(rowData.id);
              },
              disabled: !CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject),
              condition: userCanEdit(rowData),
            },
            {
              label: 'Clone',
              icon: <CloneIcon fill="currentColor" />,
              onClick: (e) => {
                e.stopPropagation();
                handleClone(rowData.id, rowData.name);
              },
              disabled: !CAN(keys.CLONE_DESIGN.action, keys.CLONE_DESIGN.subject),
              condition: visibility === VISIBILITY.PUBLISHED,
            },
            {
              label: 'Design',
              icon: <PatternConfigureIcon />,
              onClick: (e) => {
                e.stopPropagation();
                handleOpenInConfigurator(patterns[tableMeta.rowIndex].id);
              },
              disabled: !CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject),
              condition: visibility !== VISIBILITY.PUBLISHED,
            },
            {
              label: 'Validate Design',
              icon: <CheckIcon data-cy="verify-button" />,
              onClick: (e) => {
                openValidateModal(e, rowData.pattern_file, rowData.name, rowData.id);
              },
              disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
            },
            {
              label: 'Dry Run',
              icon: <DryRunIcon data-cy="verify-button" />,
              onClick: (e) => {
                openDryRunModal(e, rowData.pattern_file, rowData.name, rowData.id);
              },
              disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
            },
            {
              label: 'Undeploy',
              icon: <UndeployIcon fill="#F91313" data-cy="undeploy-button" />,
              onClick: (e) => {
                openUndeployModal(e, rowData.pattern_file, rowData.name, rowData.id);
              },
              disabled: !CAN(keys.UNDEPLOY_DESIGN.action, keys.UNDEPLOY_DESIGN.subject),
            },
            {
              label: 'Deploy',
              icon: <DoneAllIcon data-cy="deploy-button" />,
              onClick: (e) => {
                openDeployModal(e, rowData.pattern_file, rowData.name, rowData.id);
              },
              disabled: !CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject),
            },
            {
              label: 'Download',
              icon: <GetAppIcon data-cy="download-button" />,
              onClick: (e) => {
                handleDesignDownloadModal(e, rowData);
              },
              disabled: !CAN(keys.DOWNLOAD_A_DESIGN.action, keys.DOWNLOAD_A_DESIGN.subject),
            },
            {
              label: 'Design Information',
              icon: <InfoOutlinedIcon data-cy="information-button" />,
              onClick: (e) => {
                genericClickHandler(e, () => handleInfoModal(rowData));
              },
              disabled: !CAN(keys.DETAILS_OF_DESIGN.action, keys.DETAILS_OF_DESIGN.subject),
            },

            /* Publish action can be done through Info modal so we might not need separate publish action */
            /*{
              label="Publish",
              icon: <PublicIcon fill="#F91313" data-cy="publish-button" />,
              onClick: (e) => handlePublishModal(e, rowData)(),
              disabled: !CAN(keys.PUBLISH_DESIGN.action, keys.PUBLISH_DESIGN.subject),
              condition: canPublishPattern && visibility !== VISIBILITY.PUBLISHED,
            },*/

            {
              label: 'Unpublish',
              icon: <PublicIcon fill="#F91313" data-cy="unpublish-button" />,
              onClick: (e) => {
                handleUnpublishModal(e, rowData)();
              },
              disabled: !CAN(keys.UNPUBLISH_DESIGN.action, keys.UNPUBLISH_DESIGN.subject),
              condition: visibility === VISIBILITY.PUBLISHED,
            },
          ].filter((action) => action.condition === undefined || action.condition);

          return (
            <>
              <EllipsisButtonWrapper>
                <ActionPopover actions={actions} />
              </EllipsisButtonWrapper>

              <ActionWrapper>
                <Box sx={{ display: 'flex' }}>
                  {actions.map((action, index) => (
                    <TooltipIcon
                      key={index}
                      placement="top"
                      title={action.label}
                      onClick={action.onClick}
                      disabled={action.disabled}
                    >
                      {action.icon}
                    </TooltipIcon>
                  ))}
                </Box>
              </ActionWrapper>
            </>
          );
        },
      },
    },
  ];

  columns.forEach((column, idx) => {
    if (column.name === sortOrder.split(' ')[0]) {
      columns[idx].options.sortDirection = sortOrder.split(' ')[1];
    }
  });

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    let showCols = updateVisibleColumns(colViews, width);
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach((col) => {
      if (!(hideVisibility && col.name === 'visibility')) {
        initialVisibility[col.name] = showCols[col.name];
      }
    });
    return initialVisibility;
  });

  async function showModal(count, patterns) {
    let response = await modalRef.current.show({
      title: `Delete ${count ? count : ''} Design${count > 1 ? 's' : ''}?`,

      subtitle: `Are you sure you want to delete the ${patterns} design${count > 1 ? 's' : ''}?`,
      variant: PROMPT_VARIANTS.DANGER,
      primaryOption: 'DELETE',
    });
    return response;
  }

  async function deletePatterns(patterns) {
    const jsonPatterns = JSON.stringify(patterns);

    updateProgress({ showProgress: true });
    deletePattern({
      deleteBody: jsonPatterns,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        setTimeout(() => {
          notify({
            message: `${patterns.patterns.length} Designs deleted`,
            event_type: EVENT_TYPES.SUCCESS,
          });
          getPatterns();
          resetSelectedRowData()();
        }, 1200);
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.DELETE_PATTERN);
      });
  }

  const options = {
    customToolbarSelect: (selectedRows, displayData, setSelectedRows) => (
      <CustomToolbarSelect
        selectedRows={selectedRows}
        displayData={displayData}
        setSelectedRows={setSelectedRows}
        patterns={patterns}
        deletePatterns={deletePatterns}
        showModal={showModal}
      />
    ),
    filter: false,
    search: false,
    viewColumns: false,
    sort: !(user && user.user_id === 'meshery'),
    filterType: 'textField',
    responsive: 'standard',
    resizableColumns: true,
    serverSide: true,
    count,
    rowsPerPage: pageSize,
    fixedHeader: true,
    page,
    print: false,
    download: false,
    sortOrder: {
      name: sortOrder.split(' ')[0],
      direction: sortOrder.split(' ')[1],
    },
    textLabels: {
      selectedRows: {
        text: 'pattern(s) selected',
      },
    },

    onCellClick: (_, meta) =>
      meta.colIndex !== 3 && meta.colIndex !== 4 && setSelectedRowData(patterns[meta.rowIndex]),

    onRowsDelete: async function handleDelete(row) {
      const toBeDeleted = Object.keys(row.lookup).map((idx) => ({
        id: patterns[idx]?.id,
        name: patterns[idx]?.name,
      }));
      let response = await showModal(
        toBeDeleted.length,
        toBeDeleted.map((p) => ' ' + p.name),
      );
      if (response.toLowerCase() === 'DELETE') {
        deletePatterns({ patterns: toBeDeleted });
      }
      // if (response.toLowerCase() === "no")
      // fetchPatterns(page, pageSize, search, sortOrder);
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          initPatternsSubscription(tableState.page.toString(), pageSize.toString(), search, order);
          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          initPatternsSubscription(
            page.toString(),
            tableState.rowsPerPage.toString(),
            search,
            order,
          );
          setPageSize(tableState.rowsPerPage);
          break;
        case 'search':
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              setSearch(tableState.searchText);
            }
          }, 500);
          break;
        case 'sort':
          if (sortInfo.length === 2) {
            if (sortInfo[1] === 'ascending') {
              order = `${columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) {
            initPatternsSubscription(page.toString(), pageSize.toString(), search, order);
            setSortOrder(order);
          }
          break;
      }
    },
    setRowProps: (row, dataIndex, rowIndex) => {
      return {
        'data-cy': `config-row-${rowIndex}`,
      };
    },
    setTableProps: () => {
      return {
        'data-cy': 'filters-grid',
      };
    },
  };

  if (ispatternsLoading) {
    return (
      <>
        <LoadingScreen animatedIcon="AnimatedMeshPattern" message={`Loading ${pageTitle}...`} />
      </>
    );
  }

  /**
   * Gets the data of Import Filter and handles submit operation
   *
   * @param {{
   * uploadType: ("File Upload"| "URL Import");
   * name: string;
   * url: string;
   * file: string;
   * }} data
   */
  function handleImportDesign(data) {
    updateProgress({ showProgress: true });
    const { uploadType, name, url, file } = data;

    let requestBody = null;
    switch (uploadType) {
      case 'File Upload': {
        const fileElement = document.getElementById('root_file');
        const fileName = fileElement.files[0].name;
        requestBody = JSON.stringify({
          name,
          file_name: fileName,
          file: getUnit8ArrayDecodedFile(file),
        });
        break;
      }
      case 'URL Import':
        requestBody = JSON.stringify({
          url,
          name,
        });
        break;
    }

    importPattern({
      importBody: requestBody,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `"${name}" design uploaded`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        getPatterns();
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.UPLOAD_PATTERN);
      });
  }

  const filter = {
    visibility: {
      name: 'Visibility',
      //if catalog content is enabled, then show all filters including published otherwise only show public and private filters
      options: catalogVisibility
        ? [
            { label: 'Public', value: 'public' },
            { label: 'Private', value: 'private' },
            { label: 'Published', value: 'published' },
          ]
        : [
            { label: 'Public', value: 'public' },
            { label: 'Private', value: 'private' },
          ],
    },
  };

  const handleApplyFilter = () => {
    const visibilityFilter =
      selectedFilters.visibility === 'All' ? null : selectedFilters.visibility;
    setVisibilityFilter(visibilityFilter);
  };

  return (
    <>
      <NoSsr>
        {CAN(keys.VIEW_DESIGNS.action, keys.VIEW_DESIGNS.subject) ? (
          <>
            {selectedRowData && Object.keys(selectedRowData).length > 0 && (
              <YAMLEditor
                pattern={selectedRowData}
                onClose={resetSelectedRowData()}
                onSubmit={handleSubmit}
                isReadOnly={arePatternsReadOnly}
              />
            )}
            <ToolWrapper>
              {width < 600 && isSearchExpanded ? null : (
                <CreateButton style={{ display: 'flex' }}>
                  {!selectedPattern.show && (patterns.length >= 0 || viewType === 'table') && (
                    <div>
                      {disableCreateImportDesignButton ? null : (
                        <div style={{ display: 'flex', order: '1' }}>
                          <TooltipButton
                            title="Create Design"
                            aria-label="Add Pattern"
                            variant="contained"
                            color="primary"
                            size="large"
                            // @ts-ignore
                            onClick={() => router.push('designs/configurator')}
                            style={{ display: 'flex', marginRight: '2rem' }}
                            disabled={
                              !CAN(keys.CREATE_NEW_DESIGN.action, keys.CREATE_NEW_DESIGN.subject)
                            }
                          >
                            <AddIconStyled />
                            <BtnText> Create Design </BtnText>
                          </TooltipButton>
                          <TooltipButton
                            title="Import Design"
                            aria-label="Add Pattern"
                            variant="contained"
                            color="primary"
                            size="large"
                            // @ts-ignore
                            onClick={handleUploadImport}
                            style={{ display: 'flex', marginRight: '2rem', marginLeft: '-0.6rem' }}
                            disabled={!CAN(keys.IMPORT_DESIGN.action, keys.IMPORT_DESIGN.subject)}
                          >
                            <AddIconStyled>
                              <PublishIcon />
                            </AddIconStyled>
                            <BtnText> Import Design </BtnText>
                          </TooltipButton>
                        </div>
                      )}
                    </div>
                  )}
                  {!selectedPattern.show && (
                    <div style={{ display: 'flex' }}>
                      {/* <StyledCatalogFilter>
                      <CatalogFilter
                        catalogVisibility={catalogVisibility}
                        handleCatalogVisibility={handleCatalogVisibility}
                        classes={classes}
                      />
                      </StyledCatalogFilter>*/}
                    </div>
                  )}
                </CreateButton>
              )}
              <SearchWrapper style={{ display: 'flex' }}>
                <>
                  <SearchBar
                    onSearch={(value) => {
                      setSearch(value);
                      initPatternsSubscription(
                        page.toString(),
                        pageSize.toString(),
                        value,
                        sortOrder,
                      );
                    }}
                    expanded={isSearchExpanded}
                    setExpanded={setIsSearchExpanded}
                    placeholder={`Search ${pageTitle.toLowerCase()}...`}
                  />
                  {disableUniversalFilter ? null : (
                    <UniversalFilter
                      id="ref"
                      filters={filter}
                      selectedFilters={selectedFilters}
                      setSelectedFilters={setSelectedFilters}
                      handleApplyFilter={handleApplyFilter}
                    />
                  )}
                  {viewType === 'table' && (
                    <CustomColumnVisibilityControl
                      id="ref"
                      columns={columns}
                      customToolsProps={{ columnVisibility, setColumnVisibility }}
                    />
                  )}
                </>

                {!selectedPattern.show && (
                  <ViewSwitchButton>
                    <ViewSwitch view={viewType} changeView={setViewType} />
                  </ViewSwitchButton>
                )}
              </SearchWrapper>
            </ToolWrapper>
            {!selectedPattern.show && viewType === 'table' && (
              <>
                {/* <StyledRow> */}
                <ResponsiveDataTable
                  data={patterns}
                  columns={columns}
                  // @ts-ignore
                  options={options}
                  tableCols={tableCols}
                  updateCols={updateCols}
                  columnVisibility={columnVisibility}
                />
                {/* </StyledRow> */}

                <ExportModal
                  downloadModal={downloadModal}
                  handleDownloadDialogClose={handleDownloadDialogClose}
                  handleDesignDownload={handleDownload}
                />
              </>
            )}
            {!selectedPattern.show && viewType === 'grid' && (
              // grid vieww
              <MesheryPatternGrid
                selectedK8sContexts={selectedK8sContexts}
                canPublishPattern={canPublishPattern}
                patterns={patterns}
                handlePublish={handlePublish}
                handleUnpublishModal={handleUnpublishModal}
                handleClone={handleClone}
                supportedTypes="null"
                handleSubmit={handleSubmit}
                setSelectedPattern={setSelectedPattern}
                selectedPattern={selectedPattern}
                pages={Math.ceil(count / pageSize)}
                setPage={setPage}
                selectedPage={page}
                patternErrors={[]}
                publishModal={publishModal}
                setPublishModal={setPublishModal}
                publishSchema={publishSchema}
                user={user}
                fetch={() => getPatterns()}
                handleInfoModal={handleInfoModal}
                openUndeployModal={openUndeployModal}
                openValidationModal={openValidateModal}
                openDryRunModal={openDryRunModal}
                openDeployModal={openDeployModal}
                hideVisibility={hideVisibility}
                arePatternsReadOnly={arePatternsReadOnly}
              />
            )}

            <SistentModal maxWidth="sm" {...designLifecycleModal}></SistentModal>
            <SistentModal {...sistentInfoModal}>
              {CAN(keys.DETAILS_OF_DESIGN.action, keys.DETAILS_OF_DESIGN.subject) &&
                infoModal.open && (
                  <InfoModal
                    handlePublish={handlePublish}
                    infoModalOpen={true}
                    handleInfoModalClose={handleInfoModalClose}
                    dataName="patterns"
                    selectedResource={infoModal.selectedResource}
                    resourceOwnerID={infoModal.ownerID}
                    currentUser={user}
                    patternFetcher={getPatterns}
                    formSchema={publishSchema}
                    meshModels={meshModels}
                  />
                )}
            </SistentModal>

            {canPublishPattern &&
              publishModal.open &&
              CAN(keys.PUBLISH_DESIGN.action, keys.PUBLISH_DESIGN.subject) && (
                <PublishModal
                  publishFormSchema={publishSchema}
                  handleClose={handlePublishModalClose}
                  title={publishModal.pattern?.name}
                  handleSubmit={handlePublish}
                />
              )}
            {importModal.open && CAN(keys.IMPORT_DESIGN.action, keys.IMPORT_DESIGN.subject) && (
              <ImportDesignModal
                handleClose={handleUploadImportClose}
                handleImportDesign={handleImportDesign}
              />
            )}
            <_PromptComponent ref={modalRef} />
          </>
        ) : (
          <DefaultError />
        )}
      </NoSsr>
    </>
  );
}

export const ImportDesignModal = React.memo((props) => {
  const { handleClose, handleImportDesign } = props;

  return (
    <>
      <>
        <SistentModal
          open={true}
          closeModal={handleClose}
          headerIcon={
            <Pattern fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
          }
          maxWidth="sm"
          title="Import Design"
        >
          <RJSFModalWrapper
            schema={importDesignSchema}
            uiSchema={importDesignUiSchema}
            handleSubmit={handleImportDesign}
            submitBtnText="Import"
            handleClose={handleClose}
          />
        </SistentModal>
      </>
    </>
  );
});

const PublishModal = React.memo((props) => {
  const { handleClose, handleSubmit, title } = props;

  return (
    <>
      <>
        <SistentModal
          open={true}
          closeModal={handleClose}
          aria-label="catalog publish"
          title={title}
          headerIcon={
            <Pattern fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
          }
          maxWidth="sm"
        >
          <RJSFModalWrapper
            schema={publishCatalogItemSchema}
            uiSchema={publishCatalogItemUiSchema}
            handleSubmit={handleSubmit}
            submitBtnText="Submit for Approval"
            handleClose={handleClose}
            helpText="Upon submitting your catalog item, an approval flow will be initiated.[Learn more](https://docs.meshery.io/concepts/catalog)"
          />
        </SistentModal>
      </>
    </>
  );
});

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  toggleCatalogContent: bindActionCreators(toggleCatalogContent, dispatch),
});

const mapStateToProps = (state) => ({
  user: state.get('user')?.toObject(),
  selectedK8sContexts: state.get('selectedK8sContexts'),
  catalogVisibility: state.get('catalogVisibility'),
});

// @ts-ignore
export default connect(mapStateToProps, mapDispatchToProps)(MesheryPatterns);
