import React from 'react';


import { makeStyles } from '@material-ui/styles';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles({
  root: {
    marginTop: "1rem",
  },
});

const DescriptionField = ({ description }) => {
  if (description) {
    const classes = useStyles();

    return (
      <Typography variant="caption" className={classes.root}>
        {description}
      </Typography>
    );
  }

  return null;
};

export default DescriptionField;
