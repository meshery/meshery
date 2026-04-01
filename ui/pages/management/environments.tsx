import React, { useEffect } from 'react';
import { NoSsr, Box } from '@sistent/sistent';
import { connect, useDispatch } from 'react-redux';
import Head from 'next/head';
import { EnvironmentComponent } from '../../components/Lifecycle';
import { updatePage } from '@/store/slices/mesheryUi';
import { getPath } from 'lib/path';

const Environments = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Environments' }));
  }, []);

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

export default connect(null)(Environments);
