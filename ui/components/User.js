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
import Link from 'next/link';
import { withRouter } from 'next/router';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import dataFetch from '../lib/data-fetch';
import { updateUser } from '../lib/store';
import ExtensionPointSchemaValidator from '../utils/ExtensionPointSchemaValidator';
import { withNotify } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';

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

class User extends React.Component {
  state = {
    user: null,
    open: false,
    account: ExtensionPointSchemaValidator('account')(),
    providerType: '',
    capabilitiesLoaded: false,
  };

  handleToggle = () => {
    this.setState((state) => ({ open: !state.open }));
  };

  handleClose = (event) => {
    if (this.anchorEl.contains(event.target)) {
      return;
    }
    this.setState({ open: false });
  };

  handleLogout = () => {
    window.location = '/user/logout';
  };

  handleError = (error) => {
    const notify = this.props.notify;
    notify({
      message: `Error performing logout: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  handlePreference = () => {
    this.props.router.push('/user/preferences');
  };

  handleGetToken = () => {
    dataFetch(
      '/api/token',
      { credentials: 'same-origin' },
      (data) => {
        exportToJsonFile(data, 'auth.json');
      },
      (error) => ({ error }),
    );
  };

  componentDidMount() {
    dataFetch(
      '/api/user',
      {
        credentials: 'same-origin',
      },
      (user) => {
        this.setState({ user });
        this.props.updateUser({ user });
      },
      (error) => ({
        error,
      }),
    );
  }

  componentDidUpdate() {
    const { capabilitiesRegistry } = this.props;
    if (!this.state.capabilitiesLoaded && capabilitiesRegistry) {
      this.setState({
        capabilitiesLoaded: true, // to prevent re-compute
        account: ExtensionPointSchemaValidator('account')(
          capabilitiesRegistry?.extensions?.account,
        ),
        providerType: capabilitiesRegistry?.provider_type,
      });
    }
  }

  /**
   * @param {import("../utils/ExtensionPointSchemaValidator").AccountSchema[]} children
   */
  renderAccountExtension(children) {
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
                  {this.extensionPointContent(href, title)}
                </ListItem>
              </React.Fragment>
            );
          })}
        </List>
      );
    }
  }

  extensionPointContent(href, name) {
    const { classes } = this.props;

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
            onClick={() => this.props.updateExtensionType(name)}
          >
            {content}
          </span>
        </Link>
      );
    }

    return content;
  }

  render() {
    const { color, iconButtonClassName, avatarClassName, classes } = this.props;
    let avatar_url;
    if (this.state.user && this.state.user !== null) {
      avatar_url = this.state.user.avatar_url;
    }
    const { open } = this.state;
    return (
      <div>
        <NoSsr>
          <div data-test="profile-button">
            <IconButton
              color={color}
              className={iconButtonClassName}
              buttonRef={(node) => {
                this.anchorEl = node;
              }}
              aria-owns={open ? 'menu-list-grow' : undefined}
              aria-haspopup="true"
              onClick={this.handleToggle}
            >
              <Avatar
                className={avatarClassName}
                src={avatar_url}
                imgProps={{ referrerPolicy: 'no-referrer' }}
              />
            </IconButton>
          </div>
          <Popper
            open={open}
            anchorEl={this.anchorEl}
            transition
            style={{ zIndex: 10000 }}
            placement="top-end"
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
                  <ClickAwayListener onClickAway={this.handleClose}>
                    <MenuList>
                      {this.state.account && this.state.account.length ? (
                        <>{this.renderAccountExtension(this.state.account)}</>
                      ) : null}
                      <MenuItem onClick={this.handleGetToken}>Get Token</MenuItem>
                      <MenuItem onClick={this.handlePreference}>Preferences</MenuItem>
                      <MenuItem onClick={this.handleLogout}>Logout</MenuItem>
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
}

const mapDispatchToProps = (dispatch) => ({ updateUser: bindActionCreators(updateUser, dispatch) });
const mapStateToProps = (state) => ({
  capabilitiesRegistry: state.get('capabilitiesRegistry'),
});

export default withStyles(styles)(
  connect(mapStateToProps, mapDispatchToProps)(withNotify(withRouter(User))),
);
