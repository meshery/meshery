import React from 'react';
import Head from 'next/head';
import { getPath } from '../../../lib/path';
import MesheryPatterns from '../../../components/MesheryPatterns';
import { NoSsr } from '@layer5/sistent';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath } from '@/store/slices/mesheryUi';

function Patterns() {
  const dispatch = useDispatchRtk();
  dispatch(updatePagePath({ path: getPath() }));

  return (
    <NoSsr>
      <Head>
        <title>Designs | Meshery</title>
      </Head>
      <MesheryPatterns />
    </NoSsr>
  );
}

export default Patterns;
