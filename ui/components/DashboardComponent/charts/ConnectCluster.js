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
      <Typography style={{ fontSize: '1.5rem', marginBottom: '1rem' }} align="center">
        No connections found in your clusters
      </Typography>
      <ConnectClustersBtn />
    </div>
  );
};

export default ConnectCluster;
