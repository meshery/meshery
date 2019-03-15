import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, TextField, Grid, Button } from '@material-ui/core';

const grafanaStyles = theme => ({
    root: {
      padding: theme.spacing(10),
    },
    buttons: {
      display: 'flex',
    //   justifyContent: 'flex-end',
    },
    button: {
      marginTop: theme.spacing(3),
    //   marginLeft: theme.spacing(1),
    }
  });

class GrafanaConfigComponent extends Component {
    render = () => {
        const { classes, grafanaURL, grafanaAPIKey, urlError, handleChange, handleGrafanaConfigure } = this.props;
        return (
          <NoSsr>
        <React.Fragment>
            <div className={classes.root}>
            <Grid container spacing={5}>
                <Grid item xs={12} sm={4}>
                <TextField
                    required
                    id="grafanaURL"
                    name="grafanaURL"
                    label="Grafana Base URL"
                    type="url"
                    autoFocus
                    fullWidth
                    value={grafanaURL}
                    error={urlError}
                    margin="normal"
                    variant="outlined"
                    onChange={handleChange('grafanaURL')}
                />
                </Grid>
                <Grid item xs={12} sm={4}>
                <TextField
                    id="grafanaAPIKey"
                    name="grafanaAPIKey"
                    label="API Key"
                    fullWidth
                    value={grafanaAPIKey}
                    margin="normal"
                    variant="outlined"
                    onChange={handleChange('grafanaAPIKey')}
                />
                </Grid>
                <Grid item xs={12} sm={4}>
                    <div className={classes.buttons}>
                    <Button
                        type="submit"
                        variant="contained"
                        color="primary"
                        size="large"
                        onClick={handleGrafanaConfigure}
                        className={classes.button}
                    >
                    Submit
                    </Button>
                    </div>
                </Grid>
            </Grid>
            </div>
        </React.Fragment>
        </NoSsr>
        );
    }
}

GrafanaConfigComponent.propTypes = {
  classes: PropTypes.object.isRequired,
  grafanaURL: PropTypes.string.isRequired,
  grafanaAPIKey: PropTypes.string.isRequired, 
  handleChange: PropTypes.func.isRequired, 
  handleGrafanaConfigure: PropTypes.func.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaConfigComponent);