import React from 'react';
import { Grid, List, ListItem, ListItemText, Box } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import {
  FormatId,
  FormatStructuredData,
  FormattedDate,
  KeyValue,
  Link,
  createColumnUiSchema,
} from '../DataFormatter';
import useKubernetesHook, {
  useGetOperatorInfoQuery,
  useMesheryOperator,
  useMeshsSyncController,
  useNatsController,
} from '../hooks/useKubernetesHook';
import { ConnectionChip } from './ConnectionChip';

const DISABLED = 'DISABLED';
const KUBERNETES = 'kubernetes';
const MESHERY = 'meshery';

const useKubernetesStyles = makeStyles((theme) => ({
  root: {
    textTransform: 'none',
  },
  operationButton: {
    [theme?.breakpoints?.down(1180)]: {
      marginRight: '25px',
    },
  },
  icon: { width: theme.spacing(2.5) },
  operatorIcon: { width: theme.spacing(2.5), filter: theme.palette.secondary.brightness },
  column: {
    margin: theme.spacing(2),
    padding: theme.spacing(2),
    background: `${theme.palette.secondary.default}10`,
  },
  heading: { textAlign: 'center' },
  configBoxContainer: {
    [theme?.breakpoints?.down(1050)]: {
      flexGrow: 0,
      maxWidth: '100%',
      flexBasis: '100%',
    },
    [theme?.breakpoints?.down(1050)]: {
      flexDirection: 'column',
    },
  },
  clusterConfiguratorWrapper: { padding: theme.spacing(5), display: 'flex' },
  contentContainer: {
    [theme?.breakpoints?.down(1050)]: {
      flexDirection: 'column',
    },
    flexWrap: 'noWrap',
  },
  fileInputStyle: { display: 'none' },
  topToolbar: {
    margin: '1rem 0',
    paddingLeft: '1rem',
    maxWidth: '90%',
  },
  button: {
    padding: theme.spacing(1),
    borderRadius: 5,
  },
  grey: {
    background: 'WhiteSmoke',
    padding: theme.spacing(2),
    borderRadius: 'inherit',
  },
  subtitle: {
    minWidth: 400,
    overflowWrap: 'anywhere',
    textAlign: 'left',
    padding: '5px',
  },
  text: {
    width: '80%',
    wordWrap: 'break-word',
  },
  table: {
    marginTop: theme.spacing(1.5),
  },
  uploadCluster: {
    overflow: 'hidden',
  },
  OperatorSwitch: {
    pointerEvents: 'auto',
  },
}));

const customIdFormatter = (title, id) => <KeyValue Key={title} Value={<FormatId id={id} />} />;
const customDateFormatter = (title, date) => (
  <KeyValue Key={title} Value={<FormattedDate date={date} />} />
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

const KubernetesMetadataFormatter = ({ connection, metadata }) => {
  const classes = useKubernetesStyles();
  const contextID = metadata.id;

  const pingKubernetes = useKubernetesHook();
  const { ping: pingMesheryOperator } = useMesheryOperator();
  const { ping: pingMeshSync } = useMeshsSyncController();
  const { ping: pingNats } = useNatsController();

  const handleKubernetesClick = () => {
    pingKubernetes(metadata.name, metadata.server, connection.id);
  };

  const handleNATSClick = () => {
    pingNats({ connectionID: connection.id });
  };

  const handleOperatorClick = () => {
    pingMesheryOperator(contextID);
  };

  const handleMeshSyncClick = () => {
    pingMeshSync({ connectionID: connection.id });
  };

  const { operatorInfo: operatorStatus } = useGetOperatorInfoQuery({ contextID });
  const {
    isReachable: operatorState,
    meshSyncStatus: meshSyncState,
    natsStatus: natsState,
  } = operatorStatus;

  const operatorVersion = operatorStatus.operatorVersion || 'Not Available';
  const meshSyncVersion = operatorStatus.meshSyncVersion || 'Not Available';
  const NATSVersion = operatorStatus.NATSVersion || 'Not Available';

  return (
    <Grid container spacing={1} className={classes.root}>
      <Grid item xs={12} md={6}>
        <div className={classes.column}>
          <Grid container spacing={1}>
            <Grid item xs={12} md={5} className={classes.operationButton}>
              <List>
                <ListItem>
                  <ConnectionChip
                    tooltip={`Server: ${metadata.server}`}
                    title={metadata.name}
                    status={connection.status}
                    iconSrc={'/static/img/kubernetes.svg'}
                    handlePing={() => handleKubernetesClick(connection.id)}
                  />
                </ListItem>
              </List>
            </Grid>
          </Grid>
          <Grid container spacing={1} className={classes.contentContainer}>
            <Grid item xs={12} md={5}>
              <List>
                <ListItem>
                  <ListItemText primary="Name" secondary={metadata.name} />
                </ListItem>
                <ListItem>
                  <ListItemText primary="K8s Version" secondary={metadata.version} />
                </ListItem>
              </List>
            </Grid>
            <Grid item xs={12} md={5}>
              <List>
                <ListItem>
                  <ListItemText
                    primary="Created At"
                    secondary={<FormattedDate date={connection.created_at} />}
                  />
                </ListItem>
                <ListItem>
                  <ListItemText
                    primary="Updated At"
                    secondary={<FormattedDate date={connection.updated_at} />}
                  />
                </ListItem>
              </List>
            </Grid>
            <Grid item xs={12} md={5}>
              <List>
                <ListItem>
                  <ListItemText
                    className={classes.text}
                    primary="Server"
                    secondary={<Link title={metadata.server}>{metadata.server}</Link>}
                  />
                </ListItem>
              </List>
            </Grid>
          </Grid>
        </div>
      </Grid>
      <Grid item xs={12} md={6}>
        <div className={classes.column}>
          <Grid container spacing={1}>
            <Grid item xs={12} md={4} className={classes.operationButton}>
              <List>
                <ListItem>
                  <ConnectionChip
                    tooltip={operatorState ? `Version: ${operatorVersion}` : 'Not Available'}
                    title={'Operator'}
                    disabled={!operatorState}
                    status={operatorState}
                    handlePing={() => handleOperatorClick(connection.id)}
                    iconSrc="/static/img/meshery-operator.svg"
                  />
                </ListItem>
              </List>
            </Grid>

            {(meshSyncState || natsState) && (
              <>
                <Grid item xs={12} md={4}>
                  <List>
                    <ListItem>
                      <ConnectionChip
                        tooltip={meshSyncState !== DISABLED ? `Ping MeshSync` : 'Not Available'}
                        title={'MeshSync'}
                        status={meshSyncState}
                        handlePing={handleMeshSyncClick}
                        iconSrc="/static/img/meshsync.svg"
                      />
                    </ListItem>
                  </List>
                </Grid>
                <Grid item xs={12} md={4}>
                  <List>
                    <ListItem>
                      <ConnectionChip
                        tooltip={natsState === 'Not Active' ? 'Not Available' : `Reconnect NATS`}
                        title={'NATS'}
                        status={natsState}
                        handlePing={() => handleNATSClick()}
                        iconSrc="/static/img/nats-icon-color.svg"
                      />
                    </ListItem>
                  </List>
                </Grid>
              </>
            )}
          </Grid>

          <Grid container spacing={1} className={classes.contentContainer}>
            <Grid item xs={12} md={5}>
              <List>
                <ListItem>
                  <ListItemText
                    primary="Operator State"
                    secondary={operatorState ? 'Active' : 'Undeployed'}
                  />
                </ListItem>
                <ListItem>
                  <ListItemText primary="Operator Version" secondary={operatorVersion} />
                </ListItem>
              </List>
            </Grid>
            <Grid item xs={12} md={5}>
              <List>
                <ListItem>
                  <ListItemText
                    primary="MeshSync State"
                    secondary={meshSyncState || 'Undeployed'}
                  />
                </ListItem>
                <ListItem>
                  <ListItemText primary="MeshSync Version" secondary={meshSyncVersion} />
                </ListItem>
              </List>
            </Grid>
            <Grid item xs={12} md={5}>
              <List>
                <ListItem>
                  <ListItemText primary="NATS State" secondary={natsState || 'Not Connected'} />
                </ListItem>
                <ListItem>
                  <ListItemText primary="NATS Version" secondary={NATSVersion} />
                </ListItem>
              </List>
            </Grid>
          </Grid>
        </div>
      </Grid>
    </Grid>
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
  const uiSchema = createColumnUiSchema({
    metadata,
    numCols: {
      xs: 3,
      md: 5,
    },
  });

  return (
    <FormatStructuredData
      data={metadata}
      uiSchema={uiSchema}
      propertyFormatters={DefaultPropertyFormatters}
    />
  );
};

export const FormatConnectionMetadata = ({ connection }) => {
  const formatterByKind = {
    [KUBERNETES]: () => (
      <KubernetesMetadataFormatter connection={connection} metadata={connection.metadata} />
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
    <Box
      sx={{
        padding: '1rem',
      }}
    >
      {formatter()}
    </Box>
  );
};
