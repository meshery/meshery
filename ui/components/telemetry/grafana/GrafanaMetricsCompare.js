import { Component } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { MenuItem, TextField, Box, styled } from '@layer5/sistent';
import { connect } from 'react-redux';
import { UsesSistent } from '@/components/SistentWrapper';

const Root = styled(Box)(() => ({
  width: '100%',
}));

const StyledTextField = styled(TextField)(({ theme }) => ({
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(1),
}));

class GrafanaMetricsCompare extends Component {
  constructor(props) {
    super(props);
    const { chartCompare } = props;
    const panels = this.computePanels(chartCompare);
    this.state = {
      chartCompare,
      panels,
      panel: '',
      selectedSeries: '',
      series: [],
    };
  }

  computePanels(chartCompare) {
    const panels = {};
    if (chartCompare && chartCompare !== null && chartCompare.length > 0) {
      chartCompare.forEach((cc) => {
        if (
          cc.boardConfig &&
          cc.boardConfig !== null &&
          cc.boardConfig.panels &&
          cc.boardConfig.panels !== null
        ) {
          cc.boardConfig.panels.forEach((panel) => {
            // if(panels.indexOf(panel.title) === -1){
            //   panels.push(panel);
            // }
            panels[panel.title] = panel;
          });
        }
      });
    }
    return panels;
  }

  componentDidMount() {
    const { panels } = this.state;
    const panel = Object.keys(panels).length > 0 ? Object.keys(panels)[0] : '';
    let series = [];
    if (panels[panel] && panels[panel].targets) {
      series = panels[panel].targets.map((target) => target.expr);
    }
    this.setState({ panel, series, selectedSeries: series.length > 0 ? series[0] : '' });
  }

  handleChange = (name) => (event) => {
    const { panels } = this.state;
    if (name === 'panel') {
      let series = [];
      const panel = event.target.value;
      if (panels[panel] && panels[panel].targets) {
        series = panels[panel].targets.map((target) => target.expr);
      }
      this.setState({ panel, series, selectedSeries: series.length > 0 ? series[0] : '' });
    } else if (name === 'series') {
      this.setState({ selectedSeries: event.target.value });
    }
  };

  render() {
    const { panels, panel, selectedSeries, series } = this.state;

    return (
      <UsesSistent>
        <NoSsr>
          <Root>
            <StyledTextField
              select
              fullWidth
              label="Panel"
              value={panel}
              onChange={this.handleChange('panel')}
              margin="dense"
              variant="outlined"
            >
              {panels &&
                Object.keys(panels).map((p) => (
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
              onChange={this.handleChange('series')}
              margin="dense"
              variant="outlined"
            >
              {series &&
                series.map((s) => (
                  <MenuItem key={s} value={s}>
                    {s}
                  </MenuItem>
                ))}
            </StyledTextField>
          </Root>
        </NoSsr>
      </UsesSistent>
    );
  }
}

GrafanaMetricsCompare.propTypes = {
  classes: PropTypes.object.isRequired,
  chartCompare: PropTypes.array.isRequired,
};

const mapDispatchToProps = () => ({});

export default connect(null, mapDispatchToProps)(GrafanaMetricsCompare);
