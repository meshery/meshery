import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepButton from '@material-ui/core/StepButton';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';
import MeshAdapterConfigComponent from './MeshAdapterConfigComponent';
import MeshConfigComponent from './MeshConfigComponent';
import GrafanaComponent from './GrafanaComponent';
import {connect} from "react-redux";
import { Divider, StepLabel, Icon } from '@material-ui/core';
import PrometheusComponent from './PrometheusComponent';

const styles = theme => ({
  root: {
    // width: '90%',
    padding: theme.spacing(10),
  },
  button: {
    marginRight: theme.spacing(1),
  },
  completed: {
    display: 'inline-block',
  },
  instructions: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
  titleDecorate: {
    borderBottom: '1px solid',
    paddingBottom: theme.spacing(1),
  },
  stepperButtons: {
    textAlign: 'right',
    marginTop: theme.spacing(2),
    marginRight: theme.spacing(3),
  }
});

function getSteps() {
  return ['Kubernetes', 'Adapters', 'Grafana', 'Prometheus'];
}

function getRequiredSteps(){
  return [true, true, false, false];
}

function getStepContent(step) {
  switch (step) {
  case 0:
    return (
      <MeshConfigComponent />
    );
  case 1:
    return (
      <MeshAdapterConfigComponent />
    );
  case 2:
    return (
      <GrafanaComponent />
    );
  case 3:
    return (
      <PrometheusComponent />
    );
  }
}

class MesheryConfigSteps extends React.Component {
  constructor(props){
    super(props);
    const {k8sconfig, meshAdapters, grafana, prometheus} = props;
    this.state = {
      activeStep: 0,
      completed: {},
      k8sconfig, 
      meshAdapters, 
      grafana,
      prometheus,
    };
  }

  static getDerivedStateFromProps(props, state) {
    const {k8sconfig, meshAdapters, grafana, prometheus} = props;
    const { completed } = state;

    if(k8sconfig.clusterConfigured) {
      completed[0] = true;
    } else {
      completed[0] = false;
    }
    if(meshAdapters.length > 0){
      completed[1] = true;
    } else {
      completed[1] = false;
    }
    if (grafana.grafanaURL !== '' && grafana.selectedBoardsConfigs.length > 0) {
      completed[2] = true;
    } else {
      completed[2] = false;
    }
    if (prometheus.prometheusURL !== '' && prometheus.selectedPrometheusBoardsConfigs.length > 0) {
      completed[3] = true;
    } else {
      completed[3] = false;
    }

    return {k8sconfig, meshAdapters, grafana, prometheus, completed};
  }

  totalSteps = () => getSteps().length;

  handleNext = () => {
    let activeStep;

    if (this.isLastStep() && !this.allStepsCompleted()) {
      // It's the last step, but not all steps have been completed,
      // find the first step that has been completed
      const steps = getSteps();
      activeStep = steps.findIndex((step, i) => !(i in this.state.completed));
    } else {
      switch(this.state.activeStep){
      case 0:
      case 1:
      case 2:
        activeStep = this.state.activeStep + 1;
        break;
      case 3:
        activeStep = 0;
        break;
      }
    }
    this.setState({
      activeStep,
    });
  };

  handleBack = () => {
    let activeStep;
    switch(this.state.activeStep){
    case 0:
      activeStep = 3;
      break;
    case 1:
    case 2:
    case 3:
      activeStep = this.state.activeStep - 1;
      break;
    }
    this.setState({activeStep});
  };

  handleStep = step => () => {
    this.setState({
      activeStep: step,
    });
  };

  handleComplete = () => {
    const { completed } = this.state;
    completed[this.state.activeStep] = true;
    this.setState({
      completed,
    });
    this.handleNext();
  };

  handleReset = () => {
    this.setState({
      activeStep: 0,
      completed: {},
    });
  };

  completedSteps() {
    return Object.keys(this.state.completed).length;
  }

  isLastStep() {
    return this.state.activeStep === this.totalSteps() - 1;
  }

  allStepsCompleted() {
    return this.completedSteps() === this.totalSteps();
  }

  render() {
    const { classes } = this.props;
    const steps = getSteps();
    const { activeStep } = this.state;

    return (
      <div className={classes.root}>
        <Stepper nonLinear activeStep={activeStep}>
          {steps.map((label, index) => (
            <Step key={label}>
              <StepButton onClick={this.handleStep(index)} completed={this.state.completed[index]}>
                <StepLabel className={activeStep == index?classes.titleDecorate:''}
                  // icon={<Icon>k</Icon>}
                >
                  <Typography variant="h6">
                    {label}{getRequiredSteps()[index] && (<Typography variant="h5" style={{
                      display: 'inline',
                    }}><sup>*</sup></Typography>)}
                  </Typography>
                </StepLabel>
              </StepButton>
            </Step>
          ))}
        </Stepper>
        <div>
          
          <div>
            <Typography className={classes.instructions}>{getStepContent(activeStep)}</Typography>
            <Divider light variant="fullWidth" />
            <div className={classes.stepperButtons}>
              {activeStep !== steps.length &&
                  (this.state.completed[this.state.activeStep] ? (
                    <Typography variant="caption" className={classes.completed}>
                      Step {activeStep + 1} is complete
                    </Typography>
                  ) : '')}
              <Button
                size='large'
                // disabled={activeStep === 0}
                onClick={this.handleBack}
                className={classes.button}
              >
                  Back
              </Button>
              <Button
                size='large'
                variant="outlined"
                color="primary"
                onClick={this.handleNext}
                className={classes.button}
              >
                  Next
              </Button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  const k8sconfig = state.get("k8sConfig").toJS();
  const meshAdapters = state.get("meshAdapters").toJS();
  const grafana = state.get("grafana").toJS();
  const prometheus = state.get("prometheus").toJS();
  return {k8sconfig, meshAdapters, grafana, prometheus};
}

MesheryConfigSteps.propTypes = {
  classes: PropTypes.object,
};

export default withStyles(styles)(connect(
  mapStateToProps
)(MesheryConfigSteps));