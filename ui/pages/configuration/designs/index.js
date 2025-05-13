import React from 'react';
import Head from 'next/head';
import { getPath } from '../../../lib/path';
import MesheryPatterns from '../../../components/MesheryPatterns';
import { NoSsr } from '@layer5/sistent';
import { useDispatch } from 'react-redux';
import { updatePagePath } from '@/store/slices/mesheryUi';

function Patterns() {
  const dispatch = useDispatch();
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
