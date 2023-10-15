/* eslint-disable no-unused-vars */
import React from 'react';
import { makeStyles, Container, Button, Fade, Grid } from '@material-ui/core';
import { useRouter } from 'next/router';

import Stepper from './Stepper.js';
import KubernetesScreen from './Screens/KubernetesScreen';
import MesheryOperatorScreen from './Screens/MesheryOperatorScreen';
import ServiceMeshScreen from './Screens/ServiceMeshScreen';
import MetricsScreen from './Screens/MetricsScreen';
import CompletionScreen from './Screens/CompletionScreen';

const useStyles = makeStyles((theme) => ({
  container: {
    minHeight: '32.5rem',
    margin: '5rem auto',
    padding: '2rem',
    paddingRight: '4rem',
    paddingLeft: '4rem',
    background: 'white',
    boxShadow: 'lightgrey 0px 0px 10px',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
  },
  buttonContainer: { textAlign: 'right', paddingRight: '3rem', marginTop: '2rem' },
  button: {
    boxContent: 'border-box',
    marginRight: theme.spacing(1),
    padding: '0.5rem 2rem',
    textDecoration: 'none',
    background: 'white',
    color: '#647881',
    border: '1.5px solid #647881',
    '&:hover': { backgroundColor: '#647881', color: 'white' },
  },
  backButton: {
    marginRight: theme.spacing(1),
    padding: '0.5rem 2rem',
    background: '#EFEFEF',
    color: '#607D8B',
  },
  skipButton: { color: '#647881', '&:hover': { backgroundColor: 'white' } },
}));

function getSteps() {
  return ['Kubernetes', 'Meshery Operator', 'Meshery Adapters', 'Metrics'];
}

const ConfigurationWizard = () => {
  const classes = useStyles();
  const [activeStep, setActiveStep] = React.useState(0);
  const router = useRouter();
  const [stepStatus, setStepStatus] = React.useState({ kubernetes: false, operator: false });
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
        return <KubernetesScreen setStepStatus={setStepStatus} />;
      case 1:
        return <MesheryOperatorScreen setStepStatus={setStepStatus} />;
      case 2:
        return <ServiceMeshScreen />;
      case 3:
        return <MetricsScreen />;
      default:
        return null;
    }
  };

  const isNextDisabled = () => {
    // if (activeStep === 0 && !stepStatus.kubernetes) return true
    // if (activeStep === 1 && !stepStatus.operator) return true
    return false;
  };

  const handleAdvancedSettingsClick = () => {
    switch (activeStep) {
      case 0:
        return;
      case 1:
        return;
      case 2:
        router.push('/settings#service-mesh');
        return;
      case 3:
        router.push('/settings#metrics');
        return;
      default:
        return null;
    }
  };

  return (
    <Container className={classes.container}>
      <Stepper steps={steps} activeStep={activeStep} handleUserClick={handleUserClick} />
      <Fade timeout={{ enter: '1500ms' }} in="true">
        <div>
          {activeStep === steps.length ? (
            <Fade timeout={{ enter: '500ms' }} in="true">
              <CompletionScreen handleUserClick={handleUserClick} />
            </Fade>
          ) : (
            <Grid
              container
              xs={12}
              style={{ marginTop: '4rem' }}
              justify="center"
              alignItems="flex-start"
            >
              {handleStep(activeStep)}
            </Grid>
          )}
        </div>
      </Fade>
      {activeStep === steps.length ? null : (
        <div className={classes.buttonContainer}>
          {activeStep === 2 || activeStep === 3 ? (
            <Button onClick={handleAdvancedSettingsClick} className={classes.skipButton}>
              Advanced Settings
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
            // disabled={activeStep === 0 && !iskubernetesConnected}
            variant="contained"
            onClick={handleNext}
            disabled={isNextDisabled()}
            className={classes.button}
          >
            Next
          </Button>
        </div>
      )}
    </Container>
  );
};

export default ConfigurationWizard;
