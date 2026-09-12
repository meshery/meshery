//@ts-check
import React from 'react';
import {
  NoSsr,
  ErrorBoundary,
  FormControl,
  FormGroup,
  MenuItem,
  Box,
  CircularProgress,
  Typography,
  useTheme,
  type PermissionKeySpec,
} from '@sistent/sistent';
import OrgIcon from 'assets/icons/OrgIcon';
import {
  OrgName,
  StyledSelect,
  StyledFormControlLabel,
  StyledTypography,
  CustomDownIcon,
} from './styles';
import CustomErrorFallback from '../../shared/ErrorBoundary/ErrorBoundary';
import { useDispatch, useSelector } from 'react-redux';
import { setOrganization } from '@/store/slices/mesheryUi';
import {
  useUpdateSelectedOrganizationMutation,
  useGetSelectedOrganization,
} from '@/rtk-query/user';
import { useAccessibleOrgs } from '@/utils/hooks/useAccessibleOrgs';
import { useNotification } from 'utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

interface OrgSwitcherProps {
  /**
   * When supplied, the switcher only shows organizations where the user holds
   * the required permission(s). Without it, all organizations are shown
   * (backwards-compatible with existing callers).
   */
  permissionKey?: PermissionKeySpec;
}

const OrgSwitcher = ({ permissionKey }: OrgSwitcherProps) => {
  const theme = useTheme();
  const dispatch = useDispatch();
  const [updateSelectedOrg] = useUpdateSelectedOrganizationMutation();
  const { selectedOrganization } = useGetSelectedOrganization();
  const { organization } = useSelector((state) => state.ui);
  const { notify } = useNotification();

  const activeOrg = selectedOrganization || organization;

  const { accessibleOrgs, isLoading, hasNoAccessibleOrgs, isReady } =
    useAccessibleOrgs(permissionKey);

  const handleOrgSelect = async (e) => {
    const id = e.target.value;
    try {
      await updateSelectedOrg(id).unwrap();
      const selected = accessibleOrgs.find((org) => org.id === id);
      dispatch(setOrganization({ organization: selected }));
      if (typeof window !== 'undefined' && window.sessionStorage) {
        sessionStorage.removeItem('keys');
        sessionStorage.removeItem('currentWorkspace');
      }
      window.location.reload();
    } catch (err) {
      notify({
        message: 'Failed to switch organization. Please try again.',
        event_type: EVENT_TYPES.ERROR,
      });
    }
  };

  // Loading state while checking per-org permissions
  if (isLoading) {
    return (
      <FormControl fullWidth component="fieldset">
        <StyledTypography variant="h6" component="h6">
          Checking Organization Access...
        </StyledTypography>
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            py: 2,
          }}
        >
          <CircularProgress size={24} />
          <Typography variant="body2" sx={{ ml: 1, color: 'text.secondary' }}>
            Checking your permissions across organizations...
          </Typography>
        </Box>
      </FormControl>
    );
  }

  // No accessible orgs found
  if (isReady && hasNoAccessibleOrgs) {
    return (
      <FormControl fullWidth component="fieldset">
        <StyledTypography variant="h6" component="h6">
          Switch Organization
        </StyledTypography>
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            gap: 1,
            py: 1.5,
            px: 2,
            borderRadius: '10px',
            border: `1px solid ${theme.palette.divider}`,
            backgroundColor: theme.palette.background.default,
          }}
        >
          <Typography
            variant="body2"
            sx={{
              color: 'text.secondary',
              textAlign: 'center',
              fontStyle: 'italic',
            }}
          >
            You do not have the required permission in any of your organizations.
          </Typography>
        </Box>
      </FormControl>
    );
  }

  // Render the filtered switcher
  return (
    <FormControl fullWidth component="fieldset">
      <StyledTypography variant="h6" component="h6">
        Switch Organization
      </StyledTypography>
      <FormGroup>
        <StyledFormControlLabel
          key="OrgSwitcher"
          label=""
          control={
            <StyledSelect
              fullWidth
              value={activeOrg?.id ?? ''}
              onChange={handleOrgSelect}
              SelectDisplayProps={{ style: { display: 'flex' } }}
              IconComponent={CustomDownIcon}
            >
              {/* Current org shown as disabled so the dropdown is never blank */}
              {activeOrg?.id && (
                <MenuItem key={activeOrg.id} value={activeOrg.id} disabled sx={{ opacity: 0.6 }}>
                  <OrgIcon width="24" height="24" secondaryFill={theme.palette.text.secondary} />
                  <OrgName>{activeOrg.name} (current)</OrgName>
                </MenuItem>
              )}
              {accessibleOrgs?.map((org) => (
                <MenuItem key={org.id} value={org.id}>
                  <OrgIcon width="24" height="24" secondaryFill={theme.palette.icon.secondary} />
                  <OrgName>{org.name}</OrgName>
                </MenuItem>
              ))}
            </StyledSelect>
          }
        />
      </FormGroup>
    </FormControl>
  );
};

const OrgSwitcherWithErrorBoundary = ({ permissionKey }: OrgSwitcherProps) => {
  return (
    <NoSsr>
      <ErrorBoundary customFallback={CustomErrorFallback}>
        <OrgSwitcher permissionKey={permissionKey} />
      </ErrorBoundary>
    </NoSsr>
  );
};

export default OrgSwitcherWithErrorBoundary;
