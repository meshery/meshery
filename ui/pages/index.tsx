import { NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React from 'react';
import Dashboard from '../components/Dashboard';
import { usePageTitle } from '@/utils/hooks';

function Index() {
  usePageTitle('Dashboard');

  return (
    <NoSsr>
      <Head>
        <title>Dashboard | Meshery</title>
      </Head>
      <Dashboard />
    </NoSsr>
  );
}

export default Index;
