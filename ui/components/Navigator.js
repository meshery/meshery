import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { withStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import CollectionsIcon from '@material-ui/icons/Collections';
import LaptopIcon from '@material-ui/icons/Laptop';
import TimerIcon from '@material-ui/icons/Timer';
import ExitToAppIcon from '@material-ui/icons/ExitToApp';
import Link from "next/link";
import {connect} from "react-redux";
import { bindActionCreators } from 'redux'
import { updatepagepathandtitle } from '../lib/store';
import NoSsr from '@material-ui/core/NoSsr';
import Avatar from '@material-ui/core/Avatar';
import { withRouter } from 'next/router';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTerminal, faTachometerAlt, faSignal, faExternalLinkAlt, faPollH } from '@fortawesome/free-solid-svg-icons';

const categories = [
  { id: 'Dashboard', href: "/", title: 'Dashboard', show: false},
  { id: 'Performance', icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />, href: "/performance", title: 'Performance Test', show: true },
  { id: 'Settings', href: "/settings", title: 'Settings', show: false}, // title is used for comparison in the Header.js file as well
  { id: 'Playground', icon:  <FontAwesomeIcon icon={faTerminal} transform="shrink-4" fixedWidth />, href: "/playground", title: 'Manage Service Meshes', show: true },
  { id: 'Results', icon: <FontAwesomeIcon icon={faPollH} fixedWidth />, href: "/results", title: 'View & Compare Results', show: true },
]

const styles = theme => ({
  categoryHeader: {
    paddingTop: 16,
    paddingBottom: 16,
  },
  categoryHeaderPrimary: {
    color: theme.palette.common.white,
  },
  item: {
    paddingTop: 4,
    paddingBottom: 4,
    color: 'rgba(255, 255, 255, 0.7)',
  },
  itemCategory: {
    backgroundColor: '#263238',
    boxShadow: '0 -1px 0 #404854 inset',
    paddingTop: 16,
    paddingBottom: 16,
  },
  firebase: {
    fontSize: 24,
    fontFamily: theme.typography.fontFamily,
    color: theme.palette.common.white,
  },
  link: {
    display: 'inline-flex',
    width: '100%',
  },
  itemActionable: {
    '&:hover': {
      backgroundColor: 'rgba(255, 255, 255, 0.08)',
    },
  },
  itemActiveItem: {
    color: '#4fc3f7',
  },
  itemPrimary: {
    color: 'inherit',
    fontSize: theme.typography.fontSize,
    '&$textDense': {
      fontSize: theme.typography.fontSize,
    },
  },
  textDense: {},
  divider: {
    marginTop: theme.spacing(2),
  },
  mainLogo: {
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(1),
    width: 30,
    height: 30,
    borderRadius: 'unset',
  },
  community: {
    marginTop: theme.spacing(2),
  },
  settingsIcon: {
    marginLeft: theme.spacing(2),
  },
  cursorPointer: {
    cursor: 'pointer',
  },
  listIcon: {
    minWidth: theme.spacing(3.5),
    paddingTop: theme.spacing(0.5),
  }
});

class Navigator extends React.Component {

  state = {
        path: '',
    }

    static getDerivedStateFromProps(props, state) {
      let path = (typeof window !== 'undefined' ? window.location.pathname : '');
      categories.map(({title, href}) => {
          if (path.lastIndexOf('/') > 0) {
              path = path.substring(0, path.lastIndexOf('/'));
          }
          if (href === path) {
              // console.log("updating path: "+path+" and title: "+title);
              props.updatepagepathandtitle({path, title});
              return;
          }
      });
      return {path};
    }

    handleTitleClick = (event) => {
      this.props.router.push('/');
    }

    render() {
        const { classes, updatepagepathandtitle, ...other } = this.props;
        const { path } = this.state;
        // const path = this.updateTitle();
        // console.log("current page:" + path);
        return (
            <NoSsr>
            <Drawer variant="permanent" {...other}>
            <List disablePadding>
                <ListItem 
                  component="a"
                  onClick={this.handleTitleClick}
                  className={classNames(classes.firebase, classes.item, classes.itemCategory, classes.cursorPointer)}>
                  <Avatar className={classes.mainLogo} src={'/static/img/meshery-logo.png'} onClick={this.handleTitleClick} />
                  Meshery
                </ListItem>
                    {categories.map(({ id: childId, icon, href, show }) => {
                      if (typeof show !== 'undefined' && !show){
                        return '';
                      }
                      return (
                        <ListItem
                            button
                            dense
                            key={childId}
                            className={classNames(
                            classes.item,
                            classes.itemActionable,
                            path === href && classes.itemActiveItem,
                            )}
                        >
                            <Link href={href} prefetch>
                                <div className={classNames(classes.link)} >
                                    <ListItemIcon className={classes.listIcon}>{icon}</ListItemIcon>
                                    <ListItemText
                                    classes={{
                                        primary: classes.itemPrimary,
                                        textDense: classes.textDense,
                                    }}
                                    >
                                    {childId}
                                    </ListItemText>
                                </div>
                            </Link>
                        </ListItem>
                        );
                      })}
                      <Divider className={classes.divider} />
                      <ListItem
                            component="a"
                            href="https://layer5.io/meshery"
                            target="_blank"
                            key={'about'}
                            className={classNames(
                            classes.item,
                            classes.itemActionable,
                            classes.community,
                            )}
                        >
                          <div className={classNames(classes.link)} >
                              <ListItemIcon className={classes.listIcon}><FontAwesomeIcon icon={faExternalLinkAlt} transform="shrink-2" fixedWidth /></ListItemIcon>
                              <ListItemText
                              classes={{
                                  primary: classes.itemPrimary,
                                  textDense: classes.textDense,
                              }}
                              >
                              {'Community'}
                              </ListItemText>
                          </div>
                        </ListItem>
                    <Divider className={classes.divider} />
            </List>
            </Drawer>
            </NoSsr>
        );
    }
}

Navigator.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    updatepagepathandtitle: bindActionCreators(updatepagepathandtitle, dispatch)
  }
}

export default withStyles(styles)(connect(
    null,
    mapDispatchToProps
  )(withRouter(Navigator)));