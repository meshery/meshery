import React from 'react';
import { Grid2, List, ListItem, ListItemText, Box, styled, useTheme } from '@layer5/sistent';

import {
  FormatId,
  FormatStructuredData,
  FormattedDate,
  KeyValue,
  Link,
  createColumnUiSchema,
} from '../DataFormatter';
import useKubernetesHook, {
  useControllerStatus,
  useMesheryOperator,
  useMeshsSyncController,
  useNatsController,
} from '../hooks/useKubernetesHook';
import { TooltipWrappedConnectionChip } from './ConnectionChip';
import { CONTROLLER_STATES } from '../../utils/Enum';
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
                    iconSrc={'/static/img/kubernetes.svg'}
                    handlePing={() => handleKubernetesClick(connection.id)}
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
                    secondary={<FormattedDate date={connection.created_at} />}
                  />
                </ListItem>
                <ListItem>
                  <StyledListItemText
                    primary="Updated At"
                    secondary={<FormattedDate date={connection.updated_at} />}
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
                    iconSrc="/static/img/meshery-operator.svg"
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
                        status={meshSyncState?.toLowerCase()}
                        handlePing={handleMeshSyncClick}
                        iconSrc="/static/img/meshsync.svg"
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
                        title={'NATS'}
                        status={natsState?.toLowerCase()}
                        handlePing={() => handleNATSClick()}
                        iconSrc="/static/img/nats-icon-color.svg"
                        width="9rem"
                      />
                    </ListItem>
                  </List>
                </Grid2>
              </>
            )}
          </Grid2>

          <ContentContainer container spacing={1} size="grow">
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
                    primary="NATS State"
                    secondary={formatToTitleCase(natsState) || 'Not Connected'}
                  />
                </ListItem>
                <ListItem>
                  <StyledListItemText primary="NATS Version" secondary={natsVersion} />
                </ListItem>
              </List>
            </Grid2>
          </ContentContainer>
        </ColumnWrapper>
      </Grid2>
    </Grid2>
  );
};

const MesheryMetadataFormatter = ({ connection }) => {
  const metadata = connection.metadata || {};
  const uiSchema = createColumnUiSchema({
    metadata,
    numCols: {
      xs: 2,
      md: 4,
    },
  });

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
  const uiSchema = createColumnUiSchema({
    metadata,
    numCols: {
      xs: 3,
      md: 5,
    },
  });

  return (
    <Box backgroundColor={theme.palette.background.card} width="100%" padding={'1rem'}>
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
  const formatterByKind = {
    [KUBERNETES]: () => (
      <KubernetesMetadataFormatter
        meshsyncControllerState={meshsyncControllerState}
        connection={connection}
        metadata={connection.metadata}
      />
    ),
    [MESHERY]: () => <MesheryMetadataFormatter connection={connection} />,
    default: () => (
      <FormatStructuredData
        data={connection.metadata}
        propertyFormatters={DefaultPropertyFormatters}
      />
    ),
  };
  const formatter = formatterByKind[connection.kind] || formatterByKind.default;
  return (
    <Box backgroundColor={theme.palette.background.card} padding={'1rem'}>
      {formatter()}
    </Box>
  );
};

export default FormatConnectionMetadata;
