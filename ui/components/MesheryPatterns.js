// @ts-check
import React, { useState, useEffect, useRef } from "react";
import { withStyles, makeStyles, MuiThemeProvider } from "@material-ui/core/styles";
import { createTheme } from '@material-ui/core/styles';
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
  Typography } from "@material-ui/core";
import { UnControlled as CodeMirror } from "react-codemirror2";
import DeleteIcon from "@material-ui/icons/Delete";
import SaveIcon from '@material-ui/icons/Save';
import UploadIcon from "@material-ui/icons/Publish";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import MUIDataTable from "mui-datatables";
import PromptComponent from "./PromptComponent";
import Moment from "react-moment";
import { withSnackbar } from "notistack";
import CloseIcon from "@material-ui/icons/Close";
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import { updateProgress } from "../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import dataFetch from "../lib/data-fetch";
import ListAltIcon from '@material-ui/icons/ListAlt';
import URLUploader from "./URLUploader";
import { MesheryConfigurationForm } from './MesheryConfigurationForm/MesheryConfigurationForm'

const styles = (theme) => ({
  grid : {
    padding : theme.spacing(2),
  },
  tableHeader : {
    fontWeight : "bolder",
    fontSize : 18,
  },
  muiRow : {
    '& .MuiTableRow-root' : {
      cursor : 'pointer'
    }
  }
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
  yamlDialogTitle : {
    display : "flex",
    alignItems : "center"
  },
  yamlDialogTitleText : {
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

function CustomToolbar(onClick, urlOnClick) {
  return function Toolbar() {
    return (
      <>
        <label htmlFor="upload-button">
          <input type="file" accept=".yaml, .yml" hidden onChange={onClick} id="upload-button" name="upload-button" />
          <Tooltip title="Upload Pattern">
            <IconButton aria-label="Upload" component="span">
              <UploadIcon />
            </IconButton>
          </Tooltip>
        </label>
        <label htmlFor="url-upload-button">
          <URLUploader onSubmit={urlOnClick} />
        </label>
      </>
    );
  };
}

function TooltipIcon({ children, onClick, title }) {
  return (
    <Tooltip title={title} placement="top" arrow interactive >
      <IconButton onClick={onClick}>
        {children}
      </IconButton>
    </Tooltip>
  )
}

function YAMLEditor({ pattern, onClose, onSubmit }) {
  const classes = useStyles();
  const [yaml, setYaml] = useState("");
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  }

  return (
    <Dialog onClose={onClose} aria-labelledby="pattern-dialog-title" open maxWidth="md" fullScreen={fullScreen} fullWidth={!fullScreen}>
      <DialogTitle disableTypography id="pattern-dialog-title" className={classes.yamlDialogTitle}>
        <Typography variant="h6" className={classes.yamlDialogTitleText}>
          {pattern.name}
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
          value={pattern.pattern_file}
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
        <Tooltip title="Update Pattern">
          <IconButton
            aria-label="Update"
            color="primary"
            onClick={() => onSubmit(yaml, pattern.id, pattern.name, "update")}
          >
            <SaveIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Pattern">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={() => onSubmit(yaml, pattern.id, pattern.name, "delete")}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </DialogActions>
    </Dialog>
  );
}

function resetSelectedPattern() {
  return { show : false, pattern : null }
}

function MesheryPatterns({
  updateProgress, enqueueSnackbar, closeSnackbar, user, classes
}) {
  const [page, setPage] = useState(0);
  const [search] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const modalRef = useRef(null);
  const [patterns, setPatterns] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [selectedPattern, setSelectedPattern] = useState(resetSelectedPattern());

  const DEPLOY_URL = '/api/pattern/deploy';

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
          color : "#607d8b",
          marginTop : "7px",
          marginRight : "8px",
        },
        clearIcon : {
          "&:hover" : {
            color : "#607d8b"
          }
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
      }
    }
  })

  const ACTION_TYPES = {
    FETCH_PATTERNS : {
      name : "FETCH_PATTERNS",
      error_msg : "Failed to fetch patterns"
    },
    UPDATE_PATTERN : {
      name : "UPDATE_PATTERN",
      error_msg : "Failed to update pattern file"
    },
    DELETE_PATTERN : {
      name : "DELETE_PATTERN",
      error_msg : "Failed to delete pattern file"
    },
    DEPLOY_PATTERN : {
      name : "DEPLOY_PATTERN",
      error_msg : "Failed to deploy pattern file"
    },
    UPLOAD_PATTERN : {
      name : "UPLOAD_PATTERN",
      error_msg : "Failed to upload pattern file"
    },
  }

  const searchTimeout = useRef(null);
  /**
   * fetch patterns when the page loads
   */
  useEffect(() => {
    fetchPatterns(page, pageSize, search, sortOrder);
  }, []);

  const handleDeploy = (pattern_file) => {
    updateProgress({ showProgress : true })
    dataFetch(
      DEPLOY_URL,
      {
        credentials : "include",
        method : "POST",
        body : pattern_file,
      }, () => {
        console.log("PatternFile Deploy API", `/api/pattern/deploy`);
        updateProgress({ showProgress : false });
        enqueueSnackbar("Pattern Successfully Deployed!", {
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
      handleError(ACTION_TYPES.DEPLOY_PATTERN),
    )
  }

  function fetchPatterns(page, pageSize, search, sortOrder) {
    if (!search) search = "";
    if (!sortOrder) sortOrder = "";

    const query = `?page=${page}&page_size=${pageSize}&search=${encodeURIComponent(search)}&order=${encodeURIComponent(
      sortOrder
    )}`;

    updateProgress({ showProgress : true });

    dataFetch(
      `/api/pattern${query}`,
      { credentials : "include", },
      (result) => {
        console.log("PatternFile API", `/api/pattern${query}`);
        updateProgress({ showProgress : false });
        if (result) {
          setPatterns(result.patterns || []);
          setPage(result.page || 0);
          setPageSize(result.page_size || 0);
          setCount(result.total_count || 0);
        }
      },
      handleError(ACTION_TYPES.FETCH_PATTERNS)
    );
  }

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
  }

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  function handleSubmit(data, id, name, type) {
    updateProgress({ showProgress : true })
    if (type === "delete") {
      dataFetch(
        `/api/pattern/${id}`,
        {
          credentials : "include",
          method : "DELETE",
        },
        () => {
          console.log("PatternFile API", `/api/pattern/${id}`);
          updateProgress({ showProgress : false });
          fetchPatterns(page, pageSize, search, sortOrder);
          resetSelectedRowData()()
        },
        handleError(ACTION_TYPES.DELETE_PATTERN)
      );
    }

    if (type === "update") {
      dataFetch(
        `/api/pattern`,
        {
          credentials : "include",
          method : "POST",
          body : JSON.stringify({ pattern_data : { id, pattern_file : data }, save : true }),
        },
        () => {
          console.log("PatternFile API", `/api/pattern`);
          updateProgress({ showProgress : false });
          fetchPatterns(page, pageSize, search, sortOrder);
        },
        handleError(ACTION_TYPES.UPDATE_PATTERN)
      );
    }

    if (type === "upload" || type === "urlupload") {
      let body
      if (type === "upload") {
        body = JSON.stringify({ pattern_data : { pattern_file : data }, save : true })
      }
      if (type === "urlupload") {
        body = JSON.stringify({ url : data, save : true })
      }
      dataFetch(
        `/api/pattern`,
        {
          credentials : "include",
          method : "POST",
          body,
        },
        () => {
          console.log("PatternFile API", `/api/pattern`);
          updateProgress({ showProgress : false });
          fetchPatterns(page, pageSize, search, sortOrder);
        },
        handleError(ACTION_TYPES.UPLOAD_PATTERN)
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

  function urlUploadHandler(link) {
    handleSubmit(link, "", "meshery_" + Math.floor(Math.random() * 100), "urlupload");
    // console.log(link, "valid");
  }
  const columns = [
    {
      name : "name",
      label : "Pattern Name",
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
          const rowData = patterns[tableMeta.rowIndex]
          return (
            <>
              <Tooltip title="Configure">
                <IconButton onClick={() => setSelectedPattern({ pattern : patterns[tableMeta.rowIndex], show : true })}>
                  <ListAltIcon />
                </IconButton>
              </Tooltip>
              <IconButton>
                <PlayArrowIcon
                  title="Deploy"
                  aria-label="deploy"
                  color="inherit"
                  onClick={() => handleDeploy(rowData.pattern_file)} //deploy endpoint to be called here
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

  async function showModal(count) {
    let response = await modalRef.current.show({
      title : `Delete ${count ? count : ""} Pattern${count > 1 ? "s" : '' }?`,

      subtitle : `Are you sure you want to delete ${count > 1 ? "these" : 'this' }  ${count ? count : ""}  pattern${count > 1 ? "s" : '' }?`,

      options : ["Yes", "No"],
    })
    return response;
  }

  function deletePattern(id) {
    dataFetch(
      `/api/pattern/${id}`,
      {
        method : "DELETE",
        credentials : "include",
      },
      () => {
        updateProgress({ showProgress : false });

        enqueueSnackbar("Pattern deleted.", {
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
        fetchPatterns(page, pageSize, search, sortOrder);
      },
      handleError("Failed to delete pattern")
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
    customToolbar : CustomToolbar(uploadHandler, urlUploadHandler),

    onCellClick : (_, meta) => meta.colIndex !== 3 && setSelectedRowData(patterns[meta.rowIndex]),

    onRowsDelete : async function handleDelete(row) {
      let response = await showModal(Object.keys(row.lookup).length)
      console.log(response)
      if (response === "Yes") {
        const fid = Object.keys(row.lookup).map(idx => patterns[idx]?.id)
        fid.forEach(fid => deletePattern(fid))
      }
      if (response === "No")
        fetchPatterns(page, pageSize, search, sortOrder);
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
          fetchPatterns(tableState.page, pageSize, search, sortOrder);
          break;
        case "changeRowsPerPage":
          fetchPatterns(page, tableState.rowsPerPage, search, sortOrder);
          break;
        case "search":
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              fetchPatterns(page, pageSize, tableState.searchText !== null
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
            fetchPatterns(page, pageSize, search, order);
          }
          break;
      }
    },
  };

  return (
    <NoSsr>
      {selectedPattern.show &&
        <MesheryConfigurationForm onSubmit={handleSubmit} show={setSelectedPattern} pattern={selectedPattern.pattern} />}

      {selectedRowData && Object.keys(selectedRowData).length > 0 && (
        <YAMLEditor pattern={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
      )}
      {
        !selectedPattern.show && <MuiThemeProvider theme={getMuiTheme()}>
          <MUIDataTable
            title={<div className={classes.tableHeader}>Patterns</div>}
            data={patterns}
            columns={columns}
            // @ts-ignore
            options={options}
            className={classes.muiRow}
          />
        </MuiThemeProvider>
      }
      <PromptComponent ref={modalRef} />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

const mapStateToProps = (state) => {
  return { user : state.get("user")?.toObject(), };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryPatterns)));
