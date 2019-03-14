import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Chip } from '@material-ui/core';

const grafanaStyles = theme => ({
    root: {
      padding: theme.spacing(10),
    },
    panelChips: {
      display: 'flex',
      flexWrap: 'wrap',
    },
    panelChip: {
        margin: theme.spacing(0.25),
    }
  });

class GrafanaDisplaySelection extends Component {
    
    render() {
        const { classes, boardPanelConfigs } = this.props;
        return (
              <NoSsr>
              <React.Fragment>
                <table>
                  <tr>
                    <th>Board</th>
                    <th>Panels</th>
                    <th>Template Variables</th>
                  </tr>
                  {boardPanelConfigs.map(config => (
                    <tr>
                      <td>{config.board.title}</td>
                      <td>
                        <div className={classes.panelChips}>
                          {config.panels.map(panel => (
                            <Chip key={panel.id} label={panel.title} className={classes.panelChip} />
                          ))}
                        </div>
                      </td>
                      <td>{config.templateVars.join(',')}</td>
                    </tr>
                  ))}
                </table>
                
              </React.Fragment>
              </NoSsr>
            );
        }
}

GrafanaDisplaySelection.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaDisplaySelection);