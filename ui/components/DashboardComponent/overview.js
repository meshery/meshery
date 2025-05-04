import React from 'react';
import { NoSsr } from '@layer5/sistent';
import { withRouter } from 'next/router';
import { withNotify } from '../../utils/hooks/useNotification';
import { connect, Provider } from 'react-redux';
import { store } from '@/store/index';
import HoneycombComponent from './widgets/HoneyComb/HoneyCombComponent';
import { useGetMeshSyncResourceKindsQuery } from '@/rtk-query/meshsync';
import { getK8sClusterIdsFromCtxId } from '@/utils/multi-ctx';
import { bindActionCreators } from 'redux';
import { setK8sContexts, updateProgress } from 'lib/store';
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

const Overview = ({ selectedK8sContexts, k8scontext, isEditMode }) => {
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
      <div
        style={{
          background: theme.palette.background.default,
        }}
      >
        <HoneycombRoot>
          <ConnectCluster message="No clusters available. Please connect your clusters to proceed." />
        </HoneycombRoot>
      </div>
    );
  }

  const isClusterLoading = isFetching || isLoading;

  if (isError) {
    return (
      <NoSsr>
        <ErrorDisplay theme={theme} />
      </NoSsr>
    );
  }

  return (
    <NoSsr>
      <Provider store={store}>
        <HoneycombComponent
          kinds={clusterSummary?.kinds}
          isClusterLoading={isClusterLoading}
          isClusterIdsEmpty={isClusterIdsEmpty}
          isEditMode={isEditMode}
        />
      </Provider>
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
