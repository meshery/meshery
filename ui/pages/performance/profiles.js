import React, { useEffect } from 'react';
import { NoSsr } from '@layer5/sistent';
import PerformanceProfiles from '../../components/MesheryPerformance/PerformanceProfiles';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { useDispatch } from 'react-redux';

function Results() {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Profiles' }));
  }, []);
  return (
    <NoSsr>
      <Head>
        <title>Performance Profiles | Meshery</title>
      </Head>
      <PerformanceProfiles />
    </NoSsr>
  );
}

export default Results;
