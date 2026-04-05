import UserPreferences from '../../components/UserPreferences';
import { getPath } from '../../lib/path';
import Head from 'next/head';
import React, { useEffect } from 'react';
import { NoSsr } from '@sistent/sistent';
import { useDispatch, useSelector } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';
import { useGetUserPrefWithContextQuery } from '@/rtk-query/user';

const UserPref = () => {
  const dispatch = useDispatch();
  const { selectedK8sContext } = useSelector((state) => state.ui);
  const { data: prefData } = useGetUserPrefWithContextQuery(selectedK8sContext);

  const anonymousStats = prefData?.anonymousUsageStats;
  const perfResultStats = prefData?.anonymousPerfResults;

  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'User Preferences' }));
  }, []);

  return (
    <>
      {anonymousStats === undefined || perfResultStats === undefined ? (
        <div></div>
      ) : (
        <NoSsr>
          <Head>
            <title>Preferences | Meshery</title>
          </Head>

          <UserPreferences anonymousStats={anonymousStats} perfResultStats={perfResultStats} />
        </NoSsr>
      )}
    </>
  );
};

export default UserPref;
