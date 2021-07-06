import React from "react";
import clsx from "clsx";
import {
  makeStyles,
  withStyles,
  Stepper as Stepperr,
  Step,
  StepLabel,
  StepConnector,
} from "@material-ui/core/";

import svgIcons from "../icons/icons";

const ColorlibConnector = withStyles({
  alternativeLabel: {
    top: 22,
  },
  active: {
    "& $line": {
      background: "#00B39F",
      transition: "all 1s ease-in",
    },
  },
  completed: {
    "& $line": {
      background: "#00B39F",
      transition: "all 1s ease-in",
    },
  },
  line: {
    height: 3,
    border: 0,
    backgroundColor: "#eaeaf0",
    borderRadius: 1,
    transition: "all 0.5s ease-out ",
  },
})(StepConnector);

const useColorlibStepIconStyles = makeStyles({
  root: {
    backgroundColor: "#ccc",
    zIndex: 1,
    color: "#fff",
    width: 50,
    height: 50,
    display: "flex",
    border: ".2rem solid #ccc",
    borderRadius: "50%",
    justifyContent: "center",
    alignItems: "center",
    cursor: "pointer",
  },
  active: {
    background: "#fff",
    color: "#3C494E",
    boxShadow: "0 4px 10px 0 rgba(0,0,0,.25)",
    border: ".2rem solid #00B39F",
    transition: "all 0.5s ease-in",
  },
  completed: {
    border: ".2rem solid #00B39F",
    background: "#00B39F",
    transition: "all 0.5s ease-in",
  },
});

function ColorlibStepIcon(props) {
  const classes = useColorlibStepIconStyles();
  const { active, completed } = props;

  const icons = svgIcons(props.completed, props.active);

  return (
    <div
      className={clsx(classes.root, {
        [classes.active]: active,
        [classes.completed]: completed,
      })}
    >
      {icons[String(props.icon)]}
    </div>
  );
}

const Stepper = ({ steps, activeStep, handleUserClick }) => {
  const handleChange = (label, steps) => handleUserClick(steps.indexOf(label));

  return (
    <Stepperr
      alternativeLabel
      activeStep={activeStep}
      connector={<ColorlibConnector />}
    >
      {steps.map((label) => (
        <Step
          key={label}
          onClick={() => handleChange(label, steps)}
        >
          <StepLabel StepIconComponent={ColorlibStepIcon}>{label}</StepLabel>
        </Step>
      ))}
    </Stepperr>
  );
};

export default Stepper;
