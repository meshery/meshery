import React from 'react';
import Head from 'next/head';
import MesheryPatterns from '../../../components/MesheryPatterns/MesheryPatterns';
import { NoSsr } from '@sistent/sistent';
import { usePageTitle } from '@/utils/hooks';

function Patterns() {
  usePageTitle('Designs');

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
