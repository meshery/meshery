/* eslint-disable react/display-name */
import React, { useState, useEffect, useRef } from 'react';
import { NoSsr } from '@layer5/sistent';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import DeleteIcon from '@mui/icons-material/Delete';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Moment from 'react-moment';
import CloseIcon from '@mui/icons-material/Close';
import EditIcon from '@mui/icons-material/Edit';
import { toggleCatalogContent, updateProgress } from '../lib/store';
import _PromptComponent from './PromptComponent';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import { FILE_OPS, MesheryFiltersCatalog, VISIBILITY } from '../utils/Enum';
import ViewSwitch from './ViewSwitch';
import FiltersGrid from './MesheryFilters/FiltersGrid';
import { trueRandom } from '../lib/trueRandom';
import GetAppIcon from '@mui/icons-material/GetApp';
import PublicIcon from '@mui/icons-material/Public';
import PublishIcon from '@mui/icons-material/Publish';
import downloadContent from '../utils/fileDownloader';
import CloneIcon from '../public/static/img/CloneIcon';
import SaveIcon from '@mui/icons-material/Save';
import ConfigurationSubscription from './graphql/subscriptions/ConfigurationSubscription';
import fetchCatalogFilter from './graphql/queries/CatalogFilterQuery';
import { iconMedium } from '../css/icons.styles';
import { RJSFModalWrapper } from './Modal';
import { getUnit8ArrayDecodedFile, modifyRJSFSchema } from '../utils/utils';
import Filter from '../public/static/img/drawer-icons/filter_svg.js';
import { getMeshModels } from '../api/meshmodel';
import _ from 'lodash';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import {
  CustomColumnVisibilityControl,
  CustomTooltip,
  ResponsiveDataTable,
  SearchBar,
  UniversalFilter,
  importFilterSchema,
  importFilterUiSchema,
  publishCatalogItemSchema,
  publishCatalogItemUiSchema,
  IconButton,
  Dialog,
  DialogTitle,
  DialogActions,
  Divider,
  Typography,
  Button,
  Box,
  styled,
  PROMPT_VARIANTS,
} from '@layer5/sistent';
import { updateVisibleColumns } from '../utils/responsive-column';
import { useWindowDimensions } from '../utils/dimension';
import InfoModal from './Modals/Information/InfoModal';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import { DefaultTableCell, SortableTableCell } from './connections/common/index.js';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import DefaultError from './General/error-404/index';
import { Modal as SistentModal } from '@layer5/sistent';
import {
  useGetFiltersQuery,
  useCloneFilterMutation,
  usePublishFilterMutation,
  useUnpublishFilterMutation,
  useDeleteFilterMutation,
  useUpdateFilterFileMutation,
  useUploadFilterFileMutation,
} from '@/rtk-query/filter';
import LoadingScreen from './LoadingComponents/LoadingComponent';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';

const CreateButton = styled(Button)(() => ({
  width: 'fit-content',
  alignSelf: 'flex-start',
  placeSelf: 'center',
}));

const ViewSwitchButton = styled('div')(() => ({
  justifySelf: 'flex-end',
  paddingLeft: '1rem',
}));

const YmlDialogTitle = styled(DialogTitle)(() => ({
  display: 'flex',
  alignItems: 'center',
}));

const YmlDialogTitleText = styled(Typography)(() => ({
  flexGrow: 1,
}));

const BtnText = styled('span')(({ theme }) => ({
  display: 'block',
  [theme.breakpoints.down('700')]: {
    display: 'none',
  },
}));

const ActionsBox = styled(Box)(() => ({
  display: 'flex',
}));

function TooltipIcon({ children, onClick, title }) {
  return (
    <>
      <CustomTooltip title={title} placement="top" interactive>
        <div>
          <IconButton onClick={onClick}>{children}</IconButton>
        </div>
      </CustomTooltip>
    </>
  );
}

function YAMLEditor({ filter, onClose, onSubmit }) {
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const FullScreenCodeMirrorWrapper = styled('div')(() => ({
    height: '100%',
    '& .CodeMirror': {
      minHeight: '300px',
      height: fullScreen ? '80vh' : '100%',
    },
  }));

  let resourceData;
  try {
    resourceData = JSON.parse(filter.filter_resource);
  } catch (error) {
    // Handling the error or provide a default value
    console.error('Error parsing JSON:', error);
    resourceData = {}; // Setting a default value if parsing fails
  }

  const config = resourceData?.settings?.config || '';
  const [yaml, setYaml] = useState(config);

  return (
    <Dialog
      onClose={onClose}
      aria-labelledby="filter-dialog-title"
      open
      maxWidth="md"
      fullScreen={fullScreen}
      fullWidth={!fullScreen}
    >
      <YmlDialogTitle>
        <DialogTitle
          disableTypography
          id="filter-dialog-title"
          style={{ width: '100%', display: 'flex' }}
        >
          <YmlDialogTitleText variant="h6">{filter.name}</YmlDialogTitleText>
          <TooltipIcon
            title={fullScreen ? 'Exit Fullscreen' : 'Enter Fullscreen'}
            onClick={toggleFullScreen}
          >
            {fullScreen ? (
              <FullscreenExitIcon style={iconMedium} />
            ) : (
              <FullscreenIcon style={iconMedium} />
            )}
          </TooltipIcon>
          <TooltipIcon title="Exit" onClick={onClose}>
            <CloseIcon style={iconMedium} />
          </TooltipIcon>
        </DialogTitle>
      </YmlDialogTitle>
      <Divider variant="fullWidth" light />
      <FullScreenCodeMirrorWrapper>
        <CodeMirror
          value={config}
          options={{
            theme: 'material',
            lineNumbers: true,
            lineWrapping: true,
            gutters: ['CodeMirror-lint-markers'],
            lint: true,
            mode: 'text/x-yaml',
          }}
          onChange={(_, data, val) => setYaml(val)}
        />
      </FullScreenCodeMirrorWrapper>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <CustomTooltip title="Update Filter">
          <IconButton
            aria-label="Update"
            disabled={!CAN(keys.EDIT_WASM_FILTER.action, keys.EDIT_WASM_FILTER.subject)}
            onClick={() =>
              onSubmit({
                data: yaml,
                id: filter.id,
                name: filter.name,
                type: FILE_OPS.UPDATE,
                catalog_data: filter.catalog_data,
              })
            }
          >
            <SaveIcon style={iconMedium} />
          </IconButton>
        </CustomTooltip>
        <CustomTooltip title="Delete Filter">
          <IconButton
            aria-label="Delete"
            disabled={!CAN(keys.DELETE_WASM_FILTER.action, keys.DELETE_WASM_FILTER.subject)}
            onClick={() =>
              onSubmit({
                data: yaml,
                id: filter.id,
                name: filter.name,
                type: FILE_OPS.DELETE,
                catalog_data: filter.catalog_data,
              })
            }
          >
            <DeleteIcon style={iconMedium} />
          </IconButton>
        </CustomTooltip>
      </DialogActions>
    </Dialog>
  );
}

function resetSelectedFilter() {
  return { show: false, filter: null };
}

function MesheryFilters({
  updateProgress,
  user,
  catalogVisibility,
  // toggleCatalogContent,
}) {
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [count, setCount] = useState(0);
  const modalRef = useRef(null);
  const [pageSize, setPageSize] = useState(10);
  const [filters, setFilters] = useState([]);
  const [selectedFilter, setSelectedFilter] = useState(resetSelectedFilter());
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [canPublishFilter, setCanPublishFilter] = useState(false);
  const [publishSchema, setPublishSchema] = useState({});
  const { width } = useWindowDimensions();
  const [meshModels, setMeshModels] = useState([]);
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ('grid'),
  );

  //hooks
  const { notify } = useNotification();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);

  const [infoModal, setInfoModal] = useState({
    open: false,
    ownerID: '',
    selectedResource: {},
  });

  const [importModal, setImportModal] = useState({
    open: false,
  });
  const [publishModal, setPublishModal] = useState({
    open: false,
    filter: {},
    name: '',
  });

  const catalogContentRef = useRef();
  const catalogVisibilityRef = useRef();
  const disposeConfSubscriptionRef = useRef(null);

  const [visibilityFilter, setVisibilityFilter] = useState(null);

  const {
    data: filtersData,
    isLoading: isFiltersLoading,
    refetch: getFilters,
  } = useGetFiltersQuery({
    page: page,
    pagesize: pageSize,
    search: search,
    order: sortOrder,
    visibility: visibilityFilter ? JSON.stringify([visibilityFilter]) : '',
  });

  const { data: capabilitiesData } = useGetProviderCapabilitiesQuery();

  const [cloneFilter] = useCloneFilterMutation();
  const [publishFilter] = usePublishFilterMutation();
  const [unpublishFilter] = useUnpublishFilterMutation();
  const [deleteFilterFile] = useDeleteFilterMutation();
  const [updateFilterFile] = useUpdateFilterFileMutation();
  const [uploadFilterFile] = useUploadFilterFileMutation();

  useEffect(() => {
    if (filtersData) {
      const filteredWasmFilters = filtersData.filters.filter((content) => {
        if (visibilityFilter === null || content.visibility === visibilityFilter) {
          return true;
        }
        return false;
      });
      setCount(filtersData.total_count || 0);
      handleSetFilters(filteredWasmFilters);
      setVisibilityFilter(visibilityFilter);
      setFilters(filtersData.filters || []);
    }
  }, [filtersData]);

  const ACTION_TYPES = {
    FETCH_FILTERS: {
      name: 'FETCH_FILTERS',
      error_msg: 'Failed to fetch filter',
    },
    DELETE_FILTERS: {
      name: 'DELETE_FILTERS',
      error_msg: 'Failed to delete filter file',
    },
    DEPLOY_FILTERS: {
      name: 'DEPLOY_FILTERS',
      error_msg: 'Failed to deploy filter file',
    },
    UNDEPLOY_FILTERS: {
      name: 'UNDEPLOY_FILTERS',
      error_msg: 'Failed to undeploy filter file',
    },
    UPLOAD_FILTERS: {
      name: 'UPLOAD_FILTERS',
      error_msg: 'Failed to upload filter file',
    },
    CLONE_FILTERS: {
      name: 'CLONE_FILTER',
      error_msg: 'Failed to clone filter file',
    },
    PUBLISH_CATALOG: {
      name: 'PUBLISH_CATALOG',
      error_msg: 'Failed to publish catalog',
    },
    UNPUBLISH_CATALOG: {
      name: 'PUBLISH_CATALOG',
      error_msg: 'Failed to publish catalog',
    },
    SCHEMA_FETCH: {
      name: 'SCHEMA_FETCH',
      error_msg: 'failed to fetch import schema',
    },
  };

  /**
   * Checking whether users are signed in under a provider that doesn't have
   * publish filter capability and setting the canPublishFilter state accordingly
   */
  useEffect(async () => {
    try {
      const { models } = await getMeshModels();
      const modelNames = _.uniq(models?.map((model) => model.displayName));
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
      handleError(ACTION_TYPES.SCHEMA_FETCH);
    }

    if (capabilitiesData) {
      const capabilitiesRegistry = capabilitiesData;
      const filtersCatalogCapability = capabilitiesRegistry?.capabilities.filter(
        (val) => val.feature === MesheryFiltersCatalog,
      );
      if (filtersCatalogCapability?.length > 0) setCanPublishFilter(true);
    }
  }, [capabilitiesData]);

  const searchTimeout = useRef(null);

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

  const handlePublishModal = (ev, filter) => {
    if (canPublishFilter) {
      ev.stopPropagation();
      setPublishModal({
        open: true,
        filter: filter,
      });
    }
  };

  const handleInfoModalClose = () => {
    setInfoModal({
      open: false,
    });
  };

  const handleInfoModal = (filter) => {
    setInfoModal({
      open: true,
      ownerID: filter.user_id,
      selectedResource: filter,
    });
  };

  const handleUnpublishModal = (ev, filter) => {
    if (canPublishFilter) {
      return async () => {
        let response = await modalRef.current.show({
          title: `Unpublish Catalog item?`,
          subtitle: `Are you sure that you want to unpublish "${filter?.name}"?`,
          primaryOption: 'Yes',
          showInfoIcon: `Unpublishing a catolog item removes the item from the public-facing catalog (a public website accessible to anonymous visitors at meshery.io/catalog). The catalog item's visibility will change to either public (or private with a subscription). The ability to for other users to continue to access, edit, clone and collaborate on your content depends upon the assigned visibility level (public or private). Prior collaborators (users with whom you have shared your catalog item) will retain access. However, you can always republish it whenever you want.  Remember: unpublished catalog items can still be available to other users if that item is set to public visibility. For detailed information, please refer to the documentation https://docs.meshery.io/concepts/designs.`,
        });
        if (response === 'Yes') {
          updateProgress({ showProgress: true });
          unpublishFilter({ unpublishBody: JSON.stringify({ id: filter?.id }) })
            .unwrap()
            .then(() => {
              updateProgress({ showProgress: false });
              notify({
                message: `"${filter?.name}" filter unpublished`,
                event_type: EVENT_TYPES.SUCCESS,
              });
            })
            .catch(() => {
              updateProgress({ showProgress: false });
              handleError(ACTION_TYPES.UNPUBLISH_CATALOG);
            });
        }
      };
    }
  };

  const handlePublishModalClose = () => {
    setPublishModal({
      open: false,
      filter: {},
      name: '',
    });
  };

  useEffect(() => {
    if (viewType === 'grid') setSearch('');
  }, [viewType]);

  // const handleCatalogPreference = (catalogPref) => {
  //   let body = Object.assign({}, extensionPreferences);
  //   body['catalogContent'] = catalogPref;

  //   dataFetch(
  //     '/api/user/prefs',
  //     {
  //       method: 'POST',
  //       credentials: 'include',
  //       body: JSON.stringify({ usersExtensionPreferences: body }),
  //     },
  //     () => {
  //       notify({
  //         message: `Catalog Content was ${catalogPref ? 'enab' : 'disab'}led`,
  //         event_type: EVENT_TYPES.SUCCESS,
  //       });
  //     },
  //     (err) => console.error(err),
  //   );
  // };

  // const handleCatalogVisibility = () => {
  //   handleCatalogPreference(!catalogVisibilityRef.current);
  //   catalogVisibilityRef.current = !catalogVisibility;
  //   toggleCatalogContent({ catalogVisibility: !catalogVisibility });
  // };

  useEffect(() => {
    catalogVisibilityRef.current = catalogVisibility;
    const fetchCatalogFilters = fetchCatalogFilter({
      selector: {
        search: '',
        order: '',
        page: 0,
        pagesize: 0,
      },
    }).subscribe({
      next: (result) => {
        catalogContentRef.current = result?.catalogFilters;
        initFiltersSubscription();
      },
      error: (err) => console.log('There was an error fetching Catalog Filter: ', err),
    });

    return () => {
      fetchCatalogFilters.unsubscribe();
      disposeConfSubscriptionRef.current?.dispose();
    };
  }, []);

  const handlePublish = (formData) => {
    const compatibilityStore = _.uniqBy(meshModels, (model) => _.toLower(model.displayName))
      ?.filter((model) =>
        formData?.compatibility?.some((comp) => _.toLower(comp) === _.toLower(model.displayName)),
      )
      ?.map((model) => model.name);

    const payload = {
      id: publishModal.filter?.id,
      catalog_data: {
        ...formData,
        compatibility: compatibilityStore,
        type: _.toLower(formData?.type),
      },
    };
    updateProgress({ showProgress: true });
    publishFilter({ publishBody: JSON.stringify(payload) })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        if (user.role_names.includes('admin')) {
          notify({
            message: `${publishModal.filter?.name} filter published to Meshery Catalog`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        } else {
          notify({
            message:
              'filters queued for publishing into Meshery Catalog. Maintainers notified for review',
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.PUBLISH_CATALOG);
      });
  };

  function handleClone(filterID, name) {
    updateProgress({ showProgress: true });
    cloneFilter({
      body: JSON.stringify({ name: name + ' (Copy)' }),
      filterID: filterID,
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({ message: `"${name}" filter cloned`, event_type: EVENT_TYPES.SUCCESS });
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.CLONE_FILTERS);
      });
  }

  // function handleError(error) {
  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const handleSetFilters = (filters) => {
    if (filters != undefined) {
      if (catalogVisibilityRef.current && catalogContentRef.current?.length > 0) {
        setFilters([
          ...catalogContentRef.current,
          ...filters.filter((content) => content.visibility !== VISIBILITY.PUBLISHED),
        ]);
        return;
      }
      setFilters(filters.filter((content) => content.visibility !== VISIBILITY.PUBLISHED));
    }
  };

  const initFiltersSubscription = (
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
        /**
         * We are not using filter subscription and this code is commented to prevent
         * unnecessary state updates
         */
        // setPage(result.configuration?.filters?.page || 0);
        // setPageSize(result.configuration?.filters?.page_size || 10);
        // setCount(result.configuration?.filters?.total_count || 0);
        // handleSetFilters(result.configuration?.filters?.filters);
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

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  async function handleSubmit({ data, name, id, type, metadata, catalog_data }) {
    // TODO: use filter name
    updateProgress({ showProgress: true });
    if (type === FILE_OPS.DELETE) {
      const response = await showmodal(1);
      if (response !== 'Delete') {
        updateProgress({ showProgress: false });
        return;
      }
      deleteFilterFile({ id: id })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
          notify({ message: `"${name}" filter deleted`, event_type: EVENT_TYPES.SUCCESS });
          resetSelectedRowData()();
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.DELETE_FILTERS);
        });
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      // todo: remove this
      let body = { save: true };
      if (type === FILE_OPS.FILE_UPLOAD) {
        body = JSON.stringify({
          ...body,
          filter_data: { filter_file: data, name: metadata.name },
          config: metadata.config,
        });
      }
      if (type === FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({ ...body, url: data, name: metadata.name, config: metadata.config });
      }
      uploadFilterFile({ uploadBody: body })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.UPLOAD_FILTERS);
        });
    }

    if (type === FILE_OPS.UPDATE) {
      updateFilterFile({
        updateBody: JSON.stringify({
          filter_data: { id, name: name, catalog_data },
          config: data,
          save: true,
        }),
      })
        .unwrap()
        .then(() => {
          updateProgress({ showProgress: false });
        })
        .catch(() => {
          updateProgress({ showProgress: false });
          handleError(ACTION_TYPES.UPLOAD_FILTERS);
        });
    }
  }

  const handleDownload = (e, id, name) => {
    e.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      downloadContent({ id, name, type: 'filter' });
      updateProgress({ showProgress: false });
      notify({ message: `"${name}" filter downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (e) {
      console.error(e);
    }
  };

  function uploadHandler(ev, _, metadata) {
    if (!ev.target.files?.length) return;

    const file = ev.target.files[0];

    // Create a reader
    const reader = new FileReader();
    reader.addEventListener('load', (event) => {
      let uint8 = new Uint8Array(event.target.result);
      handleSubmit({
        data: Array.from(uint8),
        name: file?.name || 'meshery_' + Math.floor(trueRandom() * 100),
        type: FILE_OPS.FILE_UPLOAD,
        metadata: metadata,
      });
    });
    reader.readAsArrayBuffer(file);
  }

  let colViews = [
    ['name', 'xs'],
    ['created_at', 'm'],
    ['updated_at', 'l'],
    ['visibility', 's'],
    ['Actions', 'xs'],
  ];

  const columns = [
    {
      name: 'name',
      label: 'Filter Name',
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
      label: 'Updated At',
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
          const rowData = filters[tableMeta.rowIndex];
          const visibility = filters[tableMeta.rowIndex]?.visibility;
          return (
            <ActionsBox
              sx={{
                display: 'flex',
              }}
            >
              {visibility === VISIBILITY.PUBLISHED ? (
                <TooltipIcon
                  placement="top"
                  title={'Clone'}
                  onClick={(e) => {
                    e.stopPropagation();
                    handleClone(rowData.id, rowData.name);
                  }}
                  disabled={!CAN(keys.CLONE_WASM_FILTER.action, keys.CLONE_WASM_FILTER.subject)}
                >
                  <CloneIcon fill="currentColor" />
                </TooltipIcon>
              ) : (
                <TooltipIcon
                  title="Config"
                  onClick={(e) => {
                    e.stopPropagation();
                    setSelectedRowData(filters[tableMeta.rowIndex]);
                  }}
                  disabled={!CAN(keys.EDIT_WASM_FILTER.action, keys.EDIT_WASM_FILTER.subject)}
                >
                  <EditIcon aria-label="config" color="inherit" style={iconMedium} />
                </TooltipIcon>
              )}
              <TooltipIcon
                title="Download"
                onClick={(e) => handleDownload(e, rowData.id, rowData.name)}
                disabled={
                  !CAN(keys.DOWNLOAD_A_WASM_FILTER.action, keys.DOWNLOAD_A_WASM_FILTER.subject)
                }
              >
                <GetAppIcon data-cy="download-button" />
              </TooltipIcon>
              <TooltipIcon
                title="Filter Information"
                onClick={() => handleInfoModal(rowData)}
                disabled={
                  !CAN(keys.DETAILS_OF_WASM_FILTER.action, keys.DETAILS_OF_WASM_FILTER.subject)
                }
              >
                <InfoOutlinedIcon data-cy="information-button" />
              </TooltipIcon>
              {canPublishFilter && visibility !== VISIBILITY.PUBLISHED ? (
                <TooltipIcon
                  title="Publish"
                  onClick={(ev) => handlePublishModal(ev, rowData)}
                  disabled={!CAN(keys.PUBLISH_WASM_FILTER.action, keys.PUBLISH_WASM_FILTER.subject)}
                >
                  <PublicIcon fill="#F91313" data-cy="publish-button" />
                </TooltipIcon>
              ) : (
                <TooltipIcon
                  title="Unpublish"
                  onClick={(ev) => handleUnpublishModal(ev, rowData)()}
                  disabled={
                    !CAN(keys.UNPUBLISH_WASM_FILTER.action, keys.UNPUBLISH_WASM_FILTER.subject)
                  }
                >
                  <PublicIcon fill="#F91313" data-cy="unpublish-button" />
                </TooltipIcon>
              )}
            </ActionsBox>
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

  async function showmodal(count) {
    let response = await modalRef.current.show({
      title: `Delete ${count ? count : ''} Filter${count > 1 ? 's' : ''}?`,
      subtitle: `Are you sure you want to delete ${count > 1 ? 'these' : 'this'} ${
        count ? count : ''
      } filter${count > 1 ? 's' : ''}?`,
      primaryOption: 'Delete',
      variant: PROMPT_VARIANTS.DANGER,
    });
    return response;
  }

  function deleteFilter(id) {
    deleteFilterFile({ id: id })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({ message: `Filter deleted`, event_type: EVENT_TYPES.SUCCESS });
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.DELETE_FILTERS);
      });
  }

  const options = {
    filter: false,
    viewColumns: false,
    sort: !(user && user.user_id === 'meshery'),
    search: false,
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
      name: 'updated_at',
      direction: 'desc',
    },
    textLabels: {
      selectedRows: {
        text: 'filter(s) selected',
      },
    },

    onCellClick: (_, meta) =>
      meta.colIndex !== 3 && meta.colIndex !== 4 && setSelectedRowData(filters[meta.rowIndex]),

    onRowsDelete: async function handleDelete(row) {
      let response = await showmodal(Object.keys(row.lookup).length);

      if (response === 'Delete') {
        const fid = Object.keys(row.lookup).map((idx) => filters[idx]?.id);
        fid.forEach((fid) => deleteFilter(fid));
      }
      // if (response === "No")
      // fetchFilters(page, pageSize, search, sortOrder);
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          initFiltersSubscription(
            tableState.page.toString(),
            pageSize.toString(),
            search,
            sortOrder,
          );
          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          initFiltersSubscription(
            page.toString(),
            tableState.rowsPerPage.toString(),
            search,
            sortOrder,
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
            initFiltersSubscription(page.toString(), pageSize.toString(), search, order);
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

  /**
   * Gets the data of Import Filter and handles submit operation
   *
   * @param {{
   * uploadType: ("File Upload"| "URL Upload");
   * config: string;
   * name: string;
   * url: string;
   * file: string;
   * }} data
   */
  function handleImportFilter(data) {
    updateProgress({ showProgress: true });
    const { uploadType, name, config, url, file } = data;
    let requestBody = null;
    switch (uploadType) {
      case 'File Upload':
        requestBody = JSON.stringify({
          config,
          save: true,
          filter_data: {
            name,
            filter_file: getUnit8ArrayDecodedFile(file),
          },
        });
        break;
      case 'URL Upload':
        requestBody = JSON.stringify({
          config,
          save: true,
          url,
          filter_data: {
            name,
          },
        });
        break;
    }

    updateFilterFile({ updateBody: requestBody })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `"${name}" filter uploaded`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        getFilters();
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        handleError(ACTION_TYPES.UPLOAD_FILTERS);
      });
  }

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    let showCols = updateVisibleColumns(colViews, width);
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = showCols[col.name];
    });
    return initialVisibility;
  });

  const filter = {
    visibility: {
      name: 'visibility',
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

  const [selectedFilters, setSelectedFilters] = useState({ visibility: 'All' });

  const handleApplyFilter = () => {
    const visibilityFilter =
      selectedFilters.visibility === 'All' ? null : selectedFilters.visibility;
    // fetchFilters(page, pageSize, search, sortOrder, visibilityFilter);
    setVisibilityFilter(visibilityFilter);
  };

  if (isFiltersLoading) {
    return (
      <>
        <LoadingScreen animatedIcon="AnimatedFilter" message={`Loading Filters...`} />
      </>
    );
  }

  return (
    <>
      <>
        <NoSsr>
          {CAN(keys.VIEW_FILTERS.action, keys.VIEW_FILTERS.subject) ? (
            <>
              {selectedRowData && Object.keys(selectedRowData).length > 0 && (
                <YAMLEditor
                  filter={selectedRowData}
                  onClose={resetSelectedRowData()}
                  onSubmit={handleSubmit}
                />
              )}
              <ToolWrapper>
                {width < 600 && isSearchExpanded ? null : (
                  <div style={{ display: 'flex' }}>
                    {!selectedFilter.show && (filters.length > 0 || viewType === 'table') && (
                      <CreateButton>
                        <Button
                          aria-label="Add Filter"
                          variant="contained"
                          color="primary"
                          size="large"
                          onClick={handleUploadImport}
                          disabled={!CAN(keys.IMPORT_FILTER.action, keys.IMPORT_FILTER.subject)}
                        >
                          <PublishIcon style={iconMedium} data-cy="import-button" />
                          <BtnText> Import Filters </BtnText>
                        </Button>
                      </CreateButton>
                    )}
                    <ViewSwitchButton style={{ jdisplay: 'flex' }}>
                      {/* <CatalogFilter
                      catalogVisibility={catalogVisibility}
                      handleCatalogVisibility={handleCatalogVisibility}
                      classes={classes}
                    /> */}
                    </ViewSwitchButton>
                  </div>
                )}
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <SearchBar
                    onSearch={(value) => {
                      setSearch(value);
                      initFiltersSubscription(
                        page.toString(),
                        pageSize.toString(),
                        value,
                        sortOrder,
                      );
                    }}
                    expanded={isSearchExpanded}
                    setExpanded={setIsSearchExpanded}
                    placeholder="Search"
                  />
                  <UniversalFilter
                    id="ref"
                    filters={filter}
                    selectedFilters={selectedFilters}
                    setSelectedFilters={setSelectedFilters}
                    handleApplyFilter={handleApplyFilter}
                  />
                  {viewType === 'table' && (
                    <CustomColumnVisibilityControl
                      id="ref"
                      columns={columns}
                      customToolsProps={{ columnVisibility, setColumnVisibility }}
                    />
                  )}

                  {!selectedFilter.show && (
                    <ViewSwitch data-cy="table-view" view={viewType} changeView={setViewType} />
                  )}
                </div>
              </ToolWrapper>
              {!selectedFilter.show && viewType === 'table' && (
                <ResponsiveDataTable
                  data={filters}
                  columns={columns}
                  tableCols={tableCols}
                  updateCols={updateCols}
                  columnVisibility={columnVisibility}
                  options={options}
                />
              )}
              {!selectedFilter.show && viewType === 'grid' && (
                // grid view
                <FiltersGrid
                  filters={filters}
                  handleSubmit={handleSubmit}
                  canPublishFilter={canPublishFilter}
                  handlePublish={handlePublish}
                  handleUnpublishModal={handleUnpublishModal}
                  handleUploadImport={handleUploadImport}
                  handleClone={handleClone}
                  handleDownload={handleDownload}
                  uploadHandler={uploadHandler}
                  setSelectedFilter={setSelectedFilter}
                  selectedFilter={selectedFilter}
                  pages={Math.ceil(count / pageSize)}
                  setPage={setPage}
                  selectedPage={page}
                  publishModal={publishModal}
                  setPublishModal={setPublishModal}
                  publishSchema={publishSchema}
                  fetch={() => getFilters()}
                  handleInfoModal={handleInfoModal}
                />
              )}
              {canPublishFilter &&
                publishModal.open &&
                CAN(keys.PUBLISH_WASM_FILTER.action, keys.PUBLISH_WASM_FILTER.subject) && (
                  <PublishModal
                    handleClose={handlePublishModalClose}
                    title={publishModal.filter?.name}
                    handleSubmit={handlePublish}
                  />
                )}
              {importModal.open && CAN(keys.IMPORT_FILTER.action, keys.IMPORT_FILTER.subject) && (
                <ImportModal
                  handleClose={handleUploadImportClose}
                  handleImportFilter={handleImportFilter}
                />
              )}
              {infoModal.open &&
                CAN(keys.DETAILS_OF_WASM_FILTER.action, keys.DETAILS_OF_WASM_FILTER.subject) && (
                  <InfoModal
                    handlePublish={handlePublish}
                    infoModalOpen={true}
                    handleInfoModalClose={handleInfoModalClose}
                    dataName="filters"
                    selectedResource={infoModal.selectedResource}
                    resourceOwnerID={infoModal.ownerID}
                    currentUser={user}
                    formSchema={publishSchema}
                    meshModels={meshModels}
                    patternFetcher={getFilters}
                  />
                )}
              <_PromptComponent ref={modalRef} />
            </>
          ) : (
            <DefaultError />
          )}
        </NoSsr>
      </>
    </>
  );
}

const ImportModal = React.memo((props) => {
  const { handleClose, handleImportFilter } = props;

  return (
    <>
      <SistentModal
        open={true}
        closeModal={handleClose}
        headerIcon={
          <Filter fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
        }
        title="Import Design"
        maxWidth="sm"
      >
        <RJSFModalWrapper
          schema={importFilterSchema}
          uiSchema={importFilterUiSchema}
          handleSubmit={handleImportFilter}
          submitBtnText="Import"
          handleClose={handleClose}
        />
      </SistentModal>
    </>
  );
});

const PublishModal = React.memo((props) => {
  const { handleClose, handleSubmit, title } = props;

  return (
    <>
      <SistentModal
        open={true}
        headerIcon={
          <Filter fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
        }
        closeModal={handleClose}
        aria-label="catalog publish"
        title={title}
        maxWidth="sm"
      >
        <RJSFModalWrapper
          schema={publishCatalogItemSchema}
          uiSchema={publishCatalogItemUiSchema}
          submitBtnText="Submit for Approval"
          handleSubmit={handleSubmit}
          helpText="Upon submitting your catalog item, an approval flow will be initiated.[Learn more](https://docs.meshery.io/concepts/catalog)"
          handleClose={handleClose}
        />
      </SistentModal>
    </>
  );
});

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  toggleCatalogContent: bindActionCreators(toggleCatalogContent, dispatch),
});

const mapStateToProps = (state) => {
  return {
    user: state.get('user')?.toObject(),
    selectedK8sContexts: state.get('selectedK8sContexts'),
    catalogVisibility: state.get('catalogVisibility'),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(MesheryFilters);
