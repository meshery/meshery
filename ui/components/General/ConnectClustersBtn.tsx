import Link from 'next/link';
import { iconMedium } from '../../css/icons.styles';
import { AddCircleIcon as AddIcon, Button, useTheme } from '@sistent/sistent';

const buttonSx = { m: 1, whiteSpace: 'nowrap' } as const;

function ConnectClustersBtn() {
  const theme = useTheme();
  return (
    <Link href="/management/connections">
      <Button type="submit" variant="contained" color="primary" size="large" sx={buttonSx}>
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
