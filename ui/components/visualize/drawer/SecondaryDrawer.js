import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withSnackbar } from 'notistack';
import { withStyles } from '@material-ui/core/styles';
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
import Grid from '@material-ui/core/Grid';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import OpenInNewIcon from '@material-ui/icons/OpenInNew';
import { Divider, Link, Paper } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import GrafanaCustomCharts from '../../GrafanaCustomCharts';
// let bb = require('billboard.js');
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
        <Box padding={1}>
          {children}
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

const useStyles = () => ({
  root: {
    position: 'relative',
    flexGrow: 1,
  },
  tabroot: {
    minWidth: '20%',
  },
  paper: {
    background: 'white',
    width: '25%',
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
  metrics: {
    marginTop: 20,
  },
  actions: {
    marginTop: 20,
  },
  smp: {
    padding: '10px',
    'text-align': 'center'
  },
  button: {
    minWidth: '100%',
    display: 'block',
  },
  icon: {
    maxWidth: '80%'
  },
  modal: {
    maxWidth: '60%',
    maxHeight: '60%',
    margin: 'auto',
  },
  modalPaper: {
    backgroundColor: 'white',
  }
});


class SecondaryDrawer extends Component {
  constructor(props){
    super(props);
    this.state = {
      value: this.props.tab,
    };
  }

  render(){
    const { classes, open, data, theme, } = this.props;
    const {
      grafanaURL, grafanaAPIKey, selectedBoardsConfigs,
    } = this.props.grafana;
    const {
      value,
    } = this.state;

    if(data){
      console.log(data.data('app'));
    }

    const handleChange = (event, newValue) => {
      this.setState({ value: newValue});
    };

    const handleDrawerClose = () => {
      const { toggle } = this.props;
      toggle(null, false);
    };

    return (
      <div className={classes.root}>
        <Drawer
          classes={{ paper: classes.paper }}
          variant="persistent"
          anchor="right"
          open={open}
          PaperProps={{ elevation: 10}}
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
              {/* <Typography variant="h6">
                {data.data('app')}
              </Typography> */}
              <Paper className={classes.actions}>
                <AppBar style={{position: 'relative'}}>
                  <Toolbar>
                    <Typography component={'div'} variant="h6">
                      Actions
                    </Typography>
                  </Toolbar>
                </AppBar>
                <Typography style={{padding: '5px', textAlign: 'center'}} variant="h6">
                  Performance
                </Typography>
                <Grid className={classes.smp} container spacing={2}>
                  <Grid item xs={4}>
                    <Button className={classes.button} onClick={this.props.togglePeformanceModal}>
                      <img className={classes.icon} src="/static/img/smp-dark.svg"/>
                        Adhoc Performance Test
                    </Button>
                  </Grid>
                  <Grid item xs={4}>
                    <Button className={classes.button}>
                      <img className={classes.icon} src="/static/img/smp-dark.svg"/>
                        Select Performance Profile
                    </Button>
                  </Grid>
                  <Grid item xs={4}>
                    <Button className={classes.button}>
                      <img className={classes.icon} src="/static/img/smp-dark.svg"/>
                        Some Performance Test
                    </Button>
                  </Grid>
                </Grid>
                <Divider variant="middle"/>
                <Typography style={{padding: '5px', textAlign: 'center'}} variant="h6">
                  Conformance
                </Typography>
                <Grid className={classes.smp} container spacing={2}>
                  <Grid item xs={4}>
                    <Button className={classes.button}>
                      <img className={classes.icon} src="/static/img/smc-checklist.svg"/>
                      Start Conformance Test
                    </Button>
                  </Grid>
                </Grid>
              </Paper>
            </TabPanel>
            <TabPanel value={value} index={1}>
              <Paper className={classes.metrics}>
                <AppBar style={{position: 'relative'}}>
                  <Toolbar>
                    <Typography component={'div'} variant="h6">
                        Metrics
                    </Typography>
                    <Link href="http://localhost:3000" target="_blank">
                      <IconButton>
                        <OpenInNewIcon />  
                      </IconButton>
                    </Link>
                  </Toolbar>
                </AppBar>
                <GrafanaCustomCharts
                  boardPanelConfigs={selectedBoardsConfigs}
                  grafanaURL={grafanaURL}
                  grafanaAPIKey={grafanaAPIKey}
                />
              </Paper>
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
}

const mapStateToProps = (st) => {
  const grafana = st.get('grafana').toJS();
  return { grafana };
};

// export default withStyles(useStyles, {withTheme: true})(PersistentDrawerRight);
export default withStyles(useStyles, {withTheme: true})(connect(
  mapStateToProps,
)(withSnackbar(SecondaryDrawer)));