import { makeStyles } from '@material-ui/core';

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
    margin: 0,
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
