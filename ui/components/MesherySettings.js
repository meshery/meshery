import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { useRouter } from 'next/router';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import { AppBar, Paper, Tooltip, Typography } from '@material-ui/core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faArrowLeft,
  faCloud,
  faPoll,
  faDatabase,
  faFileInvoice,
} from '@fortawesome/free-solid-svg-icons';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import Link from 'next/link';
import MeshConfigComponent from './MeshConfigComponent';
import GrafanaComponent from './telemetry/grafana/GrafanaComponent';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import PrometheusComponent from './telemetry/prometheus/PrometheusComponent';
import { updateProgress } from '../lib/store';
import PromptComponent from './PromptComponent';
import { iconMedium } from '../css/icons.styles';
import MeshModelComponent from './MeshModelComponent';
// import CredentialIcon from '../assets/icons/CredentialIcon';
// import MesheryCredentialComponent from './MesheryCredentialComponent';
import DatabaseSummary from './DatabaseSummary';
import { getComponentsDetail, getModelsDetail, getRelationshipsDetail } from '../api/meshmodel';
// import { useNotification } from '../utils/hooks/useNotification';
// import { EVENT_TYPES } from '../lib/event-types';

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

const MesherySettings = (props) => {
  const { k8sconfig, meshAdapters, grafana, prometheus, classes } = props;
  const router = useRouter();
  // const  notify  = useNotification();
  const [modelsCount, setModelsCount] = useState(0);
  const [componentsCount, setComponentsCount] = useState(0);
  const [relationshipsCount, setRelationshipsCount] = useState(0);
  const [scannedGrafana, setScannedGrafana] = useState([]);
  const [scannedPrometheus, setScannedPrometheus] = useState([]);
  const [tabVal, setTabVal] = useState(0);
  const [subTabVal, setSubTabVal] = useState(0);
  const [isMeshConfigured] = useState(k8sconfig.clusterConfigured);
  const systemResetPromptRef = useRef(null);

  useEffect(() => {
    if (
      JSON.stringify(props.k8sconfig) !== JSON.stringify(k8sconfig) ||
      JSON.stringify(props.meshAdapters) !== JSON.stringify(meshAdapters)
    ) {
      k8sconfig(props.k8sconfig);
      meshAdapters(props.meshAdapters);
      grafana(props.grafana);
      prometheus(props.prometheus);
    }
    const compare = (arr1, arr2) => arr1.every((val, ind) => val === arr2[ind]);

    if (
      props.telemetryUrls.grafana.length !== scannedGrafana.length ||
      !compare(props.telemetryUrls.grafana, scannedGrafana)
    ) {
      setScannedGrafana(props.telemetryUrls.grafana);
    }

    if (
      props.telemetryUrls.prometheus.length !== scannedPrometheus.length ||
      !compare(props.telemetryUrls.prometheus, scannedPrometheus)
    ) {
      setScannedPrometheus(props.telemetryUrls.prometheus);
    }
  }, [props.k8sconfig, props.meshAdapters, props.grafana, props.prometheus, props.telemetryUrls]);

  const fetchData = async () => {
    try {
      const modelsResponse = await getModelsDetail();
      const componentsResponse = await getComponentsDetail();
      const relationshipsResponse = await getRelationshipsDetail();

      setModelsCount(modelsResponse.total_count);
      setComponentsCount(componentsResponse.total_count);
      setRelationshipsCount(relationshipsResponse.total_count);
    } catch (error) {
      console.error(error);
    }
  };
  useEffect(() => {
    fetchData();
    return () => {
      setModelsCount(0);
      setComponentsCount(0);
      setRelationshipsCount(0);
    };
  }, []);

  // const handleError = (msg) => (error) => {
  //   props.updateProgress({ showProgress : false });
  //   notify({ message : `${msg}: ${error}`, event_type : EVENT_TYPES.ERROR, details : error.toString() });
  // };

  const handleChange = (val) => (event, newVal) => {
    if (val === 'tabVal') {
      let newRoute = router.route;

      switch (newVal) {
        case 0:
          newRoute += '#environment';
          break;
        case 1:
          newRoute += '#service-mesh';
          break;
        case 2:
          newRoute += '#metrics';
          break;
        case 3:
          newRoute += '#system';
          break;
        case 4:
          newRoute += '#registry';
        // case 3:
        //   newRoute += '#performance'
        //   break;
      }
      if (router.route != newRoute) router.push(newRoute);
      setTabVal(newVal);
    } else if (val === 'subTabVal') {
      let newRoute = router.route;
      switch (newVal) {
        case 0:
          if (tabVal == 0) newRoute += '#environment/outclusterconfig';
          else if (tabVal == 2) newRoute += '#metrics/grafana';
          else if (tabVal == 4) newRoute += '#metrics/models';
          break;
        case 1:
          if (tabVal == 0) newRoute += '#environment/inclusterconfig';
          else if (tabVal == 2) newRoute += '#metrics/prometheus';
          else if (tabVal == 4) newRoute += '#metrics/components';
          break;
        case 2:
          if (tabVal == 0) newRoute += '#environment/inclusterconfig';
          else if (tabVal == 2) newRoute += '#metrics/prometheus';
          else if (tabVal == 4) newRoute += '#metrics/relationships';
          break;
      }
      if (router.route !== newRoute) router.push(newRoute);
      setSubTabVal(newVal);
    }
  };

  const backToPlays = () => {
    if (k8sconfig.clusterConfigured === true && meshAdapters.length > 0) {
      return (
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
  };
  return (
    <div className={classes.wrapperClss}>
      <Paper square className={classes.wrapperClss}>
        <Tabs
          value={tabVal}
          className={classes.tabs}
          onChange={handleChange('tabVal')}
          variant="fullWidth"
          indicatorColor="primary"
          textColor="primary"
        >
          <Tooltip title="Identify your cluster" placement="top">
            <Tab
              className={classes.tab}
              icon={<FontAwesomeIcon icon={faCloud} style={iconMedium} />}
              label="Environment"
              data-cy="tabEnvironment"
            />
          </Tooltip>
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
          <Tooltip title="Reset System" placement="top">
            <Tab
              className={classes.tab}
              icon={<FontAwesomeIcon icon={faDatabase} style={iconMedium} />}
              label="Reset"
              tab="systemReset"
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
          {/* <Tooltip title="Credential" placement="top">
              <Tab
                className={classes.tab}
                icon={<CredentialIcon width="1.5rem" />}
                label="Credentials"
                tab="credential"
              />
            </Tooltip> */}

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
      {tabVal === 0 && <MeshConfigComponent />}
      {tabVal === 1 && (
        <TabContainer>
          <MeshAdapterConfigComponent />
        </TabContainer>
      )}
      {tabVal === 2 && (
        <TabContainer>
          <AppBar position="static" color="default">
            <Tabs
              value={subTabVal}
              className={classes.tabs}
              onChange={handleChange('subTabVal')}
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
                scannedGrafana={scannedGrafana}
                isMeshConfigured={isMeshConfigured}
              />
            </TabContainer>
          )}
          {subTabVal === 1 && (
            <TabContainer>
              <PrometheusComponent
                scannedPrometheus={scannedPrometheus}
                isMeshConfigured={isMeshConfigured}
              />
            </TabContainer>
          )}
        </TabContainer>
      )}
      {tabVal === 3 && (
        <TabContainer>
          <DatabaseSummary promptRef={systemResetPromptRef} />
        </TabContainer>
      )}
      {tabVal === 4 && (
        <TabContainer>
          <TabContainer>
            <AppBar position="static" color="default">
              <Tabs
                value={subTabVal}
                className={classes.tabs}
                onChange={handleChange('subTabVal')}
                indicatorColor="primary"
                textColor="primary"
                variant="fullWidth"
              >
                <Tab
                  className={classes.tab}
                  label={
                    <div className={classes.iconText}>
                      Models <span style={{ fontWeight: 'bold' }}>({modelsCount})</span>
                    </div>
                  }
                />
                <Tab
                  className={classes.tab}
                  label={
                    <div className={classes.iconText}>
                      Components <span style={{ fontWeight: 'bold' }}>({componentsCount})</span>
                    </div>
                  }
                />
                <Tab
                  className={classes.tab}
                  label={
                    <div className={classes.iconText}>
                      Relationships{' '}
                      <span style={{ fontWeight: 'bold' }}>({relationshipsCount})</span>
                    </div>
                  }
                />
              </Tabs>
            </AppBar>
            {subTabVal === 0 && (
              <TabContainer>
                <MeshModelComponent view="models" />
              </TabContainer>
            )}
            {subTabVal === 1 && (
              <TabContainer>
                <MeshModelComponent view="components" />
              </TabContainer>
            )}
            {subTabVal === 2 && (
              <TabContainer>
                <MeshModelComponent view="relationships" />
              </TabContainer>
            )}
          </TabContainer>
          {/* </div> */}
        </TabContainer>
      )}

      {/* {tabVal === 3 && (
          <TabContainer>
            <MesherySettingsPerformanceComponent />

          </TabContainer>
        )} */}
      {/* {tabVal === 5 && (
          <TabContainer>
            <MesheryCredentialComponent />
          </TabContainer>
        )} */}

      {backToPlays}
      <PromptComponent ref={systemResetPromptRef} />
    </div>
  );
};

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
  connect(mapStateToProps, mapDispatchToProps)(MesherySettings),
);
