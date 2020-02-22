import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import Typography from '@material-ui/core/Typography';
import { AppBar, Paper, Tooltip } from '@material-ui/core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft, faCloud, faPoll } from '@fortawesome/free-solid-svg-icons';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import MeshConfigComponent from './MeshConfigComponent';
import GrafanaComponent from './GrafanaComponent';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import PrometheusComponent from './PrometheusComponent';
import Link from 'next/link';

const styles = theme => ({
  root: {
    flexGrow: 1,
    width: '100%'
    // backgroundColor: theme.palette.background.paper,
  },
  icon: {
    display: 'inline',
    verticalAlign: 'text-top',
    width: theme.spacing(1.75),
    marginLeft: theme.spacing(0.5)
  },
  iconText: {
    display: 'inline',
    verticalAlign: 'middle'
  },
  backToPlay: {
    margin: theme.spacing(2)
  },
  link: {
    cursor: 'pointer'
  }
});

function TabContainer (props) {
  return (
    <Typography component="div" style={{
      // paddingLeft: 8*3,
      // paddingRight: 8*3,
      // paddingBottom: 8*3,
      paddingTop: 2
    }}>
      {props.children}
    </Typography>
  );
}

TabContainer.propTypes = {
  children: PropTypes.node.isRequired
};

class MesherySettings extends React.Component {
  constructor (props) {
    super(props);
    const { k8sconfig, meshAdapters, grafana, prometheus } = props;
    this.state = {
      completed: {},
      k8sconfig,
      meshAdapters,
      grafana,
      prometheus,
      tabVal: 0,
      subTabVal: 0
    };
  }

  static getDerivedStateFromProps (props, state) {
    if (JSON.stringify(props.k8sconfig) !== JSON.stringify(state.k8sconfig) ||
        JSON.stringify(props.meshAdapters) !== JSON.stringify(state.meshAdapters)) {
      return {
        k8sconfig: props.k8sconfig,
        meshAdapters: props.meshAdapters,
        grafana: props.grafana,
        prometheus: props.prometheus
      };
    }
    return null;
  }

  handleChange (val) {
    const self = this;
    return (event, newVal) => {
      if (val === 'tabVal') {
        self.setState({ tabVal: newVal });
      } else if (val === 'subTabVal') {
        self.setState({ subTabVal: newVal });
      }
    };
  }

  render () {
    const { classes } = this.props;
    const { tabVal, subTabVal, k8sconfig, meshAdapters } = this.state;

    const mainIconScale = 'grow-10';
    let backToPlay = '';
    if (k8sconfig.clusterConfigured === true && meshAdapters.length > 0) {
      backToPlay =
        <div className={classes.backToPlay}>
          <Link href={'/management'}>
            <div className={classes.link}>
              <FontAwesomeIcon icon={faArrowLeft} transform="grow-4" fixedWidth /> You are all set to manage service meshes
            </div>
          </Link>
        </div>;
}
    return (
      <div className={classes.root}>
        <Paper square className={classes.root}>
          <Tabs
            value={tabVal}
            onChange={this.handleChange('tabVal')}
            variant="fullWidth"
            indicatorColor="primary"
            textColor="primary"
          >
            <Tooltip title="Identify your cluster" placement="top">
              <Tab icon={
                <FontAwesomeIcon icon={faCloud} transform={mainIconScale} fixedWidth />
              } label="Environment" />
            </Tooltip>
            <Tooltip title="Connect Meshery Adapters" placement="top">
              <Tab icon={
                <FontAwesomeIcon icon={faMendeley} transform={mainIconScale} fixedWidth />
              } label="Service Meshes" />
            </Tooltip>
            <Tooltip title="Configure Metrics backends" placement="top">
              <Tab icon={
                <FontAwesomeIcon icon={faPoll} transform={mainIconScale} fixedWidth />
              } label="Metrics" />
            </Tooltip>
          </Tabs>
        </Paper>
        {tabVal === 0 && <TabContainer>
          <AppBar position="static" color="default">
            <Tabs
              value={subTabVal}
              onChange={this.handleChange('subTabVal')}
              indicatorColor="primary"
              textColor="primary"
              variant="fullWidth"
            >
              <Tab label="Out of Cluster Deployment" />
              <Tab label="In Cluster Deployment"/>
            </Tabs>
          </AppBar>
          {subTabVal === 0 && <TabContainer>
            <MeshConfigComponent tabs={subTabVal} />
          </TabContainer>}
          {subTabVal === 1 && <TabContainer>
            <MeshConfigComponent tabs={subTabVal} />
          </TabContainer>}
        </TabContainer>}
        {tabVal === 1 && <TabContainer>
          <MeshAdapterConfigComponent />
        </TabContainer>}
        {tabVal === 2 &&
        <TabContainer>
          <AppBar position="static" color="default">
            <Tabs
              value={subTabVal}
              onChange={this.handleChange('subTabVal')}
              indicatorColor="primary"
              textColor="primary"
              variant="fullWidth"
            >
              <Tab label={
                <div className={classes.iconText}>Grafana
                  <img src="/static/img/grafana_icon.svg" className={classes.icon} />
                </div>
              } />
              <Tab label={
                <div className={classes.iconText}>Prometheus
                  <img src="/static/img/prometheus_logo_orange_circle.svg" className={classes.icon} />
                </div>
              } />
            </Tabs>
          </AppBar>
          {subTabVal === 0 && <TabContainer>
            <GrafanaComponent />
          </TabContainer>}
          {subTabVal === 1 && <TabContainer>
            <PrometheusComponent />
          </TabContainer>}
        </TabContainer>}

        {backToPlay}
      </div>
    );
  }
}

const mapStateToProps = state => {
  const k8sconfig = state.get('k8sConfig').toJS();
  const meshAdapters = state.get('meshAdapters').toJS();
  const grafana = state.get('grafana').toJS();
  const prometheus = state.get('prometheus').toJS();
  return { k8sconfig,
meshAdapters,
grafana,
prometheus };
};

MesherySettings.propTypes = {
  classes: PropTypes.object
};

export default withStyles(styles)(connect(mapStateToProps)(MesherySettings));