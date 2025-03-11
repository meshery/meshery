import React from 'react';
import { useGetUserOrgRolesQuery } from '@/rtk-query/orgRoles';
import { useGetUserProviderRolesQuery } from '@/rtk-query/providerRoles';
import {
  StyledBox,
  StyledChip,
  ErrorSectionContent,
  OrgNameDisabled,
  StyledTypographyDisabled,
} from './styles';
import { connect } from 'react-redux';
import { withRouter } from 'next/router';
import { Provider } from 'react-redux';
import { store } from '../../../store';
import { NoSsr } from '@layer5/sistent';
import OrgIcon from 'assets/icons/OrgIcon';
import { ErrorBoundary } from '@layer5/sistent';
import CustomErrorFallback from '../ErrorBoundary';

const CurrentSessionInfo = (props) => {
  const { organization } = props;
  const {
    data: rolesRes,
    // isSuccess: isRolesSuccess,
    // isError: isRolesError,
    // error: rolesError,
  } = useGetUserOrgRolesQuery({ orgId: organization?.id }, { skip: !organization?.id });

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
          <OrgNameDisabled>{organization?.name}</OrgNameDisabled>
        </StyledBox>
      </div>
      <div>
        <StyledTypographyDisabled variant="h6" component="h6">
          Organization Role(s)
        </StyledTypographyDisabled>
        <StyledBox>
          {rolesRes
            ? rolesRes?.roles.map((role) => <StyledChip key={role.id} label={role.role_name} />)
            : 'No roles found'}
        </StyledBox>
      </div>
      <div>
        <StyledTypographyDisabled variant="h6" component="h6">
          Provider Role(s)
        </StyledTypographyDisabled>
        <StyledBox>
          {providerRolesRes
            ? providerRolesRes?.role_names.map((role, index) => (
                <StyledChip key={index} label={role} />
              ))
            : 'No roles found'}
        </StyledBox>
      </div>
    </ErrorSectionContent>
  );
};

const CurrentSessionInfoWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary customFallback={CustomErrorFallback}>
        <Provider store={store}>
          <CurrentSessionInfo {...props} />
        </Provider>
      </ErrorBoundary>
    </NoSsr>
  );
};

const mapStateToProps = (state) => {
  const organization = state.get('organization');
  return {
    organization,
  };
};

export default connect(mapStateToProps)(withRouter(CurrentSessionInfoWithErrorBoundary));
