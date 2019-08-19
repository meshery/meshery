import {connect} from "react-redux";
import NoSsr from '@material-ui/core/NoSsr';
import { withStyles, Typography, Button, Divider, ExpansionPanelDetails } from '@material-ui/core';
import { blue } from '@material-ui/core/colors';
import PropTypes from 'prop-types';
import { withRouter } from 'next/router';
import MesheryAdapterPlayComponent from "./MesheryAdapterPlayComponent";
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import { ExpansionPanel, ExpansionPanelSummary } from "./ExpansionPanels";

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
    // flexBasis: '33.33%',
  },
  expTitleIcon: {
    width: theme.spacing(3),
    display: 'inline',
    verticalAlign: 'middle',
  },
  expIstioTitleIcon: {
    width: theme.spacing(2),
    display: 'inline',
    verticalAlign: 'middle',
    marginLeft: theme.spacing(0.5),
    marginRight: theme.spacing(0.5),
  },
  expTitle: {
    display: 'inline',
    verticalAlign: 'middle',
  }
});

class MesheryPlayComponent extends React.Component {
  constructor(props){
    super(props);
    
    const {k8sconfig, meshAdapters} = props;

    this.state = {
      k8sconfig,
      meshAdapters,
    }
  }

  handleConfigure = () => {
    this.props.router.push('/settings');
  }

  renderIndividualAdapter() {
    const { meshAdapters } = this.props;
    let adapter;
    meshAdapters.forEach(adap => {
      if (adap.adapter_location === this.props.adapter){
        adapter = adap;
      }
    });
    if(adapter){
      return (
        <MesheryAdapterPlayComponent adapter={adapter} />
      )
    }
    return '';
  }

  render() {
    const {classes, color, iconButtonClassName, avatarClassName, k8sconfig, meshAdapters} = this.props;

    if (k8sconfig.clusterConfigured === false || meshAdapters.length === 0) {
      return (
        <NoSsr>
        <React.Fragment>
          <div className={classes.alreadyConfigured}>
            {/* <Typography variant="subtitle1" gutterBottom>
            Configure service meshes
            </Typography> */}

            <Button variant="contained" color="primary" size="large" onClick={this.handleConfigure}>
              Configure playground
            </Button>
          </div>
          </React.Fragment>
          </NoSsr>
      );
    }
    if(this.props.adapter && this.props.adapter !== '') {
      const indContent = this.renderIndividualAdapter();
      if(indContent !== ''){
        return indContent;
      } // else it will render all the available adapters
    }

    var self = this;
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.root}>
            {meshAdapters.map((adapter, ind) => {
                let image = "/static/img/meshery-logo.png";
                let imageIcon = (<img src={image} className={classes.expTitleIcon} />);
                switch (adapter.name.toLowerCase()){
                  case 'istio':
                    image = "/static/img/istio.svg";
                    imageIcon = (<img src={image} className={classes.expIstioTitleIcon} />);
                    break;
                  case 'linkerd':
                    image = "/static/img/linkerd.svg";
                    imageIcon = (<img src={image} className={classes.expTitleIcon} />);
                    break;
                  case 'consul':
                    image = "/static/img/consul.svg";
                    imageIcon = (<img src={image} className={classes.expTitleIcon} />);
                    break;
                  case 'nsm':
                    image = "/static/img/nsm.svg";
                    imageIcon = (<img src={image} className={classes.expTitleIcon} />);
                    break;
                  // default:
                } 
                return (
                <ExpansionPanel key={`mplay_exp_${ind}`} square defaultExpanded={ind === 0?true:false}>
                  <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                    <div className={classes.column}>
                      <Typography variant="h6" gutterBottom>
                        <div className={classes.column}>
                          {imageIcon}{' '}
                          <span className={classes.expTitle}>{adapter.adapter_location}</span>
                        </div>
                      </Typography>
                    </div>
                  </ExpansionPanelSummary>
                  <ExpansionPanelDetails>
                      <Divider variant="fullWidth" />
                      <MesheryAdapterPlayComponent adapter={adapter} />
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
