import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, TextField, Grid, Button } from '@material-ui/core';
import ReactSelectWrapper from '../../ReactSelectWrapper';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { CONNECTION_KINDS, CONNECTION_STATES } from '@/utils/Enum';
import dataFetch from 'lib/data-fetch';

const grafanaStyles = (theme) => ({
  wrapper: {
    padding: theme.spacing(5),
    backgroundColor: theme.palette.secondary.elevatedComponents,
    borderBottomLeftRadius: theme.spacing(1),
    borderBottomRightRadius: theme.spacing(1),
    marginTop: theme.spacing(2),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  inputContainer: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(1),
  },
  button: {
    marginTop: theme.spacing(3),
    //   marginLeft: theme.spacing(1),
  },
});

function GrafanaConfigComponent({
  classes,
  grafanaURL,
  grafanaAPIKey,
  urlError,
  handleChange,
  handleGrafanaConfigure,
  handleChangeApiKey,
}) {
  const [availableGrafanaConnection, setAvailableGrafanaConnection] = useState([]);
  useEffect(() => {
    dataFetch(
      `/api/integrations/connections?page=0&pagesize=1&status=${encodeURIComponent(
        JSON.stringify([CONNECTION_STATES.CONNECTED]),
      )}&kind=${encodeURIComponent(JSON.stringify([CONNECTION_KINDS.GRAFANA]))}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (result) => {
        setAvailableGrafanaConnection(result?.connections);
      },
    );
  }, []);

  return (
    <NoSsr>
      <React.Fragment>
        <div className={classes.wrapper}>
          <Grid container spacing={1}>
            <Grid item xs={12} md={6}>
              <div className={classes.inputContainer}>
                <ReactSelectWrapper
                  onChange={(select) => handleChange('grafanaURL')(select)}
                  options={availableGrafanaConnection.map((connection) => ({
                    value: connection?.metadata?.url,
                    label: connection?.metadata?.url,
                    ...connection,
                  }))}
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
                onKeyDown={(e) => e.key == 'Enter' && handleGrafanaConfigure()}
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
              disabled={!CAN(keys.CONNECT_METRICS.action, keys.CONNECT_METRICS.subject)}
            >
              Submit
            </Button>
          </div>
        </div>
      </React.Fragment>
    </NoSsr>
  );
}

GrafanaConfigComponent.propTypes = {
  classes: PropTypes.object.isRequired,
  grafanaURL: PropTypes.object.isRequired,
  grafanaAPIKey: PropTypes.string.isRequired,
  handleChange: PropTypes.func.isRequired,
  handleGrafanaConfigure: PropTypes.func.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaConfigComponent);
