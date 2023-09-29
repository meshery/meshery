import * as React from 'react';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';
import { ColorlibConnector, useStyles } from '../styles2';
import Checkbox from '@material-ui/core/Checkbox';
import { useRouter } from 'next/router';
import TipsCarousel from '../../General/TipsCarousel';
import { ConnectionStepperTips } from '../helm-connect/constants'; //TODO: move this to common
import Stepper from '@material-ui/core/Stepper';

function StepperIcon(props) {
  const classes = useStyles();
  const { active, completed, className, stepIcons } = props;
  return (
    <Checkbox
      className={className}
      classes={{
        root: classes.colorlibStepIconRoot,
        checked: classes.active,
        disabled: classes.completed,
      }}
      checked={completed}
      icon={stepIcons[String(props.icon)]}
      checkedIcon={stepIcons[String(props.icon)]}
    />
  );
}

export default function CustomizedSteppers({ stepData }) {
  const { stepContent, stepIcons, steps } = stepData;
  const classes = useStyles();
  const [activeStep, setActiveStep] = React.useState(0);
  const router = useRouter();
  const { githubLogin, id } = router.query;
  const installationId = id;
  const [selectedRepository, setSelectedRepository] = React.useState([]);
  const [configuredRepository, setConfiguredRepository] = React.useState([]);

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
          {steps.map((label, index) => (
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
