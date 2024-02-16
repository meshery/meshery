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
    zIndex: '101',
  },
  meshModelToolbar: {
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
    display: 'flex',
    position: 'relative',
    marginTop: '1rem',
    [theme.breakpoints.down('sm')]: {
      height: '47rem',
    },
  },
  mainContainerAnimate: {
    height: '36rem',
    [theme.breakpoints.down('sm')]: {
      height: '73rem',
    },
  },
  innerContainer: {
    display: 'flex',
    flexDirection: 'row',
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    [theme.breakpoints.down('sm')]: {
      flexDirection: 'column',
    },
  },
  innerContainerAnimate: {
    width: '100%',
    top: '0%',
    paddingX: '2rem',
    transform: 'translate(0%, 0%)',
    display: 'flex',
    justifyContent: 'center',
    left: '0%',
    backgroundColor: theme.palette.secondary.tabContainer,
    [theme.breakpoints.down('sm')]: {
      flexDirection: 'row',
      paddingLeft: '1rem',
      overflowX: 'auto',
      padding: '0.4rem'
    },
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
    [theme.breakpoints.down('sm')]: {
      flexDirection: 'column-reverse',
    },
  },
  treeWrapperAnimate: {
    opacity: '1',
  },
  cardbg: {
    backgroundColor: 'white',
  },
  cardStyle: {
    background: theme.palette.secondary.card,
    color: 'white',
    height: '10rem',
    width: '13rem',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: '0rem 0.7rem',
    flexDirection: 'column',
    cursor: 'pointer',
    [theme.breakpoints.down('md')]: {
      height: '10rem',
      width: '8.5rem'
    },
    [theme.breakpoints.down('sm')]: {
      width: '13rem',
      marginTop: '0.7rem',
      marginRight: '0.5rem',
    },
  },
  cardStyleAnimate: {
    marginTop: '0.7rem',
    margin: '0rem 0.7rem',
    height: '3rem',
    width: '15rem',
    flexDirection: 'row-reverse',
    borderRadius: '8px 8px 0px 0px',
    paddingTop: '0.2rem',
    backgroundColor: theme.palette.secondary.tabCard,
    [theme.breakpoints.down('sm')]: {
      padding: '0.1rem',
      flexDirection: 'column',
      margin: '0rem 0.2rem',
      width: '10rem',
    },
  },
  detailsContainer: {
    height: '30rem',
    width: '50%',
    margin: '1rem',
    backgroundColor: theme.palette.secondary.detailsContainer,
    borderRadius: '6px',
    padding: '1rem 2rem',
    overflowY: 'auto',
    boxShadow: 'inset 0 0 6px 2px rgba(0, 0, 0,0.4)',
    [theme.breakpoints.down('sm')]: {
      width: '90%',
      padding: '1rem',
      height: 'fit-content',
      maxHeight: '30rem',
    },
  },
  emptyDetailsContainer: {
    height: '30rem',
    width: '50%',
    margin: '1rem',
    backgroundColor: theme.palette.secondary.detailsContainer,
    borderRadius: '6px',
    padding: '2.5rem',
    overflowY: 'auto',
    boxShadow: 'inset 0 0 6px 2px rgba(0, 0, 0,0.4)',
    display: 'flex',
    flexDirection: 'column',
    // justifyContent: 'center',
    alignItems: 'center',
    [theme.breakpoints.down('sm')]: {
      width: '90%',
      padding: '0.5rem',
      height: 'fit-content',
      maxHeight: '30rem',
    },
  },
  activeTab: {
    backgroundColor:
      theme.palette.type === 'dark'
        ? theme.palette.secondary.toolbarBg2
        : theme.palette.secondary.toolbarBg1,
    color: theme.palette.secondary.text,
  },
  line: {
    borderLeft: `1px dashed ${theme.palette.secondary.text}`,
  },
  treeContainer: {
    height: '30rem',
    width: '50%',
    margin: '1rem',
    display: 'flex',
    justifyContent: 'center',
    [theme.breakpoints.down('sm')]: {
      width: '90%',
    },
  },
  searchContainer: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    borderBottom: '1px solid #d2d3d4',
  },
  segment: {
    display: 'flex',
    flexDirection: 'row',
    [theme.breakpoints.down('sm')]: {
      flexDirection: 'column'
    },
  },
  fullWidth: {
    width: '50%',
    [theme.breakpoints.down('sm')]: {
      width: '70%'
    },
  }
  // Define other styles as needed
}));


export default useStyles;