import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Grid } from '@material-ui/core';

const grafanaStyles = theme => ({
    root: {
      padding: theme.spacing(10),
    },
  });

class GrafanaCharts extends Component {
    
    render() {
        const { classes, grafanaURL, boardPanelConfigs } = this.props;
        const urls = [];
        boardPanelConfigs.forEach(config => {
          const tmpVars = config.templateVars.map(tv => `var-${tv}`).join('&');
          config.panels.forEach(panel => urls.push(`${grafanaURL}/d-solo/${config.board.uid}/${config.board.slug}?theme=light&orgId=${config.board.org_id}&panelId=${panel.id}&refresh=10s&${tmpVars}`));
        })
        return (
              <NoSsr>
              <React.Fragment>
              <Grid container spacing={5}>
                {urls.map((ul, ind) => (
                <Grid item xs={12} sm={4}>
                  <iframe key={'url_-_-'+ind} src={ul} width='450' height='250' frameBorder='0'></iframe>
                </Grid>
                ))}
              </Grid>
              </React.Fragment>
              </NoSsr>
            );
        }
}

GrafanaCharts.propTypes = {
  classes: PropTypes.object.isRequired,
  grafanaURL: PropTypes.string.isRequired,
  boardPanelConfigs: PropTypes.array.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaCharts);