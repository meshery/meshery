import React from 'react';
import { NoSsr } from '@sistent/sistent';
import MesheryPerformanceComponent from '../../components/Performance/Dashboard';
import Head from 'next/head';
import { usePageTitle } from '@/utils/hooks';

function Performance() {
  usePageTitle('Performance');

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
