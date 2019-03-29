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
import SettingsInputComponentIcon from '@material-ui/icons/SettingsInputComponent';
import CollectionsIcon from '@material-ui/icons/Collections';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import TimerIcon from '@material-ui/icons/Timer';
import LinkIcon from '@material-ui/icons/Link';
import SettingsIcon from '@material-ui/icons/Settings';
import Link from "next/link";
import {connect} from "react-redux";
import { bindActionCreators } from 'redux'
import { updatepagepathandtitle } from '../lib/store';
import NoSsr from '@material-ui/core/NoSsr';
import Avatar from '@material-ui/core/Avatar';
import { withRouter } from 'next/router';

const categories = [
  { id: 'Performance', icon: <TimerIcon />, href: "/performance", title: 'Performance Test' },
  { id: 'Configure', icon: <SettingsIcon />, href: "/configure", title: 'Configure Mesh' },
  { id: 'Play', icon: <SettingsInputComponentIcon />, href: "/play", title: 'Play with Mesh' },
  { id: 'Results', icon: <CollectionsIcon />, href: "/results", title: 'View & Compare Results' },
  { id: 'Grafana', icon: <ShowChartIcon />, href: "/grafana", title: 'Grafana Charts' },
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
  }
});

class Navigator extends React.Component {

  state = {
        path: '',
    }

    // static updateTitle = (props) => {
    //     let path = (typeof window !== 'undefined' ? window.location.pathname : '');
    //     categories.map(({title, href}) => {
    //         if (path.lastIndexOf('/') > 0) {
    //             path = path.substring(0, path.lastIndexOf('/'));
    //         }
    //         if (href === path) {
    //             console.log("updating path: "+path+" and title: "+title);
    //             props.updatepagepathandtitle({path, title});
    //             return;
    //         }
    //     });
    //     return path;
    // }
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
                  className={classNames(classes.firebase, classes.item, classes.itemCategory)}>
                  <Avatar className={classes.mainLogo} src={'/static/img/meshery-logo.png'} />
                  Meshery
                </ListItem>
                    {categories.map(({ id: childId, icon, href }) => (
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
                                    <ListItemIcon>{icon}</ListItemIcon>
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
                    ))}
                      <Divider className={classes.divider} />
                      <ListItem
                            component="a"
                            href="https://layer5.io/meshery"
                            key={'about'}
                            className={classNames(
                            classes.item,
                            classes.itemActionable,
                            classes.community,
                            )}
                        >
                          <div className={classNames(classes.link)} >
                              <ListItemIcon><LinkIcon /></ListItemIcon>
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