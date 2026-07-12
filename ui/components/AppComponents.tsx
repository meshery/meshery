import React, { useEffect } from 'react';
import { FavoriteIcon, Hidden, Typography, useTheme } from '@sistent/sistent';
import Navigator from './layout/Navigator/Navigator';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import { useDispatch, useSelector } from 'react-redux';
import { connectionsToK8sContexts } from '@/rtk-query/transforms';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import { CONNECTION_KINDS } from '@/utils/Enum';
import { updateK8SConfig } from '@/store/slices/mesheryUi';
import { StyledDrawer, StyledFooterBody, StyledFooterText } from '../themes/App.styles';

type FooterProps = {
  providerCapabilities?: { restrictedAccess?: { isMesheryUiRestricted?: boolean } } | null;
  handleMesheryCommunityClick: () => void;
};

export const Footer = ({ providerCapabilities, handleMesheryCommunityClick }: FooterProps) => {
  const theme = useTheme();
  const isPlaygroundBuild = process.env.NEXT_PUBLIC_PLAYGROUND_BUILD === 'true';
  const { extensionType: extension } = useSelector((state) => state.ui);

  if (extension === 'navigator') {
    return null;
  }

  return (
    <StyledFooterBody>
      <Typography
        variant="body2"
        align="center"
        component="p"
        style={{
          color:
            theme.palette.mode === 'light'
              ? theme.palette.text.default
              : theme.palette.text.disabled,
        }}
      >
        <StyledFooterText onClick={handleMesheryCommunityClick}>
          {providerCapabilities?.restrictedAccess?.isMesheryUIRestricted || isPlaygroundBuild ? (
            'ACCESS LIMITED IN MESHERY PLAYGROUND. DEPLOY MESHERY TO ACCESS ALL FEATURES.'
          ) : (
            <>
              {' '}
              Built with{' '}
              <FavoriteIcon
                fill={theme.palette.background.brand.default}
                style={{
                  display: 'inline',
                  verticalAlign: 'bottom',
                }}
              />{' '}
              by the Meshery Community
            </>
          )}
        </StyledFooterText>
      </Typography>
    </StyledFooterBody>
  );
};

type SetAppState = (partial: Record<string, unknown>) => void;

// KubernetesSubscription keeps the app-wide k8s context list (k8sConfig) in sync
// with the user's kubernetes connections. It replaces the subscribeK8sContext
// GraphQL subscription: the list is now driven by the connections REST API
// (kind=kubernetes) and stays fresh via RTK Query cache invalidation — every
// connection mutation invalidates the `Connection_API_Connections` tag, which
// refetches this query. Everything is connection-driven.
export const KubernetesSubscription = ({ setAppState }: { setAppState: SetAppState }) => {
  const dispatch = useDispatch();
  const canViewClusters = CAN(
    Keys.IdentityAccessManagementViewAllKubernetesClusters.id,
    Keys.IdentityAccessManagementViewAllKubernetesClusters.function,
  );

  const { data: connectionData } = useGetConnectionsQuery(
    // Filter by kind via a plain repeated query param (?kind=kubernetes);
    // pageSize=all fetches every cluster in one shot.
    { kind: CONNECTION_KINDS.KUBERNETES, pageSize: 'all' },
    { skip: !canViewClusters },
  );

  useEffect(() => {
    if (!canViewClusters) {
      return;
    }

    const normalizedK8sContext = connectionsToK8sContexts(connectionData?.connections);
    const allContexts: string[] = [];
    if (normalizedK8sContext?.contexts?.length > 0) {
      normalizedK8sContext.contexts.forEach((ctx: { id: string }) => allContexts.push(ctx.id));
      allContexts.push('all');
    }

    setAppState({
      k8sContexts: normalizedK8sContext,
      activeK8sContexts: allContexts,
    });

    dispatch(updateK8SConfig({ k8sConfig: normalizedK8sContext?.contexts ?? [] }));
  }, [connectionData, canViewClusters, dispatch, setAppState]);

  return null;
};

type NavigationBarProps = {
  isDrawerCollapsed: boolean;
  mobileOpen: boolean;
  handleDrawerToggle: () => void;
  updateExtensionType: (type: string | null) => void;
  canShowNav: boolean;
};

export const NavigationBar = ({
  isDrawerCollapsed,
  mobileOpen,
  handleDrawerToggle,
  updateExtensionType,
  canShowNav,
}: NavigationBarProps) => {
  if (!canShowNav) {
    return null;
  }

  return (
    <StyledDrawer
      isDrawerCollapsed={isDrawerCollapsed}
      data-testid="navigation"
      id="left-navigation-bar"
    >
      <Hidden smUp implementation="js">
        <Navigator
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          isDrawerCollapsed={isDrawerCollapsed}
          updateExtensionType={updateExtensionType}
        />
      </Hidden>
      <Hidden xsDown implementation="css">
        <Navigator
          isDrawerCollapsed={isDrawerCollapsed}
          updateExtensionType={updateExtensionType}
        />
      </Hidden>
    </StyledDrawer>
  );
};
