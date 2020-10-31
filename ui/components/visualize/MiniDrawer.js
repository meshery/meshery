import React, { Component } from 'react';
import { withStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import TrafficIcon from '@material-ui/icons/Traffic';
import SecurityIcon from '@material-ui/icons/Security';
import InfoIcon from '@material-ui/icons/Info';
import CloseIcon from '@material-ui/icons/Close';
import { NetworkIcon } from '@patternfly/react-icons';
import { Divider } from '@material-ui/core';

const useStyles = () => ({
  root: {
    display: 'flex',
  },
  drawer: {
    width: 'auto',
    flexShrink: 0,
  },
  drawerPaper: {
    width: 'auto',
  },
  items: {
    paddingTop: 70,
  }
});

class PermanentDrawerRight extends Component{
  constructor(props){
    super(props);
  }

  handleDrawer = (tab) => {
    const { toggle } = this.props;
    toggle(tab);
  };

  render(){
    const { classes } = this.props;

    return (
      <div className={classes.root}>
        <Drawer
          className={classes.drawer}
          variant="permanent"
          classes={{
            paper: classes.drawerPaper,
          }}
          anchor="right"
        >
          <div className={classes.items}>
            <List >
              <ListItem button onClick={() => this.handleDrawer(0)}>
                <ListItemIcon>
                  <ChevronLeftIcon />
                </ListItemIcon>
              </ListItem>
              <Divider />
              <ListItem button onClick={() => this.handleDrawer(0)}>
                <ListItemIcon>
                  <InfoIcon />
                </ListItemIcon>
              </ListItem>
              <ListItem button onClick={() => this.handleDrawer(1)}>
                <ListItemIcon>
                  <TrafficIcon />
                </ListItemIcon>
              </ListItem>
              <ListItem button onClick={() => this.handleDrawer(2)}>
                <ListItemIcon>
                  <SecurityIcon />
                </ListItemIcon>
              </ListItem>
              <ListItem button onClick={() => this.handleDrawer(3)}>
                <ListItemIcon>
                  <CloseIcon />
                </ListItemIcon>
              </ListItem>
              <ListItem button onClick={() => this.handleDrawer(4)}>
                <ListItemIcon>
                  <NetworkIcon />
                </ListItemIcon>
              </ListItem>
            </List>
          </div>
        </Drawer>
      </div>
    );
  }
}

export default withStyles(useStyles, {withTheme: true})(PermanentDrawerRight);
