import React from 'react';
import { NoSsr, Box } from '@sistent/sistent';
import Head from 'next/head';
import { WorkspacesComponent } from '../../components/Lifecycle';
import { usePageTitle } from '@/utils/hooks';

const Workspaces = () => {
  usePageTitle('Workspaces');

  return (
    <NoSsr>
      <Head>
        <title>Workspaces | Meshery</title>
      </Head>
      <Box sx={{ margin: 'auto', overflow: 'hidden' }}>
        <WorkspacesComponent />
      </Box>
    </NoSsr>
  );
};

export default Workspaces;
