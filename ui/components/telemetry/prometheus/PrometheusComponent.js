import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { Typography, styled, Box } from '@layer5/sistent';
import {
  useConfigureConnectionMutation,
  useUpdateConnectionMutation,
} from '@/rtk-query/connection';
import { useGetPrometheusConfigQuery } from '@/rtk-query/telemetry';
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

  const [configurePrometheus] = useConfigureConnectionMutation();

  const submitPrometheusConfigure = async (url) => {
    updateProgress({ showProgress: true });
    try {
      await configurePrometheus({
        connectionKind: CONNECTION_KINDS.PROMETHEUS,
        body: { prometheusURL: url },
      }).unwrap();

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
    } catch (err) {
      handleError(err?.data?.error || 'Error communicating with Prometheus');
    } finally {
      updateProgress({ showProgress: false });
    }
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
      setConnectionID(data?.id);
      setConnectionName(data?.name);

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

  const [updateConnection] = useUpdateConnectionMutation();

  const handlePrometheusChipDelete = async (e) => {
    e.preventDefault();
    if (!connectionID) return;

    updateProgress({ showProgress: true });
    try {
      await updateConnection({
        connectionKind: CONNECTION_KINDS.PROMETHEUS,
        connectionPayload: { [connectionID]: CONNECTION_STATES.DISCOVERED },
      }).unwrap();

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
    } catch (err) {
      handleError(err?.data?.error || 'Failed to update Prometheus connection');
    } finally {
      updateProgress({ showProgress: false });
    }
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

  const { data: prometheusConfig, isSuccess: isPrometheusConfigLoaded } =
    useGetPrometheusConfigQuery(undefined, { skip: !props.isMeshConfigured });

  useEffect(() => {
    if (!isPrometheusConfigLoaded) return;

    updateProgress({ showProgress: false });

    if (prometheusConfig?.prometheusURL) {
      setPrometheusURL(prometheusConfig.prometheusURL);
      setPrometheusConfigSuccess(true);
      setSelectedPrometheusBoardsConfigs(prometheusConfig.selectedPrometheusBoardsConfigs || []);
    }

    const selector = {
      type: 'ALL_MESH',
      k8sClusterIDs: getK8sClusterIds(),
    };

    fetchAvailableAddons(selector).subscribe({
      next: (res) => {
        res?.addonsState?.forEach((addon) => {
          if (
            addon.name === 'prometheus' &&
            prometheusURL === '' &&
            prometheusConfig?.prometheusURL === ''
          ) {
            const newURL = `http://${addon.endpoint}`;
            setPrometheusURL(newURL);
            submitPrometheusConfigure(newURL);
          }
        });
      },
      error: (err) => console.error('Error registering Prometheus:', err),
    });
  }, [
    isPrometheusConfigLoaded,
    prometheusConfig,
    props.isMeshConfigured,
    prometheusURL,
    submitPrometheusConfigure,
    selectedK8sContexts,
    k8sConfig,
  ]);

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
