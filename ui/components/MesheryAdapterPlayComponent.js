import {
  Card, CardActions, CardHeader, Chip, Dialog, DialogActions, DialogContent,
  DialogTitle, Divider, FormControl, FormControlLabel, FormGroup, FormLabel, Grid,
  IconButton, Menu, MenuItem, Switch, Typography, withStyles } from "@material-ui/core";
import { blue } from "@material-ui/core/colors";
import NoSsr from "@material-ui/core/NoSsr";
import AddIcon from "@material-ui/icons/Add";
import DeleteIcon from "@material-ui/icons/Delete";
import PlayIcon from "@material-ui/icons/PlayArrow";
import { withRouter } from "next/router";
import PropTypes from "prop-types";
import React, { useEffect, useRef, useState } from "react";
import { Controlled as CodeMirror } from "react-codemirror2";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import dataFetch from "../lib/data-fetch";
import { setK8sContexts, updateProgress } from "../lib/store";
import { ctxUrl, getK8sClusterIdsFromCtxId } from "../utils/multi-ctx";
import fetchAvailableAddons from "./graphql/queries/AddonsStatusQuery";
import fetchAvailableNamespaces from "./graphql/queries/NamespaceQuery";
import MesheryMetrics from "./MesheryMetrics";
import MesheryResultDialog from "./MesheryResultDialog";
import ReactSelectWrapper from "./ReactSelectWrapper";
import ConfirmationMsg from "./ConfirmationModal";
import { iconMedium } from "../css/icons.styles";
import { ACTIONS } from "../utils/Enum";
import { getModelByName } from "../api/meshmodel";
import { EVENT_TYPES } from "../lib/event-types";
import { withNotify } from "../utils/hooks/useNotification";
import { usePrevious } from "../hooks/usePrevious";

const styles = (theme) => ({
  smWrapper : { backgroundColor : theme.palette.secondary.elevatedComponents2, },
  buttons : { width : "100%", },
  button : {
    marginTop : theme.spacing(3),
    marginLeft : theme.spacing(1),
  },
  margin : { margin : theme.spacing(1), },
  alreadyConfigured : {
    textAlign : "center",
    padding : theme.spacing(20),
  },
  chip : {
    height : "50px",
    fontSize : "15px",
    position : "relative",
    top : theme.spacing(0.5),
    [theme.breakpoints.down("md")] : { fontSize : "12px", },
  },
  colorSwitchBase : {
    color : blue[300],
    "&$colorChecked" : {
      color : blue[500],
      "& + $colorBar" : { backgroundColor : blue[500], },
    },
  },
  colorBar : {},
  colorChecked : {},
  uploadButton : {
    margin : theme.spacing(1),
    marginTop : theme.spacing(3),
  },
  fileLabel : { width : "100%", },
  editorContainer : { width : "100%", },
  deleteLabel : { paddingRight : theme.spacing(2), },
  alignRight : { textAlign : "right", },
  alignLeft : {
    textAlign : "left",
    marginLeft : theme.spacing(1),
  },
  padLeft : { paddingLeft : theme.spacing(0.25), },
  padRight : { paddingRight : theme.spacing(0.25), },
  deleteRight : { float : "right", },
  expTitleIcon : {
    width : theme.spacing(3),
    display : "inline",
    verticalAlign : "middle",
  },
  expIstioTitleIcon : {
    width : theme.spacing(2),
    display : "inline",
    verticalAlign : "middle",
    marginLeft : theme.spacing(0.5),
    marginRight : theme.spacing(0.5),
  },
  expTitle : {
    display : "inline",
    verticalAlign : "middle",
  },
  icon : { width : theme.spacing(2.5), },
  tableHeader : {
    fontWeight : "bolder",
    fontSize : 18,
  },
  secondaryTable : {
    borderRadius : 10,
    backgroundColor : "#f7f7f7",
  },
  paneSection : {
    backgroundColor : theme.palette.secondary.elevatedComponents,
    padding : theme.spacing(3),
    borderRadius : 4,
  },
  chipNamespace : {
    gap : '2rem',
    margin : "0px",
  },
  cardMesh : { margin : "-8px 0px", },
  inputContainer : {
    flex : '1',
    minWidth : '250px'
  },
  card : {
    height : '100%',
    display : 'flex',
    flexDirection : 'column'
  },
  ctxIcon : {
    display : 'inline',
    verticalAlign : 'text-top',
    width : theme.spacing(2.5),
    marginLeft : theme.spacing(0.5),
  },
  ctxChip : {
    backgroundColor : "white",
    cursor : "pointer",
    marginRight : theme.spacing(1),
    marginLeft : theme.spacing(1),
    marginBottom : theme.spacing(1),
    height : "100%",
    padding : theme.spacing(0.5)
  },
  text : {
    padding : theme.spacing(1)
  }
});

const MesheryAdapterPlayComponent = (props) => {
  // Ref
  const cmEditorAdd = useRef(null);
  const cmEditorDel = useRef(null);
  const addIconEles = useRef({});
  const delIconEles = useRef({});

  // States
  const [adapter,] = useState(props.adapter);
  const [namespaceList, setNameSpaceList] = useState([]);
  const [namespaceSubscription, setNamespaceSubscription] = useState(null);
  const [menuState, setMenuState] = useState({});
  const [, setActiveContexts] = useState([]);
  const [addonSwitchGroup, setAddonSwitchGroup] = useState({});
  const [versionList, setVersionList] = useState([]);
  const [version, setVersion] = useState({ labeL : "", value : "" });
  const [versionError, setVersionError] = useState(false);
  const [, setSelectionError] = useState(false);
  const [cmEditorValDel, setCmEditorValDel] = useState("");
  const [, setCmEditorValDelError] = useState(false);
  const [cmEditorValAdd, setCmEditorValAdd] = useState("");
  const [, setCmEditorValAddError] = useState(false);
  const [namespace, setNamespace] = useState({
    value : "default",
    label : "default",
  });
  const [namespaceError, setNamespaceError] = useState(false);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [customDialogAdd, setCustomDialogAdd] = useState(false);
  const [customDialogDel, setCustomDialogDel] = useState(false);
  const [category, setCategory] = useState(0);
  const [selectedOp, setSelectedOp] = useState("");
  const [isDeleteOp, setIsDeleteOp] = useState(false);
  const [operationName, setOperationName] = useState("");
  const [modalOpen, setModalOpen] = useState(false);

  const prevProps = usePrevious(props);

  // initializing menuState;
  useEffect(() => {
    if (adapter && adapter.ops) {
      // NOTE: this will have to updated to match the categories
      const menu = {};
      [0, 1, 2, 3, 4].forEach((i) => {
        menu[i] = {
          add : false,
          delete : false,
        };
      });
      setMenuState(menu);
    }
  }, [adapter]);

  useEffect(() => {
    const meshname = mapAdapterNameToMeshName(adapter.name);
    const variables = { type : meshname, k8sClusterIDs : getK8sClusterIds() };
    initSubscription();
    getMeshVersions();
    if (props.selectedK8sContexts) {
      if (props.selectedK8sContexts.includes("all")) {
        let active = [];
        props.k8sconfig.forEach((ctx) => {
          active.push(ctx.contextID);
        });
        setActiveContexts(active);
      } else {
        setActiveContexts(props.selectedK8sContexts);
      }
    }

    fetchAvailableAddons(variables).subscribe({
      next : (res) => {
        setAddonsState(res);
      },
      error : (err) => console.log("error at addon fetch: " + err),
    });
  }, []);

  useEffect(() => {
    if (prevProps?.selectedK8sContexts.length !== props?.selectedK8sContexts.length) {
      disposeSubscriptions();
      initSubscription();
    }
    if (prevProps?.adapter.name !== props?.adapter.name) {
      getMeshVersions();
    }
  }, [prevProps, props]);

  const initSubscription = () => {
    const namespaceSubscriptionList = fetchAvailableNamespaces({
      k8sClusterIDs : getK8sClusterIds(),
    }).subscribe({
      next : (res) => {
        let namespaces = [];
        res?.namespaces?.map((ns) => {
          namespaces.push({
            value : ns?.namespace,
            label : ns?.namespace,
          });
        });
        if (namespaces.length === 0) {
          namespaces.push({
            value : "default",
            label : "default",
          });
        }
        namespaces.sort((a, b) => (a.value > b.value ? 1 : -1));
        setNameSpaceList(namespaces);
      },
      error : (err) => console.log("error at namespace fetch: " + err),
    });

    setNamespaceSubscription(namespaceSubscriptionList);
  };

  const disposeSubscriptions = () => {
    if (namespaceSubscription) {
      namespaceSubscription.unsubscribe();
    }
  };

  const getMeshVersions = () => {
    const activeMesh = adapter.name;
    if (activeMesh)
      getModelByName(activeMesh.toLowerCase()).then((res) => {
        let uniqueVersions = [...new Set(res?.models?.map((model) => model?.version))].reverse();
        if (uniqueVersions.length === 0) {
          uniqueVersions = [""];
        }
        let updatedVersionList = uniqueVersions.map((version) => ({
          value : version,
          label : version,
        }));
        setVersionList(updatedVersionList);
        setVersion(versionList[0]);
      });
  };

  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(props?.selectedK8sContexts, props.k8sconfig);
  };

  const mapAdapterNameToMeshName = (name) => {
    if (name?.toLowerCase() === "istio") return "ISTIO";

    return "ALL";
  };

  const setAddonsState = (data) => {
    const meshname = adapter.name;
    const localState = {};
    data?.addonsState?.forEach((addon) => {
      if (addon.owner === meshname) {
        const name = addon.name !== "jaeger-collector" ? addon.name : "jaeger";
        localState[`${name}-addon`] = true;
      }
    });
    setAddonSwitchGroup(localState);
  };

  const handleNamespaceChange = (newValue) => {
    if (typeof newValue !== "undefined") {
      setNamespace(newValue);
      setNamespaceError(false);
    } else {
      setNamespaceError(true);
    }
  };

  const handleVersionChange = (newValue) => {
    if (typeof newValue !== "undefined") {
      setVersion(newValue);
      setNamespaceError(false);
    } else {
      setVersionError(true);
    }
  };

  const handleModalClose = (isDelete) => {
    return () => {
      if (isDelete) {
        setCustomDialogDel(false);
      } else {
        setCustomDialogAdd(false);
      }
    };
  };

  const resetSelectedRowData = () => {
    return () => {
      setSelectedRowData(null);
    };
  };

  const handleSubmit = (cat, selectedOp, deleteOp = false) => {
    return () => {
      handleOpen();

      const filteredOp = adapter.ops.filter(({ key }) => key === selectedOp);
      if (selectedOp === "" || typeof filteredOp === "undefined" || filteredOp.length === 0) {
        setSelectionError(true);
        return;
      }
      if (deleteOp) {
        if (selectedOp === "custom" && (cmEditorValDel === "" || cmEditorDel.current.state.lint.marked.length > 0)) {
          setCmEditorValDelError(true);
          setSelectionError(true);
          return;
        }
      } else if (
        selectedOp === "custom" &&
        (cmEditorValAdd === "" || cmEditorAdd.current.state.lint.marked.length > 0)
      ) {
        setCmEditorValAddError(true);
        setSelectionError(true);
        return;
      }
      if (namespace && namespace.value === "") {
        setNamespaceError(true);
        return;
      }

      if (version?.value === "") {
        setVersionError(true);
        return;
      }
      const operationName = selectedOp
        .replaceAll("_", " ")
        .split(" ")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ");
      setCategory(cat);
      setSelectedOp(selectedOp);
      setIsDeleteOp(deleteOp);
      setOperationName(operationName);
    };
  };

  const submitOp = (cat, selectedOp, deleteOp = false) => {
    const data = {
      adapter : adapter?.adapter_location,
      query : selectedOp,
      namespace : namespace?.value,
      customBody : deleteOp ? cmEditorValDel : cmEditorValAdd,
      deleteOp : deleteOp ? "on" : "",
      version : version?.value,
    };

    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
      .join("&");
    props.updateProgress({ showProgress : true });
    handleClose();
    dataFetch(
      ctxUrl("/api/system/adapter/operation", props.selectedK8sContexts),
      {
        method : "POST",
        credentials : "include",
        headers : {
          "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8",
        },
        body : params,
      },
      (result) => {
        props.updateProgress({ showProgress : false });
        const updatedMenu = { ...menuState };
        updatedMenu[cat][deleteOp ? "delete" : "add"] = false;
        setMenuState(updatedMenu);

        if (deleteOp) {
          setCustomDialogDel(false);
        } else {
          setCustomDialogAdd(false);
        }

        if (typeof result !== "undefined") {
          const notify = props.notify;
          notify({
            message : "Operation executing...",
            event_type : EVENT_TYPES.INFO,
          });
        }
      },
      handleError(cat, deleteOp, selectedOp)
    );
  };

  const handleAdapterClick = (adapterLoc) => () => {
    props.updateProgress({ showProgress : true });
    dataFetch(
      `/api/system/adapters?adapter=${encodeURIComponent(adapterLoc)}`,
      {
        credentials : "include",
      },
      (result) => {
        props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          const notify = props.notify;
          notify({
            message : "Adapter pinged!",
            event_type : EVENT_TYPES.SUCCESS,
          });
        }
      },
      handleError("Could not ping adapter.")
    );
  };

  const handleError = (cat, deleteOp, selectedOp) => {
    return (error) => {
      if (cat && deleteOp) {
        const updatedMenu = {
          ...menuState,
          [cat] : {
            [deleteOp ? "delete" : "add"] : false,
          },
        };

        setMenuState(updatedMenu);
        if (deleteOp) {
          setCustomDialogDel(false);
        } else {
          setCustomDialogAdd(false);
        }
      }
      setAddonSwitchGroup({ ...addonSwitchGroup, [selectedOp] : deleteOp });

      props.updateProgress({ showProgress : false });
      const notify = props.notify;
      notify({
        message : `Operation submission failed: ${error}`,
        event_type : EVENT_TYPES.ERROR,
        details : error.toString(),
      });
    };
  };

  /**
   * generateMenu generates the management menus for the adapater management plane
   * @param {*} cat
   * @param {boolean} isDelete if set to true, a delete menu will be generated
   * @param {{key: string, value: string, category?: number}[]} selectedAdapterOps is the array of the meshery adapaters
   *
   * @returns {JSX.Element}
   */
  const generateMenu = (cat, isDelete, selectedAdapterOps) => {
    const ele = !isDelete ? addIconEles.current[cat] : delIconEles.current[cat];
    return (
      <Menu
        id="long-menu"
        anchorEl={ele}
        keepMounted
        open={menuState[cat] && menuState[cat][isDelete ? "delete" : "add"]}
        onClose={addDelHandleClick(cat, isDelete)}
      >
        {selectedAdapterOps
          .sort((adap1, adap2) => adap1.value.localeCompare(adap2.value))
          .map(({ key, value }) => (
            <MenuItem key={`${key}_${new Date().getTime()}`} onClick={handleSubmit(cat, key, isDelete)}>
              {value}
            </MenuItem>
          ))}
      </Menu>
    );
  };

  const handleOpen = () => {
    setModalOpen(true);
  };

  const handleClose = () => {
    setModalOpen(false);
  };

  const generateYAMLEditor = (cat, isDelete) => {
    const { adapter } = props;

    return (
      <Dialog
        onClose={handleModalClose(isDelete)}
        aria-labelledby="adapter-dialog-title"
        open={isDelete ? customDialogDel : customDialogAdd}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle id="adapter-dialog-title" onClose={handleModalClose(isDelete)}>
          {adapter.name} Adapter - Custom YAML
          {isDelete ? "(delete)" : ""}
        </DialogTitle>
        <Divider variant="fullWidth" light />
        <DialogContent>
          <Grid container spacing={5}>
            <Grid item xs={6}>
              <ReactSelectWrapper
                label="Namespace"
                value={namespace}
                error={namespaceError}
                options={namespaceList}
                onChange={handleNamespaceChange}
              />
            </Grid>
            <Grid item xs={6}>
              <ReactSelectWrapper
                label="Version"
                value={version}
                error={versionError}
                options={versionList}
                onChange={handleVersionChange}
              />
            </Grid>
            <Grid item xs={12}>
              <CodeMirror
                editorDidMount={(editor) => {
                  if (isDelete) {
                    cmEditorDel.current = editor;
                  } else {
                    cmEditorAdd.current = editor;
                  }
                }}
                value={isDelete ? cmEditorValDel : cmEditorValAdd}
                options={{
                  theme : "material",
                  lineNumbers : true,
                  lineWrapping : true,
                  gutters : ["CodeMirror-lint-markers"],
                  lint : true,
                  mode : "text/x-yaml",
                }}
                onBeforeChange={(editor, data, value) => {
                  if (isDelete) {
                    setCmEditorValDel(value);
                  } else {
                    setCmEditorValAdd(value);
                  }
                  if (isDelete) {
                    if (value !== "" && cmEditorDel.current.state.lint.marked.length === 0) {
                      setSelectionError(false);
                      setCmEditorValDelError(false);
                    }
                  } else if (value !== "" && cmEditorAdd.current.state.lint.marked.length === 0) {
                    setSelectionError(false);
                    setCmEditorValAddError(false);
                  }
                }}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <Divider variant="fullWidth" light />
        <DialogActions>
          <IconButton aria-label="Apply" color="primary" onClick={handleSubmit(cat, "custom", isDelete)}>
            {/* <FontAwesomeIcon icon={faArrowRight} transform="shrink-4" fixedWidth /> */}
            {!isDelete && <PlayIcon style={iconMedium} />}
            {isDelete && <DeleteIcon style={iconMedium} />}
          </IconButton>
        </DialogActions>
      </Dialog>
    );
  };

  const addDelHandleClick = (cat, isDelete) => {
    return () => {
      const updatedMenu = { ...menuState };

      updatedMenu[cat][isDelete ? "delete" : "add"] = !menuState[cat][isDelete ? "delete" : "add"];

      setMenuState(updatedMenu);

      let dlgv = isDelete ? customDialogDel : customDialogAdd;
      if (cat === 4) {
        dlgv = !dlgv;
      }

      if (isDelete) {
        setCustomDialogDel(dlgv);
      } else {
        setCustomDialogAdd(dlgv);
      }
    };
  };

  const generateCardForCategory = (cat) => {
    if (typeof cat === "undefined") {
      cat = 0;
    }
    const { classes, adapter } = props;
    if (adapter && adapter.ops && adapter.ops.length > 0) {
      adapter.ops.forEach(({ category }) => {
        if (typeof category === "undefined") {
          category = 0;
        }
        if (filteredOps.indexOf(category) === -1) {
          filteredOps.push(category);
        }
      });
      filteredOps.sort();
    }
    let selectedAdapterOps =
      adapter && adapter.ops
        ? adapter.ops.filter(({ category }) => (typeof category === "undefined" && cat === 0) || category === cat)
        : [];
    let content;
    let description;
    switch (cat) {
      case 0:
        content = "Manage Service Mesh Lifecycle";
        description = "Deploy a service mesh or SMI adapter on your cluster.";
        break;

      case 1:
        content = "Manage Sample Application Lifecycle";
        description = "Deploy sample applications on/off the service mesh.";
        break;

      case 2:
        content = "Apply Service Mesh Configuration";
        description = "Configure your service mesh using some pre-defined options.";
        selectedAdapterOps = selectedAdapterOps.filter((ops) => !ops.value.startsWith("Add-on:"));
        break;

      case 3:
        content = "Validate Service Mesh Configuration";
        description = "Validate your service mesh configuration against best practices.";
        break;

      case 4:
        content = "Apply Custom Configuration";
        description = "Customize the configuration of your service mesh.";
        break;
    }
    return (
      <Card className={classes.card}>
        <CardHeader title={content} subheader={description} style={{ flexGrow : 1 }} />
        <CardActions disableSpacing>
          <IconButton
            aria-label="install"
            ref={(ch) => (addIconEles.current[cat] = ch)}
            onClick={addDelHandleClick(cat, false)}
          >
            {cat !== 4 ? <AddIcon style={iconMedium} /> : <PlayIcon style={iconMedium} />}
          </IconButton>
          {cat !== 4 && generateMenu(cat, false, selectedAdapterOps)}
          {cat === 4 && generateYAMLEditor(cat, false)}
          {cat !== 3 && (
            <div className={classes.fileLabel}>
              <IconButton
                aria-label="delete"
                ref={(ch) => (delIconEles.current[cat] = ch)}
                className={classes.deleteRight}
                onClick={addDelHandleClick(cat, true)}
              >
                <DeleteIcon style={iconMedium} />
              </IconButton>
              {cat !== 4 && generateMenu(cat, true, selectedAdapterOps)}
              {cat === 4 && generateYAMLEditor(cat, true)}
            </div>
          )}
        </CardActions>
      </Card>
    );
  };

  /**
   * extractAddonOperations returns an array of operations
   * which have a prefix "Addon:"
   * @param {number} addonOpsCat category for addon operations
   * @returns {{category: number, key: string, value: string}[]}
   */
  // Converted
  const extractAddonOperations = (addonOpsCat) => {
    return props.adapter.ops.filter(({ category, value }) => category === addonOpsCat && value?.startsWith("Add-on:"));
  };

  /**
   * generateAddonSwitches creates a switch based ui for the addon operations
   * @param {{category: number, key: string, value: string}[]} selectedAdapterOps available adapter operations
   * @returns {JSX.Element}
   */
  // Converted
  const generateAddonSwitches = (selectedAdapterOps) => {
    if (!selectedAdapterOps.length) return null;

    return (
      <FormControl component="fieldset" style={{ padding : "1rem" }}>
        <FormLabel component="legend">Customize Addons</FormLabel>
        <FormGroup>
          {selectedAdapterOps
            .map((ops) => ({ ...ops, value : ops.value.replace("Add-on:", "") }))
            .sort((ops1, ops2) => ops1.value.localeCompare(ops2.value))
            .map((ops) => (
              <FormControlLabel
                control={
                  <Switch
                    color="primary"
                    checked={!!addonSwitchGroup[ops.key]}
                    onChange={(ev) => {
                      const updatedAddonSwitchGroup = {
                        ...addonSwitchGroup,
                        [ev.target.name] : ev.target.checked,
                      };
                      setAddonSwitchGroup(updatedAddonSwitchGroup);
                      submitOp(ops.category, ops.key, !!addonSwitchGroup[ops.key]);
                    }}
                    name={ops.key}
                  />
                }
                label={ops.value}
                key={ops.key}
              />
            ))}
        </FormGroup>
      </FormControl>
    );
  };

  /**
   * renderGrafanaCustomCharts takes in the configuration and renders
   * the grafana boards. If the configuration is empty then it renders
   * a note directing a user to install grafana and prometheus
   * @param {Array<{ board: any, panels: Array<any>, templateVars: Array<any>}>} boardConfigs grafana board configs
   * @param {string} grafanaURL grafana URL
   * @param {string} grafanaAPIKey grafana API key
   */
  // Converted
  const renderGrafanaCustomCharts = (boardConfigs, grafanaURL, grafanaAPIKey) => {
    return (
      <MesheryMetrics
        boardConfigs={boardConfigs}
        grafanaAPIKey={grafanaAPIKey}
        grafanaURL={grafanaURL}
        handleGrafanaChartAddition={() => props.router.push("/settings/#metrics")}
      />
    );
  };

  const { classes } = props;
  const [imageSrc, setImageSrc] = useState("");
  const [filteredOps, setFilteredOps] = useState([]);

  useEffect(() => {
    if (adapter?.name) {
      const adapterName = adapter.name.split(" ").join("").toLowerCase();
      const path = "/static/img/" + adapterName + ".svg";
      setImageSrc(path);
    }

    let filter = filteredOps;
    if (adapter && adapter.ops && adapter.ops.length > 0) {
      adapter.ops.forEach(({ category }) => {
        if (typeof category === "undefined") {
          category = 0;
        }
        if (filter.indexOf(category) === -1) {
          filter.push(category);
        }
      });
      filter.sort();
      setFilteredOps(filter);
    }
  }, [adapter]);

  let adapterChip = (
    <Chip
      label={adapter.adapter_location}
      onClick={handleAdapterClick(adapter.adapter_location)}
      icon={<img src={imageSrc} className={classes.icon} />}
      className={classes.chip}
      variant="outlined"
    />
  );

  return (
    <NoSsr>
      {selectedRowData && selectedRowData !== null && Object.keys(selectedRowData).length > 0 && (
        <MesheryResultDialog rowData={selectedRowData} close={resetSelectedRowData()} />
      )}
      <React.Fragment>
        <div className={classes.smWrapper}>
          <Grid container spacing={2} direction="row" alignItems="flex-start">
            {/* SECTION 1 */}
            <Grid item xs={12}>
              <div className={classes.paneSection}>
                <Typography align="left" variant="h6" style={{ margin : "0 0 2.5rem 0" }}>
                  Manage Service Mesh
                </Typography>
                <Grid container spacing={4}>
                  <Grid
                    container
                    item
                    xs={12}
                    alignItems="flex-start"
                    justify="space-between"
                    className={classes.chipNamespace}
                  >
                    <div>{adapterChip}</div>
                    <div className={classes.inputContainer}>
                      <ReactSelectWrapper
                        label="Namespace"
                        value={namespace}
                        error={namespaceError}
                        options={namespaceList}
                        onChange={handleNamespaceChange}
                      />
                    </div>
                    <div className={classes.inputContainer}>
                      <ReactSelectWrapper
                        label="Version"
                        value={version}
                        error={versionError}
                        options={versionList}
                        onChange={handleVersionChange}
                      />
                    </div>
                  </Grid>
                  <Grid container spacing={1}>
                    <Grid
                      container
                      item
                      lg={!extractAddonOperations(2).length ? 12 : 10}
                      xs={12}
                      spacing={2}
                      className={classes.cardMesh}
                    >
                      {filteredOps.map((val, i) => (
                        <Grid item lg={3} md={4} xs={12} key={`adapter-card-${i}`}>
                          {generateCardForCategory(val)}
                        </Grid>
                      ))}
                    </Grid>
                    <Grid container item lg={2} xs={12}>
                      <Grid item xs={12} md={4}>
                        {generateAddonSwitches(extractAddonOperations(2))}
                      </Grid>
                    </Grid>
                  </Grid>
                </Grid>
              </div>
            </Grid>
            {/* SECTION 2 */}
            <Grid item xs={12}>
              <div className={classes.paneSection}>
                {renderGrafanaCustomCharts(
                  props.grafana.selectedBoardsConfigs,
                  props.grafana.grafanaURL,
                  props.grafana.grafanaAPIKey
                )}
              </div>
            </Grid>
          </Grid>
        </div>
        <ConfirmationMsg
          open={modalOpen}
          handleClose={handleClose}
          submit={{
            deploy : () => submitOp(category, selectedOp, false),
            unDeploy : () => submitOp(category, selectedOp, true),
          }}
          isDelete={isDeleteOp}
          title={operationName}
          tab={isDeleteOp ? ACTIONS.UNDEPLOY : ACTIONS.DEPLOY}
        />
      </React.Fragment>
    </NoSsr>
  );
};

MesheryAdapterPlayComponent.propTypes = {
  classes : PropTypes.object.isRequired,
  adapter : PropTypes.object.isRequired,
};

const mapStateToProps = (st) => {
  const grafana = st.get("grafana").toJS();
  const k8sconfig = st.get("k8sConfig");
  const selectedK8sContexts = st.get('selectedK8sContexts');

  return { grafana : { ...grafana, ts : new Date(grafana.ts) }, selectedK8sContexts, k8sconfig };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress : bindActionCreators(updateProgress, dispatch),
  setK8sContexts : bindActionCreators(setK8sContexts, dispatch)
});

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withNotify(MesheryAdapterPlayComponent)))
);
