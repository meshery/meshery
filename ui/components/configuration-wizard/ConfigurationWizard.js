import React from "react";
import { makeStyles, Container, Button, Fade } from "@material-ui/core/";

import Stepper from "./Stepper";
import Kubernetes from "./screens/Kubernetes";
import MesheryOperator from "./screens/MesheryOperator";
import AddServiceMesh from "./screens/AddServiceMesh";
import External from "./screens/External";
import ConfigurationDone from "./screens/ConfigurationDone";

const useStyles = makeStyles((theme) => ({
  container: {
    height: "32.5rem",
    margin: "5rem auto",
    background: "white",
    boxShadow: "lightgrey 0px 0px 10px",
    display: "flex",
    flexDirection: "column",
    justifyContent: "space-between",
  },
  buttonContainer: {
    textAlign: "right",
    paddingBottom: "5rem",
    marginRight: "11.5rem",
  },
  button: {
    marginRight: theme.spacing(1),
    padding: "0.5rem 2rem",
    textDecoration: "none",
    background: "white",
    color: "#647881",
    border: "1.5px solid #647881",
    "&:hover": {
      backgroundColor: "#647881",
      color: "white",
    },
  },
  backButton: {
    background: "white",
    color: "lightgray",
  },
  skipButton: {
    color: "#647881",
  },
}));

function getSteps() {
  return ["Kubernetes", "Meshery Operator", "Service Mesh", "External"];
}

const ConfigurationWizard = () => {
  const classes = useStyles();
  const [activeStep, setActiveStep] = React.useState(0);
  const steps = getSteps();
  
  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };
  
  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };
  const handleUserClick = (navStep) => {
    setActiveStep(navStep);
  };
  const handleStep = (step) => {
    switch (step) {
      case 0:
        return <Kubernetes />;
      case 1:
        return <MesheryOperator />;
      case 2:
        return <AddServiceMesh />;
      case 3:
        return <External />;
      default:
        return null;
    }
  };
  
  return (
    <Container className={classes.container}>
      <Stepper
        steps={steps}
        activeStep={activeStep}
        handleUserClick={handleUserClick}
      />
      <Fade timeout={{ enter: "1500ms" }} in="true">
        <div>
          {activeStep === steps.length ? (
            <Fade timeout={{ enter: "500ms" }} in="true">
              <ConfigurationDone />
            </Fade>
          ) : (
            <>
              <div>{handleStep(activeStep)}</div>
              <div className={classes.buttonContainer}>
                <Button
                  disabled={activeStep === 0}
                  onClick={handleBack}
                  className={(classes.button, classes.backButton)}
                >
                    Back
                </Button>
                {activeStep === 1 || activeStep === 2 ? (
                  <Button onClick={handleNext} className={classes.skipButton}>
                      Skip
                  </Button>
                ) : null}
                <Button
                  variant="contained"
                  onClick={handleNext}
                  className={classes.button}
                >
                  {activeStep === steps.length - 1 ? "Finish" : "Next"}
                </Button>
              </div>
            </>
          )}
        </div>
      </Fade>
    </Container>
  );
};
  
export default ConfigurationWizard;