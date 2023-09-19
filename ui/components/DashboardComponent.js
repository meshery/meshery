
import {
  Button,Chip,
  MenuItem, NoSsr, Paper, Select, TableCell, TableSortLabel, Tooltip, Typography
} from "@material-ui/core";
// import {Table, TableBody, TableContainer, TableHead, TableRow,} from "@material-ui/core"
import blue from "@material-ui/core/colors/blue";
import { makeStyles } from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import { MuiThemeProvider } from "@material-ui/core/styles";
import AddIcon from "@material-ui/icons/AddCircleOutline";
import { useRouter } from "next/router";
import PropTypes from "prop-types";
import React,{ useEffect, useState, useRef } from "react";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import dataFetch from "../lib/data-fetch";
import { updateGrafanaConfig, updateProgress, updatePrometheusConfig, updateTelemetryUrls } from "../lib/store";
import { getK8sClusterIdsFromCtxId, getK8sClusterNamesFromCtxId } from "../utils/multi-ctx";
// import { versionMapper } from "../utils/nameMapper";
import { submitGrafanaConfigure } from "./telemetry/grafana/GrafanaComponent";
import fetchAvailableAddons from "./graphql/queries/AddonsStatusQuery";
import fetchControlPlanes from "./graphql/queries/ControlPlanesQuery";
import fetchDataPlanes from "./graphql/queries/DataPlanesQuery";

import subscribeClusterResources from "./graphql/subscriptions/ClusterResourcesSubscription";
import fetchAvailableNamespaces from "./graphql/queries/NamespaceQuery";


import MUIDataTable from "mui-datatables";
import Popup from "./Popup";
import { iconMedium } from "../css/icons.styles";
import { configurationTableTheme, configurationTableThemeDark } from '../themes/configurationTableTheme';
import DashboardMeshModelGraph from './Dashboard/DashboardMeshModelGraph'
import ConnectionStatsChart from "./Dashboard/ConnectionCharts.js";
import { EVENT_TYPES } from "../lib/event-types";
import { useNotification } from "../utils/hooks/useNotification";

const useStyles = makeStyles((theme) => ({
  rootClass : {
    backgroundColor : theme.palette.secondary.elevatedComponents2,
  },
  datatable : {
    boxShadow : "none",
  },
  chip : {
    marginRight : theme.spacing(1),
    marginBottom : theme.spacing(1),
  },
  buttons : {
    display : "flex",
    justifyContent : "flex-end",
  },
  button : {
    marginTop : theme.spacing(3),
    marginLeft : theme.spacing(1),
  },
  link : {
    cursor : "pointer",
    textDecoration : "none"
  },
  metricsButton : {
    width : "240px",
  },
  alreadyConfigured : {
    textAlign : "center",
  },
  margin : {
    margin : theme.spacing(1),
  },
  colorSwitchBase : {
    color : blue[300],
    "&$colorChecked" : {
      color : blue[500],
      "& + $colorBar" : {
        backgroundColor : blue[500],
      },
    },
  },
  colorBar : {},
  colorChecked : {},
  fileLabel : {
    width : "100%",
  },
  fileLabelText : {},
  inClusterLabel : {
    paddingRight : theme.spacing(2),
  },
  alignCenter : {
    textAlign : "center",
  },
  icon : {
    width : theme.spacing(2.5),
  },
  istioIcon : {
    width : theme.spacing(1.5),
  },
  settingsIcon : {
    width : theme.spacing(2.5),
    paddingRight : theme.spacing(0.5),
  },
  addIcon : {
    width : theme.spacing(2.5),
    paddingRight : theme.spacing(0.5),
  },
  cardHeader : {
    fontSize : theme.spacing(2),
  },
  card : {
    height : "100%",
    marginTop : theme.spacing(2),
  },
  cardContent : {
    height : "100%",
  },
  redirectButton : {
    marginLeft : "-.5em",
    color : "#000",
  },
  dashboardSection : {
    backgroundColor : theme.palette.secondary.elevatedComponents,
    padding : theme.spacing(2),
    borderRadius : 4,
    height : "100%",
    marginBottom : theme.spacing(2),
  },
}));

const DashboardComponent=(props) => {

  const router = useRouter();
  const classes = useStyles();
  const { notify } = useNotification();

  const _isMounted = useRef(false);
  // const [meshAdapters, setMeshAdapters] = useState(props.meshAdapters);
  // const [contextsFromFile, setContextsFromFile] = useState([]);
  // const [availableAdapters, setAvailableAdapters] = useState([]);
  const [mts, setMts] = useState(new Date());
  // const [meshLocationURLError, setMeshLocationURLError] = useState(false);
  // const [grafanaUrl, setGrafanaUrl] = useState(props.grafana.grafanaURL);
  const [prometheusURL, setPrometheusURL] = useState(props.prometheus.prometheusURL);
  // const [k8sfileError, setK8sfileError] = useState(false);
  // const [kts, setKts] = useState(new Date());
  const [grafana, setGrafana] = useState(props.grafana);
  const [prometheus, setPrometheus] = useState(props.prometheus);
  // const [urlError, setUrlError] = useState(false);
  // const [grafanaConfigSuccess, setGrafanaConfigSuccess] = useState(props.grafana.grafanaURL !== "");
  // const [grafanaBoardSearch, setGrafanaBoardSearch] = useState("");
  const [grafanaURL, setGrafanaURL] = useState(props.grafana.grafanaURL);
  // const [grafanaAPIKey, setGrafanaAPIKey] = useState(props.grafana.grafanaAPIKey);
  // const [grafanaBoards, setGrafanaBoards] = useState(props.grafana.grafanaBoards);
  // const [selectedBoardsConfigs, setSelectedBoardsConfigs] = useState(props.grafana.selectedBoardsConfigs);
  // const [ts, setTs] = useState(props.grafana.ts);
  const [meshScan, setMeshScan] = useState([]);
  // const [activeMeshScanNamespace, setActiveMeshScanNamespace] = useState({});
  // const [meshScanNamespaces, setMeshScanNamespaces] = useState({});
  // const [isMetricsConfigured, setIsMetricsConfigured] = useState(grafana.grafanaURL !== '' && prometheus.prometheusURL !== '');
  const [controlPlaneState, setControlPlaneState] = useState("");
  const [dataPlaneState, setDataPlaneState] = useState("");
  const [clusterResources, setClusterResources] = useState([]);
  const [namespaceList, setNamespaceList] = useState([]);
  const [selectedNamespace, setSelectedNamespace] = useState("default");
  const [dataPlaneSubscription, setDataPlaneSubscription] = useState(null);
  const [controlPlaneSubscription, setControlPlaneSubscription] = useState(null);
  const [clusterResourcesSubscription, setClusterResourcesSubscription] = useState(null);
  const [clusterResourcesQuery, setClusterResourcesQuery] = useState(null);
  const [namespaceQuery, setNamespaceQuery] = useState(null);
  const [telemetryQuery, setTelemetryQuery] = useState(null);
  const prevControlPlaneState = useRef();
  const prevDataPlaneState = useRef();
  const prevSelectedNamespace = useRef();
  const prevProps = useRef(props);
  const isMetricsConfigured = grafana.grafanaURL !== "" && prometheus.prometheusURL !== "";

  useEffect(() => {
    if (props.meshAdaptersts > mts) {
      // setMeshAdapters(props.meshAdapters);
      setMts(props.meshAdaptersts);
    }
    setGrafana(props.grafana);
    setPrometheus(props.prometheus);
    // setK8sconfig(props.k8sconfig);
  }, [props.meshAdaptersts, props.meshAdapters, props.grafana, props.prometheus, props.k8sconfig, mts]);

  const disposeWorkloadWidgetSubscription = () => {
    namespaceQuery && namespaceQuery.unsubscribe();
    clusterResourcesQuery && clusterResourcesQuery.unsubscribe();
    clusterResourcesSubscription && clusterResourcesSubscription.dispose();
  };

  const disposeSubscriptions = () => {
    if (dataPlaneSubscription) {
      dataPlaneSubscription.unsubscribe();
    }
    if (controlPlaneSubscription) {
      controlPlaneSubscription.unsubscribe();
    }
    if (telemetryQuery) {
      telemetryQuery.unsubscribe();
    }
    disposeWorkloadWidgetSubscription();
  };
  const initMeshSyncControlPlaneSubscription = () => {
    const ALL_MESH = { type : "ALL_MESH", k8sClusterIDs : getK8sClusterIds() };

    if (_isMounted.current) {
      const controlPlaneSubscription = fetchControlPlanes(ALL_MESH).subscribe({
        next : (controlPlaneRes) => {
          setControlPlaneState(controlPlaneRes);
        },
        error : (err) => console.error(err),
      });

      const dataPlaneSubscription = fetchDataPlanes(ALL_MESH).subscribe({
        next : (dataPlaneRes) => {
          setDataPlaneState(dataPlaneRes);
        },
        error : (err) => console.error(err),
      });

      setDataPlaneSubscription(dataPlaneSubscription);
      setControlPlaneSubscription(controlPlaneSubscription);
    }
  };

  const initNamespaceQuery = () => {
    const namespaceQuery = fetchAvailableNamespaces({ k8sClusterIDs : getK8sClusterIds() })
      .subscribe({
        next : res => {
          let namespaces = [];
          res?.namespaces?.map(ns => {
            namespaces.push(ns?.namespace);
          });
          namespaces.sort((a, b) => (a > b ? 1 : -1));
          setNamespaceList(namespaces);
        },
        error : err => console.log("error at namespace fetch: " + err),
      });

    setNamespaceQuery(namespaceQuery);
  };

  const initDashboardClusterResourcesQuery = () => {
    let k8s = getK8sClusterIds();

    if (_isMounted.current) {
      const clusterResourcesQuery = fetchClusterResources(k8s, selectedNamespace).subscribe({
        next : res => {
          setClusterResources(res?.clusterResources);
        },
        error : err => console.log(err),
      });

      setClusterResourcesQuery(clusterResourcesQuery);
    }
  }


  const initDashboardClusterResourcesSubscription = () => {
    let k8s = getK8sClusterIds();

    if (_isMounted.current) {
      const clusterResourcesSubscription = subscribeClusterResources((res) => {
        setClusterResources(res?.clusterResources);
      }, {
        k8scontextIDs : k8s,
        namespace : selectedNamespace
      });

      setClusterResourcesSubscription(clusterResourcesSubscription);
    }
  };

  const initTelemetryComponentQuery = () => {
    const contextIDs = getK8sClusterIds();

    if (_isMounted.current) {
      const query = fetchTelemetryCompsQuery({ contexts : contextIDs }).subscribe({
        next : (components) => {
          let prometheusURLs = [];
          let grafanaURLs = [];
          components.telemetryComps?.forEach((component) => {
            const data = { spec : JSON.parse(component.spec), status : JSON.parse(component.status) };
            if (component.name === "grafana") {
              grafanaURLs = grafanaURLs.concat(extractURLFromScanData(data));
            } else {
              prometheusURLs = prometheusURLs.concat(extractURLFromScanData(data));
            }
          });
          updateTelemetryUrls({ telemetryURLs : { "grafana" : grafanaURLs, "prometheus" : prometheusURLs } });
        },
        error : (err) => console.error(err),
      });
      setTelemetryQuery(query);
    }
  };
  useEffect(() => {
    return () => {
      _isMounted.current = false;
      disposeSubscriptions();
    };
  }, []);

  useEffect(() => {
    fetchAvailableAdapters();

    if (isMetricsConfigured) {
      fetchMetricComponents();
    }

    initMeshSyncControlPlaneSubscription();
    initDashboardClusterResourcesQuery();
    initDashboardClusterResourcesSubscription();
    initNamespaceQuery();
    initTelemetryComponentQuery();
  }, []);

  useEffect(() => {
    let updateControlPlane = false;
    let updateDataPlane = false;

    if (JSON.stringify(prevControlPlaneState.current) !== JSON.stringify(controlPlaneState)) {
      updateControlPlane = true;
    }
    if (JSON.stringify(prevDataPlaneState.current) !== JSON.stringify(dataPlaneState)) {
      updateDataPlane = true;
    }

    if (updateDataPlane || updateControlPlane) {
      setMeshScanData(
        updateControlPlane ? controlPlaneState : prevControlPlaneState.current,
        updateDataPlane ? dataPlaneState : prevDataPlaneState.current
      );
    }

    // Update the refs with the current state
    prevControlPlaneState.current = controlPlaneState;
    prevDataPlaneState.current = dataPlaneState;
  }, [controlPlaneState, dataPlaneState]);

  useEffect(() => {
    if (prevProps.current.selectedK8sContexts !== props.selectedK8sContexts
            || prevProps.current.k8sconfig !== props.k8sconfig) {
      disposeSubscriptions();
      initMeshSyncControlPlaneSubscription();
      initDashboardClusterResourcesQuery();
      initDashboardClusterResourcesSubscription();
      initNamespaceQuery();
      initTelemetryComponentQuery();
    }

    // Update the ref with the current props
    prevProps.current = props;
  }, [props]);

  useEffect(() => {
    if (prevSelectedNamespace.current !== selectedNamespace) {
      disposeWorkloadWidgetSubscription();
      initDashboardClusterResourcesSubscription();
      initDashboardClusterResourcesQuery();
      initNamespaceQuery();
    }

    // Update the ref with the current state
    prevSelectedNamespace.current = selectedNamespace;
  }, [selectedNamespace]);


  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(props.selectedK8sContexts, props.k8sconfig);
  };

  const fetchMetricComponents = () => {
    let selector = { type : "ALL_MESH", k8sClusterIDs : getK8sClusterIds() };

    dataFetch(
      "/api/telemetry/metrics/config",
      {
        method : "GET",
        credentials : "include",
        headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8" },
      },
      (result) => {
        props.updateProgress({ showProgress : false });
        if (result?.prometheusURL && result?.prometheusURL !== "") {
          fetchAvailableAddons(selector).subscribe({
            next : (res) => {
              res?.addonsState?.forEach((addon) => {
                if (addon.name === "prometheus" && (!prometheusURL || prometheusURL === "")) {
                  setPrometheusURL("http://" + addon.endpoint);
                  submitPrometheusConfigure(() => console.log("Prometheus added"));
                }
              });
            },
            error : (err) => console.error("error registering prometheus:", err),
          });
        }
      },
      handleError("Error getting prometheus config")
    );

    dataFetch(
      "/api/telemetry/metrics/grafana/config",
      {
        method : "GET",
        credentials : "include",
        headers : { "Content-Type" : "application/x-www-form-urlencoded;charset=UTF-8" },
      },
      (result) => {
        props.updateProgress({ showProgress : false });
        if (result?.grafanaURL && result?.grafanaURL !== "") {
          fetchAvailableAddons(selector).subscribe({
            next : (res) => {
              res?.addonsState?.forEach((addon) => {
                if (addon.name === "grafana" && (!grafanaURL || grafanaURL === "")) {
                  setGrafanaURL("http://" + addon.endpoint);
                  submitGrafanaConfigure(() => {
                    // setSelectedBoardsConfigs((prevConfigs) => [...prevConfigs, boardConfigs]);
                    console.info("Grafana added");
                  });
                }
              });
            },
            error : (err) => console.error("error registering grafana:", err),
          });
        }
      },
      handleError("There was an error communicating with grafana config")
    );

    fetchAvailableAddons(selector).subscribe({
      next : (res) => {
        res?.addonsState?.forEach((addon) => {
          if (addon.name === "prometheus" && (!prometheusURL || prometheusURL === "")) {
            setPrometheusURL("http://" + addon.endpoint);
            submitPrometheusConfigure(() => console.log("Prometheus connected"));
          } else if (addon.name === "grafana" && (!grafanaURL || grafanaURL === "")) {
            setGrafanaURL("http://" + addon.endpoint);
            submitGrafanaConfigure(() => {
              // setSelectedBoardsConfigs((prevConfigs) => [...prevConfigs, boardConfigs]);
              console.log("Grafana added");
            });
          }
        });
      },
      error : (err) => console.error("error registering addons:", err),
    });
  };
  const fetchAvailableAdapters = () => {

    props.updateProgress({ showProgress : true });
    dataFetch(
      "/api/system/adapters",
      {
        method : "GET",
        credentials : "include",
      },
      (result) => {
        props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          /* eslint-disable no-unused-vars */
          const options = result.map((res) => ({
            value : res.adapter_location,
            label : res.adapter_location,
          }));




          // this.setState({ availableAdapters : options });
          // setAvailableAdapters(options);
        }
      },
      handleError("Unable to fetch list of adapters.")
    );
  };
  const setMeshScanData = (controlPlanesData, dataPlanesData) => {

    const namespaces = {};
    const activeNamespaces = {};

    const processedControlPlanesData = controlPlanesData?.controlPlanesState?.map((mesh) => {
      if (!mesh?.members?.length) {
        return;
      }
      let proxies = [];

      if (Array.isArray(dataPlanesData?.dataPlanesState)) {
        const dataplane = dataPlanesData.dataPlanesState.find(mesh_ => mesh_.name === mesh.name);

        if (Array.isArray(dataplane?.proxies)) proxies = dataplane.proxies;
      }

      const processedMember = mesh?.members?.map((member) => {
        if (namespaces[mesh.name]) {
          namespaces[mesh.name].add(member.namespace);
        } else {
          namespaces[mesh.name] = new Set([member.namespace]);
        }

        if (proxies.length > 0) {
          const controlPlaneMemberProxies = proxies.filter(proxy => proxy.controlPlaneMemberName === member.name);

          if (controlPlaneMemberProxies.length > 0) {
            member = {
              ...member,
              data_planes : controlPlaneMemberProxies
            };
          }
        }

        return member;
      });

      namespaces[mesh.name] = [...namespaces[mesh.name]];
      activeNamespaces[mesh.name] = namespaces[mesh.name][0] || "";

      return {
        ...mesh,
        members : processedMember
      };
    });

    setMeshScan(processedControlPlanesData?.filter(data => !!data).filter((data) => data.members?.length > 0));
    // setMeshScanNamespaces(namespaces);
    // setActiveMeshScanNamespace(activeNamespaces);
  };
  /**
   * generateMeshScanPodName takes in the podname and the hash
   * and returns the trimmed pod name
   * @param {string} podname
   * @param {string} hash
   * @param {string | undefined} custom
   * @returns {{full, trimmed}}
   */
  // const generateMeshScanPodName = (podname, hash, custom) => {
  //   const str = custom || podname;
  //   return {
  //     full : podname,
  //     trimmed : str.substring(0, (hash ? str.indexOf(hash) : str.length) - 1),
  //   };
  // };
  /**
   * generateMeshScanVersion takes in the string from which version
   * is to be extracted and returns the version. If the version string
   * is undefined then it returns "NA"
   * @param {string | undefined} versionStr is the string from which version is to be extracted
   * @returns {string}
   */
  // const generateMeshScanVersion = (versionStr) => {
  //   if (typeof versionStr !== "string") return "NA";

  //   const matchResult = versionStr.match(/\d+(\.\d+){2,}/g);
  //   if (!matchResult) return "NA";

  //   // Add "v" iff we have a valid match result
  //   return `v${matchResult[0]}`;
  // };
  const handleError = (msg) => (error) => {
    props.updateProgress({ showProgress : false });
    notify({ message : `${msg}: ${error}`, event_type : EVENT_TYPES.ERROR });
  };
  /**
   * redirectErrorToConsole returns a function which redirects
   * ther error to the console under the group labelled by the "msg"
   * param
   * @param {string} msg
   */
  // const redirectErrorToConsole = (msg) => (error) => {
  //   props.updateProgress({ showProgress : false });
  //   console.group(msg);
  //   console.error(error);
  //   console.groupEnd();
  // };
  // const handleAdapterPingError = (msg) => () => {
  //   props.updateProgress({ showProgress : false });


  //   // If you want to use the snackbar code, you might need additional logic or hooks
  //   // based on your functional component setup.

  //   notify({ message : `${msg}`, event_type : EVENT_TYPES.ERROR });
  // };
  // function handleDelete() {
  //   return false;
  // }

  // const handleAdapterClick = (adapterLoc) => () => {
  //   // const { meshAdapters } = this.state;
  //   props.updateProgress({ showProgress : true });

  //   dataFetch(
  //         `/api/system/adapters?adapter=${encodeURIComponent(adapterLoc)}`,
  //         {
  //           credentials : "include",
  //         },
  //         (result) => {
  //           props.updateProgress({ showProgress : false });
  //           if (typeof result !== "undefined") {
  //             notify({ message : `Meshery Adapter connected at ${adapterLoc}`, event_type : EVENT_TYPES.SUCCESS })
  //           }
  //         },
  //         handleAdapterPingError("Could not connect to " + adapterLoc)
  //   );
  // };

  // const handleConfigure = (val) => {
  //   router.push(`/settings#metrics/${val}`);
  // };
  const getSelectedK8sContextsNames = () => {
    return getK8sClusterNamesFromCtxId(props.selectedK8sContexts, props.k8sconfig);
  };
  const emptyStateMessageForServiceMeshesInfo = () => {
    const clusters = getSelectedK8sContextsNames();
    if (clusters.length === 0) {
      return "No Cluster is selected to show the Service Mesh Information";
    }
    if (clusters.includes("all")) {
      return `No service meshes detected in any of the cluster.`;
    }
    return `No service meshes detected in the ${clusters.join(", ")} cluster(s).`;
  };
  const emptyStateMessageForClusterResources = () => {
    const clusters = getSelectedK8sContextsNames();
    if (clusters.length === 0) {
      return "No Cluster is selected to show the discovered resources";
    }
    if (clusters.includes("all")) {
      return `No resources detected in any of the cluster.`;
    }
    return `No resources detected in the ${clusters.join(", ")} cluster(s).`;
  };

  const handleKubernetesClick = (id) => {
    props.updateProgress({ showProgress : true });
    const selectedCtx = props.k8sconfig?.find((ctx) => ctx.id === id);
    if (!selectedCtx) return;

    const { server, name } = selectedCtx;
    dataFetch(
      "/api/system/kubernetes/ping?connection_id=" + id,
      {
        credentials : "include",
      },
      (result) => {
        props.updateProgress({ showProgress : false });
        if (typeof result !== "undefined") {
          notify({ message : `${name} is connected at ${server}`, event_type : EVENT_TYPES.SUCCESS });
        }
      },
      handleError("Could not connect to Kubernetes")
    );
  };
  // const handleGrafanaClick = () => {
  //   props.updateProgress({ showProgress : true });
  //   const grafanaURL = grafana.grafanaURL;  // Assuming grafana is a state or prop

  //   dataFetch(
  //     "/api/telemetry/metrics/grafana/ping",
  //     {
  //       credentials : "include",
  //     },
  //     (result) => {
  //       props.updateProgress({ showProgress : false });
  //       if (typeof result !== "undefined") {
  //         notify({ message : `Grafana connected at ${grafanaURL}`, event_type : EVENT_TYPES.SUCCESS });
  //       }
  //     },
  //     handleError("Could not connect to Grafana")
  //   );
  // };
  /**
   * Meshcard takes in the mesh related data
   * and renders a table along with other information of
   * the mesh
   * @param {{name, icon, tag}} mesh
   * @param {{name, component, version, namespace}[]} components Array of components data
   */
  // const MeshCard = (mesh, components = []) => {
  //   let componentSort = "asc";
  //   let versionSort = "asc";
  //   let proxySort = "asc";
  //   let tempComp = [];
  //   const { theme } = props;

  //   components
  //     .filter(comp => comp.namespace === props.activeMeshScanNamespace[mesh.name])
  //     .map(component => tempComp.push(component));

  //   components = tempComp;

  //   const switchSortOrder = type => {
  //     if (type === "componentSort") {
  //       componentSort = componentSort === "asc" ? "desc" : "asc";
  //       versionSort = "asc";
  //       proxySort = "asc";
  //     } else if (type === "versionSort") {
  //       versionSort = versionSort === "asc" ? "desc" : "asc";
  //       componentSort = "asc";
  //       proxySort = "asc";
  //     } else if (type === "proxySort") {
  //       proxySort = proxySort === "asc" ? "desc" : "asc";
  //       componentSort = "asc";
  //       versionSort = "asc";
  //     }
  //   };

  //   const columns = [
  //     {
  //       name : "name",
  //       label : "Component",
  //       options : {
  //         filter : false,
  //         sort : true,
  //         searchable : true,
  //         setCellProps : () => ({ style : { textAlign : "center" } }),
  //         customHeadRender : ({ index, ...column }, sortColumn) => {
  //           return (
  //             <TableCell key={index} style={{ textAlign : "center" }} onClick={() => {
  //               sortColumn(index); switchSortOrder("componentSort");
  //             }}>
  //               <TableSortLabel active={column.sortDirection != null} direction={componentSort} >
  //                 <b>{column.label}</b>
  //               </TableSortLabel>
  //             </TableCell>

  //           )
  //         }
  //       },
  //       customBodyRender : (value) => {
  //         const modifiedName = value.replace(/-[a-zA-Z0-9]*$/, ''); // Remove last hyphen and alphanumeric characters after it
  //         return <span>{modifiedName}</span>;
  //       },
  //     },
  //     {
  //       name : "version",
  //       label : "Version",
  //       options : {
  //         filter : false,
  //         sort : true,
  //         searchable : true,
  //         setCellProps : () => ({ style : { textAlign : "center" } }),
  //         customHeadRender : ({ index, ...column }, sortColumn) => {
  //           return (
  //             <TableCell key={index} style={{ textAlign : "center" }} onClick={() => {
  //               sortColumn(index); switchSortOrder("versionSort");
  //             }}>
  //               <TableSortLabel active={column.sortDirection != null} direction={versionSort} >
  //                 <b>{column.label}</b>
  //               </TableSortLabel>
  //             </TableCell>

  //           );
  //         },
  //         customBodyRender : (value) => {
  //           return (versionMapper(value))
  //         },
  //       },
  //     },
  //     {
  //       name : "data_planes",
  //       label : "Proxy",
  //       options : {
  //         filter : false,
  //         sort : true,
  //         searchable : true,
  //         setCellProps : () => ({ style : { textAlign : "center" } }),
  //         customHeadRender : ({ index, ...column }, sortColumn) => {
  //           return (
  //             <TableCell key={index} style={{ textAlign : "center" }} onClick={() => {
  //               sortColumn(index); switchSortOrder("proxySort");
  //             }}>
  //               <TableSortLabel active={column.sortDirection != null} direction={proxySort} >
  //                 <b>{column.label}</b>
  //               </TableSortLabel>
  //             </TableCell>
  //           )
  //         },
  //         customBodyRender : (value) => {
  //           return (
  //             <>
  //               <Tooltip
  //                 key={`component-${value}`}
  //                 title={
  //                   Array.isArray(value) && value?.length > 0 ? (
  //                     value.map((cont) => {
  //                       return (
  //                         <div key={cont.name} style={{ fontSize : "15px", color : '#fff', paddingBottom : '10px', padding : '1vh' }}>
  //                           <p>Name: {cont?.containerName ? cont.containerName : 'Unspecified'}</p>
  //                           <p>Status: {cont?.status?.ready ? 'ready' : 'not ready'}</p>
  //                           {!cont?.status?.ready && (
  //                             typeof cont?.status?.lastState === 'object' && cont?.status?.lastState !== null && Object.keys(cont.status.lastState).length > 0 && (
  //                               <div>
  //                                 <p>Last state: {Object.keys(cont?.status?.lastState)[0]} <br /> Error: {Object.values(cont?.status?.lastState)[0]?.exitCode} <br /> Finished at: {Object.values(cont?.status?.lastState)[0]?.finishedAt}</p>
  //                               </div>
  //                             )
  //                           )}
  //                           {typeof cont?.status?.state === 'object' && cont?.status?.state !== null && Object.keys(cont.status.state).length > 0 && (
  //                             <p>State: {Object.keys(cont.status.state)[0]}</p>
  //                           )}
  //                           {cont?.status?.restartCount && (
  //                             <p>Restart count: {cont?.status.restartCount}</p>
  //                           )}
  //                           <p>Image: {cont.image}</p>
  //                           <p>Ports: <br /> {cont?.ports && cont.ports.map(port => `[ ${port?.name ? port.name : 'Unknown'}, ${port?.containerPort ? port.containerPort : 'Unknown'}, ${port?.protocol ? port.protocol : 'Unknown'} ]`).join(', ')}</p>
  //                           {cont?.resources && (
  //                             <div>
  //                                   Resources used: <br />

  //                               <div style={{ paddingLeft : '2vh' }}>
  //                                 {cont?.resources?.limits && (
  //                                   <div>
  //                                     <p>Limits: <br />
  //                                           CPU: {cont?.resources?.limits?.cpu} - Memory: {cont?.resources?.limits?.memory}</p>
  //                                   </div>
  //                                 )}
  //                                 {cont?.resources?.requests && (
  //                                   <div>
  //                                     <p>Requests: <br />
  //                                           CPU: {cont?.resources?.requests?.cpu} - Memory: {cont?.resources?.requests?.memory}</p>
  //                                   </div>
  //                                 )}
  //                               </div>
  //                             </div>
  //                           )}
  //                         </div>
  //                       )
  //                     })
  //                   ) : "No proxy attached"}
  //               >
  //                 <TableCell align="center">{value?.length || 0}</TableCell>
  //               </Tooltip>
  //             </>
  //           );
  //         }
  //       },
  //     },
  //   ]

  //   const options = {
  //     filter : false,
  //     selectableRows : "none",
  //     responsive : "standard",
  //     print : false,
  //     download : false,
  //     viewColumns : false,
  //     pagination : false,
  //     fixedHeader : true,
  //     customToolbar : () => {
  //       return (
  //         <>
  //           {activeMeshScanNamespace[mesh.name] && (
  //             <Select
  //               value={activeMeshScanNamespace[mesh.name]}
  //               onChange={(e) =>
  //               // self.setState((state) => ({ activeMeshScanNamespace : { ...state.activeMeshScanNamespace, [mesh.name] : e.target.value }, }))
  //                 setActiveMeshScanNamespace(prevState => ({
  //                   ...prevState,
  //                   [mesh.name] : e.target.value
  //                 }))
  //               }
  //             >
  //               {self.state.meshScanNamespaces[mesh.name] &&
  //                     self.state.meshScanNamespaces[mesh.name].map((ns) => <MenuItem key={ns.uniqueID} value={ns}>{ns}</MenuItem>)}
  //             </Select>
  //           )}
  //         </>
  //       )
  //     },
  //   }

  //   if (Array.isArray(components) && components.length)
  //     return (
  //       <Paper elevation={1} style={{ padding : "2rem", marginTop : "1rem" }}>
  //         <MuiThemeProvider theme={theme.palette.type === "dark" ? configurationTableThemeDark() : configurationTableTheme()}>
  //           <MUIDataTable
  //             className={classes.datatable}
  //             title={
  //               <>
  //                 <div style={{ display : "flex", alignItems : "center", marginBottom : "1rem" }}>
  //                   <img src={mesh.icon} className={classes.icon} style={{ marginRight : "0.75rem" }} />
  //                   <Typography variant="h6">{mesh.tag}</Typography>
  //                 </div>
  //               </>
  //             }
  //             data={components}
  //             options={options}
  //             columns={columns}
  //           />
  //         </MuiThemeProvider>
  //       </Paper>
  //     );

  //   return null;
  // };
  // const handlePrometheusClick = () => {
  //   props.updateProgress({ showProgress : true });
  //   const prometheusURL = prometheusURL;  // Assuming you have prometheusURL in state

  //   dataFetch(
  //     "/api/telemetry/metrics/ping",
  //     { credentials : "include" },
  //     result => {
  //       props.updateProgress({ showProgress : false });
  //       if (typeof result !== "undefined") {
  //         notify({ message : `Prometheus connected at ${prometheusURL}`, event_type : EVENT_TYPES.SUCCESS });
  //       }
  //     },
  //     handleError("Could not connect to Prometheus")
  //   );
  // };
  // const showCard=(title, content) => {

  //   return (
  //     <Card className={classes.card}>
  //       <CardHeader
  //         disableTypography
  //         title={title}
  //         // action={iconComponent}
  //         className={classes.cardHeader}
  //       />
  //       <CardContent className={classes.cardContent}>{content}</CardContent>
  //     </Card>
  //   );
  // }
  /**
     * ClusterResourcesCard takes in the cluster related data
     * and renders a table with cluster resources information of
     * the selected cluster and namespace
     * @param {{kind, number}[]} resources
     */
  const ClusterResourcesCard = (resources = []) => {

    let kindSort = "asc";
    let countSort = "asc";
    const { theme } = this.props;
    const switchSortOrder = (type) => {
      if (type === "kindSort") {
        kindSort = (kindSort === "asc") ? "desc" : "asc";
        countSort = "asc";
      } else if (type === "countSort") {
        countSort = (countSort === "asc") ? "desc" : "asc";
        kindSort = "asc";
      }
    }

    const columns = [
      {
        name : "kind",
        label : "Resources",
        options : {
          filter : false,
          sort : true,
          searchable : true,
          setCellProps : () => ({ style : { textAlign : "center" } }),
          customHeadRender : ({ index, ...column }, sortColumn) => {
            return (
              <TableCell key={index} style={{ textAlign : "center" }} onClick={() => {
                sortColumn(index); switchSortOrder("kindSort");
              }}>
                <TableSortLabel active={column.sortDirection != null} direction={kindSort} >
                  <b>{column.label}</b>
                </TableSortLabel>
              </TableCell>

            )
          }
        },
      },
      {
        name : "count",
        label : "Count",
        options : {
          filter : false,
          sort : true,
          searchable : true,
          setCellProps : () => ({ style : { textAlign : "center" } }),
          customHeadRender : ({ index, ...column }, sortColumn) => {
            return (
              <TableCell key={index} style={{ textAlign : "center" }} onClick={() => {
                sortColumn(index); switchSortOrder("countSort");
              }}>
                <TableSortLabel active={column.sortDirection != null} direction={countSort} >
                  <b>{column.label}</b>
                </TableSortLabel>
              </TableCell>
            )
          },
        },
      },
    ]

    const options = {
      filter : false,
      selectableRows : "none",
      responsive : "standard",namespaceList,
      print : false,
      download : false,
      viewColumns : false,
      pagination : false,
      fixedHeader : true,
      customToolbar : () => {
        return (
          <>
            {namespaceList && (
              <Select
                value={selectedNamespace}
                onChange={(e) =>
                  // self.setState({ selectedNamespace : e.target.value })
                  setSelectedNamespace(e.target.value)
                }
              >
                {/* {self.state.namespaceList && self.state.namespaceList.map((ns) => <MenuItem key={ns.uniqueID} value={ns}>{ns}</MenuItem>)} */}
                {namespaceList && namespaceList.map((ns) => <MenuItem key={ns.uniqueID} value={ns}>{ns}</MenuItem>)}
              </Select>
            )}
          </>
        )
      }
    }

    if (Array.isArray(resources) && resources.length)
      return (
        <Paper elevation={1} style={{ padding : "2rem" }}>
          <MuiThemeProvider theme={theme.palette.type === "dark" ? configurationTableThemeDark() : configurationTableTheme()}>
            <MUIDataTable
              title={
                <>
                  <div style={{ display : "flex", alignItems : "center", marginBottom : "1rem" }}>
                    <img src={"/static/img/all_mesh.svg"} className={classes.icon} style={{ marginRight : "0.75rem" }} />
                    <Typography variant="h6">All Workloads</Typography>
                  </div>
                </>
              }
              data={resources}
              options={options}
              columns={columns}
            />
          </MuiThemeProvider>
        </Paper>
      );

    return null;
  };
  const configureTemplate = () => {

    let showConfigured = "Not connected to Kubernetes.";
    let chp = (
      <div>
        {props.k8sconfig?.map(ctx => (
          <Tooltip key={ctx.uniqueID} title={`Server: ${ctx.server}`}>
            <Chip
              label={ctx?.name}
              className={classes.chip}
              onClick={() => handleKubernetesClick(ctx.connection_id)}
              icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
              variant="outlined"
              data-cy="chipContextName"
            />
          </Tooltip>
        ))}
      </div>
    );

    if (!props.k8sconfig?.length) {
      chp = showConfigured;
    }

    showConfigured = <div>{chp}</div>;

    const showServiceMesh =(
      <>
        {meshScan && Object.keys(meshScan).length
          ? (
            <>
              {meshScan.map((mesh) => {
                let tag = "";
                mesh.name
                  .split("_")
                  .forEach((element) => {
                    tag = tag + " " + element[0].toUpperCase() + element.slice(1, element.length);
                  });
                return self.Meshcard(
                  { name : mesh.name, tag : tag, icon : "/static/img/" + mesh.name + ".svg" },
                  mesh.members
                );
              })}
            </>
          )
          : (
            <div
              style={{
                padding : "2rem",
                display : "flex",
                justifyContent : "center",
                alignItems : "center",
                flexDirection : "column",
              }}
            >
              <Typography style={{ fontSize : "1.5rem", marginBottom : "2rem" }} align="center" >
                {emptyStateMessageForServiceMeshesInfo()}
              </Typography>
              <Button
                aria-label="Add Meshes"
                variant="contained"
                color="primary"
                size="large"
                onClick={() => router.push("/management")}
              >
                <AddIcon style={iconMedium} className={classes.addIcon} />
                    Install Service Mesh
              </Button>
            </div>
          )}
      </>
    );
    const showClusterResources = (
      <>
        {clusterResources && Object.keys(clusterResources) && clusterResources?.resources?.length > 0
          ? (
            ClusterResourcesCard(clusterResources?.resources)
          )
          : (
            <div
              style={{
                padding : "2rem",
                display : "flex",
                justifyContent : "center",
                alignItems : "center",
                flexDirection : "column",
              }}
            >
              <Typography style={{ fontSize : "1.5rem", marginBottom : "2rem" }} align="center" >
                {emptyStateMessageForClusterResources()}
              </Typography>
              <Button
                aria-label="Connect K8s cluster"
                variant="contained"
                color="primary"
                size="large"
                onClick={() => router.push("/settings")}
              >
                <AddIcon style={iconMedium} className={classes.addIcon} />
                    Connect Cluster
              </Button>
            </div>
          )}
      </>
    );

    return (
      <NoSsr>
        <Popup />
        <div className={classes.rootClass}>
          <Grid container spacing={2}>
            <Grid item xs={12} md={12}>
              <DashboardMeshModelGraph classes={classes} />
            </Grid>
            <Grid item xs={12} md={12}>
              <div className={classes.dashboardSection} data-test="workloads">
                <Typography variant="h6" gutterBottom className={classes.chartTitle}>
                      Workloads
                </Typography>
                {showClusterResources}
              </div>
            </Grid>

            <Grid item xs={12}>
              <Grid container spacing={2}>
                <Grid item xs={12} md={4}>
                  <ConnectionStatsChart classes={classes} />
                </Grid>
                <Grid item xs={12} md={8}>
                  <div className={classes.dashboardSection} data-test="service-mesh">
                    <Typography variant="h6" gutterBottom className={classes.chartTitle}>
                          Service Mesh
                    </Typography>
                    {showServiceMesh}
                  </div>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </div>
      </NoSsr>
    );
  };

  return configureTemplate();

}
DashboardComponent.propTypes = {
  classes : PropTypes.object.isRequired,
};
const mapDispatchToProps = (dispatch) => ({
  updateProgress : bindActionCreators(updateProgress, dispatch),
  updateGrafanaConfig : bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig : bindActionCreators(updatePrometheusConfig, dispatch),
  updateTelemetryUrls : bindActionCreators(updateTelemetryUrls, dispatch),
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig");
  const meshAdapters = state.get("meshAdapters");
  const meshAdaptersts = state.get("meshAdaptersts");
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    meshAdapters,
    meshAdaptersts,
    k8sconfig,
    grafana,
    prometheus,
    selectedK8sContexts
  };
};

export default connect(mapStateToProps, mapDispatchToProps)((DashboardComponent));
