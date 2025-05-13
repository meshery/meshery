import { NoSsr } from '@layer5/sistent';
import Head from 'next/head';
import React from 'react';
import DashboardComponent from '../components/DashboardComponent';
import { getPath } from '../lib/path';
import { useDispatch } from 'react-redux';
import { updatePagePath, updateTitle } from '@/store/slices/mesheryUi';

function Index() {
  const dispatch = useDispatch();
  dispatch(updatePagePath({ path: getPath() }));
  dispatch(updateTitle({ title: 'Dashboard' }));
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
