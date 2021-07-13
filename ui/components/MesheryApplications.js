import React, { useState, useEffect, useRef } from "react";
import { withStyles } from "@material-ui/core/styles";
import {
  NoSsr,
  TableCell,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Divider,
  Tooltip
} from "@material-ui/core";
import { UnControlled as CodeMirror } from "react-codemirror2";
import DeleteIcon from "@material-ui/icons/Delete";
import SaveIcon from '@material-ui/icons/Save';
import UploadIcon from "@material-ui/icons/Publish";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import MUIDataTable from "mui-datatables";
import Moment from "react-moment";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import EditIcon from '@material-ui/icons/Edit';
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import { updateProgress } from "../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import dataFetch from "../lib/data-fetch";

const styles = (theme) => ({
  grid: {
    padding: theme.spacing(2),
  },
  tableHeader: {
    fontWeight: "bolder",
    fontSize: 18,
  },
});

function CustomToolbar(onClick) {
  return function Toolbar() {
    return (
      <>
        <label htmlFor="upload-button">
          <input type="file" accept=".yaml, .yml" hidden onChange={onClick} id="upload-button" name="upload-button" />
          <Tooltip title="Upload Application">
            <IconButton aria-label="Upload" component="span">
              <UploadIcon />
            </IconButton>
          </Tooltip>
        </label>
      </>
    );
  };
}

function YAMLEditor({ application, onClose, onSubmit }) {
  const [yaml, setYaml] = useState("");

  return (
    <Dialog onClose={onClose} aria-labelledby="application-dialog-title" open fullWidth maxWidth="md">
      <DialogTitle id="application-dialog-title">{application.name}</DialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <CodeMirror
          value={application.application_file}
          options={{
            theme: "material",
            lineNumbers: true,
            lineWrapping: true,
            gutters: ["CodeMirror-lint-markers"],
            lint: true,
            mode: "text/x-yaml",
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
            onClick={() => onSubmit(yaml, application.id, application.name, "update")}
          >
            <SaveIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Application">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={() => onSubmit(yaml, application.id, application.name, "delete")}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </DialogActions>
    </Dialog>
  );
}

function MesheryApplications({ updateProgress, enqueueSnackbar, closeSnackbar, user, classes }) {
  const [page, setPage] = useState(0);
  const [search] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [applications, setApplications] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const DEPLOY_URL = '/api/experimental/application/deploy';

  
  const ACTION_TYPES = {
    FETCH_APPLICATIONS: {
      name: "FETCH_APPLICATION" ,
      error_msg: "Failed to fetch application" 
    },
    UPDATE_APPLICATIONS: {
      name: "UPDATEAPPLICATION",
      error_msg: "Failed to update application file"
    },
    DELETE_APPLICATIONS: {
      name: "DELETEAPPLICATION",
      error_msg: "Failed to delete application file"
    },
    DEPLOY_APPLICATIONS: {
      name: "DEPLOY_APPLICATION",
      error_msg: "Failed to deploy application file"
    },
    UPLOAD_APPLICATION: {
      name: "UPLOAD_APPLICATION",
      error_msg: "Failed to upload application file"
    },
  }


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

  const handleDeploy = (application_file) => {
    dataFetch(
      DEPLOY_URL,
      {
        credentials: "include",
        method: "POST",
        body:application_file,
      },() => {
        console.log("ApplicationFile Deploy API", `/api/experimental/application/deploy`);
        // },(e) => { 
        //   console.error(e) 
        // })
        updateProgress({showProgress : false})
      },
      handleError(ACTION_TYPES.DEPLOY_APPLICATIONS)
    ) 
  } 

  function fetchApplications(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;

    updateProgress({ showProgress: true });

    dataFetch(
      `/api/experimental/application${query}`,
      {
        credentials: "include",
      },
      (result) => {
        console.log("ApplicationFile API", `/api/experimental/application${query}`);
        updateProgress({ showProgress: false });
        if (result) {
          setApplications(result.applications|| []);
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
  const handleError = (action) => (error) =>  {  
    updateProgress({ showProgress: false });

    enqueueSnackbar(`${action.error_msg}: ${error}`, {
      variant: "error",
      action: function Action(key) {
        return (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        );
      },
      autoHideDuration: 8000,
    });
  }

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  function handleSubmit(data, id, name, type) {
    updateProgress({showProgress: true})
    if (type === "delete") {
      dataFetch(
        `/api/experimental/application/${id}`,
        {
          credentials: "include",
          method: "DELETE",
        },
        () => {
          console.log("ApplicationFile API", `/api/experimental/application/${id}`);
          updateProgress({ showProgress: false });
          fetchApplications(page, pageSize, search, sortOrder);
          resetSelectedRowData()()
        },
        // handleError
        handleError(ACTION_TYPES.FETCH_APPLICATIONS)
      );
    }

    if (type === "update") {
      dataFetch(
        `/api/experimental/application`,
        {
          credentials: "include",
          method: "POST",
          body: JSON.stringify({ application_data: { id, application_file: data }, save: true }),
        },
        () => {
          console.log("ApplicationFile API", `/api/experimental/application`);
          updateProgress({ showProgress: false });
          fetchApplications(page, pageSize, search, sortOrder);
        },
        // handleError
        handleError(ACTION_TYPES.FETCH_APPLICATIONS)
      );
    }

    if (type === "upload") {
      dataFetch(
        `/api/experimental/application`,
        {
          credentials: "include",
          method: "POST",
          body: JSON.stringify({ application_data: { application_file: data }, save: true }),
        },
        () => {
          console.log("ApplicationFile API", `/api/experimental/application`);
          updateProgress({ showProgress: false });
          fetchApplications(page, pageSize, search, sortOrder);
        },
        // handleError
        handleError(ACTION_TYPES.FETCH_APPLICATIONS)
      );
    }
  }

  function uploadHandler(ev) {
    if (!ev.target.files?.length) return;

    const file = ev.target.files[0];

    // Create a reader
    const reader = new FileReader();
    reader.addEventListener("load", (event) => {
      handleSubmit(
        event.target.result, 
        "", 
        file?.name || "meshery_" + Math.floor(Math.random() * 100), 
        "upload",
      );
    });
    reader.readAsText(file);
  }

  const columns = [
    {
      name: "name",
      label: "Application Name",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
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
      name: "created_at",
      label: "Upload Timestamp",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
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
      name: "updated_at",
      label: "Update Timestamp",
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
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
      name: "Actions",
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
          const rowData = applications[tableMeta.rowIndex]
          return (
            <>
              <IconButton>
                <EditIcon
                  title="Config"  
                  aria-label="config"
                  color="inherit"
                  onClick={() => setSelectedRowData(applications[tableMeta.rowIndex])}/>
              </IconButton>
              <IconButton>               
                <PlayArrowIcon
                  title="Deploy"  
                  aria-label="deploy"
                  color="inherit"
                  onClick={() => handleDeploy(rowData.application_file)} //deploy endpoint to be called here
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

  const options = {
    filter: false,
    sort: !(user && user.user_id === "meshery"),
    search: !(user && user.user_id === "meshery"),
    filterType: "textField",
    responsive: "scrollFullHeight",
    resizableColumns: true,
    serverSide: true,
    selectableRows: "none",
    count,
    rowsPerPage: pageSize,
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    page,
    print: false,
    download: false,
    customToolbar: CustomToolbar(uploadHandler),

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(" : ") : [];
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
              fetchApplications(page, pageSize, tableState.searchText !== null ? tableState.searchText : "", sortOrder);
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
  };

  return (
    <NoSsr>
      {selectedRowData && Object.keys(selectedRowData).length > 0 && (
        <YAMLEditor application={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
      )}
      <MUIDataTable
        title={<div className={classes.tableHeader}>Meshery Application</div>}
        data={applications}
        columns={columns}
        // @ts-ignore
        options={options}
      />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  return {
    user: state.get("user").toObject(),
  };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryApplications)));
