import React, { useEffect, useRef } from 'react';
import { FavoriteIcon, Hidden, Typography, useTheme } from '@sistent/sistent';
import Navigator from '../Navigator/Navigator';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import { useDispatch, useSelector } from 'react-redux';
import { connectionsToK8sContexts } from '@/rtk-query/transforms';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import { CONNECTION_KINDS } from '@/utils/Enum';
import { setK8sContexts, updateK8SConfig } from '@/store/slices/mesheryUi';
import { loadSelectedK8sContexts, persistSelectedK8sContexts } from '@/utils/multi-ctx';
import { StyledDrawer, StyledFooterBody, StyledFooterText } from './App.styles';

// Order-insensitive equality for context-id selections: ['a','b'] and
// ['b','a'] are the same selection.
const isSameSelection = (a: string[], b: string[]) => {
  if (a.length !== b.length) return false;
  const sortedA = [...a].sort();
  const sortedB = [...b].sort();
  return sortedA.every((id, index) => id === sortedB[index]);
};

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

  // Read the current selection through a render-synced ref so the effect can
  // consult the latest value without listing it as a dependency (the effect
  // dispatches setK8sContexts, so depending on the selection would loop).
  const selectedK8sContexts = useSelector(
    (state: { ui: { selectedK8sContexts: string[] } }) => state.ui.selectedK8sContexts,
  );
  const selectedK8sContextsRef = useRef(selectedK8sContexts);
  selectedK8sContextsRef.current = selectedK8sContexts;

  useEffect(() => {
    if (!canViewClusters) {
      return;
    }

    // Until the connections query resolves there is nothing to reconcile
    // against - proceeding would treat "still loading" as "no contexts" and
    // wipe the persisted selection with an empty one.
    if (!connectionData) {
      return;
    }

    const normalizedK8sContext = connectionsToK8sContexts(connectionData?.connections);
    const availableIds: string[] = (normalizedK8sContext?.contexts ?? []).map(
      (ctx: { id: string }) => ctx.id,
    );
    const allContexts: string[] = availableIds.length > 0 ? [...availableIds, 'all'] : [];

    // Honor the selection persisted for this browser session instead of
    // force-selecting every context on each refetch of the connections query
    // (which previously wiped the user's include/exclude choices on every
    // navigation-triggered cache refresh).
    const persisted = loadSelectedK8sContexts();
    let activeContexts = allContexts;
    if (persisted !== null && !persisted.includes('all')) {
      const valid = persisted.filter((id) => availableIds.includes(id));
      if (valid.length === availableIds.length || (valid.length === 0 && persisted.length > 0)) {
        // Every context is selected (restore the implicit 'all'), or every
        // persisted id is stale (contexts were replaced) — default to all.
        activeContexts = allContexts;
      } else {
        // Partial selection — includes the explicit "none selected" case.
        activeContexts = valid;
      }
    }

    setAppState({
      k8sContexts: normalizedK8sContext,
      activeK8sContexts: activeContexts,
    });

    dispatch(updateK8SConfig({ k8sConfig: normalizedK8sContext?.contexts ?? [] }));

    // With no contexts at all there is no selection to reconcile - leave
    // redux and storage untouched so a cluster added later starts from the
    // default all-selected view instead of an accidental explicit-empty one.
    if (availableIds.length === 0) {
      return;
    }

    // Keep the redux selection (what dashboards, deploys, and queries consume)
    // in step with the restored selection; redux boots with ['all'] and would
    // otherwise disagree with the header checkboxes after a reload. The thunk
    // also rewrites sessionStorage, resyncing it after stale ids were dropped.
    const resolvedSelection = activeContexts.includes('all') ? ['all'] : activeContexts;
    if (!isSameSelection(selectedK8sContextsRef.current, resolvedSelection)) {
      dispatch(setK8sContexts({ selectedK8sContexts: resolvedSelection }));
    } else {
      // Same selection redux-side, but storage may still hold a stale variant
      // of it (e.g. dropped ids resolved back to the implicit 'all').
      persistSelectedK8sContexts(resolvedSelection);
    }
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
