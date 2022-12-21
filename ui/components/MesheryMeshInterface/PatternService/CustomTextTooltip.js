import { withStyles } from "@material-ui/styles";
import { Tooltip } from "@material-ui/core";

export const CustomTextTooltip = ({ backgroundColor, flag, ...props }) => {
  const CustomTooltip = withStyles(() => ({
    tooltip : {
      backgroundColor : backgroundColor,
      color : "#fff",
      opacity : "100%",
      fontSize : "0.75rem",
      fontFamily : flag ? "Qanelas Soft, sans-serif" : "inherit" ,
      borderRadius : "0.9375rem",
      padding : "0.9rem",
    },
    popper : {
      zIndex : "99999 !important"
    }
  }))(Tooltip);

  return <CustomTooltip {...props} />;
};
