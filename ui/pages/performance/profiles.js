import React from 'react';
import { NoSsr } from '@layer5/sistent';
import PerformanceProfiles from '../../components/MesheryPerformance/PerformanceProfiles';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath } from '@/store/slices/mesheryUi';

function Results() {
  const dispatch = useDispatchRtk();
  dispatch(updatePagePath({ path: getPath() }));
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
