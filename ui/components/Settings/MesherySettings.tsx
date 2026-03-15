import React, { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/router';
import { NoSsr } from '@sistent/sistent';
import {
  CustomTooltip,
  AppBar,
  Typography,
  styled,
  Tabs,
  Tab,
  Paper,
  Grid2,
  LeftArrowIcon,
  PollIcon,
  DatabaseIcon,
  MendeleyIcon,
  FileIcon,
  useTheme,
  Box,
} from '@sistent/sistent';
import DashboardMeshModelGraph from '../Dashboard/charts/DashboardMeshModelGraph';
import Link from 'next/link';
import GrafanaComponent from '../telemetry/grafana/GrafanaComponent';
import MeshAdapterConfigComponent from '../MeshAdapterConfigComponent';
import PrometheusComponent from '../telemetry/prometheus/PrometheusComponent';
import _PromptComponent from '../PromptComponent';
import { iconMedium } from '../../css/icons.styles';
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
  METRICS,
  ADAPTERS,
  RESET,
  GRAFANA,
  PROMETHEUS,
  OVERVIEW,
  REGISTRY,
} from '@/constants/navigator';
import { removeDuplicateVersions } from './Registry/helper';
import MeshModelComponent from './Registry/MeshModelComponent';
import DefaultError from '../General/error-404';
import MesheryConfigurationChart from '../Dashboard/charts/MesheryConfigurationCharts';
import ConnectionStatsChart from '../Dashboard/charts/ConnectionCharts';
import { SecondaryTab, SecondaryTabs } from '../Dashboard/style';
import { useSelector } from 'react-redux';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import type { RootState } from '@/store/index';

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

type TabContainerProps = { children: React.ReactNode };

function TabContainer({ children }: TabContainerProps) {
  return (
    <Typography component="div" style={{ paddingTop: 2 }}>
      {children}
    </Typography>
  );
}

type SettingsRouter = {
  selectedSettingsCategory?: string | undefined;
  selectedTab?: string | undefined;
  handleChangeSettingsCategory: (_settingsCategory?: string) => void;
  handleChangeSelectedTab: (_tab: string) => void;
  handleChangeSelectedTabCustomCategory: (_settingsCategory: string, _tab: string) => void;
};

const settingsRouter = (router: ReturnType<typeof useRouter>): SettingsRouter => {
  const { query, push: pushRoute, route } = router;

  const selectedSettingsCategory = Array.isArray(query.settingsCategory)
    ? query.settingsCategory[0]
    : query.settingsCategory;
  const selectedTab = Array.isArray(query.tab) ? query.tab[0] : query.tab;

  const handleChangeSettingsCategory = (_settingsCategory?: string) => {
    if (query.settingsCategory === _settingsCategory) {
      return;
    }
    pushRoute(
      `${route}?settingsCategory=${_settingsCategory || query.settingsCategory}`,
      undefined,
      { shallow: true },
    );
  };

  const handleChangeSelectedTab = (_tab: string) => {
    if (query.tab === _tab) {
      return;
    }
    pushRoute(`${route}?settingsCategory=${selectedSettingsCategory}&tab=${_tab}`, undefined, {
      shallow: true,
    });
  };

  const handleChangeSelectedTabCustomCategory = (_settingsCategory: string, _tab: string) => {
    if (query.tab === _tab) {
      return;
    }
    pushRoute(`${route}?settingsCategory=${_settingsCategory}&tab=${_tab}`, undefined, {
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
  const { k8sConfig } = useSelector((state: RootState) => state.ui);
  const { prometheus } = useSelector((state: RootState) => state.telemetry);
  const { grafana } = useSelector((state: RootState) => state.telemetry);
  const { meshAdapters } = useSelector((state: RootState) => state.adapter);
  const { data: providerCapabilities } = useGetProviderCapabilitiesQuery(undefined);
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
    isMeshConfigured: k8sConfig.length > 0,
    scannedPrometheus: [],
    scannedGrafana: [],
  });

  const systemResetPromptRef = useRef<{ show: (_args: any) => Promise<string> } | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const modelsResponse = (await getMeshModels()) as { models: any[] };
        const componentsResponse = (await getComponentsDetail(1)) as { total_count: number };
        const relationshipsResponse = (await getRelationshipsDetail(1)) as { total_count: number };
        const registrantResponce = (await getMeshModelRegistrants()) as { total_count: number };

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

  const handleChange = (val: string) => {
    const {
      handleChangeSettingsCategory,
      handleChangeSelectedTab,
      handleChangeSelectedTabCustomCategory,
    } = settingsRouter(router);

    return (_event: React.SyntheticEvent, newVal: string, ..._args: any[]) => {
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
  let backToPlay: React.ReactNode = null;
  if (k8sConfig.length > 0 && meshAdapters.length > 0) {
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
          <Box sx={{ flexGrow: 1, maxWidth: '100%', height: 'auto' }}>
            <StyledPaper square>
              <Tabs
                value={tabVal}
                onChange={handleChange('tabVal')}
                variant={window.innerWidth < 900 ? 'scrollable' : 'fullWidth'}
                scrollButtons={true}
                indicatorColor="primary"
                textColor="primary"
                centered
              >
                <CustomTooltip title="Overview" placement="top">
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
                <CustomTooltip title="Configure Metrics backends" placement="top">
                  <Tab
                    icon={<PollIcon {...iconMedium} fill={theme.palette.icon.default} />}
                    label="Metrics"
                    data-testid="settings-tab-metrics"
                    // tab="tabMetrics"
                    value={METRICS}
                    disabled={!CAN(keys.VIEW_METRICS.action, keys.VIEW_METRICS.subject)}
                  />
                </CustomTooltip>

                <CustomTooltip title="Registry" placement="top">
                  <Tab
                    icon={<FileIcon {...iconMedium} fill={theme.palette.icon.default} />}
                    label="Registry"
                    data-testid="settings-tab-registry"
                    value={REGISTRY}
                    disabled={!CAN(keys.VIEW_REGISTRY.action, keys.VIEW_REGISTRY.subject)}
                  />
                </CustomTooltip>

                <CustomTooltip title="Reset System" placement="top">
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
                    <Paper
                      sx={{
                        padding: '1rem 1.5rem',
                        marginBottom: '1rem',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between',
                      }}
                    >
                      <div>
                        <Typography variant="subtitle1" sx={{ fontWeight: 500 }}>
                          Current Provider
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {providerCapabilities?.provider_name || 'Not selected'}
                          {providerCapabilities?.provider_type
                            ? ` (${providerCapabilities.provider_type})`
                            : ''}
                        </Typography>
                      </div>
                      <a
                        href="https://docs.meshery.io/extensibility/providers"
                        target="_blank"
                        rel="noopener noreferrer"
                        style={{
                          color: theme.palette.primary.main,
                          fontSize: '0.875rem',
                          textDecoration: 'none',
                        }}
                      >
                        Learn about providers
                      </a>
                    </Paper>
                    <DashboardMeshModelGraph />
                    <Grid2 container spacing={2} size="grow">
                      <Grid2 size={{ xs: 12, md: 6 }}>
                        <ConnectionStatsChart />
                      </Grid2>
                      <Grid2 size={{ xs: 12, md: 6 }}>
                        <MesheryConfigurationChart />
                      </Grid2>
                    </Grid2>
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

            {tabVal === REGISTRY && (
              <TabContainer>
                <MeshModelComponent settingsRouter={settingsRouter} />
              </TabContainer>
            )}

            {tabVal === RESET && (
              <TabContainer>
                <DatabaseSummary promptRef={systemResetPromptRef} />
              </TabContainer>
            )}
            {backToPlay}
            <_PromptComponent ref={systemResetPromptRef} />
          </Box>
        </>
      ) : (
        <DefaultError />
      )}
    </>
  );
};

export default MesherySettings;
