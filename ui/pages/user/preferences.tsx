import UserPreferences from '../../components/UserPreferences';
import Head from 'next/head';
import React from 'react';
import { NoSsr } from '@sistent/sistent';
import { useSelector } from 'react-redux';
import { useGetUserPrefWithContextQuery } from '@/rtk-query/user';
import { usePageTitle } from '@/utils/hooks';

const UserPref = () => {
  usePageTitle('User Preferences');
  const { selectedK8sContext } = useSelector((state) => state.ui);
  const { data: prefData } = useGetUserPrefWithContextQuery(selectedK8sContext);

  const anonymousStats = prefData?.anonymousUsageStats;
  const perfResultStats = prefData?.anonymousPerfResults;

  if (anonymousStats === undefined || perfResultStats === undefined) {
    return null;
  }

  return (
    <NoSsr>
      <Head>
        <title>Preferences | Meshery</title>
      </Head>
      <UserPreferences anonymousStats={anonymousStats} perfResultStats={perfResultStats} />
    </NoSsr>
  );
};

export default UserPref;
