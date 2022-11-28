import React from 'react';


import { makeStyles } from '@material-ui/styles';
import Typography from '@material-ui/core/Typography';

const useStyles = makeStyles({
  typography : {
    marginTop : "1rem",
    fontSize : "0.8rem",
  },
});

const DescriptionField = ({ description }) => {
  if (description) {
    const classes = useStyles();

    return (
      <Typography variant="caption" className={classes.typography}>
        {description}
      </Typography>
    );
  }

  return null;
};

export default DescriptionField;
