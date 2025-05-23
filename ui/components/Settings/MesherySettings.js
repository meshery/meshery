import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { useRouter } from 'next/router';
import { NoSsr } from '@layer5/sistent';
import {
  CustomTooltip,
  AppBar,
  Typography,
  styled,
  Tabs,
  Tab,
  Paper,
  Grid,
  LeftArrowIcon,
  PollIcon,
  DatabaseIcon,
  FileIcon,
  MendeleyIcon,
  useTheme,
} from '@layer5/sistent';
import DashboardMeshModelGraph from '../Dashboard/charts/DashboardMeshModelGraph';
import Link from 'next/link';
import GrafanaComponent from '../telemetry/grafana/GrafanaComponent';
import MeshAdapterConfigComponent from '../MeshAdapterConfigComponent';
import PrometheusComponent from '../telemetry/prometheus/PrometheusComponent';
import _PromptComponent from '../PromptComponent';
import { iconMedium } from '../../css/icons.styles';
import MeshModelComponent from './Registry/MeshModelComponent';
import DatabaseSummary from '../DatabaseSummary';
import {
  getComponentsDetail,
  getMeshModels,
  getRelationshipsDetail,
  getMeshModelRegistrants,
} from '../../api/meshmodel';
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
import { removeDuplicateVersions } from './Registry/helper';
import DefaultError from '../General/error-404';
import MesheryConfigurationChart from '../Dashboard/charts/MesheryConfigurationCharts';
import ConnectionStatsChart from '../Dashboard/charts/ConnectionCharts';
import { SecondaryTab, SecondaryTabs } from '../Dashboard/style';
import { useSelector } from 'react-redux';

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
const MesherySettings = () => {
  const router = useRouter();
  const { selectedSettingsCategory, selectedTab } = settingsRouter(router);
  const theme = useTheme();
  const { k8sConfig } = useSelector((state) => state.ui);
  const { prometheus } = useSelector((state) => state.telemetry);
  const { grafana } = useSelector((state) => state.telemetry);
  const { meshAdapters } = useSelector((state) => state.adapter);
  const [state, setState] = useState({
    meshAdapters,
    grafana,
    prometheus,
    tabVal: selectedSettingsCategory || OVERVIEW,
    subTabVal: selectedTab || GRAFANA,
    modelsCount: 0,
    componentsCount: 0,
    relationshipsCount: 0,
    registrantCount: 0,
    isMeshConfigured: k8sConfig.clusterConfigured,
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
    } = settingsRouter(router);

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
  if (k8sConfig.clusterConfigured === true && meshAdapters.length > 0) {
    backToPlay = (
      <div>
        <Link href="/management">
          <div>
            <LeftArrowIcon transform="grow-4" />
            You are ready to manage cloud native infrastructure
          </div>
        </Link>
      </div>
    );
  }
  return (
    <>
      {CAN(keys.VIEW_SETTINGS.action, keys.VIEW_SETTINGS.subject) ? (
        <>
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
                    icon={<MendeleyIcon {...iconMedium} fill={theme.palette.icon.default} />}
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
                    icon={<PollIcon {...iconMedium} fill={theme.palette.icon.default} />}
                    label="Metrics"
                    data-testid="settings-tab-metrics"
                    // tab="tabMetrics"
                    value={METRICS}
                    disabled={!CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject)}
                  />
                </CustomTooltip>
                <CustomTooltip title="Registry" placement="top" value={REGISTRY}>
                  <Tab
                    icon={<FileIcon {...iconMedium} fill={theme.palette.icon.default} />}
                    label="Registry"
                    data-testid="settings-tab-registry"
                    // tab="registry"
                    value={REGISTRY}
                    disabled={!CAN(keys.VIEW_REGISTRY.action, keys.VIEW_REGISTRY.subject)}
                  />
                </CustomTooltip>

                <CustomTooltip title="Reset System" placement="top" value={RESET}>
                  <Tab
                    icon={<DatabaseIcon {...iconMedium} fill={theme.palette.icon.default} />}
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
        </>
      ) : (
        <DefaultError />
      )}
    </>
  );
};

export default MesherySettings;
