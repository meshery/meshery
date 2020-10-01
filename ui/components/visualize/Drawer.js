import React from 'react';
import { makeStyles, useTheme, withStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import IconButton from '@material-ui/core/IconButton';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import ChevronRightIcon from '@material-ui/icons/ChevronRight';
import Tab from '@material-ui/core/Tab'
import Tabs from '@material-ui/core/Tabs';
// import clsx from 'clsx';


const useStyles = makeStyles(() => ({
  root: {
    position: 'relative',
    flexGrow: 1,
  },
  paper: {
    background: "white",
  },
  flex: {
    flex: 1,
  },
  list: {
    width: 750,
    paddingTop: 70,
  },
  hide: {
    display: 'none',
  },
}));

function PersistentDrawerRight(props) {
  const classes = useStyles();
  const theme = useTheme();
  const [open, setOpen] = React.useState(true);
  const [value, setValue] = React.useState(0);
  console.log(props.data.data('id'));
  
  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  const handleDrawerClose = () => {
    setOpen(false);
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
            <Tab label="Overview" />
            <Tab label="Traffic Management" />
            <Tab label="Security" />
            <Tab label="Circuit Breaker" />
            <Tab label="Gateway" />
          </Tabs>
        </div>
      </Drawer>
    </div>
  );
}

export default withStyles(useStyles)(PersistentDrawerRight);
