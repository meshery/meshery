import React from "react";
import { connect } from "react-redux";
import {
  withStyles,
  FormGroup,
  TextField,
  InputAdornment,
  MenuItem,
  Container,
  Switch,
  FormControlLabel,
  Input,
  Card,
  CardContent,
  Typography,
} from "@material-ui/core/";
import BackupIcon from "@material-ui/icons/Backup";

import KubernetesInput from "./KubernetesInput";
import KubernetesStatus from "./KubernetesStatus";
import KubernetesIcon from "../icons/KubernetesIcon";
import ConfigCard from "./ConfigCard";

const MeshySwitch = withStyles({
  switchBase: {
    color: "grey",
    "&$checked": {
      color: "#00B39F",
    },
    "&$checked + $track": {
      backgroundColor: "#00B39F",
    },
  },
  checked: {},
  track: {},
})(Switch);

const styles = () => ({
  // Container
  cardContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    padding: "2rem 6rem",
  },
  // Card
  card: {
    position: "relative",
    width: "12rem",
    minWidth: "10rem",
    border: "1px solid gray",
    borderRadius: "0.75rem",
    //top: "2rem",
    margin: "0rem 0rem 6rem 0rem",
    ["@media (max-width:1024px)"]: {
      //eslint-disable-line no-useless-computed-key
      margin: "0rem 0rem 6rem 0",
    },
  },
  cardChecked: {
    height: "15rem",
    marginBottom: "1rem",
  },
  cardUnchecked: {
    height: "10rem",
  },
  cardContent: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    alignItems: "center",
    height: "100%",
    width: "100%",
    padding: "0",
  },
  contentTop: {
    background: "#434343",
    height: "12rem",
    width: "100%",
    display: "flex",
    alignItems: "center",
  },
  contentTopUnchecked: {
    background: "#434343",
    height: "100%",
    width: "100%",
    display: "flex",
    alignItems: "center",
  },
  contentTopSwitcher: {
    paddingLeft: "2rem",
  },
  iconContainer: {
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
    marginLeft: "1rem",
  },
  cardIcon: {
    width: "3rem",
  },
  cardIconText: {
    color: "white",
    fontSize: "0.85rem",
    textAlign: "center",
    "&:first-letter": {
      textTransform: "capitalize",
    },
  },
  contentBottomInputChecked: {
    background: "white",
    height: "6rem",
    width: "100%",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    justifyContent: "center",
  },
  contentBottomInputUnchecked: {
    display: "none",
  },
  contentBottomInput: {
    border: "1px solid lightgray",
    borderRadius: "5px",
    width: "9rem",
    height: "2rem",
    marginBottom: "0.15rem",
    fontSize: "0.75rem",
    padding: "0.50rem",
  },
  topInputIcon: {
    position: "absolute",
    fontSize: "1.25rem",
    color: "lightgray",
    bottom: "4.25rem",
    left: "9rem",
    cursor: "pointer",
    zIndex: "99999",
    "&:hover": {
      color: "grey",
    },
  },
  file: {
    display: "none",
  },
  // Inputs
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

class KubernetesScreen extends React.Component {
  constructor(props) {
    super(props);
    const { k8sconfig, handleConnectToKubernetes } = this.props;
    this.state = {
      inClusterConfig: k8sconfig.inClusterConfig, // read from store
      k8sfile: k8sconfig.k8sfile, // read from store
      isClusterConfigured: k8sconfig.clusterConfigured,
      contextName: k8sconfig.contextName, // read from store
      contextNameForForm: "",
      k8sfileElementVal: "",
      k8sfileError: false,
      ts: new Date(),
      isChecked: this.props.k8sconfig ? true : false,
    };
  }
  static getDerivedStateFromProps(props, state) {
    const { inClusterConfig, contextName, clusterConfigured, k8sfile, configuredServer } = props;
    if (props.ts > state.ts) {
      return {
        inClusterConfig,
        k8sfile,
        k8sfileElementVal: "",
        contextName,
        clusterConfigured,
        configuredServer,
        ts: props.ts,
      };
    }
    return {};
  }
  handleSwitch = (name, checked) => {
    if (this.props.k8sconfig) {
      this.setState({ isChecked: checked });
      if (this.props.handleConnectToKubernetes) {
        this.props.handleConnectToKubernetes(checked);
      }
    }
  };
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
      <Container className={classes.cardContainer}>
        <Card className={`${classes.card} ${classes.cardChecked}`} variant="outlined">
          <CardContent className={classes.cardContent}>
            <div className={classes.contentTop}>
              <div className={classes.iconContainer}>
                <KubernetesIcon className={classes.cardIcon} alt="Kubernetes icon" />
                <Typography className={classes.cardIconText} color="primary">
                  Kubernetes
                </Typography>
              </div>
              <FormControlLabel
                className={classes.contentTopSwitcher}
                control={<MeshySwitch checked={this.state.isChecked} name="Kubernetes" />}
                onChange={this.handleSwitch}
              />
            </div>
            <div className={classes.contentBottomInputChecked}>
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
                        <BackupIcon style={{ marginRight: "-.8rem" }} />
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
          </CardContent>
        </Card>
        {this.props.k8sconfig ? <KubernetesStatus isChecked={this.state.isChecked}/> : null}
      </Container>
    );
  }
}

const mapStateToProps = (state) => {
  const k8sconfig = state.get("k8sConfig").toJS();
  return {
    k8sconfig,
  };
};

export default withStyles(styles)(connect(mapStateToProps)(KubernetesScreen));
