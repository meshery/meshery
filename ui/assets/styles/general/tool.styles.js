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
  mainContainer: {
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    padding: '5rem',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    borderRadius: '1rem',
  },
  treeWrapper: {
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    display: 'flex',
    padding: '1rem',
    flexDirection: 'row',
  },
  '@keyframes rotateCloseIcon': {
    form: {
      transform: 'translateY(0) translateX(0)',
      height: '10rem',
      marginBottom: '0rem',
      background: '#51636B',
      flexDirection: 'column',
    },
    to: {
      transform: 'translateY(-70px) translateX(-100px)',
      height: '4rem',
      marginBottom: '5rem',
      background: '#677a84',
      flexDirection: 'row-reverse',
    },
  },
  cardStyle: {
    background: '#51636B',
    color: 'white',
    height: '10rem',
    width: '13rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: '1rem',
    flexDirection: 'column',
    cursor: 'pointer',
  },
  cardAnimateStyle: {
    background: '#51636B',
    color: 'white',
    height: '10rem',
    width: '13rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: '1rem',
    flexDirection: 'column',
    cursor: 'pointer',
    animation: '$rotateCloseIcon 1s',
  },
  detailsContainer: {
    height: '30rem',
    width: '50%',
    margin: '1rem',
    backgroundColor: '#d9dadb80',
    borderRadius: '6px',
    padding: '2.5rem',
    overflowY: 'auto',
    boxShadow: 'inset 0 0 6px 2px rgba(0, 0, 0,0.4)',
  },
  emptyDetailsContainer: {
    height: '30rem',
    width: '50%',
    margin: '1rem',
    backgroundColor: '#d9dadb80',
    borderRadius: '6px',
    padding: '2.5rem',
    overflowY: 'auto',
    boxShadow: 'inset 0 0 6px 2px rgba(0, 0, 0,0.4)',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  },
  // Define other styles as needed
}));

export default useStyles;
