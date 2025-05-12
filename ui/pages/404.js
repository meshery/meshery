import React from 'react';
import CustomErrorMessage from '../components/ErrorPage';
import { NoSsr } from '@layer5/sistent';
import { getPath } from '../lib/path';
import Head from 'next/head';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath } from '@/store/slices/mesheryUi';

const Error = () => {
  const dispatch = useDispatchRtk();
  dispatch(updatePagePath({ path: getPath() }));

  return (
    <NoSsr>
      <Head>
        <title>404 - Page Not Found </title>
      </Head>
      <CustomErrorMessage />
    </NoSsr>
  );
};

export default Error;
