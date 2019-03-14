import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Grid, ExpansionPanel, ExpansionPanelSummary, ExpansionPanelDetails, Typography } from '@material-ui/core';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

const grafanaStyles = theme => ({
    root: {
      width: '100%',
    },
    column: {
      flexBasis: '33.33%',
    },
    heading: {
      fontSize: theme.typography.pxToRem(15),
    },
    secondaryHeading: {
      fontSize: theme.typography.pxToRem(15),
      color: theme.palette.text.secondary,
    },
  });

class GrafanaCharts extends Component {
    
    render() {
        const { classes, grafanaURL, boardPanelConfigs } = this.props;
        return (
              <NoSsr>
              <React.Fragment>
              <div className={classes.root}>
                {boardPanelConfigs.map((config, ind) => (
                  <ExpansionPanel defaultExpanded={ind === 0?true:false}>
                    <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                      <div className={classes.column}>
                        <Typography className={classes.heading}>{config.board.title}</Typography>
                      </div>
                      <div className={classes.column}>
                        <Typography className={classes.secondaryHeading}>{config.templateVars.join(', ')}</Typography>
                      </div>
                    </ExpansionPanelSummary>
                    <ExpansionPanelDetails>
                        <Grid container spacing={5}>
                          {config.panels.map(panel => (
                            <Grid item xs={12} sm={4}>
                              <iframe 
                                key={'url_-_-'+ind} 
                                src={`${grafanaURL}/d-solo/${config.board.uid}/${config.board.slug}?theme=light&orgId=${config.board.org_id}&panelId=${panel.id}&refresh=10s&${config.templateVars.map(tv => `var-${tv}`).join('&')}`} 
                                width='450' 
                                height='250' 
                                frameBorder='0'>
                              </iframe>
                            </Grid>
                          ))}
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

GrafanaCharts.propTypes = {
  classes: PropTypes.object.isRequired,
  grafanaURL: PropTypes.string.isRequired,
  boardPanelConfigs: PropTypes.array.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaCharts);