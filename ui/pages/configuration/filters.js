import React, { useEffect } from 'react';
import { NoSsr, Box } from '@sistent/sistent';
import MesheryFilters from '../../components/MesheryFilters/Filters';
import Head from 'next/head';
import { getPath } from '../../lib/path';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';

function NewFilters() {
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Filters' }));
  }, []);
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
