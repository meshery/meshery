import React from 'react';
import { NoSsr } from '@sistent/sistent';
import PerformanceProfiles from '../../components/Performance/PerformanceProfiles';
import Head from 'next/head';
import { usePageTitle } from '@/utils/hooks';

function Results() {
  usePageTitle('Profiles');

  return (
    <NoSsr>
      <Head>
        <title>Performance Profiles | Meshery</title>
      </Head>
      <PerformanceProfiles />
    </NoSsr>
  );
}

export default Results;
