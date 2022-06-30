import React from "react";
import { connect } from "react-redux";
import NoSsr from "@material-ui/core/NoSsr";
import {
  withStyles, Button, Divider, MenuItem, TextField, Grid
} from "@material-ui/core";
import { blue } from "@material-ui/core/colors";
import PropTypes from "prop-types";
import { withRouter } from "next/router";
import SettingsIcon from "@material-ui/icons/Settings";
import MesheryAdapterPlayComponent from "./MesheryAdapterPlayComponent";
import { bindActionCreators } from "redux";
import { setAdapter } from "../lib/store";

const styles = (theme) => ({
  icon : { fontSize : 20, },
  playRoot : { padding : theme.spacing(0),
    marginBottom : theme.spacing(2), },
  buttons : { display : "flex",
    justifyContent : "flex-end", },
  button : { marginTop : theme.spacing(3),
    marginLeft : theme.spacing(1), },
  margin : { margin : theme.spacing(1), },
  alreadyConfigured : { textAlign : "center",
    padding : theme.spacing(20), },
  colorSwitchBase : { color : blue[300],
    "&$colorChecked" : { color : blue[500],
      "& + $colorBar" : { backgroundColor : blue[500], }, }, },
  colorBar : {},
  colorChecked : {},
  uploadButton : { margin : theme.spacing(1),
    marginTop : theme.spacing(3), },
  fileLabel : { width : "100%", },
  editorContainer : { width : "100%", },
  deleteLabel : { paddingRight : theme.spacing(2), },
  alignRight : { textAlign : "right", },
  expTitleIcon : { width : theme.spacing(3),
    display : "inline",
    verticalAlign : "middle", },
  expIstioTitleIcon : {
    width : theme.spacing(2),
    display : "inline",
    verticalAlign : "middle",
    marginLeft : theme.spacing(0.5),
    marginRight : theme.spacing(0.5),
  },
  expTitle : { display : "inline",
    verticalAlign : "middle",
    marginLeft : theme.spacing(1), },
  paneSection : { backgroundColor : "#fff",
    padding : theme.spacing(2.5),
    borderRadius : 4, },
});

class MesheryPlayComponent extends React.Component {
  constructor(props) {
    super(props);

    const {  meshAdapters } = props;
    let adapter = {};
    if (meshAdapters && meshAdapters.size > 0) {
      adapter = meshAdapters[0];
    }
    this.state = {
      adapter,
    };
  }

  handleRouteChange =  () => {
    const queryParam = this.props?.router?.query?.adapter;
    if (queryParam) {
      const selectedAdapter = this.props.meshAdapters.find(({ adapter_location }) => adapter_location === queryParam);
      if (selectedAdapter) {
        this.setState({ adapter : selectedAdapter })
      }
    } else if (this.props.meshAdapters.size > 0) {
      this.setState({ adapter : this.props.meshAdapters.get(0) })
    }
  }

  componentDidMount() {
    const { router } = this.props;
    router.events.on('routeChangeComplete', this.handleRouteChange)
  }

  componentDidUpdate(prevProps) {
    // update the adapter when the meshadapters props are changed
    if (prevProps.meshAdapters?.size !== this.props.meshAdapters?.size
      && this.props.meshAdapters.size > 0
    ) {
      this.handleRouteChange();
    }
  }

  componentWillUnmount() {
    this.props.router.events.off('routeChangeComplete', this.handleRouteChange)
  }

  handleConfigure = () => {
    this.props.router.push("/settings#service-mesh");
  }

  pickImage(adapter) {
    const { classes } = this.props;
    let image = "/static/img/meshery-logo.png";
    let imageIcon = <img src={image} className={classes.expTitleIcon} />;
    if (adapter && adapter.name) {
      image = "/static/img/" + adapter.name.toLowerCase() + ".svg";
      imageIcon = <img src={image} className={classes.expTitleIcon} />;
    }
    return imageIcon;
  }

  handleAdapterChange = () => {
    const self = this;
    return (event) => {
      const { setAdapter,meshAdapters } = self.props;
      if (event.target.value !== "") {
        const selectedAdapter = meshAdapters.filter(({ adapter_location }) => adapter_location === event.target.value);
        if (selectedAdapter && selectedAdapter.size === 1) {
          self.setState({ adapter : selectedAdapter.get(0) });
          setAdapter({ selectedAdapter : selectedAdapter.get(0).name });
        }
      }
    };
  };

  renderIndividualAdapter() {
    const { meshAdapters } = this.props;
    let adapCount = 0;
    let adapter;
    meshAdapters.forEach((adap) => {
      if (adap.adapter_location === this.props.adapter) {
        adapter = adap;
        meshAdapters.forEach((ad) => {
          if (ad.name == adap.name) adapCount += 1;
        });
      }
    });
    if (adapter) {
      const imageIcon = this.pickImage(adapter);
      return (
        <React.Fragment>
          <MesheryAdapterPlayComponent adapter={adapter} adapCount={adapCount} adapter_icon={imageIcon} />
        </React.Fragment>
      );
    }
    return "";
  }

  render() {
    const { classes,  meshAdapters } = this.props;
    let { adapter } = this.state;

    if (meshAdapters.size === 0) {
      return (
        <NoSsr>
          <React.Fragment>
            <div className={classes.alreadyConfigured}>
              <Button variant="contained" color="primary" size="large" onClick={this.handleConfigure}>
                <SettingsIcon className={classes.icon} />
                Configure Settings
              </Button>
            </div>
          </React.Fragment>
        </NoSsr>
      );
    }

    if (this.props.adapter && this.props.adapter !== "") {
      const indContent = this.renderIndividualAdapter();
      if (indContent !== "") {
        return indContent;
      } // else it will render all the available adapters
    }

    const self = this;
    const imageIcon = self.pickImage(adapter);
    return (
      <NoSsr>
        <React.Fragment>
          <div className={classes.playRoot}>
            <Grid container>
              <Grid item xs={12} className={classes.paneSection}>
                <TextField
                  select
                  id="adapter_id"
                  name="adapter_name"
                  label="Select Service Mesh Type"
                  fullWidth
                  value={adapter && adapter.adapter_location
                    ? adapter.adapter_location
                    : ""}
                  margin="normal"
                  variant="outlined"
                  onChange={this.handleAdapterChange()}
                >
                  {meshAdapters.map((ada) => (
                    <MenuItem key={`${ada.adapter_location}_${new Date().getTime()}`} value={ada.adapter_location}>
                      {/* <ListItemIcon> */}
                      {self.pickImage(ada)}
                      {/* </ListItemIcon> */}
                      <span className={classes.expTitle}>{ada.adapter_location}</span>
                    </MenuItem>
                  ))}
                </TextField>
              </Grid>
            </Grid>
          </div>
          <Divider variant="fullWidth" light />
          {adapter && adapter.adapter_location && (
            <MesheryAdapterPlayComponent adapter={adapter} adapter_icon={imageIcon} />
          )}
        </React.Fragment>
      </NoSsr>
    );
  }
}

MesheryPlayComponent.propTypes = { classes : PropTypes.object.isRequired, };

const mapDispatchToProps = (dispatch) => ({
  setAdapter : bindActionCreators(setAdapter, dispatch)
});

const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig");
  const meshAdapters = state.get("meshAdapters");
  const meshAdaptersts = state.get("meshAdaptersts");
  const selectedAdapter = state.get("selectedAdapter");
  return { k8sconfig, meshAdapters, meshAdaptersts, selectedAdapter };
};

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withRouter(MesheryPlayComponent)));
