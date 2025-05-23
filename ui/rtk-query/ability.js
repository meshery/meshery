import { ability } from '../utils/can';
import { useGetUserKeysQuery } from './userKeys';
import _ from 'lodash';
import CustomErrorMessage from '@/components/ErrorPage';
import DefaultError from '@/components/General/error-404';
import { DynamicFullScrrenLoader } from '@/components/LoadingComponents/DynamicFullscreenLoader';
import {
  useGetSelectedOrganization,
  useGetUserPrefQuery,
  useUpdateSelectedOrganizationMutation,
  useUpdateUserPrefMutation,
} from './user';
import { useGetOrgsQuery } from './organization';
import { useEffect } from 'react';

export const useGetUserAbilities = (org, skip) => {
  const { data, ...res } = useGetUserKeysQuery(
    {
      orgId: org?.id,
    },
    {
      skip,
    },
  );

  const abilities =
    data?.keys?.map((key) => ({
      action: key.id,
      subject: _.lowerCase(key.function),
    })) || [];

  return {
    ...res,
    abilities,
  };
};

export const useGetCurrentAbilities = (org) => {
  const shouldSkip = !org || !org.id;
  const res = useGetUserAbilities(org, shouldSkip);

  if (res?.abilities) {
    ability.update(res.abilities);
  }

  return res;
};

const SelectedOrganizationProvider = ({ children }) => {
  const {
    selectedOrganization,
    didFallback,
    isLoading: isFetchingSelectedOrg,
    error: errorFetchingSelectedOrg,
  } = useGetSelectedOrganization();

  const selectedOrganizationId = selectedOrganization?.id;

  const [updatePrefs] = useUpdateSelectedOrganizationMutation();

  const { isLoading: isLoadingAbilities, error: errorLoadingAbilities } = useGetCurrentAbilities({
    id: selectedOrganizationId,
  });

  useEffect(() => {
    console.log('[loadSession] selectedOrganization', selectedOrganization);
    if (
      didFallback &&
      selectedOrganizationId &&
      !errorFetchingSelectedOrg &&
      !isFetchingSelectedOrg
    ) {
      console.log('[loadSession] setting default org');
      const res = updatePrefs(selectedOrganizationId);
      console.log('updatePrefs', res);
    }
  }, [didFallback, selectedOrganizationId, errorFetchingSelectedOrg, isFetchingSelectedOrg]);

  if (errorFetchingSelectedOrg || (!isFetchingSelectedOrg && !selectedOrganization)) {
    return (
      <>
        <DefaultError />
        <CustomErrorMessage
          message={'Error occurred while fetching your current organization'}
          showImage={false}
        />
      </>
    );
  }

  if (errorLoadingAbilities) {
    return (
      <>
        <DefaultError />
        <CustomErrorMessage
          message={
            errorLoadingAbilities.message ||
            'An error occurred while fetching your organization permissions'
          }
          showImage={false}
        />
      </>
    );
  }

  const isLoading = isFetchingSelectedOrg || isLoadingAbilities;
  return <DynamicFullScrrenLoader isLoading={isLoading}>{children}</DynamicFullScrrenLoader>;
};

export const LoadSessionGuard = ({ children }) => {
  return <SelectedOrganizationProvider>{children}</SelectedOrganizationProvider>;
};
