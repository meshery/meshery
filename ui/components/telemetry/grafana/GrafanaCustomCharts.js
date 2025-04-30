import React, { useState, useCallback } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import GrafanaDateRangePicker from './GrafanaDateRangePicker';
import { StyledAccordion, StyledAccordionSummary } from '../../StyledAccordion';
import GrafanaCustomChart from './GrafanaCustomChart';
import {
  Grid,
  AccordionDetails,
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

const GrafanaCustomCharts = (props) => {
  const {
    boardPanelConfigs,
    boardPanelData,
    grafanaURL,
    grafanaAPIKey,
    prometheusURL,
    connectionID,
    enableGrafanaChip,
    startDate: propStartDate,
    from: propFrom,
    endDate: propEndDate,
    to: propTo,
    liveTail: propLiveTail,
    sparkline: propSparkline,
  } = props;

  const newStartDate = new Date();
  newStartDate.setMinutes(newStartDate.getMinutes() - 5);

  const [startDate, setStartDate] = useState(propStartDate || newStartDate);
  const [from, setFrom] = useState(propFrom && propFrom !== null ? propFrom : 'now-5m');
  const [endDate, setEndDate] = useState(
    propEndDate && propEndDate !== null ? propEndDate : new Date(),
  );
  const [to, setTo] = useState(propTo && propTo !== null ? propTo : 'now');
  const [liveTail, setLiveTail] = useState(
    propLiveTail && propLiveTail !== null ? propLiveTail : true,
  );
  const [refresh, setRefresh] = useState('10s');
  const [sparkline] = useState(propSparkline && propSparkline !== null ? true : false);
  const [chartDialogOpen, setChartDialogOpen] = useState(false);
  const [chartDialogPanelData, setChartDialogPanelData] = useState({});
  const [chartDialogPanel, setChartDialogPanel] = useState({});
  const [chartDialogBoard, setChartDialogBoard] = useState({});

  const updateDateRange = useCallback((from, startDate, to, endDate, liveTail, refresh) => {
    setFrom(from);
    setStartDate(startDate);
    setTo(to);
    setEndDate(endDate);
    setLiveTail(liveTail);
    setRefresh(refresh);
  }, []);

  const chartDialogClose = useCallback(() => {
    return () => {
      setChartDialogOpen(false);
    };
  }, []);

  const handleChartDialogOpen = useCallback((board, panel, data) => {
    setChartDialogOpen(true);
    setChartDialogBoard(board);
    setChartDialogPanel(panel);
    setChartDialogPanelData(data);
  }, []);

  const GrafanaChip = useCallback((grafanaURL) => {
    return (
      <StyledChip
        label={grafanaURL}
        onClick={() => window.open(grafanaURL)}
        icon={<StyledIcon src="/static/img/grafana_icon.svg" />}
        variant="outlined"
      />
    );
  }, []);

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
              {enableGrafanaChip && <div>{GrafanaChip(grafanaURL)}</div>}
              <DateRangePickerContainer>
                <GrafanaDateRangePicker
                  from={from}
                  startDate={startDate}
                  to={to}
                  endDate={endDate}
                  liveTail={liveTail}
                  refresh={refresh}
                  updateDateRange={updateDateRange}
                />
              </DateRangePickerContainer>
            </ChartsHeaderOptions>
          )}

          <Dialog
            fullWidth
            maxWidth="md"
            open={chartDialogOpen}
            onClose={chartDialogClose()}
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
                    updateDateRange={updateDateRange}
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
                handleChartDialogOpen={handleChartDialogOpen}
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
                updateDateRange={updateDateRange}
                inDialog
                connectionID={connectionID}
                // testUUID={testUUID} // this is just a dialog, we dont want this series too to be persisted
                panelData={
                  chartDialogPanelData && chartDialogPanelData !== null ? chartDialogPanelData : {}
                }
              />
            </DialogContent>
            <DialogActions>
              <Button onClick={chartDialogClose()} color="primary">
                Close
              </Button>
            </DialogActions>
          </Dialog>

          {boardPanelConfigs.map((config, ind) => (
            // <ExpansionPanel defaultExpanded={ind === 0?true:false}>
            <StyledAccordion key={ind} square defaultExpanded={ind === 0}>
              <StyledAccordionSummary expandIcon={<ExpandMoreIcon />}>
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
              </StyledAccordionSummary>
              <AccordionDetails>
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
                          handleChartDialogOpen={handleChartDialogOpen}
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
                          updateDateRange={updateDateRange}
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
              </AccordionDetails>
            </StyledAccordion>
          ))}
        </GrafanaRoot>
      </React.Fragment>
    </NoSsr>
  );
};

GrafanaCustomCharts.propTypes = {
  grafanaURL: PropTypes.string.isRequired,
  grafanaAPIKey: PropTypes.string.isRequired,
  boardPanelConfigs: PropTypes.array.isRequired,
  connectionID: PropTypes.string.isRequired,
};

export default GrafanaCustomCharts;
