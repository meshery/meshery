import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { withRouter } from 'next/router';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import { AppBar, Paper, Tooltip, Typography } from '@material-ui/core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faPoll, faDatabase, faFileInvoice } from '@fortawesome/free-solid-svg-icons';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import Link from 'next/link';
import GrafanaComponent from './telemetry/grafana/GrafanaComponent';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import PrometheusComponent from './telemetry/prometheus/PrometheusComponent';
import { updateProgress } from '../lib/store';
import PromptComponent from './PromptComponent';
import { iconMedium } from '../css/icons.styles';
import MeshModelComponent from './MeshModelRegistry/MeshModelComponent';
import DatabaseSummary from './DatabaseSummary';
import {
  getComponentsDetail,
  getModelsDetail,
  getRelationshipsDetail,
  getMeshModelRegistrants,
} from '../api/meshmodel';
import { withNotify } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';

const styles = (theme) => ({
  wrapperClss: {
    flexGrow: 1,
    maxWidth: '100%',
    height: 'auto',
  },
  tab: {
    minWidth: 40,
    paddingLeft: 0,
    paddingRight: 0,
    '&.Mui-selected': {
      color: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  tabs: {
    '& .MuiTabs-indicator': {
      backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
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
  backToPlay: { margin: theme.spacing(2) },
  link: { cursor: 'pointer' },
  container: {
    display: 'flex',
    justifyContent: 'center',
    marginTop: theme.spacing(2),
  },
  paper: {
    maxWidth: '90%',
    margin: 'auto',
    overflow: 'hidden',
  },
  topToolbar: {
    marginBottom: '2rem',
    display: 'flex',
    justifyContent: 'space-between',
    paddingLeft: '1rem',
    maxWidth: '90%',
  },
  dashboardSection: {
    backgroundColor: theme.palette.secondary.elevatedComponents,
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
  },
  cardHeader: { fontSize: theme.spacing(2) },
  card: {
    height: '100%',
    marginTop: theme.spacing(2),
  },
  cardContent: { height: '100%' },
  boxWrapper: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'end',
    flexDirection: 'row',
    flexWrap: 'wrap',
    height: '60vh',
    borderRadius: 0,
    color: 'white',
    ['@media (max-width: 455px)']: {
      width: '100%',
    },
    zIndex: 5,
  },
  box: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
    width: 300,
    height: 300,
    backgroundColor: theme.palette.secondary.dark,
    border: '0px solid #000',
    boxShadow: theme.shadows[5],
    margin: theme.spacing(2),
    cursor: 'pointer',
  },
});

function TabContainer(props) {
  return (
    <Typography component="div" style={{ paddingTop: 2 }}>
      {props.children}
    </Typography>
  );
}

TabContainer.propTypes = { children: PropTypes.node.isRequired };

//TODO: Tabs are hardcoded everywhere
class MesherySettings extends React.Component {
  constructor(props) {
    super(props);
    const {
      k8sconfig,
      meshAdapters,
      grafana,
      prometheus,
      router: { asPath },
    } = props;

    this._isMounted = false;

    //TODO: Extract this logic and clean it
    let tabVal = 0,
      subTabVal = 0;
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
          break;
        case 'registry':
          tabVal = 4;
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
      modelsCount: 0,
      componentsCount: 0,
      relationshipsCount: 0,
      registrantCount: 0,
      isMeshConfigured: k8sconfig.clusterConfigured,

      // Array of scanned prometheus urls
      scannedPrometheus: [],
      // Array of scanned grafan urls
      scannedGrafana: [],
    };

    this.systemResetPromptRef = React.createRef();
  }

  static getDerivedStateFromProps(props, state) {
    let st = {};
    if (
      JSON.stringify(props.k8sconfig) !== JSON.stringify(state.k8sconfig) ||
      JSON.stringify(props.meshAdapters) !== JSON.stringify(state.meshAdapters)
    ) {
      st = {
        k8sconfig: props.k8sconfig,
        meshAdapters: props.meshAdapters,
        grafana: props.grafana,
        prometheus: props.prometheus,
      };
    }
    const compare = (arr1, arr2) => arr1.every((val, ind) => val === arr2[ind]);

    if (
      props.telemetryUrls.grafana.length !== state.scannedGrafana.length ||
      !compare(props.telemetryUrls.grafana, state.scannedGrafana)
    ) {
      st.scannedGrafana = props.telemetryUrls.grafana;
    }

    if (
      props.telemetryUrls.prometheus.length !== state.scannedPrometheus.length ||
      !compare(props.telemetryUrls.prometheus, state.scannedPrometheus)
    ) {
      st.scannedPrometheus = props.telemetryUrls.prometheus;
    }
    return st;
  }

  async componentDidMount() {
    try {
      const modelsResponse = await getModelsDetail();
      const componentsResponse = await getComponentsDetail();
      const relationshipsResponse = await getRelationshipsDetail();
      const registrantResponce = await getMeshModelRegistrants();

      const modelsCount = modelsResponse.total_count;
      const componentsCount = componentsResponse.total_count;
      const relationshipsCount = relationshipsResponse.total_count;
      const registrantCount = registrantResponce.total_count;

      this.setState({
        modelsCount,
        componentsCount,
        relationshipsCount,
        registrantCount,
      });
    } catch (error) {
      console.error(error);
    }
  }

  handleError = (msg) => (error) => {
    this.props.updateProgress({ showProgress: false });
    const notify = this.props.notify;
    notify({
      message: `${msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  handleChange = (val) => {
    const self = this;
    return (event, newVal) => {
      if (val === 'tabVal') {
        let newRoute = this.props.router.route;

        switch (newVal) {
          case 0:
            newRoute += '#adapters';
            break;
          case 1:
            newRoute += '#metrics';
            break;
          case 2:
            newRoute += '#registry';
            break;
          case 3:
            newRoute += '#reset';
        }
        if (this.props.router.route != newRoute) this.props.router.push(newRoute);
        self.setState({ tabVal: newVal });
      } else if (val === 'subTabVal') {
        let newRoute = this.props.router.route;
        switch (newVal) {
          case 0:
            if (self.state.tabVal == 0) newRoute += '#environment/outclusterconfig';
            else if (self.state.tabVal == 2) newRoute += '#metrics/grafana';
            else if (self.state.tabVal == 4) newRoute += '#metrics/models';
            break;
          case 1:
            if (self.state.tabVal == 0) newRoute += '#environment/inclusterconfig';
            else if (self.state.tabVal == 2) newRoute += '#metrics/prometheus';
            else if (self.state.tabVal == 4) newRoute += '#metrics/components';
            break;
          case 2:
            if (self.state.tabVal == 0) newRoute += '#environment/inclusterconfig';
            else if (self.state.tabVal == 2) newRoute += '#metrics/prometheus';
            else if (self.state.tabVal == 4) newRoute += '#metrics/relationships';
            break;
        }
        if (this.props.router.route != newRoute) this.props.router.push(newRoute);
        self.setState({ subTabVal: newVal });
      }
    };
  };

  render() {
    const { classes } = this.props;
    const { tabVal, subTabVal, k8sconfig, meshAdapters } = this.state;

    let backToPlay = '';
    if (k8sconfig.clusterConfigured === true && meshAdapters.length > 0) {
      backToPlay = (
        <div className={classes.backToPlay}>
          <Link href="/management">
            <div className={classes.link}>
              <FontAwesomeIcon icon={faArrowLeft} transform="grow-4" fixedWidth />
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
            <Tooltip title="Connect Meshery Adapters" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faMendeley} style={iconMedium} />}
                label="Adapters"
                data-cy="tabServiceMeshes"
              />
            </Tooltip>
            <Tooltip title="Configure Metrics backends" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faPoll} style={iconMedium} />}
                label="Metrics"
                tab="tabMetrics"
              />
            </Tooltip>
            <Tooltip title="Registry" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faFileInvoice} style={iconMedium} />}
                label="Registry"
                tab="registry"
              />
            </Tooltip>
            <Tooltip title="Reset System" placement="top">
              <Tab
                className={classes.tab}
                icon={<FontAwesomeIcon icon={faDatabase} style={iconMedium} />}
                label="Reset"
                tab="systemReset"
              />
            </Tooltip>
          </Tabs>
        </Paper>
        {tabVal === 0 && (
          <TabContainer>
            <MeshAdapterConfigComponent />
          </TabContainer>
        )}
        {tabVal === 1 && (
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
                <Tab
                  className={classes.tab}
                  label={
                    <div className={classes.iconText}>
                      Grafana
                      <img src="/static/img/grafana_icon.svg" className={classes.icon} />
                    </div>
                  }
                />
                <Tab
                  className={classes.tab}
                  label={
                    <div className={classes.iconText}>
                      Prometheus
                      <img
                        src="/static/img/prometheus_logo_orange_circle.svg"
                        className={classes.icon}
                      />
                    </div>
                  }
                />
              </Tabs>
            </AppBar>
            {subTabVal === 0 && (
              <TabContainer>
                <GrafanaComponent
                  scannedGrafana={this.state.scannedGrafana}
                  isMeshConfigured={this.state.isMeshConfigured}
                />
              </TabContainer>
            )}
            {subTabVal === 1 && (
              <TabContainer>
                <PrometheusComponent
                  scannedPrometheus={this.state.scannedPrometheus}
                  isMeshConfigured={this.state.isMeshConfigured}
                />
              </TabContainer>
            )}
          </TabContainer>
        )}
        {tabVal === 2 && (
          <TabContainer>
            <TabContainer>
              <TabContainer>
                <MeshModelComponent
                  modelsCount={this.state.modelsCount}
                  componentsCount={this.state.componentsCount}
                  relationshipsCount={this.state.relationshipsCount}
                  registrantCount={this.state.registrantCount}
                />
              </TabContainer>
            </TabContainer>
            {/* </div> */}
          </TabContainer>
        )}
        {tabVal === 3 && (
          <TabContainer>
            <DatabaseSummary promptRef={this.systemResetPromptRef} />
          </TabContainer>
        )}
        {backToPlay}
        <PromptComponent ref={this.systemResetPromptRef} />
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

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

MesherySettings.propTypes = { classes: PropTypes.object };

export default withStyles(styles, { withTheme: true })(
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withNotify(MesherySettings))),
);
