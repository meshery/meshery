import React, { useState, useRef, useEffect } from 'react';
import MenuIcon from '@mui/icons-material/Menu';
import { useRouter } from 'next/router';
import { useGetLoggedInUserQuery, useLazyGetTokenQuery } from '@/rtk-query/user';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { Popover } from '@sistent/sistent';
// @ts-expect-error - NavigationNavbar exists at runtime but types may not be exported
import { NavigationNavbar } from '@sistent/sistent';
import { IconButtonMenu } from './Header.styles';
import { useDispatch, useSelector } from 'react-redux';
import { updateExtensionType, updateUser } from '@/store/slices/mesheryUi';
import { useRegistryModal } from '@/utils/hooks/useRegistryModal';

function exportToJsonFile(jsonData, filename) {
  let dataStr = JSON.stringify(jsonData);
  let dataUri = 'data:application/json;charset=utf-8,' + encodeURIComponent(dataStr);
  let linkElement = document.createElement('a');
  linkElement.setAttribute('href', dataUri);
  linkElement.setAttribute('download', filename);
  linkElement.click();
  linkElement.remove();
}

/**
 * Extension Point: Avatar behavior for User Modes
 * Insert custom logic here to handle Single User mode, Anonymous User mode, Multi User mode behavior.
 */
const HeaderMenu = () => {
  const dispatch = useDispatch();
  const { capabilitiesRegistry } = useSelector((state: any) => state.ui);
  const [userLoaded, setUserLoaded] = useState(false);
  const [account, setAccount] = useState<
    Array<{ id: string; title: string; href?: string; show?: boolean }>
  >([]);
  const capabilitiesLoadedRef = useRef(false);
  const { notify } = useNotification();
  const router = useRouter();
  const [anchorEl, setAnchorEl] = useState(null);
  const registryModal = useRegistryModal();

  const {
    data: userData,
    isSuccess: isGetUserSuccess,
    isError: isGetUserError,
    error: getUserError,
  } = useGetLoggedInUserQuery(undefined);

  const [triggerGetToken, { isError: isTokenError, error: tokenError }] = useLazyGetTokenQuery();

  const handleLogout = () => {
    window.location.href = '/user/logout';
    handleClose();
  };

  const handlePreference = () => {
    router.push('/user/preferences');
    handleClose();
  };

  const handleSettings = () => {
    router.push('/settings');
    handleClose();
  };
  const handleRegistry = () => {
    registryModal?.openModal();
    handleClose();
  };
  const handleConnections = () => {
    router.push('/management/connections');
    handleClose();
  };

  const handleGetToken = () => {
    triggerGetToken(undefined)
      .unwrap()
      .then((data) => {
        exportToJsonFile(data, 'auth.json');
        handleClose();
      });
  };

  useEffect(() => {
    if (!userLoaded && isGetUserSuccess) {
      dispatch(updateUser({ user: userData }));
      setUserLoaded(true);
    } else if (isGetUserError) {
      notify({
        message: 'Error fetching user',
        event_type: EVENT_TYPES.ERROR,
        details: (getUserError as any)?.data,
      });
    }
  }, [userData, isGetUserSuccess, isGetUserError, getUserError, notify, dispatch]);

  useEffect(() => {
    if (isTokenError) {
      notify({
        message: 'Error fetching token',
        event_type: EVENT_TYPES.ERROR,
        details: (tokenError as any)?.data,
      });
    }
  }, [isTokenError, tokenError, notify]);

  useEffect(() => {
    if (!capabilitiesLoadedRef.current && capabilitiesRegistry) {
      capabilitiesLoadedRef.current = true;
      setAccount(
        ExtensionPointSchemaValidator('account')(capabilitiesRegistry?.extensions?.account),
      );
    }
  }, [capabilitiesRegistry]);

  const getAccountNavigationItems = () => {
    const accountItems = account.map((item) => ({
      id: item.id,
      title: item.title,
      onClick: () => {
        if (item.href) {
          dispatch(updateExtensionType({ extensionType: item.title }));
          router.push(item.href);
          handleClose();
        }
      },
      permission: typeof item.show === 'undefined' ? true : item.show,
    }));

    const defaultItems: Array<{
      id: string;
      title: string;
      onClick: () => void;
      permission?: boolean;
      showOnWeb?: boolean;
    }> = [];

    // Only add Get Token if there are no account items
    if (!account.length) {
      defaultItems.push({
        id: 'get-token',
        title: 'Get Token',
        onClick: handleGetToken,
        permission: CAN(keys.DOWNLOAD_TOKEN.action, keys.DOWNLOAD_TOKEN.subject),
      });
    }

    // Always add these items
    defaultItems.push(
      {
        id: 'Connections',
        title: 'Connections',
        onClick: handleConnections,
        showOnWeb: false,
      },
      {
        id: 'settings',
        title: 'Settings',
        onClick: handleSettings,
      },
      {
        id: 'registry',
        title: 'Registry',
        onClick: handleRegistry,
        permission: true,
      },
      {
        id: 'preferences',
        title: 'Preferences',
        onClick: handlePreference,
        permission: true,
      },
      {
        id: 'logout',
        title: 'Logout',
        onClick: handleLogout,
        permission: true,
      },
    );

    return [...accountItems, ...defaultItems];
  };

  if (userData?.status === 'anonymous') {
    return null;
  }

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'menu-popover' : undefined;

  return (
    <>
      <IconButtonMenu aria-describedby={id} onClick={handleClick}>
        <MenuIcon
          sx={{
            height: 28,
            width: 28,
          }}
        />
      </IconButtonMenu>

      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
        style={{ marginTop: '1rem' }}
      >
        <NavigationNavbar
          navigationItems={getAccountNavigationItems()}
          ListItemTextProps={{
            primaryTypographyProps: {
              sx: {
                fontSize: '1rem',
              },
            },
          }}
        />
      </Popover>
    </>
  );
};

const MenuProvider = (props) => <HeaderMenu {...props} />;

export default MenuProvider;
