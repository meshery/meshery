import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@sistent/sistent';
import { MenuItem, TextField, Box, styled } from '@sistent/sistent';
import { connect } from 'react-redux';

const Root = styled(Box)(() => ({
  width: '100%',
}));

const StyledTextField = styled(TextField)(({ theme }) => ({
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(1),
}));

const GrafanaMetricsCompare = ({ chartCompare }) => {
  const [panels, setPanels] = useState<Record<string, any>>({});
  const [panel, setPanel] = useState('');
  const [selectedSeries, setSelectedSeries] = useState('');
  const [series, setSeries] = useState<string[]>([]);

  useEffect(() => {
    const computePanels = (chartCompare) => {
      const panels: Record<string, any> = {};
      if (chartCompare && chartCompare.length > 0) {
        chartCompare.forEach((cc) => {
          if (cc.boardConfig && cc.boardConfig.panels) {
            cc.boardConfig.panels.forEach((panel) => {
              panels[panel.title] = panel;
            });
          }
        });
      }
      return panels;
    };

    const computedPanels = computePanels(chartCompare);
    setPanels(computedPanels);

    const initialPanelKey = Object.keys(computedPanels)[0] ?? '';
    let initialSeries: string[] = [];
    const initialPanel = computedPanels[initialPanelKey];
    if (initialPanel?.targets) {
      initialSeries = initialPanel.targets.map((target: any) => target.expr as string);
    }
    setPanel(initialPanelKey);
    setSeries(initialSeries);
    setSelectedSeries(initialSeries[0] ?? '');
  }, [chartCompare]);

  const handleChange = (name) => (event) => {
    if (name === 'panel') {
      const selectedPanel = String((event as any)?.target?.value ?? '');
      let newSeries: string[] = [];
      const panelObj = panels[selectedPanel];
      if (panelObj?.targets) {
        newSeries = panelObj.targets.map((target: any) => target.expr as string);
      }
      setPanel(selectedPanel);
      setSeries(newSeries);
      setSelectedSeries(newSeries[0] ?? '');
    } else if (name === 'series') {
      setSelectedSeries(String((event as any)?.target?.value ?? ''));
    }
  };

  return (
    <NoSsr>
      <Root>
        <StyledTextField
          select
          fullWidth
          label="Panel"
          value={panel}
          onChange={handleChange('panel')}
          margin="dense"
          variant="outlined"
        >
          {Object.keys(panels).map((p) => (
            <MenuItem key={p} value={p}>
              {p}
            </MenuItem>
          ))}
        </StyledTextField>
        <StyledTextField
          select
          fullWidth
          label="Series"
          value={selectedSeries}
          onChange={handleChange('series')}
          margin="dense"
          variant="outlined"
        >
          {series.map((s) => (
            <MenuItem key={s} value={s}>
              {s}
            </MenuItem>
          ))}
        </StyledTextField>
      </Root>
    </NoSsr>
  );
};

GrafanaMetricsCompare.propTypes = {
  chartCompare: PropTypes.array.isRequired,
};

const mapDispatchToProps = () => ({});

export default connect(null, mapDispatchToProps)(GrafanaMetricsCompare);
