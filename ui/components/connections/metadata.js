import React from 'react';
import { Tooltip, Grid, Chip, List, ListItem, ListItemText } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { FormatStructuredData, Link, formatDate } from '../DataFormatter';
import useKubernetesHook, {
  useGetOperatorInfoQuery,
  useMesheryOperator,
  useMeshsSyncController,
  useNatsController,
} from '../hooks/useKubernetesHook';

const DISABLED = 'DISABLED';
const KUBERNETES = 'kubernetes';

const useKubernetesStyles = makeStyles((theme) => ({
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
    <Grid container spacing={1}>
      <Grid item xs={12} md={6}>
        <div className={classes.column}>
          <Grid container spacing={1}>
            <Grid item xs={12} md={5} className={classes.operationButton}>
              <List>
                <ListItem>
                  <Tooltip title={`Server: ${metadata.server}`}>
                    <Chip
                      label={metadata.name}
                      icon={<img src="/static/img/kubernetes.svg" className={classes.icon} />}
                      variant="outlined"
                      data-cy="chipContextName"
                      onClick={() => handleKubernetesClick(connection.id)}
                    />
                  </Tooltip>
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
                    secondary={formatDate(connection.created_at)}
                  />
                </ListItem>
                <ListItem>
                  <ListItemText
                    primary="Updated At"
                    secondary={formatDate(connection.updated_at)}
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
                  <Tooltip
                    title={operatorState ? `Version: ${operatorVersion}` : 'Not Available'}
                    aria-label="meshSync"
                  >
                    <Chip
                      // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
                      label={'Operator'}
                      style={!operatorState ? { opacity: 0.5 } : {}}
                      disabled={!operatorState}
                      onClick={() => handleOperatorClick(connection.id)}
                      icon={
                        <img
                          src="/static/img/meshery-operator.svg"
                          className={classes.operatorIcon}
                        />
                      }
                      variant="outlined"
                      data-cy="chipOperator"
                    />
                  </Tooltip>
                </ListItem>
              </List>
            </Grid>

            {(meshSyncState || natsState) && (
              <>
                <Grid item xs={12} md={4}>
                  <List>
                    <ListItem>
                      <Tooltip
                        title={meshSyncState !== DISABLED ? `Ping MeshSync` : 'Not Available'}
                        aria-label="meshSync"
                      >
                        <Chip
                          label={'MeshSync'}
                          style={meshSyncState === DISABLED ? { opacity: 0.5 } : {}}
                          onClick={() => handleMeshSyncClick()}
                          icon={<img src="/static/img/meshsync.svg" className={classes.icon} />}
                          variant="outlined"
                          data-cy="chipMeshSync"
                        />
                      </Tooltip>
                    </ListItem>
                  </List>
                </Grid>
                <Grid item xs={12} md={4}>
                  <List>
                    <ListItem>
                      <Tooltip
                        title={natsState === 'Not Active' ? 'Not Available' : `Reconnect NATS`}
                        aria-label="nats"
                      >
                        <Chip
                          label={'NATS'}
                          onClick={() => handleNATSClick()}
                          style={natsState === 'Not Active' ? { opacity: 0.5 } : {}}
                          icon={
                            <img src="/static/img/nats-icon-color.svg" className={classes.icon} />
                          }
                          variant="outlined"
                          data-cy="chipNATS"
                        />
                      </Tooltip>
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

export const FormatConnectionMetadata = ({ connection }) => {
  const formatterByKind = {
    [KUBERNETES]: () => (
      <KubernetesMetadataFormatter connection={connection} metadata={connection.metadata} />
    ),
    default: () => <FormatStructuredData data={connection.metadata} />,
  };
  const formatter = formatterByKind[connection.kind] || formatterByKind.default;
  return formatter();
};
