import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  toolWrapper: {
    display: 'flex',
    justifyContent: 'space-between',
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    boxShadow: ' 0px 2px 4px -1px rgba(0,0,0,0.2)',
    height: '0rem',
    padding: '0rem',
    borderRadius: '0.5rem',
    position: 'relative',
    zIndex: '0',
    marginBottom: '0.5rem',
    marginTop: '1rem',
    transition: 'all 1s',
  },
  toolWrapperAnimate: {
    height: '4rem',
    zIndex: '125',
    padding: '0.68rem',
  },
  mainContainer: {
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    height: '25rem',
    borderRadius: '1rem',
    display: 'flex',
    position: 'relative',
    marginTop: '1rem',
    transition: 'all 1s',
  },
  mainContainerAnimate: {
    height: '36rem',
  },
  innerContainer: {
    display: 'flex',
    flexDirection: 'row',
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    transition: 'all 1s',
  },
  innerContainerAnimate: {
    width: '100%',
    top: '0%',
    paddingLeft: '2rem',
    transform: 'translate(0%, 0%)',
    borderRadius: '6px 6px 0px 0px',
    left: '0%',
    backgroundColor: '#51636B',
  },
  tabs: {
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
  },
  treeWrapper: {
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    display: 'flex',
    flexDirection: 'row',
    width: '100%',
    position: 'absolute',
    top: '3.7rem',
    opacity: '0',
    transition: 'all 1s',
  },
  treeWrapperAnimate: {
    opacity: '1',
  },
  cardbg: {
    backgroundColor: 'white'
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
    transition: 'all 1s',
  },
  cardStyleAnimate: {
    marginTop: '0.7rem',
    height: '3rem',
    backgroundColor: '#677a84',
    flexDirection: 'row-reverse',
    borderRadius: '8px 8px 0px 0px',
    paddingTop: '0.2rem',
  },
  overviewTab: {
    background: '#51636B',
    color: 'white',
    height: '10rem',
    width: '0rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: '1rem',
    flexDirection: 'column',
    cursor: 'pointer',
    visibility: 'hidden',
    opacity: '0',
    transition: 'all 1s',
  },
  overviewTabAnimate: {
    marginTop: '0.7rem',
    height: '3rem',
    width: '13rem',
    opacity: '1',
    backgroundColor: '#677a84',
    flexDirection: 'row-reverse',
    borderRadius: '8px 8px 0px 0px',
    visibility: 'visible',
    paddingTop: '0.2rem',
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
