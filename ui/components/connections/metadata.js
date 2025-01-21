import React from 'react';
import { Grid, List, ListItem, ListItemText, Box, styled } from '@layer5/sistent';
import { makeStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../../lib/store';

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
import { TootltipWrappedConnectionChip } from './ConnectionChip';
import { CONTROLLER_STATES } from '../../utils/Enum';
import { formatToTitleCase } from '../../utils/utils';
import { UsesSistent } from '../SistentWrapper';
import { ColumnWrapper, ContentContainer, OperationButton } from './styles';

const DISABLED = 'DISABLED';
const KUBERNETES = 'kubernetes';
const MESHERY = 'meshery';

const useKubernetesStyles = makeStyles((theme) => ({
  icon: { width: theme.spacing(2.5) },
  operatorIcon: { width: theme.spacing(2.5), filter: theme.palette.secondary.brightness },
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

const StyledListItemText = styled(ListItemText)(({ theme }) => ({
  '& .MuiTypography-root.MuiTypography-body2': {
    color: theme.palette.text.tertiary, // Use the secondary color from the theme
  },
}));

const KubernetesMetadataFormatter = ({ meshsyncControllerState, connection, metadata }) => {
  const classes = useKubernetesStyles();

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
    <UsesSistent>
      <Grid container spacing={1} sx={{ textTransform: 'none' }}>
        <Grid item xs={12} md={6}>
          <ColumnWrapper>
            <Grid container spacing={1}>
              <OperationButton item xs={12} md={5}>
                <List>
                  <ListItem>
                    <TootltipWrappedConnectionChip
                      tooltip={`Server: ${metadata.server}`}
                      title={metadata.name}
                      status={connection.status}
                      iconSrc={'/static/img/kubernetes.svg'}
                      handlePing={() => handleKubernetesClick(connection.id)}
                    />
                  </ListItem>
                </List>
              </OperationButton>
            </Grid>
            <ContentContainer container spacing={1}>
              <Grid item xs={12} md={5}>
                <List>
                  <ListItem>
                    <StyledListItemText primary="Name" secondary={metadata.name} />
                  </ListItem>
                  <ListItem>
                    <StyledListItemText primary="K8s Version" secondary={metadata.version} />
                  </ListItem>
                </List>
              </Grid>
              <Grid item xs={12} md={5}>
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
              </Grid>
              <Grid item xs={12} md={5}>
                <List>
                  <ListItem>
                    <StyledListItemText
                      className={classes.text}
                      primary="Server"
                      secondary={<Link title={metadata.server}>{metadata.server}</Link>}
                    />
                  </ListItem>
                </List>
              </Grid>
            </ContentContainer>
          </ColumnWrapper>
        </Grid>
        <Grid item xs={12} md={6}>
          <ColumnWrapper>
            <Grid container spacing={1}>
              <OperationButton item xs={12} md={4}>
                <List>
                  <ListItem>
                    <TootltipWrappedConnectionChip
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
                  <Grid item xs={12} md={4}>
                    <List>
                      <ListItem>
                        <TootltipWrappedConnectionChip
                          tooltip={meshSyncState !== DISABLED ? `Ping MeshSync` : 'Not Available'}
                          title={'MeshSync'}
                          status={meshSyncState?.toLowerCase()}
                          handlePing={handleMeshSyncClick}
                          iconSrc="/static/img/meshsync.svg"
                          width="9rem"
                        />
                      </ListItem>
                    </List>
                  </Grid>
                  <Grid item xs={12} md={4}>
                    <List>
                      <ListItem>
                        <TootltipWrappedConnectionChip
                          tooltip={natsState === 'Not Active' ? 'Not Available' : `Reconnect NATS`}
                          title={'NATS'}
                          status={natsState?.toLowerCase()}
                          handlePing={() => handleNATSClick()}
                          iconSrc="/static/img/nats-icon-color.svg"
                          width="9rem"
                        />
                      </ListItem>
                    </List>
                  </Grid>
                </>
              )}
            </Grid>

            <ContentContainer container spacing={1}>
              <Grid item xs={12} md={5}>
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
              </Grid>
              <Grid item xs={12} md={5}>
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
              </Grid>
              <Grid item xs={12} md={5}>
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
              </Grid>
            </ContentContainer>
          </ColumnWrapper>
        </Grid>
      </Grid>
    </UsesSistent>
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

const FormatConnectionMetadata = (props) => {
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
    <Box
      sx={{
        padding: '1rem',
      }}
    >
      {formatter()}
    </Box>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = () => ({});

export default connect(mapStateToProps, mapDispatchToProps)(FormatConnectionMetadata);
