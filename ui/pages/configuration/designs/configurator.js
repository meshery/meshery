import React, { useEffect } from 'react';
import { NoSsr } from '@layer5/sistent';
import Head from 'next/head';
import { getPath } from '../../../lib/path';
import DesignConfigurator from '../../../components/configuratorComponents/MeshModel';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath, updateTitle } from '@/store/slices/mesheryUi';

function DesignConfiguratorPage() {
  const dispatch = useDispatchRtk();
  useEffect(() => {
    dispatch(updatePagePath({ path: getPath() }));
    dispatch(updateTitle({ title: 'Configure Design' }));
  }, []);

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
