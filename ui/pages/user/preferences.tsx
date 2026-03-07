import UserPreferences from '../../components/UserPreferences';
import { getPath } from '../../lib/path';
import Head from 'next/head';
import { promisifiedDataFetch } from '../../lib/data-fetch';
import { ctxUrl } from '../../utils/multi-ctx';
import React, { useEffect, useState } from 'react';
import { NoSsr } from '@sistent/sistent';
import { useDispatch, useSelector } from 'react-redux';
import { updatePage } from '@/store/slices/mesheryUi';

const UserPref = () => {
  const dispatch = useDispatch();
  const [anonymousStats, setAnonymousStats] = useState(undefined);
  const [perfResultStats, setPerfResultStats] = useState(undefined);
  const { selectedK8sContext } = useSelector((state) => state.ui);

  useEffect(() => {
    handleFetchData(selectedK8sContext);
  }, [selectedK8sContext]);

  useEffect(() => {
    dispatch(updatePage({ path: getPath(), title: 'User Preferences' }));
  }, []);

  const handleFetchData = async (selectedK8sContexts) => {
    try {
      const result = await promisifiedDataFetch(ctxUrl('/api/user/prefs', selectedK8sContexts), {
        method: 'GET',
        credentials: 'include',
      });
      if (result) {
        setAnonymousStats(result.anonymousUsageStats);
        setPerfResultStats(result.anonymousPerfResults);
      }
    } catch (error) {
      console.error(error);
    }
  };

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
