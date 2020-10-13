import React from 'react';
import { makeStyles, useTheme, withStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import IconButton from '@material-ui/core/IconButton';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import ChevronRightIcon from '@material-ui/icons/ChevronRight';
import Tab from '@material-ui/core/Tab'
import Tabs from '@material-ui/core/Tabs';
import TrafficIcon from '@material-ui/icons/Traffic';
import SecurityIcon from '@material-ui/icons/Security';
import InfoIcon from '@material-ui/icons/Info';
import CloseIcon from '@material-ui/icons/Close';
import { NetworkIcon } from '@patternfly/react-icons';
// import clsx from 'clsx';

const useStyles = makeStyles(() => ({
  root: {
    position: 'relative',
    flexGrow: 1,
  },
  tabroot: {
    minWidth: '20%',
  },
  paper: {
    background: 'white',
    width: '25%'
  },
  flex: {
    flex: 1,
  },
  list: {
    paddingTop: 70,
  },
  hide: {
    display: 'none',
  },
}));

function PersistentDrawerRight(props) {
  const classes = useStyles();
  const theme = useTheme();
  const [value, setValue] = React.useState(0);
  const [open, setOpen] = React.useState(props.open);
  console.log(props.open);
  console.log(props.data.data('app'));
  
  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  const handleDrawerClose = () => {
    setOpen(false);
    props.toggle();
  };

  return (
    <div className={classes.root}>
      <Drawer
        classes={{paper: classes.paper}}
        variant="persistent"
        anchor="right"
        open={open}
      >
        <div className={classes.list}>
          <IconButton onClick={handleDrawerClose}>
            {theme.direction === 'rtl' ? <ChevronLeftIcon /> : <ChevronRightIcon />}
          </IconButton>
          <Tabs
            value={value}
            onChange={handleChange}
            indicatorColor="primary"
            variant="scrollable"
            scrollButtons="auto"
            aria-label="scrollable auto tabs example"
          >
            <Tab icon={<InfoIcon />} className={classes.tabroot} />
            <Tab icon={<TrafficIcon />} className={classes.tabroot} />
            <Tab icon={<SecurityIcon />} className={classes.tabroot}/>
            <Tab icon={<CloseIcon />} className={classes.tabroot}/>
            <Tab icon={<NetworkIcon />} className={classes.tabroot}/>
          </Tabs>
        </div>
      </Drawer>
    </div>
  );
}

export default withStyles(useStyles)(PersistentDrawerRight);
