import { withStyles } from "@material-ui/styles";
import { Tooltip } from "@material-ui/core";

const EnlargedTextTooltip = withStyles((theme) => ({
  tooltip : {
    backgroundColor : ' #3C494F',
    color : '#fff',
    opacity : '100%',
    fontSize : theme.typography.pxToRem(13),
    // border : '1px solid #dadde9',
    borderRadius : "0.9375rem",
    padding : "0.9rem",
  },
  popper : {
    zIndex : "99999 !important"
  }
}))(Tooltip);

const EnlargedTextErrorTooltip = withStyles((theme) => ({
  tooltip : {
    backgroundColor : '#B32700',
    color : '#fff',
    opacity : '100%',
    fontSize : theme.typography.pxToRem(13),
    // border : '1px solid #dadde9',
    borderRadius : "0.9375rem",
    padding : "0.9rem",
  },
  popper : {
    zIndex : "99999 !important"
  }
}))(Tooltip);


export { EnlargedTextTooltip, EnlargedTextErrorTooltip };