import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { Typography, Box, styled, useTheme } from '@layer5/sistent';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import dataFetch from '../../../lib/data-fetch';
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
  const boardSearchTimeoutRef = useRef(null);

  const [state, setState] = useState({
    urlError: false,
    grafanaConfigSuccess: props.grafana.grafanaURL !== '',
    grafanaBoardSearch: '',
    grafanaURL: props.grafana.grafanaURL,
    grafanaAPIKey: props.grafana.grafanaAPIKey,
    grafanaBoards: props.grafana.grafanaBoards,
    connectionID: props.grafana.connectionID,
    connectionName: props.grafana.connectionName,
    selectedBoardsConfigs: props.grafana.selectedBoardsConfigs || [],
    ts: props.grafana.ts,
  });

  const updateState = (newState) => {
    setState((prev) => ({ ...prev, ...newState }));
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
        ts: ts,
        connectionID,
        connectionName,
      });
      getGrafanaBoardsInternal();
    }
  }, [props.grafana.grafanaAPIKey, props.grafana.grafanaURL]);

  useEffect(() => {
    if (props.isMeshConfigured) {
      dataFetch(
        '/api/telemetry/metrics/grafana/config',
        {
          method: 'GET',
          credentials: 'include',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
          },
        },
        (result) => {
          props.updateProgress({ showProgress: false });
          if (!(result !== undefined && result?.grafanaURL && result?.grafanaURL !== '')) {
            let selector = {
              type: 'ALL_MESH',
              k8sClusterIDs: getK8sClusterIds(),
            };
            fetchAvailableAddons(selector).subscribe({
              next: (res) => {
                res?.addonsState?.forEach((addon) => {
                  if (addon.name === 'grafana' && state.grafanaURL === '') {
                    updateState({ grafanaURL: 'http://' + addon.endpoint });
                    submitGrafanaConfigureInternal(() =>
                      updateState({
                        selectedBoardsConfigs: state.grafanaBoards?.[2]
                          ? [state.grafanaBoards[2]]
                          : [],
                      }),
                    );
                  }
                });
              },
              error: (err) => console.log('error registering Grafana: ' + err),
            });
          }
        },
        handleError('There was an error communicating with grafana config'),
      );
    }
  }, [props.isMeshConfigured]);

  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(props.selectedK8sContexts, props.k8sconfig);
  };

  const getGrafanaBoardsInternal = (cb = () => {}) => {
    const {
      grafanaURL,
      grafanaAPIKey,
      grafanaBoardSearch,
      selectedBoardsConfigs,
      connectionID,
      connectionName,
    } = state;
    if (!grafanaURL) {
      return;
    }
    props.updateProgress({ showProgress: true });
    dataFetch(
      `/api/telemetry/metrics/grafana/boards/${connectionID}?dashboardSearch=${grafanaBoardSearch}`,
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        props.updateProgress({ showProgress: false });
        if (result !== undefined) {
          updateState({ grafanaBoards: result });
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
          cb();
        }
      },
      handleError('There was an error communicating with Grafana'),
    );
  };

  const submitGrafanaConfigureInternal = (cb = () => {}) => {
    const { grafanaURL, grafanaAPIKey, grafanaBoards, grafanaBoardSearch, selectedBoardsConfigs } =
      state;
    if (grafanaURL === '') {
      return;
    }
    const data = {
      grafanaURL,
      grafanaAPIKey,
    };
    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
      .join('&');
    props.updateProgress({ showProgress: true });
    dataFetch(
      '/api/telemetry/metrics/grafana/config',
      {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        body: params,
      },
      (result) => {
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
          getGrafanaBoardsInternal(cb);
        }
      },
      handleError('There was an error communicating with Grafana'),
    );
  };

  const handleError = (msg) => () => {
    props.updateProgress({ showProgress: false });
    props.notify({ message: msg, event_type: EVENT_TYPES.ERROR });
  };

  const handleChange = (name) => (data) => {
    console.log(grafanaAPIKey);
    if (name === 'grafanaURL' && !!data) {
      updateState({ urlError: false });
    }
    const grafanaConnectionObj = data.target?.value;
    if (name === 'grafanaBoardSearch') {
      if (boardSearchTimeoutRef.current) clearTimeout(boardSearchTimeoutRef.current);
      updateState({ [name]: data.target.value });
      boardSearchTimeoutRef.current = setTimeout(() => getGrafanaBoardsInternal(), 500);
    }

    getCredentialByID(grafanaConnectionObj.credential_id).then((res) => {
      const grafanaCfg = {
        grafanaURL: grafanaConnectionObj?.value || '',
        grafanaAPIKey: res?.secret?.secret || '',
        grafanaBoardSearch:
          grafanaConnectionObj?.metadata?.grafanaBoardSearch || state.grafanaBoardSearch,
        grafanaBoards: grafanaConnectionObj?.metadata['grafana_boards'] || [],
        selectedBoardsConfigs: grafanaConnectionObj?.metadata['selectedBoardsConfigs'],
        connectionID: grafanaConnectionObj?.id,
        connectionName: grafanaConnectionObj?.name,
      };
      props.updateGrafanaConfig({ grafana: grafanaCfg });
    });
  };

  const handleChangeApiKey = (event) => {
    updateState({ grafanaAPIKey: event.target.value });
  };

  const handleGrafanaConfigure = () => {
    const { grafanaURL } = state;
    if (
      grafanaURL === '' ||
      !(
        grafanaURL?.toLowerCase().startsWith('http://') ||
        grafanaURL.toLowerCase().startsWith('https://')
      )
    ) {
      updateState({ urlError: true });
      return;
    }
    submitGrafanaConfigureInternal();
  };

  const handleGrafanaChipDelete = () => {
    props.updateProgress({ showProgress: true });
    dataFetch(
      `/api/integrations/connections/${CONNECTION_KINDS.GRAFANA}/status`,
      {
        method: 'PUT',
        credentials: 'include',
      },
      (result) => {
        props.updateProgress({ showProgress: false });
        if (result !== undefined) {
          updateState({
            grafanaConfigSuccess: false,
            grafanaURL: '',
            grafanaAPIKey: '',
            grafanaBoardSearch: '',
            grafanaBoards: [],
            selectedBoardsConfigs: [],
          });
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
      },
      handleError('There was an error communicating with Grafana'),
    );
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

  const persistBoardSelection = (selectedBoardsConfigs) => {
    const { grafanaURL, grafanaAPIKey, grafanaBoards, grafanaBoardSearch, connectionID } = state;
    console.log(selectedBoardsConfigs);
    props.updateProgress({ showProgress: true });
    console.log('*************************', connectionID);

    dataFetch(
      `/api/telemetry/metrics/grafana/boards/${connectionID}`,
      {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json;charset=UTF-8' },
        body: JSON.stringify(selectedBoardsConfigs),
      },
      (result) => {
        console.log(result);
        props.updateProgress({ showProgress: false });
        if (result !== undefined) {
          updateState({ selectedBoardsConfigs });
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
      },
      handleError('There was an error persisting the board selection'),
    );
  };

  const {
    urlError,
    grafanaURL,
    grafanaConfigSuccess,
    grafanaAPIKey,
    grafanaBoards,
    grafanaBoardSearch,
    selectedBoardsConfigs,
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
