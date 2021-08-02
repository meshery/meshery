import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { withRouter } from 'next/router';
import { connect } from 'react-redux';
import { bindActionCreators } from "redux";
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Typography from '@material-ui/core/Typography';
import { AppBar, Paper, Tooltip, IconButton } from '@material-ui/core';
import CloseIcon from "@material-ui/icons/Close";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faCloud, faPoll, faTachometerAlt } from '@fortawesome/free-solid-svg-icons';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import Link from 'next/link';
import MeshConfigComponent from './MeshConfigComponent';
import GrafanaComponent from './GrafanaComponent';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import PrometheusComponent from './PrometheusComponent';
import MesherySettingsPerformanceComponent from "../components/MesherySettingsPerformanceComponent";
import dataFetch from '../lib/data-fetch';
import { updateProgress } from "../lib/store";
import { withSnackbar } from "notistack";

const styles = (theme) => ({
  root: {
    flexGrow: 1,
    maxWidth: '100%',
    height: 'auto',
  },
  tab: {
    minWidth: 40,
    paddingLeft: 0,
    paddingRight: 0,
  },
  icon: {
    display: 'inline',
    verticalAlign: 'text-top',
    width: theme.spacing(1.75),
    marginLeft: theme.spacing(0.5),
  },
  iconText: {
    display: 'inline',
    verticalAlign: 'middle',
  },
  backToPlay: {
    margin: theme.spacing(2),
  },
  link: {
    cursor: 'pointer',
  },
});

function TabContainer(props) {
  return (
    <Typography
      component="div"
      style={{
        paddingTop: 2,
      }}
    >
      {props.children}
    </Typography>
  );
}

TabContainer.propTypes = {
  children: PropTypes.node.isRequired,
};

class MesherySettings extends React.Component {
  constructor(props) {
    super(props);
    const {
      k8sconfig, meshAdapters, grafana, prometheus, router: { asPath }
    } = props;

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
        case 'performance':
          tabVal = 3;
          break;
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
      completed: {},
      k8sconfig,
      meshAdapters,
      grafana,
      prometheus,
      tabVal,
      subTabVal,

      // Array of scanned prometheus urls
      scannedPrometheus: [],
      // Array of scanned grafan urls
      scannedGrafana: []
    };
  }

  static getDerivedStateFromProps(props, state) {
    if (JSON.stringify(props.k8sconfig) !== JSON.stringify(state.k8sconfig)
      || JSON.stringify(props.meshAdapters) !== JSON.stringify(state.meshAdapters)) {
      return {
        k8sconfig: props.k8sconfig,
        meshAdapters: props.meshAdapters,
        grafana: props.grafana,
        prometheus: props.prometheus,
      };
    }
    return null;
  }

  componentDidMount() {
    this.fetchPromGrafanaScanData()
  }

  fetchPromGrafanaScanData = () => {
    const self = this;
    self.props.updateProgress({ showProgress: true });
    dataFetch(
      '/api/promGrafana/scan',
      {
        credentials: "same-origin",
        method: "GET",
        credentials: "include",
      },
      (result) => {
        self.props.updateProgress({ showProgress: false });
        if (!result) return;

        if (Array.isArray(result.prometheus)) {
          const urls = self.extractURLFromScanData(result.prometheus);
          self.setState(state => ({ scannedPrometheus: [...state.scannedPrometheus, ...urls] }));
        }

        if (Array.isArray(result.grafana)) {
          const urls = self.extractURLFromScanData(result.grafana);
          self.setState(state => ({ scannedGrafana: [...state.scannedGrafana, ...urls] }));
        }
      },
      self.handleError("Unable to fetch prometheus and grafana scan data")
    )
  }

  /**
   * extractURLFromScanData scans the ingress urls from the
   * mesh scan data and returns an array of the response
   * @param {object[]} scannedData 
   * @returns {string[]}
   */
  extractURLFromScanData = (scannedData) => {
    const result = [];
    scannedData.forEach(data => {
      // Add loadbalancer based url
      if (Array.isArray(data.status?.loadBalancer?.ingress)) {
        data.status.loadBalancer.ingress.forEach(lbdata => {
          let protocol = "http";

          // Iterate over ports exposed by the service
          if (Array.isArray(data.spec.ports)) {
            data.spec.ports.forEach(({ port }) => {
              if (port === 443) protocol = "https";

              // From kubernetes v1.19 docs
              // Hostname is set for load-balancer ingress points that are DNS based (typically AWS load-balancers)
              // IP is set for load-balancer ingress points that are IP based (typically GCE or OpenStack load-balancers)
              let address = lbdata.ip || lbdata.hostname;
              if (address) result.push(`${protocol}://${address}:${port}`);
            })
          }
        })
      }

      // Add clusterip based url
      // As per kubernetes v1.19 api, "None", "" as well as a valid ip is a valid clusterIP
      // Looking for valid ipv4 address
      if (data.spec.clusterIP?.match(/^((25[0-5]|(2[0-4]|1[0-9]|[1-9]|)[0-9])(\.(?!$)|$)){4}$/g)?.[0]) {
        let protocol = "http";
        if (Array.isArray(data.spec.ports)) {
          data.spec.ports.forEach(({ port }) => {
            if (port === 443) protocol = "https";
            result.push(`${protocol}://${data.spec.clusterIP}:${port}`);
          })
        }
      }
    })

    return result
  }

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress: false });
    const self = this;
    this.props.enqueueSnackbar(`${msg}: ${error}`, {
      variant: "error",
      action: (key) => (
        <IconButton key="close" aria-label="Close" color="inherit" onClick={() => self.props.closeSnackbar(key)}>
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration: 7000,
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
            newRoute += '#performance'
            break;
        }
        if (this.props.router.route != newRoute)
          this.props.router.push(newRoute)
        self.setState({ tabVal: newVal });
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
        self.setState({ subTabVal: newVal });
      }
    };
  }

  render() {
    const { classes } = this.props;
    const {
      tabVal, subTabVal, k8sconfig, meshAdapters,
    } = this.state;

    const mainIconScale = 'grow-10';
    let backToPlay = '';
    if (k8sconfig.clusterConfigured === true && meshAdapters.length > 0) {
      backToPlay = (
        <div className={classes.backToPlay}>
          <Link href="/management">
            <div className={classes.link}>
              <FontAwesomeIcon icon={faArrowLeft} transform="grow-4" fixedWidth />
              {' '}
              You are all set to manage service meshes
            </div>
          </Link>
        </div>
      );
    }
    return (
      <div className={classes.root}>
        <Paper square className={classes.root}>
          <Tabs
            value={tabVal}
            onChange={this.handleChange('tabVal')}
            variant="fullWidth"
            indicatorColor="primary"
            textColor="primary"
          >
            <Tooltip title="Identify your cluster" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faCloud} transform={mainIconScale} />
                }
                label="Environment"
                data-cy="tabEnvironment"
              />
            </Tooltip>
            <Tooltip title="Connect Meshery Adapters" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faMendeley} transform={mainIconScale} />
                }
                label="Service Meshes"
                data-cy="tabServiceMeshes"
              />
            </Tooltip>
            <Tooltip title="Configure Metrics backends" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faPoll} transform={mainIconScale} fixedWidth />
                }
                label="Metrics"
                tab="tabMetrics"
              />
            </Tooltip>
            <Tooltip title="Choose Performance Test Defaults" placement="top">
              <Tab
                className={classes.tab}
                icon={
                  <FontAwesomeIcon icon={faTachometerAlt} transform={mainIconScale} fixedWidth />
                }
                label="Performance"
                tab="tabPerformance"
              />
            </Tooltip>
          </Tabs>
        </Paper>
        {tabVal === 0 && (
          <TabContainer>
            <AppBar position="static" color="default">
              <Tabs
                value={subTabVal}
                onChange={this.handleChange('subTabVal')}
                indicatorColor="primary"
                textColor="primary"
                variant="fullWidth"
              >
                <Tab className={classes.tab} label="Out of Cluster Deployment" data-cy="tabOutOfClusterDeployment" />
                <Tab className={classes.tab} label="In Cluster Deployment" data-cy="tabInClusterDeployment" />
              </Tabs>
            </AppBar>
            {subTabVal === 0 && (
              <TabContainer>
                <MeshConfigComponent tabs={subTabVal} />
              </TabContainer>
            )}
            {subTabVal === 1 && (
              <TabContainer>
                <MeshConfigComponent tabs={subTabVal} />
              </TabContainer>
            )}
          </TabContainer>
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
                  className={classes.tab}
                  onChange={this.handleChange('subTabVal')}
                  indicatorColor="primary"
                  textColor="primary"
                  variant="fullWidth"
                >
                  <Tab label={(
                    <div className={classes.iconText}>
                      Grafana
                      <img src="/static/img/grafana_icon.svg" className={classes.icon} />
                    </div>
                  )}
                  />
                  <Tab label={(
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
                  <GrafanaComponent scannedGrafana={this.state.scannedGrafana} />
                </TabContainer>
              )}
              {subTabVal === 1 && (
                <TabContainer>
                  <PrometheusComponent scannedPrometheus={this.state.scannedPrometheus} />
                </TabContainer>
              )}
            </TabContainer>
          )}
        {tabVal === 3 && (
          <TabContainer>
            <MesherySettingsPerformanceComponent />

          </TabContainer>
        )}

        {backToPlay}
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const k8sconfig = state.get('k8sConfig').toJS();
  const meshAdapters = state.get('meshAdapters').toJS();
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  return {
    k8sconfig,
    meshAdapters,
    grafana,
    prometheus,
  };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

MesherySettings.propTypes = {
  classes: PropTypes.object,
};

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withSnackbar(MesherySettings)))
);
