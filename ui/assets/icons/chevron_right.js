import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  icon: {
    fill: theme.palette.secondary.iconMain,
  },
}));

const ChevronRightIcon = (props) => {
  const classes = useStyles();
  return (
    <svg
      className={classes.icon}
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill={props.fill}
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M9.99984 6L8.58984 7.41L13.1698 12L8.58984 16.59L9.99984 18L15.9998 12L9.99984 6Z"
      />
    </svg>
  );
};

export default ChevronRightIcon;
