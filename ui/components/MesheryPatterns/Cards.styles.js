import { makeStyles } from "@material-ui/core";

const useStyles= makeStyles(() => ({
  cardButtons : {
    display : "flex",
    justifyContent : "flex-end",
    alignItems : "center",
    marginTop : "50px"
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
    display : "flex",
    justifyContent : "center",
    paddingLeft : "0.35rem"
  },
  undeployButton : {
    backgroundColor : "#B32700",
    color : "#ffffff",
    "&:hover" : {
      backgroundColor : "#8f1f00"
    },
    marginRight : "0.5rem"
  }
}));

export default useStyles;