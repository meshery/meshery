import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import Link from 'next/link';
import { faArrowLeft, faPoll, faDatabase, faFileInvoice } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import { CustomTooltip } from '@layer5/sistent';
import { withStyles } from '@material-ui/core/styles';
import { Tabs, Tab } from '@material-ui/core';
import { AppBar, Paper, Typography } from '@material-ui/core';
import DatabaseSummary from './DatabaseSummary';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import MeshModelComponent from './MeshModelRegistry/MeshModelComponent';
import PromptComponent from './PromptComponent';
import GrafanaComponent from './telemetry/grafana/GrafanaComponent';
import PrometheusComponent from './telemetry/prometheus/PrometheusComponent';
import { removeDuplicateVersions } from './MeshModelRegistry/helper';
import DefaultError from './General/error-404';
import { updateProgress } from '../lib/store';
import { EVENT_TYPES } from '../lib/event-types';
import {
  getComponentsDetail,
  getMeshModels,
  getRelationshipsDetail,
  getMeshModelRegistrants,
} from '../api/meshmodel';
import { withNotify } from '../utils/hooks/useNotification';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { REGISTRY, METRICS, ADAPTERS, RESET, GRAFANA, PROMETHEUS } from '@/constants/navigator';
import { iconMedium } from '../css/icons.styles';

const styles = (theme) => ({
  wrapperClss: {
    flexGrow: 1,
    maxWidth: '100%',
    height: 'auto',
  },
  tab: {
    width: '25%',
    minWidth: 40,
    paddingLeft: 0,
    paddingRight: 0,
    '&.Mui-selected': {
      color: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  tabs: {
    width: '100%',
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

const settingsRouter = (router) => {
  const { query, push: pushRoute, route } = router;

  const selectedSettingsCategory = query.settingsCategory;
  const selectedTab = query.tab;

  const handleChangeSettingsCategory = (settingsCategory) => {
    if (query.settingsCategory === settingsCategory) {
      return;
    }
    pushRoute(
      `${route}?settingsCategory=${settingsCategory || query.settingsCategory}`,
      undefined,
      { shallow: true },
    );
  };

  const handleChangeSelectedTab = (tab) => {
    if (query.tab === tab) {
      return;
    }
    pushRoute(`${route}?settingsCategory=${selectedSettingsCategory}&tab=${tab}`, undefined, {
      shallow: true,
    });
  };

  const handleChangeSelectedTabCustomCategory = (settingsCategory, tab) => {
    if (query.tab === tab) {
      return;
    }
    pushRoute(`${route}?settingsCategory=${settingsCategory}&tab=${tab}`, undefined, {
      shallow: true,
    });
  };

  return {
    selectedSettingsCategory,
    selectedTab,
    handleChangeSettingsCategory,
    handleChangeSelectedTab,
    handleChangeSelectedTabCustomCategory,
  };
};

//TODO: Tabs are hardcoded everywhere
const MesherySettings = ({
  k8sconfig,
  meshAdapters,
  telemetryUrls,
  updateProgress,
  notify,
  router,
  classes,
}) => {
  const isMeshConfigured = k8sconfig.clusterConfigured;
  const { selectedSettingsCategory, selectedTab } = settingsRouter(router);
  const [tabVal, setTabVal] = useState(selectedSettingsCategory || ADAPTERS);
  const [subTabVal, setSubTabVal] = useState(selectedTab || GRAFANA);
  const [modelsCount, setModelsCount] = useState(0);
  const [componentsCount, setComponentsCount] = useState(0);
  const [relationshipsCount, setRelationshipsCount] = useState(0);
  const [registrantCount, setRegistrantCount] = useState(0);
  const [scannedGrafana, setScannedGrafana] = useState([]);
  const [scannedPrometheus, setScannedPrometheus] = useState([]);

  const systemResetPromptRef = useRef();

  const compare = (arr1, arr2) => arr1.every((val, ind) => val === arr2[ind]);

  useEffect(() => {
    if (!compare(telemetryUrls.grafana, scannedGrafana)) {
      setScannedGrafana(telemetryUrls.grafana);
    }
    if (!compare(telemetryUrls.prometheus, scannedPrometheus)) {
      setScannedPrometheus(telemetryUrls.prometheus);
    }
  }, [k8sconfig, meshAdapters, telemetryUrls, scannedGrafana, scannedPrometheus]);

  useEffect(() => {
    const { selectedSettingsCategory } = settingsRouter(router);
    setTabVal(selectedSettingsCategory);
  }, [selectedSettingsCategory]);

  useEffect(() => {
    async function fetchData() {
      try {
        const modelsResponse = await getMeshModels();
        const componentsResponse = await getComponentsDetail();
        const relationshipsResponse = await getRelationshipsDetail();
        const registrantResponce = await getMeshModelRegistrants();

        const modelsCount = removeDuplicateVersions(modelsResponse.models).length;
        const componentsCount = componentsResponse.total_count;
        const relationshipsCount = relationshipsResponse.total_count;
        const registrantCount = registrantResponce.total_count;

        setModelsCount(modelsCount);
        setComponentsCount(componentsCount);
        setRelationshipsCount(relationshipsCount);
        setRegistrantCount(registrantCount);
      } catch (error) {
        console.error(error);
      }
    }
    fetchData();
  }, []);

  /* eslint-disable no-unused-vars */
  const handleError = (msg) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const handleChange = (val) => {
    const {
      handleChangeSettingsCategory,
      handleChangeSelectedTab,
      handleChangeSelectedTabCustomCategory,
    } = settingsRouter(router);

    return (event, newVal) => {
      if (val === 'tabVal') {
        if (newVal === METRICS) {
          // If the selected tab is "Metrics", set the sub-tab to "Grafana" by default
          handleChangeSelectedTabCustomCategory(METRICS, GRAFANA);
          setTabVal(METRICS);
          setSubTabVal(GRAFANA);
        } else {
          handleChangeSettingsCategory(newVal);
          setTabVal(newVal);
          setSubTabVal(''); // Reset the sub-tab value if needed when switching tabs
        }
      } else if (val === 'subTabVal') {
        // For sub-tabs, ensure that we are in the METRICS category before updating subTabVal
        if (tabVal === METRICS) {
          handleChangeSelectedTabCustomCategory(tabVal, newVal);
        } else {
          handleChangeSelectedTab(newVal);
        }
        setSubTabVal(newVal);
      }
    };
  };

  let backToPlay;
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
    <>
      {CAN(keys.VIEW_SETTINGS.action, keys.VIEW_SETTINGS.subject) ? (
        <div className={classes.wrapperClss}>
          <Paper square className={classes.wrapperClss}>
            <Tabs
              value={tabVal}
              className={classes.tabs}
              onChange={handleChange('tabVal')}
              variant={window.innerWidth < 900 ? 'scrollable' : 'fullWidth'}
              scrollButtons="on"
              indicatorColor="primary"
              textColor="primary"
              centered
            >
              <CustomTooltip title="Connect Meshery Adapters" placement="top" value={ADAPTERS}>
                <Tab
                  className={classes.tab}
                  icon={<FontAwesomeIcon icon={faMendeley} style={iconMedium} />}
                  label="Adapters"
                  data-cy="tabServiceMeshes"
                  value={ADAPTERS}
                  disabled={!CAN(keys.VIEW_SERVICE_MESH.action, keys.VIEW_SERVICE_MESH.subject)}
                />
              </CustomTooltip>
              <CustomTooltip title="Configure Metrics backends" placement="top" value={METRICS}>
                <Tab
                  className={classes.tab}
                  icon={<FontAwesomeIcon icon={faPoll} style={iconMedium} />}
                  label="Metrics"
                  // tab="tabMetrics"
                  value={METRICS}
                  disabled={!CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject)}
                />
              </CustomTooltip>
              <CustomTooltip title="Registry" placement="top" value={REGISTRY}>
                <Tab
                  className={classes.tab}
                  icon={<FontAwesomeIcon icon={faFileInvoice} style={iconMedium} />}
                  label="Registry"
                  // tab="registry"
                  value={REGISTRY}
                  disabled={!CAN(keys.VIEW_REGISTRY.action, keys.VIEW_REGISTRY.subject)}
                />
              </CustomTooltip>
              <CustomTooltip title="Reset System" placement="top" value={RESET}>
                <Tab
                  className={classes.tab}
                  icon={<FontAwesomeIcon icon={faDatabase} style={iconMedium} />}
                  label="Reset"
                  // tab="systemReset"
                  value={RESET}
                  // disabled={!CAN(keys.VIEW_SYSTEM_RESET.action, keys.VIEW_SYSTEM_RESET.subject)} TODO: uncomment when key get seeded
                />
              </CustomTooltip>
            </Tabs>
          </Paper>
          {tabVal === ADAPTERS &&
            CAN(keys.VIEW_SERVICE_MESH.action, keys.VIEW_SERVICE_MESH.subject) && (
              <TabContainer>
                <MeshAdapterConfigComponent />
              </TabContainer>
            )}
          {tabVal === METRICS && CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject) && (
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
                    value={GRAFANA}
                    label={
                      <div className={classes.iconText}>
                        Grafana
                        <img src="/static/img/grafana_icon.svg" className={classes.icon} />
                      </div>
                    }
                  />
                  <Tab
                    className={classes.tab}
                    value={PROMETHEUS}
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
              {subTabVal === GRAFANA && (
                <TabContainer>
                  <GrafanaComponent
                    scannedGrafana={scannedGrafana}
                    isMeshConfigured={isMeshConfigured}
                  />
                </TabContainer>
              )}
              {subTabVal === PROMETHEUS && (
                <TabContainer>
                  <PrometheusComponent
                    scannedPrometheus={scannedPrometheus}
                    isMeshConfigured={isMeshConfigured}
                  />
                </TabContainer>
              )}
            </TabContainer>
          )}
          {tabVal === REGISTRY && CAN(keys.VIEW_REGISTRY.action, keys.VIEW_REGISTRY.subject) && (
            <TabContainer>
              <TabContainer>
                <TabContainer>
                  <MeshModelComponent
                    modelsCount={modelsCount}
                    componentsCount={componentsCount}
                    relationshipsCount={relationshipsCount}
                    registrantCount={registrantCount}
                    settingsRouter={settingsRouter}
                  />
                </TabContainer>
              </TabContainer>
              {/* </div> */}
            </TabContainer>
          )}
          {tabVal === RESET && (
            <TabContainer>
              <DatabaseSummary promptRef={systemResetPromptRef} />
            </TabContainer>
          )}
          {backToPlay}
          <PromptComponent ref={systemResetPromptRef} />
        </div>
      ) : (
        <DefaultError />
      )}
    </>
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
  connect(mapStateToProps, mapDispatchToProps)(withRouter(withNotify(MesherySettings))),
);
