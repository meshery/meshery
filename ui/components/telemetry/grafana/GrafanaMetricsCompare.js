import { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { MenuItem, NoSsr, TextField } from '@material-ui/core';
import { connect } from 'react-redux';

const grafanaStyles = () => ({ root : { width : '100%', }, });

class GrafanaMetricsCompare extends Component {
  constructor(props) {
    super(props);
    const { chartCompare } = props;
    const panels = this.computePanels(chartCompare);
    this.state = {
      chartCompare,
      panels,
      panel : '',
      selectedSeries : '',
      series : [],
    };
  }

  computePanels(chartCompare) {
    const panels = {};
    if (chartCompare && chartCompare !== null && chartCompare.length > 0) {
      chartCompare.forEach((cc) => {
        if (cc.boardConfig && cc.boardConfig !== null && cc.boardConfig.panels && cc.boardConfig.panels !== null) {
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
    this.setState({ panel,
      series,
      selectedSeries : series.length > 0
        ? series[0]
        : '', });
  }

  handleChange(name) {
    const { panels } = this.state;
    const self = this;
    return (event) => {
      if (name === 'panel') {
        let series = [];
        const panel = event.target.value;
        if (panels[panel] && panels[panel].targets) {
          series = panels[panel].targets.map((target) => target.expr);
        }
        self.setState({ panel,
          series,
          selectedSeries : series.length > 0
            ? series[0]
            : '', });
      } else if (name === 'series') {
        self.setState({ selectedSeries : event.target.value });
      }
    };
  }

  render() {
    const {
      panels, panel, selectedSeries, series,
    } = this.state;

    return (
      <NoSsr>
        <TextField
          select
          id="panel"
          name="panel"
          label="Panel"
          fullWidth
          value={panel}
          margin="normal"
          variant="outlined"
          onChange={this.handleChange('panel')}
        >
          {panels && Object.keys(panels).map((p) => <MenuItem key={p} value={p}>{p}</MenuItem>)}
        </TextField>
        <TextField
          select
          id="series"
          name="series"
          label="Series"
          fullWidth
          value={selectedSeries}
          margin="normal"
          variant="outlined"
          onChange={this.handleChange('series')}
        >
          {series && series.map((s) => <MenuItem key={s} value={s}>{s}</MenuItem>)}
        </TextField>
      </NoSsr>
    );
  }
}

GrafanaMetricsCompare.propTypes = { classes : PropTypes.object.isRequired,
  chartCompare : PropTypes.array.isRequired, };

const mapDispatchToProps = () => ({});

export default withStyles(grafanaStyles)(connect(
  null,
  mapDispatchToProps,
)(GrafanaMetricsCompare));
