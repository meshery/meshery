import { AddIcon, Button, styled, Typography, useTheme } from '@sistent/sistent';
import OperatorLight from '../../../assets/img/OperatorLight';
import Operator from '../../../assets/img/Operator';
import { CONNECTION_KINDS } from '@/utils/Enum';
import { useConnectionWizardModal } from '@/utils/context/ConnectionWizardContextProvider';

const TextContent = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  marginTop: '1rem',
  padding: '10px',
  borderRadius: '10px',
});

const StyledAddIcon = styled(AddIcon)(({ theme }) => ({
  width: theme.spacing(2.5),
  marginRight: theme.spacing(0.5),
}));

/**
 * Kubernetes empty state — "Connect Clusters" opens the Create Connection
 * wizard in place with Kubernetes pre-selected (no navigate away).
 */
export const K8sEmptyState = ({ message }) => {
  const theme = useTheme();
  const { openCreateConnection } = useConnectionWizardModal();

  const handleClick = () => {
    openCreateConnection({
      kind: CONNECTION_KINDS.KUBERNETES,
      skipKindSelection: true,
    });
  };

  return (
    <TextContent>
      {theme.palette.mode === 'dark' ? <OperatorLight /> : <Operator />}
      <Typography variant="h5">{message || 'No cluster connected yet'}</Typography>

      <Button
        type="button"
        variant="contained"
        color="primary"
        onClick={handleClick}
        sx={{ margin: '0.6rem 0.6rem', whiteSpace: 'nowrap' }}
      >
        <StyledAddIcon />
        Connect Clusters
      </Button>
    </TextContent>
  );
};
