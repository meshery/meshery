import React from "react";
import { makeStyles, Container, Fade, Typography } from "@material-ui/core/";

import MesheryOperatorIcon from "../icons/MesheryOperatorIcon";
import ConfigCard from "./ConfigCard";
//import MesheryOperatorStatus from './MesheryOperatorStatus'

const useStyles = makeStyles({
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
    flexDirection: "row",
    justifyContent: "space-between",
    width: "20rem",
    height: "12.5rem",
    padding: "2rem 2rem",
    marginRight: "3rem",
    boxShadow: "0px 1px 6px 1px rgba(0,0,0,0.75)",
  },
  infoContext: {
    fontSize: ".85rem",
  },
  infoKind: {
    fontSize: ".75rem",
    marginTop: ".5rem",
    color: "#CACACA",
  },
});

const MesheryOperator = () => {
  const classes = useStyles();
  const [state, setState] = React.useState(false);

  const handleSwitch = (name, checked) => {
    setState(checked);
  };

  return (
    <Fade timeout={{ enter: "500ms" }} in="true">
      <Container className={classes.cardContainer}>
        {" "}
        <ConfigCard
          handleSwitch={handleSwitch}
          name="Meshery Operator"
          Icon={MesheryOperatorIcon}
          topInputPlaceholder=""
          bottomInputPlaceholder=""
        />
        {!state ? null : ''}
      </Container>
    </Fade>
  );
};

export default MesheryOperator;
