import React from 'react';
import CustomErrorMessage from '../components/ErrorPage';
import { NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import { usePageTitle } from '@/utils/hooks';

const Error = () => {
  usePageTitle('Error');

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
