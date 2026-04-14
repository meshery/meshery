import React from 'react';
import { FormGroup, FormControlLabel, Switch } from '@sistent/sistent';
import { ThemeTogglerCore } from '@/themes/hooks';
import { FormContainerWrapper, FormGroupWrapper, FormLegend } from './style';

interface ThemeComponentProps {
  mode: 'light' | 'dark';
  toggleTheme: () => void;
}

const ThemeComponent: React.FC<ThemeComponentProps> = ({ mode, toggleTheme }) => (
  <div>
    <Switch color="primary" checked={mode === 'dark'} onChange={toggleTheme} />
    Dark Mode
  </div>
);

const ThemeToggler: React.FC = () => <ThemeTogglerCore Component={ThemeComponent} />;

interface GeneralTabProps {
  catalogContent: boolean;
  anonymousStats?: boolean;
  perfResultStats?: boolean;
  handleCatalogContentToggle: () => void;
  handleToggle: (name: string) => () => void;
}

const GeneralTab: React.FC<GeneralTabProps> = ({
  catalogContent,
  anonymousStats,
  perfResultStats,
  handleCatalogContentToggle,
  handleToggle,
}) => {
  return (
    <>
      <FormContainerWrapper>
        <FormGroupWrapper component={'fieldset' as React.ElementType}>
          <FormLegend>Extensions</FormLegend>
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
        <FormGroupWrapper component={'fieldset' as React.ElementType}>
          <FormLegend>Analytics and Improvement Program</FormLegend>
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
        <FormGroupWrapper component={'fieldset' as React.ElementType}>
          <FormLegend>Theme</FormLegend>
          <FormGroup>
            <ThemeToggler />
          </FormGroup>
        </FormGroupWrapper>
      </FormContainerWrapper>
    </>
  );
};

export default GeneralTab;
