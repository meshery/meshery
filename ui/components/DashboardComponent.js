import {
  Button,
  Card,
  CardContent,
  CardHeader,
  Chip,
  MenuItem,
  NoSsr,
  Paper,
  Select,
  TableCell,
  TableSortLabel,
  Tooltip,
  Typography,
} from '@material-ui/core';
import blue from '@material-ui/core/colors/blue';
import Grid from '@material-ui/core/Grid';
import { withStyles, MuiThemeProvider } from '@material-ui/core/styles';
import AddIcon from '@material-ui/icons/AddCircleOutline';
import { withRouter } from 'next/router';
import PropTypes from 'prop-types';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import dataFetch from '../lib/data-fetch';
import {
  updateGrafanaConfig,
  updateProgress,
  updatePrometheusConfig,
  updateTelemetryUrls,
} from '../lib/store';
import { getK8sClusterIdsFromCtxId, getK8sClusterNamesFromCtxId } from '../utils/multi-ctx';
import { versionMapper } from '../utils/nameMapper';
import fetchControlPlanes from './graphql/queries/ControlPlanesQuery';
import fetchDataPlanes from './graphql/queries/DataPlanesQuery';

import subscribeClusterResources from './graphql/subscriptions/ClusterResourcesSubscription';
import fetchAvailableNamespaces from './graphql/queries/NamespaceQuery';

import MUIDataTable from 'mui-datatables';
import Popup from './Popup';
import { iconMedium } from '../css/icons.styles';
import {
  configurationTableTheme,
  configurationTableThemeDark,
} from '../themes/configurationTableTheme';
import DashboardMeshModelGraph from './Dashboard/DashboardMeshModelGraph';
import ConnectionStatsChart from './Dashboard/ConnectionCharts.js';
import { EVENT_TYPES } from '../lib/event-types';
import { withNotify } from '../utils/hooks/useNotification';
import { useTheme } from '@emotion/react';

const styles = (theme) => ({
  rootClass: { backgroundColor: theme.palette.secondary.elevatedComponents2 },
  datatable: {
    boxShadow: 'none',
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  link: {
    cursor: 'pointer',
    textDecoration: 'none',
  },
  metricsButton: { width: '240px' },
  alreadyConfigured: { textAlign: 'center' },
  margin: { margin: theme.spacing(1) },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': { backgroundColor: blue[500] },
    },
  },
  colorBar: {},
  colorChecked: {},
  fileLabel: { width: '100%' },
  fileLabelText: {},
  inClusterLabel: { paddingRight: theme.spacing(2) },
  alignCenter: { textAlign: 'center' },
  icon: { width: theme.spacing(2.5) },
  istioIcon: { width: theme.spacing(1.5) },
  settingsIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  addIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  cardHeader: { fontSize: theme.spacing(2) },
  card: {
    height: '100%',
    marginTop: theme.spacing(2),
  },
  cardContent: { height: '100%' },
  redirectButton: {
    marginLeft: '-.5em',
    color: '#000',
  },
  dashboardSection: {
    backgroundColor: theme.palette.secondary.elevatedComponents,
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
    marginBottom: theme.spacing(2),
  },
});
class DashboardComponent extends React.Component {
  constructor(props) {
    super(props);
    const { meshAdapters, grafana, prometheus } = props;
    this._isMounted = false;
    this.state = {
      meshAdapters,
      contextsFromFile: [],
      availableAdapters: [],
      mts: new Date(),
      meshLocationURLError: false,

      grafanaUrl: grafana.grafanaURL,
      prometheusUrl: prometheus.prometheusURL,
      k8sfileError: false,
      kts: new Date(),

      grafana,
      prometheus,

      urlError: false,
      grafanaConfigSuccess: props.grafana.grafanaURL !== '',
      grafanaBoardSearch: '',
      grafanaURL: props.grafana.grafanaURL,
      grafanaAPIKey: props.grafana.grafanaAPIKey,
      grafanaBoards: props.grafana.grafanaBoards,
      selectedBoardsConfigs: props.grafana.selectedBoardsConfigs,
      ts: props.grafana.ts,

      meshScan: [],
      activeMeshScanNamespace: {},
      meshScanNamespaces: {},

      isMetricsConfigured: grafana.grafanaURL !== '' && prometheus.prometheusURL !== '',
      controlPlaneState: '',
      dataPlaneState: '',
      clusterResources: [],
      namespaceList: [],
      selectedNamespace: 'default',

      // subscriptions disposable
      dataPlaneSubscription: null,
      controlPlaneSubscription: null,
      clusterResourcesSubscription: null,
      clusterResourcesQuery: null,
      namespaceQuery: null,
      telemetryQuery: null,
    };
  }

  static getDerivedStateFromProps(props, state) {
    const { meshAdapters, meshAdaptersts, grafana, prometheus } = props;
    const st = {};
    if (meshAdaptersts > state.mts) {
      st.meshAdapters = meshAdapters;
      st.mts = meshAdaptersts;
    }
    st.grafana = grafana;
    st.prometheus = prometheus;
    st.k8sconfig = props.k8sconfig;
    return st;
  }

  disposeWorkloadWidgetSubscription = () => {
    this.state.namespaceQuery && this.state.namespaceQuery.unsubscribe();
    this.state.clusterResourcesSubscription && this.state.clusterResourcesSubscription.dispose();
  };

  disposeSubscriptions = () => {
    if (this.state.dataPlaneSubscription) {
      this.state.dataPlaneSubscription.unsubscribe();
    }
    if (this.state.controlPlaneSubscription) {
      this.state.controlPlaneSubscription.unsubscribe();
    }
    this.disposeWorkloadWidgetSubscription();
  };

  initMeshSyncControlPlaneSubscription = () => {
    /**
     * ALL_MESH indicates that we are interested in control plane
     * component of all of the service meshes supported by meshsync v2
     */
    const self = this;
    const ALL_MESH = { type: 'ALL_MESH', k8sClusterIDs: self.getK8sClusterIds() };

    if (self._isMounted) {
      const controlPlaneSubscription = fetchControlPlanes(ALL_MESH).subscribe({
        next: (controlPlaneRes) => {
          this.setState({ controlPlaneState: controlPlaneRes });
        },
        error: (err) => console.error(err),
      });

      const dataPlaneSubscription = fetchDataPlanes(ALL_MESH).subscribe({
        next: (dataPlaneRes) => {
          this.setState({ dataPlaneState: dataPlaneRes });
        },
        error: (err) => console.error(err),
      });

      this.setState({ controlPlaneSubscription, dataPlaneSubscription });
    }
  };

  initNamespaceQuery = () => {
    const self = this;

    const namespaceQuery = fetchAvailableNamespaces({
      k8sClusterIDs: self.getK8sClusterIds(),
    }).subscribe({
      next: (res) => {
        let namespaces = [];
        res?.namespaces?.map((ns) => {
          namespaces.push(ns?.namespace);
        });
        namespaces.sort((a, b) => (a > b ? 1 : -1));
        self.setState({ namespaceList: namespaces });
      },
      error: (err) => console.log('error at namespace fetch: ' + err),
    });

    this.setState({ namespaceQuery });
  };

  initDashboardClusterResourcesSubscription = () => {
    const self = this;
    let k8s = self.getK8sClusterIds();
    let namespace = self.state.selectedNamespace;

    if (self._isMounted) {
      // @ts-ignore
      const clusterResourcesSubscription = subscribeClusterResources(
        (res) => {
          this.setState({ clusterResources: res?.clusterResources });
        },
        {
          k8scontextIDs: k8s,
          namespace: namespace,
        },
      );
      this.setState({ clusterResourcesSubscription });
    }
  };

  componentWillUnmount = () => {
    this._isMounted = false;
    this.disposeSubscriptions();
  };

  componentDidMount = () => {
    this._isMounted = true;
    if (this._isMounted) {
      this.initMeshSyncControlPlaneSubscription();
      this.initDashboardClusterResourcesSubscription();
      this.initNamespaceQuery();
    }
  };

  componentDidUpdate(prevProps, prevState) {
    let updateControlPlane = false;
    let updateDataPlane = false;

    // deep compare very limited, order of object fields is important
    if (
      JSON.stringify(prevState.controlPlaneState) !== JSON.stringify(this.state.controlPlaneState)
    ) {
      updateControlPlane = true;
    }
    if (JSON.stringify(prevState.dataPlaneState) !== JSON.stringify(this.state.dataPlaneState)) {
      updateDataPlane = true;
    }

    if (updateDataPlane || updateControlPlane) {
      this.setMeshScanData(
        updateControlPlane ? this.state.controlPlaneState : prevState.controlPlaneState,
        updateDataPlane ? this.state.dataPlaneState : prevState.dataPlaneState,
      );
    }
    // handle subscriptions update on switching K8s Contexts
    if (
      prevProps?.selectedK8sContexts !== this.props?.selectedK8sContexts ||
      prevProps.k8sconfig !== this.props.k8sconfig
    ) {
      this.disposeSubscriptions();
      this.initMeshSyncControlPlaneSubscription();
      this.initDashboardClusterResourcesSubscription();
      this.initNamespaceQuery();
    }

    if (prevState?.selectedNamespace !== this.state?.selectedNamespace) {
      this.disposeWorkloadWidgetSubscription();
      this.initDashboardClusterResourcesSubscription();
      this.initNamespaceQuery();
    }
  }

  getK8sClusterIds = () => {
    const self = this;
    return getK8sClusterIdsFromCtxId(self.props?.selectedK8sContexts, self.props.k8sconfig);
  };

  setMeshScanData = (controlPlanesData, dataPlanesData) => {
    const self = this;
    const namespaces = {};
    const activeNamespaces = {};
    const processedControlPlanesData = controlPlanesData?.controlPlanesState?.map((mesh) => {
      if (!mesh?.members?.length) {
        return;
      }
      let proxies = [];

      if (Array.isArray(dataPlanesData?.dataPlanesState)) {
        const dataplane = dataPlanesData.dataPlanesState.find((mesh_) => mesh_.name === mesh.name);

        if (Array.isArray(dataplane?.proxies)) proxies = dataplane.proxies;
      }
      const processedMember = mesh?.members?.map((member) => {
        if (namespaces[mesh.name]) {
          namespaces[mesh.name].add(member.namespace);
        } else {
          namespaces[mesh.name] = new Set([member.namespace]);
        }

        // retrieve data planes according to mesh name
        if (proxies.length > 0) {
          const controlPlaneMemberProxies = proxies.filter(
            (proxy) => proxy.controlPlaneMemberName === member.name,
          );

          if (controlPlaneMemberProxies.length > 0) {
            member = {
              ...member,
              data_planes: controlPlaneMemberProxies,
            };
          }
        }

        return member;
      });
      namespaces[mesh.name] = [...namespaces[mesh.name]];
      activeNamespaces[mesh.name] = namespaces[mesh.name][0] || '';

      return {
        ...mesh,
        members: processedMember,
      };
    });
    self.setState({
      meshScan: processedControlPlanesData
        ?.filter((data) => !!data)
        .filter((data) => data.members?.length > 0),
    });
    self.setState({ meshScanNamespaces: namespaces, activeMeshScanNamespace: activeNamespaces });
  };

  /**
   * generateMeshScanPodName takes in the podname and the hash
   * and returns the trimmed pod name
   * @param {string} podname
   * @param {string} hash
   * @param {string | undefined} custom
   * @returns {{full, trimmed}}
   */
  generateMeshScanPodName = (podname, hash, custom) => {
    const str = custom || podname;
    return {
      full: podname,
      trimmed: str.substring(0, (hash ? str.indexOf(hash) : str.length) - 1),
    };
  };

  /**
   * generateMeshScanVersion takes in the string from which version
   * is to be extracted and returns the version. If the version string
   * is undefined then it returns "NA"
   * @param {string | undefined} versionStr is the string from which version is to be extracted
   * @returns {string}
   */
  generateMeshScanVersion = (versionStr) => {
    if (typeof versionStr !== 'string') return 'NA';

    const matchResult = versionStr.match(/\d+(\.\d+){2,}/g);
    if (!matchResult) return 'NA';

    // Add "v" iff we have a valid match result
    return `v${matchResult[0]}`;
  };

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress: false });
    const notify = this.props.notify;
    notify({ message: `${msg}: ${error}`, event_type: EVENT_TYPES.ERROR });
  };

  /**
   * redirectErrorToConsole returns a function which redirects
   * ther error to the console under the group labelled by the "msg"
   * param
   * @param {string} msg
   */
  redirectErrorToConsole = (msg) => (error) => {
    this.props.updateProgress({ showProgress: false });
    console.group(msg);
    console.error(error);
    console.groupEnd();
  };

  handleDelete() {
    return false;
  }

  handleConfigure = (val) => {
    this.props.router.push(`/settings#metrics/${val}`);
  };

  getSelectedK8sContextsNames = () => {
    return getK8sClusterNamesFromCtxId(this.props.selectedK8sContexts, this.props.k8sconfig);
  };

  emptyStateMessageForServiceMeshesInfo = () => {
    const clusters = this.getSelectedK8sContextsNames();
    if (clusters.length === 0) {
      return 'No Cluster is selected to show the Service Mesh Information';
    }
    if (clusters.includes('all')) {
      return `No service meshes detected in any of the cluster.`;
    }
    return `No service meshes detected in the ${clusters.join(', ')} cluster(s).`;
  };

  emptyStateMessageForClusterResources = () => {
    const clusters = this.getSelectedK8sContextsNames();
    if (clusters.length === 0) {
      return 'No Cluster is selected to show the discovered resources';
    }
    if (clusters.includes('all')) {
      return `No resources detected in any of the cluster.`;
    }
    return `No resources detected in the ${clusters.join(', ')} cluster(s).`;
  };

  handleKubernetesClick = (id) => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    const notify = this.props.notify;
    const selectedCtx = this.props.k8sconfig?.find((ctx) => ctx.id === id);
    if (!selectedCtx) return;

    const { server, name } = selectedCtx;
    dataFetch(
      '/api/system/kubernetes/ping?connection_id=' + id,
      {
        credentials: 'include',
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          notify({ message: `${name} is connected at ${server}`, event_type: EVENT_TYPES.SUCCESS });
        }
      },
      self.handleError('Could not connect to Kubernetes'),
    );
  };

  /**
   * Meshcard takes in the mesh related data
   * and renders a table along with other information of
   * the mesh
   * @param {{name, icon, tag}} mesh
   * @param {{name, component, version, namespace}[]} components Array of components data
   */
  Meshcard = (mesh, components = []) => {
    const self = this;
    let componentSort = 'asc';
    let versionSort = 'asc';
    let proxySort = 'asc';
    let tempComp = [];
    const { theme } = this.props;
    components
      .filter((comp) => comp.namespace === self.state.activeMeshScanNamespace[mesh.name])
      .map((component) => tempComp.push(component));

    components = tempComp;

    const switchSortOrder = (type) => {
      if (type === 'componentSort') {
        componentSort = componentSort === 'asc' ? 'desc' : 'asc';
        versionSort = 'asc';
        proxySort = 'asc';
      } else if (type === 'versionSort') {
        versionSort = versionSort === 'asc' ? 'desc' : 'asc';
        componentSort = 'asc';
        proxySort = 'asc';
      } else if (type === 'proxySort') {
        proxySort = proxySort === 'asc' ? 'desc' : 'asc';
        componentSort = 'asc';
        versionSort = 'asc';
      }
    };

    const columns = [
      {
        name: 'name',
        label: 'Component',
        options: {
          filter: false,
          sort: true,
          searchable: true,
          setCellProps: () => ({ style: { textAlign: 'center' } }),
          customHeadRender: ({ index, ...column }, sortColumn) => {
            return (
              <TableCell
                key={index}
                style={{ textAlign: 'center' }}
                onClick={() => {
                  sortColumn(index);
                  switchSortOrder('componentSort');
                }}
              >
                <TableSortLabel active={column.sortDirection != null} direction={componentSort}>
                  <b>{column.label}</b>
                </TableSortLabel>
              </TableCell>
            );
          },
        },
        customBodyRender: (value) => {
          const modifiedName = value.replace(/-[a-zA-Z0-9]*$/, ''); // Remove last hyphen and alphanumeric characters after it
          return <span>{modifiedName}</span>;
        },
      },
      {
        name: 'version',
        label: 'Version',
        options: {
          filter: false,
          sort: true,
          searchable: true,
          setCellProps: () => ({ style: { textAlign: 'center' } }),
          customHeadRender: ({ index, ...column }, sortColumn) => {
            return (
              <TableCell
                key={index}
                style={{ textAlign: 'center' }}
                onClick={() => {
                  sortColumn(index);
                  switchSortOrder('versionSort');
                }}
              >
                <TableSortLabel active={column.sortDirection != null} direction={versionSort}>
                  <b>{column.label}</b>
                </TableSortLabel>
              </TableCell>
            );
          },
          customBodyRender: (value) => {
            return versionMapper(value);
          },
        },
      },
      {
        name: 'data_planes',
        label: 'Proxy',
        options: {
          filter: false,
          sort: true,
          searchable: true,
          setCellProps: () => ({ style: { textAlign: 'center' } }),
          customHeadRender: ({ index, ...column }, sortColumn) => {
            return (
              <TableCell
                key={index}
                style={{ textAlign: 'center' }}
                onClick={() => {
                  sortColumn(index);
                  switchSortOrder('proxySort');
                }}
              >
                <TableSortLabel active={column.sortDirection != null} direction={proxySort}>
                  <b>{column.label}</b>
                </TableSortLabel>
              </TableCell>
            );
          },
          customBodyRender: (value) => {
            return (
              <>
                <Tooltip
                  key={`component-${value}`}
                  title={
                    Array.isArray(value) && value?.length > 0
                      ? value.map((cont) => {
                          return (
                            <div
                              key={cont.name}
                              style={{
                                fontSize: '15px',
                                color: '#fff',
                                paddingBottom: '10px',
                                padding: '1vh',
                              }}
                            >
                              <p>
                                Name: {cont?.containerName ? cont.containerName : 'Unspecified'}
                              </p>
                              <p>Status: {cont?.status?.ready ? 'ready' : 'not ready'}</p>
                              {!cont?.status?.ready &&
                                typeof cont?.status?.lastState === 'object' &&
                                cont?.status?.lastState !== null &&
                                Object.keys(cont.status.lastState).length > 0 && (
                                  <div>
                                    <p>
                                      Last state: {Object.keys(cont?.status?.lastState)[0]} <br />{' '}
                                      Error: {Object.values(cont?.status?.lastState)[0]?.exitCode}{' '}
                                      <br /> Finished at:{' '}
                                      {Object.values(cont?.status?.lastState)[0]?.finishedAt}
                                    </p>
                                  </div>
                                )}
                              {typeof cont?.status?.state === 'object' &&
                                cont?.status?.state !== null &&
                                Object.keys(cont.status.state).length > 0 && (
                                  <p>State: {Object.keys(cont.status.state)[0]}</p>
                                )}
                              {cont?.status?.restartCount && (
                                <p>Restart count: {cont?.status.restartCount}</p>
                              )}
                              <p>Image: {cont.image}</p>
                              <p>
                                Ports: <br />{' '}
                                {cont?.ports &&
                                  cont.ports
                                    .map(
                                      (port) =>
                                        `[ ${port?.name ? port.name : 'Unknown'}, ${
                                          port?.containerPort ? port.containerPort : 'Unknown'
                                        }, ${port?.protocol ? port.protocol : 'Unknown'} ]`,
                                    )
                                    .join(', ')}
                              </p>
                              {cont?.resources && (
                                <div>
                                  Resources used: <br />
                                  <div style={{ paddingLeft: '2vh' }}>
                                    {cont?.resources?.limits && (
                                      <div>
                                        <p>
                                          Limits: <br />
                                          CPU: {cont?.resources?.limits?.cpu} - Memory:{' '}
                                          {cont?.resources?.limits?.memory}
                                        </p>
                                      </div>
                                    )}
                                    {cont?.resources?.requests && (
                                      <div>
                                        <p>
                                          Requests: <br />
                                          CPU: {cont?.resources?.requests?.cpu} - Memory:{' '}
                                          {cont?.resources?.requests?.memory}
                                        </p>
                                      </div>
                                    )}
                                  </div>
                                </div>
                              )}
                            </div>
                          );
                        })
                      : 'No proxy attached'
                  }
                >
                  <TableCell align="center">{value?.length || 0}</TableCell>
                </Tooltip>
              </>
            );
          },
        },
      },
    ];

    const options = {
      filter: false,
      selectableRows: 'none',
      responsive: 'standard',
      print: false,
      download: false,
      viewColumns: false,
      pagination: false,
      fixedHeader: true,
      customToolbar: () => {
        return (
          <>
            {self.state.activeMeshScanNamespace[mesh.name] && (
              <Select
                value={self.state.activeMeshScanNamespace[mesh.name]}
                onChange={(e) =>
                  self.setState((state) => ({
                    activeMeshScanNamespace: {
                      ...state.activeMeshScanNamespace,
                      [mesh.name]: e.target.value,
                    },
                  }))
                }
              >
                {self.state.meshScanNamespaces[mesh.name] &&
                  self.state.meshScanNamespaces[mesh.name].map((ns) => (
                    <MenuItem key={ns.uniqueID} value={ns}>
                      {ns}
                    </MenuItem>
                  ))}
              </Select>
            )}
          </>
        );
      },
    };

    if (Array.isArray(components) && components.length)
      return (
        <Paper elevation={1} style={{ padding: '2rem', marginTop: '1rem' }}>
          <MuiThemeProvider
            theme={
              theme.palette.type == 'dark'
                ? configurationTableThemeDark()
                : configurationTableTheme()
            }
          >
            <MUIDataTable
              className={this.props.classes.datatable}
              title={
                <>
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
                    <img
                      src={mesh.icon}
                      className={this.props.classes.icon}
                      style={{ marginRight: '0.75rem' }}
                    />
                    <Typography variant="h6">{mesh.tag}</Typography>
                  </div>
                </>
              }
              data={components}
              options={options}
              columns={columns}
            />
          </MuiThemeProvider>
        </Paper>
      );

    return null;
  };

  /**
   * ClusterResourcesCard takes in the cluster related data
   * and renders a table with cluster resources information of
   * the selected cluster and namespace
   * @param {{kind, number}[]} resources
   */
  ClusterResourcesCard = (resources = []) => {
    const self = this;
    let kindSort = 'asc';
    let countSort = 'asc';
    const theme = useTheme();
    const switchSortOrder = (type) => {
      if (type === 'kindSort') {
        kindSort = kindSort === 'asc' ? 'desc' : 'asc';
        countSort = 'asc';
      } else if (type === 'countSort') {
        countSort = countSort === 'asc' ? 'desc' : 'asc';
        kindSort = 'asc';
      }
    };

    const columns = [
      {
        name: 'kind',
        label: 'Resources',
        options: {
          filter: false,
          sort: true,
          searchable: true,
          setCellProps: () => ({ style: { textAlign: 'center' } }),
          customHeadRender: ({ index, ...column }, sortColumn) => {
            return (
              <TableCell
                key={index}
                style={{ textAlign: 'center' }}
                onClick={() => {
                  sortColumn(index);
                  switchSortOrder('kindSort');
                }}
              >
                <TableSortLabel active={column.sortDirection != null} direction={kindSort}>
                  <b>{column.label}</b>
                </TableSortLabel>
              </TableCell>
            );
          },
        },
      },
      {
        name: 'count',
        label: 'Count',
        options: {
          filter: false,
          sort: true,
          searchable: true,
          setCellProps: () => ({ style: { textAlign: 'center' } }),
          customHeadRender: ({ index, ...column }, sortColumn) => {
            return (
              <TableCell
                key={index}
                style={{ textAlign: 'center' }}
                onClick={() => {
                  sortColumn(index);
                  switchSortOrder('countSort');
                }}
              >
                <TableSortLabel active={column.sortDirection != null} direction={countSort}>
                  <b>{column.label}</b>
                </TableSortLabel>
              </TableCell>
            );
          },
        },
      },
    ];

    const options = {
      filter: false,
      selectableRows: 'none',
      responsive: 'standard',
      print: false,
      download: false,
      viewColumns: false,
      pagination: false,
      fixedHeader: true,
      customToolbar: () => {
        return (
          <>
            {self.state.namespaceList && (
              <Select
                value={self.state.selectedNamespace}
                onChange={(e) => self.setState({ selectedNamespace: e.target.value })}
              >
                {self.state.namespaceList &&
                  self.state.namespaceList.map((ns) => (
                    <MenuItem key={ns.uniqueID} value={ns}>
                      {ns}
                    </MenuItem>
                  ))}
              </Select>
            )}
          </>
        );
      },
    };

    if (Array.isArray(resources) && resources.length)
      return (
        <Paper elevation={1} style={{ padding: '2rem' }}>
          <MuiThemeProvider
            theme={
              theme.palette.type === 'dark'
                ? configurationTableThemeDark()
                : configurationTableTheme()
            }
          >
            <MUIDataTable
              title={
                <>
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
                    <img
                      src={'/static/img/all_mesh.svg'}
                      className={this.props.classes.icon}
                      style={{ marginRight: '0.75rem' }}
                    />
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

  handlePrometheusClick = () => {
    this.props.updateProgress({ showProgress: true });
    const self = this;
    const notify = this.props.notify;
    const { prometheusURL } = this.state.prometheus;
    dataFetch(
      '/api/telemetry/metrics/ping',
      {
        credentials: 'include',
      },
      (result) => {
        this.props.updateProgress({ showProgress: false });
        if (typeof result !== 'undefined') {
          notify({
            message: `Prometheus connected at ${prometheusURL}`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      },
      self.handleError('Could not connect to Prometheus'),
    );
  };

  showCard(title, content) {
    const { classes } = this.props;
    return (
      <Card className={classes.card}>
        <CardHeader
          disableTypography
          title={title}
          // action={iconComponent}
          className={classes.cardHeader}
        />
        <CardContent className={classes.cardContent}>{content}</CardContent>
      </Card>
    );
  }

  configureTemplate = () => {
    const { classes, k8sconfig } = this.props;
    const self = this;
    let showConfigured = 'Not connected to Kubernetes.';
    let chp = (
      <div>
        {k8sconfig?.map((ctx) => (
          <Tooltip key={ctx.uniqueID} title={`Server: ${ctx.server}`}>
            <Chip
              label={ctx?.name}
              className={classes.chip}
              onClick={() => self.handleKubernetesClick(ctx.connection_id)}
              icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
              variant="outlined"
              data-cy="chipContextName"
            />
          </Tooltip>
        ))}
      </div>
    );

    if (!k8sconfig?.length) {
      chp = showConfigured;
    }

    showConfigured = <div>{chp}</div>;

    const showServiceMesh = (
      <>
        {self?.state?.meshScan && Object.keys(self?.state?.meshScan).length ? (
          <>
            {self.state.meshScan.map((mesh) => {
              let tag = '';
              mesh.name.split('_').forEach((element) => {
                tag = tag + ' ' + element[0].toUpperCase() + element.slice(1, element.length);
              });
              return self.Meshcard(
                { name: mesh.name, tag: tag, icon: '/static/img/' + mesh.name + '.svg' },
                mesh.members,
              );
            })}
          </>
        ) : (
          <div
            style={{
              padding: '2rem',
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              flexDirection: 'column',
            }}
          >
            <Typography style={{ fontSize: '1.5rem', marginBottom: '2rem' }} align="center">
              {this.emptyStateMessageForServiceMeshesInfo()}
            </Typography>
            <Button
              aria-label="Add Meshes"
              variant="contained"
              color="primary"
              size="large"
              onClick={() => self.props.router.push('/management')}
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
        {self?.state?.clusterResources &&
        Object.keys(self?.state?.clusterResources) &&
        self?.state?.clusterResources?.resources?.length > 0 ? (
          self.ClusterResourcesCard(self?.state?.clusterResources?.resources)
        ) : (
          <div
            style={{
              padding: '2rem',
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
              flexDirection: 'column',
            }}
          >
            <Typography style={{ fontSize: '1.5rem', marginBottom: '2rem' }} align="center">
              {this.emptyStateMessageForClusterResources()}
            </Typography>
            <Button
              aria-label="Connect K8s cluster"
              variant="contained"
              color="primary"
              size="large"
              onClick={() => self.props.router.push('/settings')}
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

  render() {
    // console.info("Rerendering Dashboard")
    return this.configureTemplate();
  }
}

DashboardComponent.propTypes = { classes: PropTypes.object.isRequired };

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  updatePrometheusConfig: bindActionCreators(updatePrometheusConfig, dispatch),
  updateTelemetryUrls: bindActionCreators(updateTelemetryUrls, dispatch),
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const meshAdapters = state.get('meshAdapters');
  const meshAdaptersts = state.get('meshAdaptersts');
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    meshAdapters,
    meshAdaptersts,
    k8sconfig,
    grafana,
    prometheus,
    selectedK8sContexts,
  };
};

export default withStyles(styles, { withTheme: true })(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withNotify(DashboardComponent))),
);
