import * as React from 'react';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import { ColorlibConnector, useStyles, useColorlibStepIconStyles } from '../styles';
import { useRouter } from 'next/router';
import TipsCarousel from '../../General/TipsCarousel';
import { ConnectionStepperTips } from '../connections/helm-connect/constants'; //TODO: move this to common
import Stepper from '@material-ui/core/Stepper';
import clsx from 'clsx';

function StepperIcon(props) {
  const classes = useColorlibStepIconStyles();
  const { active, completed, stepIcons } = props;

  const iconComponent = stepIcons[String(props.icon)];

  const additionalProps = {
    fill: completed ? "white" : "currentColor",
  }

  return (
    <div
      className={clsx(classes.icnlist, {
        [classes.active]: active,
        [classes.completed]: completed,
      })}
      style={{ position: 'relative', cursor: 'default' }}
    >
      {React.cloneElement(iconComponent, additionalProps)}
    </div>
  );
}


/* eslint-disable */

export default function CustomizedSteppers({ stepData }) {
  const { stepContent, stepIcons, steps } = stepData;
  const classes = useStyles();
  const [activeStep, setActiveStep] = React.useState(0);
  const router = useRouter();
  const { githubLogin } = router.query;
  // State to store data that will be shared with other steppers
  // To have less conflict, use sharedData as object with descriptive key
  const [sharedData, setSharedData] = React.useState(null);

  const ActiveStepContent = stepContent[String(activeStep + 1)].component;

  const handleNext = () => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
  };

  const stepProps = stepContent[String(activeStep + 1)]?.props?.reduce((props, propName) => {
    props[propName] = propName === 'handleNext' ? handleNext : eval(propName);
    return props;
  }, {});

  React.useEffect(() => {
    if (githubLogin === 'true') {
      handleNext();
    }
  }, [githubLogin]);

  return (
    <div style={{ display: 'flex', flexDirection: 'column' }}>
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <Stepper
          alternativeLabel
          activeStep={activeStep}
          connector={<ColorlibConnector />}
          classes={{ root: classes.stepperContainer }}
        >
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel
                StepIconComponent={(props) => <StepperIcon {...props} stepIcons={stepIcons} />}
                classes={{ label: classes.customLabelStyle }}
              >
                {label}
              </StepLabel>
            </Step>
          ))}
        </Stepper>
      </div>
      <div style={{ marginTop: '2rem', display: 'flex', minHeight: '58vh' }}>
        <TipsCarousel tips={ConnectionStepperTips} />
        {React.cloneElement(ActiveStepContent, stepProps)}
      </div>
    </div>
  );
}
