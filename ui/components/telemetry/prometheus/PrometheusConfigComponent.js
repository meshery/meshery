import React from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { Grid2, Button, styled } from '@layer5/sistent';
import ReactSelectWrapper from '../../ReactSelectWrapper';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { CONNECTION_KINDS, CONNECTION_STATES } from '@/utils/Enum';
import { useGetConnectionsQuery } from '@/rtk-query/connection';

const StyledRoot = styled('div')(({ theme }) => ({
  padding: theme.spacing(5),
  backgroundColor: theme.palette.background.card,
  borderBottomLeftRadius: theme.spacing(1),
  borderBottomRightRadius: theme.spacing(1),
  marginTop: theme.spacing(2),
}));

const ButtonContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  justifyContent: 'flex-end',
  '& .submitButton': {
    marginTop: theme.spacing(3),
  },
}));

// change this to display all connected prometheuses connecion and based on the selection updat tht erduc prom object
const PrometheusConfigComponent = ({ urlError, handleChange, handlePrometheusConfigure }) => {
  const { data } = useGetConnectionsQuery({
    page: 0,
    pagesize: 1,
    status: JSON.stringify([CONNECTION_STATES.CONNECTED]),
    kind: JSON.stringify([CONNECTION_KINDS.PROMETHEUS]),
  });

  const availablePrometheusConnection = data?.connections || [];

  return (
    <NoSsr>
      <StyledRoot>
        <Grid2 container spacing={1} size="grow">
          <Grid2 size={{ xs: 12 }}>
            <ReactSelectWrapper
              onChange={(select) => handleChange('prometheusURL')(select)}
              options={availablePrometheusConnection.map((connection) => ({
                value: connection?.metadata?.url,
                label: connection?.metadata?.url,
                ...connection,
              }))}
              label="Prometheus Base URL"
              placeholder="Address of Prometheus Server"
              noOptionsMessage="No Prometheus servers discovered"
              error={urlError}
            />
          </Grid2>
        </Grid2>
        <ButtonContainer>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={handlePrometheusConfigure}
            className="submitButton"
            disabled={!CAN(keys.CONNECT_METRICS.action, keys.CONNECT_METRICS.subject)}
          >
            Submit
          </Button>
        </ButtonContainer>
      </StyledRoot>
    </NoSsr>
  );
};

PrometheusConfigComponent.propTypes = {
  prometheusURL: PropTypes.object.isRequired,
  handleChange: PropTypes.func.isRequired,
  handlePrometheusConfigure: PropTypes.func.isRequired,
  options: PropTypes.array.isRequired,
};

export default PrometheusConfigComponent;
