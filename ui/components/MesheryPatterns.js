import React, { useState, useEffect, useRef } from "react";
import { withStyles, makeStyles } from "@material-ui/core/styles";
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
import dataFetch, { promisifiedDataFetch } from "../lib/data-fetch";
import { CircularProgress } from "@material-ui/core";
import PatternServiceForm from "./MesheryMeshInterface/PatternServiceForm";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore"
import { Button } from "@material-ui/core";
import jsYaml from "js-yaml";
import ListAltIcon from '@material-ui/icons/ListAlt';
import PascalCaseToKebab from "../utils/PascalCaseToKebab";
import ArrowBackIcon from '@material-ui/icons/ArrowBack';

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
  }
}))

function CustomToolbar(onClick) {
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
      </>
    );
  };
}

function YAMLEditor({ pattern, onClose, onSubmit }) {
  const [yaml, setYaml] = useState("");

  return (
    <Dialog onClose={onClose} aria-labelledby="pattern-dialog-title" open fullWidth maxWidth="md">
      <DialogTitle id="pattern-dialog-title">{pattern.name}</DialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <CodeMirror
          value={pattern.pattern_file}
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
  const [showForm, setShowForm] = useState(false);

  const DEPLOY_URL = '/api/pattern/deploy';

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

  /**
   * fetchPatterns constructs the queries based on the parameters given
   * and fetches the patterns
   * @param {number} page current page
   * @param {number} pageSize items per page
   * @param {string} search search string
   * @param {string} sortOrder order of sort
   */

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

    if (type === "upload") {
      dataFetch(
        `/api/pattern`,
        {
          credentials : "include",
          method : "POST",
          body : JSON.stringify({ pattern_data : { pattern_file : data }, save : true }),
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
              <IconButton onClick={() => setShowForm({ pattern : patterns[tableMeta.rowIndex], show : true })}>
                <ListAltIcon />
              </IconButton>
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

  async function showModal() {
    let response = await modalRef.current.show({ title : "Delete Pattern?",

      subtitle : "Are you sure you want to delete this pattern?",

      options : ["yes", "no"], })
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
    customToolbar : CustomToolbar(uploadHandler),

    onCellClick : (_, meta) => meta.colIndex !== 3 && setSelectedRowData(patterns[meta.rowIndex]),

    onRowsDelete : async function handleDelete(row) {
      let response = await showModal()
      console.log(response)
      if (response === "yes") {
        const fid = Object.keys(row.lookup).map(idx => patterns[idx]?.id)
        fid.forEach(fid => deletePattern(fid))
      }
      if (response === "no")
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
      {showForm &&
        <PatternForm onSubmit={handleSubmit} show={setShowForm} pattern={showForm.pattern} />}

      {selectedRowData && Object.keys(selectedRowData).length > 0 && (
        <YAMLEditor pattern={selectedRowData} onClose={resetSelectedRowData()} onSubmit={handleSubmit} />
      )}
      {
        !showForm && <MUIDataTable
          title={<div className={classes.tableHeader}>Patterns</div>}
          data={patterns}
          columns={columns}
          // @ts-ignore
          options={options}
          className={classes.muiRow}
        />
      }
      <PromptComponent ref={modalRef} />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

const mapStateToProps = (state) => {
  return { user : state.get("user").toObject(), };
};

// @ts-ignore
export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(MesheryPatterns)));

// --------------------------------------------------------------------------------------------------
// -------------------------------------------- PATTERNS FORM ---------------------------------------
// --------------------------------------------------------------------------------------------------

function PatternForm({ pattern, onSubmit, show }) {
  const [schemaSet, setSchemaSet] = useState();
  const [deployServiceConfig, setDeployServiceConfig] = useState(getPatternJson() || {});
  const [yaml, setYaml] = useState("");
  const [expanded, setExpanded] = useState([]);
  // const [changedYaml, setChangedYaml] = useState("");
  const classes = useStyles();

  function getPatternJson() {
    const patternString = pattern.pattern_file;
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
    deployConfig.name = pattern.name;
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
            Edit Pattern Configuration of <i>{`${pattern.name}`}</i>
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
            <CustomButton title="Save Pattern" onClick={() => handleSubmitFinalPattern(yaml, "", `meshery_${Math.floor(Math.random() * 100)}`, "upload")} />
            <CardActions style={{ justifyContent : "flex-end" }}>
              <Tooltip title="Update Pattern">
                <IconButton
                  aria-label="Update"
                  color="primary"
                  onClick={() => handleSubmitFinalPattern(yaml, pattern.id, pattern.name, "update")}
                >
                  <SaveIcon />
                </IconButton>
              </Tooltip>
              <Tooltip title="Delete Pattern">
                <IconButton
                  aria-label="Delete"
                  color="primary"
                  onClick={() => handleSubmitFinalPattern(yaml, pattern.id, pattern.name, "delete")}
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
