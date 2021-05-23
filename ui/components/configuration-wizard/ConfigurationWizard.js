import React from "react";
import { makeStyles, Container, Button, Fade } from "@material-ui/core/";

import Stepper from "./Stepper";
import KubernetesScreen from "./screens/KubernetesScreen";
import MesheryOperatorScreen from "./screens/MesheryOperator";
import ServiceMeshScreen from "./screens/AddServiceMesh";
import ExternalScreen from "./screens/External";
import ConfigurationDoneScreen from "./screens/ConfigurationDone";

const useStyles = makeStyles((theme) => ({
  container: {
    minHeight: "32.5rem",
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
    boxContent: "border-box",
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
    marginRight: theme.spacing(1),
    padding: "0.5rem 2rem",
    background: "#EFEFEF",
    color: "#607D8B",
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
  const [kubernetesConnected, setKubernetesConnected] = React.useState(false);
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
  const handleConnectToKubernetes = (checked) => setKubernetesConnected(checked);

  const handleStep = (step) => {
    switch (step) {
      case 0:
        return <KubernetesScreen handleConnectToKubernetes={handleConnectToKubernetes} />;
      case 1:
        return <MesheryOperatorScreen />;
      case 2:
        return <ServiceMeshScreen />;
      case 3:
        return <ExternalScreen />;
      default:
        return null;
    }
  };

  return (
    <Container className={classes.container}>
      <Stepper steps={steps} activeStep={activeStep} handleUserClick={handleUserClick} />
      <Fade timeout={{ enter: "1500ms" }} in="true">
        <div>
          {activeStep === steps.length ? (
            <Fade timeout={{ enter: "500ms" }} in="true">
              <ConfigurationDoneScreen handleUserClick={handleUserClick} />
            </Fade>
          ) : (
            <>
              <div>{handleStep(activeStep)}</div>
              <div className={classes.buttonContainer}>
                {activeStep === 2 || activeStep === 3 ? (
                  <Button onClick={handleNext} className={classes.skipButton}>
                    Skip
                  </Button>
                ) : null}
                {activeStep === 0 ? null : (
                  <Button
                    disabled={activeStep === 0}
                    onClick={handleBack}
                    className={(classes.button, classes.backButton)}
                  >
                    Back
                  </Button>
                )}
                <Button
                  disabled={activeStep === 0 && !kubernetesConnected}
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
