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

import { useEffect, useRef, useState } from 'react';

const SLOW_SESSION_BOOTSTRAP_WARNING_MS = 10000;
const SLOW_SESSION_BOOTSTRAP_LOG_INTERVAL_MS = 15000;

const formatPendingSessionBootstrapStep = ({
  isFetchingSelectedOrg,
  isLoadingAbilities,
  hasSelectedOrganization,
}) => {
  const pendingSteps = [];

  if (isFetchingSelectedOrg) {
    pendingSteps.push('your organization');
  }

  if (isLoadingAbilities) {
    pendingSteps.push(
      hasSelectedOrganization ? 'your organization permissions' : 'your session permissions',
    );
  }

  if (pendingSteps.length === 0) {
    return 'your session details';
  }

  if (pendingSteps.length === 1) {
    return pendingSteps[0];
  }

  if (pendingSteps.length === 2) {
    return `${pendingSteps[0]} and ${pendingSteps[1]}`;
  }

  return `${pendingSteps.slice(0, -1).join(', ')}, and ${pendingSteps.at(-1)}`;
};

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
  const hasSelectedOrganization = Boolean(selectedOrganization);
  const { data: providerCapabilities } = useGetProviderCapabilitiesQuery();

  const [updatePrefs] = useUpdateSelectedOrganizationMutation();
  const [showSlowLoadingNotice, setShowSlowLoadingNotice] = useState(false);

  const { isLoading: isLoadingAbilities, error: errorLoadingAbilities } = useGetCurrentAbilities({
    id: selectedOrganizationId,
  });

  const prefUpdatedToFallback = useRef(false);
  const isLoading = isFetchingSelectedOrg || isLoadingAbilities;
  const loadingDiagnosticsRef = useRef({
    providerType: providerCapabilities?.provider_type,
    hasSelectedOrganization,
    selectedOrganizationId,
    didFallback,
    isFetchingSelectedOrg,
    isLoadingAbilities,
    errorFetchingSelectedOrg,
    errorLoadingAbilities,
  });

  useEffect(() => {
    loadingDiagnosticsRef.current = {
      providerType: providerCapabilities?.provider_type,
      hasSelectedOrganization,
      selectedOrganizationId,
      didFallback,
      isFetchingSelectedOrg,
      isLoadingAbilities,
      errorFetchingSelectedOrg,
      errorLoadingAbilities,
    };
  }, [
    providerCapabilities?.provider_type,
    hasSelectedOrganization,
    selectedOrganizationId,
    didFallback,
    isFetchingSelectedOrg,
    isLoadingAbilities,
    errorFetchingSelectedOrg,
    errorLoadingAbilities,
  ]);

  useEffect(() => {
    if (!isLoading) {
      setShowSlowLoadingNotice(false);
      return;
    }

    let intervalId;
    const logSlowBootstrap = () => {
      setShowSlowLoadingNotice(true);
      console.warn(
        '[LoadSessionGuard] Session bootstrap is still loading; continuing to wait',
        loadingDiagnosticsRef.current,
      );
    };

    const timeoutId = window.setTimeout(() => {
      logSlowBootstrap();
      intervalId = window.setInterval(logSlowBootstrap, SLOW_SESSION_BOOTSTRAP_LOG_INTERVAL_MS);
    }, SLOW_SESSION_BOOTSTRAP_WARNING_MS);

    return () => {
      window.clearTimeout(timeoutId);
      if (intervalId) {
        window.clearInterval(intervalId);
      }
    };
  }, [isLoading]);

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

  const loaderMessage = showSlowLoadingNotice
    ? `Still initializing your ${providerCapabilities?.provider_type === 'local' ? 'local provider session' : 'session'}. Waiting for ${formatPendingSessionBootstrapStep(
        {
          isFetchingSelectedOrg,
          isLoadingAbilities,
          hasSelectedOrganization,
        },
      )}. This can take longer on slower machines.`
    : undefined;

  return (
    <DynamicFullScreenLoader isLoading={isLoading} message={loaderMessage}>
      {children}
    </DynamicFullScreenLoader>
  );
};

export const LoadSessionGuard = ({ children }) => {
  return <SelectedOrganizationProvider>{children}</SelectedOrganizationProvider>;
};
