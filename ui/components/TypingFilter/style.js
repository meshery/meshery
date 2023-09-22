import { makeStyles,alpha } from "@material-ui/core";

export const useStyles = makeStyles((theme) => ({
  root : {
    position : "relative",
    backgroundColor : theme.palette.secondary.elevatedComponents,
  },
  input : {
    width : "100%",
    marginBottom : ".1rem",
    "& .MuiOutlinedInput-root" : {
      borderRadius : "6px",
      backgroundColor : theme.palette.secondary.searchBackground,
      "& fieldset" : {
        borderRadius : "6px",
        border : `2px solid ${theme.palette.secondary.searchBorder}`,
      },
    },
  },

  dropDown : {
    backgroundColor : theme.palette.secondary.searchBackground,
    borderRadius : "6px",
    boxShadow :
      "0px 2px 4px 0px rgba(0, 0, 0, 0.20), 0px 1px 10px 0px rgba(0, 0, 0, 0.12), 0px 4px 5px 0px rgba(0, 0, 0, 0.14)",
    border : `2px solid ${theme.palette.secondary.searchBorder}`,
    marginTop : "0.2rem",
  },
}));

export const useFilterStyles = makeStyles((theme) => ({
  item : {
    fontFamily : "Qanelas Soft, sans-serif",
    display : "flex",
    gap : "0.3rem",
    margin : "0.3rem",
    padding : "0.3rem",
    paddingInline : "3rem",
    borderRadius : "6px",
    cursor : "pointer",
    "&:hover" : {
      backgroundColor : alpha(theme.palette.secondary.link2, 0.25),
    },
  },

  label : {
    fontWeight : 500,
    color : theme.palette.secondary.icon,
  },
  description : {
    fontWeight : 400,
    color : theme.palette.secondary.number,
  },
}));

