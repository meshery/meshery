import React from "react";
import { connect } from "react-redux";
import { withStyles, Container, Typography } from "@material-ui/core/";
import BackupIcon from "@material-ui/icons/Backup";

import KubernetesIcon from "../icons/KubernetesIcon";
import ConfigCard from "./ConfigCard";

const styles = (theme) => ({
  cardContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    padding: "2rem 6rem",
  },
  infoContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "column",
    width: "20rem",
    padding: "5rem 2rem",
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

class Kubernetes extends React.Component {
  constructor(props) {
    super(props);
    const { k8sconfig } = this.props;
    this.state = {
      isChecked: false,
      inClusterConfig: k8sconfig.inClusterConfig, // read from store
      k8sfile: k8sconfig.k8sfile, // read from store
      isClusterConfigured: k8sconfig.clusterConfigured,
      contextName: k8sconfig.contextName, // read from store
    };
  }

  handleSwitch = (name, checked) => {
    const { handleConnectToKubernetes } = this.props;
    this.setState({ isChecked: checked });
    if (handleConnectToKubernetes) {
      handleConnectToKubernetes(checked);
    }
  };
  render() {
    
    const { classes } = this.props;
    return (
      <Container className={classes.cardContainer}>
        {" "}
        <ConfigCard
          handleSwitch={this.handleSwitch}
          name="Kubernetes"
          Icon={KubernetesIcon}
          topInputPlaceholder="Upload Kubeconfig"
          TopInputIcon={BackupIcon}
          bottomInputPlaceholder="Current-Context"
        />
        {!this.state.isChecked ? null : (
          <div className={classes.infoContainer}>
            <Typography className={classes.infoStatus}>Status</Typography>
            <Typography className={classes.infoContext}>Current-Context: bob-us-east</Typography>
            <Typography className={classes.infoContext}>Clust</Typography>
            <Typography className={classes.infoContext}>Current Context: </Typography>
          </div>
        )}
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

export default withStyles(styles)(connect(mapStateToProps)(Kubernetes));
