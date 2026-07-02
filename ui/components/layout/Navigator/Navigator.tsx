import React, { useEffect, useMemo, useState, useRef } from 'react';
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
import ExtensionPointSchemaValidator from '../../../utils/ExtensionPointSchemaValidator';
import { cursorNotAllowed, disabledStyle } from '../../../css/disableComponent.styles';
import { ProviderUiAccessControl } from '../../../utils/disabledComponents';
import {
  CONFIGURATION,
  DASHBOARD,
  CATALOG,
  LIFECYCLE,
  SERVICE_MESH,
  TOGGLER,
} from '../../../constants/navigator';
import { iconSmall } from '../../../css/icons.styles';
import CAN from '@/utils/can';
import { CustomTextTooltip } from '../../meshery-mesh-interface/PatternService/CustomTextTooltip';
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
} from '../../general/style';
import { useMediaQuery } from '@sistent/sistent';
import { getProviderCapabilities, getSystemVersion } from '@/rtk-query/user';
import { useDispatch, useSelector } from 'react-redux';
import {
  toggleDrawer,
  updateBetaBadge,
  updateExtensionType,
  updateProviderCapabilities,
  updateTitle,
} from '@/store/slices/mesheryUi';
import { useRouter } from 'next/router';
import { setAdapter } from '@/store/slices/adapter';
import { getNavigatorComponents } from './navigatorComponents';

const externalLinkIconStyle = { width: '17.76px', fontSize: '1.11rem' };

const activeIconFilter = 'invert(50%) sepia(30%) saturate(1000%) hue-rotate(120deg)';

const defaultVersionDetail = {
  build: '',
  latest: '',
  outdated: false,
  commitsha: '',
  release_channel: 'NA',
};

const findNavigatorItemByPath = (items, path) => {
  for (const item of items) {
    if (item.href === path) {
      return item;
    }

    if (item.children?.length) {
      const nestedItem = findNavigatorItemByPath(item.children, path);
      if (nestedItem) {
        return nestedItem;
      }
    }
  }

  return null;
};

const buildAdapterChildren = (meshAdapters, category) => {
  const normalizedCategory = category.toLowerCase();

  return meshAdapters.reduce((children, adapter) => {
    if (adapter.name.toLowerCase() !== normalizedCategory) {
      return children;
    }

    children.push({
      id: adapter.adapter_location,
      icon: <RemoveIcon />,
      href: `/management?adapter=${adapter.adapter_location}`,
      title: `Management - ${adapter.adapter_location}`,
      link: true,
      show: true,
    });

    return children;
  }, []);
};

const buildLifecycleIcon = (adapterName, href, currentPath) => {
  const normalizedName = adapterName?.toLowerCase();
  const image = normalizedName
    ? `/static/img/${normalizedName}-light.svg`
    : '/static/img/meshery-logo/meshery-logo.png';

  return (
    <img
      src={image}
      style={{
        filter: currentPath === href ? activeIconFilter : '',
        width: '20px',
      }}
    />
  );
};

const resolveNavigatorComponents = ({
  providerUiAccessControl,
  theme,
  meshAdapters,
  catalogVisibility,
  currentPath,
}) => {
  const designPersistenceEnabled = Boolean(
    providerUiAccessControl?.providerPayload?.providerCapabilities?.some(
      (capability) => capability.feature === 'persist-meshery-patterns',
    ),
  );

  return getNavigatorComponents(providerUiAccessControl, theme).map((category) => {
    if (category.id === LIFECYCLE) {
      return {
        ...category,
        children: category.children?.map((child) => {
          if (child.id === SERVICE_MESH) {
            return child;
          }

          return {
            ...child,
            icon: child.icon ?? buildLifecycleIcon(child.id, child.href, currentPath),
            children: buildAdapterChildren(meshAdapters, child.id),
          };
        }),
      };
    }

    if (category.id === CONFIGURATION) {
      const children = category.children?.map((child) => {
        if (child.id === 'Designs') {
          return {
            ...child,
            show: designPersistenceEnabled,
          };
        }

        if (child.id === CATALOG) {
          return {
            ...child,
            show: catalogVisibility,
          };
        }

        return child;
      });

      return {
        ...category,
        children,
      };
    }

    return category;
  });
};

const NavigatorWrapper = () => {
  const isMobile = useMediaQuery('(max-width:599px)');
  const dispatch = useDispatch();
  const { isDrawerCollapsed } = useSelector((state) => state.ui);

  useEffect(() => {
    if (isMobile && !isDrawerCollapsed) {
      dispatch(toggleDrawer({ isDrawerCollapsed: true }));
    }
    if (!isMobile && isDrawerCollapsed) {
      dispatch(toggleDrawer({ isDrawerCollapsed: false }));
    }
  }, [isMobile]);

  return <NavigatorContent />;
};

const NavigatorContent = () => {
  const isMobile = useMediaQuery('(max-width:599px)');
  const { isDrawerCollapsed } = useSelector((state) => state.ui);
  const sidebarNodeRef = useRef(null);
  const { meshAdapters } = useSelector((state) => state.adapter);
  const dispatch = useDispatch();
  const { catalogVisibility } = useSelector((state) => state.ui);
  const theme = useTheme();
  const router = useRouter();
  const currentPath = useMemo(() => (router.asPath || '').split('?')[0], [router.asPath]);
  const [navigatorExtensions, setNavigatorExtensions] = useState(() =>
    ExtensionPointSchemaValidator('navigator')(),
  );
  const [showHelperButton, setShowHelperButton] = useState(false);
  const [openItems, setOpenItems] = useState<string[]>([]);
  const [hoveredId, setHoveredId] = useState<string | null>(null);
  const [providerUiAccessControl, setProviderUiAccessControl] =
    useState<ProviderUiAccessControl | null>(null);
  const [versionDetail, setVersionDetail] = useState(defaultVersionDetail);
  const navigatorComponents = useMemo(() => {
    if (!providerUiAccessControl) {
      return [];
    }

    return resolveNavigatorComponents({
      providerUiAccessControl,
      theme,
      meshAdapters,
      catalogVisibility,
      currentPath,
    });
  }, [providerUiAccessControl, catalogVisibility, currentPath, meshAdapters, theme]);

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

  async function fetchCapabilities() {
    const { data: result, isSuccess, isError, error } = await getProviderCapabilities();

    if (isSuccess) {
      setNavigatorExtensions(
        ExtensionPointSchemaValidator('navigator')(result?.extensions?.navigator),
      );
      setProviderUiAccessControl(new ProviderUiAccessControl(result));
      dispatch(updateProviderCapabilities({ providerCapabilities: result }));
    }
    if (isError) {
      console.error('Error fetching capabilities', error);
    }
  }

  async function fetchVersionDetails() {
    const { data: result, isSuccess, isError, error } = await getSystemVersion();

    if (isSuccess) {
      setVersionDetail(
        result || {
          ...defaultVersionDetail,
          build: 'Unknown',
          latest: 'Unknown',
          commitsha: 'Unknown',
        },
      );
    }
    if (isError) {
      console.error('Error fetching version details', error);
    }
  }

  useEffect(() => {
    void Promise.all([fetchCapabilities(), fetchVersionDetails()]);
  }, []);

  useEffect(() => {
    const activeNavigatorItem = findNavigatorItemByPath(navigatorComponents, currentPath);

    if (!activeNavigatorItem) {
      return;
    }

    dispatch(updateTitle({ title: activeNavigatorItem.title }));
    dispatch(updateBetaBadge({ isBeta: activeNavigatorItem.isBeta }));
  }, [currentPath, dispatch, navigatorComponents]);

  useEffect(() => {
    if (!isMobile || isDrawerCollapsed) return;

    const handleClickOutside = (e) => {
      if (sidebarNodeRef.current && !sidebarNodeRef.current.contains(e.target)) {
        dispatch(toggleDrawer({ isDrawerCollapsed: true }));
      }
    };

    const handleRouteChange = () => {
      dispatch(toggleDrawer({ isDrawerCollapsed: true }));
    };

    document.addEventListener('mousedown', handleClickOutside);
    router.events.on('routeChangeStart', handleRouteChange);

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      router.events.off('routeChangeStart', handleRouteChange);
    };
  }, [router.events, dispatch]);

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
    setShowHelperButton((previousValue) => !previousValue);
  };

  const toggleItemCollapse = (itemId) => {
    setOpenItems((activeItems) =>
      activeItems.includes(itemId)
        ? activeItems.filter((activeItem) => activeItem !== itemId)
        : [itemId],
    );
  };

  const getMesheryVersionText = () => {
    const { build, outdated, release_channel } = versionDetail;

    if (release_channel === 'edge' && outdated) return `${build}`;
    if (release_channel === 'edge' && !outdated) return `${release_channel}-latest`;
    if (release_channel === 'stable') return `${release_channel}-${build}`;

    return `${build}`;
  };

  const versionUpdateMsg = () => {
    const { outdated, latest } = versionDetail;

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
    const { release_channel, build } = versionDetail;

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

  const handleExtensionIconMouseEnter = (event: React.MouseEvent<HTMLImageElement>) => {
    const image = event.currentTarget;

    image.style.transform = 'translate(-20%, -25%)';
    image.style.top = '0';
    image.style.right = '0';
  };

  const handleExtensionIconMouseLeave = (event: React.MouseEvent<HTMLImageElement>) => {
    const image = event.currentTarget;

    image.style.transform = 'translate(0, 0)';
    image.style.top = 'auto';
    image.style.right = 'auto';
  };

  const renderNavigatorExtensions = (children, depth) => {
    if (!children || children.length === 0) {
      return null;
    }

    // Extension point: see https://docs.meshery.io/reference/extensibility
    // for how you can add your own extension.
    const extensionItems = children
      .map(({ id, icon, href, title, children: subItems, show: showc }) => {
        if (typeof showc !== 'undefined' && !showc) {
          return null;
        }

        const isActive = currentPath === href;
        const childExtensions = renderNavigatorExtensions(subItems, depth + 1);

        return (
          <RootDiv key={id} data-testid={depth === 1 ? 'extension-nav-root-item' : undefined}>
            <NavigatorListItem
              button
              depth={depth}
              isDrawerCollapsed={isDrawerCollapsed}
              isActive={isActive}
            >
              {extensionPointContent(icon, href, title, isDrawerCollapsed)}
            </NavigatorListItem>
            {childExtensions}
          </RootDiv>
        );
      })
      .filter(Boolean);

    if (extensionItems.length === 0) {
      return null;
    }

    if (depth === 1) {
      return extensionItems;
    }

    return (
      <NavigatorList disablePadding data-testid="extension-nav-submenu">
        {extensionItems}
      </NavigatorList>
    );
  };

  const extensionPointContent = (icon, href, name, drawerCollapsed) => {
    let content = (
      <>
        <NavigatorLink data-testid={name}>
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
                  filter: currentPath === href ? activeIconFilter : '',
                }}
                onMouseOver={handleExtensionIconMouseEnter}
                onMouseOut={handleExtensionIconMouseLeave}
              />
            </MainListIcon>
          </CustomTooltip>
          <SideBarText drawerCollapsed={drawerCollapsed}>{name}</SideBarText>
        </NavigatorLink>
      </>
    );

    if (href) {
      content = (
        <Link
          href={href}
          onClick={() => dispatch(updateExtensionType({ extensionType: 'navigator' }))}
        >
          <Box>{content}</Box>
        </Link>
      );
    }

    return content;
  };

  const renderChildren = (idname, children, depth) => {
    if (!children?.length) {
      return null;
    }

    const isLifecycleGroup = idname === LIFECYCLE;
    const ListItemComponent = isLifecycleGroup ? NavigatorListItemIII : NavigatorListItemII;

    return (
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
              return null;
            }

            const isActive = currentPath === hrefc;

            return (
              <div key={idc}>
                <ListItemComponent
                  {...(isLifecycleGroup
                    ? {
                        component: 'a',
                        isShow: !showc,
                      }
                    : {})}
                  button
                  data-testid={idc}
                  depth={depth}
                  isDrawerCollapsed={isDrawerCollapsed}
                  isActive={isActive}
                  onClick={() => {
                    if (isLifecycleGroup) {
                      handleAdapterClick(idc, linkc);
                    }

                    if (linkc && hrefc) {
                      router.push(hrefc);
                    }
                  }}
                  disabled={permissionc ? !CAN(permissionc.action, permissionc.subject) : false}
                >
                  {linkContent(iconc, titlec, hrefc, false, isDrawerCollapsed)}
                </ListItemComponent>
                {renderChildren(idname, childrenc, depth + 1)}
              </div>
            );
          },
        )}
      </List>
    );
  };

  const linkContent = (iconc, titlec, hrefc, linkc, drawerCollapsed) => {
    const updatedIcon =
      React.isValidElement(iconc) && typeof iconc.type !== 'string'
        ? React.cloneElement(iconc, {
            fill: currentPath === hrefc ? theme.palette.icon.brand : theme.palette.common.white,
          })
        : iconc;
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

  const Title = (
    <div
      style={
        !providerUiAccessControl?.isNavigatorComponentEnabled([DASHBOARD]) ? cursorNotAllowed : {}
      }
    >
      <>
        <StyledListItem
          component="a"
          onClick={handleTitleClick}
          disableLogo={!providerUiAccessControl?.isNavigatorComponentEnabled([DASHBOARD])}
        >
          {isDrawerCollapsed ? (
            <>
              <MainLogoCollapsed
                src="/static/img/meshery-logo/meshery-logo.png"
                onClick={handleTitleClick}
              />
              <MainLogoTextCollapsed
                src="/static/img/meshery-logo/meshery-logo-text.png"
                onClick={handleTitleClick}
              />
            </>
          ) : (
            <>
              <MainLogo
                src="/static/img/meshery-logo/meshery-logo.png"
                onClick={handleTitleClick}
              />
              <MainLogoText
                src="/static/img/meshery-logo/meshery-logo-text.png"
                onClick={handleTitleClick}
              />
            </>
          )}
        </StyledListItem>
      </>
    </div>
  );

  const Menu = (
    <>
      <HideScrollbar disablePadding>
        {navigatorComponents.map(
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
            const hasChildren = Array.isArray(children) && children.length > 0;
            return (
              <RootDiv key={childId}>
                {/* Row wraps the navigable anchor and the expand/collapse caret as
                    siblings so the caret button is never nested inside the anchor. */}
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <SideBarListItem
                    dense
                    key={childId}
                    link={!!link}
                    isActive={currentPath === href}
                    isShow={!show}
                    onClick={() => {
                      // Leaf items navigate via their link; there is nothing to expand/collapse,
                      // so never add them to openItems (doing so would strand them there and break
                      // the submenu-aware onMouseLeave below).
                      if (!hasChildren) return;
                      // Keep an already-open link submenu open when its row is clicked again.
                      if (link && openItems.includes(childId)) return;
                      toggleItemCollapse(childId);
                    }}
                    onMouseOver={() => (isDrawerCollapsed ? setHoveredId(childId) : null)}
                    onMouseLeave={() =>
                      !submenu || !openItems.includes(childId) ? setHoveredId(null) : null
                    }
                    disabled={permission ? !CAN(permission.action, permission.subject) : false}
                    style={{ flex: 1, minWidth: 0 }}
                    {...(link && href ? { component: Link, href } : {})}
                  >
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
                        (hoveredId === childId || (openItems.includes(childId) && submenu)) ? (
                          <div>
                            <CustomTooltip
                              title={title}
                              placement="right"
                              TransitionComponent={Zoom}
                            >
                              <ListItemIcon style={{ marginLeft: '20%', marginBottom: '0.4rem' }}>
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
                  </SideBarListItem>
                  {hasChildren && (
                    <ExpandMore
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        toggleItemCollapse(childId);
                      }}
                      isExpanded={openItems.includes(childId)}
                      theme={theme}
                      hasChildren={hasChildren}
                    />
                  )}
                </div>
                <Collapse
                  in={openItems.includes(childId)}
                  style={{
                    backgroundColor: theme.palette.navigation.secondary,
                    opacity: '100%',
                  }}
                >
                  {renderChildren(childId, children, 1)}
                </Collapse>
              </RootDiv>
            );
          },
        )}
        {navigatorExtensions.length ? (
          <Box component="div" data-testid="extension-navigation-region">
            <SecondaryDivider />
            {renderNavigatorExtensions(navigatorExtensions, 1)}
          </Box>
        ) : null}
        <SecondaryDivider />
      </HideScrollbar>
    </>
  );

  const [hoveredHelpIcon, setHoveredHelpIcon] = React.useState<string | null>(null);

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
              style={isDrawerCollapsed && !showHelperButton ? { display: 'none' } : {}}
              onMouseEnter={() => setHoveredHelpIcon(id)}
              onMouseLeave={() => setHoveredHelpIcon(null)}
            >
              <Grow
                in={showHelperButton || !isDrawerCollapsed}
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
        <div style={{ textAlign: 'center', width: '100%' }}>{versionDetail.build}</div>
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
        providerUiAccessControl?.isNavigatorComponentEnabled?.([TOGGLER]) ? {} : cursorNotAllowed
      }
    >
      <div
        style={
          providerUiAccessControl?.isNavigatorComponentEnabled?.([TOGGLER]) ? {} : disabledStyle
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
      <SidebarDrawer ref={sidebarNodeRef} isCollapsed={isDrawerCollapsed} variant="permanent">
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
