import { NoSsr } from '@layer5/sistent';
import Head from 'next/head';
import React from 'react';
import MesheryPlayComponent from '../../components/MesheryPlayComponent';

const Manage = () => {
  return (
    <NoSsr>
      <Head>
        <title>Management | Meshery </title>
      </Head>
      <MesheryPlayComponent />
    </NoSsr>
  );
};

export default Manage;
