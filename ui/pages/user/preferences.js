import UserPreferences from '../../components/UserPreferences';
import { connect } from 'react-redux';
import { getPath } from '../../lib/path';
import Head from 'next/head';
import { promisifiedDataFetch } from '../../lib/data-fetch';
import { ctxUrl } from '../../utils/multi-ctx';
import React, { useEffect, useState } from 'react';
import { NoSsr } from '@layer5/sistent';
import { useDispatchRtk } from '@/store/hooks';
import { updatePagePath, updateTitle } from '@/store/slices/mesheryUi';

const UserPref = (props) => {
  const dispatch = useDispatchRtk();
  const [anonymousStats, setAnonymousStats] = useState(undefined);
  const [perfResultStats, setPerfResultStats] = useState(undefined);

  useEffect(() => {
    handleFetchData(props.selectedK8sContexts);
  }, [props.selectedK8sContext]);

  useEffect(() => {
    dispatch(updatePagePath({ path: getPath() }));
    dispatch(updateTitle({ title: 'User Preferences' }));
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

const mapStateToProps = (state) => {
  const selectedK8sContexts = state.get('selectedK8sContexts');

  return {
    selectedK8sContexts,
  };
};

export default connect(mapStateToProps, null)(UserPref);
