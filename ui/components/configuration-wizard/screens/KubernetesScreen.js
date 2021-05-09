import React from "react";
import { withStyles, Container } from "@material-ui/core/";

import KubernetesInput from './KubernetesInput'
import KubernetesStatus from './KubernetesStatus'
import KubernetesIcon from "../icons/KubernetesIcon";
import ConfigCard from "./ConfigCard";

const styles = () => ({
  cardContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    padding: "2rem 6rem",
  },
});

class KubernetesScreen extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isChecked: false,
    };
  }
  // if kubernetes are not connected/configured do not switch
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
        <ConfigCard
          handleSwitch={this.handleSwitch}
          name="Kubernetes"
          Icon={KubernetesIcon}
          KubernetesInput={KubernetesInput}
        />
        {!this.state.isChecked ? null : (
          <KubernetesStatus/>
        )}
      </Container>
    );
  }
}

export default withStyles(styles)(KubernetesScreen);
