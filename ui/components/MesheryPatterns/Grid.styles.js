import { makeStyles } from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  noPaper: {
    padding: '0.5rem',
    fontSize: '3rem',
  },
  noContainer: {
    padding: '2rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
  },
  noText: {
    fontSize: '2rem',
    marginBottom: '2rem',
  },
  pagination: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: '2rem',
  },
  addIcon: {
    paddingLeft: theme.spacing(0.5),
    marginRight: theme.spacing(0.5),
  },
}));

export default useStyles;
