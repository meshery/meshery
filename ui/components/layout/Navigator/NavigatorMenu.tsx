// @ts-nocheck
import React from 'react';
import Link from 'next/link';
import { CustomTooltip, ListItemIcon, List, Collapse, Box, Zoom, useTheme } from '@sistent/sistent';
import { cursorNotAllowed } from '../../../css/disableComponent.styles';
import { LIFECYCLE } from '../../../constants/navigator';
import CAN from '@/utils/can';
import {
  HideScrollbar,
  LinkContainer,
  MainListIcon,
  NavigatorList,
  NavigatorListItem,
  NavigatorListItemII,
  NavigatorListItemIII,
  RootDiv,
  SecondaryDivider,
  SideBarListItem,
  SideBarText,
  NavigatorLink,
  ExpandMore,
} from '../../General/style';
import { updateExtensionType } from '@/store/slices/mesheryUi';

const NavigatorMenu = ({
  state,
  updateState,
  isDrawerCollapsed,
  dispatch,
  router,
  toggleItemCollapse,
  handleAdapterClick,
  updateNavigatorComponentsMenus,
  updateAdaptersLink,
}) => {
  const theme = useTheme();

  React.useEffect(() => {
    if (state.navigatorComponents.length > 0) {
      updateNavigatorComponentsMenus();
    }
  }, [state.meshAdapters, state.capabilitiesRegistryObj]);

  const extensionPointContent = (icon, href, name, drawerCollapsed) => {
    let content = (
      <>
        <LinkContainer data-testid={name}>
          <CustomTooltip
            title={name}
            placement="right"
            disableFocusListener={!drawerCollapsed}
            disableTouchListener={!drawerCollapsed}
          >
            <MainListIcon>
              <img
                src={icon}
                style={{
                  width: '20px',
                  filter:
                    window.location.pathname === href
                      ? 'invert(50%) sepia(30%) saturate(1000%) hue-rotate(120deg)'
                      : '',
                }}
                onMouseOver={(e) => {
                  e.target.style.transform = 'translate(-20%, -25%)';
                  e.target.style.top = '0';
                  e.target.style.right = '0';
                }}
                onMouseOut={(e) => {
                  e.target.style.transform = 'translate(0, 0)';
                  e.target.style.top = 'auto';
                  e.target.style.right = 'auto';
                }}
              />
            </MainListIcon>
          </CustomTooltip>
          <SideBarText drawerCollapsed={drawerCollapsed}>{name}</SideBarText>
        </LinkContainer>
      </>
    );

    if (href) {
      content = (
        <Link
          href={href}
          style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            width: '100%',
            height: '30px',
          }}
        >
          <Box
            width="100%"
            onClick={() => dispatch(updateExtensionType({ extensionType: 'navigator' }))}
          >
            {content}
          </Box>
        </Link>
      );
    }

    return content;
  };

  const renderNavigatorExtensions = (children, depth) => {
    const { path } = state;
    if (children && children.length > 0) {
      return (
        <NavigatorList disablePadding>
          {children.map(({ id, icon, href, title, children, show: showc }) => {
            if (typeof showc !== 'undefined' && !showc) {
              return '';
            }
            const isActive = path === href;
            return (
              <React.Fragment key={id}>
                <NavigatorListItem
                  button
                  depth={depth}
                  key={id}
                  isDrawerCollapsed={isDrawerCollapsed}
                  isActive={isActive}
                >
                  {extensionPointContent(icon, href, title, isDrawerCollapsed)}
                </NavigatorListItem>
                {renderNavigatorExtensions(children, depth + 1)}
              </React.Fragment>
            );
          })}
        </NavigatorList>
      );
    }
  };

  const linkContent = (iconc, titlec, hrefc, linkc, drawerCollapsed) => {
    const updatedIcon = React.cloneElement(iconc, {
      fill: state.path === hrefc ? theme.palette.icon.brand : theme.palette.common.white,
    });
    let content = (
      <>
        <LinkContainer>
          <CustomTooltip
            title={titlec}
            placement="right"
            disableFocusListener={!drawerCollapsed}
            disableHoverListener={!drawerCollapsed}
            disableTouchListener={!drawerCollapsed}
          >
            <MainListIcon>{updatedIcon}</MainListIcon>
          </CustomTooltip>
          <SideBarText drawerCollapsed={drawerCollapsed}>{titlec}</SideBarText>
        </LinkContainer>
      </>
    );

    if (linkc && hrefc) {
      content = <Link href={hrefc}>{content}</Link>;
    }
    return content;
  };

  const renderChildren = (idname, children, depth) => {
    const { path } = state;
    if (idname != LIFECYCLE && children && children.length > 0) {
      return (
        <>
          <List disablePadding>
            {children.map(
              ({
                id: idc,
                title: titlec,
                icon: iconc,
                href: hrefc,
                show: showc,
                link: linkc,
                children: childrenc,
                permission: permissionc,
              }) => {
                if (typeof showc !== 'undefined' && !showc) {
                  return '';
                }
                const isActive = path === hrefc;
                return (
                  <div key={idc}>
                    <NavigatorListItemII
                      button
                      data-testid={idc}
                      key={idc}
                      depth={depth}
                      isDrawerCollapsed={isDrawerCollapsed}
                      isActive={isActive}
                      onClick={() => {
                        if (linkc && hrefc) {
                          router.push(hrefc);
                        }
                      }}
                      disabled={permissionc ? !CAN(permissionc.action, permissionc.subject) : false}
                    >
                      {linkContent(iconc, titlec, hrefc, false, isDrawerCollapsed)}
                    </NavigatorListItemII>
                    {renderChildren(idname, childrenc, depth + 1)}
                  </div>
                );
              },
            )}
          </List>
        </>
      );
    }

    if (idname == LIFECYCLE) {
      if (children && children.length > 0) {
        return (
          <>
            <List disablePadding>
              {children.map(
                ({
                  id: idc,
                  title: titlec,
                  icon: iconc,
                  href: hrefc,
                  show: showc,
                  link: linkc,
                  children: childrenc,
                  permission: permissionc,
                }) => {
                  if (typeof showc !== 'undefined' && !showc) {
                    return '';
                  }
                  const isActive = path === hrefc;
                  return (
                    <div key={idc} style={!showc ? cursorNotAllowed : null}>
                      <NavigatorListItemIII
                        component="a"
                        data-testid={idc}
                        button
                        key={idc}
                        depth={depth}
                        isDrawerCollapsed={isDrawerCollapsed}
                        isActive={isActive}
                        isShow={!showc}
                        onClick={() => {
                          handleAdapterClick(idc, linkc);
                          if (linkc && hrefc) {
                            router.push(hrefc);
                          }
                        }}
                        disabled={
                          permissionc ? !CAN(permissionc.action, permissionc.subject) : false
                        }
                      >
                        {linkContent(iconc, titlec, hrefc, false, isDrawerCollapsed)}{' '}
                      </NavigatorListItemIII>
                      {renderChildren(idname, childrenc, depth + 1)}
                    </div>
                  );
                },
              )}
            </List>
          </>
        );
      }
      if (children && children.length === 1) {
        updateAdaptersLink();
      }
    }
    return '';
  };

  return (
    <>
      <HideScrollbar disablePadding>
        {state.navigatorComponents.map(
          ({
            id: childId,
            title,
            icon,
            href,
            show,
            link,
            children,
            hovericon,
            submenu,
            permission,
          }) => {
            return (
              <RootDiv key={childId}>
                <SideBarListItem
                  button={!!link}
                  dense
                  key={childId}
                  link={!!link}
                  isActive={state.path === href}
                  isShow={!show}
                  onClick={() => toggleItemCollapse(childId)}
                  onMouseOver={() =>
                    isDrawerCollapsed ? updateState({ hoveredId: childId }) : null
                  }
                  onMouseLeave={() =>
                    !submenu || !state.openItems.includes(childId)
                      ? updateState({ hoveredId: false })
                      : null
                  }
                  disabled={permission ? !CAN(permission.action, permission.subject) : false}
                >
                  <Link href={link ? href : ''}>
                    <NavigatorLink data-testid={childId}>
                      <CustomTooltip
                        title={childId}
                        placement="right"
                        disableFocusListener={!isDrawerCollapsed}
                        disableHoverListener={true}
                        disableTouchListener={!isDrawerCollapsed}
                        TransitionComponent={Zoom}
                      >
                        {isDrawerCollapsed &&
                        (state.hoveredId === childId ||
                          (state.openItems.includes(childId) && submenu)) ? (
                          <div>
                            <CustomTooltip
                              title={title}
                              placement="right"
                              TransitionComponent={Zoom}
                            >
                              <ListItemIcon
                                onClick={() => toggleItemCollapse(childId)}
                                style={{ marginLeft: '20%', marginBottom: '0.4rem' }}
                              >
                                {hovericon}
                              </ListItemIcon>
                            </CustomTooltip>
                          </div>
                        ) : (
                          <MainListIcon>{icon}</MainListIcon>
                        )}
                      </CustomTooltip>
                      <SideBarText drawerCollapsed={isDrawerCollapsed}>{title}</SideBarText>
                    </NavigatorLink>
                  </Link>
                  <ExpandMore
                    onClick={() => toggleItemCollapse(childId)}
                    isCollapsed={state.openItems.includes(childId)}
                    isDrawerCollapsed={isDrawerCollapsed}
                    theme={theme}
                    hasChildren={!!children}
                  />
                </SideBarListItem>
                <Collapse
                  in={state.openItems.includes(childId)}
                  style={{
                    backgroundColor: theme.palette.background.tabs,
                    opacity: '100%',
                  }}
                >
                  {renderChildren(childId, children, 1)}
                </Collapse>
              </RootDiv>
            );
          },
        )}
        {state.navigator && state.navigator.length ? (
          <React.Fragment>
            <SecondaryDivider />
            {renderNavigatorExtensions(state.navigator, 1)}
          </React.Fragment>
        ) : null}
        <SecondaryDivider />
      </HideScrollbar>
    </>
  );
};

export default NavigatorMenu;
