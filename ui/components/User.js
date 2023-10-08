import { List, ListItem } from '@material-ui/core';
import Avatar from '@material-ui/core/Avatar';
import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import Grow from '@material-ui/core/Grow';
import IconButton from '@material-ui/core/IconButton';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';
import NoSsr from '@material-ui/core/NoSsr';
import Paper from '@material-ui/core/Paper';
import Popper from '@material-ui/core/Popper';
import { withStyles } from '@material-ui/core/styles';
import classNames from 'classnames';
import Link from "next/link";
import { useRouter } from 'next/router';
import React, { useState, useEffect , useRef } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import dataFetch from '../lib/data-fetch';
import { updateUser } from '../lib/store';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';

const styles = () => ({
  link: {
    display: 'inline-flex',
    width: '100%',
    height: '30px',
    alignItems: 'self-end',
  },
});

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
  const [user, setUser] = useState(null);
  const [open, setOpen] = useState(false);
  const [account, setAccount] = useState([]);
  const [capabilitiesLoaded, setCapabilitiesLoaded] = useState(false);
  const anchorEl = useRef(null);
  const dispatch = useDispatch();
  const router = useRouter();

  const capabilitiesRegistry = useSelector(state => state.get("capabilitiesRegistry"));

  const handleToggle = () => {
    setOpen(!open);
  };

  const handleClose = (event) => {
    if (anchorEl.current && anchorEl.current.contains(event.target)) {
      return;
    }
    setOpen(false);
  }

  const handleLogout = () => {
    window.location = '/user/logout'
  };

  const handlePreference = () => {
    router.push('/user/preferences');
  }

  const handleGetToken = () => {
    dataFetch('/api/token', { credentials : 'same-origin' }, (data) => {
      exportToJsonFile(data, "auth.json");
    }, (error) => ({ error, }));
  }

  useEffect(() => {
    dataFetch('/api/user', {
      credentials : 'same-origin'
    }, (userData) => {
      setUser(userData);
      dispatch(updateUser({ user : userData }));
    }, (error) => ({
      error,
    }));
  }, [dispatch])

  useEffect(() => {
    const { capabilitiesRegistry } = props;
    if (!capabilitiesLoaded && capabilitiesRegistry) {
      setCapabilitiesLoaded(true); // to prevent re-compute
      setAccount(ExtensionPointSchemaValidator("account")(capabilitiesRegistry?.extensions?.account));
    }
  }, [capabilitiesRegistry, capabilitiesLoaded])

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
                <ListItem
                  button
                  key={id}
                >
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
      <div className={classNames(classes.link)}>
        <ListItemText classes={{ primary: classes.itemPrimary }}>{name}</ListItemText>
      </div>
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
  let avatar_url;
  if (user && user !== null) {
    avatar_url = user.avatar_url;
  }

  return (
    <div>
      <NoSsr>
        <div data-test="profile-button">
          <IconButton
            color={color}
            className={iconButtonClassName}
            ref = {anchorEl}
            aria-owns={open
              ? 'menu-list-grow'
              : undefined}
            aria-haspopup="true"
            onClick={handleToggle}
          >
            <Avatar className={avatarClassName} src={avatar_url} imgProps={{ referrerPolicy : "no-referrer" }} />
          </IconButton>
        </div>
        <Popper open={open} anchorEl={anchorEl} transition style={{ zIndex : 10000 }} placement="top-end">
          {({ TransitionProps, placement }) => (
            <Grow
              {...TransitionProps}
              id="menu-list-grow"
              style={{
                transformOrigin : placement === 'bottom'
                  ? 'left top'
                  : 'left bottom'
              }}
            >
              <Paper className={classes.popover}>
                <ClickAwayListener onClickAway={handleClose}>
                  <MenuList>
                    {account && account.length ? (
                      <>
                        {renderAccountExtension(account)}
                      </>
                    ) : null}
                    <MenuItem onClick={handleGetToken}>Get Token</MenuItem>
                    <MenuItem onClick={handlePreference}>Preferences</MenuItem>
                    <MenuItem onClick={handleLogout}>Logout</MenuItem>
                  </MenuList>
                </ClickAwayListener>
              </Paper>
            </Grow>
          )}
        </Popper>
      </NoSsr>
    </div>
  );

}

export default withStyles(styles)(User);
