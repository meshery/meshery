import React, { useEffect, useState } from "react";
import CustomizedSteppers from "./Stepper";
import { useStyles } from "../styles2";
import { helmStepContent, helmStepIcons, helmSteps } from "../connections/helm-connect/constants";
import { Paper, Box } from "@material-ui/core"

function ConnectionStepper({ connectionType }) {
  const classes = useStyles();
  const [stepData, setStepData] = useState({
    stepContent: {},
    stepIcons: {},
    steps: []
  });

  useEffect(() => {
    switch (connectionType) {
      case "helm":
        setStepData({
          stepContent: helmStepContent,
          stepIcons: helmStepIcons,
          steps: helmSteps
        });
        break;
      default:
    }
  }, [connectionType]);

  return (
    <Box className={classes.welcomeContainer}>
      <Paper className={classes.connectionStepperContent} elevation={3}>
        {stepData.steps.length > 0 && <CustomizedSteppers stepData={stepData} />}
      </Paper>
    </Box>
  );
}

export default React.memo(ConnectionStepper);
