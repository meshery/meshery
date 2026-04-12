// @ts-nocheck
import React, { useState } from 'react';
import { Tab, Tabs, CustomTooltip, NoSsr, TachometerIcon, useTheme } from '@sistent/sistent';
import SettingsRemoteIcon from '@mui/icons-material/SettingsRemote';
import SettingsCellIcon from '@mui/icons-material/SettingsCell';
import MesherySettingsPerformanceComponent from '../Settings/MesherySettingsPerformanceComponent';
import { iconMedium } from '../../css/icons.styles';
import { useWindowDimensions } from '@/utils/dimension';
import { useUserPreferences } from './hooks/useUserPreferences';
import GeneralTab from './GeneralTab';
import RemoteProviderTab from './RemoteProviderTab';
import { StatsWrapper, StyledPaper, TabLabel, IconStyled } from './style';

interface UserPreferenceProps {
  anonymousStats?: boolean;
  perfResultStats?: boolean;
  classes?: Record<string, string>;
}

const UserPreference: React.FC<UserPreferenceProps> = (props) => {
  const [tabVal, setTabVal] = useState(0);
  const { width } = useWindowDimensions();
  const theme = useTheme();

  const {
    anonymousStats,
    perfResultStats,
    catalogContent,
    userPrefs,
    providerType,
    providerInfo,
    userData,
    handleCatalogContentToggle,
    handleToggle,
    handleUpdateUserPref,
  } = useUserPreferences(props.anonymousStats, props.perfResultStats);

  const handleTabValChange = (_event: React.SyntheticEvent, newVal: number) => {
    setTabVal(newVal);
  };

  return (
    <>
      <NoSsr>
        <StyledPaper>
          <Tabs
            value={tabVal}
            onChange={handleTabValChange}
            variant={width < 600 ? 'scrollable' : 'fullWidth'}
            scrollButtons="on"
            allowScrollButtonsMobile={true}
            indicatorColor="primary"
            textColor="primary"
            centered
          >
            <CustomTooltip title="General preferences" placement="top">
              <Tab
                icon={<IconStyled as={SettingsCellIcon} />}
                label={<TabLabel>General</TabLabel>}
              />
            </CustomTooltip>
            <CustomTooltip title="Choose Performance Test Defaults" placement="top">
              <Tab
                icon={<TachometerIcon {...iconMedium} fill={theme.palette.icon.default} />}
                label={<TabLabel>Performance</TabLabel>}
              />
            </CustomTooltip>
            {userPrefs && providerType != 'local' && (
              <CustomTooltip title="Remote Provider preferences" placement="top">
                <Tab
                  icon={<SettingsRemoteIcon style={iconMedium} />}
                  label={<TabLabel>Remote Provider</TabLabel>}
                />
              </CustomTooltip>
            )}
          </Tabs>
        </StyledPaper>
        <StatsWrapper>
          {tabVal === 0 && (
            <GeneralTab
              catalogContent={catalogContent}
              anonymousStats={anonymousStats}
              perfResultStats={perfResultStats}
              handleCatalogContentToggle={handleCatalogContentToggle}
              handleToggle={handleToggle}
              handleUpdateUserPref={handleUpdateUserPref}
            />
          )}
          {tabVal === 1 && <MesherySettingsPerformanceComponent />}
          {tabVal === 2 && userPrefs && providerType !== 'local' && (
            <RemoteProviderTab
              userPrefs={userPrefs}
              providerInfo={providerInfo}
              userData={userData}
            />
          )}
        </StatsWrapper>
      </NoSsr>
    </>
  );
};

export default UserPreference;
