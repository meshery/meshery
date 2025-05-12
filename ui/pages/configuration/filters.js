import React from 'react';
import { NoSsr } from '@layer5/sistent';
import MesheryFilters from '../../components/Filters';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { Box } from '@layer5/sistent';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath } from '@/store/slices/mesheryUi';

function NewFilters() {
  const dispatch = useDispatchRtk();
  dispatch(updatePagePath({ path: getPath() }));
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
