import React from 'react';
import UserPreferences from '@/components/user-preferences';
import { useSelector } from 'react-redux';
import { useGetUserPrefWithContextQuery } from '@/rtk-query/user';
import { MesheryPage } from '../../components/general/MesheryPage';

const UserPref = () => {
  const { selectedK8sContext } = useSelector((state) => state.ui);
  const { data: prefData } = useGetUserPrefWithContextQuery(selectedK8sContext);

  const anonymousStats = prefData?.anonymousUsageStats;
  const perfResultStats = prefData?.anonymousPerfResults;

  if (anonymousStats === undefined || perfResultStats === undefined) {
    return null;
  }

  return (
    <MesheryPage title="User Preferences">
      <UserPreferences anonymousStats={anonymousStats} perfResultStats={perfResultStats} />
    </MesheryPage>
  );
};

export default UserPref;
