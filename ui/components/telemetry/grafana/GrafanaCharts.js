import React, { Component } from 'react';
import { styled } from "@mui/material/styles";
import PropTypes from 'prop-types';
import {
  NoSsr, Grid, ExpansionPanelDetails, Typography,
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import LazyLoad from 'react-lazyload';
import GrafanaDateRangePicker from './GrafanaDateRangePicker';
import { ExpansionPanel, ExpansionPanelSummary } from '../../ExpansionPanels';

const DivWrapper = styled("div")(() => ({
  width : "100%",
}));

const DivColumn = styled("div")(() => ({
  flexBasis : "33.33%",
}));

const DivDateRangePicker = styled("div")(({ theme }) => ({
  display : "flex",
  justifyContent : "flex-end",
  marginRight : theme.spacing(1),
  marginBottom : theme.spacing(2),
}));

const GridIframe = styled(Grid)(({ theme }) => ({
  minHeight : theme.spacing(55),
  minWidth : theme.spacing(55),
}));

class GrafanaCharts extends Component {
  constructor(props) {
    super(props);

    const startDate = new Date();
    startDate.setMinutes(startDate.getMinutes() - 5);
    this.state = {
      startDate,
      from : "now-5m",
      endDate : new Date(),
      to : "now",
      liveTail : true,
      refresh : "10s",
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

  render() {
    const { from, startDate, to, endDate, liveTail, refresh } = this.state;
    const { boardPanelConfigs } = this.props;
    let { grafanaURL } = this.props;
    if (grafanaURL.endsWith("/")) {
      grafanaURL = grafanaURL.substring(0, grafanaURL.length - 1);
    }
    return (
      <NoSsr>
        <React.Fragment>
          <DivWrapper>
            <DivDateRangePicker>
              <GrafanaDateRangePicker
                from={from}
                startDate={startDate}
                to={to}
                endDate={endDate}
                liveTail={liveTail}
                refresh={refresh}
                updateDateRange={this.updateDateRange}
              />
            </DivDateRangePicker>
            {boardPanelConfigs.map((config, ind) => (
              // <ExpansionPanel defaultExpanded={ind === 0?true:false}>
              <ExpansionPanel key={ind} square defaultExpanded={false}>
                <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                  <DivColumn>
                    <Typography variant="subtitle1" gutterBottom>
                      {config.board.title}
                    </Typography>
                  </DivColumn>
                  <DivColumn>
                    <Typography variant="subtitle2">
                      {config.templateVars && config.templateVars.length > 0
                        ? `Template variables: ${config.templateVars.join(" ")}`
                        : ""}
                    </Typography>
                  </DivColumn>
                </ExpansionPanelSummary>
                <ExpansionPanelDetails>
                  <Grid container spacing={5}>
                    {config.panels.map((panel, ind) => (
                      <GridIframe key={ind} item xs={12} sm={6}>
                        <LazyLoad once>
                          <iframe
                            key={`url_-_-${ind}`}
                            src={`${grafanaURL}/d-solo/${config.board.uid}/${
                              config.board.slug
                            }?theme=light&orgId=${
                              config.board.org_id
                            }&panelId=${
                              panel.id
                            }&refresh=${refresh}&from=${from}&to=${to}&${config.templateVars
                              .map((tv) => `var-${tv}`)
                              .join("&")}`}
                            // width='450'
                            width="100%"
                            // height='250'
                            height="100%"
                            frameBorder="0"
                          />
                        </LazyLoad>
                      </GridIframe>
                    ))}
                  </Grid>
                </ExpansionPanelDetails>
              </ExpansionPanel>
            ))}
          </DivWrapper>
        </React.Fragment>
      </NoSsr>
    );
  }
}

GrafanaCharts.propTypes = {
  classes : PropTypes.object.isRequired,
  grafanaURL : PropTypes.string.isRequired,
  boardPanelConfigs : PropTypes.array.isRequired,
};

export default GrafanaCharts;
