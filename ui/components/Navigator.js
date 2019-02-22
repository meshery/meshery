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
import HomeIcon from '@material-ui/icons/Home';
// import PeopleIcon from '@material-ui/icons/People';
// import DnsRoundedIcon from '@material-ui/icons/DnsRounded';
// import PermMediaOutlinedIcon from '@material-ui/icons/PhotoSizeSelectActual';
// import PublicIcon from '@material-ui/icons/Public';
// import SettingsEthernetIcon from '@material-ui/icons/SettingsEthernet';
import SettingsInputComponentIcon from '@material-ui/icons/SettingsInputComponent';
import TimerIcon from '@material-ui/icons/Timer';
import SettingsIcon from '@material-ui/icons/Settings';
// import PhonelinkSetupIcon from '@material-ui/icons/PhonelinkSetup';
import Link from "next/link";
import {connect} from "react-redux";
import { bindActionCreators } from 'redux'
import { updatepagepathandtitle } from '../lib/store';

const categories = [
    { id: 'Setup Mesh', icon: <SettingsIcon />, href: "/", title: 'Setup Mesh' },
    { id: 'Play', icon: <TimerIcon />, href: "/about", title: 'Play with Mesh' },
    { id: 'Load Test', icon: <SettingsInputComponentIcon />, href: "/post", title: 'Load Test and Charts' },
]



// [
//   {
//     // id: 'Develop',
//     children: [
//     //   { id: 'Authentication', icon: <PeopleIcon />, active: true },
//     //   { id: 'Database', icon: <DnsRoundedIcon /> },
//     //   { id: 'Storage', icon: <PermMediaOutlinedIcon /> },
//     //   { id: 'Hosting', icon: <PublicIcon /> },
//     //   { id: 'Functions', icon: <SettingsEthernetIcon /> },
//     //   { id: 'ML Kit', icon: <SettingsInputComponentIcon /> },
//     ],
//   },
// //   {
// //     id: 'Quality',
// //     children: [
// //       { id: 'Analytics', icon: <SettingsIcon /> },
// //       { id: 'Performance', icon: <TimerIcon /> },
// //       { id: 'Test Lab', icon: <PhonelinkSetupIcon /> },
// //     ],
// //   },
// ];

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
    backgroundColor: '#232f3e',
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
});

class Navigator extends React.Component {
    // constructor(props){
    //     super(props);
    //     this.state = {
    //         path: null,
    //     };
    // }
    // static async getInitialProps ({store}) {
    //     return {
    //         store: store,
    //     }
    // }

    updateTitle(){
        let path = (typeof window !== 'undefined' ? window.location.pathname : '');
        categories.map(({title, href}) => {
            if (path.lastIndexOf('/') > 0) {
                path = path.substring(0, path.lastIndexOf('/'));
            }
            if (href === path) {
                console.log("updating path: "+path+" and title: "+title);
                this.props.updatepagepathandtitle({path, title});
                return;
            }
        });
        return path;
    }

    render() {
        // accessing 'updatepagepathandtitle' to just keep it out of 'other'
        const { classes, updatepagepathandtitle, ...other } = this.props;
        const path = this.updateTitle();
        console.log("current page:" + path);
        return (
            <Drawer variant="permanent" {...other}>
            <List disablePadding>
                <ListItem className={classNames(classes.firebase, classes.item, classes.itemCategory)}>
                Meshery
                </ListItem>
                <ListItem 
                    button 
                    component="a"
                    href="https://layer5.io/meshery"
                    className={classNames(classes.item, classes.itemCategory)}>
                <ListItemIcon>
                    <HomeIcon />
                </ListItemIcon>
                <ListItemText
                    classes={{
                    primary: classes.itemPrimary,
                    }}
                >
                    Project Overview
                </ListItemText>
                </ListItem>
                {/* {categories.map(({ id, children }) => (
                <React.Fragment key={id}> */}
                    {/* <ListItem className={classes.categoryHeader}>
                    <ListItemText
                        classes={{
                        primary: classes.categoryHeaderPrimary,
                        }}
                    >
                        {id}
                    </ListItemText>
                    </ListItem> */}
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
                        <Link href={href}>
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
                {/* </React.Fragment> */}
                ))}
            </List>
            </Drawer>
        );
    }
}

Navigator.propTypes = {
  classes: PropTypes.object.isRequired,
};

// const mapStateToProps = state => ({ count: state.get('count') })

const mapDispatchToProps = dispatch => {
  return {
    updatepagepathandtitle: bindActionCreators(updatepagepathandtitle, dispatch)
  }
}

export default withStyles(styles)(connect(
    null,
    mapDispatchToProps
  )(Navigator));