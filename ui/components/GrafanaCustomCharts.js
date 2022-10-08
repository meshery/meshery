import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, Grid, ExpansionPanelDetails, Typography, Dialog, Button, DialogActions, DialogContent, DialogTitle, Chip
} from '@material-ui/core';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import GrafanaDateRangePicker from './GrafanaDateRangePicker';
import { ExpansionPanel, ExpansionPanelSummary } from './ExpansionPanels';
import GrafanaCustomChart from './GrafanaCustomChart';

const grafanaStyles = (theme) => ({
  grafanaRoot : { width : '100%', },
  column : { flex : '1', },
  heading : { fontSize : theme.typography.pxToRem(15), },
  secondaryHeading : {
    fontSize : theme.typography.pxToRem(15),
    color : theme.palette.text.secondary,
  },
  dateRangePicker : {
    display : 'flex',
    justifyContent : 'flex-end',
  },
  chartsHeaderOptions : {
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    marginBottom : "1rem",
    marginTop : "1rem",
  },
  icon : { width : theme.spacing(2.5), },
  dialogTitle : {
    '&>*' : {
      display : 'flex',
      alignItems : 'center',
      justifyContent : 'space-between',
    }
  },
});


class GrafanaCustomCharts extends Component {
  constructor(props) {
    super(props);

    const newStartDate = new Date();
    newStartDate.setMinutes(newStartDate.getMinutes() - 5);
    const {
      startDate, from, endDate, to, liveTail, sparkline
    } = props;
    this.state = {
      startDate : startDate || newStartDate,
      from : from && from !== null ? from
        : 'now-5m',
      endDate : endDate && endDate !== null ? endDate
        : new Date(),
      to : to && to !== null ? to
        : 'now',
      liveTail : liveTail && liveTail !== null ? liveTail
        : true,
      refresh : '10s',
      sparkline : sparkline && sparkline !== null ? true
        : false,
      chartDialogOpen : false,
      chartDialogPanelData : {},
      chartDialogPanel : {},
      chartDialogBoard : {},
    };
  }

  updateDateRange = (from, startDate, to, endDate, liveTail, refresh) => {
    this.setState({
      from, startDate, to, endDate, liveTail, refresh,
    });
  }

  chartDialogClose() {
    const self = this;
    return () => {
      self.setState({ chartDialogOpen : false });
    };
  }

  handleChartDialogOpen = (board, panel, data) => {
    this.setState({
      chartDialogOpen : true,
      chartDialogBoard : board,
      chartDialogPanel : panel,
      chartDialogPanelData : data,
    });
  }

  GrafanaChip(grafanaURL) {
    const { classes } = this.props;
    return (
      <Chip
        label={grafanaURL}
        onClick={() => window.open(grafanaURL)}
        icon={<img src="/static/img/grafana_icon.svg" className={classes.icon} />}
        className={classes.chip}
        variant="outlined"
      />
    )
  }
  render() {
    const {
      from, startDate, to, endDate, liveTail, refresh, chartDialogOpen, chartDialogPanel, chartDialogBoard,
      chartDialogPanelData, sparkline
    } = this.state;
    const { classes, boardPanelConfigs, boardPanelData } = this.props;
    const { grafanaURL, grafanaAPIKey, prometheusURL } = this.props;
    const { enableGrafanaChip } = this.props;
    // we are now proxying. . .
    // if (grafanaURL && grafanaURL.endsWith('/')){
    //   grafanaURL = grafanaURL.substring(0, grafanaURL.length - 1);
    // }
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.grafanaRoot}>
            {!(boardPanelData && boardPanelData !== null)
              && (
                <div className={classes.chartsHeaderOptions}>
                  {
                    enableGrafanaChip && (
                      <div>
                        {this.GrafanaChip(grafanaURL)}
                      </div>)
                  }
                  <div className={classes.dateRangePicker}>
                    <GrafanaDateRangePicker
                      from={from}
                      startDate={startDate}
                      to={to}
                      endDate={endDate}
                      liveTail={liveTail}
                      refresh={refresh}
                      updateDateRange={this.updateDateRange}
                    />
                  </div>
                </div>
              )}
            <Dialog
              fullWidth
              maxWidth="md"
              open={chartDialogOpen}
              onClose={this.chartDialogClose()}
              aria-labelledby="max-width-dialog-title"
            >
              <DialogTitle classes={{ root : classes.dialogTitle }} id="max-width-dialog-title">
                <div>
                  {chartDialogPanel.title}
                </div>
                {!(chartDialogPanelData && chartDialogPanelData !== null && Object.keys(chartDialogPanelData).length > 0)
                  ? (
                    <div className={classes.dateRangePicker}>
                      <GrafanaDateRangePicker
                        from={from}
                        startDate={startDate}
                        to={to}
                        endDate={endDate}
                        liveTail={liveTail}
                        refresh={refresh}
                        updateDateRange={this.updateDateRange}
                      />
                    </div>
                  )
                  : (<div></div>)}
              </DialogTitle>
              <DialogContent>
                <GrafanaCustomChart
                  board={chartDialogBoard}
                  panel={chartDialogPanel}
                  handleChartDialogOpen={this.handleChartDialogOpen}
                  grafanaURL={grafanaURL}
                  grafanaAPIKey={grafanaAPIKey}
                  prometheusURL={prometheusURL}
                  from={from}
                  startDate={startDate}
                  to={to}
                  endDate={endDate}
                  liveTail={liveTail}
                  refresh={refresh}
                  templateVars={chartDialogBoard.templateVars}
                  updateDateRange={this.updateDateRange}
                  inDialog
                  // testUUID={testUUID} // this is just a dialog, we dont want this series too to be persisted
                  panelData={chartDialogPanelData && chartDialogPanelData !== null
                    ? chartDialogPanelData
                    : {}}
                />
              </DialogContent>
              <DialogActions>
                <Button onClick={this.chartDialogClose()} color="primary">
                  Close
                </Button>
              </DialogActions>
            </Dialog>


            {boardPanelConfigs.map((config, ind) => (
              // <ExpansionPanel defaultExpanded={ind === 0?true:false}>
              <ExpansionPanel square defaultExpanded={ind === 0}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                  <div className={classes.column}>
                    <Typography variant="subtitle1" gutterBottom>{config.board && config.board.title
                      ? config.board.title
                      : (config.title
                        ? config.title
                        : '')}</Typography>
                  </div>
                  {(config.templateVars && config.templateVars.length > 0) && (
                    <div className={classes.column}>
                      <Typography variant="subtitle2">
                        {`Template variables: ${config.templateVars.join(' ')}`}
                      </Typography>
                    </div>)}
                </ExpansionPanelSummary>
                <ExpansionPanelDetails>
                  <Grid container spacing={3}>
                    {config.panels.map((panel, i) =>
                    // if(panel.type === 'graph'){
                      (
                        <Grid key={`grafana-chart-${i}`} item xs={12} lg={sparkline
                          ? 12
                          : 6}>
                          <GrafanaCustomChart
                            board={config}
                            sparkline={sparkline}
                            panel={panel}
                            handleChartDialogOpen={this.handleChartDialogOpen}
                            grafanaURL={grafanaURL}
                            grafanaAPIKey={grafanaAPIKey}
                            prometheusURL={prometheusURL}
                            from={from}
                            startDate={startDate}
                            to={to}
                            endDate={endDate}
                            liveTail={liveTail}
                            refresh={refresh}
                            templateVars={config.templateVars}
                            updateDateRange={this.updateDateRange}
                            inDialog={false}
                            testUUID={config.testUUID
                              ? config.testUUID
                              : ''}
                            panelData={boardPanelData && boardPanelData !== null && boardPanelData[ind] && boardPanelData[ind] !== null
                              ? boardPanelData[ind]
                              : {}}
                          />

                        </Grid>
                      ),
                      // } else return '';
                    )}
                  </Grid>
                </ExpansionPanelDetails>
              </ExpansionPanel>
            ))}
          </div>
        </React.Fragment>
      </NoSsr>
    );
  }
}

GrafanaCustomCharts.propTypes = {
  classes : PropTypes.object.isRequired,
  // grafanaURL: PropTypes.string.isRequired,
  // grafanaAPIKey: PropTypes.string.isRequired,
  boardPanelConfigs : PropTypes.array.isRequired,
  // boardPanelData:
};

export default withStyles(grafanaStyles)(GrafanaCustomCharts);
