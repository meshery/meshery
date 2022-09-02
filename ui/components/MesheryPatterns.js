// @ts-check
import {
  Avatar, Button, Dialog, DialogActions, DialogContent, DialogTitle, Divider, IconButton, NoSsr, TableCell, Tooltip, Typography
} from "@material-ui/core";
import { createTheme, makeStyles, MuiThemeProvider, withStyles } from "@material-ui/core/styles";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import CloseIcon from "@material-ui/icons/Close";
import DeleteIcon from "@material-ui/icons/Delete";
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import SaveIcon from '@material-ui/icons/Save';
import MUIDataTable from "mui-datatables";
import { withSnackbar } from "notistack";
import AddIcon from "@material-ui/icons/AddCircleOutline";
import React, { useEffect, useRef, useState } from "react";
import { UnControlled as CodeMirror } from "react-codemirror2";
import Moment from "react-moment";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import dataFetch from "../lib/data-fetch";
import FILE_OPS from "../utils/configurationFileHandlersEnum";
import { toggleCatalogContent, updateProgress } from "../lib/store";
import PatternForm from "../components/configuratorComponents/patternConfigurator";
import UploadImport from "./UploadImport";
import { ctxUrl } from "../utils/multi-ctx";
import { getComponentsinFile, randomPatternNameGenerator as getRandomName } from "../utils/utils";
import ViewSwitch from "./ViewSwitch";
import CatalogFilter from "./CatalogFilter";
import MesheryPatternGrid from "./MesheryPatterns/MesheryPatternGridView";
import UndeployIcon from "../public/static/img/UndeployIcon";
import DoneAllIcon from '@material-ui/icons/DoneAll';
import ConfirmationMsg from "./ConfirmationModal";
import PublishIcon from "@material-ui/icons/Publish";
import PromptComponent from "./PromptComponent";
import ConfigurationSubscription from "./graphql/subscriptions/ConfigurationSubscription";
import fetchCatalogPattern from "./graphql/queries/CatalogPatternQuery";

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
  },
  iconPatt : {
    width : "24px",
    height : "24px",
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
  createButton : {
    display : "flex",
    justifyContent : "flex-start",
    alignItems : "center",
    whiteSpace : "nowrap",
  },
  UploadImport : {
    marginLeft : "1.5rem",
  },
  noDesignAddButton : {
    marginTop : "0.5rem"
  },
  noDesignContainer : {
    padding : "2rem",
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    flexDirection : "column",
  },
  noDesignButtons : {
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    flexDirection : "row",
  },
  noDesignPaper : {
    padding : "0.5rem",
    fontSize : "3rem"
  },
  noDesignText : {
    fontSize : "2rem",
    marginBottom : "2rem",
  },
  addIcon : {
    paddingRight : ".35rem",
  },
  // text : {
  //   padding : "5px"
  // }
});

const useStyles = makeStyles((theme) => ({
  backButton : {
    marginRight : theme.spacing(2),
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
  autoComplete : {
    width : "120px",
    minWidth : "120px",
    maxWidth : 150,
    marginRight : "auto"
  },
  iconPatt : {
    width : "10px",
    height : "10px",
    "& .MuiAvatar-img" : {
      height : '60%',
      width : '60%'
    }
  }
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

function YAMLEditor({ pattern, onClose, onSubmit }) {
  const classes = useStyles();
  const [yaml, setYaml] = useState("");
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

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
            // @ts-ignore
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
            onClick={() => onSubmit({
              data : yaml, id : pattern.id, name : pattern.name, type : FILE_OPS.UPDATE
            })}
          >
            <SaveIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Pattern">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={() => onSubmit({
              data : yaml,
              id : pattern.id,
              name : pattern.name,
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

function resetSelectedPattern() {
  return { show : false, pattern : null };
}

function MesheryPatterns({
  updateProgress, enqueueSnackbar, closeSnackbar, user, classes, selectedK8sContexts, catalogVisibility, toggleCatalogContent
}) {
  const [page, setPage] = useState(0);
  const [search] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const modalRef = useRef();
  const [patterns, setPatterns] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [selectedPattern, setSelectedPattern] = useState(resetSelectedPattern());
  const [viewType, setViewType] = useState(
    /**  @type {TypeView} */
    ("grid")
  );
  const PATTERN_URL = '/api/pattern'
  const DEPLOY_URL = `${PATTERN_URL}/deploy`;
  const CLONE_URL = '/clone';
  const [modalOpen, setModalOpen] = useState({
    open : false,
    deploy : false,
    pattern_file : null,
    name : "",
    count : 0
  });

  const [importModal, setImportModal] = useState({
    open : false
  })

  const catalogContentRef = useRef();
  const catalogVisibilityRef = useRef();
  const disposeConfSubscriptionRef = useRef(null);

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
      MUIDataTableBodyCell : {
        root : {
          cursor : "pointer"
        },
      },
    }
  });

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
    UNDEPLOY_PATTERN : {
      name : "UNDEPLOY_PATTERN",
      error_msg : "Failed to undeploy pattern file"
    },
    UPLOAD_PATTERN : {
      name : "UPLOAD_PATTERN",
      error_msg : "Failed to upload pattern file"
    },
    CLONE_PATTERN : {
      name : "CLONE_PATTERN",
      error_msg : "Failed to clone pattern file"
    }
  };

  const searchTimeout = useRef(null);
  /**
   * fetch patterns when the page loads
   */
  // @ts-ignore
  useEffect(() => {
    document.body.style.overflowX = "hidden"

    return (() => document.body.style.overflowX = "auto")
  }, [page, pageSize, search, sortOrder]);

  const handleCatalogVisibility = () => {
    catalogVisibilityRef.current = !catalogVisibility
    toggleCatalogContent({ catalogVisibility : !catalogVisibility });
  }

  useEffect(() => {
    catalogVisibilityRef.current = catalogVisibility
    const fetchCatalogPatterns = fetchCatalogPattern({
      selector : {
        search : "",
        order : ""
      }
    }).subscribe({
      next : (result) => {
        catalogContentRef.current = result?.catalogPatterns;
      },
      error : (err) => console.log("There was an error fetching Catalog Pattern: ", err)
    });
    initPatternsSubscription();
    return () => {
      disposeConfSubscriptionRef.current.dispose();
      fetchCatalogPatterns.unsubscribe();
    }
  },[])

  useEffect(() => {
    handleSetPatterns(patterns)
  }, [catalogVisibility])

  const handleSetPatterns = (patterns) => {
    if (catalogVisibilityRef.current && catalogContentRef.current?.length > 0) {
      setPatterns([...catalogContentRef.current, ...patterns.filter(content => content.visibility !== "public")])
      return
    }
    setPatterns(patterns.filter(content => content.visibility !== "public"))
  }

  const initPatternsSubscription = (pageNo=page.toString(), pagesize=pageSize.toString(), searchText=search, order=sortOrder) => {
    if (disposeConfSubscriptionRef.current) {
      disposeConfSubscriptionRef.current.dispose();
    }
    const configurationSubscription = ConfigurationSubscription((result) => {
      setPage(result.configuration?.patterns.page || 0);
      setPageSize(result.configuration?.patterns.page_size || 0);
      setCount(result.configuration?.patterns.total_count || 0);
      handleSetPatterns(result.configuration?.patterns.patterns);
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
    // @ts-ignore
    setModalOpen({
      open : false,
      pattern_file : null,
      name : "",
      count : 0
    });
  }

  const handleModalOpen = (e, pattern_file, name, isDeploy) => {
    e.stopPropagation();
    setModalOpen({
      open : true,
      deploy : isDeploy,
      pattern_file : pattern_file,
      name : name,
      count : getComponentsinFile(pattern_file)
    });
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

  const handleDeploy = (pattern_file) => {
    updateProgress({ showProgress : true });
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      {
        credentials : "include",
        method : "POST",
        body : pattern_file,
      }, () => {
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
    );
  };

  const handleUnDeploy = (pattern_file) => {
    updateProgress({ showProgress : true });
    dataFetch(
      ctxUrl(DEPLOY_URL, selectedK8sContexts),
      {
        credentials : "include",
        method : "DELETE",
        body : pattern_file,
      }, () => {
        updateProgress({ showProgress : false });
        enqueueSnackbar("Pattern Successfully Undeployed!", {
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
      handleError(ACTION_TYPES.UNDEPLOY_PATTERN),
    );
  };

  function handleClone(patternID) {
    updateProgress({ showProgress : true });
    dataFetch(PATTERN_URL.concat(CLONE_URL, "/", patternID),
      {
        credentials : "include",
        method : "POST",
      },
      () => {
        updateProgress({ showProgress : false });
        enqueueSnackbar("Pattern Successfully Cloned!", {
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
      handleError(ACTION_TYPES.CLONE_PATTERN),
    );
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
          setPage(result.page || 0);
          setPageSize(result.page_size || 0);
          setCount(result.total_count || 0);
          handleSetPatterns(result.patterns || [])
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
  };

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  function handleSubmit({ data, id, name, type }) {
    updateProgress({ showProgress : true })
    if (type === FILE_OPS.DELETE) {
      dataFetch(
        `/api/pattern/${id}`,
        {
          credentials : "include",
          method : "DELETE",
        },
        () => {
          console.log("PatternFile API", `/api/pattern/${id}`);
          updateProgress({ showProgress : false });
          resetSelectedRowData()();
        },
        handleError(ACTION_TYPES.DELETE_PATTERN)
      );
    }

    if (type === FILE_OPS.UPDATE) {
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
        },
        handleError(ACTION_TYPES.UPDATE_PATTERN)
      );
    }

    if (type === FILE_OPS.FILE_UPLOAD || type === FILE_OPS.URL_UPLOAD) {
      let body;
      if (type === FILE_OPS.FILE_UPLOAD) {
        body = JSON.stringify({
          pattern_data : {
            name,
            pattern_file : data,
          },
          save : true
        })
      }
      if (type === FILE_OPS.URL_UPLOAD) {
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
      // @ts-ignore
      handleSubmit({
        data : event.target.result,
        name : file?.name || getRandomName(),
        type : FILE_OPS.FILE_UPLOAD
      });
    });
    reader.readAsText(file);
  }

  function urlUploadHandler(link) {
    handleSubmit({
      data : link,
      id : "",
      name : getRandomName(),
      type : FILE_OPS.URL_UPLOAD
    });
  }

  const columns = [
    {
      name : "name",
      label : "Name",
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
          const visibility = patterns[tableMeta.rowIndex].visibility
          return (
            <>
              <img src={`/static/img/${visibility}.svg`} />
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
          const rowData = patterns[tableMeta.rowIndex];
          const visibility = patterns[tableMeta.rowIndex].visibility
          return (
            <>
              { visibility === "public" ? <IconButton onClick={(e) => {
                e.stopPropagation();
                handleClone(rowData.id)
              }
              }>
                <img src="/static/img/clone.svg" />
              </IconButton> :

                <IconButton onClick={(e) => {
                  e.stopPropagation();
                  setSelectedPattern({ pattern : patterns[tableMeta.rowIndex], show : true })
                }
                }
                >
                  <Avatar src="/static/img/pattwhite.svg" className={classes.iconPatt} imgProps={{ height : "16px", width : "16px" }} />
                </IconButton> }
              {/*</Tooltip> */}
              <IconButton
                title="Deploy"
                onClick={(e) => handleModalOpen(e, rowData.pattern_file, rowData.name, true)}
              >
                <DoneAllIcon data-cy="deploy-button" />
              </IconButton>
              <IconButton
                title="Undeploy"
                onClick={(e) => handleModalOpen(e, rowData.pattern_file, rowData.name, false)}
              >
                <UndeployIcon fill="#8F1F00" data-cy="undeploy-button" />
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

  async function showModal(count, patterns) {
    console.log("patterns to be deleted", count, patterns);
    let response = await modalRef.current.show({
      title : `Delete ${count ? count : ""} Design${count > 1 ? "s" : ''}?`,

      subtitle : `Are you sure you want to delete the ${patterns} design${count > 1 ? "s" : ''}?`,

      options : ["Yes", "No"],
    });
    return response;
  }

  async function deletePatterns(patterns) {
    const jsonPatterns = JSON.stringify(patterns)


    updateProgress({ showProgress : true })
    dataFetch("/api/patterns/delete", {
      method : "POST",
      credentials : "include",
      body : jsonPatterns
    },
    () => {
      console.log("PatternFile Delete Multiple API", `/api/pattern/delete`);
      updateProgress({ showProgress : false });
      setTimeout(() => {
        enqueueSnackbar(`${patterns.patterns.length} Designs Deleted`,
          {
            variant : "success",
            autoHideDuration : 2000,
            action : function Action(key) {
              return (
                <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
                  <CloseIcon />
                </IconButton>
              );
            }
          }
        )
        resetSelectedRowData()()
      }, 1200);
    },
    handleError(ACTION_TYPES.DELETE_PATTERN)
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
        text : "pattern(s) selected"
      }
    },

    onCellClick : (_, meta) => meta.colIndex !== 3 && setSelectedRowData(patterns[meta.rowIndex]),

    onRowsDelete : async function handleDelete(row) {
      const toBeDeleted = Object.keys(row.lookup).map(idx => (
        {
          id : patterns[idx]?.id,
          name : patterns[idx]?.name,
        }
      ))
      let response = await showModal(toBeDeleted.length, toBeDeleted.map(p => " "+p.name))
      if (response.toLowerCase() === "yes") {
        deletePatterns({ patterns : toBeDeleted })
      }
      // if (response.toLowerCase() === "no")
      // fetchPatterns(page, pageSize, search, sortOrder);
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
          initPatternsSubscription(tableState.page.toString(), pageSize.toString(), search, sortOrder);
          break;
        case "changeRowsPerPage":
          initPatternsSubscription(page.toString(), tableState.rowsPerPage.toString(), search, sortOrder);
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
            initPatternsSubscription(page.toString(), pageSize.toString(), search, order);
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
          <YAMLEditor pattern={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
        )}
        {selectedPattern.show &&
          <PatternForm onSubmit={handleSubmit} show={setSelectedPattern} pattern={selectedPattern.pattern} />
        }
        <div className={classes.topToolbar} >
          {!selectedPattern.show && (patterns.length>0 || viewType==="table") && <div className={classes.createButton}>
            <div>
              <Button
                aria-label="Add Pattern"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={() => setSelectedPattern({
                  pattern : { id : null, name : "New Pattern", pattern_file : "name: New Pattern\nservices:" },
                  show : true,
                })}
                style={{ marginRight : "2rem" }}
              >
                <AddIcon className={classes.addIcon} />
              Create Design
              </Button>
              <Button
                aria-label="Add Pattern"
                variant="contained"
                color="primary"
                size="large"
                // @ts-ignore
                onClick={handleUploadImport}
                style={{ marginRight : "2rem" }}
              >
                <PublishIcon className={classes.addIcon} />
              Import Design
              </Button>
            </div>
          </div>
          }

          <div style={{ justifySelf : "flex-end", marginLeft : "auto", paddingRight : "1rem", paddingTop : "0.2rem" }}>
            <CatalogFilter catalogVisibility={catalogVisibility} handleCatalogVisibility={handleCatalogVisibility} />
          </div>

          {!selectedPattern.show &&
          <div className={classes.viewSwitchButton}>
            <ViewSwitch view={viewType} changeView={setViewType} />
          </div>
          }
        </div>
        {
          !selectedPattern.show && viewType==="table" && <MuiThemeProvider theme={getMuiTheme() }>
            <MUIDataTable
              title={<div className={classes.tableHeader}>Designs</div>}
              data={patterns}
              columns={columns}
              // @ts-ignore
              options={options}
              className={classes.muiRow}
            />
          </MuiThemeProvider>
        }
        {
          !selectedPattern.show && viewType==="grid" &&
            // grid vieww
            <MesheryPatternGrid
              patterns={patterns}
              handleDeploy={handleDeploy}
              handleUnDeploy={handleUnDeploy}
              handleClone={handleClone}
              urlUploadHandler={urlUploadHandler}
              uploadHandler={uploadHandler}
              supportedTypes="null"
              handleSubmit={handleSubmit}
              setSelectedPattern={setSelectedPattern}
              selectedPattern={selectedPattern}
              pages={Math.ceil(count / pageSize)}
              setPage={setPage}
              selectedPage={page}
              UploadImport={UploadImport}
            />
        }
        <ConfirmationMsg
          open={modalOpen.open}
          handleClose={handleModalClose}
          submit={
            { deploy : () => handleDeploy(modalOpen.pattern_file),  unDeploy : () => handleUnDeploy(modalOpen.pattern_file) }
          }
          isDelete={!modalOpen.deploy}
          title={modalOpen.name}
          componentCount={modalOpen.count}
          tab={modalOpen.deploy ? 0 : 1}
        />
        <UploadImport open={importModal.open} handleClose={handleUploadImportClose} aria-label="URL upload button" handleUrlUpload={urlUploadHandler} handleUpload={uploadHandler} configuration="Design" />
        <PromptComponent ref={modalRef} />
      </NoSsr>
    </>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), toggleCatalogContent : bindActionCreators(toggleCatalogContent, dispatch) });

const mapStateToProps = (state) => {
  return { user : state.get("user")?.toObject(), selectedK8sContexts : state.get("selectedK8sContexts"),
    catalogVisibility : state.get("catalogVisibility")
  };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryPatterns)));
