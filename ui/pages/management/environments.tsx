import React from 'react';
import { NoSsr, Box } from '@sistent/sistent';
import Head from 'next/head';
import { EnvironmentComponent } from '../../components/Lifecycle';
import { usePageTitle } from '@/utils/hooks';

const Environments = () => {
  usePageTitle('Environments');

  return (
    <NoSsr>
      <Head>
        <title>Environments | Meshery</title>
      </Head>
      <Box sx={{ margin: 'auto', overflow: 'hidden' }}>
        <EnvironmentComponent />
      </Box>
    </NoSsr>
  );
};

export default Environments;
