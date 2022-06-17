import {
  Dialog, DialogActions, DialogContent, DialogTitle, Divider, IconButton, NoSsr,
  TableCell, Tooltip, Typography
} from "@material-ui/core";
import { makeStyles, MuiThemeProvider, withStyles } from "@material-ui/core/styles";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import CloseIcon from "@material-ui/icons/Close";
import DeleteIcon from "@material-ui/icons/Delete";
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import SaveIcon from '@material-ui/icons/Save';
import MUIDataTable from "mui-datatables";
import { withSnackbar } from "notistack";
import React, { useEffect, useRef, useState } from "react";
import { UnControlled as CodeMirror } from "react-codemirror2";
import Moment from "react-moment";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import dataFetch from "../lib/data-fetch";
import { updateProgress } from "../lib/store";
import { trueRandom } from "../lib/trueRandom";
import configurationTableTheme from "../themes/configurationTableTheme";
import FILE_OPS from "../utils/configurationFileHandlersEnum";
import { ctxUrl } from "../utils/multi-ctx";
import { randomPatternNameGenerator as getRandomName } from "../utils/utils";
import PromptComponent from "./PromptComponent";
import UploadImport from "./UploadImport";
import UndeployIcon from "../public/static/img/UndeployIcon";
import DoneAllIcon from '@material-ui/icons/DoneAll';

const styles = (theme) => ({
  grid : { padding : theme.spacing(2), },
  tableHeader : {
    fontWeight : "bolder",
    fontSize : 18,
  },
  muiRow : {
    '& .MuiTableRow-root' : {
      cursor : 'pointer'
    }
  },
  createButton : {
    display : "flex",
    justifyContent : "flex-start",
    alignItems : "center",
    whiteSpace : "nowrap",
    margin : "1rem auto 2rem auto"
  },
});


const useStyles = makeStyles((theme) => ({
  codeMirror : {
    '& .CodeMirror' : {
      minHeight : "300px",
      height : '60vh',
    }
  },
  backButton : {
    marginRight : theme.spacing(2),
  },
  appBar : {
    marginBottom : "16px"
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
}));


function TooltipIcon({ children, onClick, title }) {
  return (
    <Tooltip title={title} placement="top" arrow interactive >
      <IconButton onClick={onClick}>
        {children}
      </IconButton>
    </Tooltip>
  );
}

function YAMLEditor({ application, onClose, onSubmit }) {
  const classes = useStyles();
  const [yaml, setYaml] = useState(application.application_file);
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  return (
    <Dialog onClose={onClose} aria-labelledby="application-dialog-title" open maxWidth="md" fullScreen={fullScreen} fullWidth={!fullScreen}>
      <DialogTitle disableTypography id="application-dialog-title" className={classes.ymlDialogTitle}>
        <Typography variant="h6" className={classes.ymlDialogTitleText}>
          {application.name}
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
          value={application.application_file}
          className={fullScreen ? classes.fullScreenCodeMirror : ""}
          options={{
            theme : "material",
            lineNumbers : true,
            lineWrapping : true,
            gutters : ["CodeMirror-lint-markers"],
            lint : true,
            mode : "text/x-yaml",
          }}
          onChange={(_, data, val) => setYaml(val)}
        />
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <Tooltip title="Update Application">
          <IconButton
            aria-label="Update"
            color="primary"
            onClick={() => onSubmit(yaml, application.id, application.name, FILE_OPS.UPDATE)}
          >
            <SaveIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Application">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={() => onSubmit(yaml, application.id, application.name, FILE_OPS.DELETE)}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </DialogActions>
    </Dialog>
  );
}

const ACTION_TYPES = {
  FETCH_APPLICATIONS : {
    name : "FETCH_APPLICATION",
    error_msg : "Failed to fetch application"
  },
  UPDATE_APPLICATIONS : {
    name : "UPDATEAPPLICATION",
    error_msg : "Failed to update application file"
  },
  DELETE_APPLICATIONS : {
    name : "DELETEAPPLICATION",
    error_msg : "Failed to delete application file"
  },
  DEPLOY_APPLICATIONS : {
    name : "DEPLOY_APPLICATION",
    error_msg : "Failed to deploy application file"
  },
  UNDEPLOY_APPLICATION : {
    name : "UNDEPLOY_APPLICATION",
    error_msg : "Failed to undeploy application file"
  },
  UPLOAD_APPLICATION : {
    name : "UPLOAD_APPLICATION",
    error_msg : "Failed to upload application file"
  },
};

function MesheryApplications({
  updateProgress, enqueueSnackbar, closeSnackbar, user, classes, selectedK8sContexts
}) {
  const [page, setPage] = useState(0);
  const [search] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const modalRef = useRef(null);
  const [pageSize, setPageSize] = useState(10);
  const [applications, setApplications] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [close, handleClose] = useState(true);
  const DEPLOY_URL = '/api/application/deploy';

  const searchTimeout = useRef(null);

  /**
   * fetch applications when the page loads
   */
  useEffect(() => {
    fetchApplications(page, pageSize, search, sortOrder);
  }, []);

  /**
   * fetchApplications constructs the queries based on the parameters given
   * and fetches the applications
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */
  // ASSUMPTION: APPLICATION FILES ARE ONLY K8S MANIFEST
  const handleDeploy = (application_file) => {
    updateProgress({ showProgress : true })
    dataFetch(
      "/api/pattern",
      {
        credentials : "include",
        method : "POST",
        body : JSON.stringify({ k8s_manifest : application_file }),
      },
      (res) => {
        if (res) {
          const pfile = res[0].pattern_file
          dataFetch(
            ctxUrl(DEPLOY_URL, selectedK8sContexts),
            {
              credentials : "include",
              method : "POST",
              body : pfile,
            }, () => {
              console.log("ApplicationFile Deploy API", `/api/application/deploy`);
              updateProgress({ showProgress : false });
            },
            handleError(ACTION_TYPES.DEPLOY_APPLICATIONS)
          );
        } else {
          updateProgress({ showProgress : false });
          enqueueSnackbar("Failed converting kubernetes yaml to pattern file",
            { variant : "error" }
          )
        }
      },
      handleError(ACTION_TYPES.DEPLOY_APPLICATIONS)
    );
  };

  // ASSUMPTION: APPLICATION FILES ARE ONLY K8S MANIFEST
  const handleUnDeploy = (application_file) => {
    updateProgress({ showProgress : true })
    dataFetch(
      "/api/pattern",
      {
        credentials : "include",
        method : "POST",
        body : JSON.stringify({ k8s_manifest : application_file }),
      },
      (res) => {
        const pfile = res[0].pattern_file
        if (pfile) {
          dataFetch(
            ctxUrl(DEPLOY_URL, selectedK8sContexts),
            {
              credentials : "include",
              method : "DELETE",
              body : pfile,
            }, () => {
              updateProgress({ showProgress : false });
            },
            handleError(ACTION_TYPES.UNDEPLOY_APPLICATION)
          );
        } else {
          handleError(ACTION_TYPES.UNDEPLOY_APPLICATION)
        }
      },
      handleError(ACTION_TYPES.UNDEPLOY_APPLICATION)
    );
  };

  function fetchApplications(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;

    updateProgress({ showProgress : true });

    dataFetch(
      `/api/application${query}`,
      { credentials : "include", },
      (result) => {
        console.log("ApplicationFile API", `/api/application${query}`);
        updateProgress({ showProgress : false });
        if (result) {
          setApplications(result.applications || []);
          setPage(result.page || 0);
          setPageSize(result.page_size || 0);
          setCount(result.total_count || 0);
        }
      },
      // handleError
      handleError(ACTION_TYPES.FETCH_APPLICATIONS)
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

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  function handleSubmit(data, id, name, type) {
    updateProgress({ showProgress : true })
    if (type === FILE_OPS.DELETE) {
      dataFetch(
        `/api/application/${id}`,
        {
          credentials : "include",
          method : "DELETE",
        },
        () => {
          console.log("ApplicationFile API", `/api/application/${id}`);
          updateProgress({ showProgress : false });
          fetchApplications(page, pageSize, search, sortOrder);
          resetSelectedRowData()();
        },
        // handleError
        handleError(ACTION_TYPES.DELETE_APPLICATIONS)
      );
    }

    if (type === FILE_OPS.UPDATE) {
      dataFetch(
        `/api/application`,
        {
          credentials : "include",
          method : "POST",
          body : JSON.stringify({ application_data : { id, name, application_file : data }, save : true }),
        },
        () => {
          console.log("ApplicationFile API", `/api/application`);
          updateProgress({ showProgress : false });
          fetchApplications(page, pageSize, search, sortOrder);
        },
        // handleError
        handleError(ACTION_TYPES.UPDATE_APPLICATIONS)
      );
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      let body = { save : true }
      if (type === FILE_OPS.FILE_UPLOAD) {
        console.log({ data, id, name, type })
        body = JSON.stringify({
          ...body, application_data : { name : name || getRandomName(), application_file : data }
        })
      }
      if (type === FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({ ...body, url : data })
      }
      dataFetch(
        `/api/application`,
        {
          credentials : "include",
          method : "POST",
          body
        },
        () => {
          console.log("ApplicationFile API", `/api/application`);
          updateProgress({ showProgress : false });
          handleClose(true);
          fetchApplications(page, pageSize, search, sortOrder);
        },
        // handleError
        handleError(ACTION_TYPES.UPLOAD_APPLICATION)
      );
    }
  }

  function uploadHandler(ev) {
    handleClose(false);
    if (!ev.target.files?.length) return;

    const file = ev.target.files[0];

    // Create a reader
    const reader = new FileReader();
    reader.addEventListener("load", (event) => {
      handleSubmit(
        event.target.result,
        "",
        file?.name || "meshery_" + Math.floor(trueRandom() * 100),
        FILE_OPS.FILE_UPLOAD,
      );
    });
    reader.readAsText(file);
  }

  function urlUploadHandler(link) {
    handleSubmit(link, "", "meshery_" + Math.floor(trueRandom() * 100), FILE_OPS.URL_UPLOAD);
    // console.log(link, "valid");
  }

  const columns = [
    {
      name : "name",
      label : "Application Name",
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
          const rowData = applications[tableMeta.rowIndex];
          return (
            <>
              <IconButton
                title="Deploy"
                onClick={() => handleDeploy(rowData.application_file)}
              >
                <DoneAllIcon data-cy="deploy-button" />
              </IconButton>
              <IconButton
                title="Undeploy"
                onClick={() => handleUnDeploy(rowData.application_file)}
              >
                <UndeployIcon fill="rgba(0, 0, 0, 0.54)" data-cy="undeploy-button" />
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

  async function showModal(count) {
    let response = await modalRef.current.show({
      title : `Delete ${count ? count : ""} Application${count > 1 ? "s" : ''}?`,
      subtitle : `Are you sure you want to delete ${count > 1 ? "these" : 'this'} ${count ? count : ""} application${count > 1 ? "s" : ''}?`,
      options : ["Yes", "No"],
    });
    return response;
  }

  async function deleteApplication(id) {
    dataFetch(
      `/api/application/${id}`,
      {
        method : "DELETE",
        credentials : "include",
      },
      () => {
        updateProgress({ showProgress : false });

        enqueueSnackbar("Application deleted.", {
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
        fetchApplications(page, pageSize, search, sortOrder);
      },
      handleError("Failed to delete application")
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
        text : "application(s) selected"
      }
    },

    onCellClick : (_, meta) => meta.colIndex !== 3 && setSelectedRowData(applications[meta.rowIndex]),

    onRowsDelete : async function handleDelete(row) {
      let response = await showModal(Object.keys(row.lookup).length);
      console.log(response);
      if (response === "Yes") {
        const fid = Object.keys(row.lookup).map(idx => applications[idx]?.id);
        fid.forEach(fid => deleteApplication(fid));
      }
      if (response === "No")
        fetchApplications(page, pageSize, search, sortOrder);
    },

    onTableChange : (action, tableState) => {
      const sortInfo = tableState.announceText
        ? tableState.announceText.split(" : ")
        : [];
      let order = "";
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case "changePage":
          fetchApplications(tableState.page, pageSize, search, sortOrder);
          break;
        case "changeRowsPerPage":
          fetchApplications(page, tableState.rowsPerPage, search, sortOrder);
          break;
        case "search":
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              fetchApplications(page, pageSize, tableState.searchText !== null
                ? tableState.searchText
                : "", sortOrder);
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
            fetchApplications(page, pageSize, search, order);
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
    <NoSsr>
      {selectedRowData && Object.keys(selectedRowData).length > 0 && (
        <YAMLEditor application={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
      )}
      <div className={classes.createButton}>
        <div className={classes.UploadImport}>
          <UploadImport aria-label="URL upload button" handleUpload={urlUploadHandler} handleImport={uploadHandler} configuration="Application" modalStatus={close} />
        </div>
      </div>
      <MuiThemeProvider theme={configurationTableTheme()}>
        <MUIDataTable
          title={<div className={classes.tableHeader}>Applications</div>}
          data={applications}
          columns={columns}
          // @ts-ignore
          options={options}
          className={classes.muiRow}
        />
      </MuiThemeProvider>
      <PromptComponent ref={modalRef} />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

const mapStateToProps = (state) => {
  return { user : state.get("user")?.toObject(), selectedK8sContexts : state.get("selectedK8sContexts") };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryApplications)));
