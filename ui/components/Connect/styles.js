import { makeStyles, withStyles } from '@material-ui/core';
import StepConnector from '@material-ui/core/StepConnector';

export const ColorlibConnector = withStyles({
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

export const useColorlibStepIconStyles = makeStyles((theme) => ({
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
        '@media (max-width:780px)': {
      width: 40,
      height: 40,
    },
  },
  active: {
    background: theme.palette.secondary.elevatedComponents,
    color: '#3C494E',
    border: '.2rem solid #00B39F',
    transition: 'all 0.5s ease-in',
  },
  completed: {
    border: '.2rem solid #00B39F',
    background: '#00B39F',
    transition: 'all 0.5s ease-in',
  },
}));

export const useStyles = makeStyles((theme) => ({
  welcomeContainer: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: '0 auto',
  },
  connectionStepperContent: {
    width: '60%',
    padding: '2rem',
    minHeight: '48rem',
    borderRadius: '1rem',
    marginTop: '2.5rem',
    ['@media (max-width:780px)']: {
      width: '90%',
      minHeight: '0',
    },
    ['@media (max-width:500px)']: {
      width: '95%',
      padding: '1rem',
      minHeight: '0',
    },
  },
  customLabelStyle: {
    fontSize: '0.875rem',
    ['@media (max-width:500px)']: {
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
    ['@media (max-width:780px)']: {
      width: 'auto',
      marginTop: '1rem',
    },
  },
  cancelButton: {
    background: theme.palette.secondary.penColorPrimary,
    color: theme.palette.secondary.primaryModalText,
    margin: 0
  },
  stepperButton: {
    marginTop: '1rem',
    backgroundColor: theme.palette.secondary.focused,
    color: theme.palette.secondary.primaryModalText,
    '&:hover': {
      backgroundColor: theme.palette.secondary.focused,
    },
    '&.Mui-disabled': {
      cursor: 'not-allowed',
    },
  },
  contentContainer: {
    display: 'flex',
    flexDirection: 'column',
    width: '60%',
    justifyContent: 'space-between',
    margin: '1rem 0 0 2rem',
    ['@media (max-width:780px)']: {
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
