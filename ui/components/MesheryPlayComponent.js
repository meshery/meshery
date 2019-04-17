import {connect} from "react-redux";
import NoSsr from '@material-ui/core/NoSsr';
import { withStyles, ExpansionPanel, ExpansionPanelSummary, Typography, ExpansionPanelDetails, Button, Divider } from '@material-ui/core';
import { blue } from '@material-ui/core/colors';
// import MesherySnackbarWrapper from './MesherySnackbarWrapper';
import PropTypes from 'prop-types';
import { withRouter } from 'next/router';
import { updateK8SConfig } from "../lib/store";
import { bindActionCreators } from "redux";
import MesheryAdapterPlayComponent from "./MesheryAdapterPlayComponent";
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

const styles = theme => ({
  root: {
    padding: theme.spacing(10),
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
  uploadButton: {
    margin: theme.spacing(1),
    marginTop: theme.spacing(3),
  },
  rightIcon: {
    // marginLeft: theme.spacing(1),
  },
  fileLabel: {
    width: '100%',
  },
  fileLabelText: {
    // width: '79%',
  },
  editorContainer: {
    width: '100%',
  },
  deleteLabel: {
    paddingRight: theme.spacing(2),
  },
  alignRight: {
    textAlign: 'right',
  },
  column: {
    flexBasis: '33.33%',
  },
  icon: {
    width: theme.spacing(3),
  },
  expTitle: {
    marginTop: -theme.spacing(1),
  }
});

class MesheryPlayComponent extends React.Component {
  constructor(props){
    super(props);
    
    const {k8sconfig, meshAdapters} = props;
    // if (k8sconfig.clusterConfigured === false || meshAdapters.length === 0) {
    //   props.router.push('/configure'); // TODO: need to figure out a better way to do this
    // }

    this.state = {
      k8sconfig,
      meshAdapters,
    }
  }

  handleConfigure = () => {
    // const { k8sConfig } = this.props;
    // k8sConfig.reconfigureCluster = true;
    // this.props.updateK8SConfig({k8sConfig});
    this.props.router.push('/configure');
  }

  render() {
    const {classes, color, iconButtonClassName, avatarClassName, k8sconfig, meshAdapters, ...other} = this.props;

    if (k8sconfig.clusterConfigured === false || meshAdapters.length === 0) {
      // props.router.push('/configure'); // TODO: need to figure out a better way to do this
      return (
        <NoSsr>
        <React.Fragment>
          <div className={classes.alreadyConfigured}>
            <Typography variant="subtitle1" gutterBottom>
            Connect to an existing service mesh adapter to use the playground
            </Typography>

            <Button variant="contained" color="primary" size="large" onClick={this.handleConfigure}>
              Here
              </Button>
          </div>
          </React.Fragment>
          </NoSsr>
      );
    }

    var self = this;
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.root}>
            {meshAdapters.map((adapter, ind) => {
                let image = "/static/img/meshery-logo.png";
                switch (adapter.name.toLowerCase()){
                  case 'istio':
                    image = "/static/img/istio.svg";
                    break;
                  case 'linkerd':
                    image = "/static/img/linkerd.svg";
                    break;
                  // default:
                } 
                return (
                <ExpansionPanel defaultExpanded={ind === 0?true:false}>
                  <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                    <div className={classes.column}>
                    <Typography variant="subtitle1" gutterBottom><img src={image} className={classes.icon} /> <span className={classes.expTitle}>{adapter.adapter_location}</span></Typography>
                    </div>
                    <Divider light variant="fullWidth" />
                  </ExpansionPanelSummary>
                  <ExpansionPanelDetails>
                      <MesheryAdapterPlayComponent index={ind} adapter={adapter} />
                  </ExpansionPanelDetails>
                </ExpansionPanel>
              );
              })}
          </div>
        </React.Fragment>
      </NoSsr>
    )
  }
}

MesheryPlayComponent.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapDispatchToProps = dispatch => {
  return {
      // updateK8SConfig: bindActionCreators(updateK8SConfig, dispatch),
  }
}

const mapStateToProps = state => {
    const k8sconfig = state.get("k8sConfig").toJS();
    const meshAdapters = state.get("meshAdapters").toJS();
    return {k8sconfig, meshAdapters};
}

export default withStyles(styles)(connect(
    mapStateToProps,
    mapDispatchToProps
  )(withRouter(MesheryPlayComponent)));
