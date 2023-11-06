import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  icon: {
    fill: theme.palette.secondary.iconMain,
  },
}));

const ExpandMoreIcon = (props) => {
  const classes = useStyles();
  return (
    <svg className={classes.icon} width="24" height="24" viewBox="0 0 24 24" fill={props.fill} xmlns="http://www.w3.org/2000/svg">
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M16.59 8.59003L12 13.17L7.41 8.59003L6 10L12 16L18 10L16.59 8.59003Z"
      />
    </svg>
  );
};

export default ExpandMoreIcon;
