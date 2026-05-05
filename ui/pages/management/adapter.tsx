import { NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React from 'react';
import MesheryPlayComponent from '../../components/MesheryPlayComponent';
import { usePageTitle } from '@/utils/hooks';

const Manage = () => {
  usePageTitle('Adapter');

  return (
    <NoSsr>
      <Head>
        <title>Adapter | Meshery </title>
      </Head>
      <MesheryPlayComponent />
    </NoSsr>
  );
};

export default Manage;
