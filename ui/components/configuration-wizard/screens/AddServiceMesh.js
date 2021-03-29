import React from "react";
import {
  makeStyles,
  Container,
  Fade,
} from "@material-ui/core/";


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
        <ConfigCard name="Open Service Mesh" Icon="timer" />
        <ConfigCard name="Consul" Icon="timer" />
        <ConfigCard name="Linkerd" Icon="timer" />
      </Container>
    </Fade>
  );
};

export default External;
