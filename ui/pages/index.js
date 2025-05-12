import { NoSsr } from '@layer5/sistent';
import Head from 'next/head';
import React from 'react';
import DashboardComponent from '../components/DashboardComponent';
import { getPath } from '../lib/path';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath } from '@/store/slices/mesheryUi';

function Index() {
  const dispatch = useDispatchRtk();
  dispatch(updatePagePath({ path: getPath() }));

  return (
    <NoSsr>
      <Head>
        <title>Dashboard | Meshery</title>
      </Head>
      <DashboardComponent />
    </NoSsr>
  );
}

export default Index;
