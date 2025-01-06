import ConnectClustersBtn from '@/components/General/ConnectClustersBtn';
import React from 'react';
import { Typography } from '@layer5/sistent';

const ConnectCluster = () => {
  return (
    <div
      style={{
        padding: '2rem',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        flexDirection: 'column',
      }}
    >
      <Typography variant="h5" align="center" style={{ marginBottom: '0.5rem' }}>
        Connect the clusters which have valid Kubernetes resources
      </Typography>
      <ConnectClustersBtn />
    </div>
  );
};

export default ConnectCluster;
