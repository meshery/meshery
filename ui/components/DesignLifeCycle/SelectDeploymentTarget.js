import {
  useGetEnvironmentConnectionsQuery,
  useGetEnvironmentsQuery,
} from '@/rtk-query/environments';
import { Box, Checkbox, Stack, Typography, useTheme, styled } from '@layer5/sistent';
import { Loading } from './common';
import { K8sContextConnectionChip } from '../Header';
import { createContext } from 'react';
import { useContext } from 'react';
import {
  selectIsEnvSelected,
  selectIsK8sConnectionSelected,
  toggleEnvSelection,
  toggleK8sConnection,
} from '@/store/slices/globalEnvironmentContext';
import { useSelectorRtk, useDispatchRtk } from '@/store/hooks';

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
    selectIsK8sConnectionSelected(state, environment.id, connection.id),
  );
  const dispatch = useDispatchRtk();
  const toggleConnection = () => dispatch(toggleK8sConnection(environment.id, connection.id));
  return (
    <K8sContextConnectionChip
      ctx={connection}
      onSelectChange={toggleConnection}
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
  const toggleEnv = () => dispatch(toggleEnvSelection(environment));

  console.log('theme --->', theme);

  return (
    <StyledEnvironmentCard>
      <StyledEnvironmentHeader>
        <Box gap={1} display="flex" alignItems="center">
          <Checkbox checked={isEnvSelected} onChange={toggleEnv} />
          <Typography variant="textB2SemiBold" color={theme.palette.text.default}>
            {environment.name}
          </Typography>
        </Box>
        <Typography variant="textB2SemiBold" color={theme.palette.icon.default}>
          ({connections.length})
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

export const SelectTargetEnvironments = () => {
  const organization = useContext(DeploymentTargetContext).organization;
  const { data, isLoading, isError } = useGetEnvironmentsQuery({ orgId: organization.id });
  const environments = data?.environments || [];
  const theme = useTheme();
  console.log('environments', theme, environments);
  if (isLoading) {
    return <Loading />;
  }

  if (isError) {
    return <Typography variant="textH3Medium">Error fetching environments</Typography>;
  }

  return (
    <Stack gap={2}>
      <Typography variant="textH3Medium">Identify Deployment Targets</Typography>
      <Stack spacing={2}>
        {environments.map((env) => (
          <EnvironmentCard key={env.id} environment={env} />
        ))}
      </Stack>
    </Stack>
  );
};
