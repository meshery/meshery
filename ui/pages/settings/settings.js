import MesherySettings from '../../components/Settings/MesherySettings';
import { NoSsr } from '@layer5/sistent';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import React, { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';

function Settings() {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Settings' }));
  }, []);

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
