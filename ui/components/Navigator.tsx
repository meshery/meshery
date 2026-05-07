import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import {
  CustomTooltip,
  ListItemIcon,
  Grow,
  ListItem,
  List,
  Collapse,
  Box,
  NoSsr,
  Zoom,
  HelpOutlinedIcon,
  LeftArrowIcon,
  ExternalLinkIcon as IconExternalLink,
  OpenInNewIcon,
  RemoveIcon,
  useTheme,
  SlackIcon,
  FileIcon,
  GithubIcon,
  DiscussForumIcon,
} from '@sistent/sistent';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import { cursorNotAllowed, disabledStyle } from '../css/disableComponent.styles';
import { CapabilitiesRegistry } from '../utils/disabledComponents';
import {
  CONFIGURATION,
  DASHBOARD,
  CATALOG,
  LIFECYCLE,
  SERVICE_MESH,
  TOGGLER,
} from '../constants/navigator';
import { iconSmall } from '../css/icons.styles';
import CAN from '@/utils/can';
import { CustomTextTooltip } from './MesheryMeshInterface/PatternService/CustomTextTooltip';
import {
  HideScrollbar,
  LinkContainer,
  ListIconSide,
  MainListIcon,
  MainLogo,
  MainLogoCollapsed,
  MainLogoText,
  MainLogoTextCollapsed,
  NavigatorList,
  NavigatorListItem,
  NavigatorListItemII,
  NavigatorListItemIII,
  RootDiv,
  SecondaryDivider,
  SideBarListItem,
  SideBarText,
  StyledListItem,
  NavigatorLink,
  NavigatorHelpIcons,
  HelpListItem,
  HelpButton,
  ChevronButtonWrapper,
  FixedSidebarFooter,
  SidebarDrawer,
  ExpandMore,
} from './General/style';
import { useMediaQuery } from '@sistent/sistent';
import { getProviderCapabilities, getSystemVersion } from '@/rtk-query/user';
import { useDispatch, useSelector } from 'react-redux';
import {
  toggleDrawer,
  updateBetaBadge,
  updateCapabilities,
  updateExtensionType,
  updateTitle,
} from '@/store/slices/mesheryUi';
import { useRouter } from 'next/router';
import { setAdapter } from '@/store/slices/adapter';
import { getNavigatorComponents } from './navigatorComponents';

const externalLinkIconStyle = { width: '17.76px', fontSize: '1.11rem' };

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
  const { meshAdapters, meshAdaptersts } = useSelector((state) => state.adapter);
  const dispatch = useDispatch();
  const { capabilitiesRegistry, catalogVisibility } = useSelector((state) => state.ui);
  const theme = useTheme();
  const router = useRouter();
  const [state, setState] = useState({
    path: '',
    meshAdapters,
    mts: new Date(),
    navigator: ExtensionPointSchemaValidator('navigator')(),
    showHelperButton: false,
    openItems: [],
    hoveredId: null,
    capabilitiesRegistryObj: null,
    versionDetail: {
      build: '',
      latest: '',
      outdated: false,
      commitsha: '',
      release_channel: 'NA',
    },
    navigatorComponents: [],
  });

  const updateState = (updates) => {
    setState((prevState) => ({
      ...prevState,
      ...updates,
    }));
  };

  useEffect(() => {
    fetchCapabilities();
    fetchVersionDetails();
  }, []);

  useEffect(() => {
    const path = window.location.pathname;
    if (meshAdaptersts > state.mts) {
      updateState({
        meshAdapters,
        mts: meshAdaptersts,
      });
    }

    const fetchNestedPathAndTitle = (path, title, href, children, isBeta) => {
      if (href === path) {
        dispatch(updateTitle({ title }));
        dispatch(updateBetaBadge({ isBeta }));

        return;
      }
      if (children && children.length > 0) {
        children.forEach(({ title, href, children, isBeta }) => {
          fetchNestedPathAndTitle(path, title, href, children, isBeta);
        });
      }
    };

    state.navigatorComponents.forEach(({ title, href, children, isBeta }) => {
      fetchNestedPathAndTitle(path, title, href, children, isBeta);
    });
    updateState({ path });
  }, [meshAdapters, meshAdaptersts, window.location.pathname, dispatch]);

  const ExternalLinkIcon = (
    <IconExternalLink
      {...externalLinkIconStyle}
      transform="shrink-7"
      fill={theme.palette.icon.default}
    />
  );

  const externlinks = [
    {
      id: 'doc',
      href: 'https://docs.meshery.io',
      title: 'Documentation',
      icon: <FileIcon height="20px" width="20px" />,
      hovericon: (
        <FileIcon fill={theme.palette.background.brand.default} height="20px" width="20px" />
      ),
      external_icon: ExternalLinkIcon,
    },
    {
      id: 'community',
      href: 'https://slack.meshery.io',
      title: 'Community',
      icon: (
        <SlackIcon
          primaryColor="currentColor"
          secondaryColor="currentColor"
          tertiaryColor="currentColor"
          quaternaryColor="currentColor"
          height={20}
          width={20}
        />
      ),
      hovericon: <SlackIcon height={20} width={20} />,
      external_icon: ExternalLinkIcon,
    },
    {
      id: 'forum',
      href: 'https://meshery.io/community#community-forums',
      title: 'Discussion Forum',
      icon: <DiscussForumIcon fill="currentColor" height="28px" width="28px" />,
      hovericon: <DiscussForumIcon height="28px" width="28px" />,
      external_icon: ExternalLinkIcon,
    },
    {
      id: 'issues',
      href: 'https://github.com/meshery/meshery/issues/new/choose',
      title: 'Issues',
      icon: <GithubIcon />,
      hovericon: <GithubIcon orgIcon />,
      external_icon: ExternalLinkIcon,
    },
  ];

  const fetchCapabilities = async () => {
    const { data: result, isSuccess, isError, error } = await getProviderCapabilities();

    if (isSuccess) {
      const capabilitiesRegistryObj = new CapabilitiesRegistry(result);
      const navigatorComponents = createNavigatorComponents(capabilitiesRegistryObj);
      updateState({
        navigator: ExtensionPointSchemaValidator('navigator')(result?.extensions?.navigator),
        capabilitiesRegistryObj,
        navigatorComponents,
      });
      dispatch(updateCapabilities({ capabilitiesRegistry: result }));
    }
    if (isError) {
      console.error('Error fetching capabilities', error);
    }
  };

  const fetchVersionDetails = async () => {
    const { data: result, isSuccess, isError, error } = await getSystemVersion();
    if (isSuccess) {
      updateState({
        versionDetail: result || {
          build: 'Unknown',
          latest: 'Unknown',
          outdated: false,
          commitsha: 'Unknown',
        },
      });
    }
    if (isError) {
      console.error('Error fetching version details', error);
    }
  };

  const createNavigatorComponents = (capabilityRegistryObj) => {
    return getNavigatorComponents(capabilityRegistryObj, theme);
  };

  const handleTitleClick = () => {
    router.push('/');
  };

  const handleAdapterClick = (id, link) => {
    dispatch(setAdapter({ selectedAdapter: id }));
    if (id != -1 && !link) {
      router.push('/management');
    }
  };

  const toggleMiniDrawer = () => {
    dispatch(toggleDrawer({ isDrawerCollapsed: !isDrawerCollapsed }));
  };

  const toggleSpacing = () => {
    updateState({ showHelperButton: !state.showHelperButton });
  };

  const toggleItemCollapse = (itemId) => {
    const isItemOpen = state.openItems.includes(itemId);
    const activeItems = [...state.openItems];
    if (isItemOpen) {
      updateState({ openItems: activeItems.filter((item) => item !== itemId) });
    } else {
      activeItems.push(itemId);
      updateState({ openItems: [itemId] });
    }
  };

  const updatenavigatorComponentsMenus = () => {
    const { navigatorComponents } = state;
    navigatorComponents.forEach((cat, ind) => {
      if (cat.id === LIFECYCLE) {
        cat.children.forEach((catc, ind1) => {
          if (catc.id == SERVICE_MESH) {
            return;
          }
          const icon = pickIcon(catc.id, catc.href);
          navigatorComponents[ind].children[ind1].icon = icon;

          const cr = fetchChildren(catc.id);
          navigatorComponents[ind].children[ind1].children = cr;
        });
      }

      if (cat.id === 'Configuration') {
        let show = false;
        cat.children?.forEach((ch) => {
          if (ch.id === 'Designs') {
            const idx = capabilitiesRegistry?.capabilities?.findIndex(
              (cap) => cap.feature === 'persist-meshery-patterns',
            );
            if (idx != -1) {
              ch.show = true;
              show = true;
            }
          }
        });

        cat.show = show;
      }

      //To Toggle Catalog Extension
      if (cat.id === CONFIGURATION) {
        cat.children?.forEach((ch) => {
          if (ch.id === CATALOG) {
            ch.show = catalogVisibility;
          }
        });
      }
    });
  };

  const updateAdaptersLink = () => {
    const { navigatorComponents } = state;
    navigatorComponents.forEach((cat, ind) => {
      if (cat.id === LIFECYCLE) {
        cat.children.forEach((catc, ind1) => {
          if (
            typeof navigatorComponents[ind].children[ind1].children[0] !== 'undefined' &&
            typeof navigatorComponents[ind].children[ind1].children[0].href !== 'undefined'
          ) {
            const val = true;
            const newhref = `${navigatorComponents[ind].children[ind1].children[0].href}`;
            navigatorComponents[ind].children[ind1].link = val;
            navigatorComponents[ind].children[ind1].href = newhref;
          }
        });
      }
    });
  };

  const fetchChildren = (category) => {
    const { meshAdapters } = state;
    const children = [];
    category = category.toLowerCase();
    meshAdapters.forEach((adapter) => {
      let aName = adapter.name.toLowerCase();
      if (category !== aName) {
        return;
      }
      children.push({
        id: adapter.adapter_location,
        icon: <RemoveIcon />,
        href: `/management?adapter=${adapter.adapter_location}`,
        title: `Management - ${adapter.adapter_location}`,
        link: true,
        show: true,
      });
    });
    return children;
  };

  const pickIcon = (aName, href) => {
    aName = aName.toLowerCase();
    let image = '/static/img/meshery-logo.png';
    let filter =
      window.location.pathname === href
        ? 'invert(50%) sepia(30%) saturate(1000%) hue-rotate(120deg)'
        : '';
    let logoIcon = <img src={image} style={{ width: '20px' }} />;
    if (aName) {
      image = '/static/img/' + aName + '-light.svg';
      logoIcon = <img src={image} style={{ filter: filter, width: '20px' }} />;
    }
    return logoIcon;
  };

  const getMesheryVersionText = () => {
    const { build, outdated, release_channel } = state.versionDetail;

    if (release_channel === 'edge' && outdated) return `${build}`;
    if (release_channel === 'edge' && !outdated) return `${release_channel}-latest`;
    if (release_channel === 'stable') return `${release_channel}-${build}`;

    return `${build}`;
  };

  const versionUpdateMsg = () => {
    const { outdated, latest } = state.versionDetail;

    if (outdated)
      return (
        <span style={{ marginLeft: '15px' }}>
          {'Update available '}
          <a
            href={`https://docs.meshery.io/project/releases/${latest}`}
            target="_blank"
            rel="noreferrer"
            style={{ color: 'white' }}
          >
            <CustomTextTooltip
              title={`Newer version of Meshery available: ${latest}`}
              placement="right"
            >
              <OpenInNewIcon
                fill={theme.palette.background.constant.white}
                style={{ width: '0.85rem', verticalAlign: 'middle' }}
              />
            </CustomTextTooltip>
          </a>
        </span>
      );

    return <span style={{ marginLeft: '15px' }}>Running latest</span>;
  };

  const openReleaseNotesInNew = () => {
    const { release_channel, build } = state.versionDetail;

    if (release_channel === 'edge')
      return (
        <a
          href="https://docs.meshery.io/project/releases"
          target="_blank"
          rel="noreferrer"
          style={{ color: 'white' }}
        >
          <OpenInNewIcon
            fill={theme.palette.background.constant.white}
            style={{ width: '0.85rem', verticalAlign: 'middle' }}
          />
        </a>
      );

    return (
      <a
        href={`https://docs.meshery.io/project/releases/${build}`}
        target="_blank"
        rel="noreferrer"
        style={{ color: 'white' }}
      >
        <OpenInNewIcon
          fill={theme.palette.background.constant.white}
          style={{ width: '0.85rem', verticalAlign: 'middle' }}
        />
      </a>
    );
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

  const renderChildren = (idname, children, depth) => {
    const { path } = state;
    updatenavigatorComponentsMenus();
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

  const linkContent = (iconc, titlec, hrefc, linkc, drawerCollapsed) => {
    const updatedIcon = React.cloneElement(iconc, {
      fill: state.path === hrefc ? theme.palette.icon.brand : theme.palette.common.white,
    });
    let linkContent = (
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
      linkContent = <Link href={hrefc}>{linkContent}</Link>;
    }
    return linkContent;
  };
  const { isDrawerCollapsed } = useSelector((state) => state.ui);
  const Title = (
    <div
      style={
        !state.capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD])
          ? cursorNotAllowed
          : {}
      }
    >
      <>
        <StyledListItem
          component="a"
          onClick={handleTitleClick}
          disableLogo={!state.capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD])}
        >
          {isDrawerCollapsed ? (
            <>
              <MainLogoCollapsed src="/static/img/meshery-logo.png" onClick={handleTitleClick} />
              <MainLogoTextCollapsed
                src="/static/img/meshery-logo-text.png"
                onClick={handleTitleClick}
              />
            </>
          ) : (
            <>
              <MainLogo src="/static/img/meshery-logo.png" onClick={handleTitleClick} />
              <MainLogoText src="/static/img/meshery-logo-text.png" onClick={handleTitleClick} />
            </>
          )}
        </StyledListItem>
      </>
    </div>
  );

  const Menu = (
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

  const [hoveredHelpIcon, setHoveredHelpIcon] = React.useState(null);

  const HelpIcons = (
    <>
      <NavigatorHelpIcons
        isCollapsed={isDrawerCollapsed}
        size="large"
        orientation={isDrawerCollapsed ? 'vertical' : 'horizontal'}
      >
        {externlinks.map(({ id, icon, hovericon, title, href }, index) => {
          const isHovered = hoveredHelpIcon === id;
          return (
            <HelpListItem
              key={id}
              style={isDrawerCollapsed && !state.showHelperButton ? { display: 'none' } : {}}
              onMouseEnter={() => setHoveredHelpIcon(id)}
              onMouseLeave={() => setHoveredHelpIcon(null)}
            >
              <Grow
                in={state.showHelperButton || !isDrawerCollapsed}
                timeout={{ enter: 600 - index * 200, exit: 100 * index }}
              >
                <a
                  href={href}
                  target="_blank"
                  rel="noreferrer"
                  style={
                    isDrawerCollapsed
                      ? { display: 'flex', alignItems: 'center', justifyContent: 'center' }
                      : {}
                  }
                >
                  <CustomTextTooltip title={title} placement={isDrawerCollapsed ? 'right' : 'top'}>
                    <ListIconSide>{isHovered && hovericon ? hovericon : icon}</ListIconSide>
                  </CustomTextTooltip>
                </a>
              </Grow>
            </HelpListItem>
          );
        })}
        <ListItem key="help-button" style={{ display: isDrawerCollapsed ? 'inherit' : 'none' }}>
          <CustomTextTooltip title="Help" placement={isDrawerCollapsed ? 'right' : 'top'}>
            <HelpButton isCollapsed={isDrawerCollapsed} onClick={toggleSpacing}>
              <HelpOutlinedIcon
                style={{
                  fontSize: '1.45rem',
                  ...iconSmall,
                  color: theme.palette.background.constant.white,
                  opacity: '0.7',
                  transition: 'opacity 200ms linear',
                  '&:hover': {
                    opacity: 1,
                    background: 'transparent',
                  },
                  '&:focus': {
                    opacity: 1,
                    background: 'transparent',
                  },
                }}
              />
            </HelpButton>
          </CustomTextTooltip>
        </ListItem>
      </NavigatorHelpIcons>
    </>
  );

  const Version = (
    <ListItem
      style={{
        position: 'sticky',
        paddingLeft: 0,
        paddingRight: 0,
        color: theme.palette.background.constant.white,
        fontSize: '0.75rem',
      }}
    >
      {isDrawerCollapsed ? (
        <div style={{ textAlign: 'center', width: '100%' }}>{state.versionDetail.build}</div>
      ) : (
        <Grow
          in={!isDrawerCollapsed}
          timeout={{ enter: 800, exit: 100 }}
          style={{ textAlign: 'center', width: '100%' }}
        >
          <span>
            {getMesheryVersionText()} {'  '}
            <span style={{ cursor: 'pointer' }}>{openReleaseNotesInNew()}</span>
            {versionUpdateMsg()}
          </span>
        </Grow>
      )}
    </ListItem>
  );

  const Chevron = (
    <ChevronButtonWrapper
      isCollapsed={isDrawerCollapsed}
      style={
        state?.capabilitiesRegistryObj?.isNavigatorComponentEnabled?.([TOGGLER])
          ? {}
          : cursorNotAllowed
      }
    >
      <div
        style={
          state?.capabilitiesRegistryObj?.isNavigatorComponentEnabled?.([TOGGLER])
            ? {}
            : disabledStyle
        }
        onClick={toggleMiniDrawer}
      >
        <LeftArrowIcon
          alt="Sidebar collapse toggle"
          style={{
            cursor: 'pointer',
            verticalAlign: 'middle',
          }}
          fill={theme.palette.icon.default}
          stroke={theme.palette.icon.default}
          width="1.2rem"
          height="2.8rem"
        />
      </div>
    </ChevronButtonWrapper>
  );

  return (
    <NoSsr>
      <SidebarDrawer isCollapsed={isDrawerCollapsed} variant="permanent">
        {Title}
        {Menu}
        <FixedSidebarFooter>
          {Chevron}
          {HelpIcons}
          {Version}
        </FixedSidebarFooter>
      </SidebarDrawer>
    </NoSsr>
  );
};

export const Navigator = NavigatorWrapper;

export default Navigator;
