import {
  useGetEnvironmentConnectionsQuery,
  useGetEnvironmentsQuery,
} from '@/rtk-query/environments';
import {
  Box,
  Checkbox,
  Stack,
  Typography,
  useTheme,
  styled,
  EnvironmentIcon,
  IconButton,
  CustomTooltip,
} from '@layer5/sistent';
import { Loading, StepHeading } from './common';
import { K8sContextConnectionChip } from '../Header';
import { createContext } from 'react';
import { useContext } from 'react';
import {
  selectIsEnvSelected,
  selectIsConnectionSelected,
  toggleEnvSelection,
  toggleConnection,
  selectSelectedK8sConnections,
} from '@/store/slices/globalEnvironmentContext';
import { useSelectorRtk, useDispatchRtk } from '@/store/hooks';
import { Button } from '@layer5/sistent';
import { AddIcon } from '@layer5/sistent';
import { Edit } from '@material-ui/icons';

export const DeploymentTargetContext = createContext({
  meshsyncControllerState: null,
  connectionMetadataState: null,
  organization: null,
});

const StyledEnvironmentCard = styled(Box)(({ theme }) => ({
  borderRadius: '0.5rem',
  border: `2px solid ${theme.palette.border.normal}`,
  overflow: 'hidden',
}));

const StyledEnvironmentHeader = styled(Box)(({ theme }) => ({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  backgroundColor: theme.palette.background.blur.heavy,
  padding: '0.5rem 1rem',
  borderBottom: `1px solid ${theme.palette.border.strong}`,
  fontWeight: '600',
}));

const K8sContextConnection = ({ connection, environment }) => {
  const { meshsyncControllerState, connectionMetadataState } = useContext(DeploymentTargetContext);
  const isSelected = useSelectorRtk((state) =>
    selectIsConnectionSelected(state, environment.id, connection.id),
  );
  const dispatch = useDispatchRtk();
  const toggleK8sConnection = () => dispatch(toggleConnection(environment, connection));
  return (
    <K8sContextConnectionChip
      ctx={{ ...connection.metadata, connection_id: connection.id }}
      onSelectChange={toggleK8sConnection}
      selected={isSelected}
      selectable
      meshsyncControllerState={meshsyncControllerState}
      connectionMetadataState={connectionMetadataState}
    />
  );
};

const EnvironmentConnections = ({ environment, connections }) => {
  return (
    <Box display="flex" gap={1} flexWrap={'wrap'} justifyContent="flex-start">
      {connections.map((connection) => (
        <K8sContextConnection
          key={connection.id}
          connection={connection}
          environment={environment}
        />
      ))}
    </Box>
  );
};

const EnvironmentCard = ({ environment }) => {
  const theme = useTheme();
  const { data, isLoading } = useGetEnvironmentConnectionsQuery({
    environmentId: environment.id,
  });
  const dispatch = useDispatchRtk();
  const connections =
    data?.connections?.filter((connection) => connection.kind == 'kubernetes') || [];

  const isEnvSelected = useSelectorRtk((state) => selectIsEnvSelected(state, environment.id));
  const selectedConnections = useSelectorRtk((state) =>
    selectSelectedK8sConnections(state, environment.id),
  );
  const selectedConnectionsCount = selectedConnections.length;

  const toggleEnv = () => dispatch(toggleEnvSelection(environment, connections));

  return (
    <StyledEnvironmentCard>
      <StyledEnvironmentHeader>
        <Box gap={1} display="flex" alignItems="center">
          <Checkbox
            data-testid={`env-${environment.id}`}
            checked={isEnvSelected}
            onChange={toggleEnv}
          />
          <Typography variant="textB2SemiBold" color={theme.palette.text.default}>
            {environment.name}
          </Typography>
        </Box>
        <Typography variant="textB2SemiBold" color={theme.palette.icon.default}>
          ({selectedConnectionsCount}/{connections.length})
        </Typography>
      </StyledEnvironmentHeader>
      <Box p={2}>
        {isLoading ? (
          <Loading />
        ) : (
          <EnvironmentConnections environment={environment} connections={connections} />
        )}
      </Box>
    </StyledEnvironmentCard>
  );
};

export const EnvironmentsEmptyState = ({ message, onButtonClick }) => {
  const theme = useTheme();
  return (
    <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center">
      <EnvironmentIcon height={100} width={100} fill={theme.palette.icon.default} />
      <Typography color={theme.palette.text.neutral.default} variant="textB2SemiBold">
        {message || 'No environments found'}
      </Typography>

      <Button
        type="submit"
        variant="contained"
        color="primary"
        onClick={onButtonClick}
        style={{ margin: '0.6rem 0.6rem', whiteSpace: 'nowrap' }}
      >
        <AddIcon fill={theme.palette.background.constant.white} />
        Add Environments
      </Button>
    </Box>
  );
};

export const SelectTargetEnvironments = ({ setIsEnvrionmentModalOpen }) => {
  const organization = useContext(DeploymentTargetContext).organization;
  const { data, isLoading, isError } = useGetEnvironmentsQuery({ orgId: organization.id });
  const environments = data?.environments || [];
  if (isLoading) {
    return <Loading />;
  }

  if (isError) {
    return <Typography variant="textB1Regular">Error fetching environments</Typography>;
  }

  return (
    <Stack gap={2}>
      <Box display={'flex'} justifyContent={'space-between'} alignItems={'center'}>
        <StepHeading>Identify Deployment Targets</StepHeading>
        {environments.length > 0 && (
          <CustomTooltip title="Configure Environments">
            <div>
              <IconButton
                onClick={() => setIsEnvrionmentModalOpen(true)}
                aria-label="edit-environments"
              >
                <Edit />
              </IconButton>
            </div>
          </CustomTooltip>
        )}
      </Box>
      {environments.length === 0 && (
        <EnvironmentsEmptyState
          message="No environments found. Add a new environment."
          onButtonClick={() => setIsEnvrionmentModalOpen(true)}
        />
      )}

      <Stack spacing={2}>
        {environments.map((env) => (
          <EnvironmentCard key={env.id} environment={env} />
        ))}
      </Stack>
    </Stack>
  );
};
