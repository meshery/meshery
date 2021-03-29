import React from "react";
import { makeStyles, Container, Fade } from "@material-ui/core/";

import GrafanaIcon from "../icons/GrafanaIcon.js";
import PrometheusIcon from "../icons/PrometheusIcon.js";
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
        <ConfigCard name="Grafana" Icon={GrafanaIcon} topInputPlaceholder="URL" bottomInputPlaceholder="API Key"/>
        <ConfigCard name="Prometheus" Icon={PrometheusIcon} topInputPlaceholder="" bottomInputPlaceholder=""/>
      </Container>
    </Fade>
  );
};

export default External;
