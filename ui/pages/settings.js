import MesherySettings from '../components/MesherySettings';
import { NoSsr } from '@layer5/sistent';
import Head from 'next/head';
import { getPath } from '../lib/path';
import React from 'react';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath, updateTitle } from '@/store/slices/mesheryUi';

function Settings() {
  const dispatch = useDispatchRtk();
  dispatch(updatePagePath({ path: getPath() }));
  dispatch(updateTitle({ title: 'Settings' }));

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
