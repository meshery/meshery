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
  Grid,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  CardContent,
  Card,
  CardActions,
  AppBar,
  Toolbar
} from "@material-ui/core";
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
import { updateProgress } from "../lib/store";
import TableSortLabel from "@material-ui/core/TableSortLabel";
import URLUploader from "./URLUploader";
import dataFetch, { promisifiedDataFetch } from "../lib/data-fetch";
import { CircularProgress } from "@material-ui/core";
import PatternServiceForm from "./MesheryMeshInterface/PatternServiceForm";
import FullscreenIcon from '@material-ui/icons/Fullscreen';
import FullscreenExitIcon from '@material-ui/icons/FullscreenExit';
import ExpandMoreIcon from "@material-ui/icons/ExpandMore"
import { Button } from "@material-ui/core";
import jsYaml from "js-yaml";
import PascalCaseToKebab from "../utils/PascalCaseToKebab";
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import AppsIcon from "./ConnectionWizard/icons/apps";

const styles = (theme) => ({ grid : { padding : theme.spacing(2), },
  tableHeader : { fontWeight : "bolder",
    fontSize : 18, },
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


function CustomToolbar(onClick, urlOnClick) {
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

function YAMLEditor({ application, onClose, onSubmit }) {
  const classes = useStyles();
  const [yaml, setYaml] = useState("");
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  }

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

function MesheryApplications({
  updateProgress, enqueueSnackbar, closeSnackbar, user, classes
}) {
  const [page, setPage] = useState(0);
  const [search] = useState("");
  const [sortOrder] = useState("");
  const [count, setCount] = useState(0);
  const modalRef = useRef(null);
  const [pageSize, setPageSize] = useState(10);
  const [showForm, setShowForm] = useState(false);
  const [applications, setApplications] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const DEPLOY_URL = '/api/application/deploy';

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
    FETCH_APPLICATIONS : { name : "FETCH_APPLICATION" ,
      error_msg : "Failed to fetch application" },
    UPDATE_APPLICATIONS : { name : "UPDATEAPPLICATION",
      error_msg : "Failed to update application file" },
    DELETE_APPLICATIONS : { name : "DELETEAPPLICATION",
      error_msg : "Failed to delete application file" },
    DEPLOY_APPLICATIONS : { name : "DEPLOY_APPLICATION",
      error_msg : "Failed to deploy application file" },
    UPLOAD_APPLICATION : { name : "UPLOAD_APPLICATION",
      error_msg : "Failed to upload application file" },
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
      { credentials : "include",
        method : "POST",
        body : application_file, },() => {
        console.log("ApplicationFile Deploy API", `/api/application/deploy`);
        // },(e) => {
        //   console.error(e)
        // })
        updateProgress({ showProgress : false })
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

    updateProgress({ showProgress : true });

    dataFetch(
      `/api/application${query}`,
      { credentials : "include", },
      (result) => {
        console.log("ApplicationFile API", `/api/application${query}`);
        updateProgress({ showProgress : false });
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
    updateProgress({ showProgress : false });

    enqueueSnackbar(`${action.error_msg}: ${error}`, { variant : "error",
      action : function Action(key) {
        return (
          <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
            <CloseIcon />
          </IconButton>
        );
      },
      autoHideDuration : 8000, });
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
        `/api/application/${id}`,
        { credentials : "include",
          method : "DELETE", },
        () => {
          console.log("ApplicationFile API", `/api/application/${id}`);
          updateProgress({ showProgress : false });
          fetchApplications(page, pageSize, search, sortOrder);
          resetSelectedRowData()()
        },
        // handleError
        handleError(ACTION_TYPES.DELETE_APPLICATIONS)
      );
    }

    if (type === "update") {
      dataFetch(
        `/api/application`,
        { credentials : "include",
          method : "POST",
          body : JSON.stringify({ application_data : { id, application_file : data }, save : true }), },
        () => {
          console.log("ApplicationFile API", `/api/application`);
          updateProgress({ showProgress : false });
          fetchApplications(page, pageSize, search, sortOrder);
        },
        // handleError
        handleError(ACTION_TYPES.UPDATE_APPLICATIONS)
      );
    }

    if (type === "upload" || type === "urlupload") {
      let body = { save : true }
      if (type === "upload") {
        body = JSON.stringify({ ...body,   application_data : { application_file : data }
        })
      }
      if (type === "urlupload") {
        body = JSON.stringify({ ...body, url : data })
      }
      dataFetch(
        `/api/application`,
        { credentials : "include",
          method : "POST",
          body },
        () => {
          console.log("ApplicationFile API", `/api/application`);
          updateProgress({ showProgress : false });
          fetchApplications(page, pageSize, search, sortOrder);
        },
        // handleError
        handleError(ACTION_TYPES.UPLOAD_APPLICATION)
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
    { name : "name",
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
      }, },
    { name : "created_at",
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
      }, },
    { name : "updated_at",
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
      }, },
    { name : "Actions",
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
          const rowData = applications[tableMeta.rowIndex]
          return (
            <>
              <Tooltip
                title = "configure">
                <IconButton onClick={() => setShowForm({ application : applications[tableMeta.rowIndex], show : true })}>
                  <AppsIcon />
                </IconButton>
              </Tooltip>
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
      }, },
  ];

  columns.forEach((column, idx) => {
    if (column.name === sortOrder.split(" ")[0]) {
      columns[idx].options.sortDirection = sortOrder.split(" ")[1];
    }
  });

  async function showModal() {
    let response = await modalRef.current.show({ title : "Delete Aplication?",

      subtitle : "Are you sure you want to delete this application?",

      options : ["Yes", "No"], })
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
    customToolbar : CustomToolbar(uploadHandler,urlUploadHandler),

    onCellClick : (_, meta) => meta.colIndex !== 3 && setSelectedRowData(applications[meta.rowIndex]),

    onRowsDelete : async function handleDelete(row) {
      let response = await showModal()
      console.log(response)
      if (response === "Yes") {
        const fid = Object.keys(row.lookup).map(idx => applications[idx]?.id)
        fid.forEach(fid => deleteApplication(fid))
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
  };

  return (
    <NoSsr>
      {showForm &&
        <PatternForm onSubmit={handleSubmit} show={setShowForm} application={showForm.application} />}
      {selectedRowData && Object.keys(selectedRowData).length > 0 && (
        <YAMLEditor application={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
      )}
      {
        !showForm && <MuiThemeProvider theme={getMuiTheme()}><MUIDataTable
          title={<div className={classes.tableHeader}>Applications</div>}
          data={applications}
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
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryApplications)));

// --------------------------------------------------------------------------------------------------
// -------------------------------------------- Aplications Configuration---------------------------------------
// --------------------------------------------------------------------------------------------------

function PatternForm({ application, onSubmit, show }) {
  const [schemaSet, setSchemaSet] = useState();
  const [deployServiceConfig, setDeployServiceConfig] = useState(getPatternJson() || {});
  const [yaml, setYaml] = useState("");
  const [expanded, setExpanded] = useState([]);
  // const [changedYaml, setChangedYaml] = useState("");
  const classes = useStyles();

  function getPatternJson() {
    const patternString = application.application_file;
    return jsYaml.load(patternString).services;
  }

  async function fetchWorkloadAndTraitsSchema() {
    try {
      const workloads = await promisifiedDataFetch("/api/oam/workload");
      const traits = await promisifiedDataFetch("/api/oam/trait");

      console.log({ workloads, traits });

      const workloadTraitSets = createWorkloadTraitSets(workloads, traits);

      return workloadTraitSets;
    } catch (e) {
      console.log("Error in Fetching Workload or traits", e);
      return {};
    }
  }

  function createWorkloadTraitSets(workloads, traits) {
    const sets = [];
    workloads?.forEach((w) => {
      const item = { workload : w, traits : [] };

      item.traits = traits?.filter((t) => {
        if (Array.isArray(t?.oam_definition?.spec?.appliesToWorkloads))
          return t?.oam_definition?.spec?.appliesToWorkloads?.includes(w?.oam_definition?.metadata?.name);

        return false;
      });

      sets.push(item);
    });

    return sets;
  }

  async function getJSONSchemaSets() {
    const wtSets = await fetchWorkloadAndTraitsSchema();

    return wtSets?.map((s) => {
      const item = {
        workload : JSON.parse(s.workload?.oam_ref_schema),
        traits : s.traits?.map((t) => {
          const trait = JSON.parse(t?.oam_ref_schema);

          // Attaching internal metadata to the json schema
          trait._internal = {
            patternAttributeName : t?.oam_definition.metadata.name,
          };

          return trait;
        }),
        type : s.workload?.metadata?.["ui.meshery.io/category"],
      };

      // Attaching internal metadata to the json schema
      item.workload._internal = {
        patternAttributeName : s.workload?.oam_definition.metadata.name,
      };

      return item;
    });
  }

  function getPatternAttributeName(jsonSchema) {
    return PascalCaseToKebab(jsonSchema?._internal?.patternAttributeName || "NA");
  }

  function getPatternKey(cfg) {
    return Object.keys(cfg?.services)?.[0] || undefined;
  }

  const handleSubmit = (cfg, patternName) => {
    console.log("submitted", { cfg, patternName })
    const key = getPatternKey(cfg);
    handleDeploy({ ...deployServiceConfig, [getPatternKey(cfg)] : cfg?.services?.[key] });
    if (key)
      setDeployServiceConfig({ ...deployServiceConfig, [getPatternKey(cfg)] : cfg?.services?.[key] });
    handleExpansion(patternName)
  }

  const handleChangeData = (cfg, patternName) => {
    console.log("Ran Changed", { cfg, patternName })
    const key = getPatternKey(cfg);
    handleDeploy({ ...deployServiceConfig, [getPatternKey(cfg)] : cfg?.services?.[key] });
    if (key)
      setDeployServiceConfig({ ...deployServiceConfig, [getPatternKey(cfg)] : cfg?.services?.[key] });
  }

  const handleDelete = (cfg, patternName) => {
    console.log("deleted", cfg);
    const newCfg = schemaSet.filter(schema => schema.workload.title !== patternName)
    setSchemaSet(newCfg);
  }

  const handleDeploy = (cfg) => {
    const deployConfig = {};
    deployConfig.name = application.name;
    deployConfig.services = cfg;
    const deployConfigYaml = jsYaml.dump(deployConfig);
    setYaml(deployConfigYaml);
  }

  const handleExpansion = (item) => {
    let expandedItems = [...expanded];
    if (expandedItems.includes(item)) {
      expandedItems = expandedItems.filter(el => el !== item);
    } else {
      expandedItems.push(item);
    }
    setExpanded(expandedItems);
  }

  function handleSubmitFinalPattern(yaml, id, name, action) {
    onSubmit(yaml, id, name, action);
    show(false);
  }

  const ns = "default";

  function saveCodeEditorChanges(data) {
    setYaml(data.valueOf().getValue())
  }

  useEffect(() => {
    getJSONSchemaSets().then((res) => setSchemaSet(res));
  }, []);

  if (!schemaSet) {
    return <CircularProgress />
  }

  return (
    <>
      <AppBar position="static" className={classes.appBar} elevation={0}>
        <Toolbar>
          <IconButton edge="start" className={classes.backButton} color="inherit" onClick={() => show(false)}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h6" className={classes.title}>
            Edit Application Configuration of <i>{`${application.name}`}</i>
          </Typography>
        </Toolbar>
      </AppBar>
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          {schemaSet
            .filter((s) => s.type !== "addon")
            .sort((a, b) => (a.workload?.title < b.workload?.title ? -1 : 1))
            .map((s) => (
              accordion(s)
            ))}
          <Accordion
            expanded={expanded.includes('addon')}
            onChange={() => handleExpansion('addon')}
            style={{ width : '100%' }}
          >
            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
              <Typography variant="h6">
                Configure Addons
              </Typography>
            </AccordionSummary>
            <AccordionDetails>
              {schemaSet
                .filter((s) => s.type === "addon")
                .sort((a, b) => (a.workload?.title < b.workload?.title ? -1 : 1))
                .map((s) => (
                  <Grid item>
                    <PatternServiceForm formData={deployServiceConfig[s.workload?.title]} onChange={handleChangeData} schemaSet={s} onSubmit={handleSubmit} onDelete={handleDelete} namespace={ns} />
                  </Grid>
                ))}
            </AccordionDetails>
          </Accordion>
        </Grid>
        <Grid item xs={12} md={6} >
          <CodeEditor />
        </Grid>
      </Grid>
    </>
  );

  function CustomButton({ title, onClick }) {
    return <Button
      fullWidth
      color="primary"
      variant="contained"
      onClick={onClick}
      style={{
        marginTop : "16px",
        padding : "10px"
      }}
    >
      {title}
    </Button>;
  }

  function accordion(schema) {
    const patternName = schema?.workload?.title;

    return <Accordion
      expanded={expanded.includes(patternName)}
      onChange={() => handleExpansion(patternName)}
      style={{ width : '100%' }}
    >
      <AccordionSummary expandIcon={<ExpandMoreIcon />}>
        <Typography variant="h6">
          {patternName || "Expand More"}
        </Typography>
      </AccordionSummary>
      <AccordionDetails>
        <PatternServiceForm formData={deployServiceConfig[getPatternAttributeName(schema.workload)]} onChange={handleChangeData} schemaSet={schema} onSubmit={(val) => handleSubmit(val, patternName)} onDelete={(val) => handleDelete(val, patternName)} namespace={ns} />
      </AccordionDetails>
    </Accordion>;
  }

  function CodeEditor() {
    const cardStyle = { marginBottom : "16px", position : "sticky", float : "right", minWidth : "100%" };
    const cardcontentStyle = { margin : "16px" };

    const classes = useStyles();

    return (
      <div>
        <Card style={cardStyle}>
          <CardContent style={cardcontentStyle}>
            <CodeMirror
              value={yaml}
              className={classes.codeMirror}
              options={{
                theme : "material",
                lineNumbers : true,
                lineWrapping : true,
                gutters : ["CodeMirror-lint-markers"],
                lint : true,
                mode : "text/x-yaml",
              }}
              onBlur={(a) => saveCodeEditorChanges(a)}
            />
            <CustomButton title="Save Application" onClick={() => handleSubmitFinalPattern(yaml, "", `meshery_${Math.floor(Math.random() * 100)}`, "upload")} />
            <CardActions style={{ justifyContent : "flex-end" }}>
              <Tooltip title="Update Application">
                <IconButton
                  aria-label="Update"
                  color="primary"
                  onClick={() => handleSubmitFinalPattern(yaml, application.id, application.name, "update")}
                >
                  <SaveIcon />
                </IconButton>
              </Tooltip>
              <Tooltip title="Delete Application">
                <IconButton
                  aria-label="Delete"
                  color="primary"
                  onClick={() => handleSubmitFinalPattern(yaml, application.id, application.name, "delete")}
                >
                  <DeleteIcon />
                </IconButton>
              </Tooltip>
            </CardActions>
          </CardContent>
        </Card>
      </div>
    )
  }
}