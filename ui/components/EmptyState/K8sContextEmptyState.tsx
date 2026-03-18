import { Add, Button, Link, Typography, styled, useTheme } from '@sistent/sistent';
import OperatorLight from '../../assets/img/OperatorLight';
import Operator from '../../assets/img/Operator';

const TextContent = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  marginTop: '1rem',
  padding: '10px',
  borderRadius: '10px',
});

const StyledIcon = styled(Add)(({ theme }) => ({
  width: theme.spacing(2.5),
  marginRight: theme.spacing(0.5),
}));

export const K8sEmptyState = ({ message }) => {
  const theme = useTheme();
  return (
    <TextContent>
      {theme.palette.mode === 'dark' ? <OperatorLight /> : <Operator />}
      <Typography variant="h5">{message || 'No cluster connected yet'}</Typography>

      <Link href="/management/connections">
        <Button
          type="submit"
          variant="contained"
          color="primary"
          sx={{ margin: '0.6rem 0.6rem', whiteSpace: 'nowrap' }}
        >
          <StyledIcon />
          Connect Clusters
        </Button>
      </Link>
    </TextContent>
  );
};
