import React from "react";
import Paper from "@material-ui/core/Paper";
import Stepper from "@material-ui/core/Stepper";
import StepConnector from "@material-ui/core/StepConnector";
import Checkbox from "@material-ui/core/Checkbox";
import { withStyles } from "@material-ui/core/styles";

export const ColorlibConnector = styled(StepConnector)(({ theme }) => ({
  [`&.${stepConnectorClasses.alternativeLabel}`]: {
    top: 22
  },
  [`&.${stepConnectorClasses.active}`]: {
    [`& .${stepConnectorClasses.line}`]: {
      borderColor: "#00B39F"
    }
  },
  [`&.${stepConnectorClasses.completed}`]: {
    [`& .${stepConnectorClasses.line}`]: {
      borderColor: "#00B39F"
    }
  },
  [`& .${stepConnectorClasses.line}`]: {
    borderTopWidth: 3,
    borderColor: "#d7d7d9",
    borderRadius: 1
  }
}));

export const WelcomeContainer = styled("div")(() => ({
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  margin: "0 auto"
}));


export const ColorlibStepIconRoot = styled("div")(({ theme, ownerState }) => ({
  backgroundColor: "#d7d7d9",
  zIndex: 1,
  width: 50,
  height: 50,
  display: "flex",
  borderRadius: "50%",
  justifyContent: "center",
  alignItems: "center",
  ...((ownerState.active || ownerState.completed) && {
    backgroundColor: `#fff`,
    boxShadow: "0 0 0 3px #00B39F"
  }),
  ["@media (max-width:780px)"]: {
    width: 40,
    height: 40
  }
}));
export const ConnectionStepperContent = styled(Paper)(() => ({
  width: "60%",
  padding: "2rem",
  minHeight: "48rem",
  borderRadius: "1rem",
  marginTop: "2.5rem",
  ["@media (max-width:780px)"]: {
    width: "90%",
    minHeight: "0"
  },
  ["@media (max-width:500px)"]: {
    width: "95%",
    padding: "1rem",
    minHeight: "0"
  }
}));

export const useStyles = makeStyles({
  customLabelStyle: {
    fontSize: "0.875rem",
    ["@media (max-width:500px)"]: {
      fontSize: "0.7rem"
    }
  }
});

export const VerifyContent = styled(Paper)(() => ({
  width: "15rem",
  height: "8rem",
  padding: "1rem",
  marginTop: "10%"
}));

export const VerifyContainer = styled("div")(() => ({
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  margin: "0 auto"
}));

export const StepperContainer = styled(Stepper)(() => ({
  width: "80%",
  marginTop: "2rem",
  ["@media (max-width:780px)"]: {
    width: "auto",
    marginTop: "1rem"
  }
}));

export const ContentContainer = styled("div")(() => ({
  display: "flex",
  flexDirection: "column",
  width: "60%",
  justifyContent: "space-between",
  margin: "1rem 0 0 2rem",
  ["@media (max-width:780px)"]: {
    width: "auto",
    margin: "1rem 0",
    justifyContent: "flex-start"
  }
}));

export const ContentBody = styled("div")(() => ({
  marginBottom: "2rem"
}));

export const CustomCheckbox = styled(Checkbox)(({ theme }) => ({
  "&.Mui-checked": {
    color: theme.palette.keppelGreen
  },
  "&:hover": {
    backgroundColor: "transparent"
  }
}));
