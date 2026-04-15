// @ts-nocheck
import React, { useEffect, useState } from 'react';
import { RemoveIcon } from '@sistent/sistent';
import {
  ExternalLinkIcon as IconExternalLink,
  SlackIcon,
  FileIcon,
  GithubIcon,
  DiscussForumIcon,
  useTheme,
} from '@sistent/sistent';
import ExtensionPointSchemaValidator from '../../../utils/ExtensionPointSchemaValidator';
import { CapabilitiesRegistry } from '../../../utils/disabledComponents';
import { LIFECYCLE, SERVICE_MESH, CONFIGURATION, CATALOG } from '../../../constants/navigator';
import { getProviderCapabilities, getSystemVersion } from '@/rtk-query/user';
import { useDispatch, useSelector } from 'react-redux';
import {
  toggleDrawer,
  updateBetaBadge,
  updateCapabilities,
  updateTitle,
} from '@/store/slices/mesheryUi';
import { useRouter } from 'next/router';
import { setAdapter } from '@/store/slices/adapter';
import { getNavigatorComponents, externalLinkIconStyle } from './navigationSchema';

export const useNavigatorState = () => {
  const { meshAdapters } = useSelector((state) => state.adapter);
  const { meshAdaptersts } = useSelector((state) => state.adapter);
  const dispatch = useDispatch();
  const { capabilitiesRegistry } = useSelector((state) => state.ui);
  const { catalogVisibility } = useSelector((state) => state.ui);
  const { isDrawerCollapsed } = useSelector((state) => state.ui);
  const theme = useTheme();
  const router = useRouter();
  const [state, setState] = useState({
    path: '',
    meshAdapters,
    mts: Date.now(),
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
  }, [meshAdapters, meshAdaptersts, router.asPath, dispatch]);

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
      const navigatorComponents = getNavigatorComponents(capabilitiesRegistryObj, theme);
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

  const updateNavigatorComponentsMenus = () => {
    const updated = state.navigatorComponents.map((cat) => {
      const catCopy = { ...cat };

      if (catCopy.id === LIFECYCLE) {
        catCopy.children = catCopy.children.map((catc) => {
          if (catc.id == SERVICE_MESH) return catc;
          return {
            ...catc,
            icon: pickIcon(catc.id, catc.href),
            children: fetchChildren(catc.id),
          };
        });
      }

      if (catCopy.id === 'Configuration') {
        let show = false;
        catCopy.children = catCopy.children?.map((ch) => {
          if (ch.id === 'Designs') {
            const idx = capabilitiesRegistry?.capabilities?.findIndex(
              (cap) => cap.feature === 'persist-meshery-patterns',
            );
            if (idx != -1) {
              show = true;
              return { ...ch, show: true };
            }
          }
          return ch;
        });
        catCopy.show = show;
      }

      if (catCopy.id === CONFIGURATION) {
        catCopy.children = catCopy.children?.map((ch) => {
          if (ch.id === CATALOG) return { ...ch, show: catalogVisibility };
          return ch;
        });
      }

      return catCopy;
    });
    updateState({ navigatorComponents: updated });
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

  return {
    state,
    updateState,
    theme,
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
    fetchChildren,
    pickIcon,
  };
};
