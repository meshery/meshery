import React from 'react';
import { useSelector } from 'react-redux';
import { useHasPermission } from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import DefaultError from '@/components/general/error-404';
import UserPreferences from '@/components/user-preferences';
import { useGetUserPrefWithContextQuery } from '@/rtk-query/user';
import { MesheryPage } from '../../components/general/MesheryPage';

/**
 * Renders the User Preferences page inside the MesheryPage shell while gating
 * context-scoped preference fetching and UI behind user preference permissions.
 */
const UserPref = () => {
  const canViewUserPreferences = useHasPermission(Keys.ExtensibilityViewMesheryUserPreferences);
  const { selectedK8sContext } = useSelector((state) => state.ui);
  const { data: prefData } = useGetUserPrefWithContextQuery(selectedK8sContext, {
    skip: !canViewUserPreferences,
  });

  const anonymousStats = canViewUserPreferences ? prefData?.anonymousUsageStats : undefined;
  const perfResultStats = canViewUserPreferences ? prefData?.anonymousPerfResults : undefined;

  let content = null;
  if (!canViewUserPreferences) {
    content = <DefaultError permissionKey={Keys.ExtensibilityViewMesheryUserPreferences} />;
  } else if (anonymousStats !== undefined && perfResultStats !== undefined) {
    content = <UserPreferences anonymousStats={anonymousStats} perfResultStats={perfResultStats} />;
  }

  return <MesheryPage title="User Preferences">{content}</MesheryPage>;
};

export default UserPref;
