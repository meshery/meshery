import React from "react";
import { connect } from "react-redux";
import { withStyles, Typography, Chip } from "@material-ui/core/";

const styles = (theme) => ({
  infoContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "column",
    width: "20rem",
    padding: "5rem 1rem",
    boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.75)",
  },
  infoStatus: {
    position: "absolute",
    bottom: "12.50rem",
    right: "10rem",
    color: "#647881",
    background: "#F1F3F4",
    padding: ".5rem 5rem .75rem 1.5rem",
    borderRadius: "0.25rem",
    fontSize: ".8rem",
  },
  infoContext: {
    fontSize: ".9rem",
  },
  infoKind: {
    fontSize: ".75rem",
    color: "#CACACA",
  },
});

class KubernetesStatus extends React.Component {
  constructor(props) {
    super(props);
    const {
        inClusterConfig, contextName, clusterConfigured,
      } = props;
    this.state = {
      inClusterConfig, // read from store
      contextName, // read from store
      clusterConfigured, // read from store
      ts: new Date(),
    };
    let showConfigured = "";
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

  render() {
    const { classes } = this.props;
    return (
      <div className={classes.infoContainer}>
        {this.state.clusterConfigured ? (
          <>
            <Typography className={classes.infoStatus}>Status</Typography>
            <Typography className={classes.infoContext}>
              Current-Context:{this.state.inClusterConfig ? "Using In Cluster Config" : this.state.contextName}
            </Typography>
            <Typography className={classes.infoContext}>Cluster: {this.state.inClusterConfig ? "Using In Cluster Config" : "Using Out Of Cluster Config"}</Typography>
          </>
        ) : null}
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

export default withStyles(styles)(connect(mapStateToProps)(KubernetesStatus));
