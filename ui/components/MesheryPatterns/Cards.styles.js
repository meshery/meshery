import { makeStyles } from "@material-ui/core";

const useStyles = makeStyles((theme) => ({
  cardButtons : {
    display : "grid",
    gridTemplateColumns : "repeat(5,1fr)",
    marginTop : "50px",
    height : '100%',
    width : '100%',
    gap : ".5rem"
  },
  testsButton : {
    minWidth : "100%",
    padding : "6px 9px"
  },
  perfResultsContainer : {
    marginTop : "0.5rem"
  },
  backGrid : {
    marginBottom : "0.25rem",
    minHeight : "6rem",
    position : "relative"
  },
  updateDeleteButtons : {
    width : "fit-content",
    margin : "10 0 0 auto",
    position : "absolute",
    right : 0,
    bottom : 0,
  },
  yamlDialogTitle : {
    display : "flex",
    alignItems : "center"
  },
  yamlDialogTitleText : {
    flexGrow : 1
  },
  fullScreenCodeMirror : {
    height : '100%',
    width : '100%',
    '& .CodeMirror' : {
      minHeight : "300px",
      height : '100%',
      width : '100%'
    }
  },
  maximizeButton : {
    width : "fit-content",
    margin : "0 0 0 auto",
    position : "absolute",
    right : 0,
    top : 0
  },
  noOfResultsContainer : {
    margin : "0 0 1rem",
    '& div' : {
      display : "flex",
      alignItems : "center"
    },
  },
  bottomPart : {
    display : "flex",
    justifyContent : "flex-end",
    alignItems : "center",
  },
  lastRunText : {
    marginRight : "0.5rem"
  },
  iconPatt : {
    width : "24px",
    height : "24px",
    marginRight : "5px",

  },
  btnText : {
    [theme.breakpoints.between(960, 1370)] : { display : "none" },
    display : "flex",
    justifyContent : "center",
  },
  cloneBtnText : {
    [theme.breakpoints.between(960, 1370)] : { display : "none" },
    display : "flex",
    justifyContent : "center",
    marginLeft : "3px",
  },
  undeployButton : {
    backgroundColor : "#8F1F00",
    color : "#ffffff",
    padding : "6px 9px",
    minWidth : "unset",
    "&:hover" : {
      backgroundColor : "#B32700",
    },
    '& > span' : {
      width : "unset",
    },
    '& > span > svg' : {
      marginRight : "5px"
    }
  },
  img : {
    marginRight : "0.5rem",
    filter : theme.palette.secondary.img
  },
  noPaper : {
    padding : "0.5rem",
    fontSize : "3rem"
  },
  noContainer : {
    padding : "2rem",
    display : "flex",
    justifyContent : "center",
    alignItems : "center",
    flexDirection : "column",
  },
  publishTitle : {
    display : 'flex',
    justifyContent : 'space-between',
    alignItems : 'center',
  },
  noText : {
    fontSize : "2rem",
    marginBottom : "2rem",
  },
  clonePatt : {
    width : "20px",
    height : "20px",
    marginRight : "5px",
  },
}));

export default useStyles;