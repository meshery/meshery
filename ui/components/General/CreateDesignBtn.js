import React from 'react';
import { Button, withStyles } from '@material-ui/core';
import Link from 'next/link';
import { iconMedium } from '../../css/icons.styles';
import AddIcon from '@mui/icons-material/AddCircleOutline';

const style = (theme) => ({
  addIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
});

function CreateDesignBtn({ classes }) {
  return (
    <Link href="/configuration/design">
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        style={{ margin: '0.5rem 0.5rem', whiteSpace: 'nowrap' }}
      >
        <AddIcon style={iconMedium} className={classes.addIcon} />
        Create Design
      </Button>
    </Link>
  );
}

export default withStyles(style)(CreateDesignBtn);
