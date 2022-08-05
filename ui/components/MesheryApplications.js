import {
  Dialog, DialogActions, DialogContent, DialogTitle, Divider, IconButton, NoSsr,
  TableCell, Tooltip, Typography
} from "@material-ui/core";
import { makeStyles, withStyles } from "@material-ui/core/styles";
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
import FILE_OPS from "../utils/configurationFileHandlersEnum";
import { ctxUrl } from "../utils/multi-ctx";
import { getComponentsinFile, randomPatternNameGenerator as getRandomName } from "../utils/utils";
import PromptComponent from "./PromptComponent";
import UploadImport from "./UploadImport";
import UndeployIcon from "../public/static/img/UndeployIcon";
import DoneAllIcon from '@material-ui/icons/DoneAll';
import ConfirmationMsg from "./ConfirmationModal";
import ViewSwitch from "./ViewSwitch";
import ApplicationsGrid from "./MesheryApplications/ApplicationsGrid";
import { fileDownloader } from "../utils/fileDownloader";
import { trueRandom } from "../lib/trueRandom";

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
  },
  topToolbar : {
    margin : "2rem auto",
    display : "flex",
    justifyContent : "space-between",
    paddingLeft : "1rem"
  },
  viewSwitchButton : {
    justifySelf : "flex-end",
    marginLeft : "auto",
    paddingLeft : "1rem"
  },
  // text : {
  //   padding : theme.spacing(1)
  // }
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
  const [yaml, setYaml] = useState("");
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
            onClick={() => onSubmit({
              data : yaml, id : application.id, name : application.name, type : FILE_OPS.UPDATE, source_type : application.type.String
            })}
          >
            <SaveIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Application">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={() => onSubmit({
              data : yaml,
              id : application.id,
              name : application.name,
              type : FILE_OPS.DELETE
            })}
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
  FETCH_APPLICATIONS_TYPES : {
    name : "FETCH_APPLICATION_TYPES",
    error_msg : "Failed to fetch application types"
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
  DOWNLOAD_APP : {
    name : "DOWNLOAD_APP",
    error_msg : "Failed to download application file"
  }
};

function resetSelectedApplication() {
  return { show : false, application : null };
}

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
  const [selectedApplication, setSelectedApplication] = useState(resetSelectedApplication());
  const DEPLOY_URL = '/api/pattern/deploy';
  const [types, setTypes] = useState([]);
  const [modalOpen, setModalOpen] = useState({
    open : false,
    deploy : false,
    application_file : null,
    name : "",
    count : 0
  });
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ("grid")
  );

  const searchTimeout = useRef(null);

  /**
   * fetch applications when the page loads
   */
  useEffect(() => {
    fetchApplications(page, pageSize, search, sortOrder);
  }, [page, pageSize, search, sortOrder]);

  /**
   * fetch applications when the application downloads
   */
  useEffect(() => {
    getTypes();
  },[]);

  const handleModalClose = () => {
    setModalOpen({
      open : false,
      application_file : null,
      name : "",
      count : 0
    });
  }

  const handleModalOpen = (app_file, name, isDeploy) => {
    setModalOpen({
      open : true,
      deploy : isDeploy,
      application_file : app_file,
      name : name,
      count : getComponentsinFile(app_file)
    });
  }
  /**
   * fetchApplications constructs the queries based on the parameters given
   * and fetches the applications
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */
  const handleDeploy = (application_file, name) => {
    updateProgress({ showProgress : true })
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      {
        credentials : "include",
        method : "POST",
        body : application_file,
      }, () => {
        console.log("ApplicationFile Deploy API", `/api/application/deploy`);
        enqueueSnackbar(`${name} successfully deployed!!` , {
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
        updateProgress({ showProgress : false });
      },
      handleError(ACTION_TYPES.DEPLOY_APPLICATIONS)
    );
  }

  const handleUnDeploy = (application_file, name) => {
    updateProgress({ showProgress : true })
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      {
        credentials : "include",
        method : "DELETE",
        body : application_file,
      }, () => {
        console.log("ApplicationFile Undeploy API", `/api/pattern/deploy`);
        enqueueSnackbar(`${name} successfully undeployed!!`, {
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
        updateProgress({ showProgress : false });
      },
      handleError(ACTION_TYPES.UNDEPLOY_APPLICATION)
    );
  }

  // const typesMapping = {
  //   K8s:["yaml", ".yml"]
  // }
  // typesMapping[types[type]][0]
  const handleAppDownload = (id, source_type, name) => {
    updateProgress({ showProgress : true })
    dataFetch(
        `/api/application/download/${id}/${source_type}`,
        {
          credentials : "include",
          method : "GET",
        },
        () => {
          fileDownloader(id, name, source_type);
          console.log("ApplicationFile API", `/api/application/download/${id}/${source_type}`);
          updateProgress({ showProgress : false });
        },
        handleError(ACTION_TYPES.DOWNLOAD_APP)
    );
  };

  const getTypes = () => {
    dataFetch(
      `/api/application/types`,
      {
        credentials : "include",
        method : "GET",
      },
      (res) => {
        setTypes(res)
      },
      handleError(ACTION_TYPES.FETCH_APPLICATIONS_TYPES)
    );
  }

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
        console.log(result);
        console.log("ApplicationFile API", `/api/application${query}`);
        updateProgress({ showProgress : false });
        if (result) {
          setApplications(result.applications || []);
          setPage(result.page || 0);
          setPageSize(result.page_size || 0);
          setCount(result.total_count || 0);
          // setType()
        }
      },// TODO map types
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

  function handleSubmit({ data, id, name, type, source_type }) {
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
        `/api/application/${source_type}`,
        {
          credentials : "include",
          method : "PUT",
          body : JSON.stringify({ application_data : { id, name, application_file : data }, save : true }),
        },
        () => {
          console.log("ApplicationFile API", `/api/application/${source_type}`);
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
        body = JSON.stringify({
          ...body, application_data : { name : name || getRandomName(), application_file : data }
        })
      }
      if (type === FILE_OPS.URL_UPLOAD) {
        body = JSON.stringify({ ...body, url : data })
      }
      dataFetch(
        `/api/application/${source_type}`,
        {
          credentials : "include",
          method : "POST",
          body,
        },
        () => {
          console.log("ApplicationFile API", `/api/application/${source_type}`);
          updateProgress({ showProgress : false });
          fetchApplications(page, pageSize, search, sortOrder);
        },
        handleError(ACTION_TYPES.UPLOAD_APPLICATION)
      );
    }
  }

  function uploadHandler(ev, source_type) {
    if (!ev.target.files?.length) return;

    const file = ev.target.files[0];
    // Create a reader
    const reader = new FileReader();
    reader.addEventListener("load", (event) => {
      handleSubmit({
        data : event.target.result,
        name : file?.name || "meshery_" + Math.floor(trueRandom() * 100),
        type : FILE_OPS.FILE_UPLOAD,
        source_type : source_type
      });
    });
    reader.readAsText(file);
  }

  function urlUploadHandler(link, source_type) {
    handleSubmit({
      data : link, id : "", name : "meshery_" + Math.floor(trueRandom() * 100), type : FILE_OPS.URL_UPLOAD,
      source_type : source_type
    });
    console.log(link, source_type, "valid");
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
      name : "source_type",
      label : "Source Type",
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
          console.log(rowData);
          return (
            <>
              <IconButton
                title="click to download"
                onClick={() => handleAppDownload(rowData.id ,rowData.type.String, rowData.name)}
              >
                <img src={`/static/img/${(rowData.type.String).replaceAll(" ", "_").toLowerCase()}.svg`} width="45px" height="45px" />
              </IconButton>
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
          const rowData = applications[tableMeta.rowIndex];
          return (
            <>
              <IconButton
                title="Deploy"
                onClick={() => handleModalOpen(rowData.application_file, rowData.name, true)}
              >
                <DoneAllIcon data-cy="deploy-button" />
              </IconButton>
              <IconButton
                title="Undeploy"
                onClick={() => handleModalOpen(rowData.application_file, rowData.name, false)}
              >
                <UndeployIcon fill="#B32700" data-cy="undeploy-button" />
              </IconButton>
            </>
          );
        },
      },
    }
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

    onCellClick : (_, meta) => meta.colIndex !== 3 && meta.colIndex !== 4 && setSelectedRowData(applications[meta.rowIndex]),

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
        "data-cy" : "applications-grid"
      }
    }
  };

  return (
    <>

      <NoSsr>
        {selectedRowData && Object.keys(selectedRowData).length > 0 && (
          <YAMLEditor application={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
        )}
        <div className={classes.topToolbar} >
          {!selectedApplication.show && (applications.length>0 || viewType==="table") && <div className={classes.createButton}>
            <div>
              <UploadImport supportedTypes={types} isApplication = {true} aria-label="URL upload button" handleUrlUpload={urlUploadHandler} handleUpload={uploadHandler} configuration="Application"  />
            </div>

          </div>
          }
          {!selectedApplication.show &&
          <div className={classes.viewSwitchButton}>
            <ViewSwitch view={viewType} changeView={setViewType} />
          </div>
          }
        </div>
        {
          !selectedApplication.show && viewType==="table" &&
            <MUIDataTable
              title={<div className={classes.tableHeader}>Applications</div>}
              data={applications}
              columns={columns}
              // @ts-ignore
              options={options}
              className={classes.muiRow}
            />
        }
        {
          !selectedApplication.show && viewType==="grid" &&
            // grid vieww
            <ApplicationsGrid
              applications={applications}
              handleDeploy={handleDeploy}
              handleUnDeploy={handleUnDeploy}
              handleSubmit={handleSubmit}
              urlUploadHandler={urlUploadHandler}
              uploadHandler={uploadHandler}
              setSelectedApplication={setSelectedApplication}
              selectedApplication={selectedApplication}
              pages={Math.ceil(count / pageSize)}
              setPage={setPage}
              selectedPage={page}
            />
        }
        <ConfirmationMsg
          open={modalOpen.open}
          handleClose={handleModalClose}
          submit={
            { deploy : () => handleDeploy(modalOpen.application_file, modalOpen.name),  unDeploy : () => handleUnDeploy(modalOpen.application_file, modalOpen.name) }
          }
          isDelete={!modalOpen.deploy}
          title={ modalOpen.name }
          componentCount={modalOpen.count}
          tab={modalOpen.deploy ? 0 : 1}
        />
        <PromptComponent ref={modalRef} />
      </NoSsr>
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

const mapStateToProps = (state) => {
  return { user : state.get("user")?.toObject(), selectedK8sContexts : state.get("selectedK8sContexts") };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryApplications)));