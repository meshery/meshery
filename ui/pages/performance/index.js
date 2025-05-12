import React from 'react';
import { NoSsr } from '@layer5/sistent';
import MesheryPerformanceComponent from '../../components/MesheryPerformance/Dashboard';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath } from '@/store/slices/mesheryUi';

function Performance() {
  const dispatch = useDispatchRtk();
  dispatch(updatePagePath({ path: getPath() }));

  return (
    <NoSsr>
      <Head>
        <title>Performance Dashboard | Meshery</title>
      </Head>
      <MesheryPerformanceComponent />
    </NoSsr>
  );
}

export default Performance;
