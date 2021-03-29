import React from "react";
import {
  makeStyles,
  Container,
  Fade,
} from "@material-ui/core/";

import MesheryOperatorIcon from "../icons/MesheryOperatorIcon";
import ConfigCard from "./ConfigCard";

const useStyles = makeStyles({
  cardContainer: {
    display: "flex",
    flexDirection: "row",
    justifyContent: "flex-start",
    padding: "2rem 6rem",
  },
});

const MesheryOperator = () => {
  const classes = useStyles();

  return (
    <Fade timeout={{ enter: "500ms" }} in="true">
      <Container className={classes.cardContainer}>
        {" "}
        <ConfigCard name="Meshery Operator" Icon={MesheryOperatorIcon} topInputPlaceholder="" bottomInputPlaceholder=""/>
      </Container>
    </Fade>
  );
};

export default MesheryOperator;
