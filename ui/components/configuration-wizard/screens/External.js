import React from "react";
import { makeStyles, Container, Fade } from "@material-ui/core/";

//import grafanaIcon from "../icons/grafana.svg";
//import prometheusIcon from "../icons/prometheus.svg";
import ConfigCard from "./ConfigCard";

const useStyles = makeStyles({
  cardContainer: {
    display: "flex",
    flexDirection: "row",
    justifyContent: "flex-start",
    padding: "2rem 6rem",
  },
});

const External = () => {
  const classes = useStyles();

  return (
    <Fade timeout={{ enter: "500ms" }} in="true">
      <Container className={classes.cardContainer}>
        <ConfigCard name="Grafana" icon='{grafanaIcon}' topInputPlaceholder="URL" bottomInputPlaceholder="API Key"/>
        <ConfigCard name="Prometheus" icon='{prometheusIcon}' topInputPlaceholder="" bottomInputPlaceholder=""/>
      </Container>
    </Fade>
  );
};

export default External;
