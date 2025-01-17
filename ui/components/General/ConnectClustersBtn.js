import React from 'react';
import Link from 'next/link';
import { iconMedium } from '../../css/icons.styles';
import AddIcon from '@material-ui/icons/AddCircleOutline';
import { Button, useTheme } from '@layer5/sistent';

function ConnectClustersBtn() {
  const theme = useTheme();
  return (
    <Link href="/management/connections">
      <Button
        type="submit"
        variant="contained"
        color="primary"
        size="large"
        style={{ margin: '0.5rem 0.5rem', whiteSpace: 'nowrap' }}
      >
        <AddIcon
          style={{
            width: theme.spacing(2.5),
            paddingRight: theme.spacing(0.5),
            ...iconMedium,
          }}
        />
        Connect Clusters
      </Button>
    </Link>
  );
}

export default ConnectClustersBtn;
