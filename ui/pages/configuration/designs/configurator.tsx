import React from 'react';
import { NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import DesignConfigurator from '../../../components/configuratorComponents/MeshModel';
import { usePageTitle } from '@/utils/hooks';

function DesignConfiguratorPage() {
  usePageTitle('Configure Design');

  return (
    <NoSsr>
      <Head>
        <title>Designs Configurator</title>
      </Head>
      <DesignConfigurator />
    </NoSsr>
  );
}

export default DesignConfiguratorPage;
