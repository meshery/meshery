import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { NoSsr, Chip, IconButton, Card, CardContent, Typography, CardHeader, Tooltip } from '@material-ui/core';
import dataFetch from '../lib/data-fetch';
import blue from '@material-ui/core/colors/blue';
import { updateProgress } from '../lib/store';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import { withRouter } from 'next/router';
import { withSnackbar } from 'notistack';
import CloseIcon from '@material-ui/icons/Close';
// import DoneIcon from '@material-ui/icons/Done';


const styles = theme => ({
  root: {
    padding: theme.spacing(5),
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  button: {
    marginTop: theme.spacing(3),
    marginLeft: theme.spacing(1),
  },
  margin: {
    margin: theme.spacing(1),
  },
  alreadyConfigured: {
    textAlign: 'center',
    padding: theme.spacing(20),
  },
  colorSwitchBase: {
    color: blue[300],
    '&$colorChecked': {
      color: blue[500],
      '& + $colorBar': {
        backgroundColor: blue[500],
      },
    },
  },
  colorBar: {},
  colorChecked: {},
  fileLabel: {
    width: '100%',
  },
  fileLabelText: {
  },
  inClusterLabel: {
    paddingRight: theme.spacing(2),
  },
  alignCenter: {
    textAlign: 'center',
  },
  //   alignRight: {
  //     textAlign: 'right',
  //     marginBottom: theme.spacing(2),
  //   },
  //   fileInputStyle: {
  //     opacity: '0.01',
  //   },
  icon: {
    width: theme.spacing(2.5),
  },
  istioIcon: {
    width: theme.spacing(1.5),
  },
  cardHeader: {
    fontSize: theme.spacing(2),
  },
  card: {
    height: '100%',
  },
  cardContent: {
    height: '100%',
  },
});

class DashboardComponent extends React.Component {

  constructor(props) {
    super(props);
    const {meshAdapters, k8sconfig, grafana, prometheus} = props;
    this.state = {
      meshAdapters,
      availableAdapters: [],
      mts: new Date(),
      meshLocationURLError: false,

      inClusterConfig: k8sconfig.inClusterConfig, // read from store
      k8sfile: k8sconfig.k8sfile, // read from store
      contextName: k8sconfig.contextName, // read from store
        
      clusterConfigured: k8sconfig.clusterConfigured, // read from store
      configuredServer: k8sconfig.configuredServer,
      k8sfileError: false,
      kts: new Date(),

      grafana,
      prometheus,
    };
  }
    
  static getDerivedStateFromProps(props, state){
    const { meshAdapters, meshAdaptersts, k8sconfig, grafana, prometheus } = props;
    const st = {};
    if(meshAdaptersts > state.mts) {
      st.meshAdapters = meshAdapters;
      st.mts = meshAdaptersts;
    }
    if(k8sconfig.ts > state.kts){
      st.inClusterConfig = k8sconfig.inClusterConfig;
      st.k8sfile = k8sconfig.k8sfile;
      st.contextName = k8sconfig.contextName;
      st.clusterConfigured = k8sconfig.clusterConfigured;
      st.configuredServer = k8sconfig.configuredServer;
      st.kts = props.ts;
    }

    st.grafana = props.grafana;
    st.prometheus = props.prometheus;

    return st;
  }

      componentDidMount = () => {
        this.fetchAvailableAdapters();
      }
    
      fetchAvailableAdapters = () => {
        const self = this;
        this.props.updateProgress({showProgress: true});
        dataFetch('/api/mesh/adapters', { 
          credentials: 'same-origin',
          method: 'GET',
          credentials: 'include',
        }, result => {
          this.props.updateProgress({showProgress: false});
          if (typeof result !== 'undefined'){
            const options = result.map(res => {
              return {
                value: res,
                label: res,
              };
            });
            this.setState({availableAdapters: options});
          }
        }, self.handleError("Unable to fetch list of adapters."));
      }

  handleError = (msg) => error => {
    this.props.updateProgress({showProgress: false});
    const self = this;
    this.props.enqueueSnackbar(`${msg}: ${error}`, {
      variant: 'error',
      action: (key) => (
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          onClick={() => self.props.closeSnackbar(key) }
        >
          <CloseIcon />
        </IconButton>
      ),
      autoHideDuration: 7000,
    });
  }

  handleDelete() {
    return false;
  }

  handleAdapterClick = (adapterLoc) => () => {
    // const { meshAdapters } = this.state;
    this.props.updateProgress({showProgress: true});
    let self = this;
    dataFetch(`/api/mesh/adapter/ping?adapter=${encodeURIComponent(adapterLoc)}`, { 
      credentials: 'same-origin',
      credentials: 'include',
    }, result => {
      this.props.updateProgress({showProgress: false});
      if (typeof result !== 'undefined'){
        this.props.enqueueSnackbar('Adapter successfully pinged!', {
          variant: 'success',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key) }
            >
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    }, self.handleError("Could not ping adapter."));
  }

  handleKubernetesClick = () => {
    this.props.updateProgress({showProgress: true});
    let self = this;
    dataFetch(`/api/k8sconfig/ping`, { 
      credentials: 'same-origin',
      credentials: 'include',
    }, result => {
      this.props.updateProgress({showProgress: false});
      if (typeof result !== 'undefined'){
        this.props.enqueueSnackbar('Kubernetes successfully pinged!', {
          variant: 'success',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key) }
            >
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    }, self.handleError("Could not ping Kubernetes."));
  }

  handleGrafanaClick = () => {
    this.props.updateProgress({showProgress: true});
    let self = this;
    dataFetch(`/api/grafana/ping`, { 
      credentials: 'same-origin',
      credentials: 'include',
    }, result => {
      this.props.updateProgress({showProgress: false});
      if (typeof result !== 'undefined'){
        this.props.enqueueSnackbar('Grafana successfully pinged!', {
          variant: 'success',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key) }
            >
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    }, self.handleError("Could not ping Grafana."));
  }
  
  handlePrometheusClick = () => {
    this.props.updateProgress({showProgress: true});
    let self = this;
    dataFetch(`/api/prometheus/ping`, { 
      credentials: 'same-origin',
      credentials: 'include',
    }, result => {
      this.props.updateProgress({showProgress: false});
      if (typeof result !== 'undefined'){
        this.props.enqueueSnackbar('Prometheus successfully pinged!', {
          variant: 'success',
          autoHideDuration: 2000,
          action: (key) => (
            <IconButton
              key="close"
              aria-label="Close"
              color="inherit"
              onClick={() => self.props.closeSnackbar(key) }
            >
              <CloseIcon />
            </IconButton>
          ),
        });
      }
    }, self.handleError("Could not ping Prometheus."));
  }

  showCard(title, content) {
    const { classes } = this.props;
    return (
      <Card className={classes.card}>
        <CardHeader
          disableTypography={true}
          title={title}
          // action={iconComponent}
          className={classes.cardHeader}
        />
        <CardContent className={classes.cardContent}>
          {content}
        </CardContent>
      </Card>
    )
  }


  configureTemplate = () => {
    const { classes } = this.props;
    const { inClusterConfig, contextName, clusterConfigured, configuredServer, meshAdapters, availableAdapters, grafana, prometheus } = this.state;
    const self = this;
    let showConfigured = 'Not connected to Kubernetes.';
    if (clusterConfigured) {

      let chp = (
        <Chip 
          // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
          label={inClusterConfig?'Using In Cluster Config': contextName}
          // onDelete={self.handleDelete}
          // deleteIcon={<DoneIcon />}
          onClick={self.handleKubernetesClick}
          icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />} 
          className={classes.chip}
          key='k8s-key'
          variant="outlined" />
      );

      if(configuredServer){
        chp = (
          <Tooltip title={`Server: ${configuredServer}`}>
            {chp}
          </Tooltip>
        );
      }

      showConfigured = (
        <div className={classes.alignRight}>
          {chp}  
        </div>
      )
    }

    let showAdapters = 'No adapters configured.';
    if (availableAdapters.length > 0) {

      availableAdapters.sort((a1, a2) => (a1.value < a2.value?-1:(a1.value > a2.value?1:0)));

      showAdapters = (
        <div>
          {
            availableAdapters.map((aa, ia) => {
              let isDisabled = true;
              let image = "/static/img/meshery-logo.png";
              let logoIcon = (<img src={image} className={classes.icon} />);       
              let adapterType = '';
              meshAdapters.forEach((adapter, ind) => {
                if(aa.value === adapter.adapter_location){
                  isDisabled = false;
                  adapterType = adapter.name;
                  switch (adapter.name.toLowerCase()){
                  case 'istio':
                    image = "/static/img/istio-blue.svg";
                    logoIcon = (<img src={image} className={classes.istioIcon} />);
                    break;
                  case 'linkerd':
                    image = "/static/img/linkerd.svg";
                    logoIcon = (<img src={image} className={classes.icon} />);
                    break;
                  case 'consul':
                    image = "/static/img/consul.svg";
                    logoIcon = (<img src={image} className={classes.icon} />);
                    break;
                  case 'network service mesh':
                    image = "/static/img/nsm.svg";
                    logoIcon = (<img src={image} className={classes.icon} />);
                    break;
                  case 'octarine':
                    image = "/static/img/octarine.svg";
                    logoIcon = (<img src={image} className={classes.icon} />);
                    break;
                  case 'cpx':
                    image = "/static/img/citrix.svg";
                    logoIcon = (<img src={image} className={classes.icon} />);
                    break;
                        
                  }
                }
              });
                

                
              return (
                <Tooltip title={isDisabled?"This adapter is inactive":
                   `${adapterType
                     .toLowerCase()
                     .split(' ')
                     .map((s) => s.charAt(0).toUpperCase() + s.substring(1))
                     .join(' ')} adapter on port ${aa.label.split(':')[1]}`}>
                  <Chip 
                    label={aa.label.split(':')[0]}
                    // onDelete={self.handleDelete(ind)} 
                    // onDelete={!isDisabled?self.handleDelete:null}
                    // deleteIcon={!isDisabled?<DoneIcon />:null}
                    onClick={self.handleAdapterClick(aa.value)}
                    icon={logoIcon}
                    className={classes.chip}
                    key={`adapters-${ia}`}
                    variant={isDisabled?"default":"outlined"} />
                </Tooltip>
              );
            })}
        </div>
      )

    }

    let showGrafana = 'Not connected to Grafana.';
    if(grafana && grafana.grafanaURL&& grafana.grafanaURL !== ''){
      showGrafana = (
        <Chip 
          label={grafana.grafanaURL}
          // onDelete={self.handleDelete}
          // deleteIcon={<DoneIcon />}
          onClick={self.handleGrafanaClick}
          icon={<img src="/static/img/grafana_icon.svg" className={classes.icon} />} 
          className={classes.chip}
          key='graf-key'
          variant="outlined" />
      );
    }

    let showPrometheus = 'Not connected to Prometheus.';
    if(prometheus && prometheus.prometheusURL&& prometheus.prometheusURL !== ''){
      showPrometheus = (
        <Chip 
          label={prometheus.prometheusURL}
          // onDelete={self.handleDelete}
          // deleteIcon={<DoneIcon />}
          onClick={self.handlePrometheusClick}
          icon={<img src="/static/img/prometheus_logo_orange_circle.svg" className={classes.icon} />} 
          className={classes.chip}
          key='prom-key'
          variant="outlined" />
      );
    }



    return (
      <NoSsr>
        <div className={classes.root}>
          <Typography variant="h6" gutterBottom className={classes.chartTitle}>
          Connection Status
          </Typography>  
                
          <Grid container spacing={1}>
            {/* <Grid item xs={12} sm={2}>
        Kubernetes
        </Grid>
        <Grid item xs={12} sm={10}>
        {showConfigured}
        </Grid>
        <Grid item xs={12} sm={2}>
        Adapters
        </Grid>
        <Grid item xs={12} sm={10}>
        {showAdapters}
        </Grid>
        <Grid item xs={12} sm={2}>
        Grafana
        </Grid>
        <Grid item xs={12} sm={10}>
        {showGrafana}
        </Grid>
        <Grid item xs={12} sm={2}>
        Prometheus
        </Grid>
        <Grid item xs={12} sm={10}>
        {showPrometheus}
        </Grid> */}

            <Grid item xs={12} sm={6}>
              {self.showCard('Kubernetes', showConfigured)}
            </Grid>
            <Grid item xs={12} sm={6}>
              {self.showCard('Adapters', showAdapters)}
            </Grid>
            <Grid item xs={12} sm={6}>
              {self.showCard('Grafana', showGrafana)}
            </Grid>
            <Grid item xs={12} sm={6}>
              {self.showCard('Prometheus', showPrometheus)}
            </Grid>

          </Grid>
        </div>
      </NoSsr>
    );
    
  }

  render() {
    // const { reconfigureCluster } = this.state;
    // if (reconfigureCluster) {
    return this.configureTemplate();
    // }
    // return this.alreadyConfiguredTemplate();
  }
}

DashboardComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
    // updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
    updateProgress: bindActionCreators(updateProgress, dispatch),
  }
}
const mapStateToProps = state => {
  const k8sconfig = state.get("k8sConfig").toJS();
  const meshAdapters = state.get("meshAdapters").toJS();
  const meshAdaptersts = state.get("meshAdaptersts");
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  return {meshAdapters, meshAdaptersts, k8sconfig, grafana, prometheus};
}

export default withStyles(styles)(connect(
  mapStateToProps,
  mapDispatchToProps
)(withRouter(withSnackbar(DashboardComponent))));
