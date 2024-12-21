import React from 'react';
import { NoSsr } from '@mui/material';
import { connect } from 'react-redux';
import Head from 'next/head';
import { EnvironmentComponent } from '../../components/Lifecycle';

import { Box } from '@mui/material';

const Environments = () => {
  return (
    <NoSsr>
      <Head>
        <title>Environments | Meshery</title>
      </Head>
      <Box sx={{ maxWidth: '90%', margin: 'auto', overflow: 'hidden' }}>
        <EnvironmentComponent />
      </Box>
    </NoSsr>
  );
};

export default connect(null)(Environments);
