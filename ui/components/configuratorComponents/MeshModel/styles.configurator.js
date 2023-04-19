import { makeStyles } from "@material-ui/core";

const styles = makeStyles((theme) => ({
  backButton : {
    marginRight : theme.spacing(2),
  },
  appBar : {
    marginBottom : "16px",
    backgroundColor : theme.palette.secondary.appBar,
    borderRadius : "8px"
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
    '& .CodeMirror' : {
      minHeight : "300px",
      height : '100%',
    }
  },
  patternType : {
    padding : '0px',
    paddingBottom : '5px',
    paddingTop : '5px',
    justifyContent : 'center'
  },
  formCtrl : { //
    width : "60px",
    minWidth : "60px",
    maxWidth : "60px",
    marginRight : 8,
  },
  autoComplete : {
    width : "120px",
    minWidth : "120px",
    maxWidth : 120,
  },
  autoComplete2 : {
    width : 250,
    marginLeft : 16,
    marginRight : "auto",
    padding : 0,
    // "& .MuiAutocomplete-inputRoot" : {
    //   padding : 0
    // },
    // '& .MuiInputBase-input' : {

    // }
  },
  btngroup : {
    marginLeft : "auto",
    overflowX : "auto",
    overflowY : "hidden"
  },
  paper : {
    backgroundColor : "#fcfcfc",
    padding : 8,
    height : "100%",
  },
  wrapper : {
    width : '100%'
  },
  heading : {
    fontSize : theme.typography.pxToRem(15),
    fontWeight : theme.typography.fontWeightRegular,
  },
}));

export default styles