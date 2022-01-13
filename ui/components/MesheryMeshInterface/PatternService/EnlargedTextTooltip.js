import { withStyles } from "@material-ui/styles";
import { Tooltip } from "@material-ui/core";

const EnlargedTextTooltip = withStyles((theme) => ({
  tooltip : {
    backgroundColor : ' #fcfcfc',
    color : '#1E2117',
    fontSize : theme.typography.pxToRem(12),
    border : '1px solid #dadde9',
    borderRadius : "0.9375rem",
    padding : "0.625rem",

  },
}))(Tooltip);

export default EnlargedTextTooltip;