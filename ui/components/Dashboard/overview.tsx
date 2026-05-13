import React, { useMemo } from 'react';
import { NoSsr } from '@sistent/sistent';
import HoneycombComponent from './widgets/HoneyComb/HoneyCombComponent';
import { useGetMeshSyncResourceKindsQuery } from '@/rtk-query/meshsync';
import { getK8sClusterIdsFromCtxId } from '@/utils/multi-ctx';
import ConnectCluster from './charts/ConnectCluster';
import { ErrorContainer, HoneycombRoot } from './style';
import { ErrorIcon, Typography, useTheme, Theme } from '@sistent/sistent';
import { useSelector } from 'react-redux';

const ErrorDisplay = ({ theme }: { theme: Theme }) => (
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

const Overview = ({ isEditMode }: { isEditMode?: boolean }) => {
  const { k8sConfig, selectedK8sContexts } = useSelector((state) => state.ui);
  const clusterIds = useMemo(
    () => getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig),
    [k8sConfig, selectedK8sContexts],
  );
  const hasNoClusters = clusterIds.size === 0 || clusterIds.length === 0;
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
      skip: hasNoClusters,
    },
  );

  if (hasNoClusters) {
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
      <HoneycombComponent
        kinds={clusterSummary?.kinds}
        isClusterLoading={isClusterLoading}
        isClusterIdsEmpty={hasNoClusters}
        isEditMode={isEditMode}
      />
    </NoSsr>
  );
};

export default Overview;
