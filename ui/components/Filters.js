import React, { useState, useEffect, useRef } from 'react';
import { withStyles } from '@material-ui/core/styles';
// import { createTheme } from '@material-ui/core/styles';
import {
  NoSsr,
  TableCell,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Divider,
  Tooltip,
  Typography,
  Button,
} from '@material-ui/core';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import DeleteIcon from '@material-ui/icons/Delete';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Moment from 'react-moment';
import CloseIcon from '@material-ui/icons/Close';
import EditIcon from '@material-ui/icons/Edit';
import { toggleCatalogContent, updateProgress } from '../lib/store';
import TableSortLabel from '@material-ui/core/TableSortLabel';
import dataFetch from '../lib/data-fetch';
import PromptComponent from './PromptComponent';
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import { FILE_OPS, MesheryFiltersCatalog, VISIBILITY } from '../utils/Enum';
import ViewSwitch from './ViewSwitch';
import CatalogFilter from './CatalogFilter';
import FiltersGrid from './MesheryFilters/FiltersGrid';
import { trueRandom } from '../lib/trueRandom';
import GetAppIcon from '@material-ui/icons/GetApp';
import PublicIcon from '@material-ui/icons/Public';
import { ctxUrl } from '../utils/multi-ctx';
import ConfirmationMsg from './ConfirmationModal';
import PublishIcon from '@material-ui/icons/Publish';
import downloadFile from '../utils/fileDownloader';
import CloneIcon from '../public/static/img/CloneIcon';
import SaveIcon from '@material-ui/icons/Save';
import ConfigurationSubscription from './graphql/subscriptions/ConfigurationSubscription';
import fetchCatalogFilter from './graphql/queries/CatalogFilterQuery';
import { iconMedium } from '../css/icons.styles';
import Modal from './Modal';
import { getUnit8ArrayDecodedFile, modifyRJSFSchema } from '../utils/utils';
import Filter from '../public/static/img/drawer-icons/filter_svg.js';
import { getMeshModels } from '../api/meshmodel';
import _ from 'lodash';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import SearchBar from '../utils/custom-search';
import CustomColumnVisibilityControl from '../utils/custom-column';
import { ResponsiveDataTable } from '@layer5/sistent-components';
import useStyles from '../assets/styles/general/tool.styles';
import { updateVisibleColumns } from '../utils/responsive-column';
import { useWindowDimensions } from '../utils/dimension';
import { Box } from '@mui/material';
import InfoModal from './Modals/Information/InfoModal';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';

const styles = (theme) => ({
  grid: {
    padding: theme.spacing(2),
  },
  tableHeader: {
    fontWeight: 'bolder',
    fontSize: 18,
  },
  createButton: {
    width: 'fit-content',
    alignSelf: 'flex-start',
  },
  viewSwitchButton: {
    justifySelf: 'flex-end',
    paddingLeft: '1rem',
  },
  ymlDialogTitle: {
    display: 'flex',
    alignItems: 'center',
  },
  searchWrapper: {
    justifySelf: 'flex-end',
    marginLeft: 'auto',
    paddingLeft: '1rem',
    display: 'flex',
  },

  ymlDialogTitleText: {
    flexGrow: 1,
  },
  fullScreenCodeMirror: {
    height: '100%',
    '& .CodeMirror': {
      minHeight: '300px',
      height: '100%',
    },
  },
  visibilityImg: {
    filter: theme.palette.secondary.img,
  },
  btnText: {
    display: 'block',
    '@media (max-width: 1450px)': {
      display: 'none',
    },
  },
});

function TooltipIcon({ children, onClick, title }) {
  return (
    <Tooltip title={title} placement="top" arrow interactive>
      <IconButton onClick={onClick}>{children}</IconButton>
    </Tooltip>
  );
}

function YAMLEditor({ filter, onClose, onSubmit, classes }) {
  const [yaml, setYaml] = useState('');
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  let resourceData;
  try {
    resourceData = JSON.parse(filter.filter_resource);
  } catch (error) {
    // Handling the error or provide a default value
    console.error('Error parsing JSON:', error);
    resourceData = {}; // Setting a default value if parsing fails
  }

  const config = resourceData?.settings?.config || '';

  return (
    <Dialog
      onClose={onClose}
      aria-labelledby="filter-dialog-title"
      open
      maxWidth="md"
      fullScreen={fullScreen}
      fullWidth={!fullScreen}
    >
      <DialogTitle disableTypography id="filter-dialog-title" className={classes.ymlDialogTitle}>
        <Typography variant="h6" className={classes.ymlDialogTitleText}>
          {filter.name}
        </Typography>
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
      <Divider variant="fullWidth" light />
      <DialogContent>
        <CodeMirror
          value={config}
          className={fullScreen ? classes.fullScreenCodeMirror : ''}
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
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <Tooltip title="Update Filter">
          <IconButton
            aria-label="Update"
            color="primary"
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
        </Tooltip>
        <Tooltip title="Delete Filter">
          <IconButton
            aria-label="Delete"
            color="primary"
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
        </Tooltip>
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
  classes,
  selectedK8sContexts,
  catalogVisibility,
  toggleCatalogContent,
}) {
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder] = useState('');
  const [count, setCount] = useState(0);
  const modalRef = useRef(null);
  const [pageSize, setPageSize] = useState(10);
  const [filters, setFilters] = useState([]);
  const [selectedFilter, setSelectedFilter] = useState(resetSelectedFilter());
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [extensionPreferences, setExtensionPreferences] = useState({});
  const [canPublishFilter, setCanPublishFilter] = useState(false);
  const [importSchema, setImportSchema] = useState({});
  const [publishSchema, setPublishSchema] = useState({});
  const { width } = useWindowDimensions();
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ('grid'),
  );
  const FILTER_URL = '/api/filter';
  const DEPLOY_URL = FILTER_URL + '/deploy';
  const CLONE_URL = '/clone';

  //hooks
  const { notify } = useNotification();
  const StyleClass = useStyles();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);

  const [infoModal, setInfoModal] = useState({
    open: false,
    ownerID: '',
    selectedResource: {},
  });

  const [modalOpen, setModalOpen] = useState({
    open: false,
    filter_file: null,
    deploy: false,
    name: '',
    count: 0,
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
  useEffect(() => {
    dataFetch(
      '/api/schema/resource/filter',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        setImportSchema(result);
      },
      handleError(ACTION_TYPES.SCHEMA_FETCH),
    );
    dataFetch(
      '/api/schema/resource/publish',
      {
        method: 'GET',
        credentials: 'include',
      },
      async (result) => {
        try {
          const { models } = await getMeshModels();
          const modelNames = _.uniq(models?.map((model) => model.displayName.toUpperCase()));
          modelNames.sort();

          // Modify the schema using the utility function
          const modifiedSchema = modifyRJSFSchema(
            result.rjsfSchema,
            'properties.compatibility.items.enum',
            modelNames,
          );
          setPublishSchema({ rjsfSchema: modifiedSchema, uiSchema: result.uiSchema });
        } catch (err) {
          console.error(err);
          setPublishSchema(result);
        }
      },
      handleError(ACTION_TYPES.SCHEMA_FETCH),
    );
    dataFetch(
      '/api/provider/capabilities',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (result) {
          const capabilitiesRegistry = result;
          const filtersCatalogueCapability = capabilitiesRegistry?.capabilities.filter(
            (val) => val.feature === MesheryFiltersCatalog,
          );
          if (filtersCatalogueCapability?.length > 0) setCanPublishFilter(true);
        }
      },
      (err) => console.error(err),
    );
  }, []);

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
          options: ['Yes', 'No'],
        });
        if (response === 'Yes') {
          updateProgress({ showProgress: true });
          dataFetch(
            `/api/filter/catalog/unpublish`,
            { credentials: 'include', method: 'DELETE', body: JSON.stringify({ id: filter?.id }) },
            () => {
              updateProgress({ showProgress: false });
              notify({
                message: `"${filter?.name}" filter unpublished`,
                event_type: EVENT_TYPES.SUCCESS,
              });
            },
            handleError(ACTION_TYPES.UNPUBLISH_CATALOG),
          );
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
    fetchFilters(page, pageSize, search, sortOrder);
  }, [page, pageSize, search, sortOrder]);

  useEffect(() => {
    if (viewType === 'grid') setSearch('');
  }, [viewType]);

  const handleCatalogPreference = (catalogPref) => {
    let body = Object.assign({}, extensionPreferences);
    body['catalogContent'] = catalogPref;

    dataFetch(
      '/api/user/prefs',
      {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({ usersExtensionPreferences: body }),
      },
      () => {
        notify({
          message: `Catalog Content was ${catalogPref ? 'enab' : 'disab'}led`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      },
      (err) => console.error(err),
    );
  };

  const fetchUserPrefs = () => {
    dataFetch(
      '/api/user/prefs',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (result) {
          setExtensionPreferences(result?.usersExtensionPreferences);
        }
      },
      (err) => console.error(err),
    );
  };

  const handleCatalogVisibility = () => {
    handleCatalogPreference(!catalogVisibilityRef.current);
    catalogVisibilityRef.current = !catalogVisibility;
    toggleCatalogContent({ catalogVisibility: !catalogVisibility });
  };

  useEffect(() => {
    fetchUserPrefs();
    handleSetFilters(filters);
  }, [catalogVisibility]);

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

  /**
   * fetchFilters constructs the queries based on the parameters given
   * and fetches the filters
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */
  function fetchFilters(page, pageSize, search, sortOrder) {
    if (!search) search = '';
    if (!sortOrder) sortOrder = '';

    const query = `?page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
      search,
    )}&order=${encodeURIComponent(sortOrder)}`;

    updateProgress({ showProgress: true });

    dataFetch(
      `/api/filter${query}`,
      { credentials: 'include' },
      (result) => {
        console.log('FilterFile API', `/api/filter${query}`);
        updateProgress({ showProgress: false });
        if (result) {
          handleSetFilters(result.filters || []);
          // setPage(result.page || 0);
          setPageSize(result.page_size || 0);
          setCount(result.total_count || 0);
        }
      },
      // handleError
      handleError(ACTION_TYPES.FETCH_FILTERS),
    );
  }

  const handleDeploy = (filter_file, name) => {
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      { credentials: 'include', method: 'POST', body: filter_file },
      () => {
        console.log('FilterFile Deploy API', `/api/filter/deploy`);
        notify({ message: `"${name}" filter deployed`, event_type: EVENT_TYPES.SUCCESS });
        updateProgress({ showProgress: false });
      },
      handleError(ACTION_TYPES.DEPLOY_FILTERS),
    );
  };

  const handleUndeploy = (filter_file, name) => {
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      { credentials: 'include', method: 'DELETE', body: filter_file },
      () => {
        updateProgress({ showProgress: false });
        notify({ message: `"${name}" filter undeployed`, event_type: EVENT_TYPES.SUCCESS });
      },
      handleError(ACTION_TYPES.UNDEPLOY_FILTERS),
    );
  };

  const handlePublish = (formData) => {
    const payload = {
      id: publishModal.filter?.id,
      catalog_data: formData,
    };
    updateProgress({ showProgress: true });
    dataFetch(
      `/api/filter/catalog/publish`,
      { credentials: 'include', method: 'POST', body: JSON.stringify(payload) },
      () => {
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
      },
      handleError(ACTION_TYPES.PUBLISH_CATALOG),
    );
  };

  function handleClone(filterID, name) {
    updateProgress({ showProgress: true });
    dataFetch(
      FILTER_URL.concat(CLONE_URL, '/', filterID),
      {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify({ name: name + ' (Copy)' }),
      },
      () => {
        updateProgress({ showProgress: false });
        notify({ message: `"${name}" filter cloned`, event_type: EVENT_TYPES.SUCCESS });
      },
      handleError(ACTION_TYPES.CLONE_FILTERS),
    );
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
      (result) => {
        setPage(result.configuration?.filters?.page || 0);
        setPageSize(result.configuration?.filters?.page_size || 0);
        setCount(result.configuration?.filters?.total_count || 0);
        handleSetFilters(result.configuration?.filters?.filters);
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

  const handleModalClose = () => {
    setModalOpen({
      open: false,
      filter_file: null,
      name: '',
      count: 0,
    });
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
      if (response == 'No') {
        updateProgress({ showProgress: false });
        return;
      }
      dataFetch(
        `/api/filter/${id}`,
        { credentials: 'include', method: 'DELETE' },
        () => {
          updateProgress({ showProgress: false });
          notify({ message: `"${name}" filter deleted`, event_type: EVENT_TYPES.SUCCESS });
          resetSelectedRowData()();
        },
        // handleError
        handleError(ACTION_TYPES.DELETE_FILTERS),
      );
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
      dataFetch(
        `/api/filter`,
        {
          credentials: 'include',
          headers: {
            'Content-Type': 'application/octet-stream', // Set appropriate content type for binary data
          },
          method: 'POST',
          body,
        },
        () => {
          updateProgress({ showProgress: false });
        },
        // handleError
        handleError(ACTION_TYPES.UPLOAD_FILTERS),
      );
    }

    if (type === FILE_OPS.UPDATE) {
      dataFetch(
        `/api/filter`,
        {
          credentials: 'include',
          method: 'POST',
          body: JSON.stringify({
            filter_data: { id, name: name, catalog_data },
            config: data,
            save: true,
          }),
        },
        () => {
          updateProgress({ showProgress: false });
        },
        handleError(ACTION_TYPES.UPLOAD_FILTERS),
      );
    }
  }

  const handleDownload = (e, id, name) => {
    e.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      downloadFile({ id, name, type: 'filter' });
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
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'created_at',
      label: 'Upload Timestamp',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'updated_at',
      label: 'Update Timestamp',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
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
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          const visibility = filters[tableMeta.rowIndex]?.visibility;
          return (
            <>
              <img className={classes.visibilityImg} src={`/static/img/${visibility}.svg`} />
            </>
          );
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
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          const rowData = filters[tableMeta.rowIndex];
          const visibility = filters[tableMeta.rowIndex]?.visibility;
          return (
            <Box
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
                >
                  <CloneIcon fill="currentColor" className={classes.iconPatt} />
                </TooltipIcon>
              ) : (
                <TooltipIcon
                  title="Config"
                  onClick={(e) => {
                    e.stopPropagation();
                    setSelectedRowData(filters[tableMeta.rowIndex]);
                  }}
                >
                  <EditIcon aria-label="config" color="inherit" style={iconMedium} />
                </TooltipIcon>
              )}
              <TooltipIcon
                title="Download"
                onClick={(e) => handleDownload(e, rowData.id, rowData.name)}
              >
                <GetAppIcon data-cy="download-button" />
              </TooltipIcon>
              <TooltipIcon title="Filter Information" onClick={() => handleInfoModal(rowData)}>
                <InfoOutlinedIcon data-cy="information-button" />
              </TooltipIcon>
              {canPublishFilter && visibility !== VISIBILITY.PUBLISHED ? (
                <TooltipIcon title="Publish" onClick={(ev) => handlePublishModal(ev, rowData)}>
                  <PublicIcon fill="#F91313" data-cy="publish-button" />
                </TooltipIcon>
              ) : (
                <TooltipIcon
                  title="Unpublish"
                  onClick={(ev) => handleUnpublishModal(ev, rowData)()}
                >
                  <PublicIcon fill="#F91313" data-cy="unpublish-button" />
                </TooltipIcon>
              )}
            </Box>
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

      options: ['Yes', 'No'],
    });
    return response;
  }

  function deleteFilter(id) {
    dataFetch(
      `/api/filter/${id}`,
      {
        method: 'DELETE',
        credentials: 'include',
      },
      () => {
        updateProgress({ showProgress: false });
        notify({ message: `Filter deleted`, event_type: EVENT_TYPES.SUCCESS });
      },
      handleError('Failed To Delete Filter'),
    );
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
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    page,
    print: false,
    download: false,
    textLabels: {
      selectedRows: {
        text: 'filter(s) selected',
      },
    },

    onCellClick: (_, meta) =>
      meta.colIndex !== 3 && meta.colIndex !== 4 && setSelectedRowData(filters[meta.rowIndex]),

    onRowsDelete: async function handleDelete(row) {
      let response = await showmodal(Object.keys(row.lookup).length);
      console.log(response);
      if (response === 'Yes') {
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
          break;
        case 'changeRowsPerPage':
          initFiltersSubscription(
            page.toString(),
            tableState.rowsPerPage.toString(),
            search,
            sortOrder,
          );
          break;
        case 'search':
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              fetchFilters(
                page,
                pageSize,
                tableState.searchText !== null ? tableState.searchText : '',
                sortOrder,
              );
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

    dataFetch(
      '/api/filter',
      { credentials: 'include', method: 'POST', body: requestBody },
      () => {
        updateProgress({ showProgress: false });
      },
      handleError(ACTION_TYPES.UPLOAD_FILTERS),
    );
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

  return (
    <>
      <NoSsr>
        {selectedRowData && Object.keys(selectedRowData).length > 0 && (
          <YAMLEditor
            filter={selectedRowData}
            onClose={resetSelectedRowData()}
            onSubmit={handleSubmit}
            classes={classes}
          />
        )}
        <div className={StyleClass.toolWrapper}>
          {width < 600 && isSearchExpanded ? null : (
            <div style={{ display: 'flex' }}>
              {!selectedFilter.show && (filters.length > 0 || viewType === 'table') && (
                <div className={classes.createButton}>
                  <div>
                    <Button
                      aria-label="Add Filter"
                      variant="contained"
                      color="primary"
                      size="large"
                      // @ts-ignore
                      onClick={handleUploadImport}
                      style={{ marginRight: '2rem' }}
                    >
                      <PublishIcon
                        style={iconMedium}
                        className={classes.addIcon}
                        data-cy="import-button"
                      />
                      <span className={classes.btnText}> Import Filters </span>
                    </Button>
                  </div>
                </div>
              )}
              <div style={{ jdisplay: 'flex' }}>
                <CatalogFilter
                  catalogVisibility={catalogVisibility}
                  handleCatalogVisibility={handleCatalogVisibility}
                  classes={classes}
                />
              </div>
            </div>
          )}
          <div className={classes.searchWrapper} style={{ display: 'flex' }}>
            <SearchBar
              onSearch={(value) => {
                setSearch(value);
                initFiltersSubscription(page.toString(), pageSize.toString(), value, sortOrder);
              }}
              expanded={isSearchExpanded}
              setExpanded={setIsSearchExpanded}
              placeholder="Search"
            />
            {viewType === 'table' && (
              <CustomColumnVisibilityControl
                columns={columns}
                customToolsProps={{ columnVisibility, setColumnVisibility }}
              />
            )}

            {!selectedFilter.show && (
              <ViewSwitch data-cy="table-view" view={viewType} changeView={setViewType} />
            )}
          </div>
        </div>
        {!selectedFilter.show && viewType === 'table' && (
          <ResponsiveDataTable
            data={filters}
            columns={columns}
            tableCols={tableCols}
            updateCols={updateCols}
            columnVisibility={columnVisibility}
            // @ts-ignore
            options={options}
            className={classes.muiRow}
          />
        )}
        {!selectedFilter.show && viewType === 'grid' && (
          // grid view
          <FiltersGrid
            filters={filters}
            handleDeploy={handleDeploy}
            handleUndeploy={handleUndeploy}
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
            importSchema={importSchema}
            setPage={setPage}
            selectedPage={page}
            publishModal={publishModal}
            setPublishModal={setPublishModal}
            publishSchema={publishSchema}
            fetch={() => fetchFilters(page, pageSize, search, sortOrder)}
            handleInfoModal={handleInfoModal}
          />
        )}
        <ConfirmationMsg
          open={modalOpen.open}
          handleClose={handleModalClose}
          submit={{
            deploy: () => handleDeploy(modalOpen.filter_file, modalOpen.name),
            unDeploy: () => handleUndeploy(modalOpen.filter_file, modalOpen.name),
          }}
          isDelete={!modalOpen.deploy}
          title={modalOpen.name}
          componentCount={modalOpen.count}
          tab={modalOpen.deploy ? 2 : 1}
        />
        {canPublishFilter && publishModal.open && (
          <Modal
            open={true}
            schema={publishSchema.rjsfSchema}
            uiSchema={publishSchema.uiSchema}
            title={publishModal.filter?.name}
            handleClose={handlePublishModalClose}
            handleSubmit={handlePublish}
            showInfoIcon={{
              text: 'Upon submitting your catalog item, an approval flow will be initiated.',
              link: 'https://docs.meshery.io/concepts/catalog',
            }}
            submitBtnText="Submit for Approval"
            submitBtnIcon={
              <PublicIcon style={iconMedium} className={classes.addIcon} data-cy="import-button" />
            }
          />
        )}
        <PromptComponent ref={modalRef} />
        {importModal.open && (
          <Modal
            open={true}
            schema={importSchema.rjsfSchema}
            uiSchema={importSchema.uiSchema}
            handleClose={handleUploadImportClose}
            handleSubmit={handleImportFilter}
            title="Import Filter"
            submitBtnText="Import"
            leftHeaderIcon={
              <Filter fill="#fff" style={{ height: '24px', width: '24px', fonSize: '1.45rem' }} />
            }
            submitBtnIcon={<PublishIcon />}
          />
        )}
        {infoModal.open && (
          <InfoModal
            infoModalOpen={true}
            handleInfoModalClose={handleInfoModalClose}
            dataName="filters"
            selectedResource={infoModal.selectedResource}
            resourceOwnerID={infoModal.ownerID}
            currentUserID={user?.id}
            formSchema={publishSchema}
          />
        )}
      </NoSsr>
    </>
  );
}

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

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(MesheryFilters));
