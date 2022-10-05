import { makeStyles } from "@material-ui/core";

const useStyles= makeStyles((theme) => ({
  cardButtons : {
    display : "flex",
    justifyContent : "flex-end",
    alignItems : "center",
    marginTop : "50px",
    marginLeft : "22rem",
    height : '100%',
    width : '100%',
  },
  testsButton : {
    marginRight : "0.5rem"
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
    paddingLeft : "0.35rem"
  },
  undeployButton : {
    backgroundColor : "#8F1F00",
    color : "#ffffff",
    "&:hover" : {
      backgroundColor : "#B32700",
    },
    marginRight : "0.5rem"
  },
  img : {
    marginRight : "0.5rem"
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
  noText : {
    fontSize : "2rem",
    marginBottom : "2rem",
  }
}));

export default useStyles;