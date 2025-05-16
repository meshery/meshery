import React from 'react';
import { NoSsr } from '@layer5/sistent';
import MesheryPerformanceComponent from '../../components/MesheryPerformance/Dashboard';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';
import { useEffect } from 'react';

function Performance() {
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Performance' }));
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Performance Dashboard | Meshery</title>
      </Head>
      <MesheryPerformanceComponent />
    </NoSsr>
  );
}

export default Performance;
