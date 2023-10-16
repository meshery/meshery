import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  toolWrapper: {
    marginBottom: '3rem',
    display: 'flex',
    justifyContent: 'space-between',
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    boxShadow: ' 0px 2px 4px -1px rgba(0,0,0,0.2)',
    height: '4rem',
    padding: '0.68rem',
    borderRadius: '0.5rem',
    position: 'relative',
    zIndex: '125',
  },
  // Define other styles as needed
}));

export default useStyles;
