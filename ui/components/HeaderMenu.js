import React, { useState, useRef } from 'react';
import MenuIcon from '@material-ui/icons/Menu';
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
const HeaderMenu = (props) => {
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

  const { color, iconButtonClassName, classes } = props;

  const open = Boolean(anchorEl);

  if (userData?.status == 'anonymous') {
    return null;
  }

  return (
    <div>
      <NoSsr>
        <div>
          <IconButton
            color={color}
            className={iconButtonClassName}
            ref={anchorEl}
            onClick={handleOpen}
            aria-owns={open ? 'menu-list-grow' : undefined}
            aria-haspopup="true"
          >
            <MenuIcon />
          </IconButton>
        </div>
        <div>
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
                  marginTop: '1rem',
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

const MenuProvider = (props) => {
  return (
    <Provider store={store}>
      <HeaderMenu {...props} />
    </Provider>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateUser: bindActionCreators(updateUser, dispatch),
});

const mapStateToProps = (state) => ({
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
});

export default connect(mapStateToProps, mapDispatchToProps)(MenuProvider);
