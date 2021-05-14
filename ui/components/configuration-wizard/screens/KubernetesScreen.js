import React from "react";
import { makeStyles, Container, Typography } from "@material-ui/core/";

import KubernetesInput from "./KubernetesInput";
import KubernetesStatus from "./KubernetesStatus";
import KubernetesIcon from "../icons/KubernetesIcon";
import ConfigCard from "./ConfigCard";

const useStyles = makeStyles({
  cardContainer: {
    position: "relative",
    display: "flex",
    flexDirection: "row",
    justifyContent: "space-between",
    padding: "2rem 6rem",
  },
});

const KubernetesScreen = ({ handleConnectToKubernetes }) => {
  const classes = useStyles();
  const [isChecked, setIsChecked] = React.useState(false);

  const handleSwitch = (name, checked) => {
    setIsChecked(checked)
    if (handleConnectToKubernetes) {
      handleConnectToKubernetes(checked);
    }
  };

  return (
    <Container className={classes.cardContainer}>
      <ConfigCard
        handleSwitch={handleSwitch}
        name="Kubernetes"
        Icon={KubernetesIcon}
        KubernetesInput={KubernetesInput}
      />
      {isChecked ? <KubernetesStatus /> : null}
      
    </Container>
  );
};

export default KubernetesScreen;
