import { createSlice } from '@reduxjs/toolkit';

export const globalEnvironmentContextSlice = createSlice({
  name: 'globalEnvironmentContext',
  initialState: {
    selectedEnvs: {
      // envId : {
      //   // selectedConnections : [connection]
      // }
    },
  },

  reducers: {
    selectEnv: (state, action) => {
      const { environment, selectedConnections = [] } = action.payload;
      state.selectedEnvs[environment.id] = {
        ...environment,
        selectedConnections,
      };
    },
    unselectEnv: (state, action) => {
      const { envId } = action.payload;
      delete state.selectedEnvs[envId];
    },

    selectConnection: (state, action) => {
      const { env, connection } = action.payload;
      const isEnvSelected = selectIsEnvSelected({ globalEnvironmentContext: state }, env.id);
      if (!isEnvSelected) {
        state.selectedEnvs[env.id] = {
          ...env,
          selectedConnections: [connection],
        };
        return;
      }
      state.selectedEnvs[env.id].selectedConnections.push(connection);
    },

    unselectConnection: (state, action) => {
      const { envId, connectionId } = action.payload;
      state.selectedEnvs[envId].selectedConnections = state.selectedEnvs[
        envId
      ].selectedConnections.filter((connection) => connection.id !== connectionId);
    },
  },
});

export const { selectEnv, unselectEnv, selectConnection, unselectConnection } =
  globalEnvironmentContextSlice.actions;

export const toggleConnection = (env, connection) => (dispatch, getState) => {
  const envId = env.id;
  const connectionId = connection.id;
  const isSelected = selectIsConnectionSelected(getState(), envId, connection.id);
  if (isSelected) {
    dispatch(unselectConnection({ envId, connectionId }));
  } else {
    dispatch(selectConnection({ env, connection }));
  }
};

export const toggleEnvSelection = (environment, selectedConnections) => (dispatch, getState) => {
  const isSelected = selectIsEnvSelected(getState(), environment.id);
  if (isSelected) {
    dispatch(unselectEnv({ envId: environment.id }));
    return;
  }
  dispatch(selectEnv({ environment, selectedConnections }));
};

export default globalEnvironmentContextSlice.reducer;

// selectors

const selectIsConnectionSelected = (state, envId, connectionId) => {
  if (!selectIsEnvSelected(state, envId)) {
    return false;
  }
  return selectSelectedConnections(state, envId)
    .map((connection) => connection.id)
    .includes(connectionId);
};

// Shared fallback so selectors return a stable reference when the slice is missing
// (e.g. a component rendered under a different Redux store).
const EMPTY_SELECTED_ENVS = {};

const selectSelectedEnvs = (state) =>
  state?.globalEnvironmentContext?.selectedEnvs ?? EMPTY_SELECTED_ENVS;

const selectIsEnvSelected = (state, envId) => Boolean(selectSelectedEnvs(state)[envId]);

const selectSelectedConnections = (state, envId) =>
  selectSelectedEnvs(state)[envId]?.selectedConnections || [];

const selectAllSelectedConnections = (state) => {
  const selectedEnvs = selectSelectedEnvs(state);
  return Object.values(selectedEnvs).reduce((acc, { selectedConnections }) => {
    return [...acc, ...selectedConnections];
  }, []);
};

const selectSelectedK8sConnections = (state, envId) =>
  selectSelectedConnections(state, envId).filter((connection) => connection.kind === 'kubernetes');

const selectAllSelectedK8sConnections = (state) =>
  selectAllSelectedConnections(state).filter((connection) => connection.kind === 'kubernetes');

export {
  selectIsConnectionSelected,
  selectIsEnvSelected,
  selectSelectedEnvs,
  selectSelectedConnections,
  selectAllSelectedConnections,
  selectAllSelectedK8sConnections,
  selectSelectedK8sConnections,
};
