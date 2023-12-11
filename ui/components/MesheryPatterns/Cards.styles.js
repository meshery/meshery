import { makeStyles } from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  // dialogTitle: {
  //   textAlign: 'center',
  //   minWidth: 400,
  //   padding: '10px',
  //   color: '#fff',
  //   backgroundColor:
  //     theme.palette.type === 'dark' ? theme.palette.secondary.headerColor : '#396679',
  // },
  cardButtons: {
    display: 'flex',
    flexWrap: 'wrap',
    marginTop: '50px',
    height: '100%',
    gap: '.5rem',
  },
  testsButton: {
    padding: '6px 9px',
  },
  perfResultsContainer: {
    marginTop: '0.5rem',
  },
  flipButton: {
    minWidth: 'max-content',
    padding: '6px 9px',
    twhiteSpace: 'nowwrap',
  },
  backGrid: {
    marginBottom: '0.25rem',
    minHeight: '6rem',
    position: 'relative',
  },
  updateDeleteButtons: {
    width: 'fit-content',
    margin: '10 0 0 auto',
    position: 'absolute',
    right: 0,
    bottom: 0,
  },
  // yamlDialogTitle: {
  //   display: 'flex',
  //   alignItems: 'center',
  // },
  // yamlDialogTitleText: {
  //   flexGrow: 1,
  // },
  fullScreenCodeMirror: {
    height: '100%',
    width: '100%',
    '& .CodeMirror': {
      minHeight: '300px',
      height: '100%',
      width: '100%',
    },
  },
  noOfResultsContainer: {
    margin: '0 0 1rem',
    '& div': {
      display: 'flex',
      alignItems: 'center',
    },
  },
  bottomPart: {
    display: 'flex',
    justifyContent: 'flex-start',
    alignItems: 'center',
    marginTop: '1rem',
  },
  lastRunText: {
    marginRight: '0.5rem',
  },
  cardHeaderRight: {
    display: 'flex',
    alignItems: 'center',
  },
  iconPatt: {
    width: '24px',
    height: '24px',
    marginRight: '5px',
  },
  iconDownload: {
    width: 'auto',
    height: '24px',
  },
  btnText: {
    [theme.breakpoints.down(1370)]: { display: 'none' },
    [`${theme.breakpoints.up(1920)} and (max-width: 2200px)`]: {
      display: 'none',
    },
    marginLeft: '5px',
    display: 'flex',
    justifyContent: 'center',
  },
  cloneBtnText: {
    [theme.breakpoints.down(1370)]: { display: 'none' },
    [`${theme.breakpoints.up(1920)} and (max-width: 2200px)`]: {
      display: 'none',
    },
    display: 'flex',
    justifyContent: 'center',
    marginLeft: '3px',
  },
  undeployButton: {
    backgroundColor: '#8F1F00',
    color: '#ffffff',
    padding: '6px 9px',
    minWidth: 'unset',
    '&:hover': {
      backgroundColor: '#B32700',
    },
    '& > span': {
      width: 'unset',
    },
    '& > span > svg': {
      marginRight: '5px',
    },
  },
  img: {
    marginRight: '0.5rem',
    filter: theme.palette.secondary.img,
  },
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
  clonePatt: {
    width: '20px',
    height: '20px',
    marginRight: '5px',
  },
}));

export default useStyles;
