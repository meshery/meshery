import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { Typography, Box, styled, useTheme } from '@layer5/sistent';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import GrafanaConfigComponent from './GrafanaConfigComponent';
import GrafanaSelectionComponent from './GrafanaSelectionComponent';
import GrafanaDisplaySelection from './GrafanaDisplaySelection';
// import GrafanaCharts from './GrafanaCharts';
import { updateGrafanaConfig, updateProgress } from '../../../lib/store';
import GrafanaCustomCharts from './GrafanaCustomCharts';
import fetchAvailableAddons from '../../graphql/queries/AddonsStatusQuery';
import { getK8sClusterIdsFromCtxId } from '../../../utils/multi-ctx';
import { withNotify } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import { CONNECTION_KINDS } from '@/utils/Enum';
import { withTelemetryHook } from '@/components/hooks/useTelemetryHook';
import { getCredentialByID } from '@/api/credentials';
import { useUpdateConnectionMutation } from '@/rtk-query/connection';
import {
  useConfigureGrafanaMutation,
  useGetGrafanaConfigQuery,
  useLazyGetGrafanaBoardsQuery,
  useUpdateGrafanaBoardsMutation,
} from '@/rtk-query/telemetry';
import useDebouncedCallback from '@/utils/hooks/useDebounce';

const StyledChartTitle = styled(Typography)(({ theme }) => ({
  marginLeft: theme.spacing(3),
  marginTop: theme.spacing(2),
  textAlign: 'center',
}));

const GrafanaChartsWrapper = styled(Box)(() => {
  const theme = useTheme();
  return {
    backgroundColor: theme.palette.background.default,
    padding: theme.spacing(1),
    marginTop: theme.spacing(2),
  };
});

const GrafanaComponent = (props) => {
  const [state, setState] = useState({
    urlError: false,
    grafanaConfigSuccess: props.grafana.grafanaURL !== '',
    grafanaBoardSearch: '',
    grafanaURL: props.grafana.grafanaURL,
    grafanaAPIKey: props.grafana.grafanaAPIKey,
    grafanaBoards: props.grafana.grafanaBoards,
    connectionID: props.grafana.connectionID,
    connectionName: props.grafana.connectionName,
    selectedBoardsConfigs: props.grafana.selectedBoardsConfigs,
    ts: props.grafana.ts,
  });

  //RTK Queries: Loading the Grafana configuration on component mount
  const {
    data: grafanaConfigData,
    isError: grafanaConfigIsError,
    isLoading: grafanaConfigIsLoading,
  } = useGetGrafanaConfigQuery(undefined, { skip: !props.isMeshConfigured });

  //RTK Query: This is a lazy query, it does not run on component mount, it runs when it is triggered.
  const [triggerGetGrafanaBoards] = useLazyGetGrafanaBoardsQuery();

  //RTK Mutations: These does not use cache. why? Because they are mutations, they are not supposed to be cached.
  const [configureGrafana] = useConfigureGrafanaMutation();
  const [updateGrafanaBoards] = useUpdateGrafanaBoardsMutation();
  const [updateConnection] = useUpdateConnectionMutation();

  const updateState = (newState) => {
    setState((prev) => ({ ...prev, ...newState }));
  };

  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(props.selectedK8sContexts, props.k8sconfig);
  };

  const isValidGrafanaURL = (url) =>
    Boolean(url) &&
    (url.toLowerCase().startsWith('http://') || url.toLowerCase().startsWith('https://'));

  // Validates the URL and triggers configuration submission
  const handleGrafanaConfigure = () => {
    const { grafanaURL } = state;
    // Validate URL with regex
    if (!isValidGrafanaURL(grafanaURL)) {
      updateState({ urlError: true });
      return;
    }
    submitGrafanaConfigure();
  };

  // Debounced function to fetch boards
  //TODO: Use debounce from lodash and remove custom debounce hook
  const debouncedFetchBoards = useDebouncedCallback(async () => {
    const {
      grafanaURL,
      grafanaAPIKey,
      grafanaBoardSearch,
      selectedBoardsConfigs,
      connectionID,
      connectionName,
    } = state;
    if (!grafanaURL) return;
    props.updateProgress({ showProgress: true });

    try {
      const result = await triggerGetGrafanaBoards({ connectionID, grafanaBoardSearch }).unwrap();
      props.updateProgress({ showProgress: false });
      if (result !== undefined) {
        setState((prev) => ({ ...prev, grafanaBoards: result }));
        props.updateGrafanaConfig({
          grafana: {
            grafanaURL,
            grafanaAPIKey,
            grafanaBoardSearch,
            grafanaBoards: result,
            selectedBoardsConfigs,
            connectionName,
            connectionID,
          },
        });
      }
    } catch (error) {
      props.updateProgress({ showProgress: false });
      console.error('Error fetching Grafana boards:', error);
    }
  }, 300);

  // Submits the Grafana configuration and then fetches the boards
  const submitGrafanaConfigure = async () => {
    const { grafanaURL, grafanaAPIKey, grafanaBoards, grafanaBoardSearch, selectedBoardsConfigs } =
      state;
    if (!grafanaURL) return;

    // Build URL-encoded params (using URLSearchParams for brevity)
    const params = new URLSearchParams({ grafanaURL, grafanaAPIKey }).toString();
    props.updateProgress({ showProgress: true });

    try {
      const result = await configureGrafana({ params }).unwrap();
      console.log(result);
      props.updateProgress({ showProgress: false });

      if (result !== undefined) {
        props.notify({ message: 'Grafana was configured!', event_type: EVENT_TYPES.SUCCESS });
        updateState({ grafanaConfigSuccess: true });
        props.updateGrafanaConfig({
          grafana: {
            grafanaURL,
            grafanaAPIKey,
            grafanaBoardSearch,
            grafanaBoards,
            selectedBoardsConfigs,
          },
        });
        // Fetch boards using debounced callback
        debouncedFetchBoards();
      }
    } catch (error) {
      props.updateProgress({ showProgress: false });
      console.error('There was an error communicating with Grafana', error);
    }
  };

  const handleChange = (name) => async (event) => {
    const value = event.target?.value;

    if (name === 'grafanaURL') {
      updateState({ [name]: value });
      updateState({ urlError: false });
    }

    // For board search, update state and call debounced fetch directly
    if (name === 'grafanaBoardSearch') {
      updateState({ [name]: value });
      debouncedFetchBoards();
    }
    // Get the connection object from the event value and update configuration
    const grafanaConnectionObj = value;
    try {
      const res = await getCredentialByID(grafanaConnectionObj.credential_id);
      const grafanaCfg = {
        grafanaURL: grafanaConnectionObj?.value || '',
        grafanaAPIKey: res?.secret?.secret || '',
        grafanaBoardSearch:
          grafanaConnectionObj?.metadata?.grafanaBoardSearch || state.grafanaBoardSearch,
        grafanaBoards: grafanaConnectionObj?.metadata?.grafana_boards || [],
        selectedBoardsConfigs: grafanaConnectionObj?.metadata?.selectedBoardsConfigs,
        connectionID: grafanaConnectionObj?.id,
        connectionName: grafanaConnectionObj?.name,
      };
      props.updateGrafanaConfig({ grafana: grafanaCfg });
    } catch (error) {
      console.error('Error fetching credential by ID:', error);
    }
  };

  const handleError = (msg) => () => {
    props.updateProgress({ showProgress: false });
    props.notify({ message: msg, event_type: EVENT_TYPES.ERROR });
  };

  const handleChangeApiKey = (event) => {
    updateState({ grafanaAPIKey: event.target.value });
  };

  const handleGrafanaChipDelete = async () => {
    props.updateProgress({ showProgress: true });
    try {
      // Here we use a mutation to update the connection status.
      // No payload is needed to reset; adjust as required by your API.
      const result = await updateConnection({
        connectionKind: CONNECTION_KINDS.GRAFANA,
        connectionPayload: {},
      }).unwrap();
      props.updateProgress({ showProgress: false });
      if (result !== undefined) {
        // Clear local Grafana configuration state
        setState((prev) => ({
          ...prev,
          grafanaConfigSuccess: false,
          grafanaURL: '',
          grafanaAPIKey: '',
          grafanaBoardSearch: '',
          grafanaBoards: [],
          selectedBoardsConfigs: [],
        }));
        props.updateGrafanaConfig({
          grafana: {
            grafanaURL: '',
            grafanaAPIKey: '',
            grafanaBoardSearch: '',
            grafanaBoards: [],
            selectedBoardsConfigs: [],
          },
        });
      }
    } catch (error) {
      props.updateProgress({ showProgress: false });
      console.error('There was an error communicating with Grafana', error);
    }
  };

  const handleGrafanaClick = () => {
    props.ping(state.connectionName, state.grafanaURL, state.connectionID);
  };

  const addSelectedBoardPanelConfig = (boardsSelection) => {
    const { selectedBoardsConfigs } = state;
    if (boardsSelection && boardsSelection.panels && boardsSelection.panels.length) {
      selectedBoardsConfigs?.push(boardsSelection);
      persistBoardSelection(selectedBoardsConfigs);
    }
  };

  const deleteSelectedBoardPanelConfig = (indexes) => {
    const newSelected = [...state.selectedBoardsConfigs];
    indexes.sort();
    for (let i = indexes.length - 1; i >= 0; i--) {
      newSelected.splice(indexes[i], 1);
    }
    persistBoardSelection(newSelected);
  };

  // Function to persist the board selection using RTK Query
  const persistBoardSelection = async (selectedBoardsConfigs) => {
    const { grafanaURL, grafanaAPIKey, grafanaBoards, grafanaBoardSearch, connectionID } = state;
    console.log(selectedBoardsConfigs);
    props.updateProgress({ showProgress: true });

    try {
      // Trigger the mutation with the connectionID and selected boards payload
      const result = await updateGrafanaBoards({
        connectionID,
        selectedBoardsConfigs: JSON.stringify(selectedBoardsConfigs),
      }).unwrap();
      props.updateProgress({ showProgress: false });

      if (result !== undefined) {
        // Update local state and parent configuration
        updateState((prev) => ({ ...prev, selectedBoardsConfigs }));
        props.updateGrafanaConfig({
          grafana: {
            grafanaURL,
            grafanaAPIKey,
            grafanaBoardSearch,
            grafanaBoards,
            selectedBoardsConfigs,
          },
        });
        props.notify({
          message: 'Grafana board selection updated.',
          event_type: EVENT_TYPES.SUCCESS,
        });
      }
    } catch (error) {
      props.updateProgress({ showProgress: false });
      console.error('There was an error persisting the board selection', error);
    }
  };

  useEffect(() => {
    const { grafanaURL, grafanaAPIKey, selectedBoardsConfigs, connectionID, connectionName, ts } =
      props.grafana;
    if (props.grafana.ts > state.ts) {
      updateState({
        grafanaURL,
        grafanaAPIKey,
        selectedBoardsConfigs,
        grafanaConfigSuccess: grafanaURL !== '',
        ts,
        connectionID,
        connectionName,
      });
      debouncedFetchBoards();
    }
  }, [props.grafana.grafanaAPIKey, props.grafana.grafanaURL]);

  useEffect(() => {
    if (!props.isMeshConfigured || grafanaConfigIsLoading || grafanaConfigIsError) return;

    props.updateProgress({ showProgress: grafanaConfigIsLoading });
    if (grafanaConfigData?.grafanaURL) return;

    const selector = {
      type: 'ALL_MESH',
      k8sClusterIDs: getK8sClusterIds(),
    };

    fetchAvailableAddons(selector).subscribe({
      next: (res) => {
        res?.addonsState?.forEach((addon) => {
          if (addon.name === 'grafana' && !grafanaConfigData?.grafanaURL) {
            updateState({ grafanaURL: 'http://' + addon.endpoint });
            // Submit configuration then update the selected boards if available
            submitGrafanaConfigure().then(() => {
              updateState((prev) => ({
                ...prev,
                selectedBoardsConfigs: prev.grafanaBoards?.[2] ? [prev.grafanaBoards[2]] : [],
              }));
            });
          }
        });
      },
      error: (err) => console.error('Error registering Grafana:', err),
    });
  }, [props.isMeshConfigured, grafanaConfigData, grafanaConfigIsLoading, grafanaConfigIsError]);

  const {
    urlError,
    grafanaURL,
    grafanaConfigSuccess,
    grafanaAPIKey,
    grafanaBoards,
    grafanaBoardSearch,
    selectedBoardsConfigs,
    connectionID,
  } = state;

  if (grafanaConfigSuccess) {
    let displaySelec = null;
    if (selectedBoardsConfigs?.length > 0) {
      displaySelec = (
        <>
          <GrafanaDisplaySelection
            boardPanelConfigs={selectedBoardsConfigs}
            deleteSelectedBoardPanelConfig={deleteSelectedBoardPanelConfig}
          />
          <GrafanaChartsWrapper>
            <StyledChartTitle variant="h6" gutterBottom>
              Grafana charts
            </StyledChartTitle>
            <div style={{ padding: '0 1rem' }}>
              <GrafanaCustomCharts
                boardPanelConfigs={selectedBoardsConfigs}
                grafanaURL={grafanaURL}
                grafanaAPIKey={grafanaAPIKey}
                connectionID={connectionID}
              />
            </div>
          </GrafanaChartsWrapper>
        </>
      );
    }
    return (
      <NoSsr>
        <>
          <GrafanaSelectionComponent
            grafanaURL={grafanaURL}
            grafanaBoards={grafanaBoards}
            grafanaBoardSearch={grafanaBoardSearch}
            handleGrafanaBoardSearchChange={handleChange}
            handleGrafanaChipDelete={handleGrafanaChipDelete}
            handleGrafanaClick={handleGrafanaClick}
            addSelectedBoardPanelConfig={addSelectedBoardPanelConfig}
            handleError={handleError('There was an error communicating with Grafana')}
          />
          {displaySelec}
        </>
      </NoSsr>
    );
  }
  return (
    <NoSsr>
      <GrafanaConfigComponent
        grafanaURL={grafanaURL && { label: grafanaURL, value: grafanaURL }}
        grafanaAPIKey={grafanaAPIKey}
        urlError={urlError}
        handleChange={(name) => (value) => {
          handleChange(name)({ target: { value } });
        }}
        handleChangeApiKey={handleChangeApiKey}
        handleGrafanaConfigure={handleGrafanaConfigure}
      />
    </NoSsr>
  );
};

GrafanaComponent.propTypes = {
  scannedGrafana: PropTypes.array.isRequired,
};

const mapDispatchToProps = (dispatch) => ({
  updateGrafanaConfig: bindActionCreators(updateGrafanaConfig, dispatch),
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (st) => {
  const grafana = st.get('grafana').toJS();
  const selectedK8sContexts = st.get('selectedK8sContexts');
  const k8sconfig = st.get('k8sConfig');
  return { grafana: { ...grafana, ts: new Date(grafana.ts) }, selectedK8sContexts, k8sconfig };
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(withTelemetryHook(withNotify(GrafanaComponent), CONNECTION_KINDS.GRAFANA));
