import React, { useState, useRef } from 'react';
import { Avatar } from '@layer5/sistent';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import Grow from '@material-ui/core/Grow';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';
import NoSsr from '@material-ui/core/NoSsr';
import Paper from '@material-ui/core/Paper';
import Popper from '@material-ui/core/Popper';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { Provider, connect } from 'react-redux';
import { store } from '../store';
import { bindActionCreators } from 'redux';
import { useGetLoggedInUserQuery, useLazyGetTokenQuery } from '@/rtk-query/user';
import { updateUser } from '../lib/store';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import { styled } from '@mui/material/styles';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { Button } from '@material-ui/core';

const LinkDiv = styled('div')(() => ({
  display: 'inline-flex',
  width: '100%',
  height: '30px',
  alignItems: 'self-end',
}));

function exportToJsonFile(jsonData, filename) {
  let dataStr = JSON.stringify(jsonData);
  let dataUri = 'data:application/json;charset=utf-8,' + encodeURIComponent(dataStr);

  let exportFileDefaultName = filename;

  let linkElement = document.createElement('a');
  linkElement.setAttribute('href', dataUri);
  linkElement.setAttribute('download', exportFileDefaultName);
  linkElement.click();
  linkElement.remove();
}

/**
 * Extension Point: Avatar behavior for User Modes
 * Insert custom logic here to handle Single User mode, Anonymous User mode, Multi User mode behavior.
 */
const User = (props) => {
  const [userLoaded, setUserLoaded] = useState(false);
  const [account, setAccount] = useState([]);
  const capabilitiesLoadedRef = useRef(false);
  const { notify } = useNotification();
  const [anchorEl, setAnchorEl] = useState(null);
  const router = useRouter();

  const {
    data: userData,
    isSuccess: isGetUserSuccess,
    isError: isGetUserError,
    error: getUserError,
  } = useGetLoggedInUserQuery();
  const [triggerGetToken, { isError: isTokenError, error: tokenError }] = useLazyGetTokenQuery();

  const { capabilitiesRegistry } = props;

  const getProfileUrl = () => {
    return (account || [])?.find((item) => item.title === 'Profile')?.href;
  };

  const goToProfile = () => {
    const profileUrl = getProfileUrl();
    if (profileUrl) {
      window.location = profileUrl;
      return;
    }
  };

  const handleOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    window.location = '/user/logout';
  };

  const handlePreference = () => {
    router.push('/user/preferences');
  };

  const handleGetToken = () => {
    triggerGetToken()
      .unwrap()
      .then((data) => {
        exportToJsonFile(data, 'auth.json');
      });
  };

  if (!userLoaded && isGetUserSuccess) {
    props.updateUser({ user: userData });
    setUserLoaded(true);
  } else if (isGetUserError) {
    notify({
      message: 'Error fetching user',
      event_type: EVENT_TYPES.ERROR,
      details: getUserError?.data,
    });
  }

  if (isTokenError) {
    notify({
      message: 'Error fetching token',
      event_type: EVENT_TYPES.ERROR,
      details: tokenError?.data,
    });
  }

  if (!capabilitiesLoadedRef.current && capabilitiesRegistry) {
    capabilitiesLoadedRef.current = true;
    setAccount(ExtensionPointSchemaValidator('account')(capabilitiesRegistry?.extensions?.account));
  }

  /**
   * @param {import("../utils/ExtensionPointSchemaValidator").AccountSchema[]} children
   */
  function renderAccountExtension(children) {
    if (children && children.length > 0) {
      return (
        <>
          {children.map(({ id, href, title, show: showc }) => {
            if (typeof showc !== 'undefined' && !showc) {
              return '';
            }
            return (
              <React.Fragment key={id}>
                <MenuItem button key={id}>
                  {extensionPointContent(href, title)}
                </MenuItem>
              </React.Fragment>
            );
          })}
        </>
      );
    }
  }

  function extensionPointContent(href, name) {
    const content = <LinkDiv>{name}</LinkDiv>;
    if (href) {
      return (
        <Link onClick={() => props.updateExtensionType(name)} href={href}>
          {content}
        </Link>
      );
    }

    return content;
  }

  const { color, iconButtonClassName, avatarClassName, classes } = props;

  const open = Boolean(anchorEl);

  if (userData?.status == 'anonymous') {
    const url = `${capabilitiesRegistry?.provider_url}/login?anonymousUserID=${userData?.id}&redirect=${window.location.pathname}`;

    return (
      <Link href={url}>
        <Button variant="contained" onClick={handleLogout} color="primary">
          Sign In
        </Button>
      </Link>
    );
  }

  return (
    <div>
      <NoSsr>
        <div data-test="profile-button" onMouseOver={handleOpen} onMouseLeave={handleClose}>
          <IconButton
            color={color}
            className={iconButtonClassName}
            ref={anchorEl}
            aria-owns={open ? 'menu-list-grow' : undefined}
            aria-haspopup="true"
            onClick={goToProfile}
          >
            <Avatar
              className={avatarClassName}
              src={isGetUserSuccess ? userData?.avatar_url : null}
              imgProps={{ referrerPolicy: 'no-referrer' }}
            />
          </IconButton>
        </div>
        <div onMouseOver={handleOpen} onMouseLeave={handleClose}>
          <Popper
            open={open}
            anchorEl={anchorEl}
            transition
            style={{ zIndex: 10000 }}
            placement="top-end"
            onClose={handleClose}
          >
            {({ TransitionProps, placement }) => (
              <Grow
                {...TransitionProps}
                id="menu-list-grow"
                style={{
                  transformOrigin: placement === 'bottom' ? 'left top' : 'left bottom',
                }}
              >
                <Paper className={classes.popover}>
                  <ClickAwayListener onClickAway={handleClose}>
                    <MenuList>
                      {account && account.length ? renderAccountExtension(account) : null}
                      {!account?.length && (
                        <MenuItem
                          disabled={!CAN(keys.DOWNLOAD_TOKEN.action, keys.DOWNLOAD_TOKEN.subject)}
                          onClick={handleGetToken}
                        >
                          Get Token
                        </MenuItem>
                      )}
                      <MenuItem onClick={handlePreference}>Preferences</MenuItem>
                      <MenuItem onClick={handleLogout}>Logout</MenuItem>
                    </MenuList>
                  </ClickAwayListener>
                </Paper>
              </Grow>
            )}
          </Popper>
        </div>
      </NoSsr>
    </div>
  );
};

const UserProvider = (props) => {
  return (
    <Provider store={store}>
      <User {...props} />
    </Provider>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateUser: bindActionCreators(updateUser, dispatch),
});

const mapStateToProps = (state) => ({
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
});

export default connect(mapStateToProps, mapDispatchToProps)(UserProvider);
