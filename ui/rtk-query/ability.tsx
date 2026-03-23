import { ability } from '../utils/can';
import { useGetUserKeysQuery } from './userKeys';
import _ from 'lodash';
import CustomErrorMessage from '@/components/ErrorPage';
import DefaultError from '@/components/General/error-404';
import { DynamicFullScreenLoader } from '@/components/LoadingComponents/DynamicFullscreenLoader';
import {
  useGetProviderCapabilitiesQuery,
  useGetSelectedOrganization,
  // useGetUserPrefQuery,
  useUpdateSelectedOrganizationMutation,
  // useUpdateUserPrefMutation,
} from './user';

import { useEffect, useRef } from 'react';

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
  const { data: providerCapabilities } = useGetProviderCapabilitiesQuery();

  const [updatePrefs] = useUpdateSelectedOrganizationMutation();

  const { isLoading: isLoadingAbilities, error: errorLoadingAbilities } = useGetCurrentAbilities({
    id: selectedOrganizationId,
  });

  const prefUpdatedToFallback = useRef(false);

  useEffect(() => {
    const isLoading = isFetchingSelectedOrg || isLoadingAbilities;
    if (!isLoading) {
      return;
    }

    const timeoutId = window.setTimeout(() => {
      console.warn('[LoadSessionGuard] Session bootstrap is still loading', {
        providerType: providerCapabilities?.provider_type,
        hasSelectedOrganization: Boolean(selectedOrganization),
        selectedOrganizationId,
        didFallback,
        isFetchingSelectedOrg,
        isLoadingAbilities,
        errorFetchingSelectedOrg,
        errorLoadingAbilities,
      });
    }, 5000);

    return () => window.clearTimeout(timeoutId);
  }, [
    providerCapabilities?.provider_type,
    selectedOrganization,
    selectedOrganizationId,
    didFallback,
    isFetchingSelectedOrg,
    isLoadingAbilities,
    errorFetchingSelectedOrg,
    errorLoadingAbilities,
  ]);

  useEffect(() => {
    if (prefUpdatedToFallback.current) {
      return;
    }
    if (
      didFallback &&
      selectedOrganizationId &&
      !errorFetchingSelectedOrg &&
      !isFetchingSelectedOrg
    ) {
      updatePrefs(selectedOrganizationId);
      prefUpdatedToFallback.current = true;
    }
  }, [didFallback, selectedOrganizationId, errorFetchingSelectedOrg, isFetchingSelectedOrg]);

  if (
    errorFetchingSelectedOrg ||
    (!isFetchingSelectedOrg && !selectedOrganization && !didFallback)
  ) {
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
  return <DynamicFullScreenLoader isLoading={isLoading}>{children}</DynamicFullScreenLoader>;
};

export const LoadSessionGuard = ({ children }) => {
  return <SelectedOrganizationProvider>{children}</SelectedOrganizationProvider>;
};
