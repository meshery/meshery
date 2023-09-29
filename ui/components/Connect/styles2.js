import { makeStyles, withStyles } from '@material-ui/core';
import StepConnector from '@material-ui/core/StepConnector';

export const ColorlibConnector = withStyles({
  alternativeLabel: {
    top: 22,
  },
  active: {
    '& $line': {
      borderColor: '#00B39F',
    },
  },
  completed: {
    '& $line': {
      borderColor: '#00B39F',
    },
  },
  line: {
    borderTopWidth: 3,
    borderColor: '#d7d7d9',
    borderRadius: 1,
  },
})(StepConnector);

export const useStyles = makeStyles((theme) => ({
  welcomeContainer: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: '0 auto',
  },
  colorlibStepIconRoot: {
    backgroundColor: '#d7d7d9',
    zIndex: 1,
    width: 50,
    height: 50,
    display: 'flex',
    borderRadius: '50%',
    justifyContent: 'center',
    alignItems: 'center',
    '&$active, &$completed': {
      backgroundColor: `#fff`,
      boxShadow: '0 0 0 3px #00B39F',
    },
    '@media (max-width:780px)': {
      width: 40,
      height: 40,
    },
  },
  active: {
    backgroundColor: `#fff`,
    boxShadow: '0 0 0 3px #00B39F',
  },
  completed: {
    backgroundColor: `#fff`,
    boxShadow: '0 0 0 3px #00B39F',
  },
  connectionStepperContent: {
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
  },
  customLabelStyle: {
    fontSize: '0.875rem',
    [theme.breakpoints.down('xs')]: {
      fontSize: '0.7rem',
    },
  },
  verifyContent: {
    width: '15rem',
    height: '8rem',
    padding: '1rem',
    marginTop: '10%',
  },
  verifyContainer: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: '0 auto',
  },
  stepperContainer: {
    width: '80%',
    marginTop: '2rem',
    ["@media (max-width:780px)"]: {
      width: 'auto',
      marginTop: '1rem',
    },
  },
  cancelButton: {
    marginRight: "1rem",
    color: "#000",
    backgroundColor: "#fff",
    "&:hover": {
      backgroundColor: "#fff"
    },
  },
  stepperButton: {
    marginTop: "1rem",
    backgroundColor: theme.palette.btnBg,
    color: theme.palette.white,
    "&:hover": {
      backgroundColor: theme.palette.btnHover
    },
    "&.Mui-disabled": {
      backgroundColor: theme.palette.btnDisabled,
      color: theme.palette.gray,
      cursor: "not-allowed"
    }
  },
  contentContainer: {
    display: 'flex',
    flexDirection: 'column',
    width: '60%',
    justifyContent: 'space-between',
    margin: '1rem 0 0 2rem',
    ["@media (max-width:780px)"]: {
      width: 'auto',
      margin: '1rem 0',
      justifyContent: 'flex-start',
    },
  },
  contentBody: {
    marginBottom: '2rem',
  },
  customCheckbox: {
    '&.Mui-checked': {
      color: theme.palette.keppelGreen,
    },
    '&:hover': {
      backgroundColor: 'transparent',
    },
  },
}));
