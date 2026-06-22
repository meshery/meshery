import React from 'react';
import { useGetUserProviderRolesQuery } from '@/rtk-query/providerRoles';
import {
  StyledBox,
  StyledChip,
  ErrorSectionContent,
  OrgNameDisabled,
  StyledTypographyDisabled,
  HeaderContainer,
} from './styles';
import {
  NoSsr,
  ErrorBoundary,
  CustomTooltip,
  InfoCircleIcon,
  useTheme,
  IconButton,
} from '@sistent/sistent';
import OrgIcon from '@/assets/icons/OrgIcon';
import CustomErrorFallback from '../../shared/ErrorBoundary/ErrorBoundary';
import { useGetSelectedOrganization, useGetLoggedInUserQuery } from '@/rtk-query/user';

const CurrentSessionInfo = () => {
  const { selectedOrganization, isLoading: isLoadingSelectedOrg } = useGetSelectedOrganization();
  const {
    data: userData,
    isLoading: isUserLoading,
    isError: isUserError,
  } = useGetLoggedInUserQuery();

  const currentOrgWithRoles = selectedOrganization?.id
    ? userData?.organizations?.organizationsWithRoles?.find(
        (org) => org.id === selectedOrganization.id,
      )
    : null;
  const userRoles = currentOrgWithRoles?.roleNames || [];

  const {
    data: providerRolesRes,
    // isSuccess: isProviderRolesSuccess,
    // isError: isProviderRolesError,
    // error: providerRolesError,
  } = useGetUserProviderRolesQuery();

  const theme = useTheme();

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
        <HeaderContainer>
          <StyledTypographyDisabled variant="h6" component="h6">
            Organization Role(s)
          </StyledTypographyDisabled>
          <CustomTooltip title="Organization Roles define your permissions within the selected organization (e.g., Organization Admin, Workspace Admin, Team Admin, User).">
            <IconButton size="small" style={{ padding: 4 }} aria-label="Organization roles help">
              <InfoCircleIcon
                fill={theme.palette.icon?.secondary || theme.palette.text?.secondary}
                width="16"
                height="16"
              />
            </IconButton>
          </CustomTooltip>
        </HeaderContainer>
        <StyledBox>
          {isLoadingSelectedOrg || isUserLoading
            ? 'Loading roles…'
            : isUserError
              ? 'Unable to load roles'
              : userRoles.length > 0
                ? userRoles.map((role) => <StyledChip key={role} label={role} />)
                : 'No roles found'}
        </StyledBox>
      </div>
      <div>
        <HeaderContainer>
          <StyledTypographyDisabled variant="h6" component="h6">
            Provider Role(s)
          </StyledTypographyDisabled>
          <CustomTooltip title="Provider Roles define your global roles and access rights across Meshery Cloud (e.g., Admin, MeshMap User, Curator).">
            <IconButton size="small" style={{ padding: 4 }} aria-label="Provider roles help">
              <InfoCircleIcon
                fill={theme.palette.icon?.secondary || theme.palette.text?.secondary}
                width="16"
                height="16"
              />
            </IconButton>
          </CustomTooltip>
        </HeaderContainer>
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
