import React, { useState, useEffect, useRef } from "react";
import { withStyles, makeStyles, MuiThemeProvider } from "@material-ui/core/styles";
import {  createTheme } from '@material-ui/core/styles';
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
  Typography
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
import DoneAllIcon from '@material-ui/icons/DoneAll';
import { updateProgress } from "../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import dataFetch from "../lib/data-fetch";
import UploadImport from "./UploadImport";
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import FILE_OPS from "../utils/configurationFileHandlersEnum";
import ViewSwitch from "./ViewSwitch";
import FiltersGrid from "./MesheryFilters/FiltersGrid";
import { trueRandom } from "../lib/trueRandom";
import { ctxUrl } from "../utils/multi-ctx";
import ConfirmationMsg from "./ConfirmationModal";

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
    margin : "1rem 0 2rem 1rem"
  },
  topToolbar : {
    display : "flex"
  },
  viewSwitchButton : {
    justifySelf : "flex-end",
    marginLeft : "auto",
    paddingLeft : "1rem"
  },
});

const useStyles = makeStyles(() => ({
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

}))

function TooltipIcon({ children, onClick, title }) {
  return (
    <Tooltip title={title} placement="top" arrow interactive >
      <IconButton onClick={onClick}>
        {children}
      </IconButton>
    </Tooltip>
  )
}

function YAMLEditor({ filter, onClose, onSubmit }) {
  const classes = useStyles();
  const [yaml, setYaml] = useState("");
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  }

  return (
    <Dialog onClose={onClose} aria-labelledby="filter-dialog-title" open maxWidth="md" fullScreen={fullScreen} fullWidth={!fullScreen}>
      <DialogTitle disableTypography id="filter-dialog-title" className={classes.ymlDialogTitle}>
        <Typography variant="h6" className={classes.ymlDialogTitleText}>
          {filter.name}
        </Typography>
        <TooltipIcon
          title={fullScreen ? "Exit Fullscreen" : "Enter Fullscreen"}
          onClick={toggleFullScreen}>
          {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
        </TooltipIcon>
        <TooltipIcon title="Exit" onClick={onClose}>
          <CloseIcon />
        </TooltipIcon>
      </DialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <CodeMirror
          value={filter.filter_file}
          className={fullScreen ? classes.fullScreenCodeMirror : ""}
          options={{
            theme : "material",
            lineNumbers : true,
            lineWrapping : true,
            gutters : ["CodeMirror-lint-markers"],
            lint : true,
            mode : "text/x-yaml",
          }}
          onChange={(val) => setYaml(val)}
        />
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <Tooltip title="Delete Filter">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={() => onSubmit(yaml, filter.id, filter.name, FILE_OPS.DELETE)}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </DialogActions>
    </Dialog>
  );
}

function resetSelectedFilter() {
  return { show : false, filter : null };
}

function MesheryFilters({ updateProgress, enqueueSnackbar, closeSnackbar, user, classes, selectedK8sContexts }) {
  const [page, setPage] = useState(0);
  const [search] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const modalRef = useRef(null);
  const [pageSize, setPageSize] = useState(10);
  const [filters, setFilters] = useState([]);
  const [selectedFilter, setSelectedFilter] = useState(resetSelectedFilter());
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ("grid")
  );
  const DEPLOY_URL = "/api/filter/deploy";
  const [modalOpen, setModalOpen] = useState({
    open : false,
    filter_file : null,
    deploy : false
  });

  const getMuiTheme = () => createTheme({
    overrides : {
      MuiInput : {
        underline : {
          "&:hover:not(.Mui-disabled):before" : {
            borderBottom : "2px solid #222"
          },
          "&:after" : {
            borderBottom : "2px solid #222"
          }
        }
      },
      MUIDataTableSearch : {
        searchIcon : {
          color : "#607d8b" ,
          marginTop : "7px",
          marginRight : "8px",
        },
        clearIcon : {
          "&:hover" : {
            color : "#607d8b"
          }
        },
      },
      MUIDataTableSelectCell : {
        checkboxRoot : {
          '&$checked' : {
            color : '#607d8b',
          },
        },
      },
      MUIDataTableToolbar : {
        iconActive : {
          color : "#222"
        },
        icon : {
          "&:hover" : {
            color : "#607d8b"
          }
        },
      },
    }
  })

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
    UPLOAD_FILTERS : {
      name : "UPLOAD_FILTERS",
      error_msg : "Failed to upload filter file",
    },
  };

  const searchTimeout = useRef(null);

  /**
   * fetch filters when the page loads
   */
  useEffect(() => {
    fetchFilters(page, pageSize, search, sortOrder);
  }, [page, pageSize, search, sortOrder]);

  /**
   * fetchFilters constructs the queries based on the parameters given
   * and fetches the filters
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */

  const handleDeploy = (filter_file) => {
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      { credentials : "include", method : "POST", body : filter_file },
      () => {
        console.log("FilterFile Deploy API", `/api/filter/deploy`);
        updateProgress({ showProgress : false });
      },
      handleError(ACTION_TYPES.DEPLOY_FILTERS)
    );
  };

  function fetchFilters(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
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
          setFilters(result.filters || []);
          setPage(result.page || 0);
          setPageSize(result.page_size || 0);
          setCount(result.total_count || 0);
        }
      },
      // handleError
      handleError(ACTION_TYPES.FETCH_FILTERS)
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
            <CloseIcon />
          </IconButton>
        );
      },
      autoHideDuration : 8000,
    });
  };

  const handleModalOpen = (app_file, isDeploy) => {
    setModalOpen({
      open : true,
      filter_file : app_file,
      deploy : isDeploy
    });
  }

  const handleModalClose = () => {
    setModalOpen({
      open : false,
      filter_file : null,
      // deploy: false
    });
  }

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  function handleSubmit(data, id, name, type) {
    updateProgress({ showProgress : true });
    if (type === FILE_OPS.DELETE) {
      dataFetch(
        `/api/filter/${id}`,
        { credentials : "include", method : "DELETE" },
        () => {
          console.log("FilterFile API", `/api/filter/${id}`);
          updateProgress({ showProgress : false });
          fetchFilters(page, pageSize, search, sortOrder);
          resetSelectedRowData()();
        },
        // handleError
        handleError(ACTION_TYPES.DELETE_FILTERS)
      );
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      let body = { save : true }
      if (type ===FILE_OPS.FILE_UPLOAD) {
        body = JSON.stringify({ ...body, filter_data : { filter_file : data } })
      }
      if (type ===  FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({ ...body, url : data })
      }
      dataFetch(
        `/api/filter`,
        { credentials : "include", method : "POST", body },
        () => {
          console.log("FilterFile API", `/api/filter`);
          updateProgress({ showProgress : false });
          fetchFilters(page, pageSize, search, sortOrder);
        },
        // handleError
        handleError(ACTION_TYPES.UPLOAD_FILTERS)
      );
    }
  }

  function uploadHandler(ev) {
    if (!ev.target.files?.length) return;

    const file = ev.target.files[0];

    // Create a reader
    const reader = new FileReader();
    reader.addEventListener("load", (event) => {
      handleSubmit(event.target.result, "", file?.name || "meshery_" + Math.floor(trueRandom() * 100), FILE_OPS.FILE_UPLOAD);
    });
    reader.readAsText(file);
  }

  function urlUploadHandler(link) {
    handleSubmit(link, "", "meshery_" + Math.floor(trueRandom() * 100),  FILE_OPS.URL_UPLOAD);
    console.log(link, "valid");
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
          return (
            <>
              <IconButton>
                <EditIcon
                  title="Config"
                  aria-label="config"
                  color="inherit"
                  onClick={() => setSelectedRowData(filters[tableMeta.rowIndex])}
                />
              </IconButton>
              <IconButton>
                <DoneAllIcon
                  title="Deploy"
                  aria-label="deploy"
                  color="inherit"
                  onClick={() => handleModalOpen(rowData.filter_file, true)} //deploy endpoint to be called here
                  data-cy="deploy-button"
                />
              </IconButton>
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
      title : `Delete ${count ? count : ""} Filter${count > 1 ? "s" : '' }?`,

      subtitle : `Are you sure you want to delete ${count > 1 ? "these" : 'this' } ${count ? count : ""} filter${count > 1 ? "s" : '' }?`,

      options : ["Yes", "No"], })
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

        enqueueSnackbar("Filter Successfully Deleted!", {
          variant : "success",
          autoHideDuration : 2000,
          action : function Action(key) {
            return (
              <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                <CloseIcon />
              </IconButton>
            );
          },
        });

        fetchFilters(page, pageSize, search, sortOrder);
      },
      handleError("Failed To Delete Filter")
    );
  }

  const options = {
    filter : false,
    sort : !(user && user.user_id === "meshery"),
    search : !(user && user.user_id === "meshery"),
    filterType : "textField",
    responsive : "scrollFullHeight",
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

    onRowsDelete : async function handleDelete(row) {
      let response  = await showmodal(Object.keys(row.lookup).length)
      console.log(response)
      if (response === "Yes") {
        const fid = Object.keys(row.lookup).map((idx) => filters[idx]?.id);
        fid.forEach((fid) => deleteFilter(fid));
      }
      if (response === "No")
        fetchFilters(page, pageSize, search, sortOrder);
    },

    onTableChange : (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(" : ") : [];
      let order = "";
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case "changePage":
          fetchFilters(tableState.page, pageSize, search, sortOrder);
          break;
        case "changeRowsPerPage":
          fetchFilters(page, tableState.rowsPerPage, search, sortOrder);
          break;
        case "search":
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              fetchFilters(page, pageSize, tableState.searchText !== null ? tableState.searchText : "", sortOrder);
            }
          }, 500);
          break;
        case "sort":
          if (sortInfo.length == 2) {
            if (sortInfo[1] === "ascending") {
              order = `${columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) {
            fetchFilters(page, pageSize, search, order);
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

  return (
    <>

      <NoSsr>
        {selectedRowData && Object.keys(selectedRowData).length > 0 && (
          <YAMLEditor filter={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
        )}
        <div className={classes.topToolbar} >
          {!selectedFilter.show && (filters.length>0 || viewType==="table") && <div className={classes.createButton}>
            <div>
              <UploadImport aria-label="URL upload button" handleUpload={urlUploadHandler} handleImport={uploadHandler} configuration="Filter" />
            </div>
          </div>
          }
          {!selectedFilter.show &&
          <div className={classes.viewSwitchButton}>
            <ViewSwitch view={viewType} changeView={setViewType} />
          </div>
          }
        </div>
        {
          !selectedFilter.show && viewType==="table" && <MuiThemeProvider theme={getMuiTheme() }>
            <MUIDataTable
              title={<div className={classes.tableHeader}>Filters</div>}
              data={filters}
              columns={columns}
              // @ts-ignore
              options={options}
              className={classes.muiRow}
            />
          </MuiThemeProvider>
        }
        {
          !selectedFilter.show && viewType==="grid" &&
            // grid vieww
            <FiltersGrid
              filters={filters}
              handleDeploy={handleDeploy}
              handleSubmit={handleSubmit}
              setSelectedFilter={setSelectedFilter}
              selectedFilter={selectedFilter}
              pages={Math.ceil(count / pageSize)}
              setPage={setPage}
              selectedPage={page}
            />
        }
        <ConfirmationMsg
          open={modalOpen.open}
          handleClose={handleModalClose}
          submit={() => handleDeploy(modalOpen.filter_file)}
          isDelete={!modalOpen.deploy}
          title={<Typography variant="h6" className={classes.text} >The selected operation will be applied to following contexts.</Typography>}
        />
      </NoSsr>
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch) });

const mapStateToProps = (state) => {
  return { user : state.get("user")?.toObject(), selectedK8sContexts : state.get("selectedK8sContexts") };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryFilters)));
