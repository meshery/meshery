import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@sistent/sistent';
import { Grid2, ExpansionPanelDetails, Typography, styled } from '@sistent/sistent';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import LazyLoad from 'react-lazyload';
import GrafanaDateRangePicker from './GrafanaDateRangePicker';
import { ExpansionPanel, ExpansionPanelSummary } from '../../ExpansionPanels';

const Wrapper = styled('div')({
  width: '100%',
});

const Column = styled('div')({
  flexBasis: '33.33%',
});

const DateRangePickerContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  justifyContent: 'flex-end',
  marginRight: theme.spacing(1),
  marginBottom: theme.spacing(2),
}));

const StyledHeading = styled(Typography)(({ theme }) => ({
  fontSize: theme.typography.pxToRem(15),
}));

const SecondaryHeading = styled(Typography)(({ theme }) => ({
  fontSize: theme.typography.pxToRem(15),
  color: theme.palette.text.secondary,
}));

const IframeGridItem = styled(Grid2)(({ theme }) => ({
  minHeight: theme.spacing(55),
  minWidth: theme.spacing(55),
  '& iframe': {
    width: '100%',
    height: '100%',
    border: 'none',
  },
}));

const GrafanaCharts = ({ grafanaURL, boardPanelConfigs }) => {
  const [dateRange, setDateRange] = useState({
    startDate: (() => {
      const d = new Date();
      d.setMinutes(d.getMinutes() - 5);
      return d;
    })(),
    from: 'now-5m',
    endDate: new Date(),
    to: 'now',
    liveTail: true,
    refresh: '10s',
  });

  const { from, startDate, to, endDate, liveTail, refresh } = dateRange;

  let adjustedGrafanaURL = grafanaURL;
  if (adjustedGrafanaURL.endsWith('/')) {
    adjustedGrafanaURL = adjustedGrafanaURL.slice(0, -1);
  }

  const updateDateRange = (from, startDate, to, endDate, liveTail, refresh) => {
    setDateRange({ from, startDate, to, endDate, liveTail, refresh });
  };

  return (
    <NoSsr>
      <Wrapper>
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
        {boardPanelConfigs.map((config, ind) => (
          <ExpansionPanel key={ind} square defaultExpanded={false}>
            <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
              <Column>
                <StyledHeading variant="subtitle1" gutterBottom>
                  {config.board.title}
                </StyledHeading>
              </Column>
              <Column>
                <SecondaryHeading variant="subtitle2">
                  {config.templateVars && config.templateVars.length > 0
                    ? `Template variables: ${config.templateVars.join(' ')}`
                    : ''}
                </SecondaryHeading>
              </Column>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails>
              <Grid2 container spacing={5} size="grow">
                {config.panels.map((panel, ind) => (
                  <IframeGridItem key={ind} size={{ xs: 12, sm: 6, md: 4 }}>
                    <LazyLoad once>
                      <iframe
                        key={`url_-_-${ind}`}
                        src={`${adjustedGrafanaURL}/d-solo/${config.board.uid}/${config.board.slug}?theme=light&orgId=${config.board.org_id}&panelId=${panel.id}&refresh=${refresh}&from=${from}&to=${to}&${config.templateVars
                          .map((tv) => `var-${tv}`)
                          .join('&')}`}
                      />
                    </LazyLoad>
                  </IframeGridItem>
                ))}
              </Grid2>
            </ExpansionPanelDetails>
          </ExpansionPanel>
        ))}
      </Wrapper>
    </NoSsr>
  );
};

GrafanaCharts.propTypes = {
  grafanaURL: PropTypes.string.isRequired,
  boardPanelConfigs: PropTypes.array.isRequired,
};

export default GrafanaCharts;
