import { NoSsr } from '@sistent/sistent';
import Head from 'next/head';
import React, { useEffect } from 'react';
import Dashboard from '../components/Dashboard';
import { getPath } from '../lib/path';
import { useDispatch } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';

function Index() {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'Dashboard' }));
  }, []);

  return (
    <NoSsr>
      <Head>
        <title>Dashboard | Meshery</title>
      </Head>
      <Dashboard />
    </NoSsr>
  );
}

Index.getInitialProps = async () => {
  return {};
};

export default Index;
