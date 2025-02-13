import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import GrafanaDateRangePicker from './GrafanaDateRangePicker';
import { ExpansionPanel, ExpansionPanelSummary } from '../../ExpansionPanels';
import GrafanaCustomChart from './GrafanaCustomChart';
import {
  Grid,
  ExpansionPanelDetails,
  Typography,
  Dialog,
  Button,
  DialogActions,
  DialogContent,
  DialogTitle,
  Chip,
  styled,
  useTheme,
} from '@layer5/sistent';

const GrafanaRoot = styled('div')({
  width: '100%',
});

const Column = styled('div')({
  flex: '1',
});

const StyledHeading = styled(Typography)(({ theme }) => ({
  fontSize: theme.typography.pxToRem(15),
}));

const SecondaryHeading = styled(Typography)(() => {
  const theme = useTheme();
  return {
    fontSize: theme.typography.pxToRem(15),
    color: theme.palette.text.secondary,
  };
});

const DateRangePickerContainer = styled('div')({
  display: 'flex',
  justifyContent: 'flex-end',
});

const ChartsHeaderOptions = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  marginBottom: '1rem',
  marginTop: '1rem',
}));

const StyledIcon = styled('img')(({ theme }) => ({
  width: theme.spacing(2.5),
}));

const StyledChip = styled(Chip)({
  width: '100%',
});

const StyledDialogTitle = styled(DialogTitle)({
  '& > *': {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
});

class GrafanaCustomCharts extends Component {
  constructor(props) {
    super(props);

    const newStartDate = new Date();
    newStartDate.setMinutes(newStartDate.getMinutes() - 5);
    const { startDate, from, endDate, to, liveTail, sparkline } = props;
    this.state = {
      startDate: startDate || newStartDate,
      from: from && from !== null ? from : 'now-5m',
      endDate: endDate && endDate !== null ? endDate : new Date(),
      to: to && to !== null ? to : 'now',
      liveTail: liveTail && liveTail !== null ? liveTail : true,
      refresh: '10s',
      sparkline: sparkline && sparkline !== null ? true : false,
      chartDialogOpen: false,
      chartDialogPanelData: {},
      chartDialogPanel: {},
      chartDialogBoard: {},
    };
  }

  updateDateRange = (from, startDate, to, endDate, liveTail, refresh) => {
    this.setState({
      from,
      startDate,
      to,
      endDate,
      liveTail,
      refresh,
    });
  };

  chartDialogClose() {
    const self = this;
    return () => {
      self.setState({ chartDialogOpen: false });
    };
  }

  handleChartDialogOpen = (board, panel, data) => {
    this.setState({
      chartDialogOpen: true,
      chartDialogBoard: board,
      chartDialogPanel: panel,
      chartDialogPanelData: data,
    });
  };

  GrafanaChip(grafanaURL) {
    return (
      <StyledChip
        label={grafanaURL}
        onClick={() => window.open(grafanaURL)}
        icon={<StyledIcon src="/static/img/grafana_icon.svg" />}
        variant="outlined"
      />
    );
  }
  render() {
    const {
      from,
      startDate,
      to,
      endDate,
      liveTail,
      refresh,
      chartDialogOpen,
      chartDialogPanel,
      chartDialogBoard,
      chartDialogPanelData,
      sparkline,
    } = this.state;
    const { boardPanelConfigs, boardPanelData } = this.props;
    const { grafanaURL, grafanaAPIKey, prometheusURL, connectionID } = this.props;
    const { enableGrafanaChip } = this.props;
    // we are now proxying. . .
    // if (grafanaURL && grafanaURL.endsWith('/')){
    //   grafanaURL = grafanaURL.substring(0, grafanaURL.length - 1);
    // }
    return (
      <NoSsr>
        <React.Fragment>
          <GrafanaRoot>
            {!(boardPanelData && boardPanelData !== null) && (
              <ChartsHeaderOptions>
                {enableGrafanaChip && <div>{this.GrafanaChip(grafanaURL)}</div>}
                <DateRangePickerContainer>
                  <GrafanaDateRangePicker
                    from={from}
                    startDate={startDate}
                    to={to}
                    endDate={endDate}
                    liveTail={liveTail}
                    refresh={refresh}
                    updateDateRange={this.updateDateRange}
                  />
                </DateRangePickerContainer>
              </ChartsHeaderOptions>
            )}

            <Dialog
              fullWidth
              maxWidth="md"
              open={chartDialogOpen}
              onClose={this.chartDialogClose()}
              aria-labelledby="max-width-dialog-title"
            >
              <StyledDialogTitle id="max-width-dialog-title">
                <div>{chartDialogPanel.title}</div>
                {!(
                  chartDialogPanelData &&
                  chartDialogPanelData !== null &&
                  Object.keys(chartDialogPanelData).length > 0
                ) ? (
                  <DateRangePickerContainer>
                    <GrafanaDateRangePicker
                      from={from}
                      startDate={startDate}
                      to={to}
                      endDate={endDate}
                      liveTail={liveTail}
                      refresh={refresh}
                      updateDateRange={this.updateDateRange}
                    />
                  </DateRangePickerContainer>
                ) : (
                  <div></div>
                )}
              </StyledDialogTitle>
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
                  connectionID={connectionID}
                  // testUUID={testUUID} // this is just a dialog, we dont want this series too to be persisted
                  panelData={
                    chartDialogPanelData && chartDialogPanelData !== null
                      ? chartDialogPanelData
                      : {}
                  }
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
              <ExpansionPanel key={ind} square defaultExpanded={ind === 0}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                  <Column>
                    <StyledHeading variant="subtitle1" gutterBottom>
                      {config.board && config.board.title
                        ? config.board.title
                        : config.title
                          ? config.title
                          : ''}
                    </StyledHeading>
                  </Column>
                  {config.templateVars && config.templateVars.length > 0 && (
                    <Column>
                      <SecondaryHeading variant="subtitle2">
                        {`Template variables: ${config.templateVars.join(' ')}`}
                      </SecondaryHeading>
                    </Column>
                  )}
                </ExpansionPanelSummary>
                <ExpansionPanelDetails>
                  <Grid container spacing={3}>
                    {config.panels.map(
                      (panel, i) => (
                        // if(panel.type === 'graph'){

                        <Grid key={`grafana-chart-${i}`} item xs={12} lg={sparkline ? 12 : 6}>
                          <GrafanaCustomChart
                            connectionID={connectionID}
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
                            testUUID={config.testUUID ? config.testUUID : ''}
                            panelData={
                              boardPanelData &&
                              boardPanelData !== null &&
                              boardPanelData[ind] &&
                              boardPanelData[ind] !== null
                                ? boardPanelData[ind]
                                : {}
                            }
                          />
                        </Grid>
                      ),
                      // } else return '';
                    )}
                  </Grid>
                </ExpansionPanelDetails>
              </ExpansionPanel>
            ))}
          </GrafanaRoot>
        </React.Fragment>
      </NoSsr>
    );
  }
}

GrafanaCustomCharts.propTypes = {
  // grafanaURL: PropTypes.string.isRequired,
  // grafanaAPIKey: PropTypes.string.isRequired,
  boardPanelConfigs: PropTypes.array.isRequired,
  connectionID: PropTypes.string.isRequired,
  // boardPanelData:
};

export default GrafanaCustomCharts;
