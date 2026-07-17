/* eslint-disable max-lines */
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
import telemetryPrometheusApi, {
  useLazyPingPrometheusConnectionQuery,
} from '@/rtk-query/telemetryPrometheus';
import telemetryGrafanaApi, {
  useLazyPingGrafanaConnectionQuery,
} from '@/rtk-query/telemetryGrafana';

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
import { CONNECTION_STATES, CONTROLLER_STATES, MESHSYNC_DEPLOYMENT_TYPE } from '../../utils/Enum';
import { formatToTitleCase } from '../../utils/utils';

import { ColumnWrapper, ContentContainer, OperationButton, FormatterWrapper } from './styles';

const DISABLED = 'DISABLED';
const KUBERNETES = 'kubernetes';
const MESHERY = 'meshery';
const PROMETHEUS = 'prometheus';
const GRAFANA = 'grafana';
const GITHUB = 'github';

// Canonical connection.Metadata keys (server contracts). Prefer camelCase;
// tolerate legacy snake_case for rows that predate the identifier-naming flip.
const META = {
  URL: 'url',
  NAME: 'name',
  // Meshery platform connection - BuildMesheryConnectionPayload
  SERVER_ID: 'serverId',
  SERVER_VERSION: 'serverVersion',
  SERVER_BUILD_SHA: 'serverBuildSha',
  SERVER_LOCATION: 'serverLocation',
  // Telemetry - models/telemetry/{prometheus,grafana}
  PROM_PANELS: 'telemetryPrometheusPanels',
  GRAFANA_BOARDS: 'telemetryPinnedBoards',
  // GitHub App (Layer5 Cloud) - installation + design snapshot paths
  INSTALLATION_ID: 'installationId',
  SNAPSHOT_PATHS: 'snapshotPaths',
};

const KIND_ICONS = {
  [PROMETHEUS]: '/static/img/integrations/prometheus_logo_orange_circle.svg',
  [GRAFANA]: '/static/img/integrations/grafana_icon.svg',
  [GITHUB]: '/static/img/extensions/github.svg',
  [MESHERY]: '/static/img/meshery-logo/meshery-logo.png',
  artifacthub: '/static/img/integrations/artifacthub.svg',
};

const StyledListItemText = styled(ListItemText)(({ theme }) => ({
  '& .MuiTypography-root.MuiTypography-body2': {
    color: theme.palette.text.tertiary,
  },
}));

/** First non-nullish value among metadata keys (canonical then legacy). */
const readMeta = (metadata, ...keys) => {
  if (!metadata) return undefined;
  for (const key of keys) {
    const value = metadata[key];
    if (value !== undefined && value !== null && value !== '') {
      return value;
    }
  }
  return undefined;
};

const countArray = (value) => (Array.isArray(value) ? value.length : 0);

// Schemas wire camelCase (createdAt); some table/legacy rows still use snake_case.
const connectionTimestamp = (connection, camelKey, snakeKey) =>
  connection?.[camelKey] ?? connection?.[snakeKey];

const UrlLink = ({ url }) => (url ? <Link title={url} href={url} /> : 'N/A');

const DetailListItem = ({ primary, secondary }) => (
  <ListItem>
    <StyledListItemText primary={primary} secondary={secondary ?? 'N/A'} />
  </ListItem>
);

/**
 * Shared left-panel layout used by non-Kubernetes connection detail views:
 * ping/status chip on top, then two equal detail columns.
 */
const ConnectionDetailPanel = ({ chip, leftItems, rightItems, sidePanel }) => {
  const theme = useTheme();
  return (
    <Grid2 container spacing={1} sx={{ textTransform: 'none' }} size="grow">
      <Grid2 size={{ xs: 12, md: sidePanel ? 6 : 12 }}>
        <ColumnWrapper>
          <Grid2 container spacing={1} size="grow">
            <OperationButton size={{ xs: 12, md: 5 }}>
              <List>
                <ListItem>{chip}</ListItem>
              </List>
            </OperationButton>
          </Grid2>
          <ContentContainer container spacing={1} size="grow">
            <Grid2 size={{ xs: 12, md: 5 }}>
              <List>{leftItems}</List>
            </Grid2>
            <Grid2 size={{ xs: 12, md: 5 }}>
              <List>{rightItems}</List>
            </Grid2>
          </ContentContainer>
        </ColumnWrapper>
      </Grid2>
      {sidePanel && (
        <Grid2
          size={{ xs: 12, md: 6 }}
          sx={{
            padding: '1rem',
            borderLeft: `1px solid ${theme.palette.border?.default || theme.palette.divider}`,
          }}
        >
          {sidePanel}
        </Grid2>
      )}
    </Grid2>
  );
};

/**
 * Prometheus / Grafana detail view (presentational).
 *
 * Parents call the RTK hooks at their top level and pass `triggerPing` +
 * `pingState` down - calling hooks via props would violate Rules of Hooks and
 * trip react-hooks/rules-of-hooks. Version still comes from shared RTK cache
 * (useQueryState in the parent) so a table name-column chip ping refreshes it.
 *
 * Panel/board counts come from the canonical metadata keys written by the
 * telemetry handlers; version appears only after any successful chip ping.
 */
const TelemetryMetadataFormatter = ({
  connection,
  metadata,
  productName,
  iconSrc,
  countLabel,
  countValue,
  triggerPing,
  pingState,
  diagnosticCode,
}) => {
  const connectionID = connection.id;
  const { data, isError, isFetching, isUninitialized, isSuccess } = pingState;

  const handlePing = () => {
    if (connectionID) {
      triggerPing({ connectionID });
    }
  };

  const url = readMeta(metadata, META.URL);
  const displayName = connection.name || readMeta(metadata, META.NAME) || productName;
  let version = '-';
  if (isFetching) {
    version = 'Checking…';
  } else if (!isUninitialized && isError) {
    version = 'Unavailable';
  } else if (data?.version) {
    version = data.version;
  } else if (isSuccess || !isUninitialized) {
    version = 'Unknown';
  }
  // Only override chip status after a failed chip-initiated ping.
  const displayedStatus =
    !isUninitialized && isError ? CONNECTION_STATES.DISCONNECTED : connection.status;

  const diagnostics = [];
  if (!isUninitialized && isError) {
    diagnostics.push({
      code: diagnosticCode,
      severity: 'error',
      summary: `${productName} Unreachable`,
      description: `Meshery could not communicate with the configured ${productName} URL.`,
      endpoint: url,
      remediation: [
        'Ensure the URL is correct and includes the scheme (http:// or https://).',
        `Verify ${productName} is running and reachable from the Meshery Server network.`,
        'Confirm the credential (API token / basic auth) is still valid.',
      ],
    });
  }

  const sidePanel =
    diagnostics.length > 0 ? (
      <>
        <Typography variant="body1" sx={{ fontWeight: 'bold', marginBottom: '0.5rem' }}>
          Diagnostics
        </Typography>
        {diagnostics.map((diagnostic) => (
          <DiagnosticCard key={diagnostic.code} diagnostic={diagnostic} />
        ))}
      </>
    ) : null;

  return (
    <ConnectionDetailPanel
      chip={
        <TooltipWrappedConnectionChip
          tooltip={url ? `Server: ${url}` : productName}
          title={displayName}
          status={displayedStatus}
          iconSrc={connection.kindLogo || iconSrc}
          handlePing={connectionID ? handlePing : undefined}
        />
      }
      leftItems={
        <>
          <DetailListItem primary="Name" secondary={displayName} />
          <DetailListItem
            primary="Created At"
            secondary={
              <FormattedDate date={connectionTimestamp(connection, 'createdAt', 'created_at')} />
            }
          />
          <DetailListItem
            primary="Updated At"
            secondary={
              <FormattedDate date={connectionTimestamp(connection, 'updatedAt', 'updated_at')} />
            }
          />
        </>
      }
      rightItems={
        <>
          <ListItem>
            <StyledListItemText
              style={{ width: '80%', wordWrap: 'break-word' }}
              primary="Server"
              secondary={<UrlLink url={url} />}
            />
          </ListItem>
          <DetailListItem primary="Version" secondary={version} />
          <DetailListItem primary={countLabel} secondary={countValue} />
        </>
      }
      sidePanel={sidePanel}
    />
  );
};

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
                    secondary={
                      <FormattedDate
                        date={connectionTimestamp(connection, 'createdAt', 'created_at')}
                      />
                    }
                  />
                </ListItem>
                <ListItem>
                  <StyledListItemText
                    primary="Updated At"
                    secondary={
                      <FormattedDate
                        date={connectionTimestamp(connection, 'updatedAt', 'updated_at')}
                      />
                    }
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
                    secondary={<UrlLink url={metadata?.server} />}
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

const usePrometheusPingQueryState = (arg, opts) =>
  telemetryPrometheusApi.endpoints.pingPrometheusConnection.useQueryState(arg, opts);

const useGrafanaPingQueryState = (arg, opts) =>
  telemetryGrafanaApi.endpoints.pingGrafanaConnection.useQueryState(arg, opts);

const PrometheusMetadataFormatter = ({ connection, metadata }) => {
  const connectionID = connection.id;
  const [triggerPing] = useLazyPingPrometheusConnectionQuery();
  const pingState = usePrometheusPingQueryState({ connectionID }, { skip: !connectionID });

  return (
    <TelemetryMetadataFormatter
      connection={connection}
      metadata={metadata}
      productName="Prometheus"
      iconSrc={KIND_ICONS[PROMETHEUS]}
      countLabel="Saved Panels"
      countValue={countArray(readMeta(metadata, META.PROM_PANELS))}
      triggerPing={triggerPing}
      pingState={pingState}
      diagnosticCode="prometheus-unreachable"
    />
  );
};

const GrafanaMetadataFormatter = ({ connection, metadata }) => {
  const connectionID = connection.id;
  const [triggerPing] = useLazyPingGrafanaConnectionQuery();
  const pingState = useGrafanaPingQueryState({ connectionID }, { skip: !connectionID });

  return (
    <TelemetryMetadataFormatter
      connection={connection}
      metadata={metadata}
      productName="Grafana"
      iconSrc={KIND_ICONS[GRAFANA]}
      countLabel="Pinned Boards"
      countValue={countArray(readMeta(metadata, META.GRAFANA_BOARDS))}
      triggerPing={triggerPing}
      pingState={pingState}
      diagnosticCode="grafana-unreachable"
    />
  );
};

// GitHub App connections (Layer5 Cloud): metadata carries installationId +
// snapshotPaths (design snapshot repo paths). Not owner/host/path design-file fields.
const GithubMetadataFormatter = ({ connection, metadata }) => {
  const installationId = readMeta(metadata, META.INSTALLATION_ID, 'installation_id');
  const snapshotPaths = readMeta(
    metadata,
    META.SNAPSHOT_PATHS,
    'snapshot_paths',
    'repositoryPaths',
  );
  const snapshotCount = countArray(snapshotPaths);
  const displayName = connection.name || 'GitHub App';
  const connectionType = [connection.type, connection.subType].filter(Boolean).join(' / ') || 'N/A';

  return (
    <ConnectionDetailPanel
      chip={
        <TooltipWrappedConnectionChip
          tooltip={
            installationId ? `GitHub App installation ${installationId}` : 'GitHub App connection'
          }
          title={displayName}
          status={connection.status}
          iconSrc={connection.kindLogo || KIND_ICONS[GITHUB]}
          handlePing={undefined}
        />
      }
      leftItems={
        <>
          <DetailListItem primary="Name" secondary={displayName} />
          <DetailListItem
            primary="Installation ID"
            secondary={installationId ? <FormatId id={String(installationId)} /> : 'N/A'}
          />
          <DetailListItem primary="Type" secondary={connectionType} />
        </>
      }
      rightItems={
        <>
          <DetailListItem
            primary="Snapshot Paths"
            secondary={snapshotCount > 0 ? `${snapshotCount} configured` : 'None configured'}
          />
          <DetailListItem
            primary="Created At"
            secondary={
              <FormattedDate date={connectionTimestamp(connection, 'createdAt', 'created_at')} />
            }
          />
          <DetailListItem
            primary="Updated At"
            secondary={
              <FormattedDate date={connectionTimestamp(connection, 'updatedAt', 'updated_at')} />
            }
          />
        </>
      }
    />
  );
};

// Meshery platform connection - BuildMesheryConnectionPayload camelCase keys.
const MesheryMetadataFormatter = ({ connection, metadata }) => {
  const serverLocation = readMeta(metadata, META.SERVER_LOCATION, 'server_location');
  const serverVersion = readMeta(metadata, META.SERVER_VERSION, 'server_version');
  const serverBuildSha = readMeta(metadata, META.SERVER_BUILD_SHA, 'server_build_sha');
  const serverId = readMeta(metadata, META.SERVER_ID, 'server_id');
  const displayName = connection.name || 'Meshery Server';
  const shortSha =
    typeof serverBuildSha === 'string' && serverBuildSha.length > 0
      ? serverBuildSha.substring(0, 7)
      : null;

  return (
    <ConnectionDetailPanel
      chip={
        <TooltipWrappedConnectionChip
          tooltip={serverLocation ? `Location: ${serverLocation}` : 'Meshery Server'}
          title={displayName}
          status={connection.status}
          iconSrc={connection.kindLogo || KIND_ICONS[MESHERY]}
          handlePing={undefined}
        />
      }
      leftItems={
        <>
          <DetailListItem primary="Server Name" secondary={displayName} />
          <DetailListItem primary="Server Version" secondary={serverVersion || 'N/A'} />
          <ListItem>
            <StyledListItemText
              style={{ width: '80%', wordWrap: 'break-word' }}
              primary="Server Location"
              secondary={<UrlLink url={serverLocation} />}
            />
          </ListItem>
        </>
      }
      rightItems={
        <>
          <DetailListItem
            primary="Server ID"
            secondary={serverId ? <FormatId id={String(serverId)} /> : 'N/A'}
          />
          <DetailListItem
            primary="Build SHA"
            secondary={
              shortSha ? (
                <Link
                  title={shortSha}
                  href={`https://github.com/meshery/meshery/commit/${serverBuildSha}`}
                />
              ) : (
                'N/A'
              )
            }
          />
          <DetailListItem
            primary="Created At"
            secondary={
              <FormattedDate date={connectionTimestamp(connection, 'createdAt', 'created_at')} />
            }
          />
          <DetailListItem
            primary="Updated At"
            secondary={
              <FormattedDate date={connectionTimestamp(connection, 'updatedAt', 'updated_at')} />
            }
          />
        </>
      }
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
      formatter = (
        <MesheryMetadataFormatter connection={connection} metadata={connection.metadata} />
      );
      break;
    case PROMETHEUS:
      formatter = (
        <PrometheusMetadataFormatter connection={connection} metadata={connection.metadata} />
      );
      break;
    case GRAFANA:
      formatter = (
        <GrafanaMetadataFormatter connection={connection} metadata={connection.metadata} />
      );
      break;
    case GITHUB:
      formatter = (
        <GithubMetadataFormatter connection={connection} metadata={connection.metadata} />
      );
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
