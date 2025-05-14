import { NoSsr } from '@layer5/sistent';
import Head from 'next/head';
import React, { useEffect } from 'react';
import DashboardComponent from '../components/DashboardComponent';
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
      <DashboardComponent />
    </NoSsr>
  );
}

Index.getInitialProps = async () => {
  return {};
};

export default Index;
