import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Grid, Button, } from '@material-ui/core';
import ReactSelectWrapper from "../../ReactSelectWrapper";

const promStyles = (theme) => ({ promRoot : { padding : theme.spacing(5), },
  buttons : { display : 'flex',
    justifyContent : 'flex-end', },
  button : { marginTop : theme.spacing(3),
    //   marginLeft: theme.spacing(1),
  }, });

class PrometheusConfigComponent extends Component {
    render = () => {
      const {
        classes, prometheusURL, urlError, handleChange, handlePrometheusConfigure, options = []
      } = this.props;
      return (
        <NoSsr>
          <React.Fragment>
            <div className={classes.promRoot}>
              <Grid container spacing={1}>
                <Grid item xs={12}>
                  <ReactSelectWrapper
                    onChange={(select) => handleChange('prometheusURL')(select ? select.value : '')}
                    options={options}
                    value={prometheusURL}
                    label="Prometheus Base URL"
                    placeholder="Address of Prometheus Server"
                    noOptionsMessage="No Prometheus servers discovered"
                    error={urlError}
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
  classes : PropTypes.object.isRequired,
  prometheusURL : PropTypes.object.isRequired,
  handleChange : PropTypes.func.isRequired,
  handlePrometheusConfigure : PropTypes.func.isRequired,
  options : PropTypes.array.isRequired
};

export default withStyles(promStyles)(PrometheusConfigComponent);
