// @ts-nocheck
import React, { useEffect } from 'react';
import { NoSsr, useMediaQuery } from '@sistent/sistent';
import { useDispatch, useSelector } from 'react-redux';
import { toggleDrawer } from '@/store/slices/mesheryUi';
import { SidebarDrawer } from '../../General/style';
import { useNavigatorState } from './useNavigatorState';
import NavigatorHeader from './NavigatorHeader';
import NavigatorMenu from './NavigatorMenu';
import NavigatorFooter from './NavigatorFooter';

const NavigatorWrapper = () => {
  const isMobile = useMediaQuery('(max-width:599px)');
  const dispatch = useDispatch();
  const { isDrawerCollapsed } = useSelector((state) => state.ui);
  useEffect(() => {
    if (isMobile && !isDrawerCollapsed) {
      dispatch(toggleDrawer({ isDrawerCollapsed: true }));
    }
  }, [isMobile]);

  return <Navigator_ />;
};

const Navigator_ = () => {
  const {
    state,
    updateState,
    router,
    dispatch,
    isDrawerCollapsed,
    externlinks,
    handleTitleClick,
    handleAdapterClick,
    toggleMiniDrawer,
    toggleSpacing,
    toggleItemCollapse,
    updateNavigatorComponentsMenus,
    updateAdaptersLink,
  } = useNavigatorState();

  return (
    <NoSsr>
      <SidebarDrawer isCollapsed={isDrawerCollapsed} variant="permanent">
        <NavigatorHeader
          capabilitiesRegistryObj={state.capabilitiesRegistryObj}
          isDrawerCollapsed={isDrawerCollapsed}
          handleTitleClick={handleTitleClick}
        />
        <NavigatorMenu
          state={state}
          updateState={updateState}
          isDrawerCollapsed={isDrawerCollapsed}
          dispatch={dispatch}
          router={router}
          toggleItemCollapse={toggleItemCollapse}
          handleAdapterClick={handleAdapterClick}
          updateNavigatorComponentsMenus={updateNavigatorComponentsMenus}
          updateAdaptersLink={updateAdaptersLink}
        />
        <NavigatorFooter
          isDrawerCollapsed={isDrawerCollapsed}
          state={state}
          externlinks={externlinks}
          toggleSpacing={toggleSpacing}
          toggleMiniDrawer={toggleMiniDrawer}
        />
      </SidebarDrawer>
    </NoSsr>
  );
};

export const Navigator = NavigatorWrapper;

export default Navigator;
