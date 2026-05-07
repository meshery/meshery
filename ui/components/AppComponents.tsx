import React, { useEffect } from 'react';
import { FavoriteIcon, Hidden, Typography, useTheme } from '@sistent/sistent';
import Navigator from './Navigator';
import subscribeK8sContext from './graphql/subscriptions/K8sContextSubscription';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { useDispatch, useSelector } from 'react-redux';
import { updateK8SConfig } from '@/store/slices/mesheryUi';
import { StyledDrawer, StyledFooterBody, StyledFooterText } from '../themes/App.styles';

type FooterProps = {
  capabilitiesRegistry?: { restrictedAccess?: { isMesheryUiRestricted?: boolean } } | null;
  handleMesheryCommunityClick: () => void;
};

export const Footer = ({ capabilitiesRegistry, handleMesheryCommunityClick }: FooterProps) => {
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
          {capabilitiesRegistry?.restrictedAccess?.isMesheryUIRestricted || isPlaygroundBuild ? (
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

export const KubernetesSubscription = ({ setAppState }: { setAppState: SetAppState }) => {
  const dispatch = useDispatch();

  useEffect(() => {
    if (!CAN(keys.VIEW_ALL_KUBERNETES_CLUSTERS.action, keys.VIEW_ALL_KUBERNETES_CLUSTERS.subject)) {
      return;
    }

    const subscription = subscribeK8sContext(
      (result) => {
        const allContexts: string[] = [];
        if (result.k8sContext?.contexts?.length > 0) {
          result.k8sContext.contexts.forEach((ctx: { id: string }) => allContexts.push(ctx.id));
          allContexts.push('all');
        }

        setAppState({
          k8sContexts: result.k8sContext,
          activeK8sContexts: allContexts,
        });

        dispatch(updateK8SConfig({ k8sConfig: result.k8sContext.contexts }));
      },
      {
        selector: { page: '', pageSize: '10', order: '', search: '' },
      },
    );

    const dispose = () => {
      if (subscription && typeof subscription.dispose === 'function') {
        subscription.dispose();
      }
    };
    setAppState({ disposeK8sContextSubscription: dispose });
    return dispose;
  }, [dispatch, setAppState]);

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
