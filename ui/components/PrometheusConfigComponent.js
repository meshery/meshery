import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, TextField, Grid, Button } from '@material-ui/core';

const promStyles = theme => ({
  root: {
    padding: theme.spacing(5),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    //   marginLeft: theme.spacing(1),
  }
});

class PrometheusConfigComponent extends Component {
    render = () => {
      const { classes, prometheusURL, urlError, handleChange, handlePrometheusConfigure } = this.props;
      return (
        <NoSsr>
          <React.Fragment>
            <div className={classes.root}>
              <Grid container spacing={1}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    required
                    id="prometheusURL"
                    name="prometheusURL"
                    label="Prometheus Base URL"
                    type="url"
                    autoFocus
                    fullWidth
                    value={prometheusURL}
                    error={urlError}
                    margin="normal"
                    variant="outlined"
                    onChange={handleChange('prometheusURL')}
                  />
                </Grid>
              </Grid>
              <div className={classes.buttons}>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  size="large"
                  onClick={handlePrometheusConfigure}
                  className={classes.button}
                >
                Submit
                </Button>
              </div>
            </div>
          </React.Fragment>
        </NoSsr>
      );
    }
}

PrometheusConfigComponent.propTypes = {
  classes: PropTypes.object.isRequired,
  prometheusURL: PropTypes.string.isRequired,
  handleChange: PropTypes.func.isRequired, 
  handlePrometheusConfigure: PropTypes.func.isRequired,
};

export default withStyles(promStyles)(PrometheusConfigComponent);