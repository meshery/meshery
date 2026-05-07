import React from 'react';
import { NoSsr } from '@sistent/sistent';
import MesheryFilters from '../../components/MesheryFilters/Filters';
import Head from 'next/head';
import { Box } from '@sistent/sistent';
import { usePageTitle } from '@/utils/hooks';

function NewFilters() {
  usePageTitle('Filters');

  return (
    <NoSsr>
      <Head>
        <title>Filters | Meshery</title>
      </Head>
      <Box
        sx={{
          maxWidth: '90%',
          margin: 'auto',
          overflow: 'hidden',
        }}
      >
        <MesheryFilters />
      </Box>
    </NoSsr>
  );
}

export default NewFilters;
