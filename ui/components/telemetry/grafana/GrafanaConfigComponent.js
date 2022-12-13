import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {
  NoSsr, TextField, Grid, Button,
} from '@material-ui/core';
import ReactSelectWrapper from '../../ReactSelectWrapper'

const grafanaStyles = (theme) => ({
  wrapper : { padding : theme.spacing(5), },
  buttons : { display : 'flex',
    justifyContent : 'flex-end', },
  inputContainer : { marginTop : theme.spacing(2),
    marginBottom : theme.spacing(1) },
  button : { marginTop : theme.spacing(3),
    //   marginLeft: theme.spacing(1),
  },
});

class GrafanaConfigComponent extends Component {
    render = () => {
      const {
        classes, grafanaURL, grafanaAPIKey, urlError, handleChange, handleGrafanaConfigure, options, handleChangeApiKey
      } = this.props;
      return (
        <NoSsr>
          <React.Fragment>
            <div className={classes.wrapper}>
              <Grid container spacing={1}>
                <Grid item xs={12} md={6}>
                  <div className={classes.inputContainer}>
                    <ReactSelectWrapper
                      onChange={(select) => handleChange('grafanaURL')(select ? select.value : '')}
                      options={options}
                      value={grafanaURL}
                      label="Grafana Base URL"
                      error={urlError}
                      placeholder="Address of Grafana Server"
                      noOptionsMessage="No Grafana servers discovered"
                    />
                  </div>
                </Grid>
                <Grid item xs={12} md={6}>
                  <TextField
                    id="grafanaAPIKey"
                    name="grafanaAPIKey"
                    label="API Key"
                    fullWidth
                    value={grafanaAPIKey}
                    margin="normal"
                    variant="outlined"
                    onKeyDown={(e) => {
                      if (e.keyCode == 13) {
                        handleGrafanaConfigure();
                      }
                    }}
                    onChange={handleChangeApiKey}
                  />
                </Grid>
              </Grid>
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
            </div>
          </React.Fragment>
        </NoSsr>
      );
    }
}

GrafanaConfigComponent.propTypes = {
  classes : PropTypes.object.isRequired,
  grafanaURL : PropTypes.object.isRequired,
  grafanaAPIKey : PropTypes.string.isRequired,
  handleChange : PropTypes.func.isRequired,
  handleGrafanaConfigure : PropTypes.func.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaConfigComponent);
