// @ts-nocheck
import React from 'react';
import { FormGroup, FormControlLabel, Switch } from '@sistent/sistent';
import { ThemeTogglerCore } from '@/themes/hooks';
import { FormContainerWrapper, FormGroupWrapper, FormLegend } from './style';

interface ThemeTogglerProps {
  handleUpdateUserPref: (key: string, value: string) => void;
}

interface ThemeComponentProps {
  mode: 'light' | 'dark';
  toggleTheme: () => void;
}

const ThemeToggler: React.FC<ThemeTogglerProps> = ({ handleUpdateUserPref }) => {
  const Component: React.FC<ThemeComponentProps> = ({ mode, toggleTheme }) => {
    return (
      <div>
        <Switch
          color="primary"
          checked={mode === 'dark'}
          onChange={() => {
            toggleTheme();
            handleUpdateUserPref(
              'remoteProviderPreferences.theme',
              mode === 'dark' ? 'light' : 'dark',
            );
          }}
        />
        Dark Mode
      </div>
    );
  };

  return <ThemeTogglerCore Component={Component}></ThemeTogglerCore>;
};

interface GeneralTabProps {
  catalogContent: boolean;
  anonymousStats?: boolean;
  perfResultStats?: boolean;
  handleCatalogContentToggle: () => void;
  handleToggle: (name: string) => () => void;
  handleUpdateUserPref: (key: string, value: unknown) => void;
}

const GeneralTab: React.FC<GeneralTabProps> = ({
  catalogContent,
  anonymousStats,
  perfResultStats,
  handleCatalogContentToggle,
  handleToggle,
  handleUpdateUserPref,
}) => {
  return (
    <>
      <FormContainerWrapper>
        <FormGroupWrapper component="fieldset">
          <FormLegend component="legend">Extensions</FormLegend>
          <FormGroup>
            <FormControlLabel
              key="CatalogContentPreference"
              control={
                <Switch
                  checked={catalogContent}
                  onChange={handleCatalogContentToggle}
                  color="primary"
                  data-cy="CatalogContentPreference"
                />
              }
              labelPlacement="end"
              label="Meshery Catalog Content"
            />
          </FormGroup>
        </FormGroupWrapper>
      </FormContainerWrapper>
      <FormContainerWrapper>
        <FormGroupWrapper component="fieldset">
          <FormLegend component="legend">Analytics and Improvement Program</FormLegend>
          <FormGroup>
            <FormControlLabel
              key="UsageStatsPreference"
              control={
                <Switch
                  checked={anonymousStats}
                  onChange={handleToggle('anonymousUsageStats')}
                  color="primary"
                  data-cy="UsageStatsPreference"
                />
              }
              labelPlacement="end"
              label="Send Anonymous Usage Statistics"
            />
            <FormControlLabel
              key="PerfResultPreference"
              control={
                <Switch
                  checked={perfResultStats}
                  onChange={handleToggle('anonymousPerfResults')}
                  color="primary"
                  data-cy="PerfResultPreference"
                />
              }
              labelPlacement="end"
              label="Send Anonymous Performance Results"
            />
          </FormGroup>
        </FormGroupWrapper>
      </FormContainerWrapper>
      <FormContainerWrapper>
        <FormGroupWrapper component="fieldset">
          <FormLegend component="legend">Theme</FormLegend>
          <FormGroup>
            <ThemeToggler handleUpdateUserPref={handleUpdateUserPref} />
          </FormGroup>
        </FormGroupWrapper>
      </FormContainerWrapper>
    </>
  );
};

export default GeneralTab;
