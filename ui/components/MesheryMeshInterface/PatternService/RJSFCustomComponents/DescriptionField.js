import React from 'react';


import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core';

const styles = () => ({
  typography : {
    marginTop : "1rem",
    fontSize : "0.8rem",
  },
});

const DescriptionField = ({ description, classes }) => {
  if (description) {

    return (
      <Typography variant="caption" className={classes.typography}>
        {description}
      </Typography>
    );
  }
}

export default withStyles(styles)(DescriptionField);
