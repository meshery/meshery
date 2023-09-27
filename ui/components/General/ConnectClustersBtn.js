import React from 'react';
import { Button, withStyles } from '@material-ui/core';
import Link from 'next/link';
import { iconMedium } from '../../css/icons.styles';
import AddIcon from '@material-ui/icons/AddCircleOutline';

const style = (theme) => ({
  addIcon: {
    width: theme.spacing(2.5),
    paddingRight: theme.spacing(0.5),
  },
});

function ConnectClustersBtn({ classes }) {
  return (
    <Link href="/settings">
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        style={{ margin: '0.5rem 0.5rem', whiteSpace: 'nowrap' }}
      >
        <AddIcon style={iconMedium} className={classes.addIcon} />
        Connect Clusters
      </Button>
    </Link>
  );
}

export default withStyles(style)(ConnectClustersBtn);
