import React, { useEffect } from 'react';
import { NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import { getPath } from '../../../lib/path';
import DesignConfigurator from '../../../components/configuratorComponents/MeshModel';
import { useDispatch } from 'react-redux';
import { updatePagePath, updateTitle } from '@/store/slices/mesheryUi';

function DesignConfiguratorPage() {
  const dispatch = useDispatch();
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
