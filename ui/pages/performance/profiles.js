import React, { useEffect } from 'react';
import { NoSsr } from '@sistent/sistent';
import PerformanceProfiles from '../../components/Performance/PerformanceProfiles';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';

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
