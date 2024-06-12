import React, { useState, useRef } from 'react';
import { List, ListItem } from '@material-ui/core';
import { Avatar } from '@layer5/sistent';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import Grow from '@material-ui/core/Grow';
import IconButton from '@material-ui/core/IconButton';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';
import NoSsr from '@material-ui/core/NoSsr';
import Paper from '@material-ui/core/Paper';
import Popper from '@material-ui/core/Popper';
import classNames from 'classnames';
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
        <List disablePadding>
          {children.map(({ id, href, title, show: showc }) => {
            if (typeof showc !== 'undefined' && !showc) {
              return '';
            }
            return (
              <React.Fragment key={id}>
                <ListItem button key={id}>
                  {extensionPointContent(href, title)}
                </ListItem>
              </React.Fragment>
            );
          })}
        </List>
      );
    }
  }

  function extensionPointContent(href, name) {
    const { classes } = props;

    const content = (
      <LinkDiv>
        <ListItemText classes={{ primary: classes.itemPrimary }}>{name}</ListItemText>
      </LinkDiv>
    );
    if (href) {
      return (
        <Link href={href}>
          <span
            className={classNames(classes.link)}
            onClick={() => props.updateExtensionType(name)}
          >
            {content}
          </span>
        </Link>
      );
    }

    return content;
  }

  const { color, iconButtonClassName, avatarClassName, classes } = props;

  const open = Boolean(anchorEl);
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
                      <MenuItem
                        onClick={handlePreference}
                        // disabled={
                        //   !CAN(
                        //     keys.VIEW_MESHERY_USER_PREFERENCES.action,
                        //     keys.VIEW_MESHERY_USER_PREFERENCES.subject,
                        //   )
                        // }
                      >
                        Preferences
                      </MenuItem>
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
