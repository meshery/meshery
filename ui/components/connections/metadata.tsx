import React, { useMemo, useEffect } from 'react';
import {
  Grid2,
  List,
  ListItem,
  ListItemText,
  Box,
  Typography,
  styled,
  useTheme,
} from '@sistent/sistent';
import { useGetControllerDiagnosticsQuery } from '@/rtk-query/connection';

import {
  FormatId,
  FormatStructuredData,
  FormattedDate,
  KeyValue,
  Link,
  createColumnUiSchema,
} from '../data-formatter';
import useKubernetesHook, {
  useControllerStatus,
  useMesheryOperator,
  useMeshsSyncController,
  useNatsController,
} from '@/utils/hooks/useKubernetesHook';
import { TooltipWrappedConnectionChip } from './ConnectionChip';
import { CONTROLLER_STATES, MESHSYNC_DEPLOYMENT_TYPE } from '../../utils/Enum';
import { formatToTitleCase } from '../../utils/utils';

import { ColumnWrapper, ContentContainer, OperationButton, FormatterWrapper } from './styles';

const DISABLED = 'DISABLED';
const KUBERNETES = 'kubernetes';
const MESHERY = 'meshery';

const customIdFormatter = (title, id) => (
  <FormatterWrapper>
    <KeyValue Key={title} Value={<FormatId id={id} />} />
  </FormatterWrapper>
);
const customDateFormatter = (title, date) => (
  <FormatterWrapper>
    <KeyValue Key={title} Value={<FormattedDate date={date} />} />
  </FormatterWrapper>
);

const DefaultPropertyFormatters = {
  id: (value) => customIdFormatter('Id', value),
  uid: (value) => customIdFormatter('Uid', value),
  server_id: (value) => customIdFormatter('Server Id', value),
  created_at: (value) => customDateFormatter('Created At', value),
  updated_at: (value) => customDateFormatter('Updated At', value),
  creation_timestamp: (value) => customDateFormatter('Creation Timestamp', value),
  creationTimestamp: (value) => customDateFormatter('Creation Timestamp', value),
  last_seen: (value) => customDateFormatter('Last Seen', value),
  last_reconciled: (value) => customDateFormatter('Last Reconciled', value),
  last_applied: (value) => customDateFormatter('Last Applied', value),
  last_updated: (value) => customDateFormatter('Last Updated', value),
};

const StyledListItemText = styled(ListItemText)(({ theme }) => ({
  '& .MuiTypography-root.MuiTypography-body2': {
    color: theme.palette.text.tertiary, // Use the secondary color from the theme
  },
}));

const DIAGNOSTIC_SEVERITY_PALETTE = {
  error: 'error',
  warning: 'warning',
  info: 'info',
};

// A single diagnostic: severity-colored card with an explanation, optional
// endpoint, and an ordered list of remediation steps.
const DiagnosticCard = ({ diagnostic }) => {
  const theme = useTheme();
  const paletteKey = DIAGNOSTIC_SEVERITY_PALETTE[diagnostic.severity] || 'info';
  const accent = theme.palette[paletteKey]?.main || theme.palette.text.secondary;

  return (
    <Box
      sx={{
        borderLeft: `3px solid ${accent}`,
        borderRadius: '4px',
        padding: '0.5rem 0.75rem',
        marginBottom: '0.5rem',
        backgroundColor: theme.palette.background.default,
      }}
    >
      <Typography variant="body1" sx={{ fontWeight: 600, color: accent }}>
        {diagnostic.summary}
      </Typography>
      {diagnostic.description && (
        <Typography
          variant="body2"
          sx={{ marginTop: '0.25rem', color: theme.palette.text.tertiary }}
        >
          {diagnostic.description}
        </Typography>
      )}
      {diagnostic.endpoint && (
        <Typography
          variant="body2"
          sx={{ marginTop: '0.25rem', color: theme.palette.text.tertiary }}
        >
          Endpoint:{' '}
          <Box component="code" sx={{ fontFamily: 'monospace' }}>
            {diagnostic.endpoint}
          </Box>
        </Typography>
      )}
      {Array.isArray(diagnostic.remediation) && diagnostic.remediation.length > 0 && (
        <>
          <Typography variant="body2" sx={{ marginTop: '0.5rem', fontWeight: 600 }}>
            Suggested remediation
          </Typography>
          <Box component="ol" sx={{ margin: '0.25rem 0 0', paddingInlineStart: '1.25rem' }}>
            {diagnostic.remediation.map((step, idx) => (
              <li key={idx}>
                <Typography variant="body2" sx={{ color: theme.palette.text.tertiary }}>
                  {step}
                </Typography>
              </li>
            ))}
          </Box>
        </>
      )}
    </Box>
  );
};

// ControllerDiagnosticsSection fetches human-actionable diagnostics for a
// connection's controllers on demand (separate from the live status stream) and
// refetches whenever the connection's live controller status changes, so the
// section stays in sync without bloating the status SSE payload.
const ControllerDiagnosticsSection = ({ connectionId, statusKey }) => {
  const theme = useTheme();
  const { data, isFetching, refetch } = useGetControllerDiagnosticsQuery(connectionId);

  useEffect(() => {
    // statusKey changes when any of this connection's controller states change,
    // so this refetches the diagnostics to stay in sync with the live status.
    if (connectionId) {
      refetch();
    }
  }, [statusKey, connectionId, refetch]);

  if (!connectionId) {
    return null;
  }

  const diagnostics = data?.diagnostics ?? [];

  return (
    <Grid2 size={{ xs: 12 }}>
      <ContentContainer container spacing={1} size="grow">
        <Grid2 size={{ xs: 12 }}>
          <Typography variant="body1" sx={{ fontWeight: 'bold', marginBottom: '0.5rem' }}>
            Diagnostics
          </Typography>
          {isFetching && !data ? (
            <Typography variant="body2" sx={{ color: theme.palette.text.tertiary }}>
              Checking controller health…
            </Typography>
          ) : diagnostics.length === 0 ? (
            <Typography variant="body2" sx={{ color: theme.palette.text.tertiary }}>
              No issues detected for this connection&apos;s controllers.
            </Typography>
          ) : (
            diagnostics.map((diagnostic) => (
              <DiagnosticCard key={diagnostic.code} diagnostic={diagnostic} />
            ))
          )}
        </Grid2>
      </ContentContainer>
    </Grid2>
  );
};

const KubernetesMetadataFormatter = ({ meshsyncControllerState, connection, metadata }) => {
  const pingKubernetes = useKubernetesHook();
  const { ping: pingMesheryOperator } = useMesheryOperator();
  const { ping: pingMeshSync } = useMeshsSyncController();
  const { ping: pingNats } = useNatsController();
  const { getControllerStatesByConnectionID } = useControllerStatus(meshsyncControllerState);

  const handleKubernetesClick = () => {
    pingKubernetes(metadata.name, metadata.server, connection.id);
  };

  const handleNATSClick = () => {
    pingNats({ connectionID: connection.id });
  };

  const handleOperatorClick = () => {
    pingMesheryOperator({ connectionID: connection.id });
  };

  const handleMeshSyncClick = () => {
    pingMeshSync({ connectionID: connection.id });
  };

  const { operatorState, meshSyncState, natsState, operatorVersion, meshSyncVersion, natsVersion } =
    getControllerStatesByConnectionID(connection.id);

  const meshsyncDeploymentMode =
    metadata?.meshsyncDeploymentMode ?? metadata?.meshsync_deployment_mode;
  const isEmbeddedMode = meshsyncDeploymentMode === MESHSYNC_DEPLOYMENT_TYPE.EMBEDDED;

  return (
    <Grid2 container spacing={1} sx={{ textTransform: 'none' }} size="grow">
      <Grid2 size={{ xs: 12, md: 6 }}>
        <ColumnWrapper>
          <Grid2 container spacing={1} size="grow">
            <OperationButton size={{ xs: 12, md: 5 }}>
              <List>
                <ListItem>
                  <TooltipWrappedConnectionChip
                    tooltip={`Server: ${metadata.server}`}
                    title={metadata.name}
                    status={connection.status}
                    iconSrc={'/static/img/integrations/kubernetes.svg'}
                    handlePing={handleKubernetesClick}
                  />
                </ListItem>
              </List>
            </OperationButton>
          </Grid2>
          <ContentContainer container spacing={1} size="grow">
            <Grid2 size={{ xs: 12, md: 5 }}>
              <List>
                <ListItem>
                  <StyledListItemText primary="Name" secondary={metadata.name} />
                </ListItem>
                <ListItem>
                  <StyledListItemText primary="K8s Version" secondary={metadata.version} />
                </ListItem>
              </List>
            </Grid2>
            <Grid2 size={{ xs: 12, md: 5 }}>
              <List>
                <ListItem>
                  <StyledListItemText
                    primary="Created At"
                    secondary={<FormattedDate date={connection.createdAt} />}
                  />
                </ListItem>
                <ListItem>
                  <StyledListItemText
                    primary="Updated At"
                    secondary={<FormattedDate date={connection.updatedAt} />}
                  />
                </ListItem>
              </List>
            </Grid2>
            <Grid2 size={{ xs: 12, md: 5 }}>
              <List>
                <ListItem>
                  <StyledListItemText
                    style={{
                      width: '80%',
                      wordWrap: 'break-word',
                    }}
                    primary="Server"
                    secondary={<Link title={metadata.server}>{metadata.server}</Link>}
                  />
                </ListItem>
              </List>
            </Grid2>
          </ContentContainer>
        </ColumnWrapper>
      </Grid2>
      <Grid2 size={{ xs: 12, md: 6 }}>
        <ColumnWrapper>
          {!isEmbeddedMode && (
            <Grid2 container spacing={1} size="grow">
              <OperationButton size={{ xs: 12, md: 4 }}>
                <List>
                  <ListItem>
                    <TooltipWrappedConnectionChip
                      tooltip={operatorState ? `Version: ${operatorVersion}` : 'Not Available'}
                      title={'Operator'}
                      disabled={operatorState === CONTROLLER_STATES.UNDEPLOYED}
                      status={operatorState}
                      handlePing={handleOperatorClick}
                      iconSrc="/static/img/integrations/meshery-operator.svg"
                      width="9rem"
                    />
                  </ListItem>
                </List>
              </OperationButton>

              {(meshSyncState || natsState) && (
                <>
                  <Grid2 size={{ xs: 12, md: 4 }}>
                    <List>
                      <ListItem>
                        <TooltipWrappedConnectionChip
                          tooltip={meshSyncState !== DISABLED ? `Ping MeshSync` : 'Not Available'}
                          title={'MeshSync'}
                          status={meshSyncState}
                          handlePing={handleMeshSyncClick}
                          iconSrc="/static/img/extensions/meshsync.svg"
                          width="9rem"
                        />
                      </ListItem>
                    </List>
                  </Grid2>
                  <Grid2 size={{ xs: 12, md: 4 }}>
                    <List>
                      <ListItem>
                        <TooltipWrappedConnectionChip
                          tooltip={natsState === 'Not Active' ? 'Not Available' : `Reconnect NATS`}
                          title={'BROKER'}
                          status={natsState}
                          handlePing={handleNATSClick}
                          iconSrc="/static/img/integrations/nats-icon-color.svg"
                          width="9rem"
                        />
                      </ListItem>
                    </List>
                  </Grid2>
                </>
              )}
            </Grid2>
          )}
          <ContentContainer container spacing={1} size="grow">
            {!isEmbeddedMode && (
              <>
                <Grid2 size={{ xs: 12, md: 5 }}>
                  <List>
                    <ListItem>
                      <StyledListItemText
                        primary="Operator State"
                        secondary={formatToTitleCase(operatorState)}
                      />
                    </ListItem>
                    <ListItem>
                      <StyledListItemText primary="Operator Version" secondary={operatorVersion} />
                    </ListItem>
                  </List>
                </Grid2>
                <Grid2 size={{ xs: 12, md: 5 }}>
                  <List>
                    <ListItem>
                      <StyledListItemText
                        primary="MeshSync State"
                        secondary={formatToTitleCase(meshSyncState) || 'Undeployed'}
                      />
                    </ListItem>
                    <ListItem>
                      <StyledListItemText primary="MeshSync Version" secondary={meshSyncVersion} />
                    </ListItem>
                  </List>
                </Grid2>
                <Grid2 size={{ xs: 12, md: 5 }}>
                  <List>
                    <ListItem>
                      <StyledListItemText
                        primary="Broker State"
                        secondary={formatToTitleCase(natsState) || 'Not Connected'}
                      />
                    </ListItem>
                    <ListItem>
                      <StyledListItemText primary="Broker Version" secondary={natsVersion} />
                    </ListItem>
                  </List>
                </Grid2>
              </>
            )}
            <Grid2 size={{ xs: 12, md: 8 }}>
              <List>
                <ListItem>
                  <StyledListItemText
                    primary="Deployment Mode"
                    secondary={formatToTitleCase(meshsyncDeploymentMode || 'N/A')}
                  />
                </ListItem>
              </List>
            </Grid2>
          </ContentContainer>
        </ColumnWrapper>
      </Grid2>
      <ControllerDiagnosticsSection
        connectionId={connection.id}
        statusKey={`${operatorState}|${meshSyncState}|${natsState}`}
      />
    </Grid2>
  );
};

const MesheryMetadataFormatter = ({ connection }) => {
  const uiSchema = useMemo(
    () =>
      createColumnUiSchema({
        metadata: connection.metadata || {},
        numCols: {
          xs: 2,
          md: 4,
        },
      }),
    [connection.metadata],
  );

  return (
    <FormatStructuredData
      data={connection.metadata}
      uiSchema={uiSchema}
      propertyFormatters={DefaultPropertyFormatters}
    />
  );
};

export const MeshSyncDataFormatter = ({ metadata }) => {
  const theme = useTheme();
  const uiSchema = useMemo(
    () =>
      createColumnUiSchema({
        metadata,
        numCols: {
          xs: 3,
          md: 5,
        },
      }),
    [metadata],
  );

  return (
    <Box sx={{ backgroundColor: theme.palette.background.card, width: '100%', padding: '1rem' }}>
      <FormatStructuredData
        data={metadata}
        uiSchema={uiSchema}
        propertyFormatters={DefaultPropertyFormatters}
      />
    </Box>
  );
};

const FormatConnectionMetadata = (props) => {
  const theme = useTheme();
  const { connection, meshsyncControllerState } = props;
  let formatter;

  switch (connection.kind) {
    case KUBERNETES:
      formatter = (
        <KubernetesMetadataFormatter
          meshsyncControllerState={meshsyncControllerState}
          connection={connection}
          metadata={connection.metadata}
        />
      );
      break;
    case MESHERY:
      formatter = <MesheryMetadataFormatter connection={connection} />;
      break;
    default:
      formatter = (
        <FormatStructuredData
          data={connection.metadata}
          propertyFormatters={DefaultPropertyFormatters}
        />
      );
  }

  return (
    <Box sx={{ backgroundColor: theme.palette.background.card, padding: '1rem' }}>{formatter}</Box>
  );
};

export default FormatConnectionMetadata;
