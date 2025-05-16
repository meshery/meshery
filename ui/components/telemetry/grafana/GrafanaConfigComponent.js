import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { TextField, Grid, Button, styled } from '@layer5/sistent';
import ReactSelectWrapper from '../../ReactSelectWrapper';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { CONNECTION_KINDS, CONNECTION_STATES } from '@/utils/Enum';
import dataFetch from 'lib/data-fetch';

const Wrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(5),
  backgroundColor: theme.palette.background.card,
  borderBottomLeftRadius: theme.spacing(1),
  borderBottomRightRadius: theme.spacing(1),
  marginTop: theme.spacing(2),
}));

const ButtonContainer = styled('div')({
  display: 'flex',
  justifyContent: 'flex-end',
});

const InputContainer = styled('div')(({ theme }) => ({
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(1),
}));

const StyledButton = styled(Button)(({ theme }) => ({
  marginTop: theme.spacing(3),
}));

function GrafanaConfigComponent({
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
        <Wrapper>
          <Grid container spacing={1}>
            <Grid item xs={12} md={6}>
              <InputContainer>
                <ReactSelectWrapper
                  onChange={(select) => handleChange('grafanaURL')(select)}
                  options={availableGrafanaConnection.map((connection) => ({
                    value: connection?.metadata?.url,
                    label: connection?.metadata?.url,
                    ...connection,
                  }))}
                  label="Grafana Base URL"
                  data-testid="grafana-base-url"
                  error={urlError}
                  placeholder="Address of Grafana Server"
                  noOptionsMessage="No Grafana servers discovered"
                />
              </InputContainer>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                id="grafanaAPIKey"
                name="grafanaAPIKey"
                data-testid="grafana-api-key"
                label="API Key"
                fullWidth
                value={grafanaAPIKey}
                margin="normal"
                variant="outlined"
                onKeyDown={(e) => e.key === 'Enter' && handleGrafanaConfigure()}
                onChange={handleChangeApiKey}
              />
            </Grid>
          </Grid>
          <ButtonContainer>
            <StyledButton
              type="submit"
              variant="contained"
              color="primary"
              size="large"
              onClick={handleGrafanaConfigure}
              disabled={!CAN(keys.CONNECT_METRICS.action, keys.CONNECT_METRICS.subject)}
            >
              Submit
            </StyledButton>
          </ButtonContainer>
        </Wrapper>
      </React.Fragment>
    </NoSsr>
  );
}

GrafanaConfigComponent.propTypes = {
  grafanaURL: PropTypes.object.isRequired,
  grafanaAPIKey: PropTypes.string.isRequired,
  handleChange: PropTypes.func.isRequired,
  handleGrafanaConfigure: PropTypes.func.isRequired,
};

export default GrafanaConfigComponent;
