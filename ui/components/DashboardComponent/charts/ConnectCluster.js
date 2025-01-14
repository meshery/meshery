import ConnectClustersBtn from '@/components/General/ConnectClustersBtn';
import React from 'react';
import { ConnectClusterWrapper, ConnectClusterText } from '../style';

const ConnectCluster = ({ message }) => {
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
