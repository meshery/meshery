import React, { useState } from 'react';
import { Tab, Tabs, CustomTooltip, NoSsr, TachometerIcon, useTheme } from '@sistent/sistent';
// eslint-disable-next-line no-restricted-imports
import SettingsCellIcon from '@mui/icons-material/SettingsCell';
import MesherySettingsPerformanceComponent from '../Settings/MesherySettingsPerformanceComponent';
import { iconMedium } from '../../css/icons.styles';
import { useWindowDimensions } from '@/utils/dimension';
import { useUserPreferences } from './hooks/useUserPreferences';
import GeneralTab from './GeneralTab';
import { StatsWrapper, StyledPaper, TabLabel, IconStyled } from './style';

interface UserPreferenceProps {
  anonymousStats?: boolean;
  perfResultStats?: boolean;
  classes?: Record<string, string>;
}

const UserPreference: React.FC<UserPreferenceProps> = (props) => {
  const [tabVal, setTabVal] = useState('general');
  const { width } = useWindowDimensions();
  const theme = useTheme();

  const {
    anonymousStats,
    perfResultStats,
    catalogContent,
    handleCatalogContentToggle,
    handleToggle,
  } = useUserPreferences(props.anonymousStats, props.perfResultStats);

  const handleTabValChange = (_event: React.SyntheticEvent, newVal: string) => {
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
                value="general"
                icon={<IconStyled as={SettingsCellIcon} />}
                label={<TabLabel>General</TabLabel>}
              />
            </CustomTooltip>
            <CustomTooltip title="Choose Performance Test Defaults" placement="top">
              <Tab
                value="performance"
                icon={<TachometerIcon {...iconMedium} fill={theme.palette.icon.default} />}
                label={<TabLabel>Performance</TabLabel>}
              />
            </CustomTooltip>
          </Tabs>
        </StyledPaper>
        <StatsWrapper>
          {tabVal === 'general' && (
            <GeneralTab
              catalogContent={catalogContent}
              anonymousStats={anonymousStats}
              perfResultStats={perfResultStats}
              handleCatalogContentToggle={handleCatalogContentToggle}
              handleToggle={handleToggle}
            />
          )}
          {tabVal === 'performance' && <MesherySettingsPerformanceComponent />}
        </StatsWrapper>
      </NoSsr>
    </>
  );
};

export default UserPreference;
