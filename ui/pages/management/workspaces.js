import React from 'react';
import { NoSsr } from '@mui/material';
import { connect } from 'react-redux';
import Head from 'next/head';
import { WorkspacesComponent } from '../../components/Lifecycle';
import { Box } from '@mui/material';

const Workspaces = () => {
  return (
    <NoSsr>
      <Head>
        <title>Workspaces | Meshery</title>
      </Head>
      <Box sx={{ maxWidth: '90%', margin: 'auto', overflow: 'hidden' }}>
        <WorkspacesComponent />
      </Box>
    </NoSsr>
  );
};

export default connect(null)(Workspaces);
