import { withStyles } from "@material-ui/styles";
import { Tooltip } from "@material-ui/core";

const EnlargedTextTooltip = withStyles((theme) => ({
  tooltip : {
    backgroundColor : '#f5f5f9',
    color : 'rgba(0, 0, 0, 0.87)',
    fontSize : theme.typography.pxToRem(15),
    border : '1px solid #dadde9',
  },
}))(Tooltip);

export default EnlargedTextTooltip;