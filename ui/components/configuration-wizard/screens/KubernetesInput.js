import React from "react";
import { connect } from "react-redux";
import { withStyles, FormGroup, TextField, InputAdornment } from "@material-ui/core/";
import BackupIcon from "@material-ui/icons/Backup";

const styles = (theme) => ({
  fileInputStyle: {
    opacity: "0.01",
  },
  contentBottomChecked: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "flex-end",
  },
  contentBottomUpperInput: {
     width: "11rem",
     fontSize: "0.75rem",
     marginLeft: "2.4rem",
     marginTop: "-1rem",
     marginBottom: "0rem",
   },
  contentBottomLowerInput: {
    width: "11rem",
    marginBottom: "-1rem",
    fontSize: "0.75rem",
    marginTop: "0",
  },
});

class KubernetesInput extends React.Component {
  constructor(props) {
    super(props);
    const { k8sconfig } = this.props;
    this.state = {
      inClusterConfig: k8sconfig.inClusterConfig, // read from store
      k8sfile: k8sconfig.k8sfile, // read from store
      isClusterConfigured: k8sconfig.clusterConfigured,
      contextName: k8sconfig.contextName, // read from store
      contextNameForForm: "",
      k8sfileElementVal: "",
      k8sfileError: false,
      ts: new Date(),
    };
  }
  static getDerivedStateFromProps(props, state) {
    const {
      inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer,
    } = props;
    if (props.ts > state.ts) {
      return {
        inClusterConfig,
        k8sfile,
        k8sfileElementVal: '',
        contextName,
        clusterConfigured,
        configuredServer,
        ts: props.ts,
      };
    }
    return {};
  }
  handleChange = (name) => {
    const self = this;
    return (event) => {
      if (name === "inClusterConfigForm") {
        self.setState({ [name]: event.target.checked, ts: new Date() });
        return;
      }
      if (name === "k8sfile") {
        if (event.target.value !== "") {
          self.setState({ k8sfileError: false });
        }
        self.setState({ k8sfileElementVal: event.target.value });
        self.fetchContexts();
      }
      self.setState({ [name]: event.target.value, ts: new Date() });
      this.handleSubmit();
    };
  };
  render() {
    const { classes } = this.props;
    return (
      <div className={classes.contentBottomChecked}>
        <FormGroup>
          <input
            id="k8sfile"
            type="file"
            value={this.state.k8sfileElementVal}
            onChange={this.handleChange("k8sfile")}
            className={classes.fileInputStyle}
          />
          <TextField
            id="k8sfileLabelText"
            name="k8sfileLabelText"
            className={classes.contentBottomUpperInput}
            label="Upload kubeconfig"
            variant="outlined"
            fullWidth
            value={this.state.k8sfile.replace("C:\\fakepath\\", "")}
            onClick={() => document.querySelector("#k8sfile").click()}
            margin="normal"
            InputProps={{
              readOnly: true,
              endAdornment: (
                <InputAdornment position="end">
                  <BackupIcon style={{marginRight: "-.8rem"}}/>
                </InputAdornment>
              ),
            }}
            disabled
          />
        </FormGroup>
        <TextField
          select
          id="contextName"
          name="contextName"
          label="Context Name"
          fullWidth
          value={this.state.contextNameForForm}
          margin="normal"
          variant="outlined"
          // disabled={inClusterConfigForm === true}
          onChange={this.handleChange("contextNameForForm")}
          className={classes.contentBottomLowerInput}
        >
          {this.state.contextsFromFile &&
            this.state.contextsFromFile.map((ct) => (
              <MenuItem key={`ct_---_${ct.contextName}`} value={ct.contextName}>
                {ct.contextName}
                {ct.currentContext ? " (default)" : ""}
              </MenuItem>
            ))}
        </TextField>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig").toJS();
  return {
    k8sconfig,
  };
};

export default withStyles(styles)(connect(mapStateToProps)(KubernetesInput));
