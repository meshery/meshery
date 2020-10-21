import React from 'react';
import { makeStyles, useTheme, withStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';
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
import Typography from '@material-ui/core/Typography';
import Box from '@material-ui/core/Box';
// import clsx from 'clsx';

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box p={3}>
          <Typography>{children}</Typography>
        </Box>
      )}
    </div>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
};

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

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
        classes={{ paper: classes.paper }}
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
            <Tab icon={<InfoIcon />} className={classes.tabroot} {...a11yProps(0)} />
            <Tab icon={<TrafficIcon />} className={classes.tabroot} {...a11yProps(1)} />
            <Tab icon={<SecurityIcon />} className={classes.tabroot} {...a11yProps(2)} />
            <Tab icon={<CloseIcon />} className={classes.tabroot} {...a11yProps(3)} />
            <Tab icon={<NetworkIcon />} className={classes.tabroot} {...a11yProps(4)} />
          </Tabs>
          <TabPanel value={value} index={0}>
            Item One
          </TabPanel>
          <TabPanel value={value} index={1}>
            Item Two    
          </TabPanel>
          <TabPanel value={value} index={2}>
            Item Three
          </TabPanel>
          <TabPanel value={value} index={3}>
            Item Four    
          </TabPanel>
          <TabPanel value={value} index={4}>
            Item Five
          </TabPanel>
        </div>
      </Drawer>
    </div>
  );
}

export default withStyles(useStyles)(PersistentDrawerRight);
