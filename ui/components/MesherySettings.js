import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'next/router';
import { connect, Provider } from 'react-redux';
import { NoSsr } from '@material-ui/core';
import { bindActionCreators } from 'redux';
import {
  CustomTooltip,
  AppBar,
  // Paper,
  Typography,
  styled,
  Tabs,
  Tab,
  Paper,
} from '@layer5/sistent';
import DashboardMeshModelGraph from './DashboardComponent/charts/DashboardMeshModelGraph';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faPoll, faDatabase, faFileInvoice } from '@fortawesome/free-solid-svg-icons';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import Link from 'next/link';
import GrafanaComponent from './telemetry/grafana/GrafanaComponent';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import PrometheusComponent from './telemetry/prometheus/PrometheusComponent';
import { updateProgress } from '../lib/store';
import _PromptComponent from './PromptComponent';
import { iconMedium } from '../css/icons.styles';
import MeshModelComponent from './MeshModelRegistry/MeshModelComponent';
import DatabaseSummary from './DatabaseSummary';
import {
  getComponentsDetail,
  getMeshModels,
  getRelationshipsDetail,
  getMeshModelRegistrants,
} from '../api/meshmodel';
import { withNotify } from '../utils/hooks/useNotification';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import {
  REGISTRY,
  METRICS,
  ADAPTERS,
  RESET,
  GRAFANA,
  PROMETHEUS,
  OVERVIEW,
} from '@/constants/navigator';
import Grid from '@material-ui/core/Grid';
import { removeDuplicateVersions } from './MeshModelRegistry/helper';
import DefaultError from './General/error-404';
import { store } from '../store';
import MesheryConfigurationChart from './DashboardComponent/charts/MesheryConfigurationCharts';
import ConnectionStatsChart from './DashboardComponent/charts/ConnectionCharts';
import { UsesSistent } from './SistentWrapper';
import { SecondaryTab, SecondaryTabs } from './DashboardComponent/style';

const StyledPaper = styled(Paper)(() => ({
  flexGrow: 1,
  maxWidth: '100%',
  height: 'auto',
}));

const IconText = styled('div')(() => ({
  display: 'inline',
  verticalAlign: 'middle',
}));

const StyledIcon = styled('img')(({ theme }) => ({
  display: 'inline',
  verticalAlign: 'text-top',
  width: theme.spacing(1.75),
  marginLeft: theme.spacing(0.5),
}));

const RootClass = styled('div')(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? '#303030' : '#eaeff1',
  marginTop: '1rem',
}));

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
const MesherySettings = (props) => {
  const { k8sconfig, meshAdapters, grafana, prometheus, router, classes } = props;
  const { selectedSettingsCategory, selectedTab } = settingsRouter(router);

  const [state, setState] = useState({
    k8sconfig,
    meshAdapters,
    grafana,
    prometheus,
    tabVal: selectedSettingsCategory || OVERVIEW,
    subTabVal: selectedTab || GRAFANA,
    modelsCount: 0,
    componentsCount: 0,
    relationshipsCount: 0,
    registrantCount: 0,
    isMeshConfigured: k8sconfig.clusterConfigured,
    scannedPrometheus: [],
    scannedGrafana: [],
  });

  const systemResetPromptRef = useRef(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const modelsResponse = await getMeshModels();
        const componentsResponse = await getComponentsDetail();
        const relationshipsResponse = await getRelationshipsDetail();
        const registrantResponce = await getMeshModelRegistrants();

        const modelsCount = removeDuplicateVersions(modelsResponse.models).length;
        const componentsCount = componentsResponse.total_count;
        const relationshipsCount = relationshipsResponse.total_count;
        const registrantCount = registrantResponce.total_count;

        setState((prevState) => ({
          ...prevState,
          modelsCount,
          componentsCount,
          relationshipsCount,
          registrantCount,
        }));
      } catch (error) {
        console.error(error);
      }
    };
    fetchData();
  }, []);

  useEffect(() => {
    const { selectedSettingsCategory } = settingsRouter(router);
    if (selectedSettingsCategory && selectedSettingsCategory !== state.tabVal) {
      setState((prevState) => ({
        ...prevState,
        tabVal: selectedSettingsCategory,
      }));
    }
  }, [router]);

  const handleChange = (val) => {
    const {
      handleChangeSettingsCategory,
      handleChangeSelectedTab,
      handleChangeSelectedTabCustomCategory,
    } = settingsRouter(props.router);

    return (event, newVal) => {
      if (val === 'tabVal') {
        if (newVal === METRICS) {
          handleChangeSelectedTabCustomCategory(newVal, GRAFANA);
          setState((prevState) => ({
            ...prevState,
            tabVal: newVal,
            subTabVal: GRAFANA,
          }));
        } else {
          handleChangeSettingsCategory(newVal);
          setState((prevState) => ({
            ...prevState,
            tabVal: newVal,
          }));
        }
      } else if (val === 'subTabVal') {
        handleChangeSelectedTab(newVal);
        setState((prevState) => ({
          ...prevState,
          subTabVal: newVal,
        }));
      }
    };
  };

  const { tabVal, subTabVal } = state;
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
    <>
      {CAN(keys.VIEW_SETTINGS.action, keys.VIEW_SETTINGS.subject) ? (
        <UsesSistent>
          <div sx={{ flexGrow: 1, maxWidth: '100%', height: 'auto' }}>
            <StyledPaper square>
              <Tabs
                value={tabVal}
                onChange={handleChange('tabVal')}
                variant={window.innerWidth < 900 ? 'scrollable' : 'fullWidth'}
                scrollButtons="on"
                indicatorColor="primary"
                textColor="primary"
                centered
              >
                <CustomTooltip title="Overview" placement="top" value={OVERVIEW}>
                  <Tab
                    icon={
                      <img
                        src="/static/img/meshery-logo/meshery-logo.svg"
                        alt="Meshery logo"
                        height={32}
                        width={32}
                      />
                    }
                    label="Overview"
                    // tab="Overview"
                    value={OVERVIEW}
                    // disabled={!CAN(keys.VIEW_OVERVIEW.action, keys.VIEW_OVERVIEW.subject)}
                  />
                </CustomTooltip>
                <CustomTooltip
                  title="Connect Meshery Adapters"
                  data-testid="settings-tab-adapters"
                  placement="top"
                  value={ADAPTERS}
                >
                  <Tab
                    icon={<FontAwesomeIcon icon={faMendeley} style={iconMedium} />}
                    label="Adapters"
                    data-cy="tabServiceMeshes"
                    value={ADAPTERS}
                    disabled={
                      !CAN(
                        keys.VIEW_CLOUD_NATIVE_INFRASTRUCTURE.action,
                        keys.VIEW_CLOUD_NATIVE_INFRASTRUCTURE.subject,
                      )
                    }
                  />
                </CustomTooltip>
                <CustomTooltip title="Configure Metrics backends" placement="top" value={METRICS}>
                  <Tab
                    icon={<FontAwesomeIcon icon={faPoll} style={iconMedium} />}
                    label="Metrics"
                    data-testid="settings-tab-metrics"
                    // tab="tabMetrics"
                    value={METRICS}
                    disabled={!CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject)}
                  />
                </CustomTooltip>
                <CustomTooltip title="Registry" placement="top" value={REGISTRY}>
                  <Tab
                    icon={<FontAwesomeIcon icon={faFileInvoice} style={iconMedium} />}
                    label="Registry"
                    data-testid="settings-tab-registry"
                    // tab="registry"
                    value={REGISTRY}
                    disabled={!CAN(keys.VIEW_REGISTRY.action, keys.VIEW_REGISTRY.subject)}
                  />
                </CustomTooltip>

                <CustomTooltip title="Reset System" placement="top" value={RESET}>
                  <Tab
                    icon={<FontAwesomeIcon icon={faDatabase} style={iconMedium} />}
                    label="Reset"
                    data-testid="settings-tab-reset"
                    // tab="systemReset"
                    value={RESET}
                    // disabled={!CAN(keys.VIEW_SYSTEM_RESET.action, keys.VIEW_SYSTEM_RESET.subject)} TODO: uncomment when key get seeded
                  />
                </CustomTooltip>
              </Tabs>
            </StyledPaper>

            {tabVal === OVERVIEW && (
              <TabContainer>
                <NoSsr>
                  <Provider store={store}>
                    <RootClass>
                      <DashboardMeshModelGraph />
                      <Grid container spacing={2}>
                        <Grid item xs={12} md={6}>
                          <ConnectionStatsChart />
                        </Grid>
                        <Grid item xs={12} md={6}>
                          <MesheryConfigurationChart />
                        </Grid>
                      </Grid>
                    </RootClass>
                  </Provider>
                </NoSsr>
              </TabContainer>
            )}
            {tabVal === ADAPTERS &&
              CAN(
                keys.VIEW_CLOUD_NATIVE_INFRASTRUCTURE.action,
                keys.VIEW_CLOUD_NATIVE_INFRASTRUCTURE.subject,
              ) && (
                <TabContainer>
                  <MeshAdapterConfigComponent />
                </TabContainer>
              )}
            {tabVal === METRICS && CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject) && (
              <TabContainer>
                <AppBar position="static" color="default">
                  <SecondaryTabs
                    value={subTabVal}
                    onChange={handleChange('subTabVal')}
                    indicatorColor="primary"
                    textColor="primary"
                    variant="fullWidth"
                  >
                    <SecondaryTab
                      value={GRAFANA}
                      label={
                        <IconText>
                          Grafana
                          <StyledIcon src="/static/img/grafana_icon.svg" />
                        </IconText>
                      }
                    />
                    <SecondaryTab
                      value={PROMETHEUS}
                      label={
                        <IconText>
                          Prometheus
                          <StyledIcon src="/static/img/prometheus_logo_orange_circle.svg" />
                        </IconText>
                      }
                    />
                  </SecondaryTabs>
                </AppBar>
                {subTabVal === GRAFANA && (
                  <TabContainer>
                    <GrafanaComponent
                      scannedGrafana={state.scannedGrafana}
                      isMeshConfigured={state.isMeshConfigured}
                    />
                  </TabContainer>
                )}
                {subTabVal === PROMETHEUS && (
                  <TabContainer>
                    <PrometheusComponent
                      scannedPrometheus={state.scannedPrometheus}
                      isMeshConfigured={state.isMeshConfigured}
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
                      modelsCount={state.modelsCount}
                      componentsCount={state.componentsCount}
                      relationshipsCount={state.relationshipsCount}
                      registrantCount={state.registrantCount}
                      settingsRouter={settingsRouter}
                    />
                  </TabContainer>
                </TabContainer>
              </TabContainer>
            )}

            {tabVal === RESET && (
              <TabContainer>
                <DatabaseSummary promptRef={systemResetPromptRef} />
              </TabContainer>
            )}
            {backToPlay}
            <_PromptComponent ref={systemResetPromptRef} />
          </div>
        </UsesSistent>
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

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(withNotify(MesherySettings)));
