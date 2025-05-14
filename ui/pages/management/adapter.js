import { NoSsr } from '@layer5/sistent';
import Head from 'next/head';
import React, { useEffect } from 'react';
import MesheryPlayComponent from '../../components/MesheryPlayComponent';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';
import { getPath } from 'lib/path';

const Manage = () => {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Adapter' }));
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Adapter | Meshery </title>
      </Head>
      <MesheryPlayComponent />
    </NoSsr>
  );
};

export default Manage;
