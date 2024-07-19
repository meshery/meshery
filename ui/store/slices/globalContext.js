const { createSlice } = require('@reduxjs/toolkit');
export const globalContextSlice = createSlice({
  name: 'globalContextSlice',
  initialState: {
    selectedEnvs: {
      // envId : {
      //   // selectedConnections : [connection]
      // }
    },
    selectedOrg: {},
    selectedWorkspace: {},
  },

  reducers: {
    selectOrg: (state, action) => {
      state.selectedOrg = action.payload;
    },
    selectWorkspace: (state, action) => {
      state.selectedWorkspace = action.payload;
    },
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
      const isEnvSelected = selectIsEnvSelected({ globalContextSlice: state }, env.id);
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

export const {
  selectEnv,
  unselectEnv,
  selectConnection,
  unselectConnection,
  selectOrg,
  selectWorkspace,
} = globalContextSlice.actions;

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
export const selectCurrentOrg = (orgVal) => (dispatch) => {
  dispatch(selectOrg(orgVal));
};
export const selectCurrentWorkspace = (workspaceVal) => (dispatch) => {
  dispatch(selectWorkspace(workspaceVal));
};
export const toggleEnvSelection = (environment, selectedConnections) => (dispatch, getState) => {
  const isSelected = selectIsEnvSelected(getState(), environment.id);
  if (isSelected) {
    dispatch(unselectEnv({ envId: environment.id }));
    return;
  }
  dispatch(selectEnv({ environment, selectedConnections }));
};

export const getCurrentOrg = (state) => state.globalContext.selectedOrg;
export const getCurrentWorkspace = (state) => state.globalContext.selectedWorkspace;

export default globalContextSlice.reducer;

// selectors

const selectIsConnectionSelected = (state, envId, connectionId) => {
  if (!selectIsEnvSelected(state, envId)) {
    return false;
  }
  return state.globalContext.selectedEnvs[envId].selectedConnections
    .map((connection) => connection.id)
    .includes(connectionId);
};

const selectIsEnvSelected = (state, envId) => Boolean(state.globalContext.selectedEnvs[envId]);

const selectSelectedEnvs = (state) => state.globalContext.selectedEnvs;

const selectSelectedConnections = (state, envId) =>
  state.globalContext.selectedEnvs[envId]?.selectedConnections || [];

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

export const ORG_ACTIONS_TO_PERSIST = {
  'globalContextSlice/selectOrg': ['globalContext/selectedOrg'],
};

export const WORKSPACE_ACTIONS_TO_PERSIST = {
  'globalContextSlice/selectWorkspace': ['globalContext/selectedWorkspace'],
};
