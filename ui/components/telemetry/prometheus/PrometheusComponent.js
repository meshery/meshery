import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { Typography, styled, Box } from '@layer5/sistent';
import dataFetch from '../../../lib/data-fetch';
import PrometheusSelectionComponent from './PrometheusSelectionComponent';
import GrafanaDisplaySelection from '../grafana/GrafanaDisplaySelection';
import GrafanaCustomCharts from '../grafana/GrafanaCustomCharts';
import PrometheusConfigComponent from './PrometheusConfigComponent';
import { getK8sClusterIdsFromCtxId } from '../../../utils/multi-ctx';
import fetchAvailableAddons from '../../graphql/queries/AddonsStatusQuery';
import { withNotify } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { CONNECTION_KINDS, CONNECTION_STATES } from '@/utils/Enum';
import { withTelemetryHook } from '@/components/hooks/useTelemetryHook';
import { useDispatch, useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';
import { updatePrometheusConfig } from '@/store/slices/telemetry';

const StyledBox = styled(Box)(({ theme }) => ({
  '& .buttons': {
    display: 'flex',
  },
  '& .button': {
    marginTop: theme.spacing(3),
  },
  '& .margin': {
    margin: theme.spacing(1),
  },
  '& .icon': {
    width: theme.spacing(2.5),
  },
  '& .alignRight': {
    textAlign: 'right',
  },
  '& .formControl': {
    margin: theme.spacing(1),
    minWidth: 180,
  },
  '& .panelChips': {
    display: 'flex',
    flexWrap: 'wrap',
  },
  '& .panelChip': {
    margin: theme.spacing(0.25),
  },
  '& .chartTitle': {
    marginLeft: theme.spacing(3),
    marginTop: theme.spacing(2),
    textAlign: 'center',
  },
}));

const PrometheusComponent = (props) => {
  const { prometheus: initialPrometheus } = useSelector((state) => state.telemetry);
  const [urlError, setUrlError] = useState(false);
  const [prometheusConfigSuccess, setPrometheusConfigSuccess] = useState(
    initialPrometheus.prometheusURL !== '',
  );
  const [selectedPrometheusBoardsConfigs, setSelectedPrometheusBoardsConfigs] = useState(
    initialPrometheus.selectedPrometheusBoardsConfigs || [],
  );
  const [prometheusURL, setPrometheusURL] = useState(initialPrometheus.prometheusURL);
  const [connectionID, setConnectionID] = useState(initialPrometheus.connectionID);
  const [connectionName, setConnectionName] = useState(initialPrometheus.connectionName);
  const { k8sConfig } = useSelector((state) => state.ui);
  const { selectedK8sContexts } = useSelector((state) => state.ui);
  const dispatch = useDispatch();

  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig);
  };

  const submitPrometheusConfigure = (url) => {
    const params = new URLSearchParams({ prometheusURL: url });
    updateProgress({ showProgress: true });
    dataFetch(
      `/api/telemetry/metrics/config`,
      {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        body: params,
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (result) {
          props.notify({
            message: 'Prometheus was configured!',
            event_type: EVENT_TYPES.SUCCESS,
          });
          setPrometheusConfigSuccess(true);
          dispatch(
            updatePrometheusConfig({
              prometheusURL: url,
              selectedPrometheusBoardsConfigs: [],
            }),
          );
        }
      },
      handleError,
    );
  };

  const isValidGrafanaURL = (url) =>
    Boolean(url) &&
    (url.toLowerCase().startsWith('http://') || url.toLowerCase().startsWith('https://'));

  const handleChange = (name) => (data) => {
    if (name === 'prometheusURL' && data?.value) {
      const newURL = data.value;
      if (!isValidGrafanaURL(newURL)) {
        setUrlError(true);
        return;
      }

      setPrometheusURL(newURL);
      dispatch(
        updatePrometheusConfig({
          prometheusURL: newURL,
          selectedPrometheusBoardsConfigs: data?.metadata?.prometheus_boards || [],
          connectionID: data?.id,
          connectionName: data?.name,
        }),
      );
    }
  };

  const handlePrometheusConfigure = () => {
    if (!prometheusURL || urlError) {
      setUrlError(true);
      return;
    }
    submitPrometheusConfigure(prometheusURL);
  };

  const handleError = (message = 'Error communicating with Prometheus') => {
    updateProgress({ showProgress: false });
    props.notify({ message, event_type: EVENT_TYPES.ERROR });
  };

  const handlePrometheusChipDelete = (e) => {
    e.preventDefault();
    if (!connectionID) return;

    updateProgress({ showProgress: true });
    dataFetch(
      `/api/integrations/connections/${CONNECTION_KINDS.PROMETHEUS}/status`,
      {
        method: 'PUT',
        credentials: 'include',
        body: JSON.stringify({ [connectionID]: CONNECTION_STATES.DISCOVERED }),
      },
      () => {
        updateProgress({ showProgress: false });
        setPrometheusConfigSuccess(false);
        setPrometheusURL('');
        setSelectedPrometheusBoardsConfigs([]);
        dispatch(
          updatePrometheusConfig({
            prometheusURL: '',
            selectedPrometheusBoardsConfigs: [],
          }),
        );
        props.notify({
          message: `Connection "${connectionID}" transitioned to discovered state`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      },
      (err) => handleError(err.message),
    );
  };

  const handlePrometheusClick = () => {
    props.ping(connectionName, prometheusURL, connectionID);
  };

  const addSelectedBoardPanelConfig = (boardsSelection) => {
    if (boardsSelection?.panels?.length) {
      const newConfigs = [...selectedPrometheusBoardsConfigs, boardsSelection];
      setSelectedPrometheusBoardsConfigs(newConfigs);
      dispatch(
        updatePrometheusConfig({
          prometheusURL,
          selectedPrometheusBoardsConfigs: newConfigs,
        }),
      );
    }
  };

  const deleteSelectedBoardPanelConfig = (indexes) => {
    const newConfigs = [...selectedPrometheusBoardsConfigs];
    indexes.sort((a, b) => b - a).forEach((i) => newConfigs.splice(i, 1));
    setSelectedPrometheusBoardsConfigs(newConfigs);
    dispatch(
      updatePrometheusConfig({
        prometheusURL,
        selectedPrometheusBoardsConfigs: newConfigs,
      }),
    );
  };

  useEffect(() => {
    if (initialPrometheus.prometheusURL !== prometheusURL) {
      setPrometheusConfigSuccess(initialPrometheus.prometheusURL !== '');
      setPrometheusURL(initialPrometheus.prometheusURL);
      setSelectedPrometheusBoardsConfigs(initialPrometheus.selectedPrometheusBoardsConfigs || []);
      setConnectionID(initialPrometheus.connectionID);
      setConnectionName(initialPrometheus.connectionName);
    }
  }, [initialPrometheus]);

  useEffect(() => {
    if (!props.isMeshConfigured) return;

    dataFetch(
      `/api/telemetry/metrics/config`,
      {
        method: 'GET',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
      },
      (result) => {
        updateProgress({ showProgress: false });
        if (result?.prometheusURL) {
          const selector = {
            type: 'ALL_MESH',
            k8sClusterIDs: getK8sClusterIds(),
          };
          fetchAvailableAddons(selector).subscribe({
            next: (res) => {
              res?.addonsState?.forEach((addon) => {
                if (addon.name === 'prometheus' && prometheusURL === '') {
                  const newURL = `http://${addon.endpoint}`;
                  setPrometheusURL(newURL);
                  submitPrometheusConfigure(newURL);
                }
              });
            },
            error: (err) => console.error('Error registering Prometheus:', err),
          });
        }
      },
      handleError,
    );
  }, []);

  if (prometheusConfigSuccess) {
    const displaySelec = selectedPrometheusBoardsConfigs.length > 0 && (
      <>
        <GrafanaDisplaySelection
          boardPanelConfigs={selectedPrometheusBoardsConfigs}
          deleteSelectedBoardPanelConfig={deleteSelectedBoardPanelConfig}
        />
        <Typography variant="h6" gutterBottom>
          Prometheus charts
        </Typography>
        <GrafanaCustomCharts
          boardPanelConfigs={selectedPrometheusBoardsConfigs}
          prometheusURL={prometheusURL}
          connectionID={connectionID}
        />
      </>
    );

    return (
      <NoSsr>
        <StyledBox>
          <PrometheusSelectionComponent
            prometheusURL={prometheusURL}
            handlePrometheusChipDelete={handlePrometheusChipDelete}
            addSelectedBoardPanelConfig={addSelectedBoardPanelConfig}
            handlePrometheusClick={handlePrometheusClick}
            handleError={handleError}
            connectionID={connectionID}
          />
          {displaySelec}
        </StyledBox>
      </NoSsr>
    );
  }

  return (
    <NoSsr>
      <StyledBox>
        <PrometheusConfigComponent
          prometheusURL={prometheusURL && { label: prometheusURL, value: prometheusURL }}
          urlError={urlError}
          handleChange={handleChange}
          handlePrometheusConfigure={handlePrometheusConfigure}
        />
      </StyledBox>
    </NoSsr>
  );
};

PrometheusComponent.propTypes = {
  scannedPrometheus: PropTypes.array.isRequired,
};

export default withTelemetryHook(withNotify(PrometheusComponent), CONNECTION_KINDS.PROMETHEUS);
