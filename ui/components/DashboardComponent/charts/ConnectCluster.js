import ConnectClustersBtn from '@/components/General/ConnectClustersBtn';
import React from 'react';
import { ConnectClusterWrapper, ConnectClusterText } from '../style';

const ConnectCluster = () => {
  return (
    <ConnectClusterWrapper>
      <ConnectClusterText variant="h5" align="center">
        No workloads found in your cluster(s).
      </ConnectClusterText>
      <ConnectClustersBtn />
    </ConnectClusterWrapper>
  );
};

export default ConnectCluster;
