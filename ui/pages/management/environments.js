import React from 'react';
import { NoSsr, Box } from '@layer5/sistent';
import { connect } from 'react-redux';
import Head from 'next/head';
import { EnvironmentComponent } from '../../components/Lifecycle';

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
