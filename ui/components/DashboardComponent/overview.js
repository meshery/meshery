import React from 'react';
import { NoSsr } from '@material-ui/core';
import Popup from '../Popup';
import { withStyles } from '@material-ui/core/styles';
import { withRouter } from 'next/router';
import { withNotify } from '../../utils/hooks/useNotification';
import blue from '@material-ui/core/colors/blue';
import { connect, Provider } from 'react-redux';
import { store } from '@/store/index';
import HoneycombComponent from './HoneyComb/HoneyCombComponent';
import { useGetMeshSyncResourceKindsQuery } from '@/rtk-query/meshsync';
import { getK8sClusterIdsFromCtxId } from '@/utils/multi-ctx';
import { bindActionCreators } from 'redux';
import { setK8sContexts, updateProgress } from 'lib/store';
import { UsesSistent } from '../SistentWrapper';
import ConnectCluster from './charts/ConnectCluster';
import { ErrorContainer, HoneycombRoot } from './style';
import { ErrorIcon, Typography, useTheme } from '@layer5/sistent';

const styles = (theme) => ({
  datatable: {
    boxShadow: 'none',
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  link: {
    cursor: 'pointer',
    textDecoration: 'none',
  },
  metricsButton: { width: '240px' },
  alreadyConfigured: { textAlign: 'center' },
  margin: { margin: theme.spacing(1) },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': { backgroundColor: blue[500] },
    },
  },
  colorBar: {},
  colorChecked: {},
  fileLabel: { width: '100%' },
  fileLabelText: {},
  inClusterLabel: { paddingRight: theme.spacing(2) },
  alignCenter: { textAlign: 'center' },
  icon: { width: theme.spacing(2.5) },
  istioIcon: { width: theme.spacing(1.5) },
  settingsIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  addIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
  cardHeader: { fontSize: theme.spacing(2) },
  card: {
    height: '100%',
    marginTop: theme.spacing(2),
  },
  cardContent: { height: '100%' },
  redirectButton: {
    marginLeft: '-.5em',
    color: '#000',
  },
  dashboardSection: {
    backgroundColor: theme.palette.secondary.elevatedComponents,
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
    marginBottom: theme.spacing(2),
  },
  errorContainer: {
    padding: theme.spacing(2),
    textAlign: 'center',
    backgroundColor: theme.palette.secondary.elevatedComponents2,
    marginTop: '1rem',
    borderRadius: 4,
  },
  errorIcon: {
    color: theme.palette.error.main,
    fontSize: '3rem',
    marginBottom: theme.spacing(1),
  },
  errorMessage: {
    color: theme.palette.error.main,
    marginBottom: theme.spacing(1),
  },
});

const ErrorDisplay = ({ theme }) => (
  <ErrorContainer>
    <ErrorIcon fill={theme.palette.background.error.default} />
    <Typography
      variant="h6"
      sx={{
        color: theme.palette.text.error,
      }}
    >
      Unable to fetch cluster data
    </Typography>
    <Typography variant="body1">
      There was an error retrieving cluster information. Please check your connection and try again.
    </Typography>
  </ErrorContainer>
);

const Overview = ({ selectedK8sContexts, k8scontext }) => {
  const clusterIds = getK8sClusterIdsFromCtxId(selectedK8sContexts, k8scontext);
  const isClusterIdsEmpty = clusterIds.size === 0;
  const theme = useTheme();

  const {
    data: clusterSummary,
    isFetching,
    isLoading,
    isError,
  } = useGetMeshSyncResourceKindsQuery(
    {
      page: 0,
      pagesize: 'all',
      clusterIds: clusterIds,
    },
    {
      skip: isClusterIdsEmpty || clusterIds.length === 0,
    },
  );

  if (clusterIds.length === 0) {
    return (
      <UsesSistent>
        <div
          style={{
            background: theme.palette.background.default,
            marginTop: '1rem',
          }}
        >
          <HoneycombRoot>
            <ConnectCluster message="No clusters available. Please connect your clusters to proceed." />
          </HoneycombRoot>
        </div>
      </UsesSistent>
    );
  }

  const isClusterLoading = isFetching || isLoading;

  if (isError) {
    return (
      <NoSsr>
        <UsesSistent>
          <ErrorDisplay theme={theme} />
        </UsesSistent>
      </NoSsr>
    );
  }

  return (
    <NoSsr>
      <UsesSistent>
        <Popup />
        <Provider store={store}>
          <div
            style={{
              background: theme.palette.background.default,
              marginTop: '1rem',
            }}
          >
            <HoneycombComponent
              kinds={clusterSummary?.kinds}
              isClusterLoading={isClusterLoading}
              isClusterIdsEmpty={isClusterIdsEmpty}
            />
          </div>
        </Provider>
      </UsesSistent>
    </NoSsr>
  );
};

const mapStateToProps = (state) => {
  return {
    selectedK8sContexts: state.get('selectedK8sContexts'),
    k8scontext: state.get('k8sConfig'),
  };
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
  setK8sContexts: bindActionCreators(setK8sContexts, dispatch),
});

export default withStyles(styles, { withTheme: true })(
  withRouter(withNotify(connect(mapStateToProps, mapDispatchToProps)(Overview))),
);
