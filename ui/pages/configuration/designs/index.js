import React, { useEffect } from 'react';
import Head from 'next/head';
import { getPath } from '../../../lib/path';
import MesheryPatterns from '../../../components/MesheryPatterns/MesheryPatterns';
import { NoSsr } from '@layer5/sistent';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';

function Patterns() {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Designs' }));
  }, []);

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
