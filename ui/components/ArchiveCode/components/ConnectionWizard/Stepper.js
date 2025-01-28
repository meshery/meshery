import React from 'react';
import clsx from 'clsx';
import {
  makeStyles,
  withStyles,
  Stepper as Stepperr,
  Step,
  StepLabel,
  StepConnector,
} from '@material-ui/core/';

import svgIcons from './icons/icons.js';

const ColorlibConnector = withStyles({
  alternativeLabel: { top: 22 },
  active: { '& $line': { background: '#00B39F', transition: 'all 1s ease-in' } },
  completed: { '& $line': { background: '#00B39F', transition: 'all 1s ease-in' } },
  line: {
    height: 3,
    border: 0,
    backgroundColor: '#eaeaf0',
    borderRadius: 1,
    transition: 'all 0.5s ease-out ',
  },
})(StepConnector);

const useColorlibStepIconStyles = makeStyles({
  icnlist: {
    backgroundColor: '#ccc',
    zIndex: 1,
    color: '#fff',
    width: 50,
    height: 50,
    display: 'flex',
    border: '.2rem solid #ccc',
    borderRadius: '50%',
    justifyContent: 'center',
    alignItems: 'center',
    cursor: 'pointer',
  },
  active: {
    background: '#fff',
    color: '#3C494E',
    boxShadow: '0 4px 10px 0 rgba(0,0,0,.25)',
    border: '.2rem solid #00B39F',
    transition: 'all 0.5s ease-in',
  },
  completed: {
    border: '.2rem solid #00B39F',
    background: '#00B39F',
    transition: 'all 0.5s ease-in',
  },
});

const useStepperStyles = makeStyles({ removePadding: { padding: '0 !important' } });

function ColorlibStepIcon(props) {
  const classes = useColorlibStepIconStyles();
  const { active, completed } = props;

  const icons = svgIcons(props.completed, props.active);

  return (
    <div
      className={clsx(classes.icnlist, {
        [classes.active]: active,
        [classes.completed]: completed,
      })}
      style={{ position: 'relative', cursor: 'default' }}
    >
      {icons[String(props.icon)]}
    </div>
  );
}

const Stepper = ({ steps, activeStep, handleUserClick }) => {
  const handleChange = (label, steps) => handleUserClick(steps.indexOf(label));
  const classes = useStepperStyles();
  return (
    <Stepperr
      alternativeLabel
      activeStep={activeStep}
      connector={<ColorlibConnector />}
      className={classes.removePadding}
    >
      {steps.map((label) => (
        <Step key={label} onClick={() => handleChange(label, steps)}>
          <StepLabel StepIconComponent={ColorlibStepIcon}>{label}</StepLabel>
        </Step>
      ))}
    </Stepperr>
  );
};

export default Stepper;
