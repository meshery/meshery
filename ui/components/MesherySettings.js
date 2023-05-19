import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { withRouter } from 'next/router';
import { connect } from 'react-redux';
import { bindActionCreators } from "redux";
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import {
  AppBar, Paper, Tooltip, Button, IconButton, MenuItem, Select, TableCell, TableSortLabel, Typography
} from '@material-ui/core';
import CloseIcon from "@material-ui/icons/Close";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faCloud, faPoll, faDatabase, faFileInvoice } from '@fortawesome/free-solid-svg-icons';
// import {faTachometerAlt} from '@fortawesome/free-solid-svg-icons';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import Link from 'next/link';
import MeshConfigComponent from './MeshConfigComponent';
import GrafanaComponent from './telemetry/grafana/GrafanaComponent';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import PrometheusComponent from './telemetry/prometheus/PrometheusComponent';
// import MesherySettingsPerformanceComponent from "../components/MesherySettingsPerformanceComponent";
import { updateProgress } from "../lib/store";
import { withSnackbar } from "notistack";
import PromptComponent from './PromptComponent';
import resetDatabase from './graphql/queries/ResetDatabaseQuery';
import { iconMedium } from '../css/icons.styles';
import subscribeMeshModelSummary from "./graphql/subscriptions/MeshModelSummarySubscription";
import fetchMeshModelSummary from "./graphql/queries/MeshModelSummaryQuery";
import { MuiThemeProvider } from '@material-ui/core/styles';
import MeshModelComponent from './MeshModelComponent';
import DataTable from "mui-datatables";
import { configurationTableTheme, configurationTableThemeDark } from '../themes/configurationTableTheme';


const styles = (theme) => ({
  wrapperClss : {
    flexGrow : 1,
    maxWidth : '100%',
    height : 'auto',
  },
  tab : {
    minWidth : 40,
    paddingLeft : 0,
    paddingRight : 0,
    "&.Mui-selected" : {
      color : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  tabs : {
    "& .MuiTabs-indicator" : {
      backgroundColor : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  icon : {
    display : 'inline',
    verticalAlign : 'text-top',
    width : theme.spacing(1.75),
    marginLeft : theme.spacing(0.5),
  },

  iconText : {
    display : 'inline',
    verticalAlign : 'middle',
  },
  backToPlay : { margin : theme.spacing(2), },
  link : { cursor : 'pointer', },
  DBBtn : {
    margin : theme.spacing(0.5),
    padding : theme.spacing(1),
    borderRadius : 5,
    backgroundColor : "#8F1F00",
    "&:hover" : {
      backgroundColor : "#B32700",
    },
  },
  container : {
    display : "flex",
    justifyContent : "center",
    margin : theme.spacing(2),
    // backgroundColor : "#fff",
  },
  paper : {
    maxWidth : '90%',
    margin : 'auto',
    overflow : 'hidden',
  },
  topToolbar : {
    marginBottom : "2rem",
    display : "flex",
    justifyContent : "space-between",
    paddingLeft : "1rem",
    maxWidth : "90%"
  },
  dashboardSection : {
    backgroundColor : theme.palette.secondary.elevatedComponents,
    padding : theme.spacing(2),
    borderRadius : 4,
    height : "100%",
  },
  cardHeader : { fontSize : theme.spacing(2), },
  card : {
    height : "100%",
    marginTop : theme.spacing(2),
  },
  cardContent : { height : "100%", },
});

function TabContainer(props) {
  return (
    <Typography
      component="div"
      style={{ paddingTop : 2, }}
    >
      {props.children}
    </Typography>
  );
}

TabContainer.propTypes = { children : PropTypes.node.isRequired, };

class MesherySettings extends React.Component {
  constructor(props) {
    super(props);
    const {
      k8sconfig, meshAdapters, grafana, prometheus, router : { asPath }
    } = props;

    this._isMounted = false;
    let tabVal = 0, subTabVal = 0;
    const splittedPath = asPath.split('#');
    if (splittedPath.length >= 2 && splittedPath[1]) {
      const subTabPath = splittedPath[1].split('/');
      switch (subTabPath[0]) {
        case 'environment':
          tabVal = 0;
          break;
        case 'service-mesh':
          tabVal = 1;
          break;
        case 'metrics':
          tabVal = 2;
          break;
        case 'system':
          tabVal = 3;
        // case 'performance':
        //   tabVal = 3;
        //   break;
      }
      if (subTabPath.length >= 2 && subTabPath[1]) {
        switch (subTabPath[1]) {
          case 'inclusterconfig':
            subTabVal = 0;
            break;
          case 'outclusterconfig':
            subTabVal = 1;
            break;
          case 'grafana':
            subTabVal = 0;
            break;
          case 'prometheus':
            subTabVal = 1;
            break;
        }
      }
    }
    this.state = {
      completed : {},
      k8sconfig,
      meshAdapters,
      grafana,
      prometheus,
      tabVal,
      subTabVal,

      isMeshConfigured : k8sconfig.clusterConfigured,

      // Array of scanned prometheus urls
      scannedPrometheus : [],
      // Array of scanned grafan urls
      scannedGrafana : [],

      meshmodelSummarySelector : { type : "components" },
      meshmodelSummarySelectorList : ["components", "relationships"],
      meshmodelSummary : [],
      meshmodelSummarySubscription : null,
      meshmodelSummaryQuery : null,
    };

    this.systemResetRef = React.createRef();
  }

  static getDerivedStateFromProps(props, state) {
    let st = {};
    if (JSON.stringify(props.k8sconfig) !== JSON.stringify(state.k8sconfig)
      || JSON.stringify(props.meshAdapters) !== JSON.stringify(state.meshAdapters)) {
      st = {
        k8sconfig : props.k8sconfig,
        meshAdapters : props.meshAdapters,
        grafana : props.grafana,
        prometheus : props.prometheus,
      };
    }
    const compare = (arr1, arr2) => arr1.every((val, ind) => val === arr2[ind])

    if (props.telemetryUrls.grafana.length !== state.scannedGrafana.length || !(compare(props.telemetryUrls.grafana, state.scannedGrafana))) {
      st.scannedGrafana = props.telemetryUrls.grafana
    }

    if (props.telemetryUrls.prometheus.length !== state.scannedPrometheus.length || !(compare(props.telemetryUrls.prometheus, state.scannedPrometheus))) {
      st.scannedPrometheus = props.telemetryUrls.prometheus
    }
    return st;
  }

  disposeMeshModelSummarySubscriptions = () => {
    this.state.meshmodelSummarySubscription && this.state.meshmodelSummarySubscription.dispose();
    this.state.meshmodelSummaryQuery && this.state.meshmodelSummaryQuery.unsubscribe();
  }

  initDashboardMeshModelSummaryQuery = () => {
    const self = this;
    let selector = self.state.meshmodelSummarySelector;

    if (self._isMounted) {
      // @ts-ignore
      const meshmodelSummaryQuery = fetchMeshModelSummary(selector).subscribe({
        next : (res) => {
          this.setState({ meshmodelSummary : res?.meshmodelSummary })
        },
        error : (err) => console.log(err),
      })

      this.setState({ meshmodelSummaryQuery });
    }
  }

  initDashboardMeshModelSummarySubscription = () => {
    const self = this;
    let selector = self.state.meshmodelSummarySelector;

    if (self._isMounted) {
      // @ts-ignore
      const meshmodelSummarySubscription = subscribeMeshModelSummary((res) => {
        this.setState({ meshmodelSummary : res?.meshmodelSummary })
      }, {
        selector : selector
      });
      this.setState({ meshmodelSummarySubscription });
    }
  }

  disposeSubscriptions = () => {
    this.disposeMeshModelSummarySubscriptions()
  }

  componentDidMount() {
    this._isMounted = true

    this.disposeSubscriptions()

    if (this._isMounted) {
      this.initDashboardMeshModelSummaryQuery();
      this.initDashboardMeshModelSummarySubscription();
    }
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevState?.meshmodelSummarySelector !== this.state?.meshmodelSummarySelector) {
      this.disposeMeshModelSummarySubscriptions();
      this.initDashboardMeshModelSummaryQuery();
      this.initDashboardMeshModelSummarySubscription();
    }
  }

  emptyStateMessageForMeshModelSummary = () => {
    return "No MeshModel registered."
  }

  /**
   * MeshModelSummaryCard takes in the meshmodel related data
   * and renders a table with registered meshmodel information of
   * the selected type of model (like relationships, components etc)
   * @param {
   * {
   *   kind, count
   * }[]
   * } meshmodelSummary
   */
  MeshModelSummaryCard = (meshmodelSummary = []) => {
    const self = this;
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
        name : "name",
        label : "Name",
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
          }
        },
      },
    ]

    const options = {
      filter : false,
      selectableRows : "none",
      responsive : "standard",
      print : false,
      download : false,
      viewColumns : false,
      pagination : false,
      fixedHeader : true,
      customToolbar : () => {
        return (
          <>
            {self.state.meshmodelSummarySelectorList && (
              <Select
                value={self.state.meshmodelSummarySelector.type}
                onChange={(e) =>
                  self.setState({ meshmodelSummarySelector : { type : e.target.value } })
                }
              >
                {self.state.meshmodelSummarySelectorList && self.state.meshmodelSummarySelectorList.map((opts) => <MenuItem key={opts} value={opts}>{opts}</MenuItem>)}
              </Select>
            )}
          </>
        )
      }
    }

    if (Array.isArray(meshmodelSummary) && meshmodelSummary?.length)
      return (
        <Paper elevation={1} style={{ padding : "2rem" }}>
          <MuiThemeProvider theme={ theme.palette.type == "dark" ?  configurationTableThemeDark() : configurationTableTheme() }>
            <DataTable
              title={
                <>
                  <div style={{ display : "flex", alignItems : "center", marginBottom : "1rem" }}>
                    <img src={"/static/img/all_mesh.svg"} className={this.props.classes.icon} style={{ marginRight : "0.75rem" }} />
                    <Typography variant="h6">Registered MeshModel</Typography>
                  </div>
                </>
              }
              data={meshmodelSummary}
              options={options}
              columns={columns}
            />
          </MuiThemeProvider>
        </Paper>
      );

    return null;
  };

  showMeshModelSummary = () => {
    const self = this;
    return (
      <>
        {self?.state?.meshmodelSummary[self?.state?.meshmodelSummarySelector?.type] && self?.state?.meshmodelSummary[self?.state?.meshmodelSummarySelector?.type].length > 0
          ? (
            self.MeshModelSummaryCard(self?.state?.meshmodelSummary[self?.state?.meshmodelSummarySelector?.type])
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
              <Typography style={{ fontSize : "1.5rem", marginBottom : "2rem" }} align="center" color="textSecondary">
                {this.emptyStateMessageForMeshModelSummary()}
              </Typography>
            </div>
          )}
      </>
    );
  };

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress : false });
    const self = this;
    this.props.enqueueSnackbar(`${msg}: ${error}`, {
      variant : "error",
      action : (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration : 7000,
    });
  };

  handleChange = (val) => {
    const self = this;
    return (event, newVal) => {
      if (val === 'tabVal') {
        let newRoute = this.props.router.route;
        switch (newVal) {
          case 0:
            newRoute += '#environment'
            break;
          case 1:
            newRoute += '#service-mesh'
            break;
          case 2:
            newRoute += '#metrics'
            break;
          case 3:
            newRoute += '#system'
            break;
          case 4:
            newRoute += '#meshmodel-summary'
          // case 3:
          //   newRoute += '#performance'
          //   break;
        }
        if (this.props.router.route != newRoute)
          this.props.router.push(newRoute)
        self.setState({ tabVal : newVal });
      } else if (val === 'subTabVal') {
        let newRoute = this.props.router.route;
        switch (newVal) {
          case 0:
            if (self.state.tabVal == 0)
              newRoute += '#environment/outclusterconfig'
            else if (self.state.tabVal == 2)
              newRoute += '#metrics/grafana'
            break;
          case 1:
            if (self.state.tabVal == 0)
              newRoute += '#environment/inclusterconfig'
            else if (self.state.tabVal == 2)
              newRoute += '#metrics/prometheus'
            break;
        }
        if (this.props.router.route != newRoute)
          this.props.router.push(newRoute)
        self.setState({ subTabVal : newVal });
      }
    };
  }

  handleResetDatabase = () => {
    return async () => {
      let responseOfResetDatabase = await this.systemResetRef.current.show({
        title : "Reset Meshery Database?",
        subtitle : "Are you sure that you want to purge all data?",
        options : ["RESET", "CANCEL"]
      });
      if (responseOfResetDatabase === "RESET") {
        this.props.updateProgress({ showProgress : true });
        const self = this;
        resetDatabase({
          selector : {
            clearDB : "true",
            ReSync : "true",
            hardReset : "true",
          },
          k8scontextID : ""
        }).subscribe({
          next : (res) => {
            self.props.updateProgress({ showProgress : false });
            if (res.resetStatus === "PROCESSING") {
              this.props.enqueueSnackbar(`Database reset successful.`, {
                variant : "success",
                action : (key) => (
                  <IconButton key="close" aria-label="close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
                    <CloseIcon />
                  </IconButton>
                ),
                autohideduration : 3000,
              })
            }
          },
          error : self.handleError("Database is not reachable, try restarting server.")
        });
      }
    }
  }


  render() {
    const { classes } = this.props;
    const {
      tabVal, subTabVal, k8sconfig, meshAdapters,
    } = this.state;

    let backToPlay = '';
    if (k8sconfig.clusterConfigured === true && meshAdapters.length > 0) {
      backToPlay = (
        <div className={classes.backToPlay}>
          <Link href="/management">
            <div className={classes.link}>
              <FontAwesomeIcon icon={faArrowLeft} transform="grow-4" fixedWidth />
              {' '}
              You are ready to manage cloud native infrastructure
            </div>
          </Link>
        </div>
      );
    }
    return (
      <div className={classes.wrapperClss}>
        <Paper square className={classes.wrapperClss}>
          <Tabs
            value={tabVal}
            className={classes.tabs}
            onChange={this.handleChange('tabVal')}
            variant="fullWidth"
            indicatorColor="primary"
            textColor="primary"
          >
            <Tooltip title="Identify your cluster" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faCloud}  style={iconMedium} />
                }
                label="Environment"
                data-cy="tabEnvironment"
              />
            </Tooltip>
            <Tooltip title="Connect Meshery Adapters" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faMendeley}  style={iconMedium}/>
                }
                label="Adapters"
                data-cy="tabServiceMeshes"
              />
            </Tooltip>
            <Tooltip title="Configure Metrics backends" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faPoll}   style={iconMedium}/>
                }
                label="Metrics"
                tab="tabMetrics"
              />
            </Tooltip>
            <Tooltip title="Reset System" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faDatabase}  style={iconMedium} />
                }
                label="Reset"
                tab="systemReset"
              />
            </Tooltip>
            <Tooltip title="MeshModel Summary" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faFileInvoice}  style={iconMedium} />
                }
                label="MeshModel Summary"
                tab="meshmodelSummary"
              />
            </Tooltip>

            {/*NOTE: Functionality of performance tab will be modified, until then keeping it and the related code commented */}

            {/* <Tooltip title="Choose Performance Test Defaults" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faTachometerAlt} transform={mainIconScale} fixedWidth />
                }
                label="Performance"
                tab="tabPerformance"
              />
            </Tooltip> */}
          </Tabs>
        </Paper>
        {tabVal === 0 && (
          <MeshConfigComponent />
        )}
        {tabVal === 1 && (
          <TabContainer>
            <MeshAdapterConfigComponent />
          </TabContainer>
        )}
        {tabVal === 2
          && (
            <TabContainer>
              <AppBar position="static" color="default">
                <Tabs
                  value={subTabVal}
                  className={classes.tabs}
                  onChange={this.handleChange('subTabVal')}
                  indicatorColor="primary"
                  textColor="primary"
                  variant="fullWidth"
                >
                  <Tab className={classes.tab} label={(
                    <div className={classes.iconText}>
                      Grafana
                      <img src="/static/img/grafana_icon.svg" className={classes.icon} />
                    </div>
                  )}
                  />
                  <Tab className={classes.tab} label={(
                    <div className={classes.iconText}>
                      Prometheus
                      <img src="/static/img/prometheus_logo_orange_circle.svg" className={classes.icon} />
                    </div>
                  )}
                  />
                </Tabs>
              </AppBar>
              {subTabVal === 0 && (
                <TabContainer>
                  <GrafanaComponent scannedGrafana={this.state.scannedGrafana} isMeshConfigured={this.state.isMeshConfigured} />
                </TabContainer>
              )}
              {subTabVal === 1 && (
                <TabContainer>
                  <PrometheusComponent scannedPrometheus={this.state.scannedPrometheus} isMeshConfigured={this.state.isMeshConfigured} />
                </TabContainer>
              )}
            </TabContainer>
          )}
        {tabVal === 3 && (
          <TabContainer>
            <div className={classes.container}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={this.handleResetDatabase()}
                className={classes.DBBtn}
                data-cy="btnResetDatabase"

              >
                <Typography> System Reset </Typography>
              </Button>
            </div>
          </TabContainer>
        )}
        {tabVal === 4 && (
          <TabContainer>
            <MeshModelComponent  showMeshModelSummary={this.showMeshModelSummary} />
          </TabContainer>
        )}
        {/* {tabVal === 3 && (
          <TabContainer>
            <MesherySettingsPerformanceComponent />

          </TabContainer>
        )} */}

        {backToPlay}
        <PromptComponent ref={this.systemResetRef} />
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig');
  const meshAdapters = state.get('meshAdapters').toJS();
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  const selectedK8sContexts = state.get('selectedK8sContexts');
  const telemetryUrls = state.get('telemetryURLs').toJS();
  return {
    k8sconfig,
    meshAdapters,
    grafana,
    prometheus,
    selectedK8sContexts,
    telemetryUrls,
  };
};

const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

MesherySettings.propTypes = { classes : PropTypes.object, };

export default withStyles(styles, { withTheme : true })(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(MesherySettings)))
);
