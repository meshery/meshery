import React, { useState, useEffect, useRef } from "react";
import { withStyles } from "@material-ui/core/styles";
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
  Button
} from "@material-ui/core";
import { UnControlled as CodeMirror } from "react-codemirror2";
import DeleteIcon from "@material-ui/icons/Delete";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import EditIcon from "@material-ui/icons/Edit";
import { toggleCatalogContent, updateProgress } from "../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import dataFetch from "../lib/data-fetch";
import PromptComponent from "./PromptComponent";
import UploadImport from "./UploadImport";
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import { FILE_OPS, MesheryFiltersCatalog, VISIBILITY } from "../utils/Enum";
import ViewSwitch from "./ViewSwitch";
import CatalogFilter from "./CatalogFilter";
import FiltersGrid from "./MesheryFilters/FiltersGrid";
import { trueRandom } from "../lib/trueRandom";
import GetAppIcon from '@material-ui/icons/GetApp';
import PublicIcon from '@material-ui/icons/Public';
import { ctxUrl } from "../utils/multi-ctx";
import ConfirmationMsg from "./ConfirmationModal";
import PublishIcon from "@material-ui/icons/Publish";
import downloadFile from "../utils/fileDownloader";
import CloneIcon from "../public/static/img/CloneIcon";
import SaveIcon from "@material-ui/icons/Save";
import ConfigurationSubscription from "./graphql/subscriptions/ConfigurationSubscription";
import fetchCatalogFilter from "./graphql/queries/CatalogFilterQuery";
import LoadingScreen from "./LoadingComponents/LoadingComponent";
import { iconMedium } from "../css/icons.styles";
import Modal from "./Modal";
import { publish_schema, publish_ui_schema } from "./schemas/publish_schema";
import _ from "lodash";
import SearchBar from "./searchcommon";

const styles = (theme) => ({
  grid : {
    padding : theme.spacing(2),
  },
  tableHeader : {
    fontWeight : "bolder",
    fontSize : 18,
  },
  createButton : {
    display : "flex",
    justifyContent : "flex-start",
    alignItems : "center",
    whiteSpace : "nowrap",
  },
  topToolbar : {
    margin : "2rem auto",
    display : "flex",
    justifyContent : "space-between",
    paddingLeft : "1rem"
  },
  viewSwitchButton : {
    justifySelf : "flex-end",
    paddingLeft : "1rem"
  },
  ymlDialogTitle : {
    display : "flex",
    alignItems : "center"
  },
  ymlDialogTitleText : {
    flexGrow : 1
  },
  fullScreenCodeMirror : {
    height : '100%',
    '& .CodeMirror' : {
      minHeight : "300px",
      height : '100%',
    }
  },
  visibilityImg : {
    filter : theme.palette.secondary.img,
  }
  // text : {
  //   padding : "5px"
  // }
});

function TooltipIcon({ children, onClick, title }) {
  return (
    <Tooltip title={title} placement="top" arrow interactive >
      <IconButton onClick={onClick}>
        {children}
      </IconButton>
    </Tooltip>
  )
}

function YAMLEditor({ filter, onClose, onSubmit, classes }) {
  const [yaml, setYaml] = useState("");
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  }

  const resourceData = JSON.parse(filter.filter_resource);
  const config = resourceData.settings.config;

  return (
    <Dialog onClose={onClose} aria-labelledby="filter-dialog-title" open maxWidth="md" fullScreen={fullScreen} fullWidth={!fullScreen}>
      <DialogTitle disableTypography id="filter-dialog-title" className={classes.ymlDialogTitle}>
        <Typography variant="h6" className={classes.ymlDialogTitleText}>
          {filter.name}
        </Typography>
        <TooltipIcon
          title={fullScreen ? "Exit Fullscreen" : "Enter Fullscreen"}
          onClick={toggleFullScreen}>
          {fullScreen ? <FullscreenExitIcon  style={iconMedium} /> : <FullscreenIcon  style={iconMedium} />}
        </TooltipIcon>
        <TooltipIcon title="Exit" onClick={onClose}>
          <CloseIcon style={iconMedium} />
        </TooltipIcon>
      </DialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <CodeMirror
          value={config}
          className={fullScreen ? classes.fullScreenCodeMirror : ""}
          options={{
            theme : "material",
            lineNumbers : true,
            lineWrapping : true,
            gutters : ["CodeMirror-lint-markers"],
            lint : true,
            mode : "text/x-yaml",
          }}
          onChange={(_,data,val) => setYaml(val)}
        />
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <Tooltip title="Update Filter">
          <IconButton
            aria-label="Update"
            color="primary"
            onClick={() => onSubmit({
              data : yaml,
              id : filter.id,
              name : filter.name,
              type : FILE_OPS.UPDATE
            })}
          >
            <SaveIcon  style={iconMedium} />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Filter">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={() => onSubmit({
              data : yaml,
              id : filter.id,
              name : filter.name,
              type : FILE_OPS.DELETE
            })}
          >
            <DeleteIcon  style={iconMedium} />
          </IconButton>
        </Tooltip>
      </DialogActions>
    </Dialog>
  );
}

function resetSelectedFilter() {
  return { show : false, filter : null };
}

function MesheryFilters({ updateProgress, enqueueSnackbar, closeSnackbar, user, classes, selectedK8sContexts, catalogVisibility, toggleCatalogContent }) {
  const [page, setPage] = useState(0);
  const [search,setSearch] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const modalRef = useRef(null);
  const [pageSize, setPageSize] = useState(10);
  const [filters, setFilters] = useState([]);
  const [selectedFilter, setSelectedFilter] = useState(resetSelectedFilter());
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [extensionPreferences, setExtensionPreferences] = useState({});
  const [canPublishFilter, setCanPublishFilter] = useState(false);
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ("grid")
  );
  const FILTER_URL = "/api/filter"
  const DEPLOY_URL = FILTER_URL + "/deploy";
  const CLONE_URL = "/clone";

  const [modalOpen, setModalOpen] = useState({
    open : false,
    filter_file : null,
    deploy : false,
    name : "",
    count : 0
  });

  const [importModal, setImportModal] = useState({
    open : false
  })
  const [publishModal, setPublishModal] = useState({
    open : false,
    filter : {},
    name : "",
  });
  const [payload, setPayload] = useState({
    id : "",
    catalog_data : {}
  });

  const [loading, stillLoading] = useState(true);

  const catalogContentRef = useRef();
  const catalogVisibilityRef = useRef();
  const disposeConfSubscriptionRef = useRef(null);


  const ACTION_TYPES = {
    FETCH_FILTERS : {
      name : "FETCH_FILTERS",
      error_msg : "Failed to fetch filter",
    },
    DELETE_FILTERS : {
      name : "DELETE_FILTERS",
      error_msg : "Failed to delete filter file",
    },
    DEPLOY_FILTERS : {
      name : "DEPLOY_FILTERS",
      error_msg : "Failed to deploy filter file",
    },
    UNDEPLOY_FILTERS : {
      name : "UNDEPLOY_FILTERS",
      error_msg : "Failed to undeploy filter file",
    },
    UPLOAD_FILTERS : {
      name : "UPLOAD_FILTERS",
      error_msg : "Failed to upload filter file",
    },
    CLONE_FILTERS : {
      name : "CLONE_FILTER",
      error_msg : "Failed to clone filter file"
    },
    PUBLISH_CATALOG : {
      name : "PUBLISH_CATALOG",
      error_msg : "Failed to publish catalog"
    },
    UNPUBLISH_CATALOG : {
      name : "PUBLISH_CATALOG",
      error_msg : "Failed to publish catalog"
    }
  };

  /**
   * Checking whether users are signed in under a provider that doesn't have
   * publish filter capability and setting the canPublishFilter state accordingly
   */
  useEffect(() => {
    dataFetch(
      "/api/provider/capabilities",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        if (result) {
          const capabilitiesRegistry = result;
          const filtersCatalogueCapability = capabilitiesRegistry?.capabilities.filter((val) => val.feature === MesheryFiltersCatalog);
          if (filtersCatalogueCapability?.length > 0) setCanPublishFilter(true);
        }
      },
      (err) => console.error(err)
    );
  }, [])

  const searchTimeout = useRef(null);

  const onChange = (e) => {
    setPayload({
      id : publishModal.filter?.id,
      catalog_data : e
    })
  }

  const handleUploadImport = () => {
    setImportModal({
      open : true
    });
  }

  const handleUploadImportClose = () => {
    setImportModal({
      open : false
    });
  }

  const handlePublishModal = (ev, filter) => {
    if (canPublishFilter) {
      ev.stopPropagation();
      setPublishModal({
        open : true,
        filter : filter
      });
    }
  };

  const handleUnpublishModal = (ev, filter) => {
    if (canPublishFilter) {
      return async () => {
        let response = await modalRef.current.show({
          title : `Unpublish Catalog item?`,
          subtitle : `Are you sure that you want to unpublish "${filter?.name}"?`,
          options : ["Yes", "No"]
        });
        if (response === "Yes") {
          updateProgress({ showProgress : true });
          dataFetch(
                `/api/filter/catalog/unpublish`,
                { credentials : "include", method : "DELETE", body : JSON.stringify({ "id" : filter?.id }) },
                () => {
                  updateProgress({ showProgress : false });
                  enqueueSnackbar((`"${filter?.name}" filter unpublished`), {
                    variant : "success",
                    action : function Action(key) {
                      return (
                        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                          <CloseIcon />
                        </IconButton>
                      );
                    },
                    autoHideDuration : 2000,
                  });
                },
                handleError(ACTION_TYPES.UNPUBLISH_CATALOG),
          );
        }
      }
    }
  };

  const handlePublishModalClose = () => {
    setPublishModal({
      open : false,
      filter : {},
      name : ""
    });
  };

  useEffect(() => {
    fetchFilters(page,pageSize,search,sortOrder)
  }, [page, pageSize, search, sortOrder]);

  useEffect(() => {
    if (viewType==='grid')setSearch("")
  },[viewType])

  const handleCatalogPreference = (catalogPref) => {
    let body = Object.assign({}, extensionPreferences)
    body["catalogContent"] = catalogPref

    dataFetch(
      "/api/user/prefs",
      {
        method : "POST",
        credentials : "include",
        body : JSON.stringify({ usersExtensionPreferences : body })
      },
      () => {
        enqueueSnackbar(`Catalog Content was ${catalogPref ? "enab" : "disab"}led`,
          { variant : 'success',
            autoHideDuration : 4000,
            action : (key) => (
              <IconButton
                key="close"
                aria-label="Close"
                color="inherit"
                onClick={() => closeSnackbar(key)}
              >
                <CloseIcon style={iconMedium} />
              </IconButton>
            ),
          });
      },
      err => console.error(err),
    )
  }

  const fetchUserPrefs = () => {
    dataFetch(
      "/api/user/prefs",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        if (result) {
          setExtensionPreferences(result?.usersExtensionPreferences)
        }
      },
      err => console.error(err)
    )
  }

  const handleCatalogVisibility = () => {
    handleCatalogPreference(!catalogVisibilityRef.current);
    catalogVisibilityRef.current = !catalogVisibility
    toggleCatalogContent({ catalogVisibility : !catalogVisibility });
  }

  useEffect(() => {
    fetchUserPrefs();
    handleSetFilters(filters)
  }, [catalogVisibility])

  useEffect(() => {
    catalogVisibilityRef.current = catalogVisibility
    const fetchCatalogFilters = fetchCatalogFilter({
      selector : {
        search : "",
        order : "",
        page : 0,
        pagesize : 0,
      }
    }).subscribe({
      next : (result) => {
        catalogContentRef.current = result?.catalogFilters;
        initFiltersSubscription();
      },
      error : (err) => console.log("There was an error fetching Catalog Filter: ", err)
    });

    return () => {
      fetchCatalogFilters.unsubscribe();
      disposeConfSubscriptionRef.current?.dispose();
    }
  }, [])

  /**
   * fetchFilters constructs the queries based on the parameters given
   * and fetches the filters
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */
  function fetchFilters(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;

    updateProgress({ showProgress : true });

    dataFetch(
      `/api/filter${query}`,
      { credentials : "include" },
      (result) => {
        console.log("FilterFile API", `/api/filter${query}`);
        updateProgress({ showProgress : false });
        if (result) {
          handleSetFilters(result.filters || []);
          // setPage(result.page || 0);
          setPageSize(result.page_size || 0);
          setCount(result.total_count || 0);
        }
      },
      // handleError
      handleError(ACTION_TYPES.FETCH_FILTERS)
    );
  }

  const handleDeploy = (filter_file, name) => {
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      { credentials : "include", method : "POST", body : filter_file },
      () => {
        console.log("FilterFile Deploy API", `/api/filter/deploy`);
        enqueueSnackbar(`"${name}" filter deployed`, {
          variant : "success",
          action : function Action(key) {
            return (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon style={iconMedium} />
              </IconButton>
            );
          },
          autoHideDuration : 2000,
        });
        updateProgress({ showProgress : false });
      },
      handleError(ACTION_TYPES.DEPLOY_FILTERS)
    );
  };

  const handleUndeploy = (filter_file, name) => {
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      { credentials : "include", method : "DELETE", body : filter_file },
      () => {
        updateProgress({ showProgress : false });
        enqueueSnackbar(`"${name}" filter undeployed`, {
          variant : "success",
          action : function Action(key) {
            return (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon style={iconMedium} />
              </IconButton>
            );
          },
          autoHideDuration : 2000,
        });
      },
      handleError(ACTION_TYPES.UNDEPLOY_FILTERS)
    );
  };

  const handlePublish = (catalog_data) => {
    updateProgress({ showProgress : true });
    dataFetch(
      `/api/filter/catalog/publish`,
      { credentials : "include", method : "POST", body : JSON.stringify(catalog_data) },
      () => {
        updateProgress({ showProgress : false });
        enqueueSnackbar("Filter Published!", {
          variant : "success",
          action : function Action(key) {
            return (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            );
          },
          autoHideDuration : 2000,
        });
      },
      handleError(ACTION_TYPES.PUBLISH_CATALOG),
    );
  }

  function handleClone(filterID, name) {
    updateProgress({ showProgress : true });
    dataFetch(FILTER_URL.concat(CLONE_URL, "/", filterID),
      {
        credentials : "include",
        method : "POST",
        body : JSON.stringify({ name : name + " (Copy)" }),
      },
      () => {
        updateProgress({ showProgress : false });
        enqueueSnackbar(`"${name}" filter cloned`, {
          variant : "success",
          action : function Action(key) {
            return (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon  style={iconMedium}/>
              </IconButton>
            );
          },
          autoHideDuration : 2000,
        });
      },
      handleError(ACTION_TYPES.CLONE_FILTERS),
    );
  }

  // function handleError(error) {
  const handleError = (action) => (error) => {
    updateProgress({ showProgress : false });

    enqueueSnackbar(`${action.error_msg}: ${error}`, {
      variant : "error",
      action : function Action(key) {
        return (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon  style={iconMedium}/>
          </IconButton>
        );
      },
      autoHideDuration : 8000,
    });
  };

  const handleSetFilters = (filters) => {
    if (catalogVisibilityRef.current && catalogContentRef.current?.length > 0) {
      setFilters([...catalogContentRef.current, ...filters.filter(content => content.visibility !== VISIBILITY.PUBLISHED)])
      return
    }
    setFilters(filters.filter(content => content.visibility !== VISIBILITY.PUBLISHED))
  }

  const initFiltersSubscription = (pageNo = page.toString(), pagesize = pageSize.toString(), searchText = search, order = sortOrder) => {
    if (disposeConfSubscriptionRef.current) {
      disposeConfSubscriptionRef.current.dispose();
    }
    const configurationSubscription = ConfigurationSubscription((result) => {

      stillLoading(false);
      setPage(result.configuration?.filters?.page || 0);
      setPageSize(result.configuration?.filters?.page_size || 0);
      setCount(result.configuration?.filters?.total_count || 0);
      handleSetFilters(result.configuration?.filters?.filters);
    },
    {
      applicationSelector : {
        pageSize : pagesize,
        page : pageNo,
        search : searchText,
        order : order
      },
      patternSelector : {
        pageSize : pagesize,
        page : pageNo,
        search : searchText,
        order : order
      },
      filterSelector : {
        pageSize : pagesize,
        page : pageNo,
        search : searchText,
        order : order
      }
    });
    disposeConfSubscriptionRef.current = configurationSubscription
  }

  const handleModalClose = () => {
    setModalOpen({
      open : false,
      filter_file : null,
      name : "",
      count : 0
    });
  }

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  async function handleSubmit({ data, name, id, type, metadata }) {
    // TODO: use filter name
    updateProgress({ showProgress : true });
    if (type === FILE_OPS.DELETE) {
      const response = await showmodal(1);
      if (response == "No") {
        updateProgress({ showProgress : false });
        return;
      }
      dataFetch(
        `/api/filter/${id}`,
        { credentials : "include", method : "DELETE" },
        () => {
          console.log("FilterFile API", `/api/filter/${id}`);
          updateProgress({ showProgress : false });
          enqueueSnackbar(`"${name}" filter deleted`, {
            variant : "success",
            action : function Action(key) {
              return (
                <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                  <CloseIcon  style={iconMedium}/>
                </IconButton>
              );
            },
            autoHideDuration : 2000,
          });
          resetSelectedRowData()();
        },
        // handleError
        handleError(ACTION_TYPES.DELETE_FILTERS)
      );
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      let body = { save : true }
      if (type === FILE_OPS.FILE_UPLOAD) {
        body = JSON.stringify({ ...body, filter_data : { filter_file : data, name : metadata.name },  config : metadata.config })
      }
      if (type === FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({ ...body, url : data, name : metadata.name, config : metadata.config })
      }
      dataFetch(
        `/api/filter`,
        { credentials : "include", method : "POST", body },
        () => {
          updateProgress({ showProgress : false });
        },
        // handleError
        handleError(ACTION_TYPES.UPLOAD_FILTERS)
      );
    }

    if (type === FILE_OPS.UPDATE) {
      dataFetch(
        `/api/filter`,
        {
          credentials : "include",
          method : "POST",
          body : JSON.stringify({ filter_data : { id, name : name }, config : data, save : true }),
        },
        () => {
          updateProgress({ showProgress : false });
        },
        handleError(ACTION_TYPES.UPLOAD_FILTERS)
      );
    }
  }

  const handleDownload = (e, id, name) => {
    e.stopPropagation();
    updateProgress({ showProgress : true });
    try {
      downloadFile({ id, name, type : "filter" })
      updateProgress({ showProgress : false });
      enqueueSnackbar(`"${name}" filter downloaded`, {
        variant : "success",
        action : function Action(key) {
          return (
            <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          );
        }
      });
    } catch (e) {
      console.error(e);
    }
  }

  function uploadHandler(ev, _, metadata) {
    if (!ev.target.files?.length) return;

    const file = ev.target.files[0];

    // Create a reader
    const reader = new FileReader();
    reader.addEventListener("load", (event) => {
      handleSubmit({
        data : event.target.result,
        name : file?.name || "meshery_" + Math.floor(trueRandom() * 100),
        type : FILE_OPS.FILE_UPLOAD,
        metadata : metadata
      });
    });
    reader.readAsText(file);
  }

  function urlUploadHandler(link, _, metadata,) {

    handleSubmit({
      data : link,
      name : "meshery_" + Math.floor(trueRandom() * 100),
      type : FILE_OPS.URL_UPLOAD,
      metadata : metadata
    });
  }

  const columns = [
    {
      name : "name",
      label : "Filter Name",
      options : {
        filter : false,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name : "created_at",
      label : "Upload Timestamp",
      options : {
        filter : false,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name : "updated_at",
      label : "Update Timestamp",
      options : {
        filter : false,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name : "visibility",
      label : "Visibility",
      options : {
        filter : false,
        sort : true,
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(_, tableMeta) {
          const visibility = filters[tableMeta.rowIndex].visibility
          return (
            <>
              <img className={classes.visibilityImg} src={`/static/img/${visibility}.svg`} />
            </>
          );
        },
      },
    },
    {
      name : "Actions",
      options : {
        filter : false,
        sort : false,
        searchable : false,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell key={index}>
              <b>{column.label}</b>
            </TableCell>
          );
        },
        customBodyRender : function CustomBody(_, tableMeta) {
          const rowData = filters[tableMeta.rowIndex];
          const visibility = filters[tableMeta.rowIndex].visibility
          return (
            <>
              {visibility === VISIBILITY.PUBLISHED ? <TooltipIcon
                placement ="top"
                title={"Clone"}
                onClick={(e) => {
                  e.stopPropagation();
                  handleClone(rowData.id, rowData.name)
                }
                }>
                <CloneIcon fill="currentColor" className={classes.iconPatt} />
              </TooltipIcon> :
                <TooltipIcon
                  title="Config"
                  onClick={(e) => {
                    e.stopPropagation()
                    setSelectedRowData(filters[tableMeta.rowIndex])
                  }}
                >
                  <EditIcon
                    aria-label="config"
                    color="inherit"
                    style={iconMedium}
                  />
                </TooltipIcon> }
              <TooltipIcon
                title="Download"
                onClick={(e) => handleDownload(e, rowData.id, rowData.name)}
              >
                <GetAppIcon data-cy="download-button" />
              </TooltipIcon>
              {canPublishFilter &&
                (visibility !== VISIBILITY.PUBLISHED) ?
                (<TooltipIcon
                  title="Publish"
                  onClick={(ev) => handlePublishModal(ev,rowData)}
                >
                  <PublicIcon fill="#F91313" data-cy="publish-button" />
                </TooltipIcon>)
                : (<TooltipIcon
                  title="Unpublish"
                  onClick={(ev) => handleUnpublishModal(ev, rowData)()}
                >
                  <PublicIcon fill="#F91313" data-cy="unpublish-button" />
                </TooltipIcon>)
              }
            </>
          );
        },
      },
    },
  ];

  columns.forEach((column, idx) => {
    if (column.name === sortOrder.split(" ")[0]) {
      columns[idx].options.sortDirection = sortOrder.split(" ")[1];
    }
  });

  async function showmodal(count) {
    let response = await modalRef.current.show({
      title : `Delete ${count ? count : ""} Filter${count > 1 ? "s" : ''}?`,

      subtitle : `Are you sure you want to delete ${count > 1 ? "these" : 'this'} ${count ? count : ""} filter${count > 1 ? "s" : ''}?`,

      options : ["Yes", "No"],
    })
    return response;
  }

  function deleteFilter(id) {
    dataFetch(
      `/api/filter/${id}`,
      {
        method : "DELETE",
        credentials : "include",
      },
      () => {
        updateProgress({ showProgress : false });

        enqueueSnackbar("Filter deleted", {
          variant : "success",
          autoHideDuration : 2000,
          action : function Action(key) {
            return (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon  style={iconMedium}/>
              </IconButton>
            );
          },
        });

      },
      handleError("Failed To Delete Filter")
    );
  }

  const options = {
    filter : false,
    sort : !(user && user.user_id === "meshery"),
    search : false,
    filterType : "textField",
    responsive : "standard",
    resizableColumns : true,
    serverSide : true,
    count,
    rowsPerPage : pageSize,
    rowsPerPageOptions : [10, 20, 25],
    fixedHeader : true,
    page,
    print : false,
    download : false,
    textLabels : {
      selectedRows : {
        text : "filter(s) selected"
      }
    },

    onCellClick : (_, meta) => meta.colIndex !== 3 && meta.colIndex !== 4 && setSelectedRowData(filters[meta.rowIndex]),

    onRowsDelete : async function handleDelete(row) {
      let response = await showmodal(Object.keys(row.lookup).length)
      console.log(response)
      if (response === "Yes") {
        const fid = Object.keys(row.lookup).map((idx) => filters[idx]?.id);
        fid.forEach((fid) => deleteFilter(fid));
      }
      // if (response === "No")
      // fetchFilters(page, pageSize, search, sortOrder);
    },

    onTableChange : (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(" : ") : [];
      let order = "";
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case "changePage":
          initFiltersSubscription(tableState.page.toString(), pageSize.toString(), search, sortOrder)
          break;
        case "changeRowsPerPage":
          initFiltersSubscription(page.toString(), tableState.rowsPerPage.toString(), search, sortOrder);
          break;
        case "search":
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              fetchFilters(page, pageSize, tableState.searchText !== null
                ? tableState.searchText
                : "", sortOrder);
              setSearch(tableState.searchText)
            }
          }, 500);
          break;
        case "sort":
          if (sortInfo.length === 2) {
            if (sortInfo[1] === "ascending") {
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
    setRowProps : (row, dataIndex, rowIndex) => {
      return {
        "data-cy" : `config-row-${rowIndex}`
      }
    },
    setTableProps : () => {
      return {
        "data-cy" : "filters-grid"
      }
    }
  };

  if (loading) {
    return <LoadingScreen animatedIcon="AnimatedFilter" message="Loading Filters..." />;
  }

  return (
    <>
      <NoSsr>
        {selectedRowData && Object.keys(selectedRowData).length > 0 && (
          <YAMLEditor filter={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} classes={classes} />
        )}
        <div className={classes.topToolbar} >
          {!selectedFilter.show && (filters.length > 0 || viewType === "table") && <div className={classes.createButton}>
            <div>
              <Button
                aria-label="Add Filter"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={handleUploadImport}
                style={{ marginRight : "2rem" }}
              >
                <PublishIcon  style={iconMedium} className={classes.addIcon} data-cy="import-button"/>
              Import Filters
              </Button>
            </div>
          </div>
          }
          <div
            className={classes.searchAndView}
            style={{
              display : 'flex',
              alignItems : 'center',
              justifyContent : 'center',
              margin : 'auto',
            }}
          >
            <SearchBar
              value={search}
              onChange={(e) => {
                setSearch(e.target.value);
                initFiltersSubscription(page.toString(), pageSize.toString(), e.target.value, sortOrder);
              }
              }
              label={"Search Filters"}
              width="80ch"
            />
          </div>
          <div style={{ justifySelf : "flex-end", marginLeft : "auto", paddingRight : "1rem", paddingTop : "0.2rem" }}>
            <CatalogFilter catalogVisibility={catalogVisibility} handleCatalogVisibility={handleCatalogVisibility} classes={classes} />
          </div>
          {!selectedFilter.show &&
            <div className={classes.viewSwitchButton}>
              <ViewSwitch data-cy="table-view" view={viewType} changeView={setViewType} />
            </div>
          }
        </div>
        {
          !selectedFilter.show && viewType === "table" &&
          <MUIDataTable
            title={<div className={classes.tableHeader}>Filters</div>}
            data={filters}
            columns={columns}
            // @ts-ignore
            options={options}
            className={classes.muiRow}
          />

        }
        {
          !selectedFilter.show && viewType==="grid" &&
            // grid view
            <FiltersGrid
              filters={filters}
              handleDeploy={handleDeploy}
              handleUndeploy={handleUndeploy}
              handleSubmit={handleSubmit}
              canPublishFilter={canPublishFilter}
              handlePublish={handlePublish}
              handleUnpublishModal={handleUnpublishModal}
              handleClone={handleClone}
              handleDownload={handleDownload}
              urlUploadHandler={urlUploadHandler}
              uploadHandler={uploadHandler}
              setSelectedFilter={setSelectedFilter}
              selectedFilter={selectedFilter}
              pages={Math.ceil(count / pageSize)}
              setPage={setPage}
              selectedPage={page}
              UploadImport={UploadImport}
              fetch={() => fetchFilters(page, pageSize, search, sortOrder)}
            />
        }
        <ConfirmationMsg
          open={modalOpen.open}
          handleClose={handleModalClose}
          submit={
            { deploy : () => handleDeploy(modalOpen.filter_file, modalOpen.name), unDeploy : () => handleUndeploy(modalOpen.filter_file, modalOpen.name) }
          }
          isDelete={!modalOpen.deploy}
          title={modalOpen.name}
          componentCount={modalOpen.count}
          tab={modalOpen.deploy ? 2 : 1}
        />
        {canPublishFilter &&
          <Modal open={publishModal.open} schema={publish_schema} uiSchema={publish_ui_schema} onChange={onChange} handleClose={handlePublishModalClose} formData={_.isEmpty(payload.catalog_data)? publishModal?.filter?.catalog_data : payload.catalog_data } aria-label="catalog publish" title={publishModal.filter?.name} handleSubmit={handlePublish} payload={payload} showInfoIcon={{ text : "Upon submitting your catalog item, an approval flow will be initiated.", link : "https://docs.meshery.io/concepts/catalog" }}/>
        }
        <PromptComponent ref={modalRef} />
        <UploadImport
          open={importModal.open}
          isFilter
          handleClose={handleUploadImportClose}
          handleUrlUpload={urlUploadHandler}
          handleUpload={uploadHandler}
          fetch={() => fetchFilters(page, pageSize, search, sortOrder)}
          configuration="Filter"
        />
      </NoSsr>
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress : bindActionCreators(updateProgress, dispatch),
  toggleCatalogContent : bindActionCreators(toggleCatalogContent, dispatch)
});

const mapStateToProps = (state) => {
  return {
    user : state.get("user")?.toObject(), selectedK8sContexts : state.get("selectedK8sContexts"),
    catalogVisibility : state.get("catalogVisibility")
  };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryFilters)));
