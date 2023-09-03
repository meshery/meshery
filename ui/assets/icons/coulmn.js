import React from "react";
import { makeStyles } from "@material-ui/core/styles";


const useStyles = makeStyles((theme) => ({
  icon : {
    fill : theme.palette.type === 'dark' ? "#FFFFFF" : "#3C494F"
  },
}));
const ColumnIcon = (props) => {
  const classes = useStyles();
  return (
    <svg
      className={classes.icon}
      xmlns="http://www.w3.org/2000/svg"
      height="24"
      viewBox="0 -960 960 960"
      width="24"
      fill={props.fill}
    >
      <path d="M120-200v-560h213v560H120Zm253 0v-560h213v560H373Zm253 0v-560h213v560H626Z" />
    </svg>
  );
};

export default ColumnIcon;
