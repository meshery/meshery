import React, { useEffect } from 'react';
import { NoSsr } from '@layer5/sistent';
import { connect, useDispatch } from 'react-redux';
import Head from 'next/head';
import { WorkspacesComponent } from '../../components/Lifecycle';
import { Box } from '@layer5/sistent';
import { updatePage } from '@/store/slices/mesheryUi';
import { getPath } from 'lib/path';

const Workspaces = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Workspaces' }));
  }, []);

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

export default connect(null)(Workspaces);
