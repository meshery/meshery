import MesherySettings from '../../components/Settings/MesherySettings';
import { NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React from 'react';
import { usePageTitle } from '@/utils/hooks';

function Settings() {
  usePageTitle('Settings');

  return (
    <NoSsr>
      <Head>
        <title>Settings | Meshery</title>
      </Head>
      <MesherySettings />
    </NoSsr>
  );
}

export default Settings;
