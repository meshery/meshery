import React from 'react';
import { useGetUserProviderRolesQuery } from '@/rtk-query/providerRoles';
import {
  StyledBox,
  StyledChip,
  ErrorSectionContent,
  OrgNameDisabled,
  StyledTypographyDisabled,
} from './styles';
import { NoSsr } from '@sistent/sistent';
import OrgIcon from 'assets/icons/OrgIcon';
import { ErrorBoundary } from '@sistent/sistent';
import CustomErrorFallback from '../../shared/ErrorBoundary/ErrorBoundary';
import { useGetSelectedOrganization, useGetLoggedInUserQuery } from '@/rtk-query/user';

const CurrentSessionInfo = () => {
  const { selectedOrganization } = useGetSelectedOrganization();
  const { data: userData, isLoading: isUserLoading, isError: isUserError } = useGetLoggedInUserQuery();

  const currentOrgWithRoles = userData?.organizations?.organizationsWithRoles?.find(
    (org) => org.id === selectedOrganization?.id,
  );
  const userRoles = currentOrgWithRoles?.roleNames || [];

  const {
    data: providerRolesRes,
    // isSuccess: isProviderRolesSuccess,
    // isError: isProviderRolesError,
    // error: providerRolesError,
  } = useGetUserProviderRolesQuery();

  return (
    <ErrorSectionContent>
      <div>
        <StyledTypographyDisabled variant="h6" component="h6">
          Organization
        </StyledTypographyDisabled>
        <StyledBox>
          <OrgIcon width="24" height="24" secondaryFill={'#294957'} />
          <OrgNameDisabled>{selectedOrganization?.name}</OrgNameDisabled>
        </StyledBox>
      </div>
      <div>
        <StyledTypographyDisabled variant="h6" component="h6">
          Organization Role(s)
        </StyledTypographyDisabled>
        <StyledBox>
          {userRoles.length > 0
            ? userRoles.map((role) => <StyledChip key={role} label={role} />)
            : 'No roles found'}
        </StyledBox>
      </div>
      <div>
        <StyledTypographyDisabled variant="h6" component="h6">
          Provider Role(s)
        </StyledTypographyDisabled>
        <StyledBox>
          {providerRolesRes
            ? providerRolesRes?.roleNames?.map?.((role, index) => (
                <StyledChip key={index} label={role} />
              ))
            : 'No roles found'}
        </StyledBox>
      </div>
    </ErrorSectionContent>
  );
};

const CurrentSessionInfoWithErrorBoundary = () => {
  return (
    <NoSsr>
      <ErrorBoundary customFallback={CustomErrorFallback}>
        <CurrentSessionInfo />
      </ErrorBoundary>
    </NoSsr>
  );
};

export default CurrentSessionInfoWithErrorBoundary;
