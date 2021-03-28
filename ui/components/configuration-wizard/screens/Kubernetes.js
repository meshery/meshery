import React from "react";
import { makeStyles, Container } from "@material-ui/core/";

//import kubernetesIcon from "../icons/Kubernetes.svg";
import ConfigCard from "./ConfigCard";

const useStyles = makeStyles({
  cardContainer: {
    display: "flex",
    flexDirection: "row",
    justifyContent: "flex-start",
    padding: "2rem 6rem",
  },
});

const Kubernetes = () => {
  const classes = useStyles();

  return (
    <Container className={classes.cardContainer}>
      {" "}
      <ConfigCard name="Kubernetes" icon='{kubernetesIcon}' topInputPlaceholder="Upload" bottomInputPlaceholder="Context"/>
    </Container>
  );
};

export default Kubernetes;
