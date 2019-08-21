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
  },
  nested1: {
    paddingLeft: theme.spacing(2),
  },
  nested2: {
    paddingLeft: theme.spacing(4),
  },
});

const categories = [
  { 
    id: 'Dashboard', 
    href: "/", 
    title: 'Dashboard', 
    show: false,
    link: true,
  },
  { 
    id: 'Performance', 
    icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />, 
    href: "/performance", 
    title: 'Performance Test', 
    show: true,
    link: true,
    children: [
      { 
        id: 'Results', 
        icon: <FontAwesomeIcon icon={faPollH} fixedWidth />, 
        href: "/results", 
        title: 'View & Compare Results', 
        show: true,
        link: true,
      },
    ]
  },
  { 
    id: 'Settings', 
    href: "/settings", 
    title: 'Settings', 
    show: false,
    link: true,
  }, // title is used for comparison in the Header.js file as well
  { 
    id: 'Management', 
    icon:  <FontAwesomeIcon icon={faTerminal} transform="shrink-4" fixedWidth />, 
    href: "/management", 
    title: 'Management', 
    show: true,
    link: true,
    children: [
      {
        id: 'Istio', 
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />, 
        href: "/management/istio", 
        title: 'Istio',
        link: false, 
        show: true,
      },
      {
        id: 'Linkerd', 
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />, 
        href: "/management/linkerd", 
        title: 'Linkerd',
        link: false, 
        show: true,
      },
      {
        id: 'Consul', 
        // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />, 
        href: "/management/consul", 
        title: 'Consul',
        link: false, 
        show: true,
      },
    ],
  },
]

class Navigator extends React.Component {
    constructor(props){
      super(props);
      const {meshAdapters, meshAdaptersts} = props;
      this.state = {
        path: '',
        meshAdapters,
        mts: new Date(),
      };
      
    }

    updateCategoriesMenus() {
      const self = this;
      categories.forEach((cat, ind) => {
        if(cat.id === 'Management'){
          cat.children.forEach((catc, ind1) => {
            const cr = self.fetchChildren(catc.id);
            categories[ind].children[ind1]['children'] = cr;
          });
        }
      });
    }

    static getDerivedStateFromProps(props, state) {
      const { meshAdapters, meshAdaptersts } = props;
      let path = (typeof window !== 'undefined' ? window.location.pathname + window.location.search : '');

      const st = {};
      if(meshAdaptersts > state.mts) {
        st.meshAdapters = meshAdapters;
        st.mts = meshAdaptersts;
      }

      const fetchNestedPathAndTitle = (path, title, href, children) => {
        if (href === path) {
            // console.log("updating path: "+path+" and title: "+title);
            props.updatepagepathandtitle({path, title});
            return;
        }
        if(children && children.length > 0){
          children.forEach(({title, href, children}) => {
              fetchNestedPathAndTitle(path, title, href, children);
          });
        }
      }
      if (path.lastIndexOf('/') > 0) {
          path = path.substring(0, path.lastIndexOf('/'));
      }
      categories.forEach(({title, href, children}) => {    
          fetchNestedPathAndTitle(path, title, href, children); 
      });
      st.path = path;
      return st;
    }

    fetchChildren(category) {
      const { meshAdapters } = this.state;
      const children = [];
      category = category.toLowerCase();
      meshAdapters.forEach(adapter => {
        const aName = adapter.name.toLowerCase();

        if (category !== aName) {
          return;
        }
        children.push(
          {
            id: adapter.adapter_location, 
            // icon: <FontAwesomeIcon icon={faTachometerAlt} transform="shrink-2" fixedWidth />, 
            href: `/management?adapter=${adapter.adapter_location}`, 
            title: `Management - ${adapter.adapter_location}`,
            link: true, 
            show: true,
          }
        );

        // let image = "/static/img/meshery-logo.png";
        // let logoIcon = (<img src={image} className={classes.icon} />);
        // switch (aName){
        //   case 'istio':
        //     image = "/static/img/istio.svg";
        //     logoIcon = (<img src={image} className={classes.istioIcon} />);
        //     break;
        //   case 'linkerd':
        //     image = "/static/img/linkerd.svg";
        //     logoIcon = (<img src={image} className={classes.icon} />);
        //     break;
        //   case 'consul':
        //     image = "/static/img/consul.svg";
        //     logoIcon = (<img src={image} className={classes.icon} />);
        //     break;
        //   case 'nsm':
        //     image = "/static/img/nsm.svg";
        //     logoIcon = (<img src={image} className={classes.icon} />);
        //     break;
        //   // default:
        // }
      });
      return children;
    }

    handleTitleClick = (event) => {
      this.props.router.push('/');
    }

    renderChildren(children, depth) {
      const { classes } = this.props;
      const { path } = this.state;
      
      if (children && children.length > 0){
      return (
        <List disablePadding>
          {children.map(({id: idc, icon: iconc, href: hrefc, show: showc, link: linkc, children: childrenc }) => {
            if (typeof showc !== 'undefined' && !showc){
              return '';
            }
            return (
              <React.Fragment>
              <ListItem button
                className={classNames(
                  (depth === 1?classes.nested1:classes.nested2),
                  classes.item,
                  classes.itemActionable,
                  path === hrefc && classes.itemActiveItem,
                  )}>
                {this.linkContent(iconc, idc, hrefc, linkc)}
              </ListItem>
              {this.renderChildren(childrenc, depth + 1)}
              </React.Fragment>
              );
              })}
            </List>
        );
      }
      return '';
    }

    linkContent(iconc, idc, hrefc, linkc){
      const { classes } = this.props;
      let linkContent = (
        <div className={classNames(classes.link)} >
          <ListItemIcon className={classes.listIcon}>
            {iconc}
          </ListItemIcon>
          <ListItemText
            classes={{
              primary: classes.itemPrimary,
              textDense: classes.textDense,
            }}>
            {idc}
          </ListItemText>
        </div>
      );
      if(linkc){
        linkContent= (
        <Link href={hrefc} prefetch>
        {linkContent}
        </Link>
        );
      }
      return linkContent;
    }

    render() {
        const { classes, updatepagepathandtitle, ...other } = this.props;
        const { path } = this.state;
        this.updateCategoriesMenus();
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
                    {categories.map(({ id: childId, icon, href, show, link, children }) => {
                      if (typeof show !== 'undefined' && !show){
                        return '';
                      }
                      return (
                        <React.Fragment>
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
                            <Link href={link?href:''} prefetch>
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
                        {this.renderChildren(children, 1)}
                        </React.Fragment>
                        );
                      })}
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

const mapStateToProps = state => {
  // const k8sconfig = state.get("k8sConfig").toJS();
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  // const grafana = state.get("grafana").toJS();
  // const prometheus = state.get("prometheus").toJS();
  // return {meshAdapters, meshAdaptersts, k8sconfig, grafana, prometheus};
  return {meshAdapters, meshAdaptersts};
}

export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
  )(withRouter(Navigator)));