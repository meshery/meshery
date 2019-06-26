import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import {connect} from "react-redux";
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import PhoneIcon from '@material-ui/icons/Phone';
import FavoriteIcon from '@material-ui/icons/Favorite';
import PersonPinIcon from '@material-ui/icons/PersonPin';
import HelpIcon from '@material-ui/icons/Help';
import ShoppingBasket from '@material-ui/icons/ShoppingBasket';
import ThumbDown from '@material-ui/icons/ThumbDown';
import ThumbUp from '@material-ui/icons/ThumbUp';
import Typography from '@material-ui/core/Typography';
import { Paper } from '@material-ui/core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPoll, faCloudMeatball } from '@fortawesome/free-solid-svg-icons';
// import {  } from '@fortawesome/free-regular-svg-icons';
import { faMendeley } from '@fortawesome/free-brands-svg-icons';
import MeshConfigComponent from './MeshConfigComponent';
import GrafanaComponent from './GrafanaComponent';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import PrometheusComponent from './PrometheusComponent';
// import { faArrowRight, faTrashAlt, faTerminal, faExternalLinkSquareAlt, faPoll } from '@fortawesome/free-regular-svg-icons';

const styles = theme => ({
  root: {
    flexGrow: 1,
    width: '100%',
    // backgroundColor: theme.palette.background.paper,
  },
  icon: {
    width: theme.spacing(2.5),
  },
});

function TabContainer(props) {
  return (
    <Typography component="div" style={{ padding: 8 * 3 }}>
      {props.children}
    </Typography>
  );
}

TabContainer.propTypes = {
  children: PropTypes.node.isRequired,
};

class MesherySettings extends React.Component {
  constructor(props){
    super(props);
    const {k8sconfig, meshAdapters, grafana, prometheus} = props;
    this.state = {
      completed: {},
      k8sconfig, 
      meshAdapters, 
      grafana,
      prometheus,
      tabVal: 0,
      subTabVal: 0,
    };
  }

  handleChange(val) {
    const self = this;
    return (event, newVal) => {
      if(val === 'tabVal'){
        self.setState({tabVal: newVal});
      } else if(val === 'subTabVal'){
        self.setState({subTabVal: newVal});
      }
    }
  }

  render() {
    const { classes } = this.props;
    const {tabVal, subTabVal} = this.state;

    const mainIconScale = "grow-10"
    
    return (
      <div className={classes.root}>
      <Paper square className={classes.root}>
        <Tabs
          value={tabVal}
          onChange={this.handleChange("tabVal")}
          variant="fullWidth"
          indicatorColor="secondary"
          textColor="secondary"
        >
          <Tab icon={
            <FontAwesomeIcon icon={faCloudMeatball} transform={mainIconScale} fixedWidth />
          } label="Environment" />
          <Tab icon={
            <FontAwesomeIcon icon={faMendeley} transform={mainIconScale} fixedWidth />
          } label="Service Meshes" />
          <Tab icon={
            <FontAwesomeIcon icon={faPoll} transform={mainIconScale} fixedWidth />
          } label="Metrics" />
        </Tabs>
      </Paper>
      {tabVal === 0 && <TabContainer>
        <MeshConfigComponent />
      </TabContainer>}
      {tabVal === 1 && <TabContainer>
        <MeshAdapterConfigComponent />
      </TabContainer>}
      {tabVal === 2 && 
        <TabContainer>
          <Paper square className={classes.root}>
            <Tabs
              value={subTabVal}
              onChange={this.handleChange("subTabVal")}
              indicatorColor="secondary"
              textColor="secondary"
              variant="fullWidth"
            >
              <Tab icon={
                <img src="/static/img/grafana_icon.svg" className={classes.icon} />
              } label="Grafana" />
              <Tab icon={
                <img src="/static/img/prometheus_logo_orange_circle.svg" className={classes.icon} />
              } label="Prometheus" />
            </Tabs>
          </Paper>
          {subTabVal === 0 && <TabContainer>
            <GrafanaComponent />
          </TabContainer>}
          {subTabVal === 1 && <TabContainer>
            <PrometheusComponent />
          </TabContainer>}
        </TabContainer>}
    </div>
    );
  }
}

const mapStateToProps = state => {
  const k8sconfig = state.get("k8sConfig").toJS();
  const meshAdapters = state.get("meshAdapters").toJS();
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  return {k8sconfig, meshAdapters, grafana, prometheus};
}

MesherySettings.propTypes = {
  classes: PropTypes.object,
};

export default withStyles(styles)(connect(
  mapStateToProps
)(MesherySettings));