import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'next/router';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
  CustomTooltip,
  Tab,
  Tabs,
  AppBar,
  Paper,
  Typography,
  styled,
  useTheme,
} from '@layer5/sistent';
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
  getMeshModels,
  getRelationshipsDetail,
  getMeshModelRegistrants,
} from '../api/meshmodel';
import { withNotify } from '../utils/hooks/useNotification';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { REGISTRY, METRICS, ADAPTERS, RESET, GRAFANA, PROMETHEUS } from '@/constants/navigator';
import { removeDuplicateVersions } from './MeshModelRegistry/helper';
import DefaultError from './General/error-404';

const StyledPaper = styled(Paper)(() => ({
  flexGrow: 1,
  maxWidth: '100%',
  height: 'auto',
}));

const StyledTabs = styled(Tabs)(() => {
  const theme = useTheme();
  return {
    width: '100%',
    '& .MuiTabs-indicator': {
      backgroundColor: theme.palette.mode === 'dark' ? '#00B39F' : theme.palette.primary.main,
    },
  };
});

const StyledTab = styled(Tab)(() => {
  const theme = useTheme();
  return {
    width: '25%',
    minWidth: 40,
    paddingLeft: 0,
    paddingRight: 0,
    '&.Mui-selected': {
      color: theme.palette.mode === 'dark' ? '#00B39F' : theme.palette.primary.main,
    },
  };
});

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
    tabVal: selectedSettingsCategory || ADAPTERS,
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
        <div sx={{ flexGrow: 1, maxWidth: '100%', height: 'auto' }}>
          <StyledPaper square>
            <StyledTabs
              value={tabVal}
              onChange={handleChange('tabVal')}
              variant={window.innerWidth < 900 ? 'scrollable' : 'fullWidth'}
              scrollButtons="on"
              indicatorColor="primary"
              textColor="primary"
              centered
            >
              <CustomTooltip title="Connect Meshery Adapters" placement="top" value={ADAPTERS}>
                <StyledTab
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
                <StyledTab
                  icon={<FontAwesomeIcon icon={faPoll} style={iconMedium} />}
                  label="Metrics"
                  // tab="tabMetrics"
                  value={METRICS}
                  disabled={!CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject)}
                />
              </CustomTooltip>
              <CustomTooltip title="Registry" placement="top" value={REGISTRY}>
                <StyledTab
                  icon={<FontAwesomeIcon icon={faFileInvoice} style={iconMedium} />}
                  label="Registry"
                  // tab="registry"
                  value={REGISTRY}
                  disabled={!CAN(keys.VIEW_REGISTRY.action, keys.VIEW_REGISTRY.subject)}
                />
              </CustomTooltip>
              <CustomTooltip title="Reset System" placement="top" value={RESET}>
                <StyledTab
                  icon={<FontAwesomeIcon icon={faDatabase} style={iconMedium} />}
                  label="Reset"
                  // tab="systemReset"
                  value={RESET}
                  // disabled={!CAN(keys.VIEW_SYSTEM_RESET.action, keys.VIEW_SYSTEM_RESET.subject)} TODO: uncomment when key get seeded
                />
              </CustomTooltip>
            </StyledTabs>
          </StyledPaper>
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
                <StyledTabs
                  value={subTabVal}
                  onChange={handleChange('subTabVal')}
                  indicatorColor="primary"
                  textColor="primary"
                  variant="fullWidth"
                >
                  <StyledTab
                    value={GRAFANA}
                    label={
                      <IconText>
                        Grafana
                        <StyledIcon src="/static/img/grafana_icon.svg" />
                      </IconText>
                    }
                  />
                  <StyledTab
                    value={PROMETHEUS}
                    label={
                      <IconText>
                        Prometheus
                        <StyledIcon src="/static/img/prometheus_logo_orange_circle.svg" />
                      </IconText>
                    }
                  />
                </StyledTabs>
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

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withRouter(withNotify(MesherySettings)));
