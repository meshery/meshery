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
import DoneIcon from '@material-ui/icons/Done';


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
        }, self.handleError("Unable to fetch available adapters"));
      }

//   constructor(props) {
//     super(props);
//     const {inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer } = props;
//     this.state = {
//         inClusterConfig, // read from store
//         k8sfile, // read from store
//         k8sfileElementVal: '',
//         contextName, // read from store
    
//         clusterConfigured, // read from store
//         configuredServer,
//         k8sfileError: false,
//         ts: new Date(),
//       };
//   }

//   static getDerivedStateFromProps(props, state){
//     const {inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer } = props;
//     // if(inClusterConfig !== state.inClusterConfig || clusterConfigured !== state.clusterConfigured || k8sfile !== state.k8sfile 
//         // || configuredServer !== state.configuredServer){
//     if(props.ts > state.ts){
//       return {
//         inClusterConfig,
//           k8sfile,
//           k8sfileElementVal: '',
//           contextName, 
//           clusterConfigured,
//           configuredServer,
//           ts: props.ts,
//       };
//     }
//     return {};
//   }

//   handleChange = name => {
//     const self = this;
//     return event => {
//       if (name === 'inClusterConfig'){
//         self.setState({ [name]: event.target.checked, ts: new Date() });
//         return;
//       }
//       if (name === 'k8sfile' && event.target.value !== ''){
//         self.setState({ k8sfileError: false, ts: new Date() });
//       }
//       if (name === 'k8sfile') {
//         self.setState({k8sfileElementVal: event.target.value, ts: new Date()});
//       }
//       self.setState({ [name]: event.target.value, ts: new Date() });
//     };
//   }

//   handleSubmit = () => {
//     const { inClusterConfig, k8sfile } = this.state;
//     if (!inClusterConfig && k8sfile === '') {
//         this.setState({k8sfileError: true});
//         return;
//     }
//     this.submitConfig()
//   }

//   submitConfig = () => {
//     const { inClusterConfig, k8sfile, contextName } = this.state;
//     const fileInput = document.querySelector('#k8sfile') ;
//     const formData = new FormData();
//     formData.append('inClusterConfig', inClusterConfig?"on":''); // to simulate form behaviour of a checkbox
//     if (!inClusterConfig) {
//         formData.append('contextName', contextName);
//         formData.append('k8sfile', fileInput.files[0]);
//     }
//     this.props.updateProgress({showProgress: true});
//     let self = this;
//     dataFetch('/api/k8sconfig', { 
//       credentials: 'same-origin',
//       method: 'POST',
//       credentials: 'include',
//       body: formData
//     }, result => {
//       this.props.updateProgress({showProgress: false});
//       if (typeof result !== 'undefined'){
//         this.setState({clusterConfigured: true, configuredServer: result.configuredServer, contextName: result.contextName});
//         this.props.enqueueSnackbar('Kubernetes config was successfully validated!', {
//           variant: 'success',
//           autoHideDuration: 2000,
//           action: (key) => (
//             <IconButton
//                   key="close"
//                   aria-label="Close"
//                   color="inherit"
//                   onClick={() => self.props.closeSnackbar(key) }
//                 >
//                   <CloseIcon />
//             </IconButton>
//           ),
//         });
//         this.props.updateK8SConfig({k8sConfig: {inClusterConfig, k8sfile, contextName: result.contextName, clusterConfigured: true, configuredServer: result.configuredServer}});
//       }
//     }, self.handleError);
//   }

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
      autoHideDuration: 8000,
    });
  }

//   handleTimerDialogClose = () => {
//     this.setState({timerDialogOpen: false});
//   }

//   handleReconfigure = () => {
// 	let self = this;
//     dataFetch('/api/k8sconfig', { 
//       credentials: 'same-origin',
//       method: 'DELETE',
//       credentials: 'include',
//     }, result => {
//       this.props.updateProgress({showProgress: false});
//       if (typeof result !== 'undefined'){
//         this.setState({
//         inClusterConfig: false,
//         k8sfile: '', 
//         k8sfileElementVal: '',
//         k8sfileError: false,
//         contextName: '', 
//         clusterConfigured: false,
//       })
//       this.props.updateK8SConfig({k8sConfig: {inClusterConfig: false, k8sfile:'', contextName:'', clusterConfigured: false}});
        
//       this.props.enqueueSnackbar('Kubernetes config was successfully removed!', {
//         variant: 'success',
//         autoHideDuration: 2000,
//         action: (key) => (
//           <IconButton
//                 key="close"
//                 aria-label="Close"
//                 color="inherit"
//                 onClick={() => self.props.closeSnackbar(key) }
//               >
//                 <CloseIcon />
//           </IconButton>
//         ),
//       });
//      }
//     }, self.handleError);
//   }

  handleDelete() {
    return false;
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
    
    let showConfigured = 'Kubernetes is not connected at the moment.';
    if (clusterConfigured) {
      showConfigured = (
        <div className={classes.alignRight}>
          <Chip 
              label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
            //   onDelete={self.handleReconfigure} 
              icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />} 
              className={classes.chip}
              variant="outlined" />
        </div>
      )
    }

    let showAdapters = 'No adapters are configured at the moment.';
    const self = this;
    if (availableAdapters.length > 0) {
      showAdapters = (
        <div>
           {
            availableAdapters.map(aa => {
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
                            image = "/static/img/istio.svg";
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
                        case 'nsm':
                            image = "/static/img/nsm.svg";
                            logoIcon = (<img src={image} className={classes.icon} />);
                            break;
                        }
                    }
                });
                

                
                return (
                <Tooltip title={isDisabled?"This adapter is inactive":`This is a ${adapterType.toUpperCase()} adapter.`}>
                  <Chip 
                  label={aa.label}
                  // onDelete={self.handleDelete(ind)} 
                  onDelete={!isDisabled?self.handleDelete:null}
                  deleteIcon={!isDisabled?<DoneIcon />:null}
                  icon={logoIcon}
                  className={classes.chip}
                  variant={isDisabled?"default":"outlined"} />
                </Tooltip>
                );
          })}
        </div>
      )

    }

    let showGrafana = 'Grafana is not connected at the moment.';
    if(grafana && grafana.grafanaURL&& grafana.grafanaURL !== ''){
      showGrafana = (
        <Chip 
        label={grafana.grafanaURL}
        // onDelete={handleGrafanaChipDelete} 
        icon={<img src="/static/img/grafana_icon.svg" className={classes.icon} />} 
        className={classes.chip}
        variant="outlined" />
      );
    }

    let showPrometheus = 'Prometheus is not connected at the moment.';
    if(prometheus && prometheus.prometheusURL&& prometheus.prometheusURL !== ''){
      showPrometheus = (
        <Chip 
          label={prometheus.prometheusURL}
          // onDelete={handlePrometheusChipDelete} 
          icon={<img src="/static/img/prometheus_logo_orange_circle.svg" className={classes.icon} />} 
          className={classes.chip}
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
