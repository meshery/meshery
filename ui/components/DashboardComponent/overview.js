import React from 'react';
import { NoSsr } from '@mui/material';
import Popup from '../Popup';
import { withRouter } from 'next/router';
import { withNotify } from '../../utils/hooks/useNotification';
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

export default withRouter(withNotify(connect(mapStateToProps, mapDispatchToProps)(Overview)));
