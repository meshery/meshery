import ConnectClustersBtn from '@/components/general/ConnectClustersBtn';
import React from 'react';
import { ConnectClusterWrapper, ConnectClusterText } from '../style';

const ConnectCluster = ({ message }: { message: React.ReactNode }) => {
  return (
    <ConnectClusterWrapper>
      <ConnectClusterText variant="h5" align="center">
        {message}
      </ConnectClusterText>
      <ConnectClustersBtn />
    </ConnectClusterWrapper>
  );
};

export default ConnectCluster;
