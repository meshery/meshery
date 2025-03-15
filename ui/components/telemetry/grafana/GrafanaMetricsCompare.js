import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { MenuItem, TextField, Box, styled } from '@layer5/sistent';
import { connect } from 'react-redux';

const Root = styled(Box)(() => ({
  width: '100%',
}));

const StyledTextField = styled(TextField)(({ theme }) => ({
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(1),
}));

const GrafanaMetricsCompare = ({ chartCompare }) => {
  const [panels, setPanels] = useState({});
  const [panel, setPanel] = useState('');
  const [selectedSeries, setSelectedSeries] = useState('');
  const [series, setSeries] = useState([]);

  useEffect(() => {
    const computePanels = (chartCompare) => {
      const panels = {};
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

    const initialPanel =
      Object.keys(computedPanels).length > 0 ? Object.keys(computedPanels)[0] : '';
    let initialSeries = [];
    if (computedPanels[initialPanel] && computedPanels[initialPanel].targets) {
      initialSeries = computedPanels[initialPanel].targets.map((target) => target.expr);
    }
    setPanel(initialPanel);
    setSeries(initialSeries);
    setSelectedSeries(initialSeries.length > 0 ? initialSeries[0] : '');
  }, [chartCompare]);

  const handleChange = (name) => (event) => {
    if (name === 'panel') {
      const selectedPanel = event.target.value;
      let newSeries = [];
      if (panels[selectedPanel] && panels[selectedPanel].targets) {
        newSeries = panels[selectedPanel].targets.map((target) => target.expr);
      }
      setPanel(selectedPanel);
      setSeries(newSeries);
      setSelectedSeries(newSeries.length > 0 ? newSeries[0] : '');
    } else if (name === 'series') {
      setSelectedSeries(event.target.value);
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
