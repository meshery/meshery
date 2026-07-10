import React, { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/router';
import { NoSsr } from '@sistent/sistent';
import {
  CustomTooltip,
  Typography,
  styled,
  Tabs,
  Tab,
  Paper,
  Grid2,
  LeftArrowIcon,
  DatabaseIcon,
  MendeleyIcon,
  FileIcon,
  SettingsIcon,
  useTheme,
} from '@sistent/sistent';
import DashboardMeshModelGraph from '../dashboard/charts/DashboardMeshModelGraph';
import Link from 'next/link';
import MeshAdapterConfigComponent from '../MeshAdapterConfigComponent';
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
import { Keys } from '@meshery/schemas/permissions';
import { ADAPTERS, RESET, OVERVIEW, REGISTRY, CONTROLLERS } from '@/constants/navigator';
import MesheryControllersConfig from './MesheryControllersConfig';
import { removeDuplicateVersions } from '../registry/helper';
import MeshModelComponent from '../registry/MeshModelComponent';
import DefaultError from '../general/error-404';
import MesheryConfigurationChart from '../dashboard/charts/MesheryConfigurationCharts';
import ConnectionStatsChart from '../dashboard/charts/ConnectionCharts';
import { useSelector } from 'react-redux';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';

const StyledPaper = styled(Paper)(() => ({
  flexGrow: 1,
  maxWidth: '100%',
  height: 'auto',
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
  selectedSettingsCategory?: string;
  selectedTab?: string;
  handleChangeSettingsCategory: (_settingsCategory?: string) => void;
  handleChangeSelectedTab: (_tab: string) => void;
  handleChangeSelectedTabCustomCategory: (_settingsCategory: string, _tab: string) => void;
};

const settingsRouter = (router: ReturnType<typeof useRouter>): SettingsRouter => {
  const { query, push: pushRoute, route } = router;

  const selectedSettingsCategory = query.settingsCategory;
  const selectedTab = query.tab;

  const handleChangeSettingsCategory = (_settingsCategory) => {
    if (query.settingsCategory === _settingsCategory) {
      return;
    }
    pushRoute(
      `${route}?settingsCategory=${_settingsCategory || query.settingsCategory}`,
      undefined,
      { shallow: true },
    );
  };

  const handleChangeSelectedTab = (_tab) => {
    if (query.tab === _tab) {
      return;
    }
    pushRoute(`${route}?settingsCategory=${selectedSettingsCategory}&tab=${_tab}`, undefined, {
      shallow: true,
    });
  };

  const handleChangeSelectedTabCustomCategory = (_settingsCategory, _tab) => {
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
  const { selectedSettingsCategory } = settingsRouter(router);
  const theme = useTheme();
  const { k8sConfig } = useSelector((state) => state.ui);
  const { meshAdapters } = useSelector((state) => state.adapter);
  const { data: providerCapabilities } = useGetProviderCapabilitiesQuery();
  const [state, setState] = useState({
    meshAdapters,
    tabVal: selectedSettingsCategory || OVERVIEW,
    modelsCount: 0,
    componentsCount: 0,
    relationshipsCount: 0,
    registrantCount: 0,
    isMeshConfigured: k8sConfig.clusterConfigured,
  });

  const systemResetPromptRef = useRef<{ show: (_args: any) => Promise<string> } | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const modelsResponse = await getMeshModels();
        const componentsResponse = await getComponentsDetail();
        const relationshipsResponse = await getRelationshipsDetail();
        const registrantResponce = await getMeshModelRegistrants();

        const modelsCount = removeDuplicateVersions(modelsResponse.models).length;
        const componentsCount = componentsResponse.totalCount;
        const relationshipsCount = relationshipsResponse.totalCount;
        const registrantCount = registrantResponce.totalCount;

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
    const { handleChangeSettingsCategory } = settingsRouter(router);

    return (_event, newVal, ..._args) => {
      if (val === 'tabVal') {
        handleChangeSettingsCategory(newVal);
        setState((prevState) => ({
          ...prevState,
          tabVal: newVal,
        }));
      }
    };
  };

  const { tabVal } = state;
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
      {CAN(Keys.MesherySystemViewSettings.id, Keys.MesherySystemViewSettings.function) ? (
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
                    // disabled={!CAN(Keys.VIEW_OVERVIEW.id, Keys.VIEW_OVERVIEW.function)}
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
                        Keys.InfrastructureManagementViewCloudNativeInfrastructure.id,
                        Keys.InfrastructureManagementViewCloudNativeInfrastructure.function,
                      )
                    }
                  />
                </CustomTooltip>
                <CustomTooltip title="Registry" placement="top" value={REGISTRY}>
                  <Tab
                    icon={<FileIcon {...iconMedium} fill={theme.palette.icon.default} />}
                    label="Registry"
                    data-testid="settings-tab-registry"
                    value={REGISTRY}
                    disabled={
                      !CAN(
                        Keys.MesherySystemViewRegistry.id,
                        Keys.MesherySystemViewRegistry.function,
                      )
                    }
                  />
                </CustomTooltip>

                <CustomTooltip
                  title="Meshery Operator, MeshSync & Broker defaults"
                  placement="top"
                  value={CONTROLLERS}
                >
                  <Tab
                    icon={<SettingsIcon {...iconMedium} fill={theme.palette.icon.default} />}
                    label="Operator, MeshSync & Broker"
                    data-testid="settings-tab-controllers"
                    value={CONTROLLERS}
                  />
                </CustomTooltip>

                <CustomTooltip title="Reset System" placement="top" value={RESET}>
                  <Tab
                    icon={<DatabaseIcon {...iconMedium} fill={theme.palette.icon.default} />}
                    label="Reset"
                    data-testid="settings-tab-reset"
                    // tab="systemReset"
                    value={RESET}
                    // disabled={!CAN(Keys.VIEW_SYSTEM_RESET.id, Keys.VIEW_SYSTEM_RESET.function)} TODO: uncomment when key get seeded
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
                          {providerCapabilities?.providerName || 'Not selected'}
                          {providerCapabilities?.providerType
                            ? ` (${providerCapabilities.providerType})`
                            : ''}
                        </Typography>
                      </div>
                      <Typography
                        variant="body2"
                        component="a"
                        href="https://docs.meshery.io/extensibility/providers"
                        target="_blank"
                        rel="noopener noreferrer"
                        sx={{ color: theme.palette.primary.main }}
                      >
                        Learn about providers
                      </Typography>
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
                Keys.InfrastructureManagementViewCloudNativeInfrastructure.id,
                Keys.InfrastructureManagementViewCloudNativeInfrastructure.function,
              ) && (
                <TabContainer>
                  <MeshAdapterConfigComponent />
                </TabContainer>
              )}
            {tabVal === REGISTRY && (
              <TabContainer>
                <MeshModelComponent settingsRouter={settingsRouter} />
              </TabContainer>
            )}

            {tabVal === CONTROLLERS && (
              <TabContainer>
                <MesheryControllersConfig />
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
