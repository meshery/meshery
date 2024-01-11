import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr, Grid, Button } from '@material-ui/core';
import ReactSelectWrapper from '../../ReactSelectWrapper';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useEffect } from 'react';
import { CONNECTION_KINDS, CONNECTION_STATES } from '@/utils/Enum';
import dataFetch from 'lib/data-fetch';

const promStyles = (theme) => ({
  promRoot: {
    padding: theme.spacing(5),
    backgroundColor: theme.palette.secondary.elevatedComponents,
    borderBottomLeftRadius: theme.spacing(1),
    borderBottomRightRadius: theme.spacing(1),
    marginTop: theme.spacing(2),
  },
  buttons: { display: 'flex', justifyContent: 'flex-end' },
  button: {
    marginTop: theme.spacing(3),
    //   marginLeft: theme.spacing(1),
  },
});

// change this to display all connected prometheuses connecion and based on the selection updat tht erduc prom object
const PrometheusConfigComponent = ({
  classes,
  prometheusURL,
  urlError,
  handleChange,
  handlePrometheusConfigure,
}) => {
  const [availablePrometheusConnection, setAvailablePrometheusConnection] = useState([]);

  useEffect(() => {
    dataFetch(
      `/api/integrations/connections?page=0&pagesize=1&status=${encodeURIComponent(
        JSON.stringify([CONNECTION_STATES.CONNECTED]),
      )}&kind=${encodeURIComponent(JSON.stringify([CONNECTION_KINDS.PROMETHEUS]))}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (result) => {
        setAvailablePrometheusConnection(result?.connections);
      },
    );
  }, []);

  return (
    <NoSsr>
      <React.Fragment>
        <div className={classes.promRoot}>
          <Grid container spacing={1}>
            <Grid item xs={12}>
              <ReactSelectWrapper
                onChange={(select) => handleChange('prometheusURL')(select)}
                options={availablePrometheusConnection.map((connection) => ({
                  value: connection?.metadata?.url,
                  label: connection?.metadata?.url,
                  ...connection,
                }))}
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
              disabled={!CAN(keys.CONNECT_METRICS.action, keys.CONNECT_METRICS.subject)}
            >
              Submit
            </Button>
          </div>
        </div>
      </React.Fragment>
    </NoSsr>
  );
};

PrometheusConfigComponent.propTypes = {
  classes: PropTypes.object.isRequired,
  prometheusURL: PropTypes.object.isRequired,
  handleChange: PropTypes.func.isRequired,
  handlePrometheusConfigure: PropTypes.func.isRequired,
  options: PropTypes.array.isRequired,
};

export default withStyles(promStyles)(PrometheusConfigComponent);
